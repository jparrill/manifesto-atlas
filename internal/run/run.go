// Package run implements the atlas.selection/v1 owner-controlled selection
// record and the atlas.run/v1 per-programme run record established by
// docs/design/MVP.md (TQ-0063).
//
// A selection record is a bounded, unique set of owner-selected source
// checksums; a model, PDF or contributor path cannot activate it. A run record
// binds the actual current input, report and output files of one programme by
// rooted local path and SHA-256. Verification is offline and fail-closed: it
// never performs a network lookup, follows a symlink, or treats a stored
// boolean as authority. Deterministic serialization reuses the strict
// canonical-JSON helper.
package run

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Versions.
const (
	SelectionVersion = "atlas.selection/v1"
	RunVersion       = "atlas.run/v1"
)

// Bounds.
const (
	maxInput   = 1 << 20
	maxSources = 1000
	maxRefs    = 1000
)

// Error codes. Stable and free of private data.
const (
	CodeBadJSON      = "bad-json"
	CodeBadField     = "bad-field"
	CodeBadVersion   = "bad-version"
	CodeBadHash      = "bad-hash"
	CodeBadPath      = "bad-path"
	CodeDuplicate    = "duplicate"
	CodeTooLarge     = "too-large"
	CodeHashMismatch = "hash-mismatch"
	CodeMissingFile  = "missing-file"
	CodeSymlink      = "symlink"
)

// Error is a stable run error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "run: " + e.Code
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

// Selection is the atlas.selection/v1 owner-controlled selection record.
type Selection struct {
	Version string
	Sources []SourceSelection
}

// SourceSelection is one owner-selected source binding.
type SourceSelection struct {
	ProgrammeID       string
	SourceChecksum    string
	MetadataChecksum  string
	EvidenceChecksums []string
}

// Ref is a bounded rooted local file reference.
type Ref struct {
	Path   string
	SHA256 string
}

// Run is the atlas.run/v1 per-programme run record.
type Run struct {
	Version        string
	ProgrammeID    string
	SourceChecksum string
	Inputs         map[string]*Ref
	Reports        map[string]*Ref
	Outputs        map[string]*Ref
}

