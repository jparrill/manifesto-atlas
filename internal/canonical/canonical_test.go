package canonical_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

var blockTexts = map[string]string{
	"p1-b1": "Reforzaremos la democracia y construiremos un centro de salud.",
	"p2-b1": "El texto menciona las listas de espera y la sanidad.",
	"p3-b1": "Reorganizaremos el transporte publico regional.",
}

var blockPages = map[string]int{"p1-b1": 1, "p2-b1": 2, "p3-b1": 3}

func runeIndex(hay, needle []rune) int {
	if len(needle) > len(hay) {
		return -1
	}
	for i := 0; i+len(needle) <= len(hay); i++ {
		ok := true
		for j := 0; j < len(needle); j++ {
			if hay[i+j] != needle[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func mkCite(src string, page int, block, excerpt string) canonical.Citation {
	s := runeIndex([]rune(blockTexts[block]), []rune(excerpt))
	if s < 0 {
		panic("excerpt not in block: " + excerpt)
	}
	e := s + len([]rune(excerpt))
	return canonical.Citation{
		ID:           canonical.CitationID(src, page, block, s, e, excerpt),
		SourceRecord: src,
		Page:         page,
		BlockID:      block,
		Start:        s,
		End:          e,
		Excerpt:      excerpt,
	}
}

type builder struct {
	p   canonical.Programme
	src canonical.SourceRecord
	ex  canonical.Extraction
	jur canonical.Jurisdiction
	tax canonical.Taxonomy
	ord int
}

func (b *builder) stmt(text string, cites ...string) canonical.Statement {
	b.ord++
	return canonical.Statement{ID: fmt.Sprintf("s-%d", b.ord), Text: text, Citations: cites}
}

func sp(s string) *string { return &s }

func fullBasis() canonical.Basis {
	return canonical.Basis{PageRanges: []canonical.PageRange{{Start: 1, End: 3}}}
}

func newBuilder() *builder {
	b := &builder{}
	// Records.
	manID := "src-manifest"
	progID := "prog-1"
	checksum := "sha256:" + strings.Repeat("0", 64)
	b.src = canonical.SourceRecord{
		ManifestID: manID, ProgrammeID: progID, Checksum: checksum,
		PartyID: "party-1", PartyName: "Partido de Prueba",
		Scope: "autonomous", GeographyID: "es-cl", GeographyName: "Castilla y Leon",
		DocumentType: "electoral-programme",
		Election:     sp("fixture-regional"), PublicationDate: sp("2030-01-01"),
		Rights: "Todos los derechos reservados",
	}
	b.ex = canonical.Extraction{
		ManifestID: "ex-manifest", SourceChecksum: checksum, PageCount: 3,
		BlockText: blockTexts, BlockPage: blockPages,
	}
	b.jur = canonical.Jurisdiction{
		ManifestID: "jur-manifest", Unspecified: "inst-unspecified",
		Institutions: map[string]bool{"inst-health": true, "inst-transport": true},
		Rules:        map[string]bool{"rule-health": true, "rule-defence": true},
		Relations:    map[string]bool{"direct": true, "shared": true, "delegated": true, "advocacy": true, "unspecified": true},
	}
	taxIDs := make([]string, 0, 21)
	for _, t := range taxonomy.Topics() {
		taxIDs = append(taxIDs, t.ID)
	}
	b.tax = canonical.Taxonomy{Release: "1.0.0", ManifestID: "tax-manifest", TopicIDs: taxIDs}

	// Citations.
	c0 := mkCite(manID, 1, "p1-b1", "Reforzaremos la democracia")
	c1 := mkCite(manID, 1, "p1-b1", "construiremos un centro de salud")
	c2 := mkCite(manID, 2, "p2-b1", "El texto menciona las listas de espera")
	c3 := mkCite(manID, 3, "p3-b1", "Reorganizaremos el transporte publico")
	b.p.Citations = []canonical.Citation{c0, c1, c2, c3}

	// Statements in deterministic traversal order.
	sGlobal := b.stmt("El programa propone reforzar instituciones y sanidad.", c0.ID)
	sInst := b.stmt("El texto aborda la democracia.", c0.ID)
	sHealth := b.stmt("Propone construir un centro de salud.", c1.ID)
	sTrans := b.stmt("Propone reorganizar el transporte publico.", c3.ID)
	mAct1 := b.stmt("Construir un centro de salud.", c1.ID)
	mCond1 := b.stmt("Financiado con el presupuesto autonomico.", c1.ID)
	mAct2 := b.stmt("Reorganizar el transporte publico regional.", c3.ID)

	// Metadata.
	b.p.Metadata = canonical.Metadata{
		PartyID: b.src.PartyID, PartyName: b.src.PartyName,
		DocumentType: b.src.DocumentType, Scope: b.src.Scope,
		GeographyID: b.src.GeographyID, GeographyName: b.src.GeographyName,
		Election: b.src.Election, PublicationDate: b.src.PublicationDate,
		Rights: b.src.Rights,
		SourceEvidence: map[string]string{
			"party_name": "reg:party", "scope": "src:scope",
			"geography_id": "reg:geo", "rights": "src:rights",
		},
	}

	// Measures.
	m1ID := "m-" + strings.Repeat("a", 64)
	m2ID := "m-" + strings.Repeat("b", 64)
	b.p.Measures = []canonical.Measure{
		{ID: m1ID, PrimaryTopic: "healthcare", Action: mAct1, Conditions: []canonical.Statement{mCond1},
			InstitutionID: "inst-health", ActionRelation: "shared", JurisdictionRule: sp("rule-health"),
			Citations: []string{c1.ID}},
		{ID: m2ID, PrimaryTopic: "transport-infrastructure", Action: mAct2, Conditions: nil,
			InstitutionID: "inst-transport", ActionRelation: "direct",
			Citations: []string{c3.ID}},
	}

	// Topics: exactly 21 in taxonomy order.
	topicState := map[string]string{
		"institutions-democracy":   "mentioned",
		"healthcare":               "proposals",
		"transport-infrastructure": "proposals",
		"culture-language-sport":   "unclear",
		"defence-foreign-eu":       "not-applicable",
	}
	topicSummary := map[string][]canonical.Statement{
		"institutions-democracy":   {sInst},
		"healthcare":               {sHealth},
		"transport-infrastructure": {sTrans},
	}
	topicMeasures := map[string][]string{
		"healthcare":               {m1ID},
		"transport-infrastructure": {m2ID},
	}
	for _, id := range taxIDs {
		state := "not-covered"
		if s, ok := topicState[id]; ok {
			state = s
		}
		t := canonical.Topic{ID: id, State: state, Basis: fullBasis()}
		if st, ok := topicSummary[id]; ok {
			t.Summary = st
		}
		if mids, ok := topicMeasures[id]; ok {
			t.MeasureIDs = mids
		}
		if id == "culture-language-sport" {
			t.Basis = canonical.Basis{PageRanges: []canonical.PageRange{{Start: 3, End: 3}}}
		}
		if id == "defence-foreign-eu" {
			t.Basis.ApplicabilityRule = sp("rule-defence")
		}
		b.p.Topics = append(b.p.Topics, t)
	}

	// Top-level identity.
	b.p.Version = canonical.Version
	b.p.ProgrammeID = progID
	b.p.SourceRecord = manID
	b.p.SourceChecksum = checksum
	b.p.Extraction = "ex-manifest"
	b.p.TaxonomyVersion = "1.0.0"
	b.p.Taxonomy = "tax-manifest"
	b.p.Jurisdiction = "jur-manifest"
	b.p.Summary = []canonical.Statement{sGlobal}
	return b
}

func (b *builder) bytes() []byte {
	raw, err := json.Marshal(b.p)
	if err != nil {
		panic(err)
	}
	cv, err := artifact.Canonicalize(raw)
	if err != nil {
		panic(err)
	}
	return cv
}

func (b *builder) validate() (*canonical.Programme, error) {
	return canonical.ValidateProgramme(b.bytes(), b.src, b.ex, b.jur, b.tax)
}

func mustValidate(t *testing.T, b *builder) *canonical.Programme {
	t.Helper()
	p, err := b.validate()
	if err != nil {
		t.Fatalf("valid fixture rejected: %v", err)
	}
	return p
}

func TestValidProgramme(t *testing.T) {
	p := mustValidate(t, newBuilder())
	if p.Version != canonical.Version {
		t.Fatalf("version = %q", p.Version)
	}
	if len(p.Topics) != 21 {
		t.Fatalf("topics = %d", len(p.Topics))
	}
	// Every public factual field is cited; states are all exercised.
	seen := map[string]bool{}
	for _, tp := range p.Topics {
		seen[tp.State] = true
	}
	for _, want := range []string{"proposals", "mentioned", "not-covered", "not-applicable", "unclear"} {
		if !seen[want] {
			t.Errorf("state %s not exercised", want)
		}
	}
}

func TestDeterministicBytes(t *testing.T) {
	b := newBuilder()
	first := b.bytes()
	second := newBuilder().bytes()
	if !bytes.Equal(first, second) {
		t.Fatal("identical structures produced different bytes")
	}
	canon, err := artifact.Canonicalize(first)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	if !bytes.Equal(canon, first) {
		t.Fatal("canonicalization not idempotent")
	}
	// Key order independence: same object, different key order -> same bytes.
	a, err := artifact.Canonicalize([]byte(`{"b":1,"a":2}`))
	if err != nil {
		t.Fatalf("canonicalize a: %v", err)
	}
	c, err := artifact.Canonicalize([]byte(`{"a":2,"b":1}`))
	if err != nil {
		t.Fatalf("canonicalize c: %v", err)
	}
	if !bytes.Equal(a, c) {
		t.Fatal("key order affected canonical bytes")
	}
}

// mutate runs a mutation on a fresh builder and asserts the given code is returned.
func expectCode(t *testing.T, mutate func(*builder), want string) {
	t.Helper()
	b := newBuilder()
	mutate(b)
	_, err := b.validate()
	if err == nil {
		t.Fatalf("expected code %s, got nil", want)
	}
	ce, ok := err.(*canonical.Error)
	if !ok {
		t.Fatalf("expected canonical error, got %T: %v", err, err)
	}
	if ce.Code != want {
		t.Fatalf("code = %s, want %s", ce.Code, want)
	}
}

func TestUnknownVersion(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Version = "atlas.programme/v2" }, canonical.CodeUnknownVersion)
}

func TestSourceMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.SourceChecksum = "sha256:" + strings.Repeat("f", 64) }, canonical.CodeSourceMismatch)
}

