---
id: TQ-0031
title: Produce independent text evidence and suspicious-page images
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/extraction
  - area/fidelity
depends_on:
  - TQ-0012
  - TQ-0021
  - TQ-0027
  - TQ-0028
  - TQ-0029
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

Auditors can compare independent page text and inspect original rendered pages instead of trusting one converter.

## Scope

- Implement the approved secondary extractor and pinned original-page renderer with exact PDF page mapping.
- Generate page images for textless/suspicious pages and explicit multimodal audit requests, with bounded caching and hashes.
- Implement approved OCR escalation with clear provenance; independent evidence never becomes a second canonical representation.
- Provide the focused Make gate `make test-independent-evidence` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Both golden sources yield independent page-aligned text and requested images with pinned tool identity.
- [ ] Textless, unreadable, suspicious, and OCR-escalated pages retain warnings and cannot disappear from coverage.
- [ ] Image/evidence changes invalidate dependent reports; resource-bounded repeated runs reuse only valid matching artifacts.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-independent-evidence`. Pass when automated cases prove: golden page mapping; native/textless/suspicious page; OCR escalation; image hashes; render failure; stale cache. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-independent-evidence`. Pass when automated cases prove: golden page mapping; native/textless/suspicious page; OCR escalation; image hashes; render failure; stale cache. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0031 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Sandbox rendering/OCR, enforce pixel/page/memory/time budgets, reject unsafe output paths, and retain unresolved failures rather than approving blank evidence. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0031`:

1. Record all gate results with `tq note`, then run `tq done TQ-0031`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0031` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#extraction-artifacts`
- `docs/SPEC.md#extraction-fidelity`
- `docs/SPEC.md#extraction-tests`
- `docs/QUALITY_GATES.md`
