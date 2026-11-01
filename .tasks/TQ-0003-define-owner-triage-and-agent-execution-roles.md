---
id: TQ-0003
title: Define owner triage and agent execution roles
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-01T16:40:35+01:00
updated: 2026-11-01T16:43:02+01:00
---

## Outcome

Taskqueue clearly separates owner-only inbox triage from executable todo work and restricts architecture tasks to explicitly approved high-capability models.

## Scope

- Define owner-only inbox review and promotion.
- Define developer and architect execution roles.
- Add role labels to the Taskqueue vocabulary.
- Define claim commands and model attribution.
- Define how agents file discovered follow-up work.

## Acceptance criteria

- [ ] Only Juan may move tasks from `inbox` to `todo`.
- [ ] Every executable task has exactly one `agent/*` role label.
- [ ] Developer agents only claim `agent/developer` tasks.
- [ ] Architect tasks require explicit owner approval of the executing model.
- [ ] Claimed tasks record the actual model or agent as assignee.
- [ ] Agent-created follow-up work remains in `inbox` until owner review.

## Quality Gates

### Gate 1: Acceptance criteria

- Verify every criterion is stated in `AGENTS.md` and supported by the configured labels.

### Gate 2: Focused behavior

```bash
tq label list --json
tq ready --label agent/developer --json
tq ready --label agent/architect --json
```

All commands must exit zero and both role labels must be configured.

### Gate 3: Regression

```bash
tq list --json
tq show TQ-0003 --json
```

Both commands must exit zero with no malformed queue files.

### Gate 4: Static and structural validation

- Warning-or-higher diagnostics must report zero findings for changed Markdown and YAML files.

### Gate 5: Documentation and specification consistency

- `AGENTS.md` must use Taskqueue's actual `inbox`, `todo`, label, assignee, and ready-filter behavior.
- `.tasks/AGENTS.md` remains generated and unmodified.

### Gate 6: Trust-boundary and security behavior

- N/A: this task changes development workflow policy rather than a runtime trust boundary.

### Gate 7: Artifact and task record

- Run `tq init` after changing `.taskqueue.yaml`.
- Record all gate results with `tq note` before completion.

## Specification

- `docs/SPEC.md#agent-skills`

---

## Notes

- 2026-11-01T16:42:10+01:00 — Corrected the canonical architect model and task assignees from gpt-5.6-sol to gpt-6.1-sol after owner feedback.
- 2026-11-01T16:43:02+01:00 — Quality gates:
  - Gate 1 PASS — owner-only inbox promotion, exact execution roles, architect approval, assignee tracking, and inbox follow-up rules are present in AGENTS.md.
  - Gate 2 PASS — role labels are configured and both role-filtered tq ready commands exit zero.
  - Gate 3 PASS — tq list and tq show exit zero with no malformed queue files; all three current tasks have exactly one agent role.
  - Gate 4 PASS — warning-or-higher diagnostics returned zero findings for changed Markdown and YAML files.
  - Gate 5 PASS — policy uses actual Taskqueue fields and .tasks/AGENTS.md remains generated.
  - Gate 6 PASS (N/A) — workflow documentation only; no runtime trust boundary changed.
  - Gate 7 PASS — tq init was run after configuration changes and results are recorded here.
  Canonical architect identifier and all current assignees are gpt-6.1-sol.
