---
id: TQ-0068
title: Record the 29-N 2026 programme summaries, verification and preview
status: done
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/docs
  - area/project
  - area/validation
  - agent/developer
created: 2026-11-17T23:17:00+01:00
updated: 2026-11-17T23:17:17+01:00
---

## Outcome

The repository retains the 29-N 2026 run as reviewable evidence: three programme
conversions with per-page markers, three factual summaries with source pages,
three independent-check reports, a printable PDF per summary, a navigable HTML
view, the extraction/verification tooling, and the source provenance record.

## Scope

- Included: `runs/29n2026/` with text, page-marked Markdown, summaries, verification
  reports, PDFs, HTML view, tools and `SOURCES.md`; a committed task record.
- Excluded: promoting these programmes into the canonical public catalogue, editing
  `fixtures/`, or certifying the party figures. The run is a private draft.

## Acceptance criteria

- [ ] `runs/29n2026/` contains the three summaries, three verification reports, three
      PDFs, the HTML view, the three page-marked Markdown conversions, `SOURCES.md`
      and `tools/`.
- [ ] Every numbered proposal of the three documents is covered by its summary
      (133/133, 115/115, 249/249).
- [ ] The second-gate scripts report zero open discrepancies.
- [ ] `SOURCES.md` records size and SHA-256 for each source PDF.
- [ ] `make verify` passes with the run in place.
- [ ] The run index states the limitations, including that the verifying process
      also authored the summaries.

## Quality Gates

### Gate 1: Acceptance criteria

`ls runs/29n2026` and the counts table in `runs/29n2026/README.md`.
Pass: every artifact listed above exists and the counts match the documents.

### Gate 2: Focused behavior

`python3 runs/29n2026/tools/verify.py` and `python3 runs/29n2026/tools/datos.py`.
Pass: the report shows 0 quote failures and 0 page-attribution failures (one
justified cross-page literal accepted after manual reading).

### Gate 3: Regression

`make verify` from the repository root.
Pass: `verify: OK`, including documentation structure and relative links.

### Gate 4: Static and structural validation

`node --check` over the extracted `<script>` of the HTML view and a headless Chrome
`--dump-dom` render.
Pass: no syntax error; the DOM contains 7 sections and 30 measure cards for the
default tab.

### Gate 5: Documentation and specification consistency

`runs/29n2026/README.md` and `runs/29n2026/SOURCES.md`.
Pass: the index links resolve, the limitations section is present and `make verify`
reports `docs OK` and `links OK`.

### Gate 6: Trust-boundary and security behavior

Spot check that the eight content-free PP proposals are reported as such instead of
being filled in, and that no summary claim lacks a source page.
Pass: 8/8 reported in `Resumen-PP-2026.md`; unverifiable wording corrected during
the gate (12 corrections recorded).

### Gate 7: Artifact and task record

Artifacts are committed with this record; `git log -1 --stat` shows the run files and
this task file in the same commit.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/SPEC.md#publication-flow`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-17T23:17:17+01:00 — Gate evidence (all PASS):
  - G1 artifacts: runs/29n2026 contains text/ (3), frente-obrero.md, vox.md, pp.md, Resumen-*.md (3), Resumen-*.pdf (3), Verificacion-*.md (3), Programas-electorales-2026.html, SOURCES.md, README.md, tools/ (4). Counts 133/115/249 match the documents.
  - G2 verify.py: VOX 79 quotes / PP 82 / FO 61, 0 failures except one PP literal that crosses a page break (page number printed between 'en tres' and 'meses'), read manually on pp. 41-42 and correct. datos.py: 69 numeric claims, 0 failures.
  - G3 make verify: OK (build, vet, docs OK, links OK, all go test targets).
  - G4 node --check on the extracted script: OK; headless Chrome --dump-dom: 7 sections and 30 measure cards rendered for the default tab; tab simulation shows 30/39/37 cards with 106 contexts.
  - G5 README.md index plus SOURCES.md; limitations present; make verify reports docs OK and links OK.
  - G6 the 8 content-free PP proposals (83, 109, 132, 164, 167, 173, 200, 210) are reported as empty, not filled. 12 wording corrections applied during the gate are recorded in the three Verificacion-*.md files.
  - G7 this record is committed together with the run artifacts.

  Decisions: source PDFs stay outside the repository (immutable, documented by size and SHA-256 in SOURCES.md). The HTML view is kept as a run artifact, not as an approved presentation asset: its CSS/JS/markup were written by the agent under the owner's explicit instruction, which the project contract otherwise reserves to the maintainer. The verifying process also authored the summaries; that limitation is stated in README.md.
