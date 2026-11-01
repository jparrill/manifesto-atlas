---
id: TQ-0004
title: Write repository handoff document
status: done
priority: normal
assignee: gpt-6.1-sol
labels:
  - type/docs
  - area/project
  - agent/developer
created: 2026-11-01T16:45:25+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

A concise root-level handoff document lets a fresh agent continue Manifesto Atlas without repeating decisions already stored in the specification, quality policy, or Taskqueue history.

## Scope

- Summarize current repository state and open work.
- Reference existing artifacts rather than duplicating them.
- Include suggested skills for the next agent.
- Exclude secrets and personal information.

## Acceptance criteria

- [ ] `HANDOFF.md` exists at the repository root.
- [ ] It references the canonical specification, quality policy, agent instructions, tasks, commit, and remote repository.
- [ ] It explains the next actionable steps and Taskqueue working-directory requirement.
- [ ] It includes a suggested skills section.

## Quality Gates

### Gate 1: Acceptance criteria

- Verify every acceptance criterion is present in `HANDOFF.md`.

### Gate 2: Focused behavior

```bash
test -f HANDOFF.md
rg -n "docs/SPEC.md|docs/QUALITY_GATES.md|AGENTS.md|Suggested skills|github.com/jparrill/manifesto-atlas|tq" HANDOFF.md
```

Both commands must exit zero.

### Gate 3: Regression

```bash
tq list --json
tq show TQ-0004 --json
```

Both commands must exit zero with no malformed queue files.

### Gate 4: Static and structural validation

- Warning-or-higher diagnostics on `HANDOFF.md` must report zero findings.

### Gate 5: Documentation and specification consistency

- The handoff references existing artifacts instead of restating their complete content.

### Gate 6: Trust-boundary and security behavior

- Confirm no credentials, tokens, private email addresses, or absolute user-specific home paths appear.

### Gate 7: Artifact and task record

- Record all gate results with `tq note` before completion.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0004`:

1. Record all gate results with `tq note`, then run `tq done TQ-0004`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0004` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#further-notes`

---

## Notes

- 2026-11-01T16:45:56+01:00 — Quality gates:
  - Gate 1 PASS — HANDOFF.md covers the repository, authoritative artifacts, current state, next sequence, and suggested skills.
  - Gate 2 PASS — required artifact, repository, skill, and Taskqueue references are present.
  - Gate 3 PASS — tq list and tq show TQ-0004 exited zero with no malformed queue files.
  - Gate 4 PASS — warning-or-higher diagnostics returned zero findings.
  - Gate 5 PASS — the handoff references SPEC, quality policy, and agent instructions instead of duplicating their full contents.
  - Gate 6 PASS — no credentials, private emails, or absolute user-specific home paths are present.
  - Gate 7 PASS — results are recorded here and the handoff artifact is current.
