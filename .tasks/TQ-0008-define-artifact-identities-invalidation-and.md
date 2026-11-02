---
id: TQ-0008
title: Define artifact identities, invalidation, and attestation trust
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - agent/architect
  - area/validation
  - area/source-registry
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

Every stage can identify exact inputs, verify report authority, and reject stale artifacts using one agreed contract.

## Scope

- Define the versioned artifact graph, IDs, dependency manifests, serialization/hashing boundaries, timestamps, model/prompt/tool/schema/taxonomy identity, and byte versus semantic reproducibility.
- Specify report signing, trusted public-key roots, signer authorization, key rotation/revocation, and offline verification; signatures must not substitute for semantic validation.
- Define statuses, retained failure reports, idempotency, and downstream invalidation for every changed input, including models and extraction model artifacts.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Worked manifests cover provenance, extraction, audit, draft programme, claim/coverage reports, renderings, site reports, and publication.
- [ ] One changed source, extraction, schema, taxonomy, prompt, model, tool, template, or QA configuration invalidates exactly the dependent artifacts.
- [ ] Unknown versions, forged/untrusted/revoked signatures, missing dependencies, and graph cycles have explicit non-publishable outcomes.
- [ ] docs/design/artifact-contract.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/artifact-contract.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/artifact-contract.md. Pass only when each of these cases has one explicit, internally consistent expected result: identical inputs; timestamp-only changes; each changed input class; missing edge; cyclic graph; forged report; unknown signer; rotated/revoked key. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0008 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0008`:

1. Record all gate results with `tq note`, then run `tq done TQ-0008`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0008` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T00:13:44+01:00 — Juan explicitly assigned ready architect work in the current conversation. Executing model: gpt-6.1-sol (the approved architecture model in AGENTS.md). Starting with the shared artifact contract so dependent implementation can use an explicit identity, invalidation, and offline authority boundary.
- 2026-11-02T00:18:00+01:00 — Design decisions: use strict bounded integer-preserving canonical JSON and content-addressed manifest/payload identities; keep execution timestamps outside artifact hashes but retrieval timestamps inside evidence; derive content lifecycle from approval graph to avoid cycles; use detached Ed25519 attestations with explicit offline role-authorized trust roots, current-time expiry checks and retroactive revocation. The worked manifest recipe uses reproducible synthetic label bytes and never claims official-source or domain-schema approval. No runtime or presentation implementation is included.
- 2026-11-02T00:18:16+01:00 — Quality gates:
  - Gate 1 PASS — docs/design/artifact-contract.md: Worked manifests across the pipeline defines the exact envelope construction and all 17 derived kinds, including provenance/extraction/audit/programme/claim/coverage/Markdown/HTML/site/publication; required edges and reverse-reachability scenarios cover every requested input class. Manifest/storage, attestation, and scenario sections explicitly reject unknown versions, forgery, unknown/revoked signers, missing dependencies and cycles. Decisions, interfaces, versioning and rejected alternatives are explicit.
  - Gate 2 PASS — inspected the scenario table against the edge table: identical inputs, bookkeeping versus evidence timestamps, each changed source/extraction/schema/taxonomy/prompt/model/tool/template/QA input, missing/unknown edges, cycles, forged and unknown keys, rotation overlap, expiry and retroactive revocation all have one fail-closed or precisely scoped invalidation result. Worked manifests use exact synthetic label bytes; domain acceptance remains separate.
  - Gate 3 PASS — make verify completed successfully (Git/Taskqueue, gofmt/build/vet, documentation checks and CLI tests); tq show TQ-0008 --json exited zero without warnings.
  - Gate 4 PASS — git diff --check; Python stdlib structural diagnostic confirmed all 9 required sections, working relative links AND heading anchors, the exact empty-byte SHA-256 JSON example, no trailing whitespace, worked-manifest coverage for every one of the 17 derived kinds, and the scenario table. No structural/link findings.
  - Gate 5 PASS — compared with SPEC make-targets, llm-validation and publication-and-ci: exact input invalidation, separate claim/coverage contexts, retained failures, immutable PDFs, JSON source of truth, offline CI, no presentation-code authorship and no human fallback preserved. No existing design contract conflicts; technical text is English; SPEC/public text unchanged.
  - Gate 6 PASS — reviewed trust and negative scenarios: submitted evidence cannot supply report authority; signatures bind subject/role/time, verify against explicit trusted public roots, and never turn a failed/insufficient report into approval. Bounds, strict parsing, derived safe local paths, no-follow opens, no network fallback, atomic writes and credential isolation are specified. No private keys, secrets or private absolute paths are included.
  - Gate 7 PASS — design artifact and task evidence current; static document SHA-256 a2b74e5bc6bb92f6def13c8b10cbdc471f147ef12476687f6eb1f285d78f283f; repeated make verify passed. Initial checkout clean; only the design artifact and this task record changed. No runtime implementation or unrelated user work modified.
- 2026-11-02T18:05:04+01:00 — Historical completion retained. Juan's TQ-0061 scope replan establishes docs/SPEC.md and docs/design/MVP.md as the current MVP contract; conflicting graph/PKI/network admission/legal-registry/calibrated-score/browser-pinning requirements in this task's original design are superseded for MVP work. This note does not recertify any historical implementation or erase prior execution evidence. TQ-0063 owns removal/adaptation of inactive platform code.
