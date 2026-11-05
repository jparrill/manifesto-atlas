package preview

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

const referencePath = "assets/reference/programmes.html"
const dataStart = "const parties="
const dataEnd = "const panel="

type partyView struct {
	Name  string      `json:"name"`
	Meta  string      `json:"meta"`
	Badge string      `json:"badge"`
	Note  string      `json:"note"`
	PDF   string      `json:"pdf"`
	MD    string      `json:"md"`
	Items [][2]string `json:"items"`
}

func hashBytes(raw []byte) string {
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ImportReference preserves owner-supplied presentation bytes, never authors or
// edits markup/script/styles. Replacing a different imported asset is explicit.
func ImportReference(path string) error {
	raw, err := regularRead(path, 256<<10)
	if err != nil {
		return err
	}
	if _, _, err := referenceParts(raw); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(referencePath), 0o755); err != nil {
		return err
	}
	old, err := regularRead(referencePath, 256<<10)
	if err == nil && !bytes.Equal(old, raw) {
		return errors.New("reference already exists with different bytes; select an explicit new reference version")
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := writeOutput(referencePath, raw); err != nil {
		return err
	}
	return saveJSON(filepath.Join(filepath.Dir(referencePath), "reference.json"), map[string]string{
		"authority": "owner-supplied existing visual reference", "sha256": hashBytes(raw),
		"binding": "only the existing parties data literal is replaced; markup/CSS/script unchanged"})
}

func referenceParts(raw []byte) ([]byte, []byte, error) {
	if bytes.Count(raw, []byte(dataStart)) != 1 || bytes.Count(raw, []byte(dataEnd)) != 1 {
		return nil, nil, errors.New("unsupported supplied reference: expected unique programme data boundaries")
	}
	start := bytes.Index(raw, []byte(dataStart)) + len(dataStart)
	end := bytes.Index(raw, []byte(dataEnd))
	if end <= start {
		return nil, nil, errors.New("unsupported supplied reference order")
	}
	return raw[:start], raw[end:], nil
}

func citeText(s canonical.Statement, citations map[string]canonical.Citation) string {
	pages := []string{}
	for _, id := range s.Citations {
		pages = append(pages, fmt.Sprintf("p. %d", citations[id].Page))
	}
	return s.Text + " [" + strings.Join(pages, ", ") + "]"
}

func views(p *canonical.Programme, pageCount int, status, pendingPDF string) map[string]partyView {
	citations := map[string]canonical.Citation{}
	for _, c := range p.Citations {
		citations[c.ID] = c
	}
	measures := map[string]canonical.Measure{}
	for _, m := range p.Measures {
		measures[m.ID] = m
	}
	kind := "Documento político general"
	if p.Metadata.DocumentType == "electoral-programme" {
		kind = "Programa electoral"
	}
	view := partyView{Name: p.Metadata.Party.Name,
		Meta:  fmt.Sprintf("%s · %s · %d páginas · %d temas · %d medidas resumidas", kind, p.Metadata.Geography.Name, pageCount, len(p.Topics), len(p.Measures)),
		Badge: "Borrador · Vista previa", Note: status, PDF: "source.pdf", MD: "summary.md", Items: [][2]string{}}
	if p.Metadata.DocumentTitle != nil {
		view.Note += " Documento: " + citeText(*p.Metadata.DocumentTitle, citations)
	}
	for _, s := range p.Summary {
		view.Note += " " + citeText(s, citations)
	}
	for _, t := range p.Topics {
		label, _ := taxonomy.TopicByID(t.ID)
		state, _ := taxonomy.StateLabel(taxonomy.State(t.State))
		parts := []string{state + "."}
		for _, s := range t.Summary {
			parts = append(parts, citeText(s, citations))
		}
		for _, mid := range t.MeasureIDs {
			m := measures[mid]
			parts = append(parts, citeText(m.Action, citations))
			for _, condition := range m.Conditions {
				parts = append(parts, citeText(condition, citations))
			}
			if m.Institution != nil {
				parts = append(parts, citeText(*m.Institution, citations))
			}
		}
		view.Items = append(view.Items, [2]string{label.Label, strings.Join(parts, " ")})
	}
	other := partyView{Name: "VOX · Castilla y León", Meta: "Documento pendiente de este flujo E2E", Badge: "Sin procesar",
		Note: "Esta pestaña no contiene un resumen validado. No se reutiliza el resumen antiguo del HTML de referencia.",
		PDF:  pendingPDF, MD: "pending.md", Items: [][2]string{}}
	result := map[string]partyView{"fo": view, "vox": other}
	if p.Metadata.Party.ID == "vox" {
		other.Name = "Frente Obrero"
		result = map[string]partyView{"fo": other, "vox": view}
	}
	return result
}

// renderHTML binds data only. The original reference's prefix/suffix, including
// all presentation code, remain byte-identical. Escape before JSON encoding:
// JSON's own escaping would be undone before the supplied innerHTML call.
func renderHTML(reference []byte, parties map[string]partyView) ([]byte, error) {
	prefix, suffix, err := referenceParts(reference)
	if err != nil {
		return nil, err
	}
	safe := map[string]partyView{}
	for key, value := range parties {
		if key != "fo" && key != "vox" {
			return nil, errors.New("unsupported reference party key")
		}
		for _, link := range []string{value.PDF, value.MD} {
			if strings.ContainsAny(link, "\"'<>\\\r\n") || strings.Contains(link, ":") || strings.HasPrefix(link, "/") {
				return nil, errors.New("unsafe generated source link")
			}
		}
		value.Name, value.Meta = html.EscapeString(value.Name), html.EscapeString(value.Meta)
		value.Badge, value.Note = html.EscapeString(value.Badge), html.EscapeString(value.Note)
		items := make([][2]string, len(value.Items))
		for i, item := range value.Items {
			items[i] = [2]string{html.EscapeString(item[0]), html.EscapeString(item[1])}
		}
		value.Items = items
		safe[key] = value
	}
	payload, err := json.Marshal(safe)
	if err != nil {
		return nil, err
	}
	out := append([]byte{}, prefix...)
	out = append(out, payload...)
	out = append(out, '\n')
	out = append(out, suffix...)
	return out, nil
}

func renderMarkdown(p *canonical.Programme, pageCount int, status string) []byte {
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n\n%s\n\nÁmbito: %s · %d páginas.\n\nSHA-256: `%s`\n\n", markdownText(p.Metadata.Party.Name), markdownText(status), markdownText(p.Metadata.Geography.Name), pageCount, p.SourceChecksum)
	citations := map[string]canonical.Citation{}
	for _, c := range p.Citations {
		citations[c.ID] = c
	}
	measures := map[string]canonical.Measure{}
	for _, m := range p.Measures {
		measures[m.ID] = m
	}
	write := func(s canonical.Statement) {
		fmt.Fprintf(&out, "- %s", markdownText(s.Text))
		for _, id := range s.Citations {
			c := citations[id]
			fmt.Fprintf(&out, " ([p. %d](source.pdf#page=%d), %s)", c.Page, c.Page, id)
		}
		out.WriteString("\n")
	}
	if p.Metadata.DocumentTitle != nil {
		write(*p.Metadata.DocumentTitle)
	}
	for _, s := range p.Summary {
		write(s)
	}
	for _, t := range p.Topics {
		topic, _ := taxonomy.TopicByID(t.ID)
		label, _ := taxonomy.StateLabel(taxonomy.State(t.State))
		fmt.Fprintf(&out, "## %s\n\n%s\n\n", topic.Label, label)
		for _, s := range t.Summary {
			write(s)
		}
		for _, mid := range t.MeasureIDs {
			m := measures[mid]
			write(m.Action)
			for _, c := range m.Conditions {
				write(c)
			}
			if m.Institution != nil {
				write(*m.Institution)
			}
		}
		out.WriteString("\n")
	}
	out.WriteString("## Evidencia literal\n\n")
	for _, c := range p.Citations {
		fmt.Fprintf(&out, "- **%s · p. %d:** %s\n", c.ID, c.Page, markdownText(c.Excerpt))
	}
	return []byte(out.String())
}

func markdownText(text string) string {
	return strings.NewReplacer("\\", "\\\\", "&", "&amp;", "<", "&lt;", ">", "&gt;", "`", "\\`", "*", "\\*", "_", "\\_", "[", "\\[", "]", "\\]", "\n", " ", "\r", " ").Replace(text)
}
