---
id: TQ-0045
title: Compare compatible regional programmes without merging their identities
status: todo
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/rendering
  - area/taxonomy
  - area/site
depends_on:
  - TQ-0044
  - TQ-0043
  - TQ-0036
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T18:56:23+01:00
---

## Outcome

Readers can inspect common measures and territorial differences for compatible same-party programmes while preserving source attribution.

## Scope

- Implement the approved compatibility rules and canonical comparison view data in Go for supplied comparison assets and Markdown output.
- Keep each geography/election/scope/revision and citation distinct; show all topic states without inferring ideological agreement or opposition from absence.
- Use synthetic schema-valid national/autonomous/regional fixtures when the two golden sources are not a comparable pair; do not fabricate real programmes.
- Do not author or edit presentation assets; request missing human-authored interfaces as a blocker.
- Provide the focused Make gate `make test-comparisons` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] A compatible regional fixture pair shows defined common content and differences with exact per-programme citations.
- [ ] Incompatible party/election/scope/taxonomy/revision pairs are refused or explicitly separated according to the approved contract.
- [ ] Comparison output is deterministic, Spanish, and retains omission versus non-applicability distinctions.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-comparisons`. Pass when automated cases prove: compatible pair; regional difference; identical measure; all topic states; each incompatibility class; deterministic attribution. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-comparisons`. Pass when automated cases prove: compatible pair; regional difference; identical measure; all topic states; each incompatibility class; deterministic attribution. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0045 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Reject cross-programme citation/report substitution and hostile geography/measure text; no aggregate political ranking or unsupported equivalence may be introduced. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0045`:

1. Record all gate results with `tq note`, then run `tq done TQ-0045`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0045` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#topic-taxonomy`
- `docs/SPEC.md#canonical-programme-model`
- `docs/SPEC.md#deterministic-rendering`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
