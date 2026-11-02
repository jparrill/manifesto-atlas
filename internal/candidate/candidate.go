// Package candidate implements the atlas.candidate/v1 recovery descriptor and
// the embedded atlas.source.local/v1 record established by docs/design/MVP.md
// (TQ-0028). It is an offline, fail-closed fixture check: it validates the exact
// field set, enums, hash/path/date formats, and verifies that the retained
// candidate and evidence bytes match their recorded SHA-256 digests and sizes.
// It performs no registration, no model call, and no network fetch.
package candidate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Versions of the two records handled by this package.
const (
	CandidateVersion = "atlas.candidate/v1"
	SourceVersion    = "atlas.source.local/v1"
)

// Error codes. Stable and free of private data.
const (
	CodeBadJSON      = "bad-json"
	CodeBadField     = "bad-field"
	CodeBadVersion   = "bad-version"
	CodeBadEnum      = "bad-enum"
	CodeBadHash      = "bad-hash"
	CodeBadPath      = "bad-path"
	CodeBadDate      = "bad-date"
	CodeBadURL       = "bad-url"
	CodeBadSize      = "bad-size"
	CodeHashMismatch = "hash-mismatch"
	CodeMissingFile  = "missing-file"
	CodeConflict     = "conflict"
)

// Error is a stable candidate error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "candidate: " + e.Code
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

// CodeOf returns the stable error code, or "" when err is not a *Error.
func CodeOf(err error) string {
	var e *Error
	if errorsAs(err, &e) {
		return e.Code
	}
	return ""
}

