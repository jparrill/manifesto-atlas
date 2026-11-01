---
id: TQ-0048
title: Author schema-constrained visual site QA after deterministic checks
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/site
  - area/skills
depends_on:
  - TQ-0047
  - TQ-0024
  - TQ-0016
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

manifesto-site-qa independently reviews current desktop/mobile screenshots without overriding deterministic failures.

## Scope

- Author the skill against approved canonical/report/screenshot contracts and restricted multimodal runner capabilities.
- Inspect hierarchy, readability, overflow, clipping, missing/misplaced sections, and obvious programme/comparison inconsistencies.
- Use supplied screenshot fixtures and real recorded configured-model evaluations; the model may report defects but never edit presentation code or approve a failed deterministic gate.
- Provide the focused Make gate `make test-visual-qa-skill` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Schema-valid pass/fail/insufficient reports cite exact viewports/screenshots and canonical/assets/browser/input hashes.
- [ ] Seeded visual defects and missing/stale screenshots fail or remain insufficient; failed deterministic checks always block approval.
- [ ] Recorded clean-context real-model evaluations exist for desktop/mobile views and regression uses their committed bounded fixtures.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-visual-qa-skill`. Pass when automated cases prove: good views; clipped/overflow/missing section; mobile failure; stale image; deterministic-fail/visual-pass conflict; recorded evaluations. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-visual-qa-skill`. Pass when automated cases prove: good views; clipped/overflow/missing section; mobile failure; stale image; deterministic-fail/visual-pass conflict; recorded evaluations. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0048 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Treat screenshot text and rendered source content as untrusted; deny edits, shell/network/secrets, and test instruction-bearing content and forged deterministic reports. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0048`:

1. Record all gate results with `tq note`, then run `tq done TQ-0048`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0048` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#agent-skills`
- `docs/SPEC.md#site-validation`
- `docs/SPEC.md#site-quality-tests`
- `docs/QUALITY_GATES.md`
