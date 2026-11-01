---
id: TQ-0044
title: Integrate supplied templates into deterministic Go HTML rendering
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/rendering
  - area/cli
depends_on:
  - TQ-0042
  - TQ-0041
  - TQ-0016
  - TQ-0043
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

make build generates Spanish Tokyo Night HTML from validated canonical JSON using supplied human/non-LLM templates.

## Scope

- Implement Go renderer data binding, context-aware escaping, output paths/manifests, and deterministic catalogue/programme generation.
- Consume approved supplied presentation assets without creating or editing templates, HTML, CSS, JavaScript, or embedded presentation strings.
- Coordinate make build to emit both public formats from the same validated input and record asset/tool/input hashes.
- Provide the focused Make gate `make test-html-rendering` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] All required metadata, every topic/state/summary/measure/citation, source PDF link/checksum, and methodology are rendered from canonical JSON.
- [ ] HTML is byte-deterministic for unchanged canonical data/assets/configuration and cannot drift from Markdown semantically.
- [ ] Unsafe text is escaped, unvalidated/stale data is rejected, and an LLM execution makes no presentation-code edits.
- [ ] Public artefacts exclude failed/probable/unverified sources and incompatible revisions.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-html-rendering`. Pass when automated cases prove: complete canonical fixture; all states; JSON/Markdown/DOM parity; escaping; stale data; deterministic rerun; asset-change invalidation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-html-rendering`. Pass when automated cases prove: complete canonical fixture; all states; JSON/Markdown/DOM parity; escaping; stale data; deterministic rerun; asset-change invalidation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0044 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Exercise script/attribute/URL injection, unsafe output paths, source-link schemes, and raw-markup bypasses using Go tests and supplied fixture assets only. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0044`:

1. Record all gate results with `tq note`, then run `tq done TQ-0044`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0044` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#deterministic-rendering`
- `docs/SPEC.md#rendering-tests`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`
