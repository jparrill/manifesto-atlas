package runcheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sourceMarked = `<!-- page 1 -->

INTRODUCCIÓN
El programa propone subir el salario mínimo.

<!-- page 2 -->

VIVIENDA
1. Construir vivienda pública.
2. Prohibir los pisos turísticos.
`

const summaryGood = `# Resumen de prueba

| Dato | Valor |
|---|---|
| Propuestas | 2 |

## 1. Vivienda — p. 2

1. «Construir vivienda pública» y ampliar el parque público [p. 2]
2. Prohibir los pisos turísticos [p. 2]
`

// run writes a run directory and returns its path.
func run(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	files["run.json"] = mustSpec(t, files)
	if _, ok := files["checks/claims.tsv"]; !ok {
		files["checks/claims.tsv"] = "# sin datos curados\n"
	}
	if _, ok := files["checks/allow.tsv"]; !ok {
		files["checks/allow.tsv"] = "# sin excepciones\n"
	}
	for name, body := range files {
		if name == "run.json" && body == "" {
			continue
		}
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func mustSpec(t *testing.T, files map[string]string) string {
	t.Helper()
	if spec, ok := files["run.json"]; ok && spec != "" {
		return spec
	}
	return `{
  "name": "prueba",
  "claims": "checks/claims.tsv",
  "allow": "checks/allow.tsv",
  "documents": [
    {"programme": "pp", "source": "pp.md", "summary": "Resumen-PP-2026.md", "proposals": 2}
  ]
}`
}

func codes(issues []Issue) []string {
	out := make([]string, 0, len(issues))
	for _, i := range issues {
		out = append(out, i.Code)
	}
	return out
}

func has(issues []Issue, code string) bool {
	for _, i := range issues {
		if i.Code == code {
			return true
		}
	}
	return false
}

func TestNormaliseFoldsLayout(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "When a word is split by a hyphen at line end it should join", in: "inme-\ndiato", want: "inmediato"},
		{name: "When text has accents and capitals it should fold them", in: "Indemnización PÚBLICA", want: "indemnizacion publica"},
		{name: "When text has punctuation it should become separators", in: "«sí, señor»", want: "si senor"},
		{name: "When text has runs of spaces it should collapse them", in: "a    b\n\n c", want: "a b c"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Normalise(tc.in); got != tc.want {
				t.Fatalf("Normalise(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestLoadPagesRequiresMarkers(t *testing.T) {
	t.Run("When a source has page markers it should map every page", func(t *testing.T) {
		pages, err := LoadPages(filepath.Join(run(t, map[string]string{"pp.md": sourceMarked}), "pp.md"))
		if err != nil {
			t.Fatalf("LoadPages: %v", err)
		}
		if len(pages) != 2 {
			t.Fatalf("got %d pages, want 2", len(pages))
		}
		if !strings.Contains(pages[2].Norm, "pisos turisticos") {
			t.Fatalf("page 2 not normalised: %q", pages[2].Norm)
		}
	})
	t.Run("When a source has no markers it should fail", func(t *testing.T) {
		defer func() { _ = recover() }()
		if _, err := LoadPages(filepath.Join(run(t, map[string]string{"pp.md": "sin marcas"}), "pp.md")); err == nil {
			t.Fatalf("expected an error")
		}
	})
}

func TestQuotesMustBeLiteral(t *testing.T) {
	t.Run("When a quote exists in the source it should pass", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md": sourceMarked, "Resumen-PP-2026.md": summaryGood,
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !report.OK() {
			t.Fatalf("unexpected issues: %v", report.Issues)
		}
	})
	t.Run("When a quote is a paraphrase it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": strings.Replace(summaryGood, "«Construir vivienda pública»", "«Levantar casas públicas»", 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeQuote) {
			t.Fatalf("want %s, got %v", CodeQuote, codes(report.Issues))
		}
	})
	t.Run("When a citation points at the wrong page it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": strings.Replace(summaryGood, "y ampliar el parque público [p. 2]", "y ampliar el parque público [p. 1]", 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodePage) {
			t.Fatalf("want %s, got %v", CodePage, codes(report.Issues))
		}
	})
}

func TestCoverage(t *testing.T) {
	t.Run("When a numbered proposal is missing it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": "# Resumen\n\n## 1. Vivienda — p. 2\n\n1. «Construir vivienda pública» [p. 2]\n",
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeCoverage) {
			t.Fatalf("want %s, got %v", CodeCoverage, codes(report.Issues))
		}
	})
	t.Run("When numbering restarts per section the item count must match", func(t *testing.T) {
		spec := `{
  "name": "prueba",
  "claims": "checks/claims.tsv",
  "documents": [
    {"programme": "pp", "source": "pp.md", "summary": "Resumen-PP-2026.md", "proposals": 3, "coverage": "items"}
  ]
}`
		report, err := Verify(run(t, map[string]string{
			"run.json": spec, "pp.md": sourceMarked, "Resumen-PP-2026.md": summaryGood,
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeCoverage) {
			t.Fatalf("want %s for 2 items against 3 declared, got %v", CodeCoverage, codes(report.Issues))
		}
	})
	t.Run("When an item leaves no trace in the summary it should be reported", func(t *testing.T) {
		spec := `{
  "name": "prueba",
  "claims": "checks/claims.tsv",
  "documents": [
    {"programme": "pp", "source": "pp.md", "summary": "Resumen-PP-2026.md", "proposals": 2, "coverage": "items"}
  ]
}`
		report, err := Verify(run(t, map[string]string{
			"run.json":           spec,
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": "# nada\n\nsin propuestas\n",
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeCoverage) {
			t.Fatalf("want %s, got %v", CodeCoverage, codes(report.Issues))
		}
	})
}

func TestSectionsMustExistInSource(t *testing.T) {
	t.Run("When the summary uses a section code of the source it should pass", func(t *testing.T) {
		source := strings.Replace(sourceMarked, "VIVIENDA", "3.3 Infraestructuras\nVIVIENDA", 1)
		report, err := Verify(run(t, map[string]string{
			"pp.md":              source,
			"Resumen-PP-2026.md": strings.Replace(summaryGood, "## 1. Vivienda", "## 3.3 Vivienda", 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if has(report.Issues, CodeSection) {
			t.Fatalf("unexpected %s: %v", CodeSection, report.Issues)
		}
	})
	t.Run("When the summary invents a section code it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": strings.Replace(summaryGood, "## 1. Vivienda", "## 7.2 Vivienda", 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeSection) {
			t.Fatalf("want %s, got %v", CodeSection, codes(report.Issues))
		}
	})
}

func TestClaims(t *testing.T) {
	const claims = "pp\t2\tconstruir vivienda publica\n"
	t.Run("When a curated claim sits on its page it should pass", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md": sourceMarked, "Resumen-PP-2026.md": summaryGood, "checks/claims.tsv": claims,
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !report.OK() {
			t.Fatalf("unexpected issues: %v", report.Issues)
		}
	})
	t.Run("When a curated claim sits on another page it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md": sourceMarked, "Resumen-PP-2026.md": summaryGood,
			"checks/claims.tsv": "pp\t1\tconstruir vivienda publica\n",
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeClaim) {
			t.Fatalf("want %s, got %v", CodeClaim, codes(report.Issues))
		}
	})
	t.Run("When a claim row is malformed it should fail loudly", func(t *testing.T) {
		if _, err := Verify(run(t, map[string]string{
			"pp.md": sourceMarked, "Resumen-PP-2026.md": summaryGood,
			"checks/claims.tsv": "pp\t2\n",
		})); err == nil {
			t.Fatalf("expected an error")
		}
	})
}

func TestAllowListMakesExceptionsExplicit(t *testing.T) {
	broken := strings.Replace(summaryGood,
		"1. «Construir vivienda pública» y ampliar el parque público [p. 2]",
		"1. «Levantar casas públicas» de vivienda [p. 2]", 1)
	t.Run("When an exception has a reason it should be accepted and counted", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": broken,
			"checks/allow.tsv":   "runcheck.quote\tResumen-PP-2026.md\tLevantar casas públicas\tcita revisada a mano\n",
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !report.OK() {
			t.Fatalf("unexpected issues: %v", report.Issues)
		}
		if report.Counts.Allowed != 1 {
			t.Fatalf("allowed = %d, want 1", report.Counts.Allowed)
		}
	})
	t.Run("When an exception has no reason it should be rejected", func(t *testing.T) {
		if _, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": broken,
			"checks/allow.tsv":   "runcheck.quote\tResumen-PP-2026.md\tLevantar\t\n",
		})); err == nil {
			t.Fatalf("expected an error for an exception without a reason")
		}
	})
	t.Run("When an unrelated issue appears the exception should not hide it", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"pp.md":              sourceMarked,
			"Resumen-PP-2026.md": strings.Replace(broken, "## 1. Vivienda", "## 9.9 Vivienda", 1),
			"checks/allow.tsv":   "runcheck.quote\tResumen-PP-2026.md\tLevantar casas públicas\tcita revisada\n",
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeSection) {
			t.Fatalf("want %s, got %v", CodeSection, codes(report.Issues))
		}
	})
}

