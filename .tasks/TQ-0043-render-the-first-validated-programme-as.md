---
id: TQ-0043
title: Render the first validated programme as deterministic Spanish Markdown
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/rendering
  - area/cli
depends_on:
  - TQ-0041
  - TQ-0037
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make build emits readable Spanish Markdown with all topics, measures, conditions and original-page citations from the exact validated local canonical draft.

## Scope

- Use Go Markdown rendering from canonical data, not extracted Markdown, with source-selection wording, archived PDF/checksum and methodology links.
- Preserve every factual Statement/citation and distinct Spanish topic-state labels. No HTML assets, comparison support, browser or signature is required.
- Escape source/Markdown/link data without altering factual wording; retain deterministic output hashes in the local run record.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] All 21 topics, five labels, exact summaries/measures/conditions/institutions/citations and source metadata render in stable order.
- [ ] Unvalidated/stale input and hostile link/path/Markdown content cannot produce public eligible output.
- [ ] Unchanged canonical/input configuration yields byte-identical Markdown and working local PDF/page/methodology links.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-markdown`. Pass when automated cases prove complete local fixture, five states/order, every measure/condition/citation, source authority wording, escaping, stale draft and byte determinism; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-markdown`. Pass with exit zero and named positive/negative cases for complete local fixture, five states/order, every measure/condition/citation, source authority wording, escaping, stale draft and byte determinism. Do not substitute a fake for required actual-tool/model evidence.

### Gate 3: Regression

Run `make verify` and `go test ./...`. Pass with exit zero, no Taskqueue warnings, no live model/source-site requests and no skipped required tests.

### Gate 4: Static and structural validation

Run `make check` and `git diff --check`; validate changed schemas/fixtures and documented relative links. Pass with no formatting/compiler/vet/schema/queue/whitespace errors.

### Gate 5: Documentation and specification consistency

Compare changed commands, schemas, fixtures and documentation with docs/SPEC.md and docs/design/MVP.md. Pass when actual behavior/retention/source authority match, public prose is Spanish, technical prose is English, and no LLM authors presentation assets.

### Gate 6: Trust-boundary and security behavior

Use the negative cases named in Gates 1/2 through the real changed interface. Pass only when invalid, stale, hostile, cancelled or incomplete inputs produce the stated safe outcome with no false approval, unsafe writes or leaked credentials. Retain any applicable trust-boundary evidence.

### Gate 7: Artifact and task record

Verify required retained artifacts/input hashes are current, deterministic operations repeat identically, and no unrelated user work or historical source bytes are removed. Record every gate as PASS with commands/artifact evidence using `tq note` before completion. A failed, unavailable or skipped required gate keeps this task open.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended changes and final task record, push and confirm upstream, then continue only approved ready work within the assignment and execution role. Delivery failures restore in-progress and block continuation; never bypass hooks/signing or promote inbox tasks.

## Specification

- docs/SPEC.md#deterministic-rendering
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
