---
id: TQ-0073
title: Match the hand-made view and close the false positives the first real run exposed
status: in-progress
priority: urgent
assignee: claude-sonnet-4.6[1m]
labels:
  - type/bug
  - area/rendering
  - area/validation
  - agent/developer
created: 2026-11-18T00:35:07+01:00
updated: 2026-11-18T00:35:14+01:00
---

## Outcome

The generated view looks and behaves like the hand-made reference, and the gate no
longer passes a run whose coverage target was undercounted or whose view was never
rendered.

## Scope

- Included: the coverage undercount (2-digit regex hid proposals 100..115), the
  "summary reaches further than the declared target" check, validating after build
  inside `process`, an honest embedding line when no view exists, the generated
  template rebuilt on the reference design, `scripts/view-parity.mjs`, the
  `view-parity` target and its documentation.
- Excluded: hand-authored context boxes in generated runs (they are editorial), and
  committing the owner's scratch run.

## Acceptance criteria

- [ ] A source with 115 numbered proposals records 115, and a summary that numbers
      past the declared target fails the gate.
- [ ] `process` checks the embedding after building the view.
- [ ] The generated view uses the reference design and passes parity with zero style
      differences and every behaviour check.
- [ ] The template script is validated in the test suite so a stray backtick cannot
      ship again.

## Quality Gates

### Gate 1: Acceptance criteria

`make process` on a real PDF plus `make validate-content`.
Pass: proposals = 115 for the VOX document; the undercount case fails with
`runcheck.coverage`.

### Gate 2: Focused behavior

`go test ./internal/runcheck/... ./internal/runview/... ./internal/cli/...`
Pass: the undercount subtest, the section/label parsing, Build payload/refusal and
`node --check` of the template script.

### Gate 3: Regression

`make verify` → `verify: OK`.

### Gate 4: Static and structural validation

`gofmt -l`, `go vet ./...`, `make view-parity`.
Pass: no findings; parity reports 0 differences and all behaviour checks ok.

### Gate 5: Documentation

`docs/design/programme-run.md#6` documents `make view-parity`.
Pass: `make verify` reports `docs OK` and `links OK`.

### Gate 6: Trust-boundary

The gate still refuses a summary the model wrote if the target is wrong or the
view is missing; `build` still refuses to overwrite a hand-made view.
Pass: subtests and the parity run.

### Gate 7: Artifact and task record

Committed with the note holding the evidence.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/design/programme-run.md`

---

## Notes

- 2026-11-18T00:35:14+01:00 — Gate evidence (all PASS):
  - G1: re-running extract on the VOX PDF now reports 'proposals 115 numbered items in the source' (was 99 because the counter regex accepted only two digits). Sanity check: validate-content on a run whose summary reaches proposal 115 while the manifest declares 99 now reports 'the summary reaches proposal 115 but the run declares 99: the target was undercounted and the rest is unverified' and exits 1.
  - G2: go test ./internal/runcheck ./internal/runview ./internal/cli ./internal/summarize -> ok. New subtests: undercount detection, section title/page parsing for '[p. 7]' and '— pp. 23-25', card label taken from the literal quote, Build payload contents, refusal to overwrite a hand-made view, template payload placeholder, and 'node --check' over the template script (this is the guard against the stray backtick that broke the previous template).
  - G3: make verify -> verify: OK.
  - G4: gofmt -l empty, go vet clean. make view-parity CURATED=runs/29n2026/Programas-electorales-2026.html GENERATED=runs/prueba-total-yatusabes/Programas-electorales.html -> 'diferencias de estilo: 0' plus 7/7 behaviour checks ok (tabs, open all, close all, document button, viewer opens, viewer scrolls, Escape closes).
  - G5: docs/design/programme-run.md gains the parity section and the Makefile target; make verify reports docs OK and links OK.
  - G6: the gate keeps refusing undercounted targets and missing views; build still refuses to overwrite a hand-made file without -force (subtests).
  - G7: committed with this record.

  Decisions: the reference CSS is reused verbatim instead of re-implementing it, so the generated view cannot drift from the owner's design; parity is measured with computed styles and user-visible behaviour, not by eye. The owner's scratch run runs/prueba-total-yatusabes is left uncommitted on purpose.
