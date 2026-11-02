// Package attestation implements the atlas.attestation/v1 and atlas.trust/v1
// contracts established by docs/design/artifact-contract.md (TQ-0008) for
// TQ-0020.
//
// It provides Ed25519 signing of detached report attestations and offline,
// fail-closed verification against an explicit owner-reviewed trust bundle.
// Verification uses public material only and never touches the network.
// Signing accepts key bytes only; it never reads key paths from a report and
// never includes key material in errors.
package attestation

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Contract versions.
const (
	AttestationContract = "atlas.attestation/v1"
	TrustContract       = "atlas.trust/v1"
)

// Limits.
const (
	maxEnvelope = 16 << 10 // 16 KiB attestation envelope
	maxBundle   = 1 << 20  // 1 MiB trust bundle
	maxKeys     = 256      // keys per trust bundle
)

// Error codes. Stable and free of key material.
const (
	CodeMalformed     = "malformed"
	CodeTrailing      = "trailing-value"
	CodeTooLarge      = "too-large"
	CodeBadContract   = "unsupported-version"
	CodeMissingField  = "missing-field"
	CodeBadSubject    = "bad-subject"
	CodeBadRole       = "bad-role"
	CodeBadKeyID      = "bad-key-id"
	CodeBadTime       = "bad-time"
	CodeBadSignature  = "bad-signature"
	CodeBadPublicKey  = "bad-public-key"
	CodeKeyIDMismatch = "key-id-mismatch"
	CodeDupKey        = "duplicate-key"
	CodeDupRole       = "duplicate-role"
	CodeTooManyKeys   = "too-many-keys"
	CodeUnknownKey    = "unknown-key"
	CodeRoleMismatch  = "role-mismatch"
	CodeNotYetValid   = "not-yet-valid"
	CodeExpired       = "expired"
	CodeRevoked       = "revoked"
	CodeSignature     = "signature"
)

// Error is a stable attestation error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "attestation: " + e.Code
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

// kindRole is the required report-kind to signer-role mapping.
var kindRole = map[string]string{
	"provenance":                "provenance",
	"extraction-audit":          "extraction-audit",
	"extraction-approval":       "extraction-gate",
	"claim-report":              "claim-validator",
	"coverage-report":           "coverage-validator",
	"content-approval":          "content-gate",
	"site-deterministic-report": "site-deterministic",
	"site-visual-report":        "site-visual",
	"site-approval":             "site-gate",
	"publication":               "publication-gate",
}

// Roles returns the set of known signer roles.
func Roles() map[string]bool {
	out := map[string]bool{}
	for _, r := range kindRole {
		out[r] = true
	}
	return out
}

// RoleForKind returns the required signer role for a report kind, or an error
// for a kind that is not an attested report kind.
func RoleForKind(kind string) (string, error) {
	r, ok := kindRole[kind]
	if !ok {
		return "", errf(CodeBadRole, kind)
	}
	return r, nil
}

// Attestation is a detached Ed25519 attestation for a report manifest.
type Attestation struct {
	Contract  string
	Subject   string
	Role      string
	KeyID     string
	IssuedAt  string
	Signature string
}

// TrustKey is one trusted signer entry in an atlas.trust/v1 bundle.
type TrustKey struct {
	KeyID     string
	PublicKey []byte
	Roles     []string
	NotBefore time.Time
	NotAfter  time.Time
	Revoked   bool
}

// TrustBundle is the owner-reviewed trusted signer set.
type TrustBundle struct {
	Contract string
	Keys     []TrustKey
}

