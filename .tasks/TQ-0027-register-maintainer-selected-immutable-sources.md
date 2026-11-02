---
id: TQ-0027
title: Register maintainer-selected immutable sources and revisions
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/source-registry
  - area/cli
depends_on:
  - TQ-0021
  - TQ-0063
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make verify-source/add validate the trusted local selection and register exact immutable source bytes with factual metadata and origin evidence.

## Scope

- Implement atlas.source.local/v1 metadata and atlas.selection/v1 binding, explicit maintainer-selected authority, year/programme/checksum storage, idempotent duplicate and appended corrected revisions.
- Verify the separately operator-selected source/metadata/evidence hashes; contributor files and model-authored flags do not activate selection.
- Use safe rooted atomic no-replace installation and committed-source mutation checks; no network trust registry or report signatures.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] A selected inspected source registers with correct factual identity, origin evidence and publisher rights; a repeated identical registration is idempotent.
- [ ] A changed selected source creates an explicit new revision and never replaces historical bytes or metadata.
- [ ] Unselected/stale/mismatched bytes or evidence, unsafe paths/symlinks, conflicting metadata and interrupted writes cannot create an accepted record.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-source-registry`. Pass when automated cases prove trusted selection vs contributor self-selection, exact hashes, first/duplicate/corrected registration, committed-byte overwrite, unsafe paths and interrupted installation; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-source-registry`. Pass with exit zero and named positive/negative cases for trusted selection vs contributor self-selection, exact hashes, first/duplicate/corrected registration, committed-byte overwrite, unsafe paths and interrupted installation. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#immutable-source-storage
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
