---
id: TQ-0019
title: Validate and hash versioned artifact manifests and reports
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/validation
depends_on:
  - TQ-0018
  - TQ-0008
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T10:35:28+01:00
---

## Outcome

All stages can persist deterministic versioned artifacts and detect changed or malformed inputs offline.

## Scope

- Implement approved serialization, hashes, schema/manifest validation, graph traversal, and required input/version records.
- Provide positive/negative fixtures for report outcomes and every stage identity; retain failed/insufficient reports without treating them as approvals.
- Reject missing/unknown schema versions, dangling references, duplicate identities, cycles, and stale dependencies.
- Provide the focused Make gate `make test-artifacts` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Repeated identical logical inputs produce identical canonical bytes/hashes under the approved timestamp policy.
- [ ] Changing any source/extraction/schema/taxonomy/prompt/model/tool/configuration identity invalidates precisely its descendants.
- [ ] Malformed or incomplete manifests cannot advance a stage; independent secondary evidence never replaces canonical data.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-artifacts`. Pass when automated cases prove: determinism; each input mutation; unknown version; malformed report; dangling reference; graph cycle. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-artifacts`. Pass when automated cases prove: determinism; each input mutation; unknown version; malformed report; dangling reference; graph cycle. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0019 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Bound JSON sizes/depth, reject unsafe artifact paths and external references, and never follow a manifest into arbitrary files or network locations. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0019`:

1. Record all gate results with `tq note`, then run `tq done TQ-0019`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0019` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T10:35:16+01:00 — Quality gates:
  - Gate 1 PASS — AC1 (identical logical inputs -> identical canonical bytes/hashes) proven by TestDeterminism and TestCanonicalReferenceVectors (key-order/whitespace variants hash to one ID). AC2 (each source/extraction/schema/taxonomy/prompt/model/tool/configuration mutation invalidates precisely its descendants) proven by TestStalenessPrecise across seven leaf kinds, each asserted equal to computed reverse-reachable descendants. AC3 (malformed/incomplete manifests cannot advance; secondary evidence never replaces canonical) proven by TestManifestStructure (unknown version/kind/slot/field and bad hash rejected), TestValidatePayload, TestReportResult (only pass advances), and graph validation (dangling/cycle/kind-mismatch rejected).
  - Gate 2 PASS — make test-artifacts: automated cases prove determinism; each input mutation (staleness); unknown version; malformed report; dangling reference; graph cycle; plus duplicate-id, id-mismatch, payload integrity, and all 26 kinds.
  - Gate 3 PASS — make verify: git whitespace, taskqueue integrity, gofmt, go build, go vet, docs structure, relative links, go test (test-cli + test-artifacts) all pass.
  - Gate 4 PASS — gofmt -l clean; go build and go vet ./... clean via make check; git diff --check clean.
  - Gate 5 PASS — docs/MAKEFILE.md updated (test-artifacts added to baseline surface, Go-tests bullet, and internal/artifact layout); Makefile help lists the target; interfaces match docs/design/artifact-contract.md; technical text English; no public (Spanish) content changed.
  - Gate 6 PASS — bounded sizes (16MiB payload, 1MiB manifest, 64MiB graph, 10000 nodes, 1024 edges) and depth 64 (TestLimits, TestCanonicalDepthLimit); rejects BOM, trailing values, duplicate members, fractional/exponent/negative-zero numbers, out-of-range integers, invalid UTF-8, and unpaired/invalid surrogates. Package is purely in-memory: no filesystem, network, model, or path following, so a manifest cannot be followed into arbitrary files/network. No secrets or private absolute paths committed.
  - Gate 7 PASS — fixtures are deterministic pure functions of inputs (no timestamps); make test-artifacts and make verify pass; all gate results recorded before done.
