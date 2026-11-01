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
updated: 2026-11-01T16:41:59+01:00
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

## Specification

- `docs/SPEC.md#further-notes`

---

## Notes

- 2026-11-01T16:33:18+01:00 — Initialized Taskqueue at the repository root, replaced the generic labels with Manifesto Atlas type/area vocabularies, and added AGENTS.md plus a CLAUDE.md include. Verified label, list, and ready commands and checked the new files with LSP diagnostics.
