---
id: TQ-0050
title: Verify the complete publishable artifact graph offline
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/site
  - area/cli
depends_on:
  - TQ-0049
  - TQ-0041
  - TQ-0035
  - TQ-0027
  - TQ-0020
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T18:56:23+01:00
---

## Outcome

make publish-check accepts only complete, current, authenticated graphs and emits an exact static publication manifest without deploying.

## Scope

- Verify immutable source bytes, verified provenance, approved extraction, validated canonical data, both independent content reports, deterministic outputs, and all site gates.
- Validate every hash/version/signature/dependency and reject stale downstream artifacts, including historical sources modified in place.
- Build a catalogue/distribution manifest containing only eligible programmes and atomically transition eligible records to published according to the approved lifecycle.
- No CI workflow, GitHub Pages deployment, or runtime model execution is included.
- Provide the focused Make gate `make test-publication` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Complete passing stored graphs verify offline without a model, secret signing key, or live source website.
- [ ] Probable/unverified/rejected provenance, failed/insufficient audits, draft content, disagreement, stale/tampered reports, or failed rendering/site checks block publication.
- [ ] No incomplete intermediate artifact enters the public catalogue; generated distribution is deterministic and source bytes stay unchanged.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-publication`. Pass when automated cases prove: complete graph; each missing/failing/stale gate; source mutation; model/prompt/taxonomy/template change; forged signer; deterministic manifest. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-publication`. Pass when automated cases prove: complete graph; each missing/failing/stale gate; source mutation; model/prompt/taxonomy/template change; forged signer; deterministic manifest. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0050 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Exercise graph/report substitution, traversal, omitted blockers, stale signed replay, unauthorized keys, and partial publication writes; fail closed without network/model fallback. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0050`:

1. Record all gate results with `tq note`, then run `tq done TQ-0050`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0050` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#publication-and-ci`
- `docs/SPEC.md#end-to-end-tests`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
