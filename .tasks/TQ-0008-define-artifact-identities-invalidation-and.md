---
id: TQ-0008
title: Define artifact identities, invalidation, and attestation trust
status: inbox
priority: high
labels:
  - type/docs
  - agent/architect
  - area/validation
  - area/source-registry
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:29+01:00
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
