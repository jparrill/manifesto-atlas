package mdrender

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
		deny []string
	}{
		{name: "When a line starts with # it should become a heading",
			in: "# Título", want: []string{"<h1>Título</h1>"}},
		{name: "When a line starts with ## it should become a second-level heading",
			in: "## 1.1 Sección", want: []string{"<h2>1.1 Sección</h2>"}},
		{name: "When two lines are separated by a blank line they should be two paragraphs",
			in: "uno\n\ndos", want: []string{"<p>uno</p>", "<p>dos</p>"}},
		{name: "When a dash starts consecutive lines they should be one list",
			in: "- uno\n- dos", want: []string{"<ul>", "<li>uno</li>", "<li>dos</li>", "</ul>"}},
		{name: "When numbers start consecutive lines they should be an ordered list",
			in: "1. uno\n2. dos", want: []string{"<ol>", "<li>uno</li>", "<li>dos</li>"}},
		{name: "When a table is present it should render header and body",
			in:   "| Dato | Valor |\n|---|---|\n| Páginas | 64 |",
			want: []string{"<th>Dato</th>", "<td>64</td>"}},
		{name: "When text has bold it should render strong",
			in: "**grueso**", want: []string{"<strong>grueso</strong>"}},
		{name: "When text has a quote mark it should stay literal",
			in: "«texto literal»", want: []string{"«texto literal»"}},
		{name: "When text contains HTML it should be escaped",
			in:   `uno <script>alert(1)</script> dos`,
			want: []string{"&lt;script&gt;"}, deny: []string{"<script>"}},
		{name: "When a link points at a sibling file it should not render as a link",
			in: "[resumen](Resumen-PP.md)", want: []string{"resumen"}, deny: []string{"href"}},
		{name: "When a link points at https it should render as a link",
			in: "[fuente](https://example.org/p.pdf)", want: []string{`href="https://example.org/p.pdf"`}},
		{name: "When a line is a comment it should be dropped",
			in: "<!-- page 4 -->\ntexto", want: []string{"texto"}, deny: []string{"page"}},
		{name: "When a horizontal rule is present it should render hr",
			in: "---", want: []string{"<hr>"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Render(tc.in)
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Fatalf("Render(%q) = %q, want %q", tc.in, got, want)
				}
			}
			for _, deny := range tc.deny {
				if strings.Contains(got, deny) {
					t.Fatalf("Render(%q) = %q, must not contain %q", tc.in, got, deny)
				}
			}
		})
	}
}
