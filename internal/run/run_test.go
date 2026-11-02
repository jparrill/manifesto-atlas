package run

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

func hash(t *testing.T, b []byte) string {
	t.Helper()
	return artifact.HashBytes(b)
}

func write(t *testing.T, dir, name string, b []byte) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func codeOf(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	for err != nil {
		if ce, ok := err.(*Error); ok {
			return ce.Code
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			t.Fatalf("not a *Error: %v", err)
		}
		err = u.Unwrap()
	}
	t.Fatalf("not a *Error: %v", err)
	return ""
}

func TestSelectionRoundTrip(t *testing.T) {
	sel := &Selection{
		Version: SelectionVersion,
		Sources: []SourceSelection{{
			ProgrammeID:       "frente-obrero-programa-politico",
			SourceChecksum:    hash(t, []byte("pdf-bytes")),
			MetadataChecksum:  hash(t, []byte("meta-bytes")),
			EvidenceChecksums: []string{hash(t, []byte("evidence-bytes"))},
		}},
	}
	first, err := sel.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	second, err := sel.Marshal()
	if err != nil {
		t.Fatalf("marshal 2: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("selection marshal not deterministic:\n%s\n%s", first, second)
	}
	got, err := ParseSelection(first)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(got.Sources) != 1 || got.Sources[0].ProgrammeID != sel.Sources[0].ProgrammeID {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
	if len(got.Sources[0].EvidenceChecksums) != 1 || got.Sources[0].EvidenceChecksums[0] != sel.Sources[0].EvidenceChecksums[0] {
		t.Fatalf("evidence round-trip mismatch: %+v", got)
	}
}

func fullInputs(overrides map[string]*Ref) map[string]*Ref {
	m := make(map[string]*Ref, len(inputKeys))
	for _, k := range inputKeys {
		m[k] = nil
	}
	for k, v := range overrides {
		m[k] = v
	}
	return m
}

func fullReports(overrides map[string]*Ref) map[string]*Ref {
	m := make(map[string]*Ref, len(reportKeys))
	for _, k := range reportKeys {
		m[k] = nil
	}
	for k, v := range overrides {
		m[k] = v
	}
	return m
}

func fullOutputs(overrides map[string]*Ref) map[string]*Ref {
	m := make(map[string]*Ref, len(outputKeys))
	for _, k := range outputKeys {
		m[k] = nil
	}
	for k, v := range overrides {
		m[k] = v
	}
	return m
}

func TestRunRoundTrip(t *testing.T) {
	run := &Run{
		Version:        RunVersion,
		ProgrammeID:    "frente-obrero-programa-politico",
		SourceChecksum: hash(t, []byte("pdf-bytes")),
		Inputs:         fullInputs(map[string]*Ref{"selection": {Path: "selection.json", SHA256: hash(t, []byte("sel"))}}),
		Reports:        fullReports(map[string]*Ref{"claims": {Path: "claims.json", SHA256: hash(t, []byte("cl"))}}),
		Outputs:        fullOutputs(nil),
	}
	first, err := run.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got, err := ParseRun(first)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.Inputs["selection"] == nil || got.Inputs["selection"].Path != "selection.json" {
		t.Fatalf("inputs round-trip mismatch: %+v", got.Inputs)
	}
	if got.Inputs["draft"] != nil {
		t.Fatalf("null input not preserved: %+v", got.Inputs["draft"])
	}
	if got.Reports["claims"] == nil || got.Outputs["html"] != nil {
		t.Fatalf("report/output round-trip mismatch")
	}
	second, err := got.Marshal()
	if err != nil {
		t.Fatalf("re-marshal: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("run marshal not deterministic after round-trip")
	}
}

func TestParseSelectionRejects(t *testing.T) {
	valid := `{"version":"atlas.selection/v1","sources":[{"programme_id":"p-one","source_checksum":"` + hash(t, []byte("a")) + `","metadata_checksum":"` + hash(t, []byte("b")) + `","evidence_checksums":[]}]}`
	if _, err := ParseSelection([]byte(valid)); err != nil {
		t.Fatalf("valid selection rejected: %v", err)
	}
	// Duplicate programme_id.
	dup := `{"version":"atlas.selection/v1","sources":[` +
		`{"programme_id":"p-one","source_checksum":"` + hash(t, []byte("a")) + `","metadata_checksum":"` + hash(t, []byte("b")) + `","evidence_checksums":[]},` +
		`{"programme_id":"p-one","source_checksum":"` + hash(t, []byte("c")) + `","metadata_checksum":"` + hash(t, []byte("d")) + `","evidence_checksums":[]}]}`
	if _, err := ParseSelection([]byte(dup)); codeOf(t, err) != CodeDuplicate {
		t.Fatal("duplicate programme_id not rejected")
	}
	// Malformed JSON.
	if _, err := ParseSelection([]byte(`{"version":`)); codeOf(t, err) != CodeBadJSON {
		t.Fatal("malformed JSON not rejected")
	}
	// Bad version.
	badVer := `{"version":"atlas.selection/v2","sources":[]}`
	if _, err := ParseSelection([]byte(badVer)); codeOf(t, err) != CodeBadVersion {
		t.Fatal("bad version not rejected")
	}
	// Bad hash.
	badHash := `{"version":"atlas.selection/v1","sources":[{"programme_id":"p","source_checksum":"nope","metadata_checksum":"` + hash(t, []byte("b")) + `","evidence_checksums":[]}]}`
	if _, err := ParseSelection([]byte(badHash)); codeOf(t, err) != CodeBadHash {
		t.Fatal("bad hash not rejected")
	}
	// Unknown field.
	unknown := `{"version":"atlas.selection/v1","sources":[],"extra":1}`
	if _, err := ParseSelection([]byte(unknown)); codeOf(t, err) != CodeBadField {
		t.Fatal("unknown field not rejected")
	}
}

func TestParseRunRejectsUnknownStage(t *testing.T) {
	h := hash(t, []byte("x"))
	valid := `{"version":"atlas.run/v1","programme_id":"p","source_checksum":"` + h + `","inputs":{` +
		`"source-metadata":null,"selection":null,"extraction":null,"taxonomy":null,"programme-schema":null,"claim-schema":null,"coverage-schema":null,"draft":null,"normalizer-prompt":null,"claim-prompt":null,"coverage-prompt":null,"model-configuration":null,"tool-configuration":null` +
		`},"reports":{"claims":null,"coverage":null,"site":null},"outputs":{"markdown":null,"html":null}}`
	if _, err := ParseRun([]byte(valid)); err != nil {
		t.Fatalf("valid run rejected: %v", err)
	}
	// Unknown input stage key.
	unknownInput := `{"version":"atlas.run/v1","programme_id":"p","source_checksum":"` + h + `","inputs":{` +
		`"source-metadata":null,"selection":null,"extraction":null,"taxonomy":null,"programme-schema":null,"claim-schema":null,"coverage-schema":null,"draft":null,"normalizer-prompt":null,"claim-prompt":null,"coverage-prompt":null,"model-configuration":null,"tool-configuration":null,"bogus":null` +
		`},"reports":{"claims":null,"coverage":null,"site":null},"outputs":{"markdown":null,"html":null}}`
	if _, err := ParseRun([]byte(unknownInput)); codeOf(t, err) != CodeBadField {
		t.Fatal("unknown input stage not rejected")
	}
	// Malformed JSON.
	if _, err := ParseRun([]byte(`{"version":"atlas.run/v1"`)); codeOf(t, err) != CodeBadJSON {
		t.Fatal("malformed JSON not rejected")
	}
	// Bad version.
	if _, err := ParseRun([]byte(`{"version":"atlas.run/v2","programme_id":"p","source_checksum":"` + h + `","inputs":{},"reports":{},"outputs":{}}`)); codeOf(t, err) != CodeBadVersion {
		t.Fatal("bad version not rejected")
	}
	// Bad reference hash.
	badRef := `{"version":"atlas.run/v1","programme_id":"p","source_checksum":"` + h + `","inputs":{` +
		`"source-metadata":null,"selection":{"path":"s.json","sha256":"bad"},"extraction":null,"taxonomy":null,"programme-schema":null,"claim-schema":null,"coverage-schema":null,"draft":null,"normalizer-prompt":null,"claim-prompt":null,"coverage-prompt":null,"model-configuration":null,"tool-configuration":null` +
		`},"reports":{"claims":null,"coverage":null,"site":null},"outputs":{"markdown":null,"html":null}}`
	if _, err := ParseRun([]byte(badRef)); codeOf(t, err) != CodeBadHash {
		t.Fatal("bad reference hash not rejected")
	}
	// Traversal reference path rejected.
	trav := `{"version":"atlas.run/v1","programme_id":"p","source_checksum":"` + h + `","inputs":{` +
		`"source-metadata":null,"selection":{"path":"../escape.json","sha256":"` + h + `"},"extraction":null,"taxonomy":null,"programme-schema":null,"claim-schema":null,"coverage-schema":null,"draft":null,"normalizer-prompt":null,"claim-prompt":null,"coverage-prompt":null,"model-configuration":null,"tool-configuration":null` +
		`},"reports":{"claims":null,"coverage":null,"site":null},"outputs":{"markdown":null,"html":null}}`
	if _, err := ParseRun([]byte(trav)); codeOf(t, err) != CodeBadPath {
		t.Fatal("traversal reference path not rejected")
	}
}

func TestVerifyRunOK(t *testing.T) {
	tmp := t.TempDir()
	selBytes := []byte("selection-bytes")
	write(t, tmp, "selection.json", selBytes)
	h := hash(t, selBytes)
	run := &Run{
		Version:        RunVersion,
		ProgrammeID:    "p",
		SourceChecksum: hash(t, []byte("src")),
		Inputs:         map[string]*Ref{"selection": {Path: "selection.json", SHA256: h}},
		Reports:        map[string]*Ref{"claims": nil, "coverage": nil, "site": nil},
		Outputs:        map[string]*Ref{"markdown": nil, "html": nil},
	}
	if err := VerifyRun(tmp, run); err != nil {
		t.Fatalf("valid run files rejected: %v", err)
	}
}

func TestVerifyRunStaleMutation(t *testing.T) {
	tmp := t.TempDir()
	write(t, tmp, "draft.json", []byte("original"))
	h := hash(t, []byte("original"))
	run := &Run{
		Version:        RunVersion,
		ProgrammeID:    "p",
		SourceChecksum: hash(t, []byte("src")),
		Inputs:         map[string]*Ref{"draft": {Path: "draft.json", SHA256: h}},
		Reports:        map[string]*Ref{"claims": nil, "coverage": nil, "site": nil},
		Outputs:        map[string]*Ref{"markdown": nil, "html": nil},
	}
	if err := VerifyRun(tmp, run); err != nil {
		t.Fatalf("pre-mutation verify failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "draft.json"), []byte("mutated"), 0o644); err != nil {
		t.Fatalf("mutate: %v", err)
	}
	if codeOf(t, VerifyRun(tmp, run)) != CodeHashMismatch {
		t.Fatal("stale mutation not detected")
	}
}

func TestVerifyRunMissingReference(t *testing.T) {
	tmp := t.TempDir()
	run := &Run{
		Version:        RunVersion,
		ProgrammeID:    "p",
		SourceChecksum: hash(t, []byte("src")),
		Inputs:         map[string]*Ref{"selection": {Path: "absent.json", SHA256: hash(t, []byte("x"))}},
		Reports:        map[string]*Ref{"claims": nil, "coverage": nil, "site": nil},
		Outputs:        map[string]*Ref{"markdown": nil, "html": nil},
	}
	if codeOf(t, VerifyRun(tmp, run)) != CodeMissingFile {
		t.Fatal("missing reference not detected")
	}
}

func TestVerifyRunSymlink(t *testing.T) {
	tmp := t.TempDir()
	target := write(t, tmp, "target.json", []byte("data"))
	link := filepath.Join(tmp, "link.json")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not supported: %v", err)
	}
	h := hash(t, []byte("data"))
	run := &Run{
		Version:        RunVersion,
		ProgrammeID:    "p",
		SourceChecksum: hash(t, []byte("src")),
		Inputs:         map[string]*Ref{"selection": {Path: "link.json", SHA256: h}},
		Reports:        map[string]*Ref{"claims": nil, "coverage": nil, "site": nil},
		Outputs:        map[string]*Ref{"markdown": nil, "html": nil},
	}
	if codeOf(t, VerifyRun(tmp, run)) != CodeSymlink {
		t.Fatal("symlink reference not rejected")
	}
}
