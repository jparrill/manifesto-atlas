// Package canonical implements the deterministic atlas.programme/v1 validator
// for TQ-0037. It checks a canonical draft against its strict schema, the
// approved taxonomy, the bound source/extraction/jurisdiction records, and the
// citation/page/excerpt evidence contract.
//
// This is a structural/evidence check only: it can mark a draft structurally
// valid, never semantically validated. Only independent claim/coverage passes
// and exact input-bound approvals do that. All objects reject unknown fields;
// the draft cannot assert its own validation state.
package canonical

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Version is the exact atlas.programme/v1 contract.
const Version = "atlas.programme/v1"

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
	CodeUnknownVersion     = "unknown-version"
	CodeSourceMismatch     = "source-mismatch"
	CodeMissingField       = "missing-field"
	CodeUnknownField       = "unknown-field"
	CodeBadID              = "bad-id"
	CodeDuplicateID        = "duplicate-id"
	CodeDanglingRef        = "dangling-reference"
	CodeBadPage            = "bad-page"
	CodeExcerptMismatch    = "excerpt-mismatch"
	CodeBadState           = "bad-state"
	CodeStateConflict      = "state-conflict"
	CodeIncompleteBasis    = "incomplete-basis"
	CodeMissingRule        = "missing-rule"
	CodeUnknownInstitution = "unknown-institution"
	CodeUnknownRelation    = "unknown-relation"
	CodeUnknownRule        = "unknown-rule"
	CodeOverLimit          = "over-limit"
	CodeUnsafeText         = "unsafe-text"
	CodeBadCount           = "bad-count"
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

// SourceRecord is the validated source identity the draft must bind to.
type SourceRecord struct {
	ManifestID      string
	ProgrammeID     string
	Checksum        string
	PartyID         string
	PartyName       string
	Scope           string
	GeographyID     string
	GeographyName   string
	DocumentType    string
	Election        *string
	PublicationDate *string
	Rights          string
}

// Extraction is the validated page-aware extraction the citations bind to.
type Extraction struct {
	ManifestID     string
	SourceChecksum string
	PageCount      int
	BlockText      map[string]string
	BlockPage      map[string]int
}

// Jurisdiction is the validated competence profile the measures bind to.
type Jurisdiction struct {
	ManifestID   string
	Unspecified  string
	Institutions map[string]bool
	Rules        map[string]bool
	Relations    map[string]bool
}

// Taxonomy is the approved taxonomy release the topics follow.
type Taxonomy struct {
	Release    string
	ManifestID string
	TopicIDs   []string
}

// Programme is the atlas.programme/v1 draft under validation.
type Programme struct {
	Version         string      `json:"version"`
	ProgrammeID     string      `json:"programme_id"`
	SourceRecord    string      `json:"source_record"`
	SourceChecksum  string      `json:"source_checksum"`
	Extraction      string      `json:"extraction"`
	TaxonomyVersion string      `json:"taxonomy_version"`
	Taxonomy        string      `json:"taxonomy"`
	Jurisdiction    string      `json:"jurisdiction"`
	Metadata        Metadata    `json:"metadata"`
	Summary         []Statement `json:"summary"`
	Topics          []Topic     `json:"topics"`
	Measures        []Measure   `json:"measures"`
	Citations       []Citation  `json:"citations"`
}

