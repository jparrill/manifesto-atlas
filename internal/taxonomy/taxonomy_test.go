package taxonomy

import (
	"testing"
)

var wantIDs = []string{
	"institutions-democracy", "administration-territory", "economy-industry-trade",
	"taxation-budgets", "employment-labour", "pensions", "healthcare",
	"education-research", "housing", "social-inclusion", "family-demography",
	"equality-rights", "immigration-asylum", "agriculture-rural",
	"environment-climate", "energy", "transport-infrastructure",
	"digital-science-innovation", "culture-language-sport", "justice-security",
	"defence-foreign-eu",
}

func TestTopicsRegistry(t *testing.T) {
	ts := Topics()
	if len(ts) != 21 {
		t.Fatalf("topics = %d, want 21", len(ts))
	}
	seen := map[string]bool{}
	for i, want := range wantIDs {
		if ts[i].Order != i+1 {
			t.Errorf("topic[%d].Order = %d, want %d", i, ts[i].Order, i+1)
		}
		if ts[i].ID != want {
			t.Errorf("topic[%d].ID = %q, want %q", i, ts[i].ID, want)
		}
		if ts[i].Label == "" {
			t.Errorf("topic[%d].Label empty", i)
		}
		if ts[i].Boundary == "" {
			t.Errorf("topic[%d].Boundary empty", i)
		}
		if seen[ts[i].ID] {
			t.Errorf("duplicate ID %q", ts[i].ID)
		}
		seen[ts[i].ID] = true
	}
	// TopicByID round-trips and rejects unknown.
	if _, ok := TopicByID("healthcare"); !ok {
		t.Error("TopicByID(healthcare) not found")
	}
	if _, ok := TopicByID("nonexistent"); ok {
		t.Error("TopicByID(nonexistent) should not be found")
	}
}

func TestValidateRelease(t *testing.T) {
	if err := ValidateRelease(Release); err != nil {
		t.Fatalf("approved release rejected: %v", err)
	}
	for _, v := range []string{"", "1.1.0", "0.9.0", "v1"} {
		if err := ValidateRelease(v); err == nil {
			t.Errorf("unknown/stale release %q accepted", v)
		}
	}
}

func TestStatesAndLabels(t *testing.T) {
	want := map[State]string{
		StateProposals:     "Contiene propuestas",
		StateMentioned:     "Mencionado sin propuestas concretas",
		StateNotCovered:    "No se aborda en el documento",
		StateNotApplicable: "Fuera del ámbito del documento",
		StateUnclear:       "No se puede determinar con fiabilidad",
	}
	if len(States()) != 5 {
		t.Fatalf("states = %d, want 5", len(States()))
	}
	for s, label := range want {
		got, err := StateLabel(s)
		if err != nil || got != label {
			t.Errorf("StateLabel(%s) = %q, %v; want %q", s, got, err, label)
		}
	}
}

func profileWith(relations map[string]Relation, sources []LegalSource) Profile {
	return Profile{
		Version:       Release,
		Scope:         ScopeNational,
		Geography:     "ES",
		Institution:   "test",
		EffectiveDate: "2026-01-01",
		Sources:       sources,
		Relations:     relations,
	}
}

func src() LegalSource {
	return LegalSource{Document: "BOE-A-1978-31229", Article: "149", Version: "1978", Hash: "sha256:ab"}
}

func TestClassifyStates(t *testing.T) {
	base := profileWith(map[string]Relation{"healthcare": RelShared}, []LegalSource{src()})
	cases := []struct {
		name string
		in   ClassifyInput
		want State
	}{
		{"proposals", ClassifyInput{TopicID: "healthcare", HasAction: true, FullCoverage: true, EvidenceReadable: true, Profile: base}, StateProposals},
		{"mentioned", ClassifyInput{TopicID: "healthcare", HasDiscussion: true, FullCoverage: true, EvidenceReadable: true, Profile: base}, StateMentioned},
		{"not-covered", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true, Profile: base}, StateNotCovered},
		{"unclear-unreadable", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: false, Profile: base}, StateUnclear},
		{"unclear-no-coverage", ClassifyInput{TopicID: "healthcare", FullCoverage: false, EvidenceReadable: true, Profile: base}, StateUnclear},
		{"unclear-unknown-relation", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true, Profile: profileWith(map[string]Relation{}, nil)}, StateUnclear},
	}
	for _, tc := range cases {
		got, err := Classify(tc.in)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%s: got %s, want %s", tc.name, got, tc.want)
		}
	}
}

func TestClassifyNotApplicable(t *testing.T) {
	// outside with positive versioned evidence and no content -> not-applicable.
	p := profileWith(map[string]Relation{"defence-foreign-eu": RelOutside}, []LegalSource{src()})
	got, err := Classify(ClassifyInput{TopicID: "defence-foreign-eu", FullCoverage: true, EvidenceReadable: true, Profile: p})
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if got != StateNotApplicable {
		t.Fatalf("got %s, want not-applicable", got)
	}
	// outside without legal evidence must be rejected (no implicit all-outside rule).
	pNoSrc := profileWith(map[string]Relation{"defence-foreign-eu": RelOutside}, nil)
	if _, err := Classify(ClassifyInput{TopicID: "defence-foreign-eu", FullCoverage: true, EvidenceReadable: true, Profile: pNoSrc}); err == nil {
		t.Fatal("outside without evidence accepted")
	}
}

