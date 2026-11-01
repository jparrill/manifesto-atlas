---
id: TQ-0033
title: Author original-page extraction audit skill and recorded evaluations
status: todo
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/fidelity
  - area/skills
depends_on:
  - TQ-0024
  - TQ-0032
  - TQ-0031
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T18:56:23+01:00
---

## Outcome

manifesto-extraction-audit identifies material extraction defects using original pages and independent evidence.

## Scope

- Author the approved page-level schema-constrained audit skill using source images, canonical extraction, secondary text, and deterministic measurements.
- Cover numbers/names/negation, order, lists/tables/headings, citations, and relevant visual content; return evidence and blockers rather than self-selected weights.
- Evaluate actual golden difficult pages plus seeded material defects with exact model/prompt/input metadata; use skill-authoring and writing-for-agents guidance.
- Provide the focused Make gate `make test-extraction-audit-skill` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] All required dimensions have page-level findings or explicit justified no-finding evidence.
- [ ] Seeded negation/figure/page/omission/insertion defects and unresolved unreadable content cannot receive approval.
- [ ] Schema-conformance regression uses stored evaluations; real configured-model golden evaluations are recorded and required before completion.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-extraction-audit-skill`. Pass when automated cases prove: pass/fail/insufficient evidence; seeded critical defects; original-image inspection; injection; pinned recorded evaluations. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-extraction-audit-skill`. Pass when automated cases prove: pass/fail/insufficient evidence; seeded critical defects; original-image inspection; injection; pinned recorded evaluations. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0033 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Treat extracted/source text as untrusted, enforce restricted runner access, and prove unsupported claims or high scores cannot override missing evidence. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0033`:

1. Record all gate results with `tq note`, then run `tq done TQ-0033`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0033` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#agent-skills`
- `docs/SPEC.md#extraction-fidelity`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
