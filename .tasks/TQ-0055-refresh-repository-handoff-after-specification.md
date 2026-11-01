---
id: TQ-0055
title: Refresh repository handoff after specification backlog planning
status: todo
priority: normal
assignee: gpt-6.1-sol
labels:
  - type/docs
  - area/project
  - agent/developer
depends_on:
  - TQ-0056
created: 2026-11-01T17:13:37+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

HANDOFF.md lets a fresh agent continue from the actual planned backlog and current execution frontier without treating inbox proposals as approved implementation work.

## Scope

- Update the existing root handoff requested by Juan; reference task records instead of duplicating their contracts.
- Record current Git state, completed planning/ignore work, ready Makefile task, inbox backlog, dependency/coverage evidence, CI-design-only boundary, missing golden PDFs, and non-LLM presentation restriction.
- Replace unavailable suggested skill names with skills available in this environment.

## Acceptance criteria

- [ ] HANDOFF.md reflects current Taskqueue states and points to TQ-0005 coverage notes, TQ-0007 next executable work, and TQ-0017 CI design in inbox.
- [ ] It summarizes all backlog groups and native dependencies without promoting inbox work.
- [ ] It correctly describes .gitignore, missing implementation/PDF artifacts, current baseline commit, and uncommitted session work.
- [ ] It names available suggested skills, authoritative artifacts, and the safe Taskqueue claim/close workflow.
- [ ] It explicitly records TQ-0042 human/non-LLM presentation execution and contains no secrets or private absolute paths.

## Quality Gates

### Gate 1: Acceptance criteria

Compare HANDOFF.md with tq list --json, tq ready --label agent/developer --json, git status --short --branch, and git log --oneline -1. Pass when counts, task states, ready frontier, and baseline agree, accounting explicitly for this handoff task being completed after verification.

### Gate 2: Focused behavior

Assert required references and all relative file links resolve, each mentioned TQ ID exists, and old unavailable skill names are absent. Run documented non-mutating tq list/ready/show examples. Pass with zero missing references or inconsistent statements.

### Gate 3: Regression

Run tq list --json, tq show for TQ-0005/TQ-0006/TQ-0007/TQ-0017/TQ-0042 and this task, plus git diff --check. Run make verify if a Makefile exists; otherwise these exact pre-Makefile checks apply. All must exit zero with empty queue-warning stderr.

### Gate 4: Static and structural validation

Run git diff --no-index --check /dev/null HANDOFF.md and inspect Markdown headings, fenced blocks, and relative links. Pass with no whitespace or structural findings; a no-index difference exit alone is not an error when check output is empty.

### Gate 5: Documentation and specification consistency

Compare authoritative links and pipeline/rules summaries against SPEC.md, QUALITY_GATES.md, AGENTS.md, .tasks/AGENTS.md, and actual task bodies. Pass when no contract is contradicted and the handoff references details rather than duplicating them.

### Gate 6: Trust-boundary and security behavior

Check for private absolute paths, credentials, private keys, and misleading approval authority. Pass when none exist, inbox remains owner-only, CI remains design-only, and the presentation-code prohibition is explicit.

### Gate 7: Artifact and task record

Confirm HANDOFF.md and .gitignore are current, the coverage/graph audit remains valid, and all gate results are recorded with tq note before completion. Pass when no required check is skipped or unresolved.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0055`:

1. Record all gate results with `tq note`, then run `tq done TQ-0055`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0055` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- docs/SPEC.md#further-notes
- docs/SPEC.md#module-seams
- docs/SPEC.md#publication-and-ci
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T17:17:52+01:00 — Owner requested an explicit automatic done/commit/push/continuation policy while this refresh was in progress. Final handoff delivery depends on TQ-0056 so the handoff can reference the authoritative policy and actual pushed planning/ignore/workflow state.
