---
id: TQ-0049
title: Require basic deterministic site checks for HTML eligibility
status: todo
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/site
  - area/validation
  - area/cli
depends_on:
  - TQ-0046
  - TQ-0047
  - TQ-0050
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:09:33+01:00
---

## Outcome

make validate-site and HTML publish-check require current parity/link/structure and basic browser/accessibility results, without model QA or report PKI.

## Scope

- Combine the simple input-bound site reports and extend TQ-0050's Markdown-only eligibility to HTML after these actual checks exist.
- Integrate actual default public HTML rendering with the current content gate; TQ-0044 supplies private previews only until this integration. Public mode must not accept a preview or placeholder pass.
- No pixel baseline, visual skill, signature, comparison or generic graph closure is required; any required component failure/missing/stale result blocks HTML.
- Preserve independent content passes and exclude explicit private preview roots from public output.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] HTML becomes eligible only with matching current content, outputs/assets, parity and basic accessibility/mobile checks.
- [ ] Each failed/missing/stale/swapped component and preview-root substitution blocks HTML without affecting valid Markdown-only eligibility.
- [ ] Offline replay and unchanged deterministic checks reproduce the local eligibility/output manifest; no model/site-network call occurs.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-site-gate`. Pass when automated cases prove all basic checks pass, each failed/missing/stale component, swapped output/assets, preview exclusion and independent Markdown eligibility; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-site-gate`. Pass with exit zero and named positive/negative cases for all basic checks pass, each failed/missing/stale component, swapped output/assets, preview exclusion and independent Markdown eligibility. Do not substitute a fake for required actual-tool/model evidence.

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
- 2026-11-02T18:07:29+01:00 — Under TQ-0061, TQ-0044 initially supplies safe labelled previews only. This task also integrates actual default public HTML rendering with current TQ-0041 content evidence and TQ-0050/0046/0047 eligibility, before reporting HTML available. Private previews remain excluded; no placeholder pass may substitute for that integration.
