package provenance_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/provenance"
)

func sp(s string) *string { return &s }

func mkHash(b byte) string { return "sha256:" + strings.Repeat(string(b), 64) }

var (
	grantReviewed = "2030-01-01T00:00:00Z"
	grantUntil    = "2030-06-01T00:00:00Z"
	nowMid        = time.Date(2030, 3, 1, 0, 0, 0, 0, time.UTC)
)

func newRegistry() provenance.Registry {
	return provenance.Registry{
		Version:    provenance.ProvenanceVersion,
		RegistryID: "reg-1",
		ReviewedAt: grantReviewed,
		Parties: []provenance.Party{
			{
				Version:     provenance.PartyVersion,
				PartyID:     "fixture-party",
				LegalName:   "Partido de Prueba S.L.",
				DisplayName: "Partido de Prueba",
				DomainGrants: []provenance.DomainGrant{
					{Hostname: "oficial.fixture-party.es", PartyID: "fixture-party", Scope: "autonomous",
						Evidence: []string{mkHash('0')}, ReviewedAt: grantReviewed, ValidUntil: grantUntil, State: "active"},
					{Hostname: "castillayleon.fixture-party.es", PartyID: "fixture-party", Scope: "autonomous",
						Evidence: []string{mkHash('1')}, ReviewedAt: grantReviewed, ValidUntil: grantUntil, State: "active"},
				},
			},
		},
	}
}

func code(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	ce, ok := err.(*provenance.Error)
	if !ok {
		t.Fatalf("expected provenance error, got %T: %v", err, err)
	}
	return ce.Code
}

func TestCheckAnchorResolves(t *testing.T) {
	a, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "autonomous", newRegistry(), nowMid)
	if err != nil {
		t.Fatalf("anchor rejected: %v", err)
	}
	if a.Hostname != "oficial.fixture-party.es" || a.PartyID != "fixture-party" || a.Scope != "autonomous" {
		t.Fatalf("unexpected anchor: %+v", a)
	}
	if len(a.Evidence) != 1 {
		t.Fatalf("evidence = %d", len(a.Evidence))
	}
}

func TestRegionalDomain(t *testing.T) {
	a, err := provenance.CheckAnchor("fixture-party", "castillayleon.fixture-party.es", "autonomous", newRegistry(), nowMid)
	if err != nil {
		t.Fatalf("regional anchor rejected: %v", err)
	}
	if a.Hostname != "castillayleon.fixture-party.es" {
		t.Fatalf("hostname = %q", a.Hostname)
	}
}

func TestUnknownParty(t *testing.T) {
	_, err := provenance.CheckAnchor("other-party", "oficial.fixture-party.es", "autonomous", newRegistry(), nowMid)
	if code(t, err) != provenance.CodeUnknownParty {
		t.Fatalf("wrong code for unknown party")
	}
}

func TestLookalike(t *testing.T) {
	reg := newRegistry()
	for _, host := range []string{
		"oficial.fixture-partyx.es",         // typo
		"www.oficial.fixture-party.es",      // subdomain, no automatic www trust
		"fixture-party.es",                  // suffix, no inference
		"evil.oficial.fixture-party.es",     // lookalike subdomain
		"oficial.fixture-party.es.evil.com", // appended domain
	} {
		if _, err := provenance.CheckAnchor("fixture-party", host, "autonomous", reg, nowMid); code(t, err) != provenance.CodeUnknownDomain {
			t.Fatalf("lookalike %q not rejected as unknown domain", host)
		}
	}
}

func TestScopeMismatch(t *testing.T) {
	_, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "national", newRegistry(), nowMid)
	if code(t, err) != provenance.CodeScopeMismatch {
		t.Fatal("wrong code for scope mismatch")
	}
}

func TestRevoked(t *testing.T) {
	reg := newRegistry()
	reg.Parties[0].DomainGrants[0].State = "revoked"
	if _, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "autonomous", reg, nowMid); code(t, err) != provenance.CodeRevoked {
		t.Fatal("wrong code for revoked grant")
	}
}

func TestExpired(t *testing.T) {
	later := time.Date(2030, 7, 1, 0, 0, 0, 0, time.UTC)
	if _, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "autonomous", newRegistry(), later); code(t, err) != provenance.CodeExpired {
		t.Fatal("wrong code for expired grant")
	}
}

func TestNotYetReviewed(t *testing.T) {
	earlier := time.Date(2029, 12, 1, 0, 0, 0, 0, time.UTC)
	if _, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "autonomous", newRegistry(), earlier); code(t, err) != provenance.CodeNotYetReviewed {
		t.Fatal("wrong code for not-yet-reviewed grant")
	}
}

