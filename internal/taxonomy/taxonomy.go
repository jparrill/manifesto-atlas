// Package taxonomy implements the atlas.taxonomy/v1 contract (release 1.0.0)
// for TQ-0036 under the docs/design/MVP.md local-first scope.
//
// It provides the immutable ordered 21-topic registry, the five mutually
// exclusive topic states with their Spanish public labels, release validation,
// and deterministic evidence-based document-scope classification. It contains
// no general legal-competence registry, comparison engine, or migration
// framework: those are deferred. All checks are offline and fail-closed; no
// model may add, remove, rename, or reorder topics.
package taxonomy

// Release is the single approved v1 taxonomy release.
const Release = "1.0.0"

// Error codes. Stable and free of private data.
const (
	CodeUnknownTopic   = "unknown-topic"
	CodeMissingField   = "missing-field"
	CodeBadRelation    = "bad-relation"
	CodeUnknownVersion = "unknown-version"
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

// Applicability is a document's explicit scope for one topic. It is established
// from the document itself, not from a general legal-competence registry.
type Applicability string

// Document-scope applicabilities.
const (
	// AppInScope: the document's domain addresses this topic's domain. Shared
	// powers and advocacy topics are in-scope, so their absence is not-covered,
	// never a blanket exclusion.
	AppInScope Applicability = "in-scope"
	// AppOutOfScope: the document is explicitly about a different domain; this
	// requires retained scope evidence.
	AppOutOfScope Applicability = "out-of-scope"
	// AppUnknown: applicability is ambiguous or not established. It is never
	// treated as a guessed exclusion.
	AppUnknown Applicability = ""
)

// ScopeEvidence is a retained reference establishing an explicit whole-topic
// document scope. It is required for a not-applicable decision.
type ScopeEvidence struct {
	Path   string
	SHA256 string
}

// ClassifyInput is the evidence needed to classify one topic.
type ClassifyInput struct {
	TopicID          string
	HasAction        bool // at least one identifiable, cited action/commitment
	HasDiscussion    bool // cited topic discussion without an identifiable action
	FullCoverage     bool // whole-source coverage was completed (required for absence)
	EvidenceReadable bool // evidence is present and readable (not ambiguous/missing)
	Applicability    Applicability
	// ScopeEvidence is required when Applicability == AppOutOfScope.
	ScopeEvidence *ScopeEvidence
}

// Classify deterministically maps evidence and the document's explicit scope to
// one of the five topic states. Priority: unresolved material evidence ->
// unclear; otherwise concrete action -> proposals; otherwise discussion ->
// mentioned; only then is absence evaluated against applicability. Omission is
// never opposition: absence in an in-scope topic is not-covered; not-applicable
// requires explicit retained whole-topic scope evidence; ambiguous
// applicability is unclear.
func Classify(in ClassifyInput) (State, error) {
	if _, ok := TopicByID(in.TopicID); !ok {
		return "", errf(CodeUnknownTopic, in.TopicID)
	}
	switch in.Applicability {
	case AppInScope, AppOutOfScope, AppUnknown:
	default:
		return "", errf(CodeBadRelation, string(in.Applicability))
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
	switch in.Applicability {
	case AppOutOfScope:
		if in.ScopeEvidence == nil || in.ScopeEvidence.Path == "" || in.ScopeEvidence.SHA256 == "" {
			return "", errf(CodeMissingField, "scope-evidence")
		}
		return StateNotApplicable, nil
	case AppInScope:
		return StateNotCovered, nil
	default: // AppUnknown
		return StateUnclear, nil
	}
}