// Sign creates a valid attestation for subject (a manifest ID) under role at
// issuedAt, signed with the Ed25519 private key. issuedAt must be whole-second
// UTC RFC3339. On any validation failure it returns an error and no attestation,
// so a signing failure cannot leave an apparent approval. Key bytes are never
// echoed in errors.
func Sign(subject, role, issuedAt string, priv ed25519.PrivateKey) (*Attestation, error) {
	if !artifact.ValidHash(subject) {
		return nil, errf(CodeBadSubject, "subject")
	}
	if !Roles()[role] {
		return nil, errf(CodeBadRole, role)
	}
	if _, err := parseWholeSecondUTC(issuedAt); err != nil {
		return nil, errf(CodeBadTime, "issued_at")
	}
	if len(priv) != ed25519.PrivateKeySize {
		return nil, errf(CodeMalformed, "key-size")
	}
	pub := priv.Public().(ed25519.PublicKey)
	keyID := artifact.HashBytes(pub)
	stmt := signedStatement(AttestationContract, subject, role, keyID, issuedAt)
	sig := ed25519.Sign(priv, stmt)
	return &Attestation{
		Contract:  AttestationContract,
		Subject:   subject,
		Role:      role,
		KeyID:     keyID,
		IssuedAt:  issuedAt,
		Signature: base64.StdEncoding.EncodeToString(sig),
	}, nil
}

// signedStatement builds the exact bytes that are signed: the contract name, a
// newline, then the canonical encoding of the other five fields (excluding
// signature).
func signedStatement(contract, subject, role, keyID, issuedAt string) []byte {
	canon, err := artifact.CanonicalValue(map[string]string{
		"contract":  contract,
		"subject":   subject,
		"role":      role,
		"key_id":    keyID,
		"issued_at": issuedAt,
	})
	if err != nil {
		// All inputs are validated plain strings; canonicalization cannot fail.
		panic(err)
	}
	return append(append([]byte(contract), '\n'), canon...)
}