func TestEvidenceIntegrity(t *testing.T) {
	reg := newRegistry()
	reg.Parties[0].DomainGrants[0].Evidence = []string{"not-a-hash"}
	if _, err := provenance.CheckAnchor("fixture-party", "oficial.fixture-party.es", "autonomous", reg, nowMid); code(t, err) != provenance.CodeEvidenceIntegrity {
		t.Fatal("wrong code for bad evidence")
	}
}

func TestConflictingIdentity(t *testing.T) {
	reg := newRegistry()
	// A second party claims the same hostname.
	reg.Parties = append(reg.Parties, provenance.Party{
		Version: provenance.PartyVersion, PartyID: "rival-party", LegalName: "Rival",
		DomainGrants: []provenance.DomainGrant{
			{Hostname: "oficial.fixture-party.es", PartyID: "rival-party", Scope: "autonomous",
				Evidence: []string{mkHash('2')}, ReviewedAt: grantReviewed, ValidUntil: grantUntil, State: "active"},
		},
	})
	// Requesting the rival for the hostname owned by fixture-party is a conflict.
	if _, err := provenance.CheckAnchor("rival-party", "oficial.fixture-party.es", "autonomous", reg, nowMid); code(t, err) != provenance.CodeConflictingIdentity {
		t.Fatal("wrong code for conflicting identity")
	}
}

func TestValidateRegistryValid(t *testing.T) {
	if err := provenance.ValidateRegistry(newRegistry()); err != nil {
		t.Fatalf("valid registry rejected: %v", err)
	}
}

func TestValidateRegistryRejects(t *testing.T) {
	// Unknown version.
	reg := newRegistry()
	reg.Version = "atlas.provenance/v2"
	if code(t, provenance.ValidateRegistry(reg)) != provenance.CodeUnknownVersion {
		t.Fatal("unknown version not rejected")
	}
	// Duplicate party.
	reg = newRegistry()
	reg.Parties = append(reg.Parties, reg.Parties[0])
	if code(t, provenance.ValidateRegistry(reg)) != provenance.CodeDuplicateParty {
		t.Fatal("duplicate party not rejected")
	}
	// Unicode hostname (no IDNA guessing).
	reg = newRegistry()
	reg.Parties[0].DomainGrants[0].Hostname = "ofícial.fixture-party.es"
	if code(t, provenance.ValidateRegistry(reg)) != provenance.CodeBadHostname {
		t.Fatal("unicode hostname not rejected")
	}
	// Cross-party hostname conflict.
	reg = newRegistry()
	reg.Parties = append(reg.Parties, provenance.Party{
		Version: provenance.PartyVersion, PartyID: "rival-party", LegalName: "Rival",
		DomainGrants: []provenance.DomainGrant{
			{Hostname: "oficial.fixture-party.es", PartyID: "rival-party", Scope: "autonomous",
				Evidence: []string{mkHash('2')}, ReviewedAt: grantReviewed, ValidUntil: grantUntil, State: "active"},
		},
	})
	if code(t, provenance.ValidateRegistry(reg)) != provenance.CodeConflictingIdentity {
		t.Fatal("cross-party hostname conflict not rejected")
	}
	// Grant longer than 180 days.
	reg = newRegistry()
	reg.Parties[0].DomainGrants[0].ValidUntil = "2031-01-01T00:00:00Z" // 365 days
	if code(t, provenance.ValidateRegistry(reg)) != provenance.CodeGrantTooLong {
		t.Fatal("over-long grant not rejected")
	}
}

func TestVersionInvalidation(t *testing.T) {
	party, host, scope := "fixture-party", "oficial.fixture-party.es", "autonomous"
	v1 := newRegistry()

	// Identical registry: no change.
	if provenance.AnchorChanged(v1, newRegistry(), party, host, scope) {
		t.Fatal("identical registries reported as changed")
	}

	// v2 supersedes v1 and revokes the grant: dependent provenance invalidated.
	v2 := newRegistry()
	v2.RegistryID = "reg-2"
	v2.Supersedes = sp("reg-1")
	v2.Parties[0].DomainGrants[0].State = "revoked"
	if !provenance.AnchorChanged(v1, v2, party, host, scope) {
		t.Fatal("revoked grant not reported as changed")
	}

	// v2 removes the grant entirely: invalidated.
	v3 := newRegistry()
	v3.RegistryID = "reg-3"
	v3.Supersedes = sp("reg-1")
	v3.Parties[0].DomainGrants = v3.Parties[0].DomainGrants[1:]
	if !provenance.AnchorChanged(v1, v3, party, host, scope) {
		t.Fatal("removed grant not reported as changed")
	}

	// v2 keeps the exact grant: not invalidated.
	v4 := newRegistry()
	v4.RegistryID = "reg-4"
	v4.Supersedes = sp("reg-1")
	if provenance.AnchorChanged(v1, v4, party, host, scope) {
		t.Fatal("unchanged grant reported as changed")
	}
}
