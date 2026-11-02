---
id: TQ-0020
title: Sign reports and verify authorized offline attestations
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/provenance
depends_on:
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T11:20:57+01:00
---

## Outcome

Stored approvals are authenticated against explicit trusted signers without committing private credentials.

## Scope

- Implement the approved signing/verification envelope, trusted public-key configuration, authorized report roles, rotation/revocation, and exact manifest binding.
- Keep signing credentials outside the repository and semantic job capabilities; use test-only keys for fixtures.
- Expose offline verification for later source registration and publication checks.
- Provide the focused Make gate `make test-attestations` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Valid authorized signatures verify against the exact report and artifact inputs.
- [ ] Modified payloads, mismatched inputs, unauthorized/unknown/revoked keys, and malformed signatures are rejected.
- [ ] Verification succeeds without secret keys or network access; signing failure cannot leave an apparent approval.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-attestations`. Pass when automated cases prove: valid signature; tampering; wrong/unknown/revoked signer; rotation; missing credential; offline verification. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-attestations`. Pass when automated cases prove: valid signature; tampering; wrong/unknown/revoked signer; rotation; missing credential; offline verification. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0020 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Exercise signature substitution, role confusion, oversized envelopes, and unsafe key paths; logs/errors must not expose private key material. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0020`:

1. Record all gate results with `tq note`, then run `tq done TQ-0020`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0020` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#module-seams`
- `docs/SPEC.md#publication-and-ci`
- `docs/SPEC.md#llm-validation`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T11:20:48+01:00 — Quality gates:
  - Gate 1 PASS — AC1 (valid authorized signature verifies against exact report/artifact inputs) proven by TestSignVerifyRoundTrip: Sign -> ParseAttestation (raw envelope) -> VerifyAttestation succeeds, and the signed statement binds contract/subject/role/key_id/issued_at. AC2 (modified payloads, mismatched inputs, unauthorized/unknown/revoked keys, malformed signatures rejected) proven by TestTamperingRejected (subject/role/issued_at/key_id/signature flips), TestRoleMismatch, TestUnknownKey, TestRevokedKey, TestMalformedSignatureEncoding, and TestTimeBounds. AC3 (verification without secret keys or network; signing failure leaves no approval) proven by VerifyAttestation using only the public trust bundle (no private key referenced) and TestMissingCredential (nil/short key returns an error and no attestation).
  - Gate 2 PASS — make test-attestations: automated cases prove valid signature; tampering; wrong/unknown/revoked signer; rotation overlap; missing credential; offline verification.
  - Gate 3 PASS — make verify: git whitespace, taskqueue integrity, gofmt, go build, go vet, docs structure, relative links, go test (test-cli + test-artifacts + test-attestations) all pass.
  - Gate 4 PASS — gofmt -l clean; go build and go vet ./... clean via make check; git diff --check clean.
  - Gate 5 PASS — docs/MAKEFILE.md updated (test-attestations added to baseline surface, Go-tests bullet, internal/attestation layout); Makefile help lists the target; interfaces match the attestation/trust sections of docs/design/artifact-contract.md; technical text English; no public (Spanish) content changed.
  - Gate 6 PASS — signature substitution and role confusion rejected (TestTamperingRejected, TestRoleMismatch); envelope/bundle/key-count limits enforced (16KiB, 1MiB, 256 keys); Sign accepts key bytes only and never reads key paths from a report; all errors carry stable codes and never echo key material; no network used. No secrets or private absolute paths committed.
  - Gate 7 PASS — verification logic is deterministic over its public inputs; make test-attestations and make verify pass; all gate results recorded before done.
