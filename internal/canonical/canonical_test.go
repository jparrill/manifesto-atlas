package canonical_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func mkCite(page int, block, excerpt string) canonical.Citation {
	s := runeIndex([]rune(blockTexts[block]), []rune(excerpt))
	if s < 0 {
		panic("excerpt not in block: " + excerpt)
	}
	return canonical.Citation{Page: page, BlockID: block, Start: s, End: s + len([]rune(excerpt)), Excerpt: excerpt}
}

func sp(s string) *string { return &s }

type builder struct {
	p             canonical.Programme
	src           canonical.SourceRecord
	ex            canonical.Extraction
	root          string
	scopeEvidence canonical.RetainedFile
	ord           int
}

// stmt assigns the next sequential local statement ordinal in traversal order.
func (b *builder) stmt(text string, cites ...string) canonical.Statement {
	b.ord++
	return canonical.Statement{ID: fmt.Sprintf("s-%d", b.ord), Text: text, Citations: cites}
}

func fullBasis() canonical.Basis {
	return canonical.Basis{PageRanges: []canonical.PageRange{{Start: 1, End: 3}}}
}

func newBuilder(t *testing.T) *builder {
	t.Helper()
	b := &builder{}

	// Retained scope-evidence file for the not-applicable topic.
	root := t.TempDir()
	evPath := "evidence/scope-defence.txt"
	evContent := []byte("Alcance: documento autonómico; la defensa es competencia estatal.")
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, evPath)), 0o755); err != nil {
		t.Fatalf("mkdir evidence: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, evPath), evContent, 0o644); err != nil {
		t.Fatalf("write evidence: %v", err)
	}
	b.root = root
	evSum := sha256.Sum256(evContent)
	b.scopeEvidence = canonical.RetainedFile{Path: evPath, SHA256: "sha256:" + hex.EncodeToString(evSum[:])}

	// Selected source + extraction records.
	progID := "frente-obrero-programa-politico"
	checksum := "sha256:" + strings.Repeat("0", 64)
	metaChecksum := "sha256:" + strings.Repeat("1", 64)
	exChecksum := "sha256:" + strings.Repeat("2", 64)
	rightsStmt := "Todos los derechos reservados"
	b.src = canonical.SourceRecord{
		ProgrammeID: progID, SourceChecksum: checksum, SourceMetadataChecksum: metaChecksum,
		Party:        canonical.Party{ID: "party-1", Name: "Partido de Prueba"},
		DocumentType: "electoral-programme", Scope: "autonomous",
		Geography:       canonical.Geography{ID: "es-cl", Name: "Castilla y Leon"},
		Election:        &canonical.Election{ID: "fixture-regional", Type: "autonomous", Date: "2030-05-01"},
		PublicationDate: sp("2030-01-01"),
		Rights:          canonical.Rights{Publisher: "Partido de Prueba", Statement: &rightsStmt},
	}
	b.ex = canonical.Extraction{
		ExtractionChecksum: exChecksum, SourceChecksum: checksum, PageCount: 3,
		BlockText: blockTexts, BlockPage: blockPages,
	}

	// Citations, array order c-1..c-4.
	b.p.Citations = []canonical.Citation{
		{ID: "c-1", Page: 1, BlockID: "p1-b1", Start: 0, End: 0, Excerpt: "Reforzaremos la democracia"},
		{ID: "c-2", Page: 1, BlockID: "p1-b1", Start: 0, End: 0, Excerpt: "construiremos un centro de salud"},
		{ID: "c-3", Page: 2, BlockID: "p2-b1", Start: 0, End: 0, Excerpt: "El texto menciona las listas de espera"},
		{ID: "c-4", Page: 3, BlockID: "p3-b1", Start: 0, End: 0, Excerpt: "Reorganizaremos el transporte publico"},
	}
	for i := range b.p.Citations {
		c := mkCite(b.p.Citations[i].Page, b.p.Citations[i].BlockID, b.p.Citations[i].Excerpt)
		b.p.Citations[i].Start = c.Start
		b.p.Citations[i].End = c.End
	}

	// Statements in canonical traversal order: document title, global summary,
	// each topic summary in topic order, then each measure in array order.
	sTitle := b.stmt("Programa politico de prueba", "c-1")
	sGlobal := b.stmt("El programa propone reforzar instituciones y sanidad.", "c-1")
	sInst := b.stmt("El texto aborda la democracia.", "c-1")
	sHealth := b.stmt("Propone construir un centro de salud.", "c-2")
	sTrans := b.stmt("Propone reorganizar el transporte publico.", "c-4")
	m1Act := b.stmt("Construir un centro de salud.", "c-2")
	m1Cond := b.stmt("Financiado con el presupuesto autonomico.", "c-2")
	m1Inst := b.stmt("Gestionado por la sanidad publica.", "c-2")
	m2Act := b.stmt("Reorganizar el transporte publico regional.", "c-4")

	// Metadata equals the selected source record, plus a cited document title.
	b.p.Metadata = canonical.Metadata{
		Party: b.src.Party, DocumentType: b.src.DocumentType, Scope: b.src.Scope,
		Geography: b.src.Geography, Election: b.src.Election,
		PublicationDate: b.src.PublicationDate, Rights: b.src.Rights,
		DocumentTitle: &sTitle,
	}

	// Measures, array order m-1, m-2.
	b.p.Measures = []canonical.Measure{
		{ID: "m-1", PrimaryTopic: "healthcare", Action: m1Act,
			Conditions: []canonical.Statement{m1Cond}, Institution: &m1Inst, Citations: []string{"c-2"}},
		{ID: "m-2", PrimaryTopic: "transport-infrastructure", Action: m2Act,
			Conditions: []canonical.Statement{}, Institution: nil, Citations: []string{"c-4"}},
	}

	// Topics: exactly 21 in taxonomy order.
	topicSummary := map[string][]canonical.Statement{
		"institutions-democracy":   {sInst},
		"healthcare":               {sHealth},
		"transport-infrastructure": {sTrans},
	}
	topicMeasures := map[string][]string{
		"healthcare":               {"m-1"},
		"transport-infrastructure": {"m-2"},
	}
	stateOverride := map[string]string{
		"institutions-democracy":   "mentioned",
		"healthcare":               "proposals",
		"transport-infrastructure": "proposals",
		"culture-language-sport":   "unclear",
		"defence-foreign-eu":       "not-applicable",
	}
	for _, tp := range taxonomy.Topics() {
		state := "not-covered"
		if s, ok := stateOverride[tp.ID]; ok {
			state = s
		}
		basis := fullBasis()
		if tp.ID == "culture-language-sport" {
			basis.PageRanges = []canonical.PageRange{{Start: 3, End: 3}}
		}
		if tp.ID == "defence-foreign-eu" {
			basis.ScopeEvidence = &b.scopeEvidence
		}
		topic := canonical.Topic{
			ID: tp.ID, State: state,
			Summary: []canonical.Statement{}, MeasureIDs: []string{}, SecondaryMeasureIDs: []string{},
			Basis: basis,
		}
		if st, ok := topicSummary[tp.ID]; ok {
			topic.Summary = st
		}
		if mids, ok := topicMeasures[tp.ID]; ok {
			topic.MeasureIDs = mids
		}
		b.p.Topics = append(b.p.Topics, topic)
	}

	// Top-level identity binding.
	b.p.Version = canonical.Version
	b.p.ProgrammeID = progID
	b.p.SourceChecksum = checksum
	b.p.SourceMetadataChecksum = metaChecksum
	b.p.ExtractionChecksum = exChecksum
	b.p.TaxonomyVersion = taxonomy.Release
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
	return canonical.ValidateProgramme(b.bytes(), b.root, b.src, b.ex)
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
	p := mustValidate(t, newBuilder(t))
	if p.Version != canonical.Version {
		t.Fatalf("version = %q", p.Version)
	}
	if len(p.Topics) != 21 {
		t.Fatalf("topics = %d", len(p.Topics))
	}
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
	first := newBuilder(t).bytes()
	second := newBuilder(t).bytes()
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

func expectCode(t *testing.T, mutate func(*builder), want string) {
	t.Helper()
	b := newBuilder(t)
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
	expectCode(t, func(b *builder) { b.p.Version = "atlas.programme.local/v2" }, canonical.CodeUnknownVersion)
}

func TestSourceMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.SourceChecksum = "sha256:" + strings.Repeat("f", 64) }, canonical.CodeSourceMismatch)
}

func TestSourceMetadataChecksumMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.SourceMetadataChecksum = "sha256:" + strings.Repeat("f", 64) }, canonical.CodeSourceMismatch)
}

func TestExtractionChecksumMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.ExtractionChecksum = "sha256:" + strings.Repeat("f", 64) }, canonical.CodeSourceMismatch)
}

func TestMetadataMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Metadata.Party.Name = "Otro partido" }, canonical.CodeSourceMismatch)
}

func TestMissingCitation(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Summary[0].Citations = nil }, canonical.CodeMissingField)
}

func TestInvalidPage(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Citations[0].Page = 99 }, canonical.CodeBadPage)
}

func TestExcerptMismatch(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Citations[0].Excerpt = "completely different text" }, canonical.CodeExcerptMismatch)
}

func TestBadCitationOrdinal(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Citations[1].ID = "c-1" }, canonical.CodeBadID)
}

func TestDuplicateMeasure(t *testing.T) {
	expectCode(t, func(b *builder) { b.p.Measures[1].ID = b.p.Measures[0].ID }, canonical.CodeBadID)
}

func TestUnknownField(t *testing.T) {
	b := newBuilder(t)
	raw := b.bytes()
	idx := bytes.LastIndexByte(raw, '}')
	mut := append([]byte{}, raw[:idx]...)
	mut = append(mut, []byte(`,"extra":1`)...)
	mut = append(mut, raw[idx:]...)
	_, err := canonical.ValidateProgramme(mut, b.root, b.src, b.ex)
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
		b.p.Topics[6].Summary = append(b.p.Topics[6].Summary, b.p.Topics[6].Summary[0])
	}, canonical.CodeBadID)
}

