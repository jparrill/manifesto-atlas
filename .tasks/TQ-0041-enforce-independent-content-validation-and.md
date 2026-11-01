---
id: TQ-0041
title: Enforce independent content validation and draft lifecycle
status: inbox
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/normalization
  - area/cli
depends_on:
  - TQ-0038
  - TQ-0039
  - TQ-0040
  - TQ-0020
  - TQ-0035
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:30+01:00
---

## Outcome

make validate-content advances a draft only when independent claim and complete-source coverage reports both pass.

## Scope

- Implement draft -> validating -> validated transitions with input-bound independent jobs, signed reports, retained failures, and deterministic offline replay.
- Require current provenance/extraction/canonical/taxonomy/schema identities and both semantic approvals.
- Different models are preferred; one model is permitted only through the approved isolated invocation policy. No human-review bypass exists.
- Provide the focused Make gate `make test-content-gate` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Two matching complete passing reports advance the exact draft to validated.
- [ ] Fail, disagreement, insufficient evidence, missing report, reused generator context, bad signature, and any stale input keep content unpublishable.
- [ ] Every source/extraction/draft/schema/taxonomy/prompt/model/tool mutation invalidates relevant approvals; failed attempts remain auditable.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-content-gate`. Pass when automated cases prove: pass/pass; pass/fail; disagreement; insufficient; missing/stale/forged report; one-model isolation; lifecycle rollback. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-content-gate`. Pass when automated cases prove: pass/pass; pass/fail; disagreement; insufficient; missing/stale/forged report; one-model isolation; lifecycle rollback. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0041 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Test payload state forgery, swapping another draft's report, replay after input change, unauthorized signing, and manual approval injection; all must fail closed. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0041`:

1. Record all gate results with `tq note`, then run `tq done TQ-0041`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0041` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#validator-contract-tests`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`
