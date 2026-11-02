// Package taxonomy implements the atlas.taxonomy/v1 contract (release 1.0.0)
// established by docs/design/taxonomy-jurisdiction.md (TQ-0013) for TQ-0036.
//
// It provides the immutable ordered 21-topic registry, the five mutually
// exclusive topic states with their Spanish public labels, versioned
// jurisdiction/scope profiles with pinned legal evidence, deterministic
// classification, comparison compatibility, and explicit migration validation.
// All checks are offline and fail-closed; no model may add, remove, rename, or
// reorder topics.
package taxonomy

import (
	"fmt"
)

// Release is the single approved v1 taxonomy release.
const Release = "1.0.0"

// Error codes. Stable and free of private data.
const (
	CodeUnknownTopic     = "unknown-topic"
	CodeMissingField     = "missing-field"
	CodeBadScope         = "bad-scope"
	CodeBadRelation      = "bad-relation"
	CodeUnknownVersion   = "unknown-version"
	CodeVersionMismatch  = "version-mismatch"
	CodeScopeMismatch    = "scope-mismatch"
	CodeMigrationIncompl = "migration-incomplete"
	CodeIDReuse          = "id-reuse"
	CodeBadHash          = "bad-hash"
)

// Error is a stable taxonomy error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "taxonomy: " + e.Code
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

// Topic is one immutable v1 topic.
type Topic struct {
	Order    int
	ID       string
	Label    string // Spanish public label
	Boundary string // English technical boundary
}

// topics is the exact, ordered v1 registry. It must not be modified in place.
var topics = []Topic{
	{1, "institutions-democracy", "Instituciones y democracia", "Electoral institutions, participation, accountability"},
	{2, "administration-territory", "Administración pública y organización territorial", "Administrative services and territorial organization"},
	{3, "economy-industry-trade", "Economía, industria y comercio", "Economic production, industrial and trade policy"},
	{4, "taxation-budgets", "Fiscalidad y presupuestos públicos", "Taxes, revenue, budget allocation and public debt"},
	{5, "employment-labour", "Empleo y relaciones laborales", "Employment policy, working conditions and labour relations"},
	{6, "pensions", "Pensiones", "Retirement and other pension systems"},
	{7, "healthcare", "Sanidad", "Health systems, public health and care provision"},
	{8, "education-research", "Educación, universidades e investigación", "Education institutions and academic research"},
	{9, "housing", "Vivienda", "Housing access, supply, tenure and residential policy"},
	{10, "social-inclusion", "Servicios sociales e inclusión", "Social assistance, dependency and inclusion"},
	{11, "family-demography", "Familia y demografía", "Family support and demographic policy"},
	{12, "equality-rights", "Igualdad y derechos civiles", "Equal treatment, freedoms and civil rights"},
	{13, "immigration-asylum", "Inmigración, asilo e integración", "Migration status, asylum and integration"},
	{14, "agriculture-rural", "Agricultura, pesca y medio rural", "Agriculture, fisheries and rural affairs"},
	{15, "environment-climate", "Medio ambiente y clima", "Environmental protection and climate policy"},
	{16, "energy", "Energía", "Energy production, supply and consumption policy"},
	{17, "transport-infrastructure", "Transporte e infraestructuras", "Mobility, transport networks and infrastructure"},
	{18, "digital-science-innovation", "Política digital, ciencia e innovación", "Digital systems, technological development and innovation"},
	{19, "culture-language-sport", "Cultura, lengua y deporte", "Cultural, linguistic and sporting policy"},
	{20, "justice-security", "Justicia y seguridad", "Justice institutions, public safety and security"},
	{21, "defence-foreign-eu", "Defensa, asuntos exteriores y Unión Europea", "Defence, external relations and European affairs"},
}

// Topics returns the ordered v1 registry (21 topics). The slice must not be
// mutated.
func Topics() []Topic {
	out := make([]Topic, len(topics))
	copy(out, topics)
	return out
}

// TopicByID returns the topic with the given immutable ID.
func TopicByID(id string) (Topic, bool) {
	for _, t := range topics {
		if t.ID == id {
			return t, true
		}
	}
	return Topic{}, false
}

// ValidateRelease reports whether v is an approved v1 taxonomy release. Unknown
// or stale versions are rejected; there is no best-effort reinterpretation.
func ValidateRelease(v string) error {
	if v != Release {
		return errf(CodeUnknownVersion, v)
	}
	return nil
}

