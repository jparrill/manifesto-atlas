// Package runsource turns a published PDF into the deterministic half of a run:
// provenance, page-marked extraction, the manifest entry and an unfilled summary
// skeleton. It never writes prose and never invents content.
package runsource

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jparrill/manifesto-atlas/internal/runmanifest"
)

// ScaffoldMarker marks a summary skeleton that a human or agent still has to
// write. The gate refuses to pass while it is present.
const ScaffoldMarker = "PARA-COMPLETAR"

// Runner executes an external tool and returns its stdout. Tests replace it.
type Runner func(dir, name string, args ...string) ([]byte, error)

// OS returns the production runner.
func OS() Runner {
	return func(dir, name string, args ...string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		var stderr strings.Builder
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("%s: %w: %s", name, err, strings.TrimSpace(stderr.String()))
		}
		return out, nil
	}
}

// Pages maps a PDF page number to its extracted text.
type Pages []string

// Provenance records the immutable identity of the source.
type Provenance struct {
	File   string
	URL    string
	Bytes  int64
	SHA256 string
	Pages  int
}

// Extract reads the PDF, records its provenance and returns the text page by page.
// It needs pdftotext and pdfinfo on PATH; nothing else, and no network.
func Extract(dir, pdfPath, url string, run Runner) (Provenance, Pages, error) {
	if !strings.HasPrefix(url, "https://") {
		return Provenance{}, nil, fmt.Errorf("the provenance URL must be https")
	}
	info, err := os.Stat(pdfPath)
	if err != nil {
		return Provenance{}, nil, err
	}
	sum, err := fileSHA256(pdfPath)
	if err != nil {
		return Provenance{}, nil, err
	}

	raw, err := run(filepath.Dir(pdfPath), "pdfinfo", filepath.Base(pdfPath))
	if err != nil {
		return Provenance{}, nil, fmt.Errorf("pdfinfo: %w", err)
	}
	pages, err := pdfPages(string(raw))
	if err != nil {
		return Provenance{}, nil, err
	}

	text, err := run(filepath.Dir(pdfPath), "pdftotext", "-layout", filepath.Base(pdfPath), "-")
	if err != nil {
		return Provenance{}, nil, fmt.Errorf("pdftotext: %w", err)
	}
	split := strings.Split(strings.ReplaceAll(string(text), "\r\n", "\n"), "\f")
	// pdftotext terminates every page with a form feed, so the split ends with an
	// empty fragment; drop exactly that one and compare with the real page count.
	if n := len(split); n > 0 && strings.TrimSpace(split[n-1]) == "" {
		split = split[:n-1]
	}
	if len(split) != pages {
		return Provenance{}, nil, fmt.Errorf("extracted %d pages, pdfinfo reports %d", len(split), pages)
	}
	return Provenance{File: filepath.Base(pdfPath), URL: url, Bytes: info.Size(),
		SHA256: sum, Pages: pages}, Pages(split), nil
}

var pagesLine = regexp.MustCompile(`(?m)^Pages:\s+(\d+)\s*$`)

func pdfPages(info string) (int, error) {
	m := pagesLine.FindStringSubmatch(info)
	if m == nil {
		return 0, fmt.Errorf("pdfinfo did not report the page count")
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("pdfinfo reported an unusable page count %q", m[1])
	}
	return n, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Marked renders the page-marked Markdown the gate reads. Empty pages get an
// explicit placeholder so the numbering never shifts.
func Marked(pages Pages) []byte {
	var b strings.Builder
	for i, page := range pages {
		fmt.Fprintf(&b, "<!-- page %d -->\n\n", i+1)
		body := strings.Trim(page, "\n")
		if strings.TrimSpace(body) == "" {
			body = "_(sin texto extraíble)_"
		}
		b.WriteString(body)
		b.WriteString("\n\n")
	}
	return []byte(b.String())
}

// Sections returns the numbered headings found in the extraction, in order, with
// the page where each one appears. It is the deterministic skeleton of a summary.
type Section struct {
	Code  string
	Title string
	Page  int
}

var headingLine = regexp.MustCompile(`^(\d{1,2})(?:\.(\d{1,2}))?\.?\s+([A-ZÁÉÍÓÚÑ][^.]{3,80})$`)

// Sections detects "01 TÍTULO" and "1.1 Título" style headings.
func Sections(pages Pages) []Section {
	var out []Section
	seen := map[string]bool{}
	for i, page := range pages {
		for _, line := range strings.Split(page, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || len(line) > 90 {
				continue
			}
			m := headingLine.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			code := m[1]
			if m[2] != "" {
				code += "." + m[2]
			}
			if seen[code] {
				continue
			}
			seen[code] = true
			out = append(out, Section{Code: code, Title: strings.TrimSpace(m[3]), Page: i + 1})
		}
	}
	return out
}

// Scaffold writes the summary skeleton: the metadata the run already knows, one
// section per detected heading and one placeholder per section that the author
// must replace with cited proposals.
func Scaffold(programme, title string, prov Provenance, sections []Section) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — resumen pendiente de redacción\n\n", title)
	b.WriteString("Este fichero lo genera `manifesto-atlas extract`. Contiene la procedencia ya\n")
	b.WriteString("verificada y el esqueleto del documento. Falta el trabajo que no es mecánico:\n")
	b.WriteString("leer el programa y escribir, propuesta a propuesta, qué dice y en qué página.\n\n")
	b.WriteString("| Dato | Valor |\n|---|---|\n")
	fmt.Fprintf(&b, "| Documento | `%s` |\n", prov.File)
	fmt.Fprintf(&b, "| Procedencia | <%s> |\n", prov.URL)
	fmt.Fprintf(&b, "| Tamaño | %d B |\n", prov.Bytes)
	fmt.Fprintf(&b, "| SHA-256 | `%s` |\n", prov.SHA256)
	fmt.Fprintf(&b, "| Páginas | %d |\n\n", prov.Pages)
	b.WriteString("Reglas: una línea por propuesta con su página `[p. N]`, comillas «» solo para\n")
	b.WriteString("texto literal, nada fuera del documento y las propuestas vacías se declaran.\n")
	b.WriteString("Cuando esté escrito, borra todos los marcadores `" + ScaffoldMarker + "`.\n\n---\n\n")

	if len(sections) == 0 {
		b.WriteString("## Estructura no detectada\n\n")
		b.WriteString("La extracción no tiene encabezados numerados reconocibles. Escribe las\n")
		b.WriteString("secciones a mano, con su rango de páginas.\n\n")
		b.WriteString("<!-- " + ScaffoldMarker + ": resumen completo del programa -->\n")
	} else {
		for _, s := range sections {
			fmt.Fprintf(&b, "## %s. %s — p. %d\n\n", s.Code, s.Title, s.Page)
			fmt.Fprintf(&b, "<!-- %s: propuestas de la sección %s (p. %d) -->\n\n", ScaffoldMarker, s.Code, s.Page)
		}
	}
	return []byte(b.String())
}

