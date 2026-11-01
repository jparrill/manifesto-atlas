---
id: TQ-0056
title: Require automatic per-task commit push and approved-work continuation
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-01T17:17:33+01:00
updated: 2026-11-01T18:42:08+01:00
---

## Outcome

Every task has an explicit automatic delivery phase: successful gates, done status, task-scoped commit, confirmed push, and continuation to approved ready work within the assigned scope.

## Scope

- Make AGENTS.md the authoritative closure and continuation contract; align QUALITY_GATES.md and HANDOFF.md references.
- Add a final delivery phase to every current task through the Taskqueue CLI, preserving metadata, acceptance criteria, and historical notes.
- Specify failed commit/push recovery, exact remote confirmation, safe staging, existing hooks, and context/role/scope restrictions.
- Deliver the current planning and requested gitignore/handoff work with task-scoped commits and confirmed pushes; historical completed work is not reimplemented.

## Acceptance criteria

- [ ] AGENTS.md explicitly authorizes automatic commit/push after all gates pass and defines the complete ordered loop without another routine approval prompt.
- [ ] Every Taskqueue task contains a final phase referencing that contract and its own task ID.
- [ ] A failed commit or push prevents selecting the next task and restores the task to in-progress with a blocker note; normal hooks and non-force Git behavior remain required.
- [ ] Automatic continuation considers only approved ready todo work matching execution role and current assigned scope, never inbox or prohibited presentation authorship.
- [ ] Git history and the confirmed remote commit provide delivery evidence without a recursive post-push task-note commit.
- [ ] The procedure applies to the current planning, gitignore, workflow, and handoff deliveries using intended files only; each delivery is confirmed before selecting more work.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect the closure contract, all task final phases, and the task-scoped delivery history. Pass when each acceptance criterion has matching document, lifecycle, or remote evidence; completion of the separate handoff task remains its own delivery.

### Gate 2: Focused behavior

Assert through tq list/show that every task contains exactly one final delivery heading with its own ID and preserved dependencies/roles; new task records retain all seven gates, and historical verification records remain unchanged. Inspect scenario outcomes for gate failure, commit rejection, push rejection, successful push, empty ready queue, role mismatch, inbox-only work, and out-of-scope ready work. Pass when only confirmed delivery permits the next eligible claim.

### Gate 3: Regression

Run tq list --json, tq ready for both execution roles, tq show for every task, and git diff --check. Run make verify if a Makefile exists; otherwise the named commands are the exact available regression. Pass with exit zero and no malformed/skipped queue warnings.

### Gate 4: Static and structural validation

Check changed tracked and untracked Markdown with git diff --check or git diff --no-index --check, plus relative document links, required headings, declared labels, existing dependency references, and acyclic graph assertions. Pass with no whitespace/structural findings or unresolved task IDs.

### Gate 5: Documentation and specification consistency

Compare AGENTS.md, QUALITY_GATES.md, HANDOFF.md, and task final phases. Pass when all reference one authoritative procedure, preserve mandatory gates and owner-only triage, and do not duplicate or edit generated .tasks/AGENTS.md.

### Gate 6: Trust-boundary and security behavior

Inspect staging and failure scenarios. Pass when secrets/unrelated work are excluded, hooks remain enabled, force-push/amend/history rewriting is absent, unknown upstream/credentials are blockers rather than reasons to weaken checks, and inbox/presentation restrictions remain enforced.

### Gate 7: Artifact and task record

Verify current policy/task artifacts and record all seven pre-delivery gate results with tq note before marking done. Pass when the task record and artifacts are complete; the separate final phase then commits, pushes, and confirms delivery. Git history and remote confirmation provide delivery evidence; any delivery failure reopens this task and blocks continuation.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0056`:

1. Record all gate results with `tq note`, then run `tq done TQ-0056`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0056` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- docs/SPEC.md#testing-decisions
- docs/SPEC.md#publication-and-ci
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:42:08+01:00 — Quality gates:
  - Gate 1 PASS — AGENTS.md defines the ordered automatic gate/done/commit/push/confirmed-delivery/next-task loop. Planning dd313cd and gitignore 74f0c8b are already pushed; handoff remains its separately tracked dependent delivery.
  - Gate 2 PASS — deterministic audit reread all 56 task records; each has exactly one final phase with its own task ID. Existing metadata, dependencies, notes, and seven gates on new tasks are preserved. Documented failure/success/role/inbox/scope scenarios permit continuation only after confirmed delivery.
  - Gate 3 PASS — tq list/show for every task and ready queries for both roles exited zero with empty stderr; git diff --check passed. No Makefile exists, so the exact pre-Makefile regression applies.
  - Gate 4 PASS — Taskqueue round trips, exact-ID/heading assertions, existing-reference and acyclic-graph checks passed. Markdown changes and handoff no-index whitespace checks have zero findings. Relative final-phase links resolve to the existing AGENTS.md closure heading.
  - Gate 5 PASS — AGENTS.md is authoritative; QUALITY_GATES.md, HANDOFF.md, and all final phases point to it. Owner-only inbox review and mandatory gates remain intact; generated .tasks/AGENTS.md is unchanged.
  - Gate 6 PASS — procedure requires intended-file staging, normal hooks/signing, existing upstream, non-force push, failure reopening, and no history rewriting. Role/scope/presentation restrictions remain explicit. Changed policy/task content contains no secrets or private home paths.
  - Gate 7 PASS — current policy artifacts and every task final phase are complete, all pre-delivery evidence is recorded here, and remote confirmation is performed after this task-scoped commit without recursively editing the task to store its own hash.
