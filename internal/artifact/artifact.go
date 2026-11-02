// Package artifact provides the strict canonical-JSON serialization and SHA-256
// content-identity helpers reused across the pipeline (TQ-0063). The former
// envelope/manifest/graph machinery was retired; only the bounded strict parser,
// canonical re-encoding, and hash helpers remain. Verification is fail-closed
// and never retrieves files, follows URLs, or invokes a model.
package artifact

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Contract is the canonical-JSON serialization profile implemented here.
const Contract = "atlas.artifact/v1"

// Limits. None may be disabled by values inside untrusted input.
const (
	// maxCanonicalInput bounds a JSON payload fed to Canonicalize.
	maxCanonicalInput = 16 << 20
	// maxDepth is the maximum JSON nesting depth (root is depth 1).
	maxDepth = 64
)

// v1 integer bounds: [-(2^53-1), 2^53-1].
const (
	maxInt = 9007199254740991
	minInt = -9007199254740991
)

// Error codes. They are stable and carry no private data.
const (
	CodeMalformed   = "malformed-json"
	CodeTrailing    = "trailing-value"
	CodeDupMember   = "duplicate-member"
	CodeFractional  = "fractional-number"
	CodeNegZero     = "negative-zero"
	CodeOutOfRange  = "integer-out-of-range"
	CodeTooDeep     = "too-deep"
	CodeTooLarge    = "too-large"
	CodeInvalidUTF8 = "invalid-utf8"
)

// Error is a stable, fail-closed artifact error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "artifact: " + e.Code
	if e.Subject != "" {
		s += " (" + e.Subject + ")"
	}
	if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

// Unwrap exposes the underlying error when present.
func (e *Error) Unwrap() error { return e.Err }

func errf(code, subject string) error { return &Error{Code: code, Subject: subject} }

// HashBytes returns the atlas H(bytes) identity: "sha256:" plus 64 lowercase
// hex digits of the SHA-256 digest.
func HashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ValidHash reports whether h is a full H identity: "sha256:" followed by 64
// lowercase hexadecimal digits.
func ValidHash(h string) bool { return validHash(h) }

