---
id: TQ-0032
title: 'Superseded: Measure extraction fidelity and detect deterministic blockers'
status: rejected
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/fidelity
depends_on:
  - TQ-0030
  - TQ-0031
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## MVP disposition

Superseded as a separate platform task by TQ-0030 deterministic page/warning blockers and TQ-0031/TQ-0039/TQ-0040 original-evidence/claim/coverage checks. Numeric formula and separate signed extraction approval are outside MVP; retained critical negative cases are explicitly owned by the replacement tasks. This is cancellation of the old scope, not a successful implementation.

## Outcome

Page-level extraction defects are measurable and blocking failures survive any aggregate score.

## Scope

- Compute the approved page coverage, text preservation, figures/entities/negations, reading order, structure, page alignment, and relevant visual coverage observations.
- Compare secondary evidence and preserve findings with page/block citations; distinguish deterministic findings from semantic judgment still required.
- Provide controlled mutations of golden artifacts for missing pages, changed numbers/negation, wrong mapping, unsupported insertions, reordered content, and unreadable regions.
- Do not invent fidelity weights or approve extraction before calibrated formula and semantic audit exist.
- Provide the focused Make gate `make test-fidelity-measurements` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Each controlled critical mutation produces a page-specific blocker, including a local error in otherwise high-quality output.
- [ ] Measurements are deterministic and bind exact source/extraction/secondary/image/tool inputs.
- [ ] Reports expose all required dimensions and unresolved evidence rather than claiming literal PDF-to-Markdown identity.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-fidelity-measurements`. Pass when automated cases prove: all golden pages; each critical mutation; deterministic rerun; secondary mismatch; unresolved visual region. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-fidelity-measurements`. Pass when automated cases prove: all golden pages; each critical mutation; deterministic rerun; secondary mismatch; unresolved visual region. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0032 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Bound adversarial text/table sizes and reject corrupted page maps, unknown blocks, or incomplete evidence; low aggregate error cannot erase a blocker. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0032`:

1. Record all gate results with `tq note`, then run `tq done TQ-0032`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0032` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#extraction-fidelity`
- `docs/SPEC.md#extraction-audit-tests`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rejected/superseded under Juan's explicit TQ-0061 replan. Superseded as a separate platform task by TQ-0030 deterministic page/warning blockers and TQ-0031/TQ-0039/TQ-0040 original-evidence/claim/coverage checks. Numeric formula and separate signed extraction approval are outside MVP; retained critical negative cases are explicitly owned by the replacement tasks. Historical acceptance/gate text and notes are preserved; no old gate is declared passed.
