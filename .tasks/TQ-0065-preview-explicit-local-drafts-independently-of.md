---
id: TQ-0065
title: Preview explicit local drafts independently of publication
status: todo
priority: normal
labels:
  - type/feature
  - area/site
  - agent/developer
depends_on:
  - TQ-0044
created: 2026-11-02T18:05:04+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make serve previews an explicitly generated local fixture/draft root without requiring final publication approval or exposing processing files.

## Scope

- Use the standard Go HTTP file-serving mechanisms over a declared generated preview root, loopback only; no upload/API/runtime generation or filesystem browsing outside that root.
- Consume already supplied human/non-LLM assets and explicitly labelled draft/fixture output; do not author HTML/CSS/JavaScript/templates or embedded presentation strings.
- Keep preview storage separate from eligible public output; no native dependency on publish-check or content approval is needed for an explicit private preview.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] An explicit draft/fixture preview works before publish-check and is clearly non-public; normal publication selection excludes it.
- [ ] The server binds only loopback, serves only existing generated files, and shuts down on cancellation.
- [ ] Traversal/symlinks, source/quarantine/credential paths, directory escapes, uploads and non-loopback exposure are denied.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-preview`. Pass when automated cases prove draft preview before publication, loopback binding, read-only files, cancellation, traversal/symlink denial and public-catalogue exclusion; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-preview`. Pass with exit zero and named positive/negative cases for draft preview before publication, loopback binding, read-only files, cancellation, traversal/symlink denial and public-catalogue exclusion. Do not substitute a fake for required actual-tool/model evidence.

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

- 2026-11-02T18:05:04+01:00 — Created directly in todo under Juan's explicit TQ-0061 request to revise existing tasks and create the missing MVP outcomes. Outcome, scope, criteria, seven gates, priority, dependencies and execution role are declared. No authority to promote unrelated inbox proposals is granted.
