---
id: TQ-0001
title: Initialize project task management
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/developer
created: 2026-11-01T16:32:47+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

The repository has an initialized, versioned Taskqueue workflow that all supported agents can discover and follow.

## Scope

- Initialize Taskqueue at the repository root.
- Define project-specific task labels.
- Add shared agent instructions for creating and managing tasks.
- Connect Claude Code to the shared instructions and generated Taskqueue guide.

## Acceptance criteria

- [ ] `.taskqueue.yaml` and `.tasks/AGENTS.md` exist and are valid.
- [ ] The label vocabulary covers task types and Manifesto Atlas areas.
- [ ] `AGENTS.md` defines when and how tasks are declared.
- [ ] `CLAUDE.md` loads the shared project and Taskqueue instructions.
- [ ] Taskqueue commands work from the repository root.

## Verification

- `tq label list --json`
- `tq list --json`
- `tq ready --json`


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0001`:

1. Record all gate results with `tq note`, then run `tq done TQ-0001`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0001` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#further-notes`

---

## Notes

- 2026-11-01T16:33:18+01:00 — Initialized Taskqueue at the repository root, replaced the generic labels with Manifesto Atlas type/area vocabularies, and added AGENTS.md plus a CLAUDE.md include. Verified label, list, and ready commands and checked the new files with LSP diagnostics.
