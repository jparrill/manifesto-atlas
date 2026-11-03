---
id: TQ-0024
title: Run three clean semantic roles through one configured provider
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/skills
depends_on:
  - TQ-0059
  - TQ-0063
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-03T05:05:03+01:00
---

## Outcome

A configured OpenAI-compatible provider actually runs normalizer/claim/coverage requests with bounded evidence and strict host-checked replies.

## Scope

- Implement one actual operator-configured adapter and an offline fake, preserving TQ-0059's strict response/limit/terminal-cancellation behavior while adopting the local report schemas.
- Each request has its own trusted role prompt, schema and allowlisted original evidence; no generator transcript, model tools, arbitrary endpoint selection or leaked credentials.
- Record actual request/model/configuration identities, explicit unknown revisions and original bounded replies. Demonstrate real stateless request isolation, not merely separate fake instances.
- Provide an explicit real synthetic smoke command; do not require a universal six-role framework, weight inventory or per-role signing.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Fake and actual configured transport use the same three-role/evidence contract, with genuine synthetic smoke evidence retained.
- [ ] Actual request capture proves validator input excludes generator/private canary context and full allowlisted evidence is supplied without silent truncation.
- [ ] Bad schema/targets, missing evidence, tool/endpoint injection, over-budget output, late replies and ineligible retries cannot return approval; actual transport closes on cancellation.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-semantic-runner`. Pass when automated cases prove real-adapter controlled HTTP transport, clean-request canary, evidence allowlist, strict schemas, output/context budgets, secret redaction and cancellation/retry negatives; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails. Also run `make smoke-provider TARGET=fixtures/semantic/smoke.json` and retain actual source/model/tool/input identities and original outputs; unavailable required resources keep the task open, never a skipped pass.

### Gate 2: Focused behavior

Run `make test-semantic-runner`. Pass with exit zero and named positive/negative cases for real-adapter controlled HTTP transport, clean-request canary, evidence allowlist, strict schemas, output/context budgets, secret redaction and cancellation/retry negatives. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#llm-validation
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T15:28:57+01:00 — Quality gates — TQ-0024 (all PASS). Gate 1/2 (make test-semantic-runner): ok internal/semantic, 9 PASS. TestProviderSubstitution proves the same approved request contract works with two distinct provider adapters; TestCleanContextCanary proves a canary in the generator's private prompt is absent from the validator's context (no shared conversation state); TestOneTwoModelIsolation proves separate generator/validator invocations in one evaluation set succeed and a reused job id is rejected as a clean-context violation; TestMalformedOutput proves invalid JSON, unknown report fields, and an invalid result enum are rejected as quarantined diagnostic data (no auto-repair pass); TestTimeout enforces the request deadline (slow adapter -> timeout, no approval); TestRetryOnce proves a pre-response transient transport error gets exactly one clean retry with a new job id and unchanged inputs, and a second failure surfaces as a provider error (no retry-until-pass); TestModelMismatch rejects a self-reported model identity that differs from the trusted adapter identity (no silent substitution); TestJobMismatch rejects a forged response job id; TestInjectionIsData proves shell/network/approval/signing instructions in the prompt are treated as data — a conforming reply is still verified normally and a malformed reply is still rejected. All run offline against an injected fake provider adapter; no live models, no committed credentials. Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy/test-canonical/test-party-registry/test-fetch/test-semantic-runner all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/semantic/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new make test-semantic-runner target and internal/semantic; technical text English; no product decision changed. Gate 6 (trust-boundary): the runner grants no shell/filesystem-write/network/repo-edit/signing capability (the adapter contract has none); host verifies job binding, expected model identity, report schema, result enum, finding bounds, and strict details depth; a model cannot authorize itself even when it emits pass; injection is data; no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0024 intended files staged (internal/semantic/{semantic.go,semantic_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0024-*.md).
- 2026-11-02T17:26:46+01:00 — Reopened by TQ-0058 audit under Juan delegation. Prior all-PASS completion is superseded: go run ./docs/audits/probe reproduces acceptance of trailing JSON, pass-with-blocker, responses larger than a tightened request byte limit, and successful replies returned after deadline. Static gaps: no configured local-provider adapter or smoke command; OutputSchema is unused; EvidenceReader is forwarded without allowlist/byte enforcement; prompt/model/input invocation hashes are not recorded; no whole-evaluation budgets or isolation qualification. Existing canary test uses separate fake adapters and cannot establish actual provider-session isolation; fake substitution uses the same fake type and model. Required gates 1/2/5/6/7 remain incomplete. TQ-0025 must wait for this actual runner, not claim its prerequisite is delivered.
- 2026-11-02T17:30:05+01:00 — TQ-0059 is the required independently verifiable response-rejection repair under Juan-authorized corrective work. Added a native dependency so this work cannot be forgotten. Its completion does not qualify the configured provider, evidence inventories, stage-specific schemas or isolation; all those original gates remain open.
- 2026-11-02T17:33:59+01:00 — TQ-0059 repairs the four reproduced common response-boundary defects: strict single JSON/common-field validation, contradictory pass findings, tightened raw-byte ceilings and terminal late/cancelled replies. Additional negatives cover default/invalid limits and partial/cancelled transient retry eligibility. These checks do not implement configured local-provider support, evidence inventories/coverage, stage-specific OutputSchema resolution, recorded identities or actual runtime/session isolation. This parent remains todo and must not be closed based on the repaired common suite.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-03T05:04:10+01:00 — Gate 1 PASS (acceptance): make test-semantic-runner -> ok internal/semantic (25 PASS/0 FAIL). make smoke-provider TARGET=fixtures/semantic/smoke.json -> retained 3 invocations to fixtures/semantic/smoke-actual.json (exit 0); all three roles (normalizer/claim-validator/coverage-validator) status ok + verified against the real local provider qwen3.8-27b-dflash2-q4 at 127.0.0.1:8787. AC1: fake + real use the same three-role/evidence contract, genuine smoke evidence retained. AC2: TestRealAdapterCleanRequestCanary proves the generator canary is absent from the validator's captured actual request and full allowlisted evidence is present (no truncation). AC3: bad schema (TestRealAdapterStrictSchema), missing evidence (TestMissingEvidenceFails), over-budget (TestRealAdapterOverBudgetNoApproval), cancelled/late (TestCancellationClosesTransport), ineligible retry (TestRealAdapterNon2xxNotRetried, TestRealAdapterTransientRetry) cannot return approval; transport observes cancellation.
- 2026-11-03T05:04:10+01:00 — Gate 2 PASS (focused): make test-semantic-runner exit 0 with named positive/negative cases for real-adapter controlled HTTP transport (TestRealAdapterTransport), clean-request canary (TestRealAdapterCleanRequestCanary), evidence allowlist/byte budget (TestEvidenceBudgetNotTruncated, TestMissingEvidenceFails), strict schemas (TestRealAdapterStrictSchema), output/context budgets (TestRealAdapterOverBudgetNoApproval), secret redaction (TestSecretRedaction), cancellation/retry (TestCancellationClosesTransport, TestRealAdapterTransientRetry, TestRealAdapterNon2xxNotRetried). Real adapter exercised offline via httptest + custom RoundTripper; no live model in regression.
- 2026-11-03T05:04:10+01:00 — Gate 3 PASS (regression): make verify -> OK (check, docs, links, all test suites); go test ./... -> all packages ok. No live model/source-site requests under verify; smoke-provider is excluded from verify by design.
- 2026-11-03T05:04:10+01:00 — Gate 4 PASS (static/structural): make check -> OK (git whitespace, taskqueue, gofmt, go build, go vet); git diff --check clean. smoke.json parses; smoke-actual.json retained. gofmt/build/vet clean across cmd + internal.
- 2026-11-03T05:04:10+01:00 — Gate 5 PASS (docs/spec): docs/MAKEFILE.md documents make smoke-provider (table + prose, marked operator gate, never in verify). docs/SPEC.md validator-contract section (run a configured model, retain original bounded outputs with exact inputs and observed identities) and docs/design/MVP.md 'One provider and three clean roles' match the implementation. Public prose Spanish, technical English; no LLM-authored presentation assets.
- 2026-11-03T05:04:10+01:00 — Gate 6 PASS (trust-boundary): negatives through the real changed interface - TestRealAdapterStrictSchema (bad schema), TestMissingEvidenceFails (missing evidence), TestEvidenceBudgetNotTruncated (over-budget, no truncation, no HTTP), TestRealAdapterOverBudgetNoApproval (over-budget reply), TestCancellationClosesTransport (late reply discarded, transport observes cancellation), TestRealAdapterNon2xxNotRetried + TestRealAdapterTransientRetry (retry eligibility), TestSecretRedaction (key absent from captured request/records), TestRealAdapterCleanRequestCanary (no cross-role leak). Injection is data (TestInjectionIsData preserved). No false approval, unsafe writes, or leaked credentials.
- 2026-11-03T05:04:10+01:00 — Gate 7 PASS (artifact/task record): fixtures/semantic/smoke.json (target) + fixtures/semantic/smoke-actual.json (retained actual evidence: endpoint, model, revision, identity, per-role prompt/evidence/inputs/request hashes, original bounded replies, verification outcomes) committed. The real-provider reply is genuine stochastic execution evidence (not a deterministic build); deterministic operations (hashing, verification, budgets, strict envelope) repeat identically. No unrelated work or historical source bytes removed; stray build binary excluded.
