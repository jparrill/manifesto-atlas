---
id: TQ-0030
title: Extract pinned Docling page-aware JSON and derived Markdown
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/extraction
  - area/cli
depends_on:
  - TQ-0012
  - TQ-0021
  - TQ-0027
  - TQ-0028
  - TQ-0029
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T18:56:22+01:00
---

## Outcome

make extract produces canonical page-aware extraction from registered PDFs without unnecessary OCR.

## Scope

- Implement the approved pinned Docling/model/environment adapter inside enforced sandbox limits.
- Normalize page/block types, coordinates/order, tables, warnings, coverage, tool metadata, and derived intermediate Markdown.
- Disable OCR for native text by default; record escalation needs without silently approving suspicious pages.
- Exercise both actual golden PDFs and define the heavy local fixture gate separately from fast stored-output regression checks.
- Provide the focused Make gate `make test-extraction` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Both golden PDFs produce complete page-mapped structured extraction and derived Markdown with observed headings/numbers/names/negations.
- [ ] Native-text runs do not invoke OCR without explicit approved signals; page warnings and converter failures remain visible.
- [ ] Pinned repeated extraction satisfies the approved deterministic-content contract and model/tool changes invalidate downstream reports.
- [ ] Gate evidence records real runtime/peak resource observations; inability to run either required fixture blocks completion.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-extraction`. Pass when automated cases prove: both golden adapter outputs; all pages; known text anchors; OCR-off; repeated pinned run; warnings; tool change. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-extraction`. Pass when automated cases prove: both golden adapter outputs; all pages; known text anchors; OCR-off; repeated pinned run; warnings; tool change. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0030 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Run corrupt/over-budget PDFs, converter timeout and memory exhaustion, missing pinned assets, and write/network-denial probes; never overwrite source bytes. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0030`:

1. Record all gate results with `tq note`, then run `tq done TQ-0030`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0030` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#extraction-artifacts`
- `docs/SPEC.md#extraction-tests`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