// State is one of the five mutually exclusive topic states.
type State string

// Topic states.
const (
	StateProposals     State = "proposals"
	StateMentioned     State = "mentioned"
	StateNotCovered    State = "not-covered"
	StateNotApplicable State = "not-applicable"
	StateUnclear       State = "unclear"
)

// States returns the five states.
func States() []State {
	return []State{StateProposals, StateMentioned, StateNotCovered, StateNotApplicable, StateUnclear}
}

// StateLabel returns the Spanish public label for a state.
func StateLabel(s State) (string, error) {
	switch s {
	case StateProposals:
		return "Contiene propuestas", nil
	case StateMentioned:
		return "Mencionado sin propuestas concretas", nil
	case StateNotCovered:
		return "No se aborda en el documento", nil
	case StateNotApplicable:
		return "Fuera del ámbito del documento", nil
	case StateUnclear:
		return "No se puede determinar con fiabilidad", nil
	default:
		return "", errf(CodeBadRelation, string(s))
	}
}

// Relation is a topic/subtopic competence relation in a scope profile.
type Relation string

// Competence relations.
const (
	RelDirect    Relation = "direct"
	RelShared    Relation = "shared"
	RelDelegated Relation = "delegated"
	RelAdvocacy  Relation = "advocacy-only"
	RelOutside   Relation = "outside"
	RelUnknown   Relation = "unknown"
)

// Relations returns the six competence relations.
func Relations() []Relation {
	return []Relation{RelDirect, RelShared, RelDelegated, RelAdvocacy, RelOutside, RelUnknown}
}

func knownRelation(r Relation) bool {
	for _, k := range Relations() {
		if k == r {
			return true
		}
	}
	return false
}

// Scope is a document scope category.
type Scope string

// Document scopes.
const (
	ScopeNational   Scope = "national"
	ScopeAutonomous Scope = "autonomous"
	ScopeProvincial Scope = "provincial"
	ScopeMunicipal  Scope = "municipal"
	ScopeEuropean   Scope = "european"
	ScopeGeneral    Scope = "general-political"
)

func knownScope(s Scope) bool {
	switch s {
	case ScopeNational, ScopeAutonomous, ScopeProvincial, ScopeMunicipal, ScopeEuropean, ScopeGeneral:
		return true
	default:
		return false
	}
}

// LegalSource is a pinned legal evidence reference.
type LegalSource struct {
	Document string
	Article  string
	Version  string
	Hash     string
}

// Profile is a versioned jurisdiction/scope configuration resource.
type Profile struct {
	Version       string
	Scope         Scope
	Geography     string
	Institution   string
	EffectiveDate string
	Sources       []LegalSource
	Relations     map[string]Relation
}

// ValidateProfile checks that a profile is versioned, has a known scope and
// effective date, pinned legal sources, and only known topics/relations. It does
// not accept a missing profile as an implicit all-outside rule.
func ValidateProfile(p Profile) error {
	if p.Version == "" {
		return errf(CodeMissingField, "version")
	}
	if !knownScope(p.Scope) {
		return errf(CodeBadScope, string(p.Scope))
	}
	if p.EffectiveDate == "" {
		return errf(CodeMissingField, "effective_date")
	}
	for i, s := range p.Sources {
		if s.Document == "" || s.Article == "" || s.Version == "" || s.Hash == "" {
			return errf(CodeMissingField, fmt.Sprintf("sources[%d]", i))
		}
	}
	// An outside decision requires positive, versioned legal evidence.
	hasOutside := false
	for id, r := range p.Relations {
		if _, ok := TopicByID(id); !ok {
			return errf(CodeUnknownTopic, id)
		}
		if !knownRelation(r) {
			return errf(CodeBadRelation, string(r))
		}
		if r == RelOutside {
			hasOutside = true
		}
	}
	if hasOutside && len(p.Sources) == 0 {
		return errf(CodeMissingField, "sources")
	}
	return nil
}

// ClassifyInput is the evidence needed to classify one topic.
type ClassifyInput struct {
	TopicID          string
	HasAction        bool // at least one identifiable, cited action/commitment
	HasDiscussion    bool // cited topic discussion without an identifiable action
	FullCoverage     bool // whole-source coverage was completed (required for absence)
	EvidenceReadable bool // evidence is present and readable (not ambiguous/missing)
	Profile          Profile
}