// Exact reference key sets.
var (
	inputKeys  = []string{"source-metadata", "selection", "extraction", "taxonomy", "programme-schema", "claim-schema", "coverage-schema", "draft", "normalizer-prompt", "claim-prompt", "coverage-prompt", "model-configuration", "tool-configuration"}
	reportKeys = []string{"claims", "coverage", "site"}
	outputKeys = []string{"markdown", "html"}
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func validSlug(s string) bool { return slugRe.MatchString(s) }
func nonempty(s string) bool  { return s != "" }

// ParseSelection validates the strict atlas.selection/v1 record.
func ParseSelection(data []byte) (*Selection, error) {
	if len(data) > maxInput {
		return nil, errf(CodeTooLarge, "input")
	}
	m, err := decodeExact(data, []string{"version", "sources"})
	if err != nil {
		return nil, err
	}
	version, err := fieldString(m, "version")
	if err != nil {
		return nil, err
	}
	if version != SelectionVersion {
		return nil, errf(CodeBadVersion, version)
	}
	var srcRaw []json.RawMessage
	if err := json.Unmarshal(m["sources"], &srcRaw); err != nil {
		return nil, errf(CodeBadJSON, "sources")
	}
	if len(srcRaw) > maxSources {
		return nil, errf(CodeTooLarge, "sources")
	}
	seen := make(map[string]bool, len(srcRaw))
	sources := make([]SourceSelection, 0, len(srcRaw))
	for _, sr := range srcRaw {
		sm, err := decodeExact(sr, []string{"programme_id", "source_checksum", "metadata_checksum", "evidence_checksums"})
		if err != nil {
			return nil, err
		}
		pid, err := fieldString(sm, "programme_id")
		if err != nil {
			return nil, err
		}
		if !validSlug(pid) {
			return nil, errf(CodeBadField, "programme_id")
		}
		if seen[pid] {
			return nil, errf(CodeDuplicate, "programme_id")
		}
		seen[pid] = true
		sc, err := fieldString(sm, "source_checksum")
		if err != nil {
			return nil, err
		}
		if !artifact.ValidHash(sc) {
			return nil, errf(CodeBadHash, "source_checksum")
		}
		mc, err := fieldString(sm, "metadata_checksum")
		if err != nil {
			return nil, err
		}
		if !artifact.ValidHash(mc) {
			return nil, errf(CodeBadHash, "metadata_checksum")
		}
		var evRaw []json.RawMessage
		if err := json.Unmarshal(sm["evidence_checksums"], &evRaw); err != nil {
			return nil, errf(CodeBadJSON, "evidence_checksums")
		}
		if len(evRaw) > maxRefs {
			return nil, errf(CodeTooLarge, "evidence_checksums")
		}
		ev := make([]string, 0, len(evRaw))
		for _, er := range evRaw {
			var s string
			if err := json.Unmarshal(er, &s); err != nil {
				return nil, errf(CodeBadHash, "evidence_checksums")
			}
			if !artifact.ValidHash(s) {
				return nil, errf(CodeBadHash, "evidence_checksums")
			}
			ev = append(ev, s)
		}
		sources = append(sources, SourceSelection{
			ProgrammeID:       pid,
			SourceChecksum:    sc,
			MetadataChecksum:  mc,
			EvidenceChecksums: ev,
		})
	}
	return &Selection{Version: version, Sources: sources}, nil
}

// ParseRun validates the strict atlas.run/v1 record.
func ParseRun(data []byte) (*Run, error) {
	if len(data) > maxInput {
		return nil, errf(CodeTooLarge, "input")
	}
	m, err := decodeExact(data, []string{"version", "programme_id", "source_checksum", "inputs", "reports", "outputs"})
	if err != nil {
		return nil, err
	}
	version, err := fieldString(m, "version")
	if err != nil {
		return nil, err
	}
	if version != RunVersion {
		return nil, errf(CodeBadVersion, version)
	}
	pid, err := fieldString(m, "programme_id")
	if err != nil {
		return nil, err
	}
	if !validSlug(pid) {
		return nil, errf(CodeBadField, "programme_id")
	}
	sc, err := fieldString(m, "source_checksum")
	if err != nil {
		return nil, err
	}
	if !artifact.ValidHash(sc) {
		return nil, errf(CodeBadHash, "source_checksum")
	}
	inputs, err := parseRefMap(m["inputs"], inputKeys)
	if err != nil {
		return nil, err
	}
	reports, err := parseRefMap(m["reports"], reportKeys)
	if err != nil {
		return nil, err
	}
	outputs, err := parseRefMap(m["outputs"], outputKeys)
	if err != nil {
		return nil, err
	}
	return &Run{
		Version:        version,
		ProgrammeID:    pid,
		SourceChecksum: sc,
		Inputs:         inputs,
		Reports:        reports,
		Outputs:        outputs,
	}, nil
}

func parseRefMap(raw json.RawMessage, keys []string) (map[string]*Ref, error) {
	m, err := decodeExact(raw, keys)
	if err != nil {
		return nil, err
	}
	out := make(map[string]*Ref, len(keys))
	for k, v := range m {
		ref, err := parseRef(v)
		if err != nil {
			return nil, err
		}
		out[k] = ref
	}
	return out, nil
}

func parseRef(raw json.RawMessage) (*Ref, error) {
	if string(raw) == "null" {
		return nil, nil
	}
	m, err := decodeExact(raw, []string{"path", "sha256"})
	if err != nil {
		return nil, err
	}
	p, err := fieldString(m, "path")
	if err != nil {
		return nil, err
	}
	if err := checkSafeRelPath(p); err != nil {
		return nil, err
	}
	h, err := fieldString(m, "sha256")
	if err != nil {
		return nil, err
	}
	if !artifact.ValidHash(h) {
		return nil, errf(CodeBadHash, "sha256")
	}
	return &Ref{Path: p, SHA256: h}, nil
}

// Marshal returns the deterministic canonical bytes of a selection record.
func (s *Selection) Marshal() ([]byte, error) {
	srcs := make([]any, 0, len(s.Sources))
	for _, ss := range s.Sources {
		ev := make([]any, 0, len(ss.EvidenceChecksums))
		for _, e := range ss.EvidenceChecksums {
			ev = append(ev, e)
		}
		srcs = append(srcs, map[string]any{
			"programme_id":       ss.ProgrammeID,
			"source_checksum":    ss.SourceChecksum,
			"metadata_checksum":  ss.MetadataChecksum,
			"evidence_checksums": ev,
		})
	}
	return artifact.CanonicalValue(map[string]any{
		"version": s.Version,
		"sources": srcs,
	})
}

// Marshal returns the deterministic canonical bytes of a run record.
func (r *Run) Marshal() ([]byte, error) {
	return artifact.CanonicalValue(map[string]any{
		"version":         r.Version,
		"programme_id":    r.ProgrammeID,
		"source_checksum": r.SourceChecksum,
		"inputs":          refMapVal(r.Inputs),
		"reports":         refMapVal(r.Reports),
		"outputs":         refMapVal(r.Outputs),
	})
}

func refMapVal(refs map[string]*Ref) map[string]any {
	out := make(map[string]any, len(refs))
	for k, v := range refs {
		if v == nil {
			out[k] = nil
			continue
		}
		out[k] = map[string]any{"path": v.Path, "sha256": v.SHA256}
	}
	return out
}

// VerifyFiles checks that every non-nil reference in refs resolves to a bounded
// regular file under root (no symlink, no traversal) whose bytes match the
// recorded SHA-256. It performs no network lookup.
func VerifyFiles(root string, refs map[string]*Ref) error {
	for _, k := range sortedKeys(refs) {
		ref := refs[k]
		if ref == nil {
			continue
		}
		if err := verifyOne(root, ref, k); err != nil {
			return err
		}
	}
	return nil
}

func verifyOne(root string, ref *Ref, subject string) error {
	p, err := safeChild(root, ref.Path)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(p)
	if err != nil {
		return errf(CodeMissingFile, subject)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return errf(CodeSymlink, subject)
	}
	if !fi.Mode().IsRegular() {
		return errf(CodeMissingFile, subject)
	}
	id, err := fileIdentity(p)
	if err != nil {
		return errf(CodeMissingFile, subject)
	}
	if id != ref.SHA256 {
		return errf(CodeHashMismatch, subject)
	}
	return nil
}

// VerifyRun checks all referenced input/report/output files of a run record.
func VerifyRun(root string, r *Run) error {
	combined := make(map[string]*Ref, len(r.Inputs)+len(r.Reports)+len(r.Outputs))
	for k, v := range r.Inputs {
		combined["inputs."+k] = v
	}
	for k, v := range r.Reports {
		combined["reports."+k] = v
	}
	for k, v := range r.Outputs {
		combined["outputs."+k] = v
	}
	return VerifyFiles(root, combined)
}

func fileIdentity(p string) (string, error) {
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
	if !nonempty(p) {
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
	if fullAbs != rootAbs && !hasPrefix(fullAbs, rootAbs+string(os.PathSeparator)) {
		return "", errf(CodeBadPath, "escape")
	}
	return fullAbs, nil
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func sortedKeys(m map[string]*Ref) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
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

func contains(keys []string, k string) bool {
	for _, x := range keys {
		if x == k {
			return true
		}
	}
	return false
}