func validHash(h string) bool {
	if len(h) != 7+64 || h[:7] != "sha256:" {
		return false
	}
	for i := 7; i < len(h); i++ {
		c := h[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// CanonicalValue returns the canonical bytes for a Go value restricted to the
// strict v1 profile (objects, arrays, strings, booleans, null, and integer
// json.Number). It is a convenience for building canonical statements without
// hand-serializing JSON.
func CanonicalValue(v any) ([]byte, error) { return encodeCanonical(v) }

// Canonicalize validates that data is exactly one strict JSON value under the
// atlas.artifact/v1 serialization rules and returns its canonical bytes.
//
// The rules reject a BOM, trailing data, duplicate object members, fractional
// or exponent numbers, negative zero, integers outside the v1 range, nesting
// deeper than maxDepth, invalid UTF-8, and unpaired surrogates. Keys are sorted
// by UTF-8 byte order; array order and string code points are preserved. The
// output uses Go encoding/json escaping with HTML escaping disabled and has no
// trailing newline.
func Canonicalize(data []byte) ([]byte, error) {
	if len(data) > maxCanonicalInput {
		return nil, errf(CodeTooLarge, "input")
	}
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return nil, errf(CodeMalformed, "bom")
	}
	p := &parser{b: data}
	p.skipWS()
	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	p.skipWS()
	if p.i != len(data) {
		return nil, errf(CodeTrailing, fmt.Sprintf("%d", p.i))
	}
	return encodeCanonical(v)
}

// parser is a strict, bounded, integer-preserving JSON scanner. It validates
// while building a Go tree so that canonicalization is a pure re-encoding.
type parser struct {
	b     []byte
	i     int
	depth int
}

func (p *parser) skipWS() {
	for p.i < len(p.b) {
		switch p.b[p.i] {
		case ' ', '\t', '\n', '\r':
			p.i++
		default:
			return
		}
	}
}

func (p *parser) parseValue() (any, error) {
	p.skipWS()
	if p.i >= len(p.b) {
		return nil, errf(CodeMalformed, "eof")
	}
	switch c := p.b[p.i]; c {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't':
		return p.parseLit("true", true)
	case 'f':
		return p.parseLit("false", false)
	case 'n':
		return p.parseLit("null", nil)
	default:
		return p.parseNumber()
	}
}

func (p *parser) parseLit(word string, val any) (any, error) {
	if p.i+len(word) <= len(p.b) && string(p.b[p.i:p.i+len(word)]) == word {
		p.i += len(word)
		return val, nil
	}
	return nil, errf(CodeMalformed, string(p.b[p.i]))
}

func (p *parser) parseObject() (any, error) {
	if p.depth+1 > maxDepth {
		return nil, errf(CodeTooDeep, fmt.Sprintf("%d", p.depth+1))
	}
	p.i++ // consume '{'
	p.depth++
	defer func() { p.depth-- }()
	obj := make(map[string]any)
	p.skipWS()
	if p.i < len(p.b) && p.b[p.i] == '}' {
		p.i++
		return obj, nil
	}
	for {
		p.skipWS()
		if p.i >= len(p.b) || p.b[p.i] != '"' {
			return nil, errf(CodeMalformed, "object-key")
		}
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		if _, dup := obj[key]; dup {
			return nil, errf(CodeDupMember, key)
		}
		p.skipWS()
		if p.i >= len(p.b) || p.b[p.i] != ':' {
			return nil, errf(CodeMalformed, "object-colon")
		}
		p.i++ // consume ':'
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = val
		p.skipWS()
		if p.i >= len(p.b) {
			return nil, errf(CodeMalformed, "object-end")
		}
		switch p.b[p.i] {
		case ',':
			p.i++
		case '}':
			p.i++
			return obj, nil
		default:
			return nil, errf(CodeMalformed, "object-sep")
		}
	}
}

func (p *parser) parseArray() (any, error) {
	if p.depth+1 > maxDepth {
		return nil, errf(CodeTooDeep, fmt.Sprintf("%d", p.depth+1))
	}
	p.i++ // consume '['
	p.depth++
	defer func() { p.depth-- }()
	arr := make([]any, 0, 4)
	p.skipWS()
	if p.i < len(p.b) && p.b[p.i] == ']' {
		p.i++
		return arr, nil
	}
	for {
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
		p.skipWS()
		if p.i >= len(p.b) {
			return nil, errf(CodeMalformed, "array-end")
		}
		switch p.b[p.i] {
		case ',':
			p.i++
		case ']':
			p.i++
			return arr, nil
		default:
			return nil, errf(CodeMalformed, "array-sep")
		}
	}
}

// parseString reads a JSON string into a Go string, rejecting unescaped control
// characters, invalid UTF-8, and unpaired or invalid surrogate escapes.
func (p *parser) parseString() (string, error) {
	p.i++ // consume opening quote
	var sb strings.Builder
	for {
		if p.i >= len(p.b) {
			return "", errf(CodeMalformed, "string-eof")
		}
		c := p.b[p.i]
		switch {
		case c == '"':
			p.i++
			return sb.String(), nil
		case c == '\\':
			p.i++
			if p.i >= len(p.b) {
				return "", errf(CodeMalformed, "string-escape-eof")
			}
			e := p.b[p.i]
			switch e {
			case '"', '\\', '/':
				sb.WriteByte(e)
				p.i++
			case 'b':
				sb.WriteByte('\b')
				p.i++
			case 'f':
				sb.WriteByte('\f')
				p.i++
			case 'n':
				sb.WriteByte('\n')
				p.i++
			case 'r':
				sb.WriteByte('\r')
				p.i++
			case 't':
				sb.WriteByte('\t')
				p.i++
			case 'u':
				r, err := p.parseUnicode()
				if err != nil {
					return "", err
				}
				if utf16.IsSurrogate(r) {
					// p.i is at the second escape's backslash; advance to its 'u'
					// so parseUnicode can consume the \uXXXX for the low surrogate.
					if p.i+1 < len(p.b) && p.b[p.i] == '\\' && p.b[p.i+1] == 'u' {
						p.i++
						r2, err := p.parseUnicode()
						if err != nil {
							return "", err
						}
						combined := utf16.DecodeRune(r, r2)
						if combined == utf8.RuneError {
							return "", errf(CodeInvalidUTF8, "surrogate-pair")
						}
						sb.WriteRune(combined)
					} else {
						return "", errf(CodeInvalidUTF8, "unpaired-surrogate")
					}
				} else {
					sb.WriteRune(r)
				}
			default:
				return "", errf(CodeMalformed, "string-escape")
			}
		case c < 0x20:
			return "", errf(CodeMalformed, "string-control")
		default:
			r, size := utf8.DecodeRune(p.b[p.i:])
			if r == utf8.RuneError && size == 1 {
				return "", errf(CodeInvalidUTF8, "")
			}
			sb.WriteRune(r)
			p.i += size
		}
	}
}

func (p *parser) parseUnicode() (rune, error) {
	p.i++ // consume 'u'
	if p.i+4 > len(p.b) {
		return 0, errf(CodeMalformed, "unicode-eof")
	}
	var r rune
	for j := 0; j < 4; j++ {
		h := p.b[p.i+j]
		var d rune
		switch {
		case h >= '0' && h <= '9':
			d = rune(h - '0')
		case h >= 'a' && h <= 'f':
			d = rune(h-'a') + 10
		case h >= 'A' && h <= 'F':
			d = rune(h-'A') + 10
		default:
			return 0, errf(CodeMalformed, "unicode-hex")
		}
		r = r*16 + d
	}
	p.i += 4
	return r, nil
}

// parseNumber reads a strict v1 integer, rejecting fractional/exponent tokens
// and negative zero, and returning an integer-preserving json.Number.
func (p *parser) parseNumber() (any, error) {
	start := p.i
	neg := false
	if p.i < len(p.b) && p.b[p.i] == '-' {
		neg = true
		p.i++
	}
	if p.i >= len(p.b) || p.b[p.i] < '0' || p.b[p.i] > '9' {
		return nil, errf(CodeMalformed, "number-int")
	}
	if p.b[p.i] == '0' {
		p.i++
		if p.i < len(p.b) && p.b[p.i] >= '0' && p.b[p.i] <= '9' {
			return nil, errf(CodeMalformed, "leading-zero")
		}
	} else {
		for p.i < len(p.b) && p.b[p.i] >= '0' && p.b[p.i] <= '9' {
			p.i++
		}
	}
	if p.i < len(p.b) {
		switch p.b[p.i] {
		case '.', 'e', 'E':
			return nil, errf(CodeFractional, "")
		}
	}
	tok := string(p.b[start:p.i])
	if neg {
		allZero := true
		for k := 1; k < len(tok); k++ {
			if tok[k] != '0' {
				allZero = false
				break
			}
		}
		if allZero {
			return nil, errf(CodeNegZero, "")
		}
	}
	n, err := strconv.ParseInt(tok, 10, 64)
	if err != nil || n > maxInt || n < minInt {
		return nil, errf(CodeOutOfRange, tok)
	}
	return json.Number(tok), nil
}

// encodeCanonical encodes a validated Go tree to canonical bytes: objects with
// keys sorted by UTF-8 byte order (Go's encoding/json map ordering), arrays and
// strings preserved, Go HTML-escaping disabled, and no trailing newline.
func encodeCanonical(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}
