---
id: TQ-0007
title: Establish a truthful Makefile interface for current and planned stages
status: todo
priority: high
labels:
  - type/chore
  - agent/developer
  - area/project
  - area/cli
created: 2026-11-01T17:07:32+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

A fresh checkout exposes discoverable Make commands and real baseline checks, with a clear expansion path for every specified pipeline stage.

## Scope

- Inspect the complete Spec and newly created backlog; inventory submit, verify-source, add, extract, validate-extraction, normalize, validate-content, build, validate-site, publish-check, process, and serve.
- Provide help, check, and verify now using actual available documentation/Git/Taskqueue checks; reject malformed-queue warnings even if tq exits zero.
- List future stage targets and their owning task IDs. Invocation of an unavailable stage must fail non-zero with an actionable explanation, never return a fake success.
- Keep deterministic implementation in Go once introduced and orchestration in Make. Later implementation tasks replace the unavailable-stage failure and extend real gates.
- No CI workflow or pipeline implementation is included.

## Acceptance criteria

- [ ] make help lists all 12 Spec pipeline targets and distinguishes implemented checks from planned stages.
- [ ] make check and make verify perform meaningful current checks and exit zero on a valid checkout without requiring a nonexistent Go module.
- [ ] Whitespace errors and a malformed task fixture cause the applicable baseline gate to fail non-zero, including tq warning-only failures.
- [ ] Every unavailable pipeline target exits non-zero; invoking make without arguments is safe and does not publish or fetch anything.
- [ ] Technical usage documentation maps planned targets to backlog owners and explains how Go/model/browser checks will be integrated.

## Quality Gates

### Gate 1: Acceptance criteria

Run make help, make check, and make verify. In a disposable repository copy inject a whitespace error and a malformed task, then assert failure. Invoke each unavailable stage and assert non-zero with its owning task reference. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run make help, make check, and make verify. In a disposable repository copy inject a whitespace error and a malformed task, then assert failure. Invoke each unavailable stage and assert non-zero with its owning task reference.

### Gate 3: Regression

Run `make verify`, `git diff --check`, `tq list --json`, and `tq show TQ-0007 --json`. Require exit zero and no malformed/skipped-queue warnings on stderr; repeat regression after any gate fix.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Check shell quoting with a checkout path containing spaces, propagate command failures, and verify default/help/check targets perform no network retrieval, model invocation, or publication. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0007`:

1. Record all gate results with `tq note`, then run `tq done TQ-0007`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0007` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#module-seams`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`
