---
id: TQ-0031
title: Supply bounded original-page evidence for independent validators
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/extraction
  - area/fidelity
depends_on:
  - TQ-0021
  - TQ-0030
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

Claim and coverage jobs can resolve actual original PDF page images and necessary corroborating text without trusting only the generator's extraction.

## Scope

- Render requested original pages through an existing bounded renderer and retain source/page/image checksums and actual version/configuration.
- Cited and suspicious/unreadable regions require original evidence; full-source coverage can request every page. A secondary text tool may corroborate actual needs, not become a mandatory second canonical extractor.
- No automatic OCR escalation platform, image cache DAG or second golden-source prerequisite.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] All requested first-source original pages map to the exact registered checksum/PDF indices and resolve to bounded retained evidence.
- [ ] Scanned/unreadable/render-failed pages stay explicit and block a reliable pass when required vision/evidence is unavailable.
- [ ] Wrong page/checksum, stale evidence, unsafe paths, image/resource overflow and cancellation fail without installing apparent valid evidence.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-independent-evidence`. Pass when automated cases prove original page mapping, native/textless/suspicious regions, image hashes, wrong checksum/page, render failure, limits and unsafe output; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails. Also run `make evaluate-page-evidence TARGET=fixtures/frente/source.json` and retain actual source/model/tool/input identities and original outputs; unavailable required resources keep the task open, never a skipped pass.

### Gate 2: Focused behavior

Run `make test-independent-evidence`. Pass with exit zero and named positive/negative cases for original page mapping, native/textless/suspicious regions, image hashes, wrong checksum/page, render failure, limits and unsafe output. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#extraction-fidelity
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