func TestSharedPowersNationalAutonomous(t *testing.T) {
	// Shared healthcare competence: national and autonomous both applicable;
	// absence is not-covered in each (not not-applicable, not opposition).
	national := profileWith(map[string]Relation{"healthcare": RelShared}, []LegalSource{src()})
	national.Scope = ScopeNational
	autonomous := profileWith(map[string]Relation{"healthcare": RelShared}, []LegalSource{src()})
	autonomous.Scope = ScopeAutonomous
	for name, p := range map[string]Profile{"national": national, "autonomous": autonomous} {
		got, err := Classify(ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true, Profile: p})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got != StateNotCovered {
			t.Fatalf("%s: got %s, want not-covered", name, got)
		}
	}
}

func TestGeneralScopeApplicable(t *testing.T) {
	// A broad general-political programme keeps all topics applicable; absence is
	// not-covered.
	p := profileWith(map[string]Relation{"energy": RelDirect}, []LegalSource{src()})
	p.Scope = ScopeGeneral
	got, err := Classify(ClassifyInput{TopicID: "energy", FullCoverage: true, EvidenceReadable: true, Profile: p})
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if got != StateNotCovered {
		t.Fatalf("got %s, want not-covered", got)
	}
}

func TestValidateProfile(t *testing.T) {
	valid := profileWith(map[string]Relation{"healthcare": RelDirect}, []LegalSource{src()})
	if err := ValidateProfile(valid); err != nil {
		t.Fatalf("valid profile rejected: %v", err)
	}
	// Missing version.
	noVer := valid
	noVer.Version = ""
	if err := ValidateProfile(noVer); err == nil {
		t.Error("missing version accepted")
	}
	// Bad scope.
	badScope := valid
	badScope.Scope = Scope("galactic")
	if err := ValidateProfile(badScope); err == nil {
		t.Error("bad scope accepted")
	}
	// Unknown topic in relations.
	unknownTopic := valid
	unknownTopic.Relations = map[string]Relation{"bogus": RelDirect}
	if err := ValidateProfile(unknownTopic); err == nil {
		t.Error("unknown topic in relations accepted")
	}
	// Unknown relation.
	unknownRel := valid
	unknownRel.Relations = map[string]Relation{"healthcare": Relation("maybe")}
	if err := ValidateProfile(unknownRel); err == nil {
		t.Error("unknown relation accepted")
	}
	// Unpinned source.
	unpinned := valid
	unpinned.Sources = []LegalSource{LegalSource{Document: "x", Article: "", Version: "1", Hash: "h"}}
	if err := ValidateProfile(unpinned); err == nil {
		t.Error("unpinned source accepted")
	}
}

func TestCompatibleComparison(t *testing.T) {
	a := profileWith(nil, nil)
	b := profileWith(nil, nil)
	if err := CompatibleComparison(a, b); err != nil {
		t.Fatalf("same version+scope rejected: %v", err)
	}
	diffVer := b
	diffVer.Version = "1.1.0"
	if err := CompatibleComparison(a, diffVer); err == nil {
		t.Error("different taxonomy version accepted")
	}
	diffScope := b
	diffScope.Scope = ScopeAutonomous
	if err := CompatibleComparison(a, diffScope); err == nil {
		t.Error("national vs regional comparison accepted")
	}
}

func sameMappings() []Mapping {
	maps := make([]Mapping, 0, len(topics))
	for _, t := range topics {
		maps = append(maps, Mapping{OldTopicID: t.ID, NewTopicID: t.ID, Kind: "same"})
	}
	return maps
}

func TestMigration(t *testing.T) {
	// Unchanged hash: not a real new release.
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: Release, OldHash: "sha256:a", NewHash: "sha256:a", Mappings: sameMappings()}); err == nil {
		t.Error("unchanged hash accepted")
	}
	// Empty mappings.
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: "1.1.0", OldHash: "sha256:a", NewHash: "sha256:b", Mappings: nil}); err == nil {
		t.Error("empty migration accepted")
	}
	// Incomplete mapping (drop one topic).
	incomplete := sameMappings()[:20]
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: "1.1.0", OldHash: "sha256:a", NewHash: "sha256:b", Mappings: incomplete}); err == nil {
		t.Error("incomplete migration accepted")
	}
	// ID reuse: two old topics (non-merge) map to the same new ID.
	reuse := sameMappings()
	reuse[1].NewTopicID = reuse[0].NewTopicID
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: "1.1.0", OldHash: "sha256:a", NewHash: "sha256:b", Mappings: reuse}); err == nil {
		t.Error("id reuse accepted")
	}
	// Structurally complete but unknown target release: fails closed.
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: "1.1.0", OldHash: "sha256:a", NewHash: "sha256:b", Mappings: sameMappings()}); err == nil {
		t.Error("unknown target release accepted")
	}
	// Identical release.
	if err := ValidateMigration(Migration{FromRelease: Release, ToRelease: Release, OldHash: "sha256:a", NewHash: "sha256:b", Mappings: sameMappings()}); err == nil {
		t.Error("identical release migration accepted")
	}
}