// ParseAttestation validates the strict envelope of an attestation from raw
// bytes: exactly the six fields, the v1 contract, a valid subject/role/key-id,
// a whole-second UTC issued_at, and a standard padded base64 signature of 64
// bytes. It does not verify the signature.
func ParseAttestation(raw []byte) (*Attestation, error) {
	if len(raw) > maxEnvelope {
		return nil, errf(CodeTooLarge, "envelope")
	}
	canon, err := artifact.Canonicalize(raw)
	if err != nil {
		var ae *artifact.Error
		if errorAsAE(err, &ae) {
			return nil, &Error{Code: ae.Code, Subject: ae.Subject}
		}
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	var wire struct {
		Contract  *string `json:"contract"`
		Subject   *string `json:"subject"`
		Role      *string `json:"role"`
		KeyID     *string `json:"key_id"`
		IssuedAt  *string `json:"issued_at"`
		Signature *string `json:"signature"`
	}
	dec := json.NewDecoder(bytes.NewReader(canon))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&wire); err != nil {
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	for name, p := range map[string]*string{
		"contract": wire.Contract, "subject": wire.Subject, "role": wire.Role,
		"key_id": wire.KeyID, "issued_at": wire.IssuedAt, "signature": wire.Signature,
	} {
		if p == nil {
			return nil, errf(CodeMissingField, name)
		}
	}
	if *wire.Contract != AttestationContract {
		return nil, errf(CodeBadContract, "contract")
	}
	if !artifact.ValidHash(*wire.Subject) {
		return nil, errf(CodeBadSubject, "subject")
	}
	if !Roles()[*wire.Role] {
		return nil, errf(CodeBadRole, "role")
	}
	if !artifact.ValidHash(*wire.KeyID) {
		return nil, errf(CodeBadKeyID, "key_id")
	}
	if _, err := parseWholeSecondUTC(*wire.IssuedAt); err != nil {
		return nil, errf(CodeBadTime, "issued_at")
	}
	sig, err := decodeSignature(*wire.Signature)
	if err != nil {
		return nil, errf(CodeBadSignature, "signature")
	}
	_ = sig
	return &Attestation{
		Contract:  *wire.Contract,
		Subject:   *wire.Subject,
		Role:      *wire.Role,
		KeyID:     *wire.KeyID,
		IssuedAt:  *wire.IssuedAt,
		Signature: *wire.Signature,
	}, nil
}

// VerifyAttestation checks the attestation structure, recomputes and verifies
// the Ed25519 signature against the trusted bundle, and enforces role, the
// inclusive validity interval, and revocation at the explicit evaluation time.
// It uses public material only.
func VerifyAttestation(a *Attestation, bundle *TrustBundle, evalTime time.Time) error {
	if a == nil {
		return errf(CodeMalformed, "nil")
	}
	if a.Contract != AttestationContract {
		return errf(CodeBadContract, "contract")
	}
	if !artifact.ValidHash(a.Subject) {
		return errf(CodeBadSubject, "subject")
	}
	if !Roles()[a.Role] {
		return errf(CodeBadRole, "role")
	}
	if !artifact.ValidHash(a.KeyID) {
		return errf(CodeBadKeyID, "key_id")
	}
	issued, err := parseWholeSecondUTC(a.IssuedAt)
	if err != nil {
		return errf(CodeBadTime, "issued_at")
	}
	sig, err := decodeSignature(a.Signature)
	if err != nil {
		return errf(CodeBadSignature, "signature")
	}
	key, ok := findKey(bundle, a.KeyID)
	if !ok {
		return errf(CodeUnknownKey, "key_id")
	}
	if artifact.HashBytes(key.PublicKey) != a.KeyID {
		return errf(CodeKeyIDMismatch, "key_id")
	}
	if !hasRole(key, a.Role) {
		return errf(CodeRoleMismatch, "role")
	}
	if !ed25519.Verify(key.PublicKey, signedStatement(AttestationContract, a.Subject, a.Role, a.KeyID, a.IssuedAt), sig) {
		return errf(CodeSignature, "signature")
	}
	// Temporal and revocation checks at the explicit evaluation time.
	if issued.Before(key.NotBefore) || issued.After(key.NotAfter) {
		return errf(CodeNotYetValid, "issued_at")
	}
	if issued.After(evalTime) {
		return errf(CodeNotYetValid, "issued_at")
	}
	if evalTime.Before(key.NotBefore) || evalTime.After(key.NotAfter) {
		return errf(CodeExpired, "eval_time")
	}
	if key.Revoked {
		return errf(CodeRevoked, "key_id")
	}
	return nil
}

// ParseTrustBundle validates a strict atlas.trust/v1 bundle: exactly the
// contract and keys fields, distinct keys with matching computed key IDs,
// distinct non-empty known roles, valid inclusive time ranges, and at most
// maxKeys. No private material is accepted.
func ParseTrustBundle(raw []byte) (*TrustBundle, error) {
	if len(raw) > maxBundle {
		return nil, errf(CodeTooLarge, "bundle")
	}
	canon, err := artifact.Canonicalize(raw)
	if err != nil {
		var ae *artifact.Error
		if errorAsAE(err, &ae) {
			return nil, &Error{Code: ae.Code, Subject: ae.Subject}
		}
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	var wire struct {
		Contract *string    `json:"contract"`
		Keys     *[]keyWire `json:"keys"`
	}
	dec := json.NewDecoder(bytes.NewReader(canon))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&wire); err != nil {
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	if wire.Contract == nil {
		return nil, errf(CodeMissingField, "contract")
	}
	if *wire.Contract != TrustContract {
		return nil, errf(CodeBadContract, "contract")
	}
	if wire.Keys == nil {
		return nil, errf(CodeMissingField, "keys")
	}
	if len(*wire.Keys) > maxKeys {
		return nil, errf(CodeTooManyKeys, fmt.Sprintf("%d", len(*wire.Keys)))
	}
	bundle := &TrustBundle{Contract: TrustContract, Keys: make([]TrustKey, 0, len(*wire.Keys))}
	seen := map[string]bool{}
	for i, k := range *wire.Keys {
		if k.KeyID == nil || k.PublicKey == nil || k.Roles == nil ||
			k.NotBefore == nil || k.NotAfter == nil || k.Revoked == nil {
			return nil, errf(CodeMissingField, fmt.Sprintf("keys[%d]", i))
		}
		pub, err := base64.StdEncoding.DecodeString(*k.PublicKey)
		if err != nil || len(pub) != ed25519.PublicKeySize {
			return nil, errf(CodeBadPublicKey, fmt.Sprintf("keys[%d]", i))
		}
		if !artifact.ValidHash(*k.KeyID) {
			return nil, errf(CodeBadKeyID, fmt.Sprintf("keys[%d]", i))
		}
		if artifact.HashBytes(pub) != *k.KeyID {
			return nil, errf(CodeKeyIDMismatch, fmt.Sprintf("keys[%d]", i))
		}
		if seen[*k.KeyID] {
			return nil, errf(CodeDupKey, fmt.Sprintf("keys[%d]", i))
		}
		seen[*k.KeyID] = true
		if len(*k.Roles) == 0 {
			return nil, errf(CodeMissingField, fmt.Sprintf("keys[%d].roles", i))
		}
		roleSeen := map[string]bool{}
		for _, r := range *k.Roles {
			if !Roles()[r] {
				return nil, errf(CodeBadRole, fmt.Sprintf("keys[%d]", i))
			}
			if roleSeen[r] {
				return nil, errf(CodeDupRole, fmt.Sprintf("keys[%d]", i))
			}
			roleSeen[r] = true
		}
		nb, err := parseWholeSecondUTC(*k.NotBefore)
		if err != nil {
			return nil, errf(CodeBadTime, fmt.Sprintf("keys[%d].not_before", i))
		}
		na, err := parseWholeSecondUTC(*k.NotAfter)
		if err != nil {
			return nil, errf(CodeBadTime, fmt.Sprintf("keys[%d].not_after", i))
		}
		if nb.After(na) {
			return nil, errf(CodeBadTime, fmt.Sprintf("keys[%d]", i))
		}
		bundle.Keys = append(bundle.Keys, TrustKey{
			KeyID:     *k.KeyID,
			PublicKey: pub,
			Roles:     *k.Roles,
			NotBefore: nb,
			NotAfter:  na,
			Revoked:   *k.Revoked,
		})
	}
	return bundle, nil
}

type keyWire struct {
	KeyID     *string   `json:"key_id"`
	PublicKey *string   `json:"public_key"`
	Roles     *[]string `json:"roles"`
	NotBefore *string   `json:"not_before"`
	NotAfter  *string   `json:"not_after"`
	Revoked   *bool     `json:"revoked"`
}

func findKey(b *TrustBundle, keyID string) (TrustKey, bool) {
	if b == nil {
		return TrustKey{}, false
	}
	for _, k := range b.Keys {
		if k.KeyID == keyID {
			return k, true
		}
	}
	return TrustKey{}, false
}

func hasRole(k TrustKey, role string) bool {
	for _, r := range k.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// decodeSignature decodes a canonical padded standard base64 signature and
// requires exactly 64 bytes.
func decodeSignature(s string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(raw) != ed25519.SignatureSize {
		return nil, errf(CodeBadSignature, fmt.Sprintf("%d", len(raw)))
	}
	return raw, nil
}

// parseWholeSecondUTC parses a strict whole-second UTC RFC3339 timestamp of the
// form YYYY-MM-DDTHH:MM:SSZ.
func parseWholeSecondUTC(s string) (time.Time, error) {
	if len(s) != 20 || s[19] != 'Z' {
		return time.Time{}, errf(CodeBadTime, "format")
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		return time.Time{}, errf(CodeBadTime, "parse")
	}
	if t.Location() != time.UTC {
		return time.Time{}, errf(CodeBadTime, "not-utc")
	}
	return t, nil
}

// errorAsAE reports whether err is or wraps an *artifact.Error, storing it.
func errorAsAE(err error, target **artifact.Error) bool {
	var e *artifact.Error
	if errors.As(err, &e) {
		*target = e
		return true
	}
	return false
}
