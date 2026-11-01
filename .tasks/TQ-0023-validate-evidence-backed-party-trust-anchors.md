---
id: TQ-0023
title: Validate evidence-backed party trust anchors
status: inbox
priority: high
labels:
  - type/feature
  - agent/developer
  - area/provenance
depends_on:
  - TQ-0010
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:30+01:00
---

## Outcome

Provenance resolves party-controlled domains through a versioned evidence-backed registry rather than arbitrary submitted URLs.

## Scope

- Implement the approved party/territory/domain evidence registry and authority/update validation.
- Represent expiry/revalidation and linked regional/campaign domains; include fixture anchors without asserting unverified production facts.
- Bind trust-registry versions to provenance reports and invalidate affected approvals when anchors change.
- Provide the focused Make gate `make test-party-registry` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Approved exact domains and evidence chains resolve deterministically for a party and scope.
- [ ] Lookalikes, conflicting identities, expired/revoked anchors, unauthorized updates, and unknown domains cannot yield verified provenance.
- [ ] Registry changes invalidate dependent evidence without rewriting historical records.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-party-registry`. Pass when automated cases prove: known/unknown party; regional domain; lookalike; expired/revoked anchor; unauthorized update; version invalidation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-party-registry`. Pass when automated cases prove: known/unknown party; regional domain; lookalike; expired/revoked anchor; unauthorized update; version invalidation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0023 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Treat anchor updates as authority changes; test contributor self-approval, Unicode/URL hostname ambiguity, and malicious external evidence references. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0023`:

1. Record all gate results with `tq note`, then run `tq done TQ-0023`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0023` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#provenance-trust-model`
- `docs/SPEC.md#further-notes`
- `docs/QUALITY_GATES.md`
