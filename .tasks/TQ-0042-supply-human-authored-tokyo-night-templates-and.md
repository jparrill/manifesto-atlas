---
id: TQ-0042
title: Supply human-authored Tokyo Night templates and presentation assets
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/rendering
  - area/site
depends_on:
  - TQ-0016
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-05T00:20:21+01:00
---

## Outcome

Approved human/non-LLM presentation assets make the Spanish programme and comparison layouts available for deterministic rendering.

## Scope

- Execution restriction: Juan or a non-LLM author must create/edit all templates, HTML, CSS, and JavaScript. LLM agents must not claim the presentation-authoring work or generate those assets through code strings or generators.
- Implement the approved Tokyo Night programme/catalogue/comparison layouts, semantic landmarks, keyboard behavior, responsive styles, topic states, source/methodology/citation areas, and long-content handling.
- Keep semantic content supplied by canonical data rather than hardcoded party summaries; provide licensed, pinned assets and representative render examples.
- Go plumbing may be implemented later by an agent against these supplied interfaces.

## Acceptance criteria

- [ ] Human/non-LLM-authored assets and authorship evidence exist and implement the approved presentation interface.
- [ ] Representative programme and comparison views meet the site contract at its concrete desktop/mobile sizes, with usable keyboard focus and all topic states.
- [ ] No LLM-authored template/HTML/CSS/JavaScript or equivalent embedded presentation generator is introduced.
- [ ] Source text interpolation uses approved context-aware escaping, without unsafe raw markup insertion.

## Quality Gates

### Gate 1: Acceptance criteria

Use the site-quality contract's representative data, viewport matrix, required landmark/keyboard checklist, and approved screenshot examples. Pass only when every defined layout/state and human-authorship condition is satisfied; record evidence before renderer integration. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Use the site-quality contract's representative data, viewport matrix, required landmark/keyboard checklist, and approved screenshot examples. Pass only when every defined layout/state and human-authorship condition is satisfied; record evidence before renderer integration.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0042 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Verify Spanish visible labels, Tokyo Night requirements, licence/asset provenance, and consistency with docs/design/site-quality.md. This is explicitly not executable by an LLM coding agent.

### Gate 6: Trust-boundary and security behavior

Render hostile text samples through the approved template interface; verify escaping and no script execution, unsafe links, hidden citations, or keyboard traps. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0042`:

1. Record all gate results with `tq note`, then run `tq done TQ-0042`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0042` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#deterministic-rendering`
- `docs/SPEC.md#site-validation`
- `docs/SPEC.md#out-of-scope`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:24+01:00 — TQ-0057 review: retained in inbox because presentation assets require Juan or a non-LLM author. An LLM execution agent must not claim or author this work. Downstream HTML/site dependencies remain intact until approved human-authored assets are supplied.
- 2026-11-05T00:20:21+01:00 — Owner supplied an existing HTML reference during TQ-0067. Its bytes are imported unchanged under assets/reference and the private renderer binds only programme JSON data; desktop/mobile evidence is retained with the real preview. This does not claim all production presentation acceptance criteria or promote this inbox task.