// Metadata is the source-derived party/election/scope/geography block.
type Metadata struct {
	PartyID         string            `json:"party_id"`
	PartyName       string            `json:"party_name"`
	DocumentTitle   *Statement        `json:"document_title"`
	DocumentType    string            `json:"document_type"`
	Scope           string            `json:"scope"`
	GeographyID     string            `json:"geography_id"`
	GeographyName   string            `json:"geography_name"`
	Election        *string           `json:"election"`
	PublicationDate *string           `json:"publication_date"`
	Rights          string            `json:"rights"`
	SourceEvidence  map[string]string `json:"source_evidence"`
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

// Basis is the coverage/applicability evidence for a topic state.
type Basis struct {
	PageRanges        []PageRange `json:"page_ranges"`
	ApplicabilityRule *string     `json:"applicability_rule"`
}

// PageRange is a one-based inclusive PDF page interval.
type PageRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Measure is one atomic proposal.
type Measure struct {
	ID               string      `json:"id"`
	PrimaryTopic     string      `json:"primary_topic"`
	Action           Statement   `json:"action"`
	Conditions       []Statement `json:"conditions"`
	InstitutionID    string      `json:"institution_id"`
	ActionRelation   string      `json:"action_relation"`
	JurisdictionRule *string     `json:"jurisdiction_rule"`
	Citations        []string    `json:"citations"`
}

// Citation is one exact page/block/offset excerpt.
type Citation struct {
	ID           string `json:"id"`
	SourceRecord string `json:"source_record"`
	Page         int    `json:"page"`
	BlockID      string `json:"block_id"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	Excerpt      string `json:"excerpt"`
}

// ValidateProgramme runs the deterministic structural/evidence check over a
// canonical draft bound to its source, extraction, jurisdiction, and taxonomy.
// It returns the validated draft on success or a typed failure.
func ValidateProgramme(raw []byte, src SourceRecord, ex Extraction, jur Jurisdiction, tax Taxonomy) (*Programme, error) {
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
	if err := validateProgramme(&p, src, ex, jur, tax); err != nil {
		return nil, err
	}
	return &p, nil
}

func validateProgramme(p *Programme, src SourceRecord, ex Extraction, jur Jurisdiction, tax Taxonomy) error {
	if p.Version != Version {
		return errf(CodeUnknownVersion, p.Version)
	}
	// Exact identity binding to the validated records.
	if p.ProgrammeID != src.ProgrammeID {
		return errf(CodeSourceMismatch, "programme_id")
	}
	if p.SourceRecord != src.ManifestID {
		return errf(CodeSourceMismatch, "source_record")
	}
	if p.SourceChecksum != src.Checksum {
		return errf(CodeSourceMismatch, "source_checksum")
	}
	if p.Extraction != ex.ManifestID {
		return errf(CodeSourceMismatch, "extraction")
	}
	if p.Taxonomy != tax.ManifestID {
		return errf(CodeSourceMismatch, "taxonomy")
	}
	if p.TaxonomyVersion != tax.Release || p.TaxonomyVersion != "1.0.0" {
		return errf(CodeUnknownVersion, p.TaxonomyVersion)
	}
	if p.Jurisdiction != jur.ManifestID {
		return errf(CodeSourceMismatch, "jurisdiction")
	}
	// Citations first: they are the evidence every statement references.
	citeIDs, err := validateCitations(p.Citations, src, ex)
	if err != nil {
		return err
	}
	// Statement IDs are globally unique across the whole document.
	usedStmt := make(map[string]bool)
	if err := validateMetadata(p.Metadata, src, citeIDs, usedStmt); err != nil {
		return err
	}
	// Measures: unique, evidenced, bound to jurisdiction.
	measureTopic, err := validateMeasures(p.Measures, citeIDs, jur, usedStmt)
	if err != nil {
		return err
	}
	// Topics: exactly 21 in taxonomy order with per-state rules.
	if err := validateTopics(p.Topics, p.Measures, measureTopic, tax, ex.PageCount, citeIDs, usedStmt); err != nil {
		return err
	}
	// Global summary: cited, bounded.
	if err := validateStatementList(p.Summary, citeIDs, maxGlobalSummary, "summary", usedStmt); err != nil {
		return err
	}
	return nil
}

func validateMetadata(m Metadata, src SourceRecord, citeIDs map[string]bool, usedStmt map[string]bool) error {
	if m.PartyID != src.PartyID {
		return errf(CodeSourceMismatch, "party_id")
	}
	if m.PartyName != src.PartyName {
		return errf(CodeSourceMismatch, "party_name")
	}
	if m.Scope != src.Scope {
		return errf(CodeSourceMismatch, "scope")
	}
	if m.GeographyID != src.GeographyID {
		return errf(CodeSourceMismatch, "geography_id")
	}
	if m.GeographyName != src.GeographyName {
		return errf(CodeSourceMismatch, "geography_name")
	}
	if m.DocumentType != src.DocumentType {
		return errf(CodeSourceMismatch, "document_type")
	}
	if m.Rights != src.Rights {
		return errf(CodeSourceMismatch, "rights")
	}
	if !sameNullable(m.Election, src.Election) {
		return errf(CodeSourceMismatch, "election")
	}
	if !sameNullable(m.PublicationDate, src.PublicationDate) {
		return errf(CodeSourceMismatch, "publication_date")
	}
	if err := checkGeneratedText("party_name", m.PartyName, maxDisplayField); err != nil {
		return err
	}
	if err := checkGeneratedText("geography_name", m.GeographyName, maxDisplayField); err != nil {
		return err
	}
	if err := checkGeneratedText("rights", m.Rights, maxRights); err != nil {
		return err
	}
	if len(m.SourceEvidence) == 0 {
		return errf(CodeMissingField, "source_evidence")
	}
	for k, v := range m.SourceEvidence {
		if k == "" || v == "" {
			return errf(CodeMissingField, "source_evidence")
		}
	}
	if m.DocumentTitle != nil {
		return validateStatement(*m.DocumentTitle, citeIDs, usedStmt)
	}
	return nil
}

func sameNullable(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func validateCitations(cs []Citation, src SourceRecord, ex Extraction) (map[string]bool, error) {
	if len(cs) > maxCitations {
		return nil, errf(CodeOverLimit, "citations")
	}
	ids := make(map[string]bool, len(cs))
	for i := range cs {
		c := &cs[i]
		if c.SourceRecord != src.ManifestID {
			return nil, errf(CodeSourceMismatch, "citations")
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
		if want := CitationID(c.SourceRecord, c.Page, c.BlockID, c.Start, c.End, c.Excerpt); c.ID != want {
			return nil, errf(CodeBadID, "citations")
		}
		if ids[c.ID] {
			return nil, errf(CodeDuplicateID, c.ID)
		}
		ids[c.ID] = true
	}
	return ids, nil
}

func validateMeasures(ms []Measure, citeIDs map[string]bool, jur Jurisdiction, usedStmt map[string]bool) (map[string]string, error) {
	if len(ms) > maxMeasures {
		return nil, errf(CodeOverLimit, "measures")
	}
	primary := make(map[string]string, len(ms)) // measure id -> primary topic
	seen := make(map[string]bool, len(ms))
	for i := range ms {
		m := &ms[i]
		if err := checkMeasureID(m.ID); err != nil {
			return nil, err
		}
		if seen[m.ID] {
			return nil, errf(CodeDuplicateID, m.ID)
		}
		seen[m.ID] = true
		if m.PrimaryTopic == "" {
			return nil, errf(CodeMissingField, "primary_topic")
		}
		if !jur.Relations[m.ActionRelation] {
			return nil, errf(CodeUnknownRelation, m.ActionRelation)
		}
		if m.InstitutionID == "" {
			return nil, errf(CodeMissingField, "institution_id")
		}
		if !jur.Institutions[m.InstitutionID] && m.InstitutionID != jur.Unspecified {
			return nil, errf(CodeUnknownInstitution, m.InstitutionID)
		}
		if m.JurisdictionRule != nil && !jur.Rules[*m.JurisdictionRule] {
			return nil, errf(CodeUnknownRule, *m.JurisdictionRule)
		}
		// Action statement is required and cited.
		if err := validateStatement(m.Action, citeIDs, usedStmt); err != nil {
			return nil, err
		}
		if len(m.Conditions) > maxConditions {
			return nil, errf(CodeOverLimit, "conditions")
		}
		// measure.citations == dedup(source-order) union of action+conditions.
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
		for _, c := range m.Conditions {
			if err := validateStatement(c, citeIDs, usedStmt); err != nil {
				return nil, err
			}
			add(c.Citations)
		}
		if !equalStrings(m.Citations, union) {
			return nil, errf(CodeDanglingRef, "measure_citations")
		}
		primary[m.ID] = m.PrimaryTopic
	}
	return primary, nil
}

func validateTopics(ts []Topic, ms []Measure, measureTopic map[string]string, tax Taxonomy, pageCount int, citeIDs map[string]bool, usedStmt map[string]bool) error {
	if len(ts) != topicCount {
		return errf(CodeBadCount, "topics")
	}
	measureSet := make(map[string]bool, len(ms))
	for i := range ms {
		measureSet[ms[i].ID] = true
	}
	claimed := make(map[string]bool) // primary measure claimed by a topic
	for i, t := range ts {
		if t.ID != tax.TopicIDs[i] {
			return errf(CodeBadCount, "topics")
		}
		if !validState(t.State) {
			return errf(CodeBadState, t.State)
		}
		if err := validatePageRanges(t.Basis.PageRanges, pageCount); err != nil {
			return err
		}
		if err := validateStatementList(t.Summary, citeIDs, maxTopicSummary, "topic_summary", usedStmt); err != nil {
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
		// not-covered, not-applicable, unclear: no asserted proposal or
		// factual summary, and no secondary references (absence state).
		if hasPrimary || hasSecondary || hasSummary {
			return errf(CodeStateConflict, t.ID)
		}
		if t.State == "not-applicable" {
			if t.Basis.ApplicabilityRule == nil {
				return errf(CodeMissingRule, t.ID)
			}
		}
		// Absence states (not-covered, not-applicable) require complete source
		// coverage; unclear retains affected references but need not cover all.
		if t.State == "not-covered" || t.State == "not-applicable" {
			if !coversAll(t.Basis.PageRanges, pageCount) {
				return errf(CodeIncompleteBasis, t.ID)
			}
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

func validateStatementList(ss []Statement, citeIDs map[string]bool, max int, field string, usedStmt map[string]bool) error {
	if len(ss) > max {
		return errf(CodeOverLimit, field)
	}
	for i := range ss {
		if err := validateStatement(ss[i], citeIDs, usedStmt); err != nil {
			return err
		}
	}
	return nil
}

func validateStatement(s Statement, citeIDs map[string]bool, usedStmt map[string]bool) error {
	if err := checkStatementID(s.ID); err != nil {
		return err
	}
	if usedStmt[s.ID] {
		return errf(CodeDuplicateID, s.ID)
	}
	usedStmt[s.ID] = true
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

// checkStatementID enforces the stable "s-<ordinal>" local ASCII identifier.
func checkStatementID(id string) error {
	if len(id) < 3 || len(id) > maxIDBytes || !strings.HasPrefix(id, "s-") {
		return errf(CodeBadID, "statement")
	}
	for i := 2; i < len(id); i++ {
		if id[i] < '0' || id[i] > '9' {
			return errf(CodeBadID, "statement")
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

func checkMeasureID(id string) error {
	const prefix = "m-"
	if len(id) != len(prefix)+64 || !strings.HasPrefix(id, prefix) {
		return errf(CodeBadID, "measures")
	}
	return isHex64(id[len(prefix):])
}

func isHex64(s string) error {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return errf(CodeBadID, "measures")
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

// CitationID is the stable content identity of an exact citation: "c-" plus the
// full SHA-256 hex of the canonical {source_record, page, block_id, start, end,
// excerpt} object. It never hashes prose alone.
func CitationID(sourceRecord string, page int, blockID string, start, end int, excerpt string) string {
	cv, err := artifact.CanonicalValue(map[string]any{
		"block_id":      blockID,
		"end":           end,
		"excerpt":       excerpt,
		"page":          page,
		"source_record": sourceRecord,
		"start":         start,
	})
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(cv)
	return "c-" + hex.EncodeToString(sum[:])
}
