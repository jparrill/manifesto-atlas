// Package canonical implements the deterministic atlas.programme.local/v1
// validator for TQ-0037. It checks a local draft against its strict schema, the
// approved 21-topic taxonomy, the bound source/extraction revisions, the cited
// page/block/code-point excerpt evidence, and retained scope-evidence files.
//
// This is a structural/evidence check only: it can mark a draft structurally
// valid, never semantically validated. Only independent claim/coverage passes do
// that. All objects reject unknown fields; the draft cannot assert its own
// validation state. Identity is by local stable ordinals within the source
// revision, not by a hashed-identity protocol per measure/citation.
package canonical

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

// Version is the exact atlas.programme.local/v1 contract.
const Version = "atlas.programme.local/v1"

// Limits from the approved contract.
const (
	maxPayloadBytes  = 16 << 20 // 16 MiB canonical payload
	topicCount       = 21
	maxMeasures      = 5000
	maxCitations     = 20000
	maxGlobalSummary = 10
	maxTopicSummary  = 10
	maxStatement     = 2000
	maxConditions    = 20
	maxExcerpt       = 2000
	maxIDBytes       = 80
	maxDisplayField  = 500
	maxRights        = 2000
	maxCitePerStmt   = 32
)

// Error codes for the canonical validator.
const (
	CodeUnknownVersion  = "unknown-version"
	CodeSourceMismatch  = "source-mismatch"
	CodeMissingField    = "missing-field"
	CodeUnknownField    = "unknown-field"
	CodeBadID           = "bad-id"
	CodeDuplicateID     = "duplicate-id"
	CodeDanglingRef     = "dangling-reference"
	CodeBadPage         = "bad-page"
	CodeExcerptMismatch = "excerpt-mismatch"
	CodeBadState        = "bad-state"
	CodeStateConflict   = "state-conflict"
	CodeIncompleteBasis = "incomplete-basis"
	CodeMissingScope    = "missing-scope"
	CodeOverLimit       = "over-limit"
	CodeUnsafeText      = "unsafe-text"
	CodeBadCount        = "bad-count"
	CodeBadPath         = "bad-path"
	CodeMissingFile     = "missing-file"
	CodeHashMismatch    = "hash-mismatch"
	CodeSymlink         = "symlink"
)

// Error is a stable canonical error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "canonical: " + e.Code
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

// Party is the source-bound party identity.
type Party struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Geography is the source-bound territorial scope.
type Geography struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Election is a source-bound election reference, or null for general-political.
type Election struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Date string `json:"date"`
}

// Rights is the source-bound rights block. Statement is explicit null when the
// source declares no rights statement.
type Rights struct {
	Publisher string  `json:"publisher"`
	Statement *string `json:"statement"`
}

// Metadata is the source-derived block. Every field except DocumentTitle must
// equal the selected source record; DocumentTitle is a cited Statement or null.
type Metadata struct {
	Party           Party      `json:"party"`
	DocumentType    string     `json:"document_type"`
	Scope           string     `json:"scope"`
	Geography       Geography  `json:"geography"`
	Election        *Election  `json:"election"`
	PublicationDate *string    `json:"publication_date"`
	Rights          Rights     `json:"rights"`
	DocumentTitle   *Statement `json:"document_title"`
}

// Statement is a cited plain-text assertion.
type Statement struct {
	ID        string   `json:"id"`
	Text      string   `json:"text"`
	Citations []string `json:"citations"`
}

// Topic is one taxonomy section with its state and measures.
type Topic struct {
	ID                  string      `json:"id"`
	State               string      `json:"state"`
	Summary             []Statement `json:"summary"`
	MeasureIDs          []string    `json:"measure_ids"`
	SecondaryMeasureIDs []string    `json:"secondary_measure_ids"`
	Basis               Basis       `json:"basis"`
}

// Basis is the coverage/scope evidence for a topic state. PageRanges are sorted
// disjoint inclusive PDF intervals; ScopeEvidence is null or a retained
// source-evidence reference establishing explicit entire-topic document scope.
type Basis struct {
	PageRanges    []PageRange   `json:"page_ranges"`
	ScopeEvidence *RetainedFile `json:"scope_evidence"`
}

