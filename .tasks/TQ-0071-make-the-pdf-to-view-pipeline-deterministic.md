---
id: TQ-0071
title: Make the PDF-to-view pipeline deterministic from the CLI
status: done
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/feature
  - area/extraction
  - area/rendering
  - agent/developer
created: 2026-11-17T23:57:41+01:00
updated: 2026-11-17T23:59:23+01:00
---

## Outcome

A published PDF plus its official URL produce a run with one command each for
provenance and extraction, verification and rendering: `make extract`,
`make validate-content` and `make build`. The three are offline, deterministic and
covered by tests, including an end-to-end flow with fake `pdftotext`/`pdfinfo`.

## Scope

- Included: `internal/runsource` (provenance, extraction, page markers, sections,
  scaffold, manifest registration, provenance document), `internal/mdrender`
  (closed Markdown subset to HTML), `internal/runview` (single-file view with the
  documents embedded), `templates/view.html`, `internal/runmanifest` (shared
  manifest), the `extract`, `validate-content` and `build` stages in the CLI, the
  new gate checks (`runcheck.provenance`, `runcheck.scaffold`) and the report
  written by `validate-content`.
- Excluded: writing summaries (semantic work), calling models, downloading the PDF
  from the URL, and the publication stages (`normalize`, `publish-check`, …), which
  stay declared and unavailable.

## Acceptance criteria

- [ ] `extract` records URL, size, SHA-256 and page count, and refuses a non-https
      origin or a page-count mismatch.
- [ ] `extract` writes the raw text, the page-marked extraction, `SOURCES.md`,
      `run.json`, empty `checks/` files and a summary skeleton with placeholders.
- [ ] The gate refuses a run whose summary still has placeholders or whose source
      has no verifiable origin.
- [ ] `validate-content` writes `Verificacion.md` and fails with one line per issue.
- [ ] `build` writes one HTML with every document embedded, records it in the
      manifest, and refuses to overwrite a hand-made view without `-force`.
- [ ] `make verify` passes with the pipeline tests.

## Quality Gates

### Gate 1: Acceptance criteria

`go test ./internal/cli/...` runs the whole flow with fake tools and asserts every
artifact. Pass: all subtests pass.

### Gate 2: Focused behavior

`go test ./internal/runsource/... ./internal/mdrender/... ./internal/runview/...`
Pass: provenance, page markers, sections, scaffold, Markdown subset and escaping,
and `Register` idempotency.

### Gate 3: Regression

`make verify`
Pass: `verify: OK`, with the run gate and the pipeline tests in the output.

### Gate 4: Static and structural validation

`gofmt -l`, `go vet ./...`, and a real render of the generated view in headless
Chrome. Pass: no findings; the DOM contains the tabs and the payload.

### Gate 5: Documentation

`docs/design/programme-run.md#8` and `docs/MAKEFILE.md`.
Pass: both describe the same three commands, their guarantees and their limits;
`make verify` reports `docs OK` and `links OK`.

### Gate 6: Trust-boundary

The origin must be https; the gate rejects a run without provenance; `build` cannot
overwrite a hand-made view; the Markdown renderer escapes HTML and drops local
links. Pass: subtests for each case.

### Gate 7: Artifact and task record

Committed with the task record, gate output in the note.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/design/programme-run.md`

---

## Notes

- 2026-11-17T23:59:23+01:00 — Gate evidence (all PASS):
  - G1: go test ./internal/cli/... -> TestRunStagesEndToEnd exercises the whole flow with fake pdftotext/pdfinfo: extract writes run.json, programa.md with '<!-- page 2 -->', text/programa.txt, Resumen-PP.md skeleton and SOURCES.md; validate-content fails with runcheck.scaffold while the skeleton is unfilled; after writing a cited summary it passes and writes Verificacion.md; build writes the embedded view and the gate then reports 'embedding'; extract without -url exits 2.
  - G2: go test ./internal/runsource/... ./internal/mdrender/... -> provenance (https, bytes, sha256, pages), page-count mismatch, empty-page placeholder, section detection, scaffold content, Register idempotency, Markdown subset and HTML escaping, local links not rendered as links.
  - G3: make verify -> verify: OK with 'verify: programme run gate' and 'verify: go test (run pipeline)'.
  - G4: gofmt -l empty, go vet ./... clean. Real render of the generated view (Chrome --dump-dom): tabs Frente Obrero/VOX/Partido Popular, 23 sections, 153 cards, buttons 'Resumen completo', 'Informe de verificación' and 'Fuentes y procedencia'. Probe: the viewer opens the real report (5188 chars, 30 table rows, scrollable).
  - G5: docs/design/programme-run.md#8 documents the three commands, and docs/MAKEFILE.md lists extract/validate-content/build; make verify reports docs OK and links OK.
  - G6: non-https origin refused (unit test), gate rejects a run without provenance (unit test), build refuses to overwrite a hand-made view without -force, mdrender escapes HTML and drops sibling links.
  - G7: committed with this record.

  Decisions: the manifest definition moved to internal/runmanifest so the gate and the builders share one source of truth. The view payload keeps the shape the embedding check already validated (programme id -> documents, plus sources and a chrome block). validate-content writes a single Verificacion.md for generated runs and the view also accepts the hand-named Verificacion-<summary>.md of the 29-N run. extract never downloads the PDF: the origin is recorded and checked separately.
