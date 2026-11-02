package attestation

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

const (
	nb       = "2026-01-01T00:00:00Z"
	na       = "2027-01-01T00:00:00Z"
	issued   = "2026-06-01T12:00:00Z"
	evalTime = "2026-06-02T00:00:00Z"
)

var (
	subject = "sha256:" + strings.Repeat("ab", 32)
	otherID = "sha256:" + strings.Repeat("cd", 32)
)

func eval(t *testing.T) time.Time {
	t.Helper()
	ts, err := time.Parse("2006-01-02T15:04:05Z07:00", evalTime)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}

type keySpec struct {
	pub       ed25519.PublicKey
	roles     []string
	notBefore string
	notAfter  string
	revoked   bool
}

func newKey(t *testing.T) (ed25519.PrivateKey, ed25519.PublicKey) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return priv, pub
}

func bundleBytes(t *testing.T, specs []keySpec) []byte {
	t.Helper()
	entries := make([]map[string]any, 0, len(specs))
	for _, s := range specs {
		entries = append(entries, map[string]any{
			"key_id":     artifact.HashBytes(s.pub),
			"public_key": base64.StdEncoding.EncodeToString(s.pub),
			"roles":      s.roles,
			"not_before": s.notBefore,
			"not_after":  s.notAfter,
			"revoked":    s.revoked,
		})
	}
	raw, err := artifact.CanonicalValue(map[string]any{
		"contract": TrustContract,
		"keys":     entries,
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func parseBundle(t *testing.T, raw []byte) *TrustBundle {
	t.Helper()
	b, err := ParseTrustBundle(raw)
	if err != nil {
		t.Fatalf("ParseTrustBundle: %v", err)
	}
	return b
}

func attestationRaw(t *testing.T, a *Attestation) []byte {
	t.Helper()
	raw, err := artifact.CanonicalValue(map[string]string{
		"contract":  a.Contract,
		"subject":   a.Subject,
		"role":      a.Role,
		"key_id":    a.KeyID,
		"issued_at": a.IssuedAt,
		"signature": a.Signature,
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func codeOf(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("error is not *Error: %v", err)
	}
	return e.Code
}

func TestSignVerifyRoundTrip(t *testing.T) {
	priv, pub := newKey(t)
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na,
	}}))
	a, err := Sign(subject, "claim-validator", issued, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	// Round-trip through the raw envelope parser.
	parsed, err := ParseAttestation(attestationRaw(t, a))
	if err != nil {
		t.Fatalf("ParseAttestation: %v", err)
	}
	if err := VerifyAttestation(parsed, bundle, eval(t)); err != nil {
		t.Fatalf("VerifyAttestation: %v", err)
	}
	// Verification uses only public material; the private key is not referenced.
	if err := VerifyAttestation(a, bundle, eval(t)); err != nil {
		t.Fatalf("VerifyAttestation (direct): %v", err)
	}
}

func TestTamperingRejected(t *testing.T) {
	priv, pub := newKey(t)
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator", "coverage-validator"},
		notBefore: nb, notAfter: na,
	}}))
	base, err := Sign(subject, "claim-validator", issued, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	cases := map[string]*Attestation{
		"subject":   {Subject: otherID, Role: base.Role, KeyID: base.KeyID, IssuedAt: base.IssuedAt, Signature: base.Signature, Contract: base.Contract},
		"role":      {Subject: base.Subject, Role: "coverage-validator", KeyID: base.KeyID, IssuedAt: base.IssuedAt, Signature: base.Signature, Contract: base.Contract},
		"issued_at": {Subject: base.Subject, Role: base.Role, KeyID: base.KeyID, IssuedAt: "2026-06-01T12:00:01Z", Signature: base.Signature, Contract: base.Contract},
		"key_id":    {Subject: base.Subject, Role: base.Role, KeyID: "sha256:" + strings.Repeat("ef", 32), IssuedAt: base.IssuedAt, Signature: base.Signature, Contract: base.Contract},
	}
	for name, a := range cases {
		if err := VerifyAttestation(a, bundle, eval(t)); err == nil {
			t.Errorf("tampered %s accepted", name)
		}
	}
	// Flip a signature character (still valid base64, wrong bytes).
	flip := *base
	flip.Signature = flip.Signature[:len(flip.Signature)-2] + "AA"
	if err := VerifyAttestation(&flip, bundle, eval(t)); err == nil {
		t.Error("flipped signature accepted")
	}
}

