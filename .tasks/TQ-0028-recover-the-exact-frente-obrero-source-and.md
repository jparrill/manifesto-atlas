---
id: TQ-0028
title: Recover the exact Frente Obrero source and origin evidence
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/source-registry
  - area/extraction
depends_on:
  - TQ-0061
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

The exact owner-approved Frente Obrero general-political document is available as preserved candidate bytes and factual origin evidence before registration infrastructure exists.

## Scope

- Locate the exact original and authoritative source evidence using explicit maintainer-authorized acquisition tools. Retain raw bytes/checksum, retrieval/origin evidence, rights and evidenced metadata as fixtures/frente/source.json and local evidence files.
- This is candidate recovery and source selection, not an automated official certification or registered/validated programme. Do not parse the PDF on an unbounded host; page/heading expectations are added during bounded extraction.
- Do not invent URLs/checksums/scope/election dates or substitute another revision/PP/2023 material. If exact recovery is unavailable, retain concrete blocker evidence and leave open.
- Provide a small offline hash/metadata fixture check without depending on pipeline registration, runners, domains or signatures.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Exact candidate bytes, SHA-256, actual official-origin evidence and factual metadata exist without requiring TQ-0027 or any model.
- [ ] The document is correctly general-political with no invented election; unknown dates/rights are explicit and display strings are evidence-backed.
- [ ] Offline fixture checking detects mutated bytes, missing evidence and metadata conflicts without parsing or fetching live resources.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-fixture-frente`. Pass when automated cases prove exact candidate/evidence hashes, factual non-electoral metadata, unknown dates/rights, missing evidence and mutated copies; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-fixture-frente`. Pass with exit zero and named positive/negative cases for exact candidate/evidence hashes, factual non-electoral metadata, unknown dates/rights, missing evidence and mutated copies. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#testing-decisions
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