func TestMissingCitation(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Summary[0].Citations = nil // uncited global summary
	}, canonical.CodeMissingField)
}

func TestInvalidPage(t *testing.T) {
	expectCode(t, func(b *builder) {
		c := b.p.Citations[0]
		c.Page = 99
		c.ID = canonical.CitationID(c.SourceRecord, c.Page, c.BlockID, c.Start, c.End, c.Excerpt)
		b.p.Citations[0] = c
	}, canonical.CodeBadPage)
}

func TestExcerptMismatch(t *testing.T) {
	expectCode(t, func(b *builder) {
		c := b.p.Citations[0]
		c.Excerpt = "completely different text"
		c.ID = canonical.CitationID(c.SourceRecord, c.Page, c.BlockID, c.Start, c.End, c.Excerpt)
		b.p.Citations[0] = c
	}, canonical.CodeExcerptMismatch)
}

func TestDuplicateMeasure(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Measures[1].ID = b.p.Measures[0].ID
	}, canonical.CodeDuplicateID)
}

func TestUnknownField(t *testing.T) {
	b := newBuilder()
	raw := b.bytes()
	idx := bytes.LastIndexByte(raw, '}')
	mut := append([]byte{}, raw[:idx]...)
	mut = append(mut, []byte(`,"extra":1`)...)
	mut = append(mut, raw[idx:]...)
	_, err := canonical.ValidateProgramme(mut, b.src, b.ex, b.jur, b.tax)
	if err == nil {
		t.Fatal("unknown field accepted")
	}
	ce, ok := err.(*canonical.Error)
	if !ok || ce.Code != canonical.CodeUnknownField {
		t.Fatalf("code = %v, want unknown-field", err)
	}
}

