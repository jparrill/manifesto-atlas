package runsource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeRun(info, text string) Runner {
	return func(dir, name string, args ...string) ([]byte, error) {
		if name == "pdfinfo" {
			return []byte(info), nil
		}
		return []byte(text), nil
	}
}

const infoTwo = "Title: prueba\nPages:          2\n"

func pdf(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "programa.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4 prueba"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	return path
}

func TestExtract(t *testing.T) {
	t.Run("When the PDF has two pages it should record provenance and both pages", func(t *testing.T) {
		prov, pages, err := Extract(".", pdf(t), "https://example.org/p.pdf",
			fakeRun(infoTwo, "portada\fsegunda página\f"))
		if err != nil {
			t.Fatalf("Extract: %v", err)
		}
		if prov.Pages != 2 || len(pages) != 2 {
			t.Fatalf("pages: prov=%d extracted=%d", prov.Pages, len(pages))
		}
		if len(prov.SHA256) != 64 {
			t.Fatalf("sha256 = %q", prov.SHA256)
		}
		if prov.Bytes == 0 {
			t.Fatalf("byte size not recorded")
		}
	})
	t.Run("When the origin is not https it should refuse", func(t *testing.T) {
		if _, _, err := Extract(".", pdf(t), "http://example.org/p.pdf", fakeRun(infoTwo, "a\f")); err == nil {
			t.Fatalf("expected an error")
		}
	})
	t.Run("When pdfinfo and pdftotext disagree on the page count it should fail", func(t *testing.T) {
		infoThree := "Pages:          3\n"
		if _, _, err := Extract(".", pdf(t), "https://example.org/p.pdf", fakeRun(infoThree, "una\fdos\f")); err == nil {
			t.Fatalf("expected a page-count error")
		}
	})
}

func TestMarked(t *testing.T) {
	t.Run("When a page has no text it should keep the marker and say so", func(t *testing.T) {
		got := string(Marked(Pages{"primera", "   "}))
		if !strings.Contains(got, "<!-- page 1 -->") || !strings.Contains(got, "<!-- page 2 -->") {
			t.Fatalf("markers missing: %q", got)
		}
		if !strings.Contains(got, "sin texto extraíble") {
			t.Fatalf("empty page not declared: %q", got)
		}
	})
}

func TestSections(t *testing.T) {
	t.Run("When the document has numbered headings it should detect them in order", func(t *testing.T) {
		got := Sections(Pages{
			"01 LA ESPAÑA QUE DEBEMOS SER\n",
			"1.1 La Administración será rigurosa\n1.2 El Estado dejará de asfixiar\n",
		})
		if len(got) != 3 {
			t.Fatalf("sections = %+v", got)
		}
		if got[1].Code != "1.1" || got[1].Page != 2 {
			t.Fatalf("second section = %+v", got[1])
		}
	})
	t.Run("When there are no numbered headings it should return nothing", func(t *testing.T) {
		if got := Sections(Pages{"texto suelto sin encabezados\n"}); len(got) != 0 {
			t.Fatalf("sections = %+v", got)
		}
	})
}

func TestScaffold(t *testing.T) {
	prov := Provenance{File: "p.pdf", URL: "https://example.org/p.pdf", Bytes: 12,
		SHA256: strings.Repeat("a", 64), Pages: 2}
	t.Run("When a scaffold is written it should carry provenance and placeholders", func(t *testing.T) {
		got := string(Scaffold("pp", "PP", prov, []Section{{Code: "1.1", Title: "Administración", Page: 2}}))
		for _, want := range []string{"https://example.org/p.pdf", strings.Repeat("a", 64), ScaffoldMarker, "## 1.1. Administración"} {
			if !strings.Contains(got, want) {
				t.Fatalf("scaffold misses %q:\n%s", want, got)
			}
		}
	})
	t.Run("When no section was detected it should ask for the summary by hand", func(t *testing.T) {
		got := string(Scaffold("pp", "PP", prov, nil))
		if !strings.Contains(got, "Estructura no detectada") || !strings.Contains(got, ScaffoldMarker) {
			t.Fatalf("scaffold = %s", got)
		}
	})
}

func TestRegister(t *testing.T) {
	dir := t.TempDir()
	prov := Provenance{File: "p.pdf", URL: "https://example.org/p.pdf", Bytes: 12,
		SHA256: strings.Repeat("b", 64), Pages: 2}
	t.Run("When a programme is registered twice it should not duplicate it", func(t *testing.T) {
		if _, err := Register(dir, "pp", "PP", "pp.md", "Resumen-PP.md", prov, 2); err != nil {
			t.Fatalf("Register: %v", err)
		}
		if _, err := Register(dir, "pp", "PP", "pp.md", "Resumen-PP.md", prov, 2); err != nil {
			t.Fatalf("Register twice: %v", err)
		}
		raw, err := os.ReadFile(filepath.Join(dir, "run.json"))
		if err != nil {
			t.Fatalf("read manifest: %v", err)
		}
		if strings.Count(string(raw), `"programme": "pp"`) != 2 { // one document, one source
			t.Fatalf("manifest duplicated entries:\n%s", raw)
		}
	})
}

func TestCountProposals(t *testing.T) {
	t.Run("When a summary numbers its proposals it should count the highest", func(t *testing.T) {
		got := CountProposals([]byte("1. uno\n2. dos\n**3.** tres\n[pp. 9-10]\n"))
		if got != 3 {
			t.Fatalf("CountProposals = %d, want 3", got)
		}
	})
}
