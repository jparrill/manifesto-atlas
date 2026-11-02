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

func TestClassifyStates(t *testing.T) {
	cases := []struct {
		name string
		in   ClassifyInput
		want State
	}{
		{"proposals", ClassifyInput{TopicID: "healthcare", HasAction: true, FullCoverage: true, EvidenceReadable: true, Applicability: AppInScope}, StateProposals},
		{"mentioned", ClassifyInput{TopicID: "healthcare", HasDiscussion: true, FullCoverage: true, EvidenceReadable: true, Applicability: AppInScope}, StateMentioned},
		{"not-covered", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true, Applicability: AppInScope}, StateNotCovered},
		{"unclear-unreadable", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: false, Applicability: AppInScope}, StateUnclear},
		{"unclear-no-coverage", ClassifyInput{TopicID: "healthcare", FullCoverage: false, EvidenceReadable: true, Applicability: AppInScope}, StateUnclear},
		{"unclear-unknown-applicability", ClassifyInput{TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true, Applicability: AppUnknown}, StateUnclear},
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
	// Out-of-scope with retained scope evidence and no content -> not-applicable.
	got, err := Classify(ClassifyInput{
		TopicID: "defence-foreign-eu", FullCoverage: true, EvidenceReadable: true,
		Applicability: AppOutOfScope,
		ScopeEvidence: &ScopeEvidence{Path: "evidence/scope.html", SHA256: "sha256:" + "ab"},
	})
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if got != StateNotApplicable {
		t.Fatalf("got %s, want not-applicable", got)
	}
	// Out-of-scope without scope evidence must be rejected (no implicit
	// all-outside rule, no guessed exclusion).
	if _, err := Classify(ClassifyInput{
		TopicID: "defence-foreign-eu", FullCoverage: true, EvidenceReadable: true,
		Applicability: AppOutOfScope,
	}); err == nil {
		t.Fatal("out-of-scope without evidence accepted")
	}
	// Out-of-scope with an empty evidence reference must be rejected.
	if _, err := Classify(ClassifyInput{
		TopicID: "defence-foreign-eu", FullCoverage: true, EvidenceReadable: true,
		Applicability: AppOutOfScope,
		ScopeEvidence: &ScopeEvidence{},
	}); err == nil {
		t.Fatal("out-of-scope with empty evidence accepted")
	}
}

func TestSharedAndGeneralScopeNotExcluded(t *testing.T) {
	// Shared powers (healthcare) and broad general-political topics are in-scope:
	// absence is not-covered, never a blanket not-applicable exclusion.
	for _, tc := range []struct {
		name    string
		topicID string
	}{
		{"shared-healthcare", "healthcare"},
		{"general-political-energy", "energy"},
	} {
		got, err := Classify(ClassifyInput{
			TopicID: tc.topicID, FullCoverage: true, EvidenceReadable: true,
			Applicability: AppInScope,
		})
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got != StateNotCovered {
			t.Fatalf("%s: got %s, want not-covered", tc.name, got)
		}
	}
}

func TestClassifyUnknownTopic(t *testing.T) {
	if _, err := Classify(ClassifyInput{
		TopicID: "bogus-topic", FullCoverage: true, EvidenceReadable: true,
		Applicability: AppInScope,
	}); err == nil {
		t.Fatal("unknown topic accepted")
	}
}

func TestClassifyBadApplicability(t *testing.T) {
	if _, err := Classify(ClassifyInput{
		TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true,
		Applicability: Applicability("maybe"),
	}); err == nil {
		t.Fatal("unknown applicability value accepted")
	}
}
