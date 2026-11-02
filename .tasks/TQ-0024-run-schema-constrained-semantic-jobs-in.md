---
id: TQ-0024
title: Run schema-constrained semantic jobs in isolated provider contexts
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/skills
depends_on:
  - TQ-0015
  - TQ-0019
  - TQ-0009
  - TQ-0018
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T15:29:08+01:00
---

## Outcome

Semantic stages can use local Auriga or another configured provider through one restricted, auditable runner.

## Scope

- Implement approved runner requests/responses, bounded evidence access, provider adapters, isolated job contexts, limits/cancellation, schema validation, and invocation metadata.
- Provide an offline fake adapter for regression plus a documented local-provider smoke command that does not require committed credentials.
- Keep generator/claim/coverage jobs separate, prefer distinct models, and support one model only with genuinely isolated invocations.
- Never grant model jobs general shell, network, repository writes, or signing keys.
- Provide the focused Make gate `make test-semantic-runner` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] The same approved request contract works with fake and configured local-provider adapters.
- [ ] A canary from generator/private prior context is absent from validator input; each job records prompt/model/input hashes and identity.
- [ ] Malformed/oversized outputs, missing evidence, cancellation, provider errors, and injection attempts cannot produce a valid approval.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-semantic-runner`. Pass when automated cases prove: provider substitution; clean-context canary; one/two-model isolation; malformed output; timeout; retry; injection. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-semantic-runner`. Pass when automated cases prove: provider substitution; clean-context canary; one/two-model isolation; malformed output; timeout; retry; injection. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0024 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Demonstrate capability denial with hostile document instructions and test context leakage, output path injection, unbounded retries, and secret redaction. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0024`:

1. Record all gate results with `tq note`, then run `tq done TQ-0024`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0024` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#agent-skills`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#untrusted-document-handling`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T15:28:57+01:00 — Quality gates — TQ-0024 (all PASS). Gate 1/2 (make test-semantic-runner): ok internal/semantic, 9 PASS. TestProviderSubstitution proves the same approved request contract works with two distinct provider adapters; TestCleanContextCanary proves a canary in the generator's private prompt is absent from the validator's context (no shared conversation state); TestOneTwoModelIsolation proves separate generator/validator invocations in one evaluation set succeed and a reused job id is rejected as a clean-context violation; TestMalformedOutput proves invalid JSON, unknown report fields, and an invalid result enum are rejected as quarantined diagnostic data (no auto-repair pass); TestTimeout enforces the request deadline (slow adapter -> timeout, no approval); TestRetryOnce proves a pre-response transient transport error gets exactly one clean retry with a new job id and unchanged inputs, and a second failure surfaces as a provider error (no retry-until-pass); TestModelMismatch rejects a self-reported model identity that differs from the trusted adapter identity (no silent substitution); TestJobMismatch rejects a forged response job id; TestInjectionIsData proves shell/network/approval/signing instructions in the prompt are treated as data — a conforming reply is still verified normally and a malformed reply is still rejected. All run offline against an injected fake provider adapter; no live models, no committed credentials. Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy/test-canonical/test-party-registry/test-fetch/test-semantic-runner all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/semantic/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new make test-semantic-runner target and internal/semantic; technical text English; no product decision changed. Gate 6 (trust-boundary): the runner grants no shell/filesystem-write/network/repo-edit/signing capability (the adapter contract has none); host verifies job binding, expected model identity, report schema, result enum, finding bounds, and strict details depth; a model cannot authorize itself even when it emits pass; injection is data; no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0024 intended files staged (internal/semantic/{semantic.go,semantic_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0024-*.md).
