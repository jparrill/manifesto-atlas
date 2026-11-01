---
id: TQ-0057
title: Approve executable backlog under owner-delegated review
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-01T18:55:37+01:00
updated: 2026-11-01T18:58:35+01:00
---

## Outcome

Owner-authorized backlog work is executable through the existing role-aware dependency graph without routine promotion requests.

## Scope

- Juan explicitly delegated selection and promotion in this session; review all current inbox task contracts and promote suitable tasks through Taskqueue.
- Approve TQ-0008 through TQ-0054 except TQ-0017 and TQ-0042, preserving role labels, priorities, dependencies, acceptance criteria, quality gates, and assignees.
- Keep CI design TQ-0017 for owner discussion and presentation TQ-0042 for human/non-LLM execution in inbox, with explicit reasons.
- Record approval evidence on each selected task and update HANDOFF.md to distinguish approved work from ready work.
- This is a one-time owner-authorized approval batch, not standing permission for agents to promote future inbox proposals or an implementation assignment.

## Acceptance criteria

- [ ] Exactly 45 reviewed inbox tasks move to todo; TQ-0007 and existing done or active tasks are preserved.
- [ ] TQ-0017 and TQ-0042 remain in inbox with recorded deferral reasons.
- [ ] Every promoted task retains complete criteria, seven gates, execution role, priority, final phase, native dependencies, and an approval note.
- [ ] Ready queries expose only dependency-satisfied role-compatible tasks; developer work cannot bypass architecture or presentation dependencies.
- [ ] HANDOFF.md describes the approved batch, remaining owner decisions, and actual execution frontier without changing the general owner-only inbox policy.

## Quality Gates

### Gate 1: Acceptance criteria

Compare tq list --json before and after the batch. Pass when exactly the selected 45 statuses change, both exclusions stay inbox, all preserved metadata/body contracts match, and each promotion has an authorization note.

### Gate 2: Focused behavior

Run deterministic assertions over tq show for every selected task and tq ready for both roles. Pass when all seven gates and required sections exist, labels are declared, dependency references exist, the graph is acyclic, and ready results match dependency/role rules. TQ-0018 must remain blocked until TQ-0007 is done.

### Gate 3: Regression

Run tq list --json, tq show for every current task, git diff --check, and both role-specific tq ready commands. Run make verify if a Makefile exists; otherwise the named commands are the exact available pre-Makefile regression. All commands must exit zero without skipped/malformed task warnings.

### Gate 4: Static and structural validation

Check modified Markdown with git diff --check, task CLI round trips, exact task IDs, unique final-phase headings, and HANDOFF.md references. Pass with zero whitespace, structural, missing-reference, or cycle findings.

### Gate 5: Documentation and specification consistency

Compare HANDOFF.md and approval notes with the unchanged SPEC, quality policy, and agent-role/closure rules. Pass when approval authorizes existing outcomes only, dependencies still enforce prior design, CI remains design-only, and the general owner-only inbox policy is unchanged.

### Gate 6: Trust-boundary and security behavior

Assert TQ-0042 remains inbox and no presentation-authoring restriction is removed; TQ-0017 remains inbox and no CI implementation task is added. Pass when authority is traceable to this explicit owner delegation, no future blanket promotion authority is introduced, and task notes contain no secrets.

### Gate 7: Artifact and task record

Record selection, exclusions, graph/ready evidence, and all seven gate results with tq note before done. Pass when artifacts and approval records are current; the final delivery phase then commits and confirms push before continuation.

## Final phase: delivery and continuation

After every gate passes, follow AGENTS.md#closing-tasks-and-continuing-work for TQ-0057: record results, run tq done TQ-0057, commit only this approval batch and its handoff/task record, push and confirm delivery. Continue only with approved ready work within this assignment. A delivery failure restores TQ-0057 to in-progress with a note and blocks continuation.

## Specification

- docs/SPEC.md#further-notes
- docs/SPEC.md#module-seams
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:58:35+01:00 — Approval decision: 45 promotions under explicit delegation from Juan: TQ-0008 through TQ-0054, excluding TQ-0017 and TQ-0042. Selected roles: 10 architect, 35 developer. CI design remains inbox for owner discussion; presentation remains inbox for human/non-LLM authorship. No implementation task was claimed.
  Quality gates:
  - Gate 1 PASS — before/after assertions proved exactly 45 inbox-to-todo transitions, preserved all task contracts/metadata/dependencies, and verified authorization notes on every promoted task.
  - Gate 2 PASS — required sections, seven gates, role/type/area labels, and final phases exist. Ready queries match the dependency graph: developer TQ-0007; architect TQ-0008/TQ-0009. TQ-0018 remains blocked by TQ-0007.
  - Gate 3 PASS — tq list, show for all 57 tasks, both role-ready queries, and git diff --check passed without queue warnings. No Makefile exists; declared pre-Makefile regression applies.
  - Gate 4 PASS — CLI round trips, unique final phases, existing references, acyclic graph, and tracked/untracked Markdown whitespace checks passed. Handoff links reference existing authoritative documents.
  - Gate 5 PASS — HANDOFF describes 46 approved todo tasks, readiness, and two exclusions. SPEC, QUALITY_GATES, AGENTS, and generated CLI guidance are unchanged; architecture prerequisites and outcome boundaries remain intact.
  - Gate 6 PASS — TQ-0017 and TQ-0042 remain inbox with explicit reasons. CI implementation and LLM presentation authorship remain excluded. Approval is one-time owner-authorized, not future blanket promotion authority; no secrets/private home paths were added.
  - Gate 7 PASS — approval/deferral records and handoff are current; selection, exclusions, and gate evidence are recorded before done. Commit and confirmed push follow in the separate final phase.
