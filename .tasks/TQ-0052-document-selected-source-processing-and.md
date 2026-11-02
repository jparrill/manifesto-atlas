---
id: TQ-0052
title: Document selected-source processing and proposal-only contributions
status: todo
priority: normal
labels:
  - type/docs
  - agent/developer
  - area/project
  - area/provenance
depends_on:
  - TQ-0051
  - TQ-0027
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

A maintainer can follow the real local workflow and a contributor can propose source evidence without self-approving origin or content.

## Scope

- Write Spanish public instructions and English technical commands for selected local candidates, safe inspection, immutable registration, generation/validation and Markdown output.
- Explain the separately owner-controlled selection and trusted local report admission; no automatic arbitrary-URL verification, publisher PKI, report self-approval or workflow is implied.
- Walk through actual implemented commands in a disposable checkout and retain source-rights/revision/failure guidance.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] The exact documented fixture walkthrough yields the expected source/run/checksum/output states with existing commands.
- [ ] Contributor source/report proposals remain untrusted until separate operator selection/processing; a supplied pass cannot grant public eligibility.
- [ ] Public prose is Spanish, links/paths work and no secret/private path, invented origin certification or automatic deployment is required.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-contribution-guide`. Pass when automated cases prove disposable documented selected-source walkthrough, mismatched/unselected candidate, contributor pass rejection, revisions/rights and relative links; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-contribution-guide`. Pass with exit zero and named positive/negative cases for disposable documented selected-source walkthrough, mismatched/unselected candidate, contributor pass rejection, revisions/rights and relative links. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#language-and-licensing
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:24+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
