// Package provenance implements the atlas.provenance/v1 trust model. This file
// (TQ-0023) provides the versioned party/territory/domain evidence registry and
// the deterministic, fail-closed anchor authority check.
//
// Provenance resolves party-controlled domains through an explicitly selected,
// evidence-backed registry rather than arbitrary submitted URLs. No wildcard,
// substring, suffix, www, or subdomain trust is ever inferred. A resolved anchor
// is an authority assertion, not a semantic approval; independent download and
// byte-identity checks (other provenance tasks) still apply.
package provenance

import (
	"net"
	"strings"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Contract versions.
const (
	ProvenanceVersion = "atlas.provenance/v1"
	PartyVersion      = "atlas.party/v1"
)

// maxGrantLifetime bounds a domain grant to no more than 180 days after review.
const maxGrantLifetime = 180 * 24 * time.Hour

// Error codes.
const (
	CodeUnknownVersion      = "unknown-version"
	CodeBadPartyID          = "bad-party-id"
	CodeDuplicateParty      = "duplicate-party"
	CodeBadHostname         = "bad-hostname"
	CodeMissingField        = "missing-field"
	CodeBadState            = "bad-state"
	CodeBadDate             = "bad-date"
	CodeGrantTooLong        = "grant-too-long"
	CodeEvidenceIntegrity   = "evidence-integrity"
	CodeUnknownParty        = "unknown-party"
	CodeUnknownDomain       = "unknown-domain"
	CodeConflictingIdentity = "conflicting-identity"
	CodeScopeMismatch       = "scope-mismatch"
	CodeRevoked             = "revoked"
	CodeExpired             = "expired"
	CodeNotYetReviewed      = "not-yet-reviewed"
)

// Error is a stable provenance error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "provenance: " + e.Code
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

// Relationship records an explicit merger, coalition, or renamed legal entity.
type Relationship struct {
	Kind         string `json:"kind"`
	OtherPartyID string `json:"other_party_id"`
}

// DomainGrant is one exact-hostname authority grant for a party and scope.
type DomainGrant struct {
	Hostname   string   `json:"hostname"`
	PartyID    string   `json:"party_id"`
	Scope      string   `json:"scope"`
	Evidence   []string `json:"evidence"`
	ReviewedAt string   `json:"reviewed_at"`
	ValidUntil string   `json:"valid_until"`
	State      string   `json:"state"`
}

// KeyGrant is one publisher-signature key grant for a party and scope.
type KeyGrant struct {
	KeyID      string   `json:"key_id"`
	PublicKey  string   `json:"public_key"`
	Scope      string   `json:"scope"`
	Evidence   []string `json:"evidence"`
	ReviewedAt string   `json:"reviewed_at"`
	ValidUntil string   `json:"valid_until"`
	State      string   `json:"state"`
}

// Party is one immutable atlas.party/v1 registry entry.
type Party struct {
	Version      string        `json:"version"`
	PartyID      string        `json:"party_id"`
	LegalName    string        `json:"legal_name"`
	DisplayName  string        `json:"display_name"`
	ElectoralID  *string       `json:"electoral_id"`
	Relationship *Relationship `json:"relationship"`
	DomainGrants []DomainGrant `json:"domain_grants"`
	KeyGrants    []KeyGrant    `json:"key_grants"`
}

// Registry is the versioned trusted party registry the caller explicitly selects.
type Registry struct {
	Version    string  `json:"version"`
	RegistryID string  `json:"registry_id"`
	ReviewedAt string  `json:"reviewed_at"`
	Supersedes *string `json:"supersedes"`
	Parties    []Party `json:"parties"`
}

// Anchor is a resolved, validated domain grant.
type Anchor struct {
	PartyID    string
	Hostname   string
	Scope      string
	Evidence   []string
	ReviewedAt string
	ValidUntil string
}

// CheckAnchor resolves publisher authority for (party, hostname, scope) against
// an explicitly selected registry at evaluation time. It is deterministic and
// fail-closed: only an exact active in-window grant with intact evidence
// resolves; lookalikes, conflicts, revocation, expiry, and unknown domains do
// not yield verified provenance.
func CheckAnchor(partyID, hostname, scope string, reg Registry, now time.Time) (*Anchor, error) {
	if reg.Version != ProvenanceVersion {
		return nil, errf(CodeUnknownVersion, reg.Version)
	}
	if err := validPartyID(partyID); err != nil {
		return nil, err
	}
	if err := validAHostname(hostname); err != nil {
		return nil, err
	}
	partyKnown := false
	for i := range reg.Parties {
		if reg.Parties[i].PartyID == partyID {
			partyKnown = true
			break
		}
	}
	if !partyKnown {
		return nil, errf(CodeUnknownParty, partyID)
	}
	var matches []*DomainGrant
	var owner string
	for i := range reg.Parties {
		p := &reg.Parties[i]
		for j := range p.DomainGrants {
			g := &p.DomainGrants[j]
			if g.Hostname == hostname {
				matches = append(matches, g)
				owner = p.PartyID
			}
		}
	}
	if len(matches) == 0 {
		return nil, errf(CodeUnknownDomain, hostname)
	}
	// The hostname is owned by a different party, or duplicated: a conflict.
	if owner != partyID || len(matches) > 1 {
		return nil, errf(CodeConflictingIdentity, hostname)
	}
	g := matches[0]
	if g.Scope != scope {
		return nil, errf(CodeScopeMismatch, scope)
	}
	if g.State != "active" {
		return nil, errf(CodeRevoked, hostname)
	}
	if err := checkEvidence(g.Evidence); err != nil {
		return nil, err
	}
	reviewed, err := parseTime(g.ReviewedAt)
	if err != nil {
		return nil, errf(CodeBadDate, g.ReviewedAt)
	}
	until, err := parseTime(g.ValidUntil)
	if err != nil {
		return nil, errf(CodeBadDate, g.ValidUntil)
	}
	if now.Before(reviewed) {
		return nil, errf(CodeNotYetReviewed, hostname)
	}
	if now.After(until) {
		return nil, errf(CodeExpired, hostname)
	}
	return &Anchor{
		PartyID:    partyID,
		Hostname:   hostname,
		Scope:      g.Scope,
		Evidence:   g.Evidence,
		ReviewedAt: g.ReviewedAt,
		ValidUntil: g.ValidUntil,
	}, nil
}

// ValidateRegistry checks the structural integrity of a trusted registry:
// versions, unique party IDs, valid hostnames/scopes/dates, the 180-day grant
// bound, evidence integrity, and cross-party hostname uniqueness.
func ValidateRegistry(reg Registry) error {
	if reg.Version != ProvenanceVersion {
		return errf(CodeUnknownVersion, reg.Version)
	}
	if reg.RegistryID == "" {
		return errf(CodeMissingField, "registry_id")
	}
	if _, err := parseTime(reg.ReviewedAt); err != nil {
		return errf(CodeBadDate, reg.ReviewedAt)
	}
	if reg.Supersedes != nil && *reg.Supersedes == reg.RegistryID {
		return errf(CodeMissingField, "supersedes")
	}
	seenParty := map[string]bool{}
	hostnameOwner := map[string]string{}
	for i := range reg.Parties {
		p := &reg.Parties[i]
		if p.Version != PartyVersion {
			return errf(CodeUnknownVersion, p.Version)
		}
		if err := validPartyID(p.PartyID); err != nil {
			return err
		}
		if seenParty[p.PartyID] {
			return errf(CodeDuplicateParty, p.PartyID)
		}
		seenParty[p.PartyID] = true
		if p.LegalName == "" {
			return errf(CodeMissingField, "legal_name")
		}
		for j := range p.DomainGrants {
			g := &p.DomainGrants[j]
			if err := validAHostname(g.Hostname); err != nil {
				return err
			}
			if g.PartyID != p.PartyID {
				return errf(CodeMissingField, "grant_party")
			}
			if g.Scope == "" {
				return errf(CodeMissingField, "scope")
			}
			if g.State != "active" && g.State != "revoked" {
				return errf(CodeBadState, g.State)
			}
			if err := checkEvidence(g.Evidence); err != nil {
				return err
			}
			reviewed, err := parseTime(g.ReviewedAt)
			if err != nil {
				return errf(CodeBadDate, g.ReviewedAt)
			}
			until, err := parseTime(g.ValidUntil)
			if err != nil {
				return errf(CodeBadDate, g.ValidUntil)
			}
			if until.Before(reviewed) {
				return errf(CodeBadDate, g.ValidUntil)
			}
			if until.Sub(reviewed) > maxGrantLifetime {
				return errf(CodeGrantTooLong, g.Hostname)
			}
			if owner, ok := hostnameOwner[g.Hostname]; ok && owner != p.PartyID {
				return errf(CodeConflictingIdentity, g.Hostname)
			}
			hostnameOwner[g.Hostname] = p.PartyID
		}
	}
	return nil
}

// AnchorChanged reports whether the resolved grant for (party, hostname, scope)
// differs between two registry versions. Any difference — absence, revocation,
// scope, evidence, or validity window — invalidates provenance that depended on
// the older grant. Historical records are preserved; this only decides
// invalidation, it never rewrites them.
func AnchorChanged(from, to Registry, partyID, hostname, scope string) bool {
	a, okA := resolveGrant(from, partyID, hostname, scope)
	b, okB := resolveGrant(to, partyID, hostname, scope)
	if okA != okB {
		return true
	}
	if !okA {
		return false
	}
	return !grantEqual(a, b)
}

func resolveGrant(reg Registry, partyID, hostname, scope string) (*DomainGrant, bool) {
	for i := range reg.Parties {
		p := &reg.Parties[i]
		if p.PartyID != partyID {
			continue
		}
		for j := range p.DomainGrants {
			g := &p.DomainGrants[j]
			if g.Hostname == hostname && g.Scope == scope {
				return g, true
			}
		}
	}
	return nil, false
}

func grantEqual(a, b *DomainGrant) bool {
	return a.Hostname == b.Hostname &&
		a.PartyID == b.PartyID &&
		a.Scope == b.Scope &&
		a.ReviewedAt == b.ReviewedAt &&
		a.ValidUntil == b.ValidUntil &&
		a.State == b.State &&
		equalStrings(a.Evidence, b.Evidence)
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkEvidence(e []string) error {
	if len(e) == 0 {
		return errf(CodeEvidenceIntegrity, "empty")
	}
	for _, id := range e {
		if !artifact.ValidHash(id) {
			return errf(CodeEvidenceIntegrity, id)
		}
	}
	return nil
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func validPartyID(id string) error {
	if len(id) < 2 || len(id) > 64 {
		return errf(CodeBadPartyID, id)
	}
	if id[0] < 'a' || id[0] > 'z' {
		return errf(CodeBadPartyID, id)
	}
	for i := 1; i < len(id); i++ {
		c := id[i]
		if !(c == '-' || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return errf(CodeBadPartyID, id)
		}
	}
	return nil
}

// validAHostname enforces an exact lowercase ASCII A-label: no wildcard,
// substring, IP literal, Unicode (no IDNA guessing), or malformed labels.
func validAHostname(h string) error {
	if h == "" || len(h) > 253 {
		return errf(CodeBadHostname, h)
	}
	for i := 0; i < len(h); i++ {
		c := h[i]
		if c >= 0x80 {
			return errf(CodeBadHostname, h)
		}
		if !(c == '.' || c == '-' || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return errf(CodeBadHostname, h)
		}
	}
	if strings.HasPrefix(h, ".") || strings.HasSuffix(h, ".") {
		return errf(CodeBadHostname, h)
	}
	for _, label := range strings.Split(h, ".") {
		if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
			return errf(CodeBadHostname, h)
		}
	}
	if net.ParseIP(h) != nil {
		return errf(CodeBadHostname, h)
	}
	return nil
}
