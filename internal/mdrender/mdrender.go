// Package mdrender renders the Markdown subset used by programme summaries.
//
// The subset is deliberately small and closed: headings, paragraphs, ordered and
// unordered lists, tables, blockquotes, horizontal rules, bold, italic, inline
// code and links. Anything else is escaped and shown as text, so a summary
// cannot inject markup into the view or the printable PDF.
package mdrender

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Render converts the supported Markdown subset to HTML.
func Render(md string) string {
	var b strings.Builder
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == "":
			i++
		case strings.HasPrefix(trimmed, "<!--"):
			// HTML comments are not content.
			for i < len(lines) && !strings.Contains(lines[i], "-->") {
				i++
			}
			i++
		case strings.HasPrefix(trimmed, "### "):
			fmt.Fprintf(&b, "<h3>%s</h3>\n", inline(strings.TrimPrefix(trimmed, "### ")))
			i++
		case strings.HasPrefix(trimmed, "## "):
			fmt.Fprintf(&b, "<h2>%s</h2>\n", inline(strings.TrimPrefix(trimmed, "## ")))
			i++
		case strings.HasPrefix(trimmed, "# "):
			fmt.Fprintf(&b, "<h1>%s</h1>\n", inline(strings.TrimPrefix(trimmed, "# ")))
			i++
		case trimmed == "---" || trimmed == "***":
			b.WriteString("<hr>\n")
			i++
		case strings.HasPrefix(trimmed, ">"):
			var quote []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), ">") {
				quote = append(quote, strings.TrimPrefix(strings.TrimSpace(lines[i]), ">"))
				i++
			}
			fmt.Fprintf(&b, "<blockquote><p>%s</p></blockquote>\n", inline(strings.Join(quote, " ")))
		case looksLikeTable(lines, i):
			i = table(&b, lines, i)
		case isList(trimmed):
			i = list(&b, lines, i)
		default:
			var para []string
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				if t == "" || isList(t) || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "<!--") ||
					looksLikeTable(lines, i) || strings.HasPrefix(t, ">") || t == "---" {
					break
				}
				para = append(para, t)
				i++
			}
			if len(para) > 0 {
				fmt.Fprintf(&b, "<p>%s</p>\n", inline(strings.Join(para, " ")))
			} else {
				i++
			}
		}
	}
	return b.String()
}

func isList(t string) bool {
	if strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") {
		return true
	}
	_, _, ok := orderedItem(t)
	return ok
}

var ordered = regexp.MustCompile(`^(\d{1,3})[.)]\s+(.*)$`)

func orderedItem(t string) (string, string, bool) {
	m := ordered.FindStringSubmatch(t)
	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}

func list(b *strings.Builder, lines []string, i int) int {
	tag := "ul"
	if _, _, ok := orderedItem(strings.TrimSpace(lines[i])); ok {
		tag = "ol"
	}
	fmt.Fprintf(b, "<%s>\n", tag)
	for i < len(lines) {
		t := strings.TrimSpace(lines[i])
		if !isList(t) {
			break
		}
		text := t
		if _, rest, ok := orderedItem(t); ok {
			text = rest
		} else {
			text = strings.TrimPrefix(strings.TrimPrefix(t, "- "), "* ")
		}
		fmt.Fprintf(b, "<li>%s</li>\n", inline(text))
		i++
	}
	fmt.Fprintf(b, "</%s>\n", tag)
	return i
}

func looksLikeTable(lines []string, i int) bool {
	if i+1 >= len(lines) {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(lines[i]), "|") &&
		strings.Contains(lines[i+1], "---") && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "|")
}

func table(b *strings.Builder, lines []string, i int) int {
	b.WriteString("<table>\n<thead>\n<tr>")
	for _, cell := range cells(lines[i]) {
		fmt.Fprintf(b, "<th>%s</th>", inline(cell))
	}
	b.WriteString("</tr>\n</thead>\n<tbody>\n")
	i += 2
	for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
		b.WriteString("<tr>")
		for _, cell := range cells(lines[i]) {
			fmt.Fprintf(b, "<td>%s</td>", inline(cell))
		}
		b.WriteString("</tr>\n")
		i++
	}
	b.WriteString("</tbody>\n</table>\n")
	return i
}

func cells(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.TrimPrefix(row, "|")
	row = strings.TrimSuffix(row, "|")
	parts := strings.Split(row, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

var (
	bold   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italic = regexp.MustCompile(`(^|[\s(])\*([^*\s][^*]*)\*`)
	code   = regexp.MustCompile("`([^`]+)`")
	link   = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
)

// inline escapes text and then applies the inline marks.
func inline(s string) string {
	out := html.EscapeString(s)
	out = code.ReplaceAllString(out, "<code>$1</code>")
	out = bold.ReplaceAllString(out, "<strong>$1</strong>")
	out = italic.ReplaceAllString(out, "$1<em>$2</em>")
	out = link.ReplaceAllStringFunc(out, func(m string) string {
		parts := link.FindStringSubmatch(m)
		target := parts[2]
		if !safeLink(target) {
			return parts[1]
		}
		return fmt.Sprintf(`<a href="%s" rel="noopener">%s</a>`, target, parts[1])
	})
	return out
}

// safeLink keeps the rendered documents free of file-system references: only
// http(s) and fragment targets survive; anything else is left as plain text.
func safeLink(target string) bool {
	return strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "http://") ||
		strings.HasPrefix(target, "#")
}
