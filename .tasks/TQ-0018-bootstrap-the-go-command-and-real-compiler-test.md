---
id: TQ-0018
title: Bootstrap the Go command and real compiler/test gates
status: inbox
priority: high
labels:
  - type/chore
  - agent/developer
  - area/project
  - area/cli
depends_on:
  - TQ-0007
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:30+01:00
---

## Outcome

The repository builds a Go entry point with predictable CLI behavior and enforced Go quality checks.

## Scope

- Add the Go module and smallest useful command dispatcher; do not pre-create empty module layers or pretend pipeline stages work.
- Pin/document the supported Go version and dependencies; implement safe help and explicit unknown/unavailable-command failure.
- Extend make check with formatting/static analysis and make verify with actual Go tests while retaining baseline gates.
- Provide the focused Make gate `make test-cli` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] A clean checkout builds the Go command and displays documented help.
- [ ] Unknown commands, malformed flags, and missing required arguments exit non-zero without creating artifacts.
- [ ] make check rejects unformatted/invalid Go and make verify runs the real suite; stage implementations can register commands without changing user-facing Make names.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-cli`. Pass when automated cases prove: help/version behavior; unknown command; invalid flag; path with spaces; propagated failure. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-cli`. Pass when automated cases prove: help/version behavior; unknown command; invalid flag; path with spaces; propagated failure. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0018 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Test hostile flag/path values are treated as data, no shell interpolation, and help/default operation cannot fetch, modify archive content, or publish. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0018`:

1. Record all gate results with `tq note`, then run `tq done TQ-0018`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0018` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#module-seams`
- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#testing-decisions`
- `docs/QUALITY_GATES.md`