// Register adds or replaces the document and its provenance in the manifest.
func Register(dir, programme, title, sourceFile, summaryFile string, prov Provenance, proposals int) (runmanifest.Spec, error) {
	spec, err := runmanifest.Load(dir)
	if err != nil {
		spec = runmanifest.Spec{Name: filepath.Base(dir), Claims: "checks/claims.tsv",
			Allow: "checks/allow.tsv"}
	}
	doc := runmanifest.Document{Programme: programme, Source: sourceFile, Summary: summaryFile,
		Coverage: "numbers", Title: title, PDF: prov.File, Proposals: proposals}
	replaced := false
	for i := range spec.Documents {
		if spec.Documents[i].Programme == programme {
			if doc.Proposals == 0 {
				doc.Proposals = spec.Documents[i].Proposals
			}
			spec.Documents[i] = doc
			replaced = true
		}
	}
	if !replaced {
		spec.Documents = append(spec.Documents, doc)
	}
	src := runmanifest.Source{Programme: programme, URL: prov.URL, File: prov.File,
		Bytes: prov.Bytes, SHA256: prov.SHA256, Pages: prov.Pages}
	for i := range spec.Sources {
		if spec.Sources[i].Programme == programme {
			spec.Sources[i] = src
			src.Programme = ""
		}
	}
	if src.Programme != "" {
		spec.Sources = append(spec.Sources, src)
	}
	if spec.Sources == nil {
		spec.Sources = []runmanifest.Source{}
	}
	return spec, runmanifest.Save(dir, spec)
}

// SourcesMarkdown writes the provenance document of a run. It is regenerated from
// the manifest so the record and the gate can never disagree.
func SourcesMarkdown(spec runmanifest.Spec) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# Fuentes del run %s\n\n", spec.Name)
	b.WriteString("Los PDF de origen son inmutables: no se copian al repositorio ni se modifican.\n")
	b.WriteString("Se registran su URL oficial, tamaño y SHA-256 para poder comprobar que la\n")
	b.WriteString("extracción y los resúmenes corresponden a esos bytes.\n\n")
	b.WriteString("| Programa | URL oficial | Tamaño | SHA-256 | Páginas |\n|---|---|---|---|---|\n")
	for _, s := range spec.Sources {
		fmt.Fprintf(&b, "| %s | %s | %d B | `%s` | %d |\n", s.Programme, s.URL, s.Bytes, s.SHA256, s.Pages)
	}
	b.WriteString("\nComprobación reproducible de una fuente:\n\n```bash\n")
	b.WriteString("curl -sSL \"<url oficial>\" | shasum -a 256\n```\n\n")
	b.WriteString("La comprobación es puntual, no una certificación: si el partido sustituye el\n")
	b.WriteString("fichero en la misma ruta, la huella deja de coincidir y el run queda invalidado\n")
	b.WriteString("para ese documento.\n")
	return []byte(b.String())
}

// CountSourceProposals counts the numbered items of an extraction. It is the
// total a summary must cover, whatever numbering the document uses.
func CountSourceProposals(pages Pages) int {
	re := regexp.MustCompile(`(?m)^\s{0,3}(\d{1,2})\.\s+\S`)
	n := 0
	for _, page := range pages {
		n += len(re.FindAllStringSubmatch(page, -1))
	}
	return n
}

// CountProposals counts the numbered proposals of a summary and stores the total
// in the manifest, so coverage is measured against what was written.
func CountProposals(summary []byte) int {
	re := regexp.MustCompile(`(?m)^\s{0,3}(?:\*\*)?(\d{1,3})\.(?:\*\*)?\s`)
	n := 0
	for _, m := range re.FindAllStringSubmatch(string(summary), -1) {
		if v, err := strconv.Atoi(m[1]); err == nil && v > n {
			n = v
		}
	}
	return n
}
