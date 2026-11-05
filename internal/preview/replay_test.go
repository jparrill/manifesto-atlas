package preview

import (
	"bytes"
	"encoding/json"
	"html"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/semantic"
)

// Replay the genuinely retained run; no PDF tools, model, browser or network are
// needed. Each case uses private copies so historical evidence is never edited.
func TestStoredRealPreview(t *testing.T) {
	source, err := filepath.Abs("../../runs/preview/frente-e2e")
	if err != nil {
		t.Fatal(err)
	}
	reference, err := os.ReadFile("../../assets/reference/programmes.html")
	if err != nil {
		t.Fatal(err)
	}
	prepare := func(t *testing.T) {
		t.Helper()
		t.Chdir(t.TempDir())
		if err := os.MkdirAll("assets/reference", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(referencePath, reference, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir("run", 0o755); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(source)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			raw, err := os.ReadFile(filepath.Join(source, entry.Name()))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join("run", entry.Name()), raw, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Run("When a real stored run is replayed offline, it should preserve every statement and deterministic private output", func(t *testing.T) {
		prepare(t)
		if err := Render("run"); err != nil {
			t.Fatal(err)
		}
		first, _ := os.ReadFile("run/index.html")
		markdown, _ := os.ReadFile("run/summary.md")
		if err := Render("run"); err != nil {
			t.Fatal(err)
		}
		second, _ := os.ReadFile("run/index.html")
		if !bytes.Equal(first, second) {
			t.Fatal("offline HTML changed")
		}
		if !bytes.Contains(first, []byte("BORRADOR")) || !bytes.Contains(markdown, []byte("evidencia insuficiente")) {
			t.Fatal("the real insufficient outcome was hidden")
		}
		var p canonical.Programme
		raw, _ := os.ReadFile("run/programme.json")
		if err := json.Unmarshal(raw, &p); err != nil {
			t.Fatal(err)
		}
		prefix, suffix, err := referenceParts(first)
		if err != nil {
			t.Fatal(err)
		}
		var parties map[string]partyView
		if err := json.Unmarshal(bytes.TrimSpace(first[len(prefix):len(first)-len(suffix)]), &parties); err != nil {
			t.Fatal(err)
		}
		if len(parties["fo"].Items) != 21 || len(parties["vox"].Items) != 0 {
			t.Fatal("topics dropped or old unprocessed-party content reused")
		}
		visible := parties["fo"].Note
		for _, item := range parties["fo"].Items {
			visible += " " + item[1]
		}
		check := func(s canonical.Statement) {
			t.Helper()
			if !bytes.Contains(markdown, []byte(s.Text)) || !strings.Contains(html.UnescapeString(visible), s.Text) {
				t.Fatalf("public statement lost: %s", s.ID)
			}
		}
		for _, topic := range p.Topics {
			for _, s := range topic.Summary {
				check(s)
			}
		}
		for _, m := range p.Measures {
			check(m.Action)
			for _, s := range m.Conditions {
				check(s)
			}
		}
	})
	for _, name := range []string{"programme.json", "source.pdf", "model.json", "manifesto-claim-validator.json"} {
		t.Run("When retained "+name+" changes, it should reject stale output", func(t *testing.T) {
			prepare(t)
			if err := os.WriteFile(filepath.Join("run", name), []byte("tampered"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := Render("run"); err == nil || !strings.Contains(err.Error(), "stale/missing run input") {
				t.Fatalf("tampering not rejected: %v", err)
			}
		})
	}
	t.Run("When a semantic report is absent, it should reject rather than manufacture approval", func(t *testing.T) {
		prepare(t)
		if err := os.Remove("run/manifesto-coverage-validator.json"); err != nil {
			t.Fatal(err)
		}
		if err := Render("run"); err == nil {
			t.Fatal("missing whole-source report accepted")
		}
	})
	t.Run("When generated HTML is a symlink to the original, it should reject without changing source bytes", func(t *testing.T) {
		prepare(t)
		original, _ := os.ReadFile("run/source.pdf")
		if err := os.Remove("run/index.html"); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("source.pdf", "run/index.html"); err != nil {
			t.Fatal(err)
		}
		if err := Render("run"); err == nil || !strings.Contains(err.Error(), "unsafe output leaf") {
			t.Fatalf("output symlink not rejected: %v", err)
		}
		current, _ := os.ReadFile("run/source.pdf")
		if !bytes.Equal(original, current) {
			t.Fatal("immutable source changed through an output symlink")
		}
	})
	t.Run("When output is hard-linked to the source, it should replace only the generated output inode", func(t *testing.T) {
		prepare(t)
		original, _ := os.ReadFile("run/source.pdf")
		if err := os.Remove("run/index.html"); err != nil {
			t.Fatal(err)
		}
		if err := os.Link("run/source.pdf", "run/index.html"); err != nil {
			t.Fatal(err)
		}
		if err := Render("run"); err != nil {
			t.Fatal(err)
		}
		current, _ := os.ReadFile("run/source.pdf")
		if !bytes.Equal(original, current) {
			t.Fatal("immutable source inode was modified")
		}
	})
	t.Run("When a non-passing report belongs to another draft, it should reject stale outcomes too", func(t *testing.T) {
		prepare(t)
		raw, _ := os.ReadFile("run/run.json")
		var record Record
		if err := json.Unmarshal(raw, &record); err != nil {
			t.Fatal(err)
		}
		role := semantic.RoleCoverageValidator
		outcome := record.Reports[role]
		outcome.Response.Report.Inputs["draft"] = hashBytes([]byte("a different draft"))
		record.Reports[role] = outcome
		name := role + ".json"
		if err := saveJSON(filepath.Join("run", name), outcome); err != nil {
			t.Fatal(err)
		}
		changed, _ := os.ReadFile(filepath.Join("run", name))
		record.Inputs[name] = hashBytes(changed)
		if err := saveJSON("run/run.json", record); err != nil {
			t.Fatal(err)
		}
		if err := Render("run"); err == nil || !strings.Contains(err.Error(), "different inputs") {
			t.Fatalf("stale insufficient report was accepted: %v", err)
		}
	})
}
