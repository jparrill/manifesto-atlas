---
id: TQ-0002
title: Require quality gates for every task
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-01T16:35:23+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

Every new Taskqueue task declares objective quality gates and cannot be closed until all gates pass and their results are recorded.

## Scope

- Adapt the proven Rotris verification-gate workflow to Manifesto Atlas.
- Define project-wide and change-specific gate categories.
- Update the task declaration template and closing workflow.
- Keep the generated `.tasks/AGENTS.md` untouched.

## Acceptance criteria

- [ ] `AGENTS.md` requires a numbered `Quality Gates` section in every task.
- [ ] A project quality-gate document defines mandatory and conditional gates.
- [ ] Gates are binary, reproducible, and include commands or artifacts plus expected results.
- [ ] Agents must record results before moving a task to `done`.
- [ ] The policy works before and after the project Makefile exists.

## Quality Gates

### Gate 1: Acceptance criteria

- Verify every acceptance-criteria statement appears in `AGENTS.md` or `docs/QUALITY_GATES.md`.
- Passes when no criterion is missing or contradicted.

### Gate 2: Focused behavior

```bash
rg -n "Quality Gates|tq note|make verify|binary and reproducible|failed, flaky, skipped, or unavailable" AGENTS.md docs/QUALITY_GATES.md
```

Passes when the command finds the task template, result recording, future Makefile seam, reproducibility requirement, and blocking failure states.

### Gate 3: Regression

```bash
tq list --json
tq show TQ-0002 --json
tq ready --json
```

Passes when all commands exit zero, this task remains readable, and Taskqueue reports no malformed queue files.

### Gate 4: Static and structural validation

- Run warning-or-higher diagnostics on `AGENTS.md`, `CLAUDE.md`, `.taskqueue.yaml`, and `docs/QUALITY_GATES.md`.
- Passes with zero diagnostics.

### Gate 5: Documentation and specification consistency

- Confirm `AGENTS.md` points to the generated Taskqueue guide and the quality-gate policy without duplicating the generated CLI reference.
- Confirm terminology is consistently `Quality Gates`.

### Gate 6: Trust-boundary and security behavior

- N/A: this task changes agent workflow documentation and does not alter a runtime trust boundary.

### Gate 7: Artifact and task record

- Confirm no generated Taskqueue guide content was edited manually.
- Record all gate results with `tq note` before moving the task to `done`.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0002`:

1. Record all gate results with `tq note`, then run `tq done TQ-0002`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0002` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#testing-decisions`

---

## Notes

- 2026-11-01T16:37:16+01:00 — Quality gates:
  - Gate 1 PASS — every acceptance criterion is represented in AGENTS.md or docs/QUALITY_GATES.md.
  - Gate 2 PASS — policy marker search found the task template, tq note evidence requirement, make verify seam, reproducibility rule, and blocking failure states.
  - Gate 3 PASS — tq list, show, and ready exited zero with no malformed queue files.
  - Gate 4 PASS — warning-or-higher diagnostics returned zero findings for AGENTS.md, CLAUDE.md, .taskqueue.yaml, and docs/QUALITY_GATES.md.
  - Gate 5 PASS — AGENTS.md points to .tasks/AGENTS.md and uses Quality Gates consistently without copying the generated CLI guide.
  - Gate 6 PASS (N/A) — documentation-only workflow change; no runtime trust boundary changed.
  - Gate 7 PASS — generated .tasks/AGENTS.md was not manually edited and this evidence is recorded.
  Rotris review: 104 of 109 tasks contain a Verification Gate heading or variant; all 86 done and all 4 todo tasks have gates. The five tasks without gates are parked inbox ideas. Strong tasks use exact commands and pass conditions; gate-result notes are not consistently recorded, which this policy now requires.
