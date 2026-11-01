---
id: TQ-0006
title: Ignore local tooling state and reproducible transient artifacts
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - agent/developer
  - area/project
created: 2026-11-01T17:07:32+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

Local editor, Serena, operating-system, and disposable build state stay out of Git without hiding archive evidence.

## Scope

- Inspect actual tracked/untracked paths and planned outputs before selecting narrow ignore rules.
- Ignore the repository-local .serena/ directory, OS/editor temporary files, and justified disposable caches, quarantine scratch space, and distribution output.
- Keep Taskqueue, specifications, schemas, canonical JSON, immutable PDFs, committed reports, checksums, and required golden fixtures trackable; document any output whose retention policy is not yet decided.
- Do not remove or overwrite existing local files or untrack user work.

## Acceptance criteria

- [ ] git status --short no longer lists .serena/; local files still exist.
- [ ] git check-ignore --no-index -v .serena/project.local.yml .DS_Store succeeds.
- [ ] Representative source PDFs, programme JSON, reports, fixtures, .tasks files, HANDOFF.md, and .taskqueue.yaml are not ignored.
- [ ] Every additional rule has an observed or specification-backed purpose; no blanket PDF/JSON/Markdown exclusion exists.

## Quality Gates

### Gate 1: Acceptance criteria

Run git check-ignore --no-index for positive and negative path tables, including nested source/report paths. Pass when only the explicitly disposable cases match and no tested evidence path matches. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run git check-ignore --no-index for positive and negative path tables, including nested source/report paths. Pass when only the explicitly disposable cases match and no tested evidence path matches.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0006 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Use path-only probes: do not read local credentials or delete scratch files. Assert committed evidence paths are not hidden by broad rules. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0006`:

1. Record all gate results with `tq note`, then run `tq done TQ-0006`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0006` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#publication-and-ci`
- `docs/SPEC.md#immutable-source-storage`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T17:12:56+01:00 — Implemented the owner-requested .serena/ exclusion and macOS .DS_Store metadata rule. Additional editor/cache/build/quarantine patterns are not justified by current repository artifacts and are left for their owning implementation tasks once retention contracts exist; immutable PDFs, canonical JSON, reports, fixtures, and Taskqueue remain trackable. No files were removed or untracked.
  Quality gates:
  - Gate 1 PASS — root .gitignore contains the requested .serena/ rule and justified macOS metadata rule; git status --short --untracked-files=all -- .serena is empty and local Serena files remain present.
  - Gate 2 PASS — git check-ignore --no-index probes passed for 5 ignored paths and 10 non-ignored evidence/configuration paths, including nested PDFs, canonical JSON, reports, fixtures, task files, HANDOFF.md, and .taskqueue.yaml.
  - Gate 3 PASS — tq list --json, tq show TQ-0006 --json, and git diff --check exited zero with empty stderr. No Makefile exists; the declared pre-Makefile regression applies.
  - Gate 4 PASS — git diff --no-index --check checked actual untracked .gitignore and task-file bytes with zero whitespace findings; ignore-rule source was the new root .gitignore.
  - Gate 5 PASS — English comments explain each rule; evidence preservation matches immutable-source and publication requirements. No speculative output retention decision is hidden in broad ignore rules.
  - Gate 6 PASS — path-only positive/negative checks performed without reading credentials, deleting scratch files, or modifying source evidence. No broad PDF/JSON/Markdown exclusion exists.
  - Gate 7 PASS — rules are deterministic, no generated artifacts are required, local files are preserved, and all evidence is recorded before completion.