func errorsAs(err error, target **Error) bool {
	for err != nil {
		if e, ok := err.(*Error); ok {
			*target = e
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// Candidate is the atlas.candidate/v1 recovery descriptor.
type Candidate struct {
	Version       string
	ProgrammeID   string
	CandidatePath string
	Metadata      SourceRecord
}

// SourceRecord is the atlas.source.local/v1 record embedded in a candidate.
type SourceRecord struct {
	Version         string
	ProgrammeID     string
	Party           Party
	DocumentType    string
	Scope           string
	Geography       Geography
	Election        *Election
	PublicationDate *string
	ArchiveYear     int
	YearBasis       string
	Source          SourceRef
	Origin          Origin
	Rights          Rights
	Supersedes      *string
}

type Party struct {
	ID   string
	Name string
}

type Geography struct {
	ID   string
	Name string
}

type Election struct {
	ID   string
	Type string
	Date string
}

type SourceRef struct {
	Path             string
	SHA256           string
	Size             int
	OriginalFilename string
	OriginURL        string
	RetrievedAt      string
}

type Origin struct {
	Authority string
	Evidence  []EvidenceRef
}

type EvidenceRef struct {
	Path   string
	SHA256 string
}

type Rights struct {
	Publisher string
	Statement *string
}

var (
	slugRe    = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	hashRe    = regexp.MustCompile(`^[0-9a-f]{64}$`)
	isoDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// candidateKeys is the exact top-level field set of atlas.candidate/v1.
var candidateKeys = []string{"version", "programme_id", "candidate_path", "metadata"}

// sourceKeys is the exact field set of atlas.source.local/v1.
var sourceKeys = []string{
	"version", "programme_id", "party", "document_type", "scope", "geography",
	"election", "publication_date", "archive_year", "year_basis", "source",
	"origin", "rights", "supersedes",
}

// Parse validates the structure, enums and formats of a candidate descriptor.
// It performs no file I/O; use VerifyFiles to check the retained bytes.
func Parse(data []byte) (*Candidate, error) {
	top, err := decodeExact(data, candidateKeys)
	if err != nil {
		return nil, err
	}

	version, err := fieldString(top, "version")
	if err != nil {
		return nil, err
	}
	if version != CandidateVersion {
		return nil, errf(CodeBadVersion, version)
	}

	programmeID, err := fieldString(top, "programme_id")
	if err != nil {
		return nil, err
	}
	if !validSlug(programmeID) {
		return nil, errf(CodeBadField, "programme_id")
	}

	candidatePath, err := fieldString(top, "candidate_path")
	if err != nil {
		return nil, err
	}
	if err := checkRelativePath(candidatePath); err != nil {
		return nil, err
	}

	meta, err := parseSource(top["metadata"])
	if err != nil {
		return nil, err
	}
	if meta.ProgrammeID != programmeID {
		return nil, errf(CodeConflict, "programme_id")
	}

	return &Candidate{
		Version:       version,
		ProgrammeID:   programmeID,
		CandidatePath: candidatePath,
		Metadata:      *meta,
	}, nil
}

func parseSource(raw json.RawMessage) (*SourceRecord, error) {
	m, err := decodeExact(raw, sourceKeys)
	if err != nil {
		return nil, err
	}

	version, err := fieldString(m, "version")
	if err != nil {
		return nil, err
	}
	if version != SourceVersion {
		return nil, errf(CodeBadVersion, version)
	}

	programmeID, err := fieldString(m, "programme_id")
	if err != nil {
		return nil, err
	}
	if !validSlug(programmeID) {
		return nil, errf(CodeBadField, "programme_id")
	}

	party, err := parseParty(m["party"])
	if err != nil {
		return nil, err
	}

	documentType, err := fieldString(m, "document_type")
	if err != nil {
		return nil, err
	}
	if !validDocumentType(documentType) {
		return nil, errf(CodeBadEnum, "document_type")
	}

	scope, err := fieldString(m, "scope")
	if err != nil {
		return nil, err
	}
	if !validScope(scope) {
		return nil, errf(CodeBadEnum, "scope")
	}

	geography, err := parseGeography(m["geography"])
	if err != nil {
		return nil, err
	}

	election, err := parseElection(m["election"])
	if err != nil {
		return nil, err
	}
	if documentType == "electoral-programme" && election == nil {
		return nil, errf(CodeConflict, "election-required")
	}

	var publicationDate *string
	if !isNull(m["publication_date"]) {
		s, err := fieldString(m, "publication_date")
		if err != nil {
			return nil, err
		}
		if err := checkISODate(s); err != nil {
			return nil, err
		}
		publicationDate = &s
	}

	archiveYear, err := fieldInt(m, "archive_year")
	if err != nil {
		return nil, err
	}
	if archiveYear < 1970 || archiveYear > 2100 {
		return nil, errf(CodeBadField, "archive_year")
	}

	yearBasis, err := fieldString(m, "year_basis")
	if err != nil {
		return nil, err
	}
	if yearBasis != "publication" && yearBasis != "retrieval" {
		return nil, errf(CodeBadEnum, "year_basis")
	}
	if yearBasis == "publication" && publicationDate == nil {
		return nil, errf(CodeConflict, "year-basis")
	}

	source, err := parseSourceRef(m["source"])
	if err != nil {
		return nil, err
	}

	origin, err := parseOrigin(m["origin"])
	if err != nil {
		return nil, err
	}

	rights, err := parseRights(m["rights"])
	if err != nil {
		return nil, err
	}

	var supersedes *string
	if !isNull(m["supersedes"]) {
		s, err := fieldString(m, "supersedes")
		if err != nil {
			return nil, err
		}
		if !validHash(s) {
			return nil, errf(CodeBadHash, "supersedes")
		}
		supersedes = &s
	}

	return &SourceRecord{
		Version:         version,
		ProgrammeID:     programmeID,
		Party:           *party,
		DocumentType:    documentType,
		Scope:           scope,
		Geography:       *geography,
		Election:        election,
		PublicationDate: publicationDate,
		ArchiveYear:     archiveYear,
		YearBasis:       yearBasis,
		Source:          *source,
		Origin:          *origin,
		Rights:          *rights,
		Supersedes:      supersedes,
	}, nil
}

func parseParty(raw json.RawMessage) (*Party, error) {
	m, err := decodeExact(raw, []string{"id", "name"})
	if err != nil {
		return nil, err
	}
	id, err := fieldString(m, "id")
	if err != nil {
		return nil, err
	}
	if !validSlug(id) {
		return nil, errf(CodeBadField, "party.id")
	}
	name, err := fieldString(m, "name")
	if err != nil {
		return nil, err
	}
	if !nonempty(name) {
		return nil, errf(CodeBadField, "party.name")
	}
	return &Party{ID: id, Name: name}, nil
}

func parseGeography(raw json.RawMessage) (*Geography, error) {
	m, err := decodeExact(raw, []string{"id", "name"})
	if err != nil {
		return nil, err
	}
	id, err := fieldString(m, "id")
	if err != nil {
		return nil, err
	}
	if !nonempty(id) {
		return nil, errf(CodeBadField, "geography.id")
	}
	name, err := fieldString(m, "name")
	if err != nil {
		return nil, err
	}
	if !nonempty(name) {
		return nil, errf(CodeBadField, "geography.name")
	}
	return &Geography{ID: id, Name: name}, nil
}

func parseElection(raw json.RawMessage) (*Election, error) {
	if isNull(raw) {
		return nil, nil
	}
	m, err := decodeExact(raw, []string{"id", "type", "date"})
	if err != nil {
		return nil, err
	}
	id, err := fieldString(m, "id")
	if err != nil {
		return nil, err
	}
	typ, err := fieldString(m, "type")
	if err != nil {
		return nil, err
	}
	date, err := fieldString(m, "date")
	if err != nil {
		return nil, err
	}
	if !nonempty(id) || !nonempty(typ) {
		return nil, errf(CodeBadField, "election")
	}
	if err := checkISODate(date); err != nil {
		return nil, err
	}
	return &Election{ID: id, Type: typ, Date: date}, nil
}

func parseSourceRef(raw json.RawMessage) (*SourceRef, error) {
	m, err := decodeExact(raw, []string{
		"path", "sha256", "size", "original_filename", "origin_url", "retrieved_at",
	})
	if err != nil {
		return nil, err
	}
	p, err := fieldString(m, "path")
	if err != nil {
		return nil, err
	}
	if !nonempty(p) {
		return nil, errf(CodeBadField, "source.path")
	}
	sha, err := fieldString(m, "sha256")
	if err != nil {
		return nil, err
	}
	if !validHash(sha) {
		return nil, errf(CodeBadHash, "source.sha256")
	}
	size, err := fieldInt(m, "size")
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		return nil, errf(CodeBadSize, "source.size")
	}
	fn, err := fieldString(m, "original_filename")
	if err != nil {
		return nil, err
	}
	if !nonempty(fn) {
		return nil, errf(CodeBadField, "source.original_filename")
	}
	ou, err := fieldString(m, "origin_url")
	if err != nil {
		return nil, err
	}
	if err := checkURL(ou); err != nil {
		return nil, err
	}
	ra, err := fieldString(m, "retrieved_at")
	if err != nil {
		return nil, err
	}
	if err := checkRetrievedAt(ra); err != nil {
		return nil, err
	}
	return &SourceRef{
		Path:             p,
		SHA256:           sha,
		Size:             size,
		OriginalFilename: fn,
		OriginURL:        ou,
		RetrievedAt:      ra,
	}, nil
}

func parseOrigin(raw json.RawMessage) (*Origin, error) {
	m, err := decodeExact(raw, []string{"authority", "evidence"})
	if err != nil {
		return nil, err
	}
	authority, err := fieldString(m, "authority")
	if err != nil {
		return nil, err
	}
	if authority != "maintainer-selected" {
		return nil, errf(CodeBadEnum, "origin.authority")
	}
	var evRaw []json.RawMessage
	if err := json.Unmarshal(m["evidence"], &evRaw); err != nil || len(evRaw) == 0 {
		return nil, errf(CodeBadField, "origin.evidence")
	}
	evidence := make([]EvidenceRef, 0, len(evRaw))
	for _, er := range evRaw {
		em, err := decodeExact(er, []string{"path", "sha256"})
		if err != nil {
			return nil, err
		}
		ep, err := fieldString(em, "path")
		if err != nil {
			return nil, err
		}
		if err := checkRelativePath(ep); err != nil {
			return nil, err
		}
		esh, err := fieldString(em, "sha256")
		if err != nil {
			return nil, err
		}
		if !validHash(esh) {
			return nil, errf(CodeBadHash, "origin.evidence.sha256")
		}
		evidence = append(evidence, EvidenceRef{Path: ep, SHA256: esh})
	}
	return &Origin{Authority: authority, Evidence: evidence}, nil
}

func parseRights(raw json.RawMessage) (*Rights, error) {
	m, err := decodeExact(raw, []string{"publisher", "statement"})
	if err != nil {
		return nil, err
	}
	publisher, err := fieldString(m, "publisher")
	if err != nil {
		return nil, err
	}
	if !nonempty(publisher) {
		return nil, errf(CodeBadField, "rights.publisher")
	}
	var statement *string
	if !isNull(m["statement"]) {
		s, err := fieldString(m, "statement")
		if err != nil {
			return nil, err
		}
		statement = &s
	}
	return &Rights{Publisher: publisher, Statement: statement}, nil
}

// VerifyFiles checks that the retained candidate and evidence bytes under root
// match the recorded SHA-256 digests and that the candidate size matches. It
// rejects missing files, unsafe paths and byte mutations.
func VerifyFiles(root string, c *Candidate) error {
	candPath, err := safeChild(root, c.CandidatePath)
	if err != nil {
		return err
	}
	sha, size, err := hashFile(candPath)
	if err != nil {
		return errf(CodeMissingFile, c.CandidatePath)
	}
	if sha != c.Metadata.Source.SHA256 {
		return errf(CodeHashMismatch, "candidate")
	}
	if int64(c.Metadata.Source.Size) != size {
		return errf(CodeBadSize, "candidate")
	}
	for i, ev := range c.Metadata.Origin.Evidence {
		evPath, err := safeChild(root, ev.Path)
		if err != nil {
			return err
		}
		sha, _, err := hashFile(evPath)
		if err != nil {
			return errf(CodeMissingFile, fmt.Sprintf("evidence[%d]", i))
		}
		if sha != ev.SHA256 {
			return errf(CodeHashMismatch, fmt.Sprintf("evidence[%d]", i))
		}
	}
	return nil
}

// Check loads, parses and byte-verifies the candidate descriptor at root.
func Check(root string) (*Candidate, error) {
	data, err := os.ReadFile(filepath.Join(root, "source.json"))
	if err != nil {
		return nil, errf(CodeMissingFile, "source.json")
	}
	c, err := Parse(data)
	if err != nil {
		return nil, err
	}
	if err := VerifyFiles(root, c); err != nil {
		return nil, err
	}
	return c, nil
}

func decodeExact(data []byte, required []string) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errf(CodeBadJSON, "object")
	}
	if m == nil {
		return nil, errf(CodeBadJSON, "object")
	}
	got := make(map[string]bool, len(m))
	for k := range m {
		got[k] = true
	}
	for _, k := range required {
		if !got[k] {
			return nil, errf(CodeBadField, "missing:"+k)
		}
	}
	for k := range got {
		if !contains(required, k) {
			return nil, errf(CodeBadField, "unknown:"+k)
		}
	}
	return m, nil
}

func fieldString(m map[string]json.RawMessage, key string) (string, error) {
	var s string
	if err := json.Unmarshal(m[key], &s); err != nil {
		return "", errf(CodeBadField, key)
	}
	return s, nil
}

func fieldInt(m map[string]json.RawMessage, key string) (int, error) {
	var n int
	if err := json.Unmarshal(m[key], &n); err != nil {
		return 0, errf(CodeBadField, key)
	}
	return n, nil
}

func isNull(raw json.RawMessage) bool { return string(raw) == "null" }

func contains(keys []string, k string) bool {
	for _, x := range keys {
		if x == k {
			return true
		}
	}
	return false
}

func validSlug(s string) bool { return slugRe.MatchString(s) }

func nonempty(s string) bool { return strings.TrimSpace(s) != "" }

func validHash(s string) bool { return hashRe.MatchString(s) }

func validDocumentType(s string) bool {
	return s == "electoral-programme" || s == "general-political"
}

func validScope(s string) bool {
	switch s {
	case "european", "national", "autonomous", "provincial", "municipal", "general-political":
		return true
	}
	return false
}

func checkRelativePath(p string) error {
	if !nonempty(p) {
		return errf(CodeBadPath, "empty")
	}
	if strings.HasPrefix(p, "/") {
		return errf(CodeBadPath, "absolute")
	}
	clean := path.Clean(p)
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return errf(CodeBadPath, "traversal")
	}
	return nil
}