func TestRoleMismatch(t *testing.T) {
	priv, pub := newKey(t)
	// Key is only authorized for claim-validator.
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na,
	}}))
	// Sign under coverage-validator (a known role) but the key lacks it.
	a, err := Sign(subject, "coverage-validator", issued, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if err := VerifyAttestation(a, bundle, eval(t)); err == nil {
		t.Fatal("role mismatch accepted")
	}
}

func TestUnknownKey(t *testing.T) {
	privA, _ := newKey(t)
	_, pubB := newKey(t)
	// Bundle only knows key B, but the attestation is from key A.
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pubB, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na,
	}}))
	a, err := Sign(subject, "claim-validator", issued, privA)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if codeOf(t, VerifyAttestation(a, bundle, eval(t))) != CodeUnknownKey {
		t.Fatal("expected unknown-key")
	}
}

func TestRevokedKey(t *testing.T) {
	priv, pub := newKey(t)
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na, revoked: true,
	}}))
	a, err := Sign(subject, "claim-validator", issued, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if codeOf(t, VerifyAttestation(a, bundle, eval(t))) != CodeRevoked {
		t.Fatal("expected revoked")
	}
}

func TestRotationOverlap(t *testing.T) {
	privA, pubA := newKey(t)
	privB, pubB := newKey(t)
	// Two overlapping authorized keys (rotation in progress).
	bundle := parseBundle(t, bundleBytes(t, []keySpec{
		{pub: pubA, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na},
		{pub: pubB, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na},
	}))
	for name, priv := range map[string]ed25519.PrivateKey{"A": privA, "B": privB} {
		a, err := Sign(subject, "claim-validator", issued, priv)
		if err != nil {
			t.Fatalf("Sign %s: %v", name, err)
		}
		if err := VerifyAttestation(a, bundle, eval(t)); err != nil {
			t.Fatalf("rotation key %s rejected: %v", name, err)
		}
	}
}

func TestMissingCredential(t *testing.T) {
	if _, err := Sign(subject, "claim-validator", issued, nil); err == nil {
		t.Fatal("nil key accepted")
	}
	short := make([]byte, 31)
	if _, err := Sign(subject, "claim-validator", issued, short); err == nil {
		t.Fatal("short key accepted")
	}
}

func TestTimeBounds(t *testing.T) {
	priv, pub := newKey(t)
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na,
	}}))
	// issued_at before the key became valid.
	early, _ := Sign(subject, "claim-validator", "2025-01-01T00:00:00Z", priv)
	if codeOf(t, VerifyAttestation(early, bundle, eval(t))) != CodeNotYetValid {
		t.Fatal("expected not-yet-valid for early issued_at")
	}
	// Evaluation time after the key expired.
	lateEval, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2027-06-01T00:00:00Z")
	a, _ := Sign(subject, "claim-validator", issued, priv)
	if codeOf(t, VerifyAttestation(a, bundle, lateEval)) != CodeExpired {
		t.Fatal("expected expired for late eval time")
	}
}