func TestStateConflict(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[6].MeasureIDs = []string{}
	}, canonical.CodeStateConflict)
}

func TestNotCoveredIncompleteBasis(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[1].Basis.PageRanges = []canonical.PageRange{{Start: 1, End: 2}}
	}, canonical.CodeIncompleteBasis)
}

func TestNotApplicableMissingScope(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[20].Basis.ScopeEvidence = nil
	}, canonical.CodeMissingScope)
}

func TestScopeEvidenceHashMismatch(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[20].Basis.ScopeEvidence = &canonical.RetainedFile{Path: b.scopeEvidence.Path, SHA256: "sha256:" + strings.Repeat("f", 64)}
	}, canonical.CodeHashMismatch)
}

func TestScopeEvidenceTraversal(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Topics[20].Basis.ScopeEvidence = &canonical.RetainedFile{Path: "../escape.txt", SHA256: b.scopeEvidence.SHA256}
	}, canonical.CodeBadPath)
}

func TestInstitutionCitationUnion(t *testing.T) {
	expectCode(t, func(b *builder) {
		// Institution cites c-2; dropping it from the measure union must fail.
		b.p.Measures[0].Citations = nil
	}, canonical.CodeDanglingRef)
}

func TestUnsafeText(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Summary[0].Text = "Propone <b>algo</b>"
	}, canonical.CodeUnsafeText)
}

func TestNationalScope(t *testing.T) {
	b := newBuilder(t)
	b.src.Scope = "national"
	b.src.Geography = canonical.Geography{ID: "es", Name: "España"}
	b.p.Metadata.Scope = "national"
	b.p.Metadata.Geography = canonical.Geography{ID: "es", Name: "España"}
	if _, err := b.validate(); err != nil {
		t.Fatalf("national scope rejected: %v", err)
	}
}

func TestGeneralPoliticalNullElection(t *testing.T) {
	b := newBuilder(t)
	b.src.Scope = "general-political"
	b.src.DocumentType = "general-political"
	b.src.Election = nil
	b.p.Metadata.Scope = "general-political"
	b.p.Metadata.DocumentType = "general-political"
	b.p.Metadata.Election = nil
	if _, err := b.validate(); err != nil {
		t.Fatalf("general-political null election rejected: %v", err)
	}
}

func TestMeasureCitationsUnion(t *testing.T) {
	expectCode(t, func(b *builder) {
		b.p.Measures[1].Citations = nil
	}, canonical.CodeDanglingRef)
}

func TestWrongTopicOrder(t *testing.T) {
	// The 21-topic release is fixed and ordered; a topic ID outside the approved
	// order rejects.
	expectCode(t, func(b *builder) {
		b.p.Topics[3].ID = "bogus-topic"
	}, canonical.CodeBadCount)
}

func TestUnicodeCodePointOffsets(t *testing.T) {
	// Excerpts are exact Unicode code-point slices, not byte slices.
	b := newBuilder(t)
	text := "añade émojis 😀 aquí"
	b.ex.BlockText["p1-b2"] = text
	b.ex.BlockPage["p1-b2"] = 1
	excerpt := "émojis"
	r := []rune(text)
	s := runeIndex(r, []rune(excerpt))
	b.p.Citations[0] = canonical.Citation{
		ID: "c-1", Page: 1, BlockID: "p1-b2", Start: s, End: s + len([]rune(excerpt)), Excerpt: excerpt,
	}
	if _, err := b.validate(); err != nil {
		t.Fatalf("multi-byte code-point citation rejected: %v", err)
	}
}
