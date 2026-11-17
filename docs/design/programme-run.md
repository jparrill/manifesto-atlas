# Programme runs: from three PDFs to a verified summary

This is the contract for a **programme run**: taking published programmes,
converting them, writing factual summaries with page citations, proving those
summaries against the source, and shipping a view that needs nothing but
itself. It records the process that produced [`runs/29n2026/`](../../runs/29n2026/README.md)
so the next run follows the same shape instead of reinventing it.

A run is a **private draft**. It is not a publication and it does not enter the
canonical catalogue: it exists to produce evidence someone can review, with the
trust boundary in `docs/design/provenance-trust.md` still in force.

## Where each responsibility lives

| Layer | Who does it | Why |
| --- | --- | --- |
| Provenance (URL, bytes, SHA-256) | The run, checked by a human command | The source must be identifiable months later |
| Extraction, page markers, checks | Go and shell, deterministic | Same input, same output, no judgement |
| Reading the programme and writing the summary | An agent, with the rules below | Requires judgement; is reviewed, not trusted |
| Verifying the summary | `make check-run`, deterministic | Mechanical claims have mechanical proofs |
| Presentation | Owner-supplied template plus a generator | The project does not let an agent invent design |

## 1. Run directory

```
runs/<run>/
  README.md            index: scope, artifacts, gate result, limits
  SOURCES.md           official URL, size and SHA-256 per source
  run.json             manifest read by the gate
  checks/claims.tsv    curated factual data: programme, page, phrase
  checks/allow.tsv     reviewed exceptions, each with a reason
  text/*.txt           raw extraction (pdftotext -layout)
  <source>.md          extraction with <!-- page N --> markers
  Resumen-*.md         factual summary with page citations
  Resumen-*.pdf        printable summary
  Verificacion-*.md    what the gate checked, what it corrected, what remains
  Programas-*.html     the view (self-contained)
  tools/               extraction, generators, reference scripts
```

`run.json` declares the run:

```json
{
  "name": "29-N 2026",
  "claims": "checks/claims.tsv",
  "allow": "checks/allow.tsv",
  "view": "Programas-electorales-2026.html",
  "documents": [
    {"programme": "vox", "source": "vox.md", "summary": "Resumen-VOX-2026.md",
     "proposals": 115, "coverage": "numbers"}
  ]
}
```

`coverage` is `numbers` when the document numbers its proposals once (1..N) and
`items` when the numbering restarts in every section.

## 2. Provenance first

Register the official URL of each source and the SHA-256 of the exact bytes
analysed, then prove the URL still serves them:

```bash
curl -sSL "<official-url>" | shasum -a 256   # must equal SOURCES.md
```

Rules:

- Source PDFs are **immutable and never committed**. They are identified by URL
  and hash. A run without a verifiable URL is a run without provenance and must
  say so.
- The check is a point-in-time verification by the maintainer, not a
  certification. If the publisher replaces the file at the same path, the hash
  stops matching and the run is invalid for that document.
- Record what the file itself declares (pages, `ModDate`, producer). If a
  document claims to be something other than an electoral programme, say it: the
  Frente Obrero document states it is a political programme, and the summary
  repeats that.

## 3. Extraction

```bash
pdftotext -layout "<pdf>" text/<name>.txt
python3 tools/pagemark.py text/<name>.txt <name>.md
```

- One `<!-- page N -->` marker per PDF page, no reflow, no OCR. If a page has no
  text layer, the marker stays with an explicit placeholder and the summary
  enumerates those pages.
- Never edit the extraction to make a quote fit. If the text is odd, the summary
  says what the text says.

## 4. Writing the summary

A summary states what the document says, in Spanish, for the public. It does not
evaluate viability, legality or convenience, and it does not add outside data.

- **One line per proposal, with its page**: `[p. 4]` or `[pp. 41-43]`.
- **Quotes are literal**: anything inside «» must exist in the document, character
  for character after folding layout. Do not put a paraphrase, a negation or a
  title rewrite inside «».
- **Numbers as written**: if the document says «de cero a tres años», do not quote
  «de 0 a 3 años».
- **Empty proposals are reported, not filled**: the PP document leaves eight
  proposals ("permanezca atento a la campaña"); the summary lists them as empty.
- **Keep the document's own structure.** If it has four blocks and 28
  subapartados, the summary uses them. Invented chapter codes are rejected by the
  gate.
- **State the nature and scope**: national versus regional, electoral programme
  versus political document, and any proposal the document conditions to another
  process.

## 5. The second gate

```bash
make check-run            # RUN=runs/29n2026 by default
```

`cmd/runcheck` proves, without network and without the source PDFs:

| Check | Proves | Issue code |
| --- | --- | --- |
| Quotes | Every «…» of the summary exists in the source | `runcheck.quote` |
| Page references | Every `[p. N]` line has something from that page | `runcheck.page` |
| Coverage (`numbers`) | Every proposal 1..N appears as an item of the summary | `runcheck.coverage` |
| Coverage (`items`) | The source item count matches and every item leaves a trace | `runcheck.coverage` |
| Claims | Each curated datum sits on the page it declares | `runcheck.claim` |
| Sections | Every `## N.M` of the summary exists in the source | `runcheck.section` |
| Embedding | Every document travels inside the view and no markup points at a file | `runcheck.embed` |

