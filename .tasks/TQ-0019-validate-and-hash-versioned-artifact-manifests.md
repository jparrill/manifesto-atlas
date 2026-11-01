---
id: TQ-0019
title: Validate and hash versioned artifact manifests and reports
status: inbox
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
depends_on:
  - TQ-0018
  - TQ-0008
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:30+01:00
---

## Outcome

All stages can persist deterministic versioned artifacts and detect changed or malformed inputs offline.

## Scope

- Implement approved serialization, hashes, schema/manifest validation, graph traversal, and required input/version records.
- Provide positive/negative fixtures for report outcomes and every stage identity; retain failed/insufficient reports without treating them as approvals.
- Reject missing/unknown schema versions, dangling references, duplicate identities, cycles, and stale dependencies.
- Provide the focused Make gate `make test-artifacts` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Repeated identical logical inputs produce identical canonical bytes/hashes under the approved timestamp policy.
- [ ] Changing any source/extraction/schema/taxonomy/prompt/model/tool/configuration identity invalidates precisely its descendants.
- [ ] Malformed or incomplete manifests cannot advance a stage; independent secondary evidence never replaces canonical data.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-artifacts`. Pass when automated cases prove: determinism; each input mutation; unknown version; malformed report; dangling reference; graph cycle. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-artifacts`. Pass when automated cases prove: determinism; each input mutation; unknown version; malformed report; dangling reference; graph cycle. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0019 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Bound JSON sizes/depth, reject unsafe artifact paths and external references, and never follow a manifest into arbitrary files or network locations. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0019`:

1. Record all gate results with `tq note`, then run `tq done TQ-0019`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0019` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`
