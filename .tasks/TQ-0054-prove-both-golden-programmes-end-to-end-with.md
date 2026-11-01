---
id: TQ-0054
title: Prove both golden programmes end to end with offline replay
status: todo
priority: high
labels:
  - type/chore
  - agent/developer
  - area/project
  - area/validation
  - area/extraction
  - area/site
depends_on:
  - TQ-0051
  - TQ-0028
  - TQ-0029
  - TQ-0045
  - TQ-0052
  - TQ-0053
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T18:56:24+01:00
---

## Outcome

Both approved real PDFs have complete auditable artifact graphs and deterministic public output, with regressions verifiable without rerunning models.

## Scope

- Execute real golden processing with approved provider/sandbox/tool versions and all provenance/extraction/content/site gates; retain exact input identities and signed reports.
- Implement fast offline integration fixtures and explicit heavy real-PDF/model/browser verification commands with reproducible resource budgets.
- Exercise deliberately failed/stale/mutated graphs to prove exclusion from catalogue and publication.
- This is local end-to-end verification, not CI or deployment implementation.
- Provide the focused Make gate `make test-golden-e2e` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Frente Obrero and VOX Castilla y León each complete verified-source -> approved extraction -> validated canonical -> deterministic Markdown/HTML -> site approval -> publish-check.
- [ ] make verify validates stored golden contracts/reports/outputs without live models or official websites; documented heavy gates actually run both real sources and required browser/model checks.
- [ ] Seeded provenance/extraction/claim/coverage/render/site failures and any changed recorded input prevent publication.
- [ ] Clean repeated builds reproduce output and all task-specific gates have recorded PASS evidence; unavailable required tools/models leave this task open.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-golden-e2e`. Pass when automated cases prove: both stored complete graphs; deterministic replay; every gate failure; stale dependency; immutable source mutation; no live-model call in offline verification. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-golden-e2e`. Pass when automated cases prove: both stored complete graphs; deterministic replay; every gate failure; stale dependency; immutable source mutation; no live-model call in offline verification. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0054 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Run controlled hostile-input/failure paths across actual trust boundaries; verify no sandbox escape, hidden publication fallback, leaked secrets, or unsigned substituted approval. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0054`:

1. Record all gate results with `tq note`, then run `tq done TQ-0054`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0054` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#end-to-end-tests`
- `docs/SPEC.md#testing-decisions`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:24+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