// PageRange is a one-based inclusive PDF page interval.
type PageRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// RetainedFile is a retained source-evidence reference by safe relative path and
// SHA-256.
type RetainedFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// Measure is one atomic proposal. Institution is null or a cited Statement
// naming the source-supported responsible/target institution; legal authority is
// never inferred. Citations is the ordered deduplicated union of the action,
// conditions, and institution Statement references.
type Measure struct {
	ID           string      `json:"id"`
	PrimaryTopic string      `json:"primary_topic"`
	Action       Statement   `json:"action"`
	Conditions   []Statement `json:"conditions"`
	Institution  *Statement  `json:"institution"`
	Citations    []string    `json:"citations"`
}

// Citation is one exact page/block/code-point excerpt, bound to the top-level
// source/extraction revision.
type Citation struct {
	ID      string `json:"id"`
	Page    int    `json:"page"`
	BlockID string `json:"block_id"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Excerpt string `json:"excerpt"`
}

// Programme is the atlas.programme.local/v1 draft under validation.
type Programme struct {
	Version                string      `json:"version"`
	ProgrammeID            string      `json:"programme_id"`
	SourceChecksum         string      `json:"source_checksum"`
	SourceMetadataChecksum string      `json:"source_metadata_checksum"`
	ExtractionChecksum     string      `json:"extraction_checksum"`
	TaxonomyVersion        string      `json:"taxonomy_version"`
	Metadata               Metadata    `json:"metadata"`
	Summary                []Statement `json:"summary"`
	Topics                 []Topic     `json:"topics"`
	Measures               []Measure   `json:"measures"`
	Citations              []Citation  `json:"citations"`
}

// SourceRecord is the selected local source the draft must bind to. The draft's
// metadata (all except document_title) must equal these fields; the checksums
// bind the draft to the exact selected source and its metadata bytes.
type SourceRecord struct {
	ProgrammeID            string
	SourceChecksum         string
	SourceMetadataChecksum string
	Party                  Party
	DocumentType           string
	Scope                  string
	Geography              Geography
	Election               *Election
	PublicationDate        *string
	Rights                 Rights
}

// Extraction is the validated page-aware extraction the citations bind to.
type Extraction struct {
	ExtractionChecksum string
	SourceChecksum     string
	PageCount          int
	BlockText          map[string]string
	BlockPage          map[string]int
}

// ValidateProgramme runs the deterministic structural/evidence check over a local
// draft bound to its selected source record and extraction. root is the directory
// holding retained scope-evidence files ("" when none are referenced). It returns
// the validated draft on success or a typed failure.
func ValidateProgramme(raw []byte, root string, src SourceRecord, ex Extraction) (*Programme, error) {
	canon, err := artifact.Canonicalize(raw)
	if err != nil {
		return nil, err
	}
	if len(canon) > maxPayloadBytes {
		return nil, errf(CodeOverLimit, "payload")
	}
	var p Programme
	dec := json.NewDecoder(bytes.NewReader(canon))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return nil, errf(CodeUnknownField, err.Error())
	}
	if err := validateProgramme(&p, root, src, ex); err != nil {
		return nil, err
	}
	return &p, nil
}

func validateProgramme(p *Programme, root string, src SourceRecord, ex Extraction) error {
	if p.Version != Version {
		return errf(CodeUnknownVersion, p.Version)
	}
	if err := taxonomy.ValidateRelease(p.TaxonomyVersion); err != nil {
		return errf(CodeUnknownVersion, p.TaxonomyVersion)
	}
	// Exact identity binding to the selected source and extraction revisions.
	if p.ProgrammeID != src.ProgrammeID {
		return errf(CodeSourceMismatch, "programme_id")
	}
	if p.SourceChecksum != src.SourceChecksum {
		return errf(CodeSourceMismatch, "source_checksum")
	}
	if p.SourceMetadataChecksum != src.SourceMetadataChecksum {
		return errf(CodeSourceMismatch, "source_metadata_checksum")
	}
	if p.ExtractionChecksum != ex.ExtractionChecksum {
		return errf(CodeSourceMismatch, "extraction_checksum")
	}
	// Required top-level arrays must be present (arrays, not missing/null).
	if p.Citations == nil {
		return errf(CodeMissingField, "citations")
	}
	if p.Summary == nil {
		return errf(CodeMissingField, "summary")
	}
	if p.Measures == nil {
		return errf(CodeMissingField, "measures")
	}
	// Citations first: they are the evidence every statement references.
	citeIDs, err := validateCitations(p.Citations, ex)
	if err != nil {
		return err
	}
	// Statements are numbered sequentially in canonical traversal order.
	if err := checkStatementOrder(p); err != nil {
		return err
	}
	// Metadata must equal the selected source record (except document_title).
	if err := validateMetadata(p.Metadata, src, citeIDs); err != nil {
		return err
	}
	// Measures: unique ordinals, cited, bounded.
	measureTopic, err := validateMeasures(p.Measures, citeIDs)
	if err != nil {
		return err
	}
	// Topics: exactly 21 in taxonomy order with per-state rules.
	if err := validateTopics(p.Topics, root, p.Measures, measureTopic, ex.PageCount, citeIDs); err != nil {
		return err
	}
	// Global summary: cited, bounded.
	if err := validateStatementList(p.Summary, citeIDs, maxGlobalSummary, "summary"); err != nil {
		return err
	}
	return nil
}

func validateMetadata(m Metadata, src SourceRecord, citeIDs map[string]bool) error {
	if m.Party.ID != src.Party.ID || m.Party.Name != src.Party.Name {
		return errf(CodeSourceMismatch, "party")
	}
	if m.DocumentType != src.DocumentType {
		return errf(CodeSourceMismatch, "document_type")
	}
	if m.Scope != src.Scope {
		return errf(CodeSourceMismatch, "scope")
	}
	if m.Geography.ID != src.Geography.ID || m.Geography.Name != src.Geography.Name {
		return errf(CodeSourceMismatch, "geography")
	}
	if !sameElection(m.Election, src.Election) {
		return errf(CodeSourceMismatch, "election")
	}
	if !sameNullable(m.PublicationDate, src.PublicationDate) {
		return errf(CodeSourceMismatch, "publication_date")
	}
	if m.Rights.Publisher != src.Rights.Publisher || !sameNullable(m.Rights.Statement, src.Rights.Statement) {
		return errf(CodeSourceMismatch, "rights")
	}
	if err := checkGeneratedText("party_name", m.Party.Name, maxDisplayField); err != nil {
		return err
	}
	if err := checkGeneratedText("geography_name", m.Geography.Name, maxDisplayField); err != nil {
		return err
	}
	if err := checkGeneratedText("rights_publisher", m.Rights.Publisher, maxDisplayField); err != nil {
		return err
	}
	if m.Rights.Statement != nil {
		if err := checkGeneratedText("rights_statement", *m.Rights.Statement, maxRights); err != nil {
			return err
		}
	}
	if m.DocumentTitle != nil {
		return validateStatement(*m.DocumentTitle, citeIDs)
	}
	return nil
}

func sameElection(a, b *Election) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.ID == b.ID && a.Type == b.Type && a.Date == b.Date
}

func sameNullable(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func validateCitations(cs []Citation, ex Extraction) (map[string]bool, error) {
	if len(cs) > maxCitations {
		return nil, errf(CodeOverLimit, "citations")
	}
	ids := make(map[string]bool, len(cs))
	for i := range cs {
		c := &cs[i]
		if c.ID != fmt.Sprintf("c-%d", i+1) {
			return nil, errf(CodeBadID, c.ID)
		}
		if c.Page < 1 || c.Page > ex.PageCount {
			return nil, errf(CodeBadPage, c.BlockID)
		}
		text, ok := ex.BlockText[c.BlockID]
		if !ok {
			return nil, errf(CodeDanglingRef, c.BlockID)
		}
		if ex.BlockPage[c.BlockID] != c.Page {
			return nil, errf(CodeBadPage, c.BlockID)
		}
		if c.Start < 0 || c.End < c.Start {
			return nil, errf(CodeExcerptMismatch, c.ID)
		}
		if utf8.RuneCountInString(c.Excerpt) < 1 || utf8.RuneCountInString(c.Excerpt) > maxExcerpt {
			return nil, errf(CodeOverLimit, "excerpt")
		}
		// Excerpt must equal the exact block slice [start:end] in code points.
		cp := []rune(text)
		if c.End > len(cp) {
			return nil, errf(CodeExcerptMismatch, c.ID)
		}
		if string(cp[c.Start:c.End]) != c.Excerpt {
			return nil, errf(CodeExcerptMismatch, c.ID)
		}
		ids[c.ID] = true
	}
	return ids, nil
}

func validateMeasures(ms []Measure, citeIDs map[string]bool) (map[string]string, error) {
	if len(ms) > maxMeasures {
		return nil, errf(CodeOverLimit, "measures")
	}
	primary := make(map[string]string, len(ms)) // measure id -> primary topic
	seen := make(map[string]bool, len(ms))
	for i := range ms {
		m := &ms[i]
		if m.ID != fmt.Sprintf("m-%d", i+1) {
			return nil, errf(CodeBadID, m.ID)
		}
		if seen[m.ID] {
			return nil, errf(CodeDuplicateID, m.ID)
		}
		seen[m.ID] = true
		if m.PrimaryTopic == "" {
			return nil, errf(CodeMissingField, "primary_topic")
		}
		if err := validateStatement(m.Action, citeIDs); err != nil {
			return nil, err
		}
		if len(m.Conditions) > maxConditions {
			return nil, errf(CodeOverLimit, "conditions")
		}
		// measure.citations == dedup(source-order) union of action+conditions+institution.
		var union []string
		used := make(map[string]bool)
		add := func(ids []string) {
			for _, id := range ids {
				if !used[id] {
					used[id] = true
					union = append(union, id)
				}
			}
		}
		add(m.Action.Citations)
		for i := range m.Conditions {
			if err := validateStatement(m.Conditions[i], citeIDs); err != nil {
				return nil, err
			}
			add(m.Conditions[i].Citations)
		}
		if m.Institution != nil {
			if err := validateStatement(*m.Institution, citeIDs); err != nil {
				return nil, err
			}
			add(m.Institution.Citations)
		}
		if !equalStrings(m.Citations, union) {
			return nil, errf(CodeDanglingRef, "measure_citations")
		}
		primary[m.ID] = m.PrimaryTopic
	}
	return primary, nil
}

func validateTopics(ts []Topic, root string, ms []Measure, measureTopic map[string]string, pageCount int, citeIDs map[string]bool) error {
	if len(ts) != topicCount {
		return errf(CodeBadCount, "topics")
	}
	wantIDs := make([]string, 0, topicCount)
	for _, t := range taxonomy.Topics() {
		wantIDs = append(wantIDs, t.ID)
	}
	measureSet := make(map[string]bool, len(ms))
	for i := range ms {
		measureSet[ms[i].ID] = true
	}
	claimed := make(map[string]bool) // primary measure claimed by a topic
	for i, t := range ts {
		if t.ID != wantIDs[i] {
			return errf(CodeBadCount, "topics")
		}
		if !validState(t.State) {
			return errf(CodeBadState, t.State)
		}
		if t.Summary == nil || t.MeasureIDs == nil || t.SecondaryMeasureIDs == nil || t.Basis.PageRanges == nil {
			return errf(CodeMissingField, "topic_"+t.ID)
		}
		if err := validatePageRanges(t.Basis.PageRanges, pageCount); err != nil {
			return err
		}
		if t.Basis.ScopeEvidence != nil {
			if err := verifyRetainedFile(root, *t.Basis.ScopeEvidence); err != nil {
				return err
			}
		}
		if err := validateStatementList(t.Summary, citeIDs, maxTopicSummary, "topic_summary"); err != nil {
			return err
		}
		// Primary measures: must exist and be owned by exactly this topic.
		for _, mid := range t.MeasureIDs {
			if !measureSet[mid] {
				return errf(CodeDanglingRef, mid)
			}
			if measureTopic[mid] != t.ID {
				return errf(CodeStateConflict, mid)
			}
			if claimed[mid] {
				return errf(CodeDuplicateID, mid)
			}
			claimed[mid] = true
		}
		// Secondary measures: must exist.
		for _, mid := range t.SecondaryMeasureIDs {
			if !measureSet[mid] {
				return errf(CodeDanglingRef, mid)
			}
		}
		if err := validateStateRules(t, pageCount); err != nil {
			return err
		}
	}
	// Every primary measure must be claimed by exactly one topic.
	if len(claimed) != len(ms) {
		return errf(CodeDanglingRef, "unassigned_measure")
	}
	return nil
}

// validateStateRules enforces the per-state content constraints.
func validateStateRules(t Topic, pageCount int) error {
	hasPrimary := len(t.MeasureIDs) > 0
	hasSecondary := len(t.SecondaryMeasureIDs) > 0
	hasSummary := len(t.Summary) > 0
	switch t.State {
	case "proposals":
		if !hasPrimary || !hasSummary {
			return errf(CodeStateConflict, t.ID)
		}
	case "mentioned":
		if hasPrimary || !hasSummary {
			return errf(CodeStateConflict, t.ID)
		}
	default:
		// not-covered, not-applicable, unclear: no asserted proposal or factual
		// summary, and no secondary references (absence state).
		if hasPrimary || hasSecondary || hasSummary {
			return errf(CodeStateConflict, t.ID)
		}
		// Absence states (not-covered, not-applicable) require complete source
		// coverage; unclear retains affected references but need not cover all.
		if t.State == "not-covered" || t.State == "not-applicable" {
			if !coversAll(t.Basis.PageRanges, pageCount) {
				return errf(CodeIncompleteBasis, t.ID)
			}
		}
		// not-applicable additionally requires explicit entire-topic scope evidence.
		if t.State == "not-applicable" && t.Basis.ScopeEvidence == nil {
			return errf(CodeMissingScope, t.ID)
		}
	}
	return nil
}

func validatePageRanges(ranges []PageRange, pageCount int) error {
	prevEnd := 0
	for _, r := range ranges {
		if r.Start < 1 || r.End < r.Start || r.End > pageCount {
			return errf(CodeBadPage, "range")
		}
		if r.Start <= prevEnd {
			return errf(CodeBadPage, "overlap")
		}
		prevEnd = r.End
	}
	return nil
}

func coversAll(ranges []PageRange, pageCount int) bool {
	if pageCount == 0 || len(ranges) == 0 {
		return false
	}
	if ranges[0].Start != 1 || ranges[len(ranges)-1].End != pageCount {
		return false
	}
	for i := 1; i < len(ranges); i++ {
		if ranges[i].Start != ranges[i-1].End+1 {
			return false
		}
	}
	return true
}

func validateStatementList(ss []Statement, citeIDs map[string]bool, max int, field string) error {
	if len(ss) > max {
		return errf(CodeOverLimit, field)
	}
	for i := range ss {
		if err := validateStatement(ss[i], citeIDs); err != nil {
			return err
		}
	}
	return nil
}

func validateStatement(s Statement, citeIDs map[string]bool) error {
	if err := checkStatementID(s.ID); err != nil {
		return err
	}
	if err := checkGeneratedText("text", s.Text, maxStatement); err != nil {
		return err
	}
	if len(s.Citations) < 1 || len(s.Citations) > maxCitePerStmt {
		return errf(CodeMissingField, "citations")
	}
	if citeIDs != nil {
		for _, id := range s.Citations {
			if !citeIDs[id] {
				return errf(CodeDanglingRef, id)
			}
		}
	}
	return nil
}

// checkStatementOrder enforces that every statement is numbered s-1, s-2, ... in
// the canonical traversal order: document title, global summary, each topic
// summary in topic order, then each measure in array order (its action,
// conditions, institution). This yields a complete gap-free sequence and rejects
// duplicate or out-of-order statement IDs.
func checkStatementOrder(p *Programme) error {
	n := 0
	visit := func(s Statement) error {
		n++
		if s.ID != fmt.Sprintf("s-%d", n) {
			return errf(CodeBadID, s.ID)
		}
		return nil
	}
	if p.Metadata.DocumentTitle != nil {
		if err := visit(*p.Metadata.DocumentTitle); err != nil {
			return err
		}
	}
	for i := range p.Summary {
		if err := visit(p.Summary[i]); err != nil {
			return err
		}
	}
	for ti := range p.Topics {
		for si := range p.Topics[ti].Summary {
			if err := visit(p.Topics[ti].Summary[si]); err != nil {
				return err
			}
		}
	}
	for mi := range p.Measures {
		m := &p.Measures[mi]
		if err := visit(m.Action); err != nil {
			return err
		}
		for ci := range m.Conditions {
			if err := visit(m.Conditions[ci]); err != nil {
				return err
			}
		}
		if m.Institution != nil {
			if err := visit(*m.Institution); err != nil {
				return err
			}
		}
	}
	if n != countStatements(p) {
		return errf(CodeBadID, "statement-order")
	}
	return nil
}

func countStatements(p *Programme) int {
	n := 0
	if p.Metadata.DocumentTitle != nil {
		n++
	}
	n += len(p.Summary)
	for i := range p.Topics {
		n += len(p.Topics[i].Summary)
	}
	for i := range p.Measures {
		n++ // action
		n += len(p.Measures[i].Conditions)
		if p.Measures[i].Institution != nil {
			n++
		}
	}
	return n
}

// checkStatementID enforces the "s-<ordinal>" local ASCII identifier format.
func checkStatementID(id string) error {
	const prefix = "s-"
	if len(id) <= len(prefix) || !strings.HasPrefix(id, prefix) || len(id) > len(prefix)+maxIDBytes {
		return errf(CodeBadID, id)
	}
	for i := len(prefix); i < len(id); i++ {
		if id[i] < '0' || id[i] > '9' {
			return errf(CodeBadID, id)
		}
	}
	return nil
}

// checkGeneratedText enforces the plain-data rule for generated Statement and
// metadata display fields: bounded length, no control characters except
// newline/tab, no raw markup delimiters, no Markdown link/image syntax, and no
// executable URL schemes. Quoted source excerpts are NOT subject to this rule.
func checkGeneratedText(field, s string, max int) error {
	n := utf8.RuneCountInString(s)
	if n < 1 || n > max {
		return errf(CodeOverLimit, field)
	}
	for _, r := range s {
		if (r < 0x20 && r != '\n' && r != '\t') || r == 0x7f {
			return errf(CodeUnsafeText, field)
		}
	}
	if strings.ContainsAny(s, "<>") {
		return errf(CodeUnsafeText, field)
	}
	if strings.Contains(s, "](") {
		return errf(CodeUnsafeText, field)
	}
	lower := strings.ToLower(s)
	for _, scheme := range []string{"javascript:", "vbscript:", "data:", "file:"} {
		if strings.Contains(lower, scheme) {
			return errf(CodeUnsafeText, field)
		}
	}
	return nil
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

func validState(s string) bool {
	switch s {
	case "proposals", "mentioned", "not-covered", "not-applicable", "unclear":
		return true
	}
	return false
}

// verifyRetainedFile checks that a retained scope-evidence reference resolves to a
// regular, non-symlink file under root whose bytes hash to the recorded SHA-256.
func verifyRetainedFile(root string, f RetainedFile) error {
	p, err := safeChild(root, f.Path)
	if err != nil {
		return errf(CodeBadPath, f.Path)
	}
	fi, err := os.Lstat(p)
	if err != nil {
		return errf(CodeMissingFile, f.Path)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return errf(CodeSymlink, f.Path)
	}
	if !fi.Mode().IsRegular() {
		return errf(CodeMissingFile, f.Path)
	}
	id, err := hashFile(p)
	if err != nil {
		return errf(CodeMissingFile, f.Path)
	}
	if id != f.SHA256 {
		return errf(CodeHashMismatch, f.Path)
	}
	return nil
}

func hashFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

func checkSafeRelPath(p string) error {
	if p == "" {
		return errf(CodeBadPath, "empty")
	}
	if filepath.IsAbs(p) || path.IsAbs(p) {
		return errf(CodeBadPath, "absolute")
	}
	clean := path.Clean(filepath.ToSlash(p))
	if clean == ".." || len(clean) > 3 && clean[:3] == "../" {
		return errf(CodeBadPath, "traversal")
	}
	return nil
}

func safeChild(root, rel string) (string, error) {
	if err := checkSafeRelPath(rel); err != nil {
		return "", err
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", errf(CodeBadPath, "root")
	}
	full := filepath.Join(rootAbs, filepath.FromSlash(path.Clean(filepath.ToSlash(rel))))
	fullAbs, err := filepath.Abs(full)
	if err != nil {
		return "", errf(CodeBadPath, "resolve")
	}
	if fullAbs != rootAbs && !strings.HasPrefix(fullAbs, rootAbs+string(os.PathSeparator)) {
		return "", errf(CodeBadPath, "escape")
	}
	return fullAbs, nil
}
