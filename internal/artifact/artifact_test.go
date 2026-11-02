package artifact

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// --- fixtures / helpers -------------------------------------------------

func mustCanonic(t *testing.T, v any) []byte {
	t.Helper()
	b, err := encodeCanonical(v)
	if err != nil {
		t.Fatalf("encodeCanonical: %v", err)
	}
	return b
}

func codeOf(t *testing.T, err error) string {
	t.Helper()
	var ae *Error
	if !errors.As(err, &ae) {
		t.Fatalf("error is not *Error: %v", err)
	}
	return ae.Code
}

// --- canonicalization ---------------------------------------------------

func TestCanonicalReferenceVectors(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		code string
	}{
		{"sorted", `{"b":2,"a":1}`, `{"a":1,"b":2}`, ""},
		{"whitespace", `{"a":1, "b":2}`, `{"a":1,"b":2}`, ""},
		{"noHTMLEscape", `{"quote":"<salud>&"}`, `{"quote":"<salud>&"}`, ""},
		{"dupMember", `{"a":1,"a":2}`, "", CodeDupMember},
		{"fractional", `{"score":0.5}`, "", CodeFractional},
		{"exponent", `{"n":1e0}`, "", CodeFractional},
		{"negZero", `{"n":-0}`, "", CodeNegZero},
		{"outOfRange", `{"n":9007199254740992}`, "", CodeOutOfRange},
		{"trailing", `{} {}`, "", CodeTrailing},
		{"twoValues", `[1] 2`, "", CodeTrailing},
		{"bom", "\ufeff{}", "", CodeMalformed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Canonicalize([]byte(tc.in))
			if tc.code != "" {
				if err == nil {
					t.Fatalf("want error %s, got none (out=%s)", tc.code, out)
				}
				if codeOf(t, err) != tc.code {
					t.Fatalf("want code %s, got %v", tc.code, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(out) != tc.want {
				t.Fatalf("canon = %s, want %s", out, tc.want)
			}
		})
	}
}

func TestCanonicalInvalidUTF8AndSurrogates(t *testing.T) {
	// Invalid raw UTF-8 byte in a string.
	if _, err := Canonicalize([]byte(`{"s":"\xff"}`)); err == nil {
		t.Error("want invalid-utf8 rejection")
	}
	// Unpaired high surrogate.
	if _, err := Canonicalize([]byte(`{"s":"\ud800"}`)); err == nil {
		t.Error("want unpaired-surrogate rejection")
	}
	// Invalid surrogate pair (high then high, not a low surrogate).
	if _, err := Canonicalize([]byte(`{"s":"\ud800\ud800"}`)); err == nil {
		t.Error("want invalid-surrogate-pair rejection")
	}
	// A valid pair encodes U+10000.
	out, err := Canonicalize([]byte(`{"s":"\ud800\udc00"}`))
	if err != nil {
		t.Fatalf("valid pair rejected: %v", err)
	}
	if !bytes.Contains(out, []byte("\U00010000")) {
		t.Fatalf("valid pair not encoded as U+10000: %s", out)
	}
}

func TestCanonicalDepthLimit(t *testing.T) {
	// 64 nested arrays is allowed (root is depth 1).
	ok := strings.Repeat("[", 64) + strings.Repeat("]", 64)
	if _, err := Canonicalize([]byte(ok)); err != nil {
		t.Fatalf("depth 64 rejected: %v", err)
	}
	// 65 nested arrays is rejected.
	too := strings.Repeat("[", 65) + strings.Repeat("]", 65)
	if _, err := Canonicalize([]byte(too)); err == nil {
		t.Fatal("depth 65 not rejected")
	}
}

func TestLimits(t *testing.T) {
	// Canonicalize rejects input over the 16 MiB payload bound.
	over := make([]byte, 16<<20+1)
	over[0] = '{'
	if _, err := Canonicalize(over); err == nil {
		t.Fatal("oversized payload accepted")
	}
}

func TestCanonicalIntegerRange(t *testing.T) {
	if _, err := Canonicalize([]byte(`{"n":9007199254740991}`)); err != nil {
		t.Fatalf("max int rejected: %v", err)
	}
	if _, err := Canonicalize([]byte(`{"n":-9007199254740991}`)); err != nil {
		t.Fatalf("min int rejected: %v", err)
	}
	if _, err := Canonicalize([]byte(`{"n":9007199254740992}`)); err == nil {
		t.Fatal("max+1 not rejected")
	}
}

func TestDeterminism(t *testing.T) {
	// Same logical content, different key order and indentation -> same bytes/ID.
	a := `{"kind":"source","contract":"atlas.artifact/v1","inputs":{},"payload":{"size":3,"sha256":"sha256:0000000000000000000000000000000000000000000000000000000000000000","media_type":"text/plain"}}`
	b := "{\n  \"payload\": { \"media_type\": \"text/plain\", \"sha256\": \"sha256:0000000000000000000000000000000000000000000000000000000000000000\", \"size\": 3 },\n  \"kind\": \"source\",\n  \"inputs\": {},\n  \"contract\": \"atlas.artifact/v1\"\n}"
	ca, errA := Canonicalize([]byte(a))
	cb, errB := Canonicalize([]byte(b))
	if errA != nil || errB != nil {
		t.Fatalf("canonicalize failed: %v / %v", errA, errB)
	}
	if !bytes.Equal(ca, cb) {
		t.Fatalf("canonical bytes differ:\n%s\n%s", ca, cb)
	}
	if HashBytes(ca) != HashBytes(cb) {
		t.Fatal("content IDs differ")
	}
}

func TestValidHash(t *testing.T) {
	if !ValidHash("sha256:" + strings.Repeat("a", 64)) {
		t.Fatal("valid hash rejected")
	}
	for _, bad := range []string{
		"sha256:" + strings.Repeat("a", 63),
		"sha256:" + strings.Repeat("a", 65),
		"sha512:" + strings.Repeat("a", 64),
		"sha256:" + strings.Repeat("A", 64), // uppercase is not lowercase hex
		"sha256:" + strings.Repeat("g", 64),
		"",
	} {
		if ValidHash(bad) {
			t.Fatalf("invalid hash accepted: %q", bad)
		}
	}
}