Exit code 0 means the run passes; 1 means failure with one line per issue. The
report prints how many elements each check examined, so a green run is also
evidence of coverage, not of an empty check.

### Exceptions are explicit

A check can legitimately fire (a literal that crosses a page break, a statement
about the file itself). Those go to `checks/allow.tsv`, four tab-separated
fields — code, file, a substring of the detail and **a mandatory reason**:

```
runcheck.quote	Resumen-PP-2026.md	en tres meses, con silencio positivo	Literal exacto que cruza el salto entre las páginas 41 y 42…
```

An empty reason is rejected. An exception never silences a different issue: the
match also requires the code and the file.

### What the gate cannot prove

It proves literalness, location, coverage and self-containment. It cannot prove
that the interpretation is right, that the party's figures are true, or that no
proposal of the same section was merged away. Reporting the counts and reading
the summary is still required.

## 6. The view

The view is a **single file**: `tools/template.html` plus
`tools/build_standalone.py`, which embeds the summaries, the verification
reports and the provenance document as HTML inside a JSON payload and renders
them in an overlay viewer.

- Nothing in the markup may reference a sibling file. External `http(s)` links
  are allowed because they are provenance, not content.
- Context boxes that explain a law, an institution or an event are **editorial**
  and must be labelled as such, in the page and in the footer. They never appear
  inside a «…» quote.
- The source PDFs are **not** embedded: base64 of 28 MB produces a 38 MB document
  nobody can read comfortably. They are linked by official URL and identified in
  `SOURCES.md`.
- What the gate requires embedded is the **content**: the summary and the
  verification report of every programme, plus the provenance document. The
  printable PDFs of the summaries are generated from that same text, so they are
  derivable artifacts and are not embedded.

Portable companions, both generated:

```bash
python3 tools/build_standalone.py     # the self-contained view
python3 tools/build_compendium.py     # A4 PDF: cover, index, summaries, reports
```

## 7. Failure modes seen in practice

Each of these happened in `runs/29n2026/` and each is now either a check or a
line in this document. Read them before writing a summary.

| Failure | Symptom | Now prevented by |
| --- | --- | --- |
| Paraphrase inside «» | "quien entre de forma ilegal" for "Quien entra" | `runcheck.quote` |
| Negation rewritten as quote | "sin anexos reservados" for "No habrá anexos reservados" | `runcheck.quote` |
| Digits for words | "de 0 a 3 años" for "de cero a tres años" | `runcheck.quote` |
| Dropped article in a quote | "socio, competidor y rival" for "un socio, un competidor…" | `runcheck.quote` |
| Literal split by a page break | page number printed between two halves of the quote | `checks/allow.tsv` with a reason |
| Wrong page in a citation | cites the exposition page instead of the proposal page | `runcheck.page` |
| Missing subapartado | four blocks shown, one of the 28 sections lost | `runcheck.section`, review |
| Card in the wrong block | infrastructures under "ambición" instead of "Reunir la Nación" | review against the source index |
| Invented section code | a `## 7.2` the document does not have | `runcheck.section` |
| Proposal silently dropped | coverage gap in a 249-item programme | `runcheck.coverage` |
| Coverage fooled by page numbers | "2" found in `[p. 2]`, not in an item | `summaryItem` only matches list items |
| View needing sibling files | the HTML opens but the buttons do nothing | `runcheck.embed` |
| Viewer without scroll | `overflow: hidden` on the body, long text unreachable | stated in `docs/design/site-quality.md` review; probe before shipping |
| Template drift in the generator | the payload merges into nothing after an edit | `runcheck.embed` fails; regenerate and re-check |

## 8. Adding a run

1. Create `runs/<run>/`, copy the source PDFs out of Git, and write `SOURCES.md`
   with URL, size and SHA-256.
2. Extract with `pdftotext -layout` plus `tools/pagemark.py`.
3. Write `run.json` with the programme list, the coverage mode and the
   `claims`/`allow`/`view` paths.
4. Write the summaries with page citations and literal quotes.
5. Add the curated claims to `checks/claims.tsv` (phrases must appear on the cited
   page).
6. Run `make check-run`, fix, and record every exception in `checks/allow.tsv`
   with its reason.
7. Generate the view and the compendium, then write the run `README.md` with the
   counts, the exceptions and the limits.
8. Close the Taskqueue task with the gate output as evidence.

## 9. Gate mapping

| Task gate | Evidence from this run |
| --- | --- |
| Acceptance criteria | `make check-run` output plus the counts in the run `README.md` |
| Focused behavior | `go test ./internal/runcheck/...` — one subtest per check and per failure mode |
| Regression | `make verify`, which now runs `check-run` before the documentation checks |
| Static and structural | `node --check` over the view script, headless render, `pdftotext` over the compendium |
| Documentation | this file, the run `README.md` and `SOURCES.md` |
| Trust boundary | provenance is verified but not certified; empty proposals reported; exceptions carry a reason |
| Artifact and task record | the run, the generator and the task record in the same commit |
