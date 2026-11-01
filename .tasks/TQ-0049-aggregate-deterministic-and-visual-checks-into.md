---
id: TQ-0049
title: Aggregate deterministic and visual checks into the site validation gate
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/site
  - area/validation
  - area/cli
depends_on:
  - TQ-0046
  - TQ-0047
  - TQ-0048
  - TQ-0020
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

make validate-site produces one current authenticated approval only when every required site check passes.

## Scope

- Orchestrate parity/structure/links, accessibility/keyboard, responsive/browser/screenshot regression, and independent visual QA in the specified order.
- Persist signed stage reports with exact input graphs and retained failures; provide offline replay of stored results without model invocations.
- No deterministic failure can be waived by a visual pass, threshold averaging, missing tool, or manual fallback.
- Provide the focused Make gate `make test-site-gate` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] A current complete passing report set yields site approval for the exact canonical programme/rendered assets/views.
- [ ] Every individual failure, absent tool/report/image, stale configuration/input, bad signature, or insufficient visual evidence blocks approval.
- [ ] Repeated unchanged checks are idempotent and asset/browser/viewport/baseline changes invalidate the correct approvals.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-site-gate`. Pass when automated cases prove: all pass; each failed/missing component; deterministic fail plus visual pass; stale screenshot/config; forged report; offline replay. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-site-gate`. Pass when automated cases prove: all pass; each failed/missing component; deterministic fail plus visual pass; stale screenshot/config; forged report; offline replay. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0049 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Test report swapping, omitted viewport checks, malicious artifact paths, forged visual approval, and signing-role confusion; no partial report set is publishable. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0049`:

1. Record all gate results with `tq note`, then run `tq done TQ-0049`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0049` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#site-validation`
- `docs/SPEC.md#publication-and-ci`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`
