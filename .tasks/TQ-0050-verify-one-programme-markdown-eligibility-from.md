---
id: TQ-0050
title: Verify one programme Markdown eligibility from local run evidence
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/site
  - area/cli
depends_on:
  - TQ-0041
  - TQ-0027
  - TQ-0043
  - TQ-0063
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make publish-check verifies the exact selected immutable source and current independent reports offline and emits eligible Markdown output without deployment.

## Scope

- Check actual source/metadata/evidence, extraction, draft/schema/prompt/model/tool hashes, complete pass/pass role reports and deterministic Markdown references in atlas.run/v1.
- Use the trusted local job/checkout admission boundary; arbitrary incoming reports and draft booleans are not approval. No report signatures or DAG graph engine.
- Initially reject HTML eligibility as unavailable; TQ-0049 later extends it. Do not make first Markdown depend on supplied HTML, browser, comparisons or deployment.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] One complete current selected-source Markdown run verifies without a model, source website or signing key and emits a deterministic output manifest.
- [ ] Missing/failed/insufficient/disagreeing/stale/swapped reports, changed original bytes, incomplete checked sets and unselected/contributor approval records reject.
- [ ] No draft/preview/partial artifact enters public selection; HTML cannot be silently approved before its additional checks exist.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-publication`. Pass when automated cases prove current local Markdown run, every missing/non-pass/stale input, source mutation, checked-set omission, contributor/manual approval, partial installation and unsupported HTML mode; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-publication`. Pass with exit zero and named positive/negative cases for current local Markdown run, every missing/non-pass/stale input, source mutation, checked-set omission, contributor/manual approval, partial installation and unsupported HTML mode. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#publication-and-ci
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
