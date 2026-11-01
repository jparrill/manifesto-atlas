---
id: TQ-0017
title: Design CI verification and GitHub Pages publication before implementation
status: inbox
priority: normal
labels:
  - type/docs
  - agent/architect
  - area/ci
  - area/validation
depends_on:
  - TQ-0008
  - TQ-0009
  - TQ-0016
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

Juan can review a complete CI and deployment design before any workflow implementation is proposed.

## Scope

- Design PR versus trusted-main execution, source immutability checks, offline committed-report/signature/hash verification, deterministic builds, cache invalidation, permissions, dependency pinning, and Pages deployment.
- Define fork/untrusted-PDF handling, signed-report trust provisioning, environment protections, concurrency, rollback, artifact retention, and cost/resource budgets.
- CI must invoke the same Make gates as local execution and must not rerun expensive LLM jobs.
- Keep this task in inbox pending Juan's review. Do not create workflows, deployment configuration, or CI implementation tasks in this session.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Design includes a pipeline/job diagram, permission matrix, trust-boundary analysis, exact local-target mapping, and failure/rollback scenarios.
- [ ] Stale reports, modified committed PDFs, forged approvals, and failed local gates stop deployment.
- [ ] No workflow is implemented; downstream CI work requires separate owner-approved tasks after this design.
- [ ] docs/design/ci-design.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/ci-design.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/ci-design.md. Pass only when each of these cases has one explicit, internally consistent expected result: fork PR; main merge; changed PDF; stale report; bad signer; poisoned cache; failed build; cancelled deployment; rollback. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0017 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0017`:

1. Record all gate results with `tq note`, then run `tq done TQ-0017`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0017` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#publication-and-ci`
- `docs/SPEC.md#untrusted-document-handling`
- `docs/SPEC.md#end-to-end-tests`
- `docs/QUALITY_GATES.md`
