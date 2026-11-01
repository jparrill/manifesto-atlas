---
id: TQ-0051
title: Orchestrate idempotent programme processing and local preview
status: todo
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/cli
  - area/site
  - area/project
depends_on:
  - TQ-0050
  - TQ-0007
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T18:56:24+01:00
---

## Outcome

make process runs applicable stages for one programme and make serve previews only eligible deterministic output.

## Scope

- Replace remaining planned-stage failures only when their real implementation exists; keep all 12 Spec targets mapped to stable Go seams.
- Implement graph-aware stage selection/reuse/invalidation, bounded retries/cancellation, retained failure reports, and resumable processing.
- Provide local static preview without runtime generation/backend, exposed upload endpoints, or CI/deployment behavior.
- Document explicit expensive model/fixture commands and separate offline verification from fresh semantic generation.
- Provide the focused Make gate `make test-process` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] A programme can traverse submit/provenance/registration/extraction/audit/normalization/content/render/site/publication-check through documented Make interfaces.
- [ ] Unchanged inputs reuse only valid current artifacts; changing one upstream input reruns exactly required descendants and never overwrites source bytes.
- [ ] Failures/cancellation stop downstream work and return non-zero; local preview serves no quarantined/draft/failed source or runtime-generated content.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-process`. Pass when automated cases prove: full fake-adapter path; idempotent repeat; each input invalidation; failure propagation; cancellation/resume; static preview filtering. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-process`. Pass when automated cases prove: full fake-adapter path; idempotent repeat; each input invalidation; failure propagation; cancellation/resume; static preview filtering. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0051 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Validate programme IDs/paths, concurrent runs and partial writes, localhost binding, path traversal in preview, and denied access to quarantine/secrets/source-processing internals. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0051`:

1. Record all gate results with `tq note`, then run `tq done TQ-0051`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0051` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#module-seams`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:24+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
