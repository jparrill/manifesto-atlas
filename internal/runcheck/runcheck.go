// Package runcheck verifies a programme run deterministically.
//
// A run is a directory holding, per programme, a page-marked extraction of the
// source PDF and a factual summary that cites it. The checks here are the
// second gate of the runbook (docs/design/programme-run.md):
//
//   - every literal quote of the summary exists in the source;
//   - every page citation points at a page that talks about that line;
//   - every numbered proposal of the source is covered by the summary;
//   - every curated claim in the run's claims file is found where it says;
//   - every section code used by the summary exists in the source.
//
// The checks are offline, read only the run directory and never call a model.
package runcheck

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/jparrill/manifesto-atlas/internal/runmanifest"
)

// Issue codes. They are stable so a run note can cite one.
const (
	CodeSpec       = "runcheck.spec"
	CodeQuote      = "runcheck.quote"
	CodePage       = "runcheck.page"
	CodeCoverage   = "runcheck.coverage"
	CodeClaim      = "runcheck.claim"
	CodeSection    = "runcheck.section"
	CodeEmbed      = "runcheck.embed"
	CodeProvenance = "runcheck.provenance"
	CodeScaffold   = "runcheck.scaffold"
)

// Issue is one failed check.
type Issue struct {
	Code   string
	File   string
	Detail string
}

func (i Issue) String() string { return i.Code + ": " + i.File + ": " + i.Detail }

// Counts are the elements each check examined.
type Counts struct {
	Quotes    int
	Pages     int
	Proposals int
	Claims    int
	Sections  int
	Allowed   int
	Embedded  int
	Refs      int
}

// Report is the outcome of a run verification.
type Report struct {
	Spec   runmanifest.Spec
	Counts Counts
	Issues []Issue
}

// OK reports whether the run passed every check.
func (r Report) OK() bool { return len(r.Issues) == 0 }

// Load reads and validates the run manifest.
func Load(dir string) (runmanifest.Spec, error) {
	spec, err := runmanifest.Load(dir)
	if err != nil {
		return runmanifest.Spec{}, fmt.Errorf("%s: %w", CodeSpec, err)
	}
	return spec, nil
}

// Verify runs every check against the run directory.
func Verify(dir string) (Report, error) {
	spec, err := Load(dir)
	if err != nil {
		return Report{}, err
	}
	rep := Report{Spec: spec}

	var claims []runmanifest.Claim
	if spec.Claims != "" {
		claims, err = runmanifest.LoadClaims(filepath.Join(dir, spec.Claims))
		if err != nil {
			return Report{}, fmt.Errorf("%s: %w", CodeSpec, err)
		}
	}

	var allow allowList
	if spec.Allow != "" {
		allow, err = loadAllow(filepath.Join(dir, spec.Allow))
		if err != nil {
			return Report{}, err
		}
	}

	for _, doc := range spec.Documents {
		source, err := LoadPages(filepath.Join(dir, doc.Source))
		if err != nil {
			return Report{}, fmt.Errorf("%s: %s: %w", CodeSpec, doc.Source, err)
		}
		summaryRaw, err := os.ReadFile(filepath.Join(dir, doc.Summary))
		if err != nil {
			return Report{}, fmt.Errorf("%s: %s: %w", CodeSpec, doc.Summary, err)
		}
		summary := string(summaryRaw)

		quotes, pages := CheckQuotesAndPages(doc.Summary, summary, source)
		rep.Counts.Quotes += quotes.total
		rep.Counts.Pages += pages.total
		rep.Issues = append(rep.Issues, quotes.issues...)
		rep.Issues = append(rep.Issues, pages.issues...)

		coverage := CheckCoverage(doc.Summary, summary, doc.Proposals, doc.Coverage, source)
		rep.Counts.Proposals += doc.Proposals
		rep.Issues = append(rep.Issues, coverage...)

		rep.Issues = append(rep.Issues, CheckScaffold(doc.Summary, summary)...)

		sections := CheckSections(doc.Summary, summary, source)
		rep.Counts.Sections += sections.total
		rep.Issues = append(rep.Issues, sections.issues...)

		for _, claim := range claims {
			if claim.Programme != doc.Programme {
				continue
			}
			rep.Counts.Claims++
			if !claimFound(claim, source) {
				rep.Issues = append(rep.Issues, Issue{
					Code: CodeClaim, File: doc.Summary,
					Detail: fmt.Sprintf("claim on page %s not found: %q", claim.Pages, claim.Phrase),
				})
			}
		}
	}
	rep.Issues = append(rep.Issues, CheckProvenance(spec)...)

	if spec.View != "" {
		embed, refs, issues := CheckView(filepath.Join(dir, spec.View), programmeNames(spec))
		rep.Counts.Embedded += embed
		rep.Counts.Refs += refs
		rep.Issues = append(rep.Issues, issues...)
	}

	kept := rep.Issues[:0]
	for _, issue := range rep.Issues {
		if rule, ok := allow.match(issue); ok {
			rep.Counts.Allowed++
			_ = rule
			continue
		}
		kept = append(kept, issue)
	}
	rep.Issues = kept
	return rep, nil
}