// Classify deterministically maps evidence and the jurisdiction profile to one
// of the five topic states. Priority: unresolved material evidence -> unclear;
// otherwise concrete action -> proposals; otherwise discussion -> mentioned;
// only then absence is evaluated against applicability. Omission is never
// opposition: absence in an applicable topic is not-covered, and not-applicable
// requires explicit whole-topic outside-scope evidence.
func Classify(in ClassifyInput) (State, error) {
	if _, ok := TopicByID(in.TopicID); !ok {
		return "", errf(CodeUnknownTopic, in.TopicID)
	}
	if err := ValidateProfile(in.Profile); err != nil {
		return "", err
	}
	if !in.EvidenceReadable {
		return StateUnclear, nil
	}
	if in.HasAction {
		return StateProposals, nil
	}
	if in.HasDiscussion {
		return StateMentioned, nil
	}
	if !in.FullCoverage {
		return StateUnclear, nil
	}
	switch in.Profile.Relations[in.TopicID] {
	case RelOutside:
		return StateNotApplicable, nil
	case RelUnknown, "":
		return StateUnclear, nil
	default: // direct, shared, delegated, advocacy-only establish relevance
		return StateNotCovered, nil
	}
}

// CompatibleComparison reports whether two programmes may be compared under the
// same taxonomy. It requires the exact same taxonomy version and the same scope;
// a national programme is never equated to its regional counterpart.
func CompatibleComparison(a, b Profile) error {
	if a.Version != b.Version {
		return errf(CodeVersionMismatch, a.Version+"/"+b.Version)
	}
	if a.Scope != b.Scope {
		return errf(CodeScopeMismatch, string(a.Scope)+"/"+string(b.Scope))
	}
	return nil
}

// Mapping is one old->new topic mapping in a migration.
type Mapping struct {
	OldTopicID string
	NewTopicID string
	Kind       string // "same", "split", or "merge"
}

// Migration records an explicit taxonomy release change with revalidation.
type Migration struct {
	FromRelease string
	ToRelease   string
	OldHash     string
	NewHash     string
	Mappings    []Mapping
}

// ValidateMigration checks that a migration is explicit and complete: distinct
// approved releases, differing content hashes (a new immutable release), and a
// mapping covering every old topic with no silent reuse of an old ID as a new
// one. Historical classifications are never rewritten.
func ValidateMigration(m Migration) error {
	// 1. Content must actually change: a migration is a new immutable release.
	if m.OldHash == m.NewHash {
		return errf(CodeBadHash, "unchanged")
	}
	// 2. Structure: explicit, complete, and no silent reuse of an ID.
	if len(m.Mappings) == 0 {
		return errf(CodeMigrationIncompl, "empty")
	}
	seenOld := map[string]bool{}
	seenNew := map[string]bool{}
	for _, mp := range m.Mappings {
		if _, ok := TopicByID(mp.OldTopicID); !ok {
			return errf(CodeUnknownTopic, mp.OldTopicID)
		}
		if mp.NewTopicID == "" {
			return errf(CodeMissingField, "new_topic_id")
		}
		if mp.OldTopicID == mp.NewTopicID && mp.Kind != "same" {
			return errf(CodeIDReuse, mp.OldTopicID)
		}
		if mp.Kind != "merge" {
			if seenNew[mp.NewTopicID] {
				return errf(CodeIDReuse, mp.NewTopicID)
			}
			seenNew[mp.NewTopicID] = true
		}
		if seenOld[mp.OldTopicID] {
			return errf(CodeIDReuse, mp.OldTopicID)
		}
		seenOld[mp.OldTopicID] = true
	}
	// Complete coverage: every v1 topic must be accounted for.
	if len(seenOld) != len(topics) {
		return errf(CodeMigrationIncompl, fmt.Sprintf("%d/%d", len(seenOld), len(topics)))
	}
	// 3. Both releases must be approved and distinct. In v1 only Release is
	//    known, so any migration fails closed here until a successor is approved.
	if err := ValidateRelease(m.FromRelease); err != nil {
		return errf(CodeUnknownVersion, m.FromRelease)
	}
	if err := ValidateRelease(m.ToRelease); err != nil {
		return errf(CodeUnknownVersion, m.ToRelease)
	}
	if m.FromRelease == m.ToRelease {
		return errf(CodeVersionMismatch, m.FromRelease)
	}
	return nil
}
