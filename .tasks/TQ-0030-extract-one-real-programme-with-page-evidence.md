---
id: TQ-0030
title: Extract one real programme with page evidence and deterministic blockers
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/extraction
  - area/cli
depends_on:
  - TQ-0021
  - TQ-0027
  - TQ-0028
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make extract/validate-extraction produce complete page-aware evidence for the first registered real source without OCR surprises or a calibrated fidelity score.

## Scope

- Wire one actual Docling adapter into the bounded deployment and atlas.extraction.local/v1; retain native raw output, page/block text/order, warnings and actual tool/configuration/resource observations.
- Implement page/count/reference and material-warning blockers formerly split across TQ-0032/TQ-0035. No signature, numeric score, independent extraction-audit skill or second golden PDF prerequisite.
- Keep native-text OCR disabled; mark unreadable/scanned/materially ambiguous regions unresolved rather than silently generating text. Deterministic checks do not claim semantic fidelity.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] The first real source yields all contiguous original PDF pages and known observed headings/numbers/negations with correctly linked blocks and warnings.
- [ ] Missing pages, invalid mappings, corrupt output, material unresolved warnings and over-budget/cancelled converter jobs fail with retained evidence and no source mutation.
- [ ] Actual repeated unchanged converter/configuration runs reproduce the declared normalized deterministic output, or report drift without overwriting it.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-extraction`. Pass when automated cases prove native/table/list/column fixtures, missing/wrong page, original source mismatch, OCR-off, warning/blocker propagation, tool/config change and failed cleanup; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails. Also run `make evaluate-extraction TARGET=fixtures/frente/source.json` and retain actual source/model/tool/input identities and original outputs; unavailable required resources keep the task open, never a skipped pass.

### Gate 2: Focused behavior

Run `make test-extraction`. Pass with exit zero and named positive/negative cases for native/table/list/column fixtures, missing/wrong page, original source mismatch, OCR-off, warning/blocker propagation, tool/config change and failed cleanup. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#extraction-artifacts
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