func TestViewMustBeSelfContained(t *testing.T) {
	withPayload := func(payload string) string {
		return `<!doctype html><html><body><div id="viewer-body"></div>` +
			`<script type="application/json" id="embedded-docs">` + payload + `</script>` +
			`<script>var x = "<a href='sibling.md'>x</a>";</script></body></html>`
	}
	docs := `{"pp":[{"label":"Resumen completo","title":"t","html":"<h1>Resumen</h1>"}],"sources":{"html":"<h1>Fuentes</h1>"}}`
	spec := func(view string) string {
		return `{"name":"prueba","claims":"checks/claims.tsv","view":"` + view + `",
  "documents":[{"programme":"pp","source":"pp.md","summary":"Resumen-PP-2026.md","proposals":2}]}`
	}

	t.Run("When every document is embedded and nothing points outside it should pass", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood, "vista.html": withPayload(docs),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !report.OK() {
			t.Fatalf("unexpected issues: %v", report.Issues)
		}
		if report.Counts.Embedded != 1 {
			t.Fatalf("embedded = %d, want 1", report.Counts.Embedded)
		}
	})
	t.Run("When the view has no payload it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood,
			"vista.html":         `<html><body><div id="viewer-body"></div></body></html>`,
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeEmbed) {
			t.Fatalf("want %s, got %v", CodeEmbed, codes(report.Issues))
		}
	})
	t.Run("When the markup links a sibling file it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood,
			"vista.html": strings.Replace(withPayload(docs), "<div id=\"viewer-body\">",
				`<div id="viewer-body"><a href="Resumen-PP-2026.md">resumen</a>`, 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeEmbed) {
			t.Fatalf("want %s, got %v", CodeEmbed, codes(report.Issues))
		}
	})
	t.Run("When the view loads an external stylesheet it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood,
			"vista.html": strings.Replace(withPayload(docs), "</body>",
				`<link rel="stylesheet" href="estilo.css"></body>`, 1),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeEmbed) {
			t.Fatalf("want %s, got %v", CodeEmbed, codes(report.Issues))
		}
	})
	t.Run("When a programme has no embedded document it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood,
			"vista.html":         withPayload(`{"pp":[{"label":"x","html":""}],"sources":{"html":"<h1>f</h1>"}}`),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeEmbed) {
			t.Fatalf("want %s, got %v", CodeEmbed, codes(report.Issues))
		}
	})
	t.Run("When the embedded document links a sibling file it should be reported", func(t *testing.T) {
		report, err := Verify(run(t, map[string]string{
			"run.json": spec("vista.html"), "pp.md": sourceMarked,
			"Resumen-PP-2026.md": summaryGood,
			"vista.html":         withPayload(`{"pp":[{"label":"x","html":"<a href='notas.md'>x</a>"}],"sources":{"html":"f"}}`),
		}))
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !has(report.Issues, CodeEmbed) {
			t.Fatalf("want %s, got %v", CodeEmbed, codes(report.Issues))
		}
	})
}