// ---------------------------------------------------------------- pages

// Page keeps both forms of a page: Raw for structure (list numbering, headings)
// and Norm for tolerant matching.
type Page struct {
	Raw  string
	Norm string
}

// Pages maps a PDF page number to its text.
type Pages map[int]Page

var pageMarker = regexp.MustCompile(`<!-- page (\d+) -->`)

// LoadPages reads a page-marked extraction.
func LoadPages(path string) (Pages, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(raw)
	marks := pageMarker.FindAllStringSubmatchIndex(text, -1)
	if len(marks) == 0 {
		return nil, fmt.Errorf("no <!-- page N --> markers")
	}
	pages := make(Pages, len(marks))
	for i, m := range marks {
		n, err := strconv.Atoi(text[m[2]:m[3]])
		if err != nil {
			return nil, fmt.Errorf("bad page marker %q", text[m[2]:m[3]])
		}
		end := len(text)
		if i+1 < len(marks) {
			end = marks[i+1][0]
		}
		raw := text[m[1]:end]
		pages[n] = Page{Raw: raw, Norm: Normalise(raw)}
	}
	return pages, nil
}

func (p Pages) all() string {
	nums := make([]int, 0, len(p))
	for n := range p {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	var b strings.Builder
	for _, n := range nums {
		b.WriteString(p[n].Norm)
		b.WriteByte(' ')
	}
	return b.String()
}

// ---------------------------------------------------------------- normalise

var hyphenBreak = regexp.MustCompile(`(\p{L})-\s*\n\s*(\p{L})`)
var spaceRun = regexp.MustCompile(`\s+`)

// Normalise folds a text so that layout differences (hyphenation, line breaks,
// capitals, accents, punctuation) cannot hide a match. It is the single
// definition used by every check.
func Normalise(s string) string {
	s = hyphenBreak.ReplaceAllString(s, "$1$2")
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(s) {
		switch r {
		case 'á', 'à', 'ä', 'â':
			r = 'a'
		case 'é', 'è', 'ë', 'ê':
			r = 'e'
		case 'í', 'ì', 'ï', 'î':
			r = 'i'
		case 'ó', 'ò', 'ö', 'ô':
			r = 'o'
		case 'ú', 'ù', 'ü', 'û':
			r = 'u'
		case 'ñ':
			r = 'n'
		case 'ç':
			r = 'c'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.TrimSpace(spaceRun.ReplaceAllString(b.String(), " "))
}

// ---------------------------------------------------------------- quotes

var quotedFragment = regexp.MustCompile(`«([^»]{8,200})»`)

type counted struct {
	total  int
	issues []Issue
}

// CheckQuotesAndPages reports quotes that are not literal and page citations
// that point at pages that do not talk about the line.
func CheckQuotesAndPages(file, summary string, source Pages) (counted, counted) {
	full := source.all()

	var quotes counted
	for _, m := range quotedFragment.FindAllStringSubmatch(summary, -1) {
		fragment := Normalise(m[1])
		if fragment == "" {
			continue
		}
		quotes.total++
		if !strings.Contains(full, fragment) {
			quotes.issues = append(quotes.issues, Issue{
				Code: CodeQuote, File: file,
				Detail: "not literal: «" + m[1] + "»",
			})
		}
	}

	var pages counted
	for _, line := range strings.Split(summary, "\n") {
		cited := citedPages(line)
		if len(cited) == 0 {
			continue
		}
		pages.total++
		words := keywords(line)
		if len(words) == 0 {
			continue
		}
		if !mentions(source, cited, words) {
			pages.issues = append(pages.issues, Issue{
				Code: CodePage, File: file,
				Detail: fmt.Sprintf("no cited page (%v) talks about: %s", cited, truncate(line, 90)),
			})
		}
	}
	return quotes, pages
}

var citation = regexp.MustCompile(`\[pp?\.\s*([0-9,\s–-]+)\]`)

func citedPages(line string) []int {
	var out []int
	for _, m := range citation.FindAllStringSubmatch(line, -1) {
		for _, part := range strings.Split(m[1], ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if lo, hi, ok := splitRange(part); ok {
				for p := lo; p <= hi; p++ {
					out = append(out, p)
				}
				continue
			}
			if n, err := strconv.Atoi(part); err == nil {
				out = append(out, n)
			}
		}
	}
	return out
}

func splitRange(s string) (int, int, bool) {
	for _, sep := range []string{"–", "-"} {
		if lo, hi, ok := strings.Cut(s, sep); ok {
			a, err1 := strconv.Atoi(strings.TrimSpace(lo))
			b, err2 := strconv.Atoi(strings.TrimSpace(hi))
			if err1 == nil && err2 == nil && a <= b {
				return a, b, true
			}
		}
	}
	return 0, 0, false
}

var stopWords = map[string]bool{}

func init() {
	for _, w := range strings.Fields(`de la el los las que en y a del al con por para se su sus un una unos
	unas es son como mas menos entre sin sobre tras durante desde hasta cuando donde quien quienes cual cuales
	este esta estos estas ese esa eso otro otra otros otras todo toda todos todas ser esta estan sera no ni o
	pero si ya lo le les nos me te mi tu mis tus nuestro nuestra`) {
		stopWords[w] = true
	}
}

var word = regexp.MustCompile(`[\p{L}]+`)

// keywords returns up to four distinct meaningful words of a line.
func keywords(line string) []string { return keywordsN(line, 4) }

// keywordsN returns up to n distinct meaningful words of a line.
func keywordsN(line string, n int) []string {
	var out []string
	seen := map[string]bool{}
	for _, w := range word.FindAllString(Normalise(line), -1) {
		if len([]rune(w)) < 6 || stopWords[w] || seen[w] {
			continue
		}
		seen[w] = true
		out = append(out, w)
		if len(out) == n {
			break
		}
	}
	return out
}

func mentions(source Pages, cited []int, words []string) bool {
	for _, p := range cited {
		body := source[p].Norm
		for _, w := range words {
			if strings.Contains(body, w) {
				return true
			}
		}
	}
	return false
}

// ---------------------------------------------------------------- coverage

// CheckCoverage proves that the summary covers the source proposals. With
// "numbers" the source numbers its proposals once (1..N) and every number must
// appear in the summary. With "items" the numbering restarts in each section,
// so the check counts the source items and requires each one to leave a trace:
// at least one distinctive word of its opening must appear in the summary.
func CheckCoverage(file, summary string, proposals int, mode string, source Pages) []Issue {
	var issues []Issue
	if mode == "items" {
		items := numberedItems(source)
		if len(items) != proposals {
			issues = append(issues, Issue{
				Code: CodeCoverage, File: file,
				Detail: fmt.Sprintf("source has %d numbered items, run.json declares %d", len(items), proposals),
			})
		}
		folded := Normalise(summary)
		var lost []string
		for _, item := range items {
			words := keywordsN(item, 10)
			if len(words) == 0 {
				continue
			}
			found := false
			for _, w := range words {
				if strings.Contains(folded, w) {
					found = true
					break
				}
			}
			if !found {
				lost = append(lost, truncate(item, 60))
			}
		}
		if len(lost) > 0 {
			issues = append(issues, Issue{
				Code: CodeCoverage, File: file,
				Detail: fmt.Sprintf("%d source items leave no trace in the summary: %s",
					len(lost), strings.Join(lost, " | ")),
			})
		}
		return issues
	}

	have := map[int]bool{}
	highest := 0
	for _, m := range summaryItem.FindAllStringSubmatch(summary, -1) {
		if n, err := strconv.Atoi(m[1]); err == nil {
			have[n] = true
			if n > highest {
				highest = n
			}
		}
	}
	// A summary that numbers further than the declared target proves the target
	// was undercounted, so anything past it was never verified.
	if highest > proposals {
		issues = append(issues, Issue{
			Code: CodeCoverage, File: file,
			Detail: fmt.Sprintf("the summary reaches proposal %d but the run declares %d: the target was undercounted and the rest is unverified",
				highest, proposals),
		})
	}
	var missing []string
	for n := 1; n <= proposals; n++ {
		if !have[n] {
			missing = append(missing, strconv.Itoa(n))
		}
	}
	if len(missing) > 0 {
		issues = append(issues, Issue{
			Code: CodeCoverage, File: file,
			Detail: "proposals missing from the summary: " + strings.Join(missing, ", "),
		})
	}
	return issues
}

var numberedItem = regexp.MustCompile(`(?m)^\s{0,3}(\d{1,2})\.\s+(\S.*)$`)

// summaryItem matches a proposal of a summary, with or without bold markers
// ("1. …", "**1.** …" or "**1. Title.** …"). Page citations never match.
var summaryItem = regexp.MustCompile(`(?m)^\s{0,3}(?:\*\*)?(\d{1,3})\.(?:\*\*)?\s`)

func numberedItems(source Pages) []string {
	nums := make([]int, 0, len(source))
	for n := range source {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	var items []string
	for _, n := range nums {
		for _, m := range numberedItem.FindAllStringSubmatch(source[n].Raw, -1) {
			items = append(items, m[2])
		}
	}
	return items
}

// ---------------------------------------------------------------- sections

var sectionHeading = regexp.MustCompile(`(?m)^#{2,3}\s+(\d+\.\d+)\s`)

// CheckSections requires every numbered section the summary uses to exist in
// the source, so a summary cannot invent a chapter.
func CheckSections(file, summary string, source Pages) counted {
	full := source.all()
	var c counted
	for _, m := range sectionHeading.FindAllStringSubmatch(summary, -1) {
		c.total++
		if !strings.Contains(full, Normalise(m[1])) {
			c.issues = append(c.issues, Issue{
				Code: CodeSection, File: file,
				Detail: "section " + m[1] + " does not appear in the source",
			})
		}
	}
	return c
}

// ---------------------------------------------------------------- allow

// AllowRule is a reviewed exception: a check that legitimately fires.
type AllowRule struct {
	Code   string
	File   string
	Detail string
	Reason string
}

type allowList []AllowRule

func loadAllow(path string) (allowList, error) {
	if path == "" {
		return nil, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: read %s: %w", CodeSpec, filepath.Base(path), err)
	}
	var list allowList
	for i, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 4 || f[3] == "" {
			return nil, fmt.Errorf("%s: %s:%d: want 4 tab-separated fields (code, file, detail, reason)",
				CodeSpec, filepath.Base(path), i+1)
		}
		list = append(list, AllowRule{Code: f[0], File: f[1], Detail: f[2], Reason: f[3]})
	}
	return list, nil
}

func (l allowList) match(i Issue) (AllowRule, bool) {
	for _, r := range l {
		if r.Code == i.Code && r.File == i.File && strings.Contains(i.Detail, r.Detail) {
			return r, true
		}
	}
	return AllowRule{}, false
}

// ---------------------------------------------------------------- provenance

var hexHash = regexp.MustCompile(`^[0-9a-f]{64}$`)

// CheckProvenance requires every document to carry a verifiable origin: an
// https URL, the byte size and the SHA-256 of the exact PDF analysed.
func CheckProvenance(spec runmanifest.Spec) []Issue {
	var issues []Issue
	origin := map[string]runmanifest.Source{}
	for _, s := range spec.Sources {
		origin[s.Programme] = s
	}
	for _, doc := range spec.Documents {
		src, ok := origin[doc.Programme]
		if !ok {
			issues = append(issues, Issue{Code: CodeProvenance, File: "run.json",
				Detail: "programme " + doc.Programme + " has no source entry"})
			continue
		}
		switch {
		case !strings.HasPrefix(src.URL, "https://"):
			issues = append(issues, Issue{Code: CodeProvenance, File: "run.json",
				Detail: "programme " + doc.Programme + ": origin URL is not https"})
		case !hexHash.MatchString(src.SHA256):
			issues = append(issues, Issue{Code: CodeProvenance, File: "run.json",
				Detail: "programme " + doc.Programme + ": sha256 is not a lowercase hex digest"})
		case src.Bytes <= 0:
			issues = append(issues, Issue{Code: CodeProvenance, File: "run.json",
				Detail: "programme " + doc.Programme + ": byte size is missing"})
		case src.Pages <= 0:
			issues = append(issues, Issue{Code: CodeProvenance, File: "run.json",
				Detail: "programme " + doc.Programme + ": page count is missing"})
		}
	}
	return issues
}

// ---------------------------------------------------------------- scaffold

// scaffoldMarker is written by `extract` while a summary is still unwritten.
const scaffoldMarker = "PARA-COMPLETAR"

// CheckScaffold fails while the deterministic skeleton of a summary is still
// unfilled, so an empty run cannot pass as a finished one.
func CheckScaffold(file, summary string) []Issue {
	if !strings.Contains(summary, scaffoldMarker) {
		return nil
	}
	count := strings.Count(summary, scaffoldMarker)
	return []Issue{{Code: CodeScaffold, File: file,
		Detail: fmt.Sprintf("%d placeholders %s left in the summary", count, scaffoldMarker)}}
}

// ---------------------------------------------------------------- embedding

var (
	scriptBlock = regexp.MustCompile(`(?s)<script[^>]*>.*?</script>`)
	payloadTag  = regexp.MustCompile(`(?s)<script type="application/json" id="embedded-docs">(.*?)</script>`)
	reference   = regexp.MustCompile(`(?i)\b(?:href|src|action|poster|data-src)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	stylesheet  = regexp.MustCompile(`(?i)<link[^>]+rel\s*=\s*"stylesheet"`)
)

// embeddedDoc is one document carried inside the view payload.
type embeddedDoc struct {
	Label  string `json:"label"`
	Title  string `json:"title"`
	HTML   string `json:"html"`
	Name   string `json:"name"`
	Base64 string `json:"b64"`
}

// CheckView proves the view is self-contained: every programme document travels
// inside it and the page needs no sibling file. External http(s) links are
// allowed because they are provenance references, not content.
func CheckView(path string, programmes []string) (int, int, []Issue) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, []Issue{{Code: CodeEmbed, File: filepath.Base(path),
			Detail: "cannot read the declared view"}}
	}
	page := string(raw)
	file := filepath.Base(path)
	var issues []Issue

	m := payloadTag.FindStringSubmatch(page)
	if m == nil {
		return 0, 0, []Issue{{Code: CodeEmbed, File: file,
			Detail: "no embedded-docs payload: the view would depend on sibling files"}}
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.ReplaceAll(m[1], `<\/`, `</`)), &payload); err != nil {
		return 0, 0, []Issue{{Code: CodeEmbed, File: file,
			Detail: "embedded-docs is not valid JSON: " + err.Error()}}
	}

	embedded := 0
	for _, name := range programmes {
		entry, ok := payload[name]
		if !ok {
			issues = append(issues, Issue{Code: CodeEmbed, File: file,
				Detail: "programme " + name + " has no embedded documents"})
			continue
		}
		var docs []embeddedDoc
		if err := json.Unmarshal(entry, &docs); err != nil {
			issues = append(issues, Issue{Code: CodeEmbed, File: file,
				Detail: "programme " + name + ": bad embedded document list"})
			continue
		}
		var withText int
		for _, d := range docs {
			if strings.TrimSpace(d.HTML) != "" {
				withText++
			}
		}
		if withText == 0 {
			issues = append(issues, Issue{Code: CodeEmbed, File: file,
				Detail: "programme " + name + " embeds no readable document"})
		}
		embedded += withText
	}
	if _, ok := payload["sources"]; !ok {
		issues = append(issues, Issue{Code: CodeEmbed, File: file,
			Detail: "the provenance document (sources) is not embedded"})
	}
	if !strings.Contains(page, `id="viewer-body"`) {
		issues = append(issues, Issue{Code: CodeEmbed, File: file,
			Detail: "no viewer: the embedded documents could not be displayed"})
	}

	// References outside the scripts must never point at the file system.
	static := scriptBlock.ReplaceAllString(page, "")
	refs := 0
	for _, r := range reference.FindAllStringSubmatch(static, -1) {
		refs++
		if ref := refValue(r); !isExternal(ref) {
			issues = append(issues, Issue{Code: CodeEmbed, File: file,
				Detail: "local reference in markup: " + ref})
		}
	}
	if stylesheet.MatchString(static) {
		issues = append(issues, Issue{Code: CodeEmbed, File: file,
			Detail: "external stylesheet: the view would depend on a sibling file"})
	}

	// The embedded documents are scanned too: a summary linking a sibling file
	// would break the moment the view travels alone.
	for _, name := range programmes {
		var docs []embeddedDoc
		if err := json.Unmarshal(payload[name], &docs); err != nil {
			continue
		}
		for _, d := range docs {
			for _, r := range reference.FindAllStringSubmatch(d.HTML, -1) {
				refs++
				if ref := refValue(r); !isExternal(ref) {
					issues = append(issues, Issue{Code: CodeEmbed, File: file,
						Detail: fmt.Sprintf("local reference inside the embedded document %q: %s", d.Label, ref)})
				}
			}
		}
	}
	return embedded, refs, issues
}

// refValue returns the attribute value whichever quote style was used.
func refValue(m []string) string {
	if len(m) > 2 && m[2] != "" {
		return m[2]
	}
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func isExternal(ref string) bool {
	ref = strings.TrimSpace(ref)
	if ref == "" || strings.HasPrefix(ref, "#") {
		return true
	}
	for _, prefix := range []string{"http://", "https://", "mailto:", "data:", "blob:", "tel:"} {
		if strings.HasPrefix(ref, prefix) {
			return true
		}
	}
	return false
}

func programmeNames(spec runmanifest.Spec) []string {
	names := make([]string, 0, len(spec.Documents))
	for _, d := range spec.Documents {
		names = append(names, d.Programme)
	}
	return names
}

// ---------------------------------------------------------------- claims

func claimFound(c runmanifest.Claim, source Pages) bool {
	pages := citedPages("[p. " + c.Pages + "]")
	if len(pages) == 0 {
		return false
	}
	return mentions(source, pages, []string{Normalise(c.Phrase)})
}

func truncate(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n]) + "…"
}

// Summary is a one-line-per-check human report.
func (r Report) Summary() []string {
	status := func(n int) string {
		if n == 0 {
			return "ok"
		}
		return strconv.Itoa(n) + " issues"
	}
	byCode := map[string]int{}
	for _, i := range r.Issues {
		byCode[i.Code]++
	}
	embeddingLine := "no view declared yet (run build first)"
	if r.Spec.View != "" {
		embeddingLine = fmt.Sprintf("%d documents embedded, %d references checked, %s",
			r.Counts.Embedded, r.Counts.Refs, status(byCode[CodeEmbed]))
	}
	return []string{
		fmt.Sprintf("run: %s (%d programmes)", r.Spec.Name, len(r.Spec.Documents)),
		fmt.Sprintf("  quotes       %d checked, %s", r.Counts.Quotes, status(byCode[CodeQuote])),
		fmt.Sprintf("  page refs    %d checked, %s", r.Counts.Pages, status(byCode[CodePage])),
		fmt.Sprintf("  proposals    %d expected, %s", r.Counts.Proposals, status(byCode[CodeCoverage])),
		fmt.Sprintf("  claims       %d checked, %s", r.Counts.Claims, status(byCode[CodeClaim])),
		fmt.Sprintf("  sections     %d checked, %s", r.Counts.Sections, status(byCode[CodeSection])),
		"  embedding    " + embeddingLine,
		fmt.Sprintf("  provenance   %d sources declared, %s", len(r.Spec.Sources), status(byCode[CodeProvenance])),
		fmt.Sprintf("  scaffold     %s", status(byCode[CodeScaffold])),
		fmt.Sprintf("  exceptions   %d reviewed (see the run's allow file)", r.Counts.Allowed),
	}
}

// Markdown renders the report as the per-run verification document. It is what
// `validate-content` writes next to the summary, so the evidence and the gate
// never disagree.
func (r Report) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Verificación del run %s\n\n", r.Spec.Name)
	b.WriteString("Generado por `manifesto-atlas validate-content`. Comprueba literalidad de citas,\n")
	b.WriteString("atribución de página, cobertura de propuestas, datos curados, códigos de sección,\n")
	b.WriteString("procedencia de las fuentes y que la vista embeba todos los documentos.\n\n")
	b.WriteString("## Resultado\n\n")
	b.WriteString("| Comprobación | Elementos | Resultado |\n|---|---|---|\n")
	rows := []struct {
		name   string
		total  int
		issues int
	}{
		{"Citas literales", r.Counts.Quotes, countCode(r.Issues, CodeQuote)},
		{"Atribución de página", r.Counts.Pages, countCode(r.Issues, CodePage)},
		{"Cobertura de propuestas", r.Counts.Proposals, countCode(r.Issues, CodeCoverage)},
		{"Datos curados", r.Counts.Claims, countCode(r.Issues, CodeClaim)},
		{"Códigos de sección", r.Counts.Sections, countCode(r.Issues, CodeSection)},
		{"Procedencia", len(r.Spec.Sources), countCode(r.Issues, CodeProvenance)},
		{"Andamiaje sin completar", 1, countCode(r.Issues, CodeScaffold)},
		{"Embebido en la vista", r.Counts.Embedded, countCode(r.Issues, CodeEmbed)},
	}
	for _, row := range rows {
		result := "ok"
		if row.issues > 0 {
			result = fmt.Sprintf("%d fallos", row.issues)
		}
		fmt.Fprintf(&b, "| %s | %d | %s |\n", row.name, row.total, result)
	}
	fmt.Fprintf(&b, "\nExcepciones revisadas: %d (ver `checks/allow.tsv`).\n", r.Counts.Allowed)
	if len(r.Issues) > 0 {
		b.WriteString("\n## Fallos\n\n")
		for _, i := range r.Issues {
			fmt.Fprintf(&b, "- `%s` en `%s`: %s\n", i.Code, i.File, i.Detail)
		}
	}
	b.WriteString("\n## Límites\n\n")
	b.WriteString("- La verificación es mecánica: no juzga la interpretación ni la exactitud de\n")
	b.WriteString("  las cifras que declara el partido.\n")
	b.WriteString("- El proceso que redactó los resúmenes es el mismo que ejecuta esta\n")
	b.WriteString("  comprobación, así que no sustituye una revisión humana independiente.\n")
	return b.String()
}

func countCode(issues []Issue, code string) int {
	n := 0
	for _, i := range issues {
		if i.Code == code {
			n++
		}
	}
	return n
}
