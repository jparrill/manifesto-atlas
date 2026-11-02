package candidate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// fixtureRoot is the retained Frente Obrero candidate, relative to this package.
const fixtureRoot = "../../fixtures/frente"

func TestCheckFrenteFixture(t *testing.T) {
	c, err := Check(fixtureRoot)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if c.Version != CandidateVersion {
		t.Errorf("version = %q", c.Version)
	}
	if c.Metadata.DocumentType != "general-political" {
		t.Errorf("document_type = %q, want general-political", c.Metadata.DocumentType)
	}
	if c.Metadata.Scope != "general-political" {
		t.Errorf("scope = %q, want general-political", c.Metadata.Scope)
	}
	if c.Metadata.Election != nil {
		t.Errorf("election = %+v, want null (non-electoral)", c.Metadata.Election)
	}
	if c.Metadata.PublicationDate != nil {
		t.Errorf("publication_date = %q, want null (unknown, not guessed)", *c.Metadata.PublicationDate)
	}
	if c.Metadata.Rights.Statement != nil {
		t.Errorf("rights.statement = %q, want null (unknown, not guessed)", *c.Metadata.Rights.Statement)
	}
	if c.Metadata.Origin.Authority != "maintainer-selected" {
		t.Errorf("origin.authority = %q, want maintainer-selected", c.Metadata.Origin.Authority)
	}
	if len(c.Metadata.Origin.Evidence) == 0 {
		t.Error("origin.evidence is empty, want at least one retained reference")
	}
}

func loadBase(t *testing.T) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixtureRoot, "source.json"))
	if err != nil {
		t.Fatalf("read source.json: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal source.json: %v", err)
	}
	return m
}

func parseMutated(t *testing.T, mutate func(base, meta map[string]any)) error {
	t.Helper()
	base := loadBase(t)
	meta, ok := base["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata is not an object")
	}
	mutate(base, meta)
	data, err := json.Marshal(base)
	if err != nil {
		t.Fatalf("marshal mutated: %v", err)
	}
	_, err = Parse(data)
	return err
}

func wantCode(t *testing.T, err error, code, name string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected an error, got nil", name)
	}
	if got := CodeOf(err); got != code {
		t.Fatalf("%s: code = %q, want %q (err: %v)", name, got, code, err)
	}
}

func TestParseRejects(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(base, meta map[string]any)
		code   string
	}{
		{"unknown top-level field", func(b, m map[string]any) { b["extra"] = 1 }, CodeBadField},
		{"missing metadata", func(b, m map[string]any) { delete(b, "metadata") }, CodeBadField},
		{"bad candidate version", func(b, m map[string]any) { b["version"] = "atlas.candidate/v2" }, CodeBadVersion},
		{"bad source version", func(b, m map[string]any) { m["version"] = "atlas.source.local/v2" }, CodeBadVersion},
		{"programme id conflict", func(b, m map[string]any) { m["programme_id"] = "other" }, CodeConflict},
		{"bad document type", func(b, m map[string]any) { m["document_type"] = "party-list" }, CodeBadEnum},
		{"bad scope", func(b, m map[string]any) { m["scope"] = "planetary" }, CodeBadEnum},
		{"bad year basis", func(b, m map[string]any) { m["year_basis"] = "guess" }, CodeBadEnum},
		{"bad origin authority", func(b, m map[string]any) { m["origin"].(map[string]any)["authority"] = "model" }, CodeBadEnum},
		{"bad source hash", func(b, m map[string]any) { m["source"].(map[string]any)["sha256"] = "zzz" }, CodeBadHash},
		{"electoral without election", func(b, m map[string]any) {
			m["document_type"] = "electoral-programme"
			m["election"] = nil
		}, CodeConflict},
		{"publication basis without date", func(b, m map[string]any) {
			m["year_basis"] = "publication"
			m["publication_date"] = nil
		}, CodeConflict},
		{"unsafe candidate path", func(b, m map[string]any) { b["candidate_path"] = "../escape.pdf" }, CodeBadPath},
		{"insecure origin url", func(b, m map[string]any) {
			m["source"].(map[string]any)["origin_url"] = "http://insecure.example/p.pdf"
		}, CodeBadURL},
		{"invalid retrieved_at", func(b, m map[string]any) { m["source"].(map[string]any)["retrieved_at"] = "not-a-time" }, CodeBadDate},
		{"empty evidence list", func(b, m map[string]any) { m["origin"].(map[string]any)["evidence"] = []any{} }, CodeBadField},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseMutated(t, tc.mutate)
			wantCode(t, err, tc.code, tc.name)
		})
	}
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("readdir %s: %v", src, err)
	}
	for _, e := range entries {
		sp := filepath.Join(src, e.Name())
		dp := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := os.MkdirAll(dp, 0o755); err != nil {
				t.Fatalf("mkdir %s: %v", dp, err)
			}
			copyTree(t, sp, dp)
			continue
		}
		data, err := os.ReadFile(sp)
		if err != nil {
			t.Fatalf("read %s: %v", sp, err)
		}
		if err := os.WriteFile(dp, data, 0o644); err != nil {
			t.Fatalf("write %s: %v", dp, err)
		}
	}
}

func TestVerifyRejectsMutatedCandidate(t *testing.T) {
	tmp := t.TempDir()
	copyTree(t, fixtureRoot, tmp)
	pdf := filepath.Join(tmp, "programa.pdf")
	data, err := os.ReadFile(pdf)
	if err != nil {
		t.Fatalf("read candidate: %v", err)
	}
	data[len(data)-1] ^= 0x01
	if err := os.WriteFile(pdf, data, 0o644); err != nil {
		t.Fatalf("write mutated candidate: %v", err)
	}
	_, err = Check(tmp)
	wantCode(t, err, CodeHashMismatch, "mutated candidate")
}

func TestVerifyRejectsMissingEvidence(t *testing.T) {
	tmp := t.TempDir()
	copyTree(t, fixtureRoot, tmp)
	ev := filepath.Join(tmp, "evidence", "documentos-politicos.html")
	if err := os.Remove(ev); err != nil {
		t.Fatalf("remove evidence: %v", err)
	}
	_, err := Check(tmp)
	wantCode(t, err, CodeMissingFile, "missing evidence")
}