func TestMalformedSignatureEncoding(t *testing.T) {
	priv, pub := newKey(t)
	bundle := parseBundle(t, bundleBytes(t, []keySpec{{
		pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na,
	}}))
	a, _ := Sign(subject, "claim-validator", issued, priv)
	// Wrong-length base64.
	bad := *a
	bad.Signature = base64.StdEncoding.EncodeToString([]byte("short"))
	if codeOf(t, VerifyAttestation(&bad, bundle, eval(t))) != CodeBadSignature {
		t.Fatal("expected bad-signature")
	}
	// Non-base64.
	bad2 := *a
	bad2.Signature = "!!!!not-base64!!!!"
	if codeOf(t, VerifyAttestation(&bad2, bundle, eval(t))) != CodeBadSignature {
		t.Fatal("expected bad-signature for non-base64")
	}
}

func TestRoleForKind(t *testing.T) {
	if r, err := RoleForKind("claim-report"); err != nil || r != "claim-validator" {
		t.Fatalf("claim-report role = %q, %v", r, err)
	}
	if r, err := RoleForKind("extraction-approval"); err != nil || r != "extraction-gate" {
		t.Fatalf("extraction-approval role = %q, %v", r, err)
	}
	if _, err := RoleForKind("source"); err == nil {
		t.Fatal("non-report kind accepted")
	}
}

func TestTrustBundleValidation(t *testing.T) {
	_, pub := newKey(t)
	// Duplicate keys.
	parseFail(t, bundleBytes(t, []keySpec{
		{pub: pub, roles: []string{"claim-validator"}, notBefore: nb, notAfter: na},
		{pub: pub, roles: []string{"coverage-validator"}, notBefore: nb, notAfter: na},
	}), CodeDupKey)
	// Mismatched key_id.
	raw := map[string]any{
		"contract": TrustContract,
		"keys": []map[string]any{{
			"key_id":     "sha256:" + strings.Repeat("00", 32),
			"public_key": base64.StdEncoding.EncodeToString(pub),
			"roles":      []string{"claim-validator"},
			"not_before": nb, "not_after": na, "revoked": false,
		}},
	}
	expectCode(t, raw, CodeKeyIDMismatch)
	// Unknown role.
	expectCode(t, map[string]any{
		"contract": TrustContract,
		"keys": []map[string]any{{
			"key_id": artifact.HashBytes(pub), "public_key": base64.StdEncoding.EncodeToString(pub),
			"roles": []string{"not-a-role"}, "not_before": nb, "not_after": na, "revoked": false,
		}},
	}, CodeBadRole)
	// Empty roles.
	expectCode(t, map[string]any{
		"contract": TrustContract,
		"keys": []map[string]any{{
			"key_id": artifact.HashBytes(pub), "public_key": base64.StdEncoding.EncodeToString(pub),
			"roles": []string{}, "not_before": nb, "not_after": na, "revoked": false,
		}},
	}, CodeMissingField)
	// not_before after not_after.
	expectCode(t, map[string]any{
		"contract": TrustContract,
		"keys": []map[string]any{{
			"key_id": artifact.HashBytes(pub), "public_key": base64.StdEncoding.EncodeToString(pub),
			"roles": []string{"claim-validator"}, "not_before": na, "not_after": nb, "revoked": false,
		}},
	}, CodeBadTime)
	// Unknown contract.
	expectCode(t, map[string]any{
		"contract": "atlas.trust/v2",
		"keys":     []map[string]any{},
	}, CodeBadContract)
	// Unknown top-level field.
	expectCode(t, map[string]any{
		"contract": TrustContract,
		"keys":     []map[string]any{},
		"extra":    1,
	}, CodeMalformed)
}

func expectCode(t *testing.T, v any, want string) {
	t.Helper()
	raw, err := artifact.CanonicalValue(v)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	if _, err := ParseTrustBundle(raw); err == nil {
		t.Fatalf("expected %s, accepted", want)
	} else if codeOf(t, err) != want {
		t.Fatalf("want %s, got %s", want, codeOf(t, err))
	}
}

func parseFail(t *testing.T, raw []byte, want string) error {
	t.Helper()
	if _, err := ParseTrustBundle(raw); err == nil {
		t.Fatalf("expected %s, accepted", want)
	}
	return nil
}