func TestDuplicateStatement(t *testing.T) {
	expectCode(t, func(b *builder) {
		// Reuse a statement ID in a topic summary.
		b.p.Topics[6].Summary = append(b.p.Topics[6].Summary, b.p.Topics[6].Summary[0])
	}, canonical.CodeDuplicateID)
}

func TestStateConflict(t *testing.T) {
	// proposals without primary measures.
	expectCode(t, func(b *builder) {
		b.p.Topics[6].MeasureIDs = nil
	}, canonical.CodeStateConflict)
}

func TestNotCoveredIncompleteBasis(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[1].Basis.PageRanges = []canonical.PageRange{{Start: 1, End: 2}}
	}, canonical.CodeIncompleteBasis)
}

func TestNotApplicableMissingRule(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[20].Basis.ApplicabilityRule = nil
	}, canonical.CodeMissingRule)
}

func TestUnknownInstitution(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Measures[0].InstitutionID = "inst-bogus"
	}, canonical.CodeUnknownInstitution)
}

func TestUnknownRelation(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Measures[0].ActionRelation = "dictatorial"
	}, canonical.CodeUnknownRelation)
}

func TestUnsafeText(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Summary[0].Text = "Propone <b>algo</b>"
	}, canonical.CodeUnsafeText)
}

func TestNationalScope(t *testing.T) {
	b := newBuilder()
	b.src.Scope = "national"
	b.src.GeographyID = "es"
	b.src.GeographyName = "España"
	b.p.Metadata.Scope = "national"
	b.p.Metadata.GeographyID = "es"
	b.p.Metadata.GeographyName = "España"
	if _, err := b.validate(); err != nil {
		t.Fatalf("national scope rejected: %v", err)
	}
}

func TestMeasureCitationsUnion(t *testing.T) {
	// Wrong measure.citations (missing a condition citation) must fail.
	expectCode(t, func(b *builder) {
		b.p.Measures[0].Citations = nil
	}, canonical.CodeDanglingRef)
}