func checkISODate(s string) error {
	if !isoDateRe.MatchString(s) {
		return errf(CodeBadDate, s)
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return errf(CodeBadDate, s)
	}
	_ = t
	return nil
}

func checkRetrievedAt(s string) error {
	if _, err := time.Parse(time.RFC3339, s); err != nil {
		return errf(CodeBadDate, "retrieved_at")
	}
	return nil
}

func checkURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return errf(CodeBadURL, "parse")
	}
	if u.Scheme != "https" || u.Host == "" {
		return errf(CodeBadURL, "scheme")
	}
	return nil
}

// safeChild joins a declared relative path to root, rejecting absolute paths,
// traversal and any result that escapes the root directory.
func safeChild(root, rel string) (string, error) {
	if err := checkRelativePath(rel); err != nil {
		return "", err
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", errf(CodeBadPath, "root")
	}
	full := filepath.Join(rootAbs, filepath.FromSlash(path.Clean(rel)))
	fullAbs, err := filepath.Abs(full)
	if err != nil {
		return "", errf(CodeBadPath, "resolve")
	}
	if fullAbs != rootAbs && !strings.HasPrefix(fullAbs, rootAbs+string(os.PathSeparator)) {
		return "", errf(CodeBadPath, "escape")
	}
	return fullAbs, nil
}

func hashFile(p string) (string, int64, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}
