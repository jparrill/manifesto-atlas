---
id: TQ-0020
title: Sign reports and verify authorized offline attestations
status: inbox
priority: high
labels:
  - type/feature
  - agent/developer
  - area/validation
  - area/provenance
depends_on:
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:30+01:00
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
