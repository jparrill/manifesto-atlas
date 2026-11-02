---
id: TQ-0041
title: Require both current independent content passes without report PKI
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/normalization
  - area/cli
depends_on:
  - TQ-0039
  - TQ-0040
  - TQ-0063
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make validate-content marks only the exact current draft eligible when both actual independent claim and complete-source coverage reports pass.

## Scope

- Join current selected source, extraction, canonical schema, prompt/model/tool configuration and checked inventories using the local run record.
- Persist failures/disagreement and distinguish draft/validating/validated derived state; no mutable payload approval, report signing or generic DAG.
- One model is allowed only in separate clean jobs. Admission comes from the trusted local job/checkout, never a contributor-provided pass report or a human content override.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Matching complete pass/pass reports advance only their exact source/extraction/draft and role inputs.
- [ ] Fail/insufficient/error, disagreement, missing/stale/swapped report, omitted checked target or reused generation context keeps the draft unpublishable.
- [ ] Input changes invalidate the affected programme's current result while historical reports remain inspectable; neither a checksum nor manual flag grants authority.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-content-gate`. Pass when automated cases prove pass/pass, every non-pass combination, disagreement, checked-set omission, stale/swapped/contributor report, payload/manual approval forgery and clean one-model jobs; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-content-gate`. Pass with exit zero and named positive/negative cases for pass/pass, every non-pass combination, disagreement, checked-set omission, stale/swapped/contributor report, payload/manual approval forgery and clean one-model jobs. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#llm-validation
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
