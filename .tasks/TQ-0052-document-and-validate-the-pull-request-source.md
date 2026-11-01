---
id: TQ-0052
title: Document and validate the pull-request source contribution flow
status: inbox
priority: normal
labels:
  - type/docs
  - agent/developer
  - area/project
  - area/provenance
depends_on:
  - TQ-0051
  - TQ-0027
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

A contributor can propose a programme through Git without claiming that uploaded bytes are official.

## Scope

- Write a Spanish public contribution guide plus English technical command reference using the actual implemented Make interface.
- Explain required party/election/geography/source-URL metadata, quarantine, independent canonical download/checksum comparison, revisions/duplicates, publisher rights, and failure states.
- Provide safe sample submission metadata and a documented PR checklist; no GitHub Actions/workflow/upload service is introduced.
- Update README/HANDOFF references to actual capabilities and available skills instead of stale planned names.

## Acceptance criteria

- [ ] A disposable-checkout walkthrough with fixture source evidence follows the documented submit/verify-source/add path and produces the expected immutable record.
- [ ] The guide explicitly distinguishes candidate bytes, official evidence, verified/probable/unverified/rejected states, and public eligibility.
- [ ] All documented commands/paths/links exist, visible public prose is Spanish, and no private paths/credentials are required.

## Quality Gates

### Gate 1: Acceptance criteria

Run the guide's exact commands with the controlled fixture transport in a disposable checkout, then run make verify. Pass when the expected record/checksums/states match and all relative links resolve. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run the guide's exact commands with the controlled fixture transport in a disposable checkout, then run make verify. Pass when the expected record/checksums/states match and all relative links resolve.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0052 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Check README, HANDOFF, contribution guide, CLI help, public language, MIT-versus-publisher rights, and the actual stage/status semantics agree.

### Gate 6: Trust-boundary and security behavior

Demonstrate mismatch and unverified examples stay unpublishable; ensure the contribution workflow grants no executable authority to PDFs, submitted metadata, or external URLs. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0052`:

1. Record all gate results with `tq note`, then run `tq done TQ-0052`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0052` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#solution`
- `docs/SPEC.md#user-stories`
- `docs/SPEC.md#language-and-licensing`
- `docs/QUALITY_GATES.md`
