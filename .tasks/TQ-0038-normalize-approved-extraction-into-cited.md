---
id: TQ-0038
title: Normalize approved extraction into cited Spanish programme JSON
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/normalization
  - area/skills
  - area/cli
depends_on:
  - TQ-0037
  - TQ-0024
  - TQ-0035
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

make normalize invokes manifesto-normalizer to create a schema-valid cited draft using the approved taxonomy.

## Scope

- Author the project skill after schema/taxonomy stability, and wire the restricted runner to canonical validation and draft artifact storage.
- Produce factual Spanish summaries and atomic measures, responsible institutions, all topic states, and original-page citations from approved extraction only.
- Record prompt/model/tool/source/extraction/schema/taxonomy hashes; do not generate public HTML or approve the draft.
- Evaluate both golden sources and controlled unsupported-content/omission examples with skill-authoring guidance.
- Provide the focused Make gate `make test-normalizer` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Both golden inputs yield schema-valid draft records with every topic, traceable public statements, and no invented positions or future assessment scores.
- [ ] Unapproved/stale extraction, unavailable evidence, unknown taxonomy, or malformed output stops normalization.
- [ ] Recorded real-model evaluations include exact inputs/prompts/models; deterministic regression validates stored outputs without rerunning models.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-normalizer`. Pass when automated cases prove: golden recorded drafts; all topic states; uncited/invented content; stale/unapproved extraction; malformed output; injection. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-normalizer`. Pass when automated cases prove: golden recorded drafts; all topic states; uncited/invented content; stale/unapproved extraction; malformed output; injection. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0038 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Test source instructions requesting shell access, hidden approval, markup generation, or unsupported political conclusions; capability restrictions and draft-only status must hold. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0038`:

1. Record all gate results with `tq note`, then run `tq done TQ-0038`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0038` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#agent-skills`
- `docs/SPEC.md#canonical-programme-model`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`
