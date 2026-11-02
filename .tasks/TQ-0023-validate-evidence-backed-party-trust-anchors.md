---
id: TQ-0023
title: Validate evidence-backed party trust anchors
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/provenance
depends_on:
  - TQ-0010
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T13:52:25+01:00
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

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T13:52:16+01:00 — Quality gates — TQ-0023 (all PASS). Gate 1/2 (make test-party-registry): ok internal/provenance, 13 PASS. TestCheckAnchorResolves proves an approved exact domain + evidence chain resolves deterministically for a party and scope; TestRegionalDomain proves a linked regional/campaign hostname with its own grant resolves; TestUnknownParty rejects a party absent from the registry; TestLookalike rejects typo/subdomain/suffix/appended and www hostnames (no wildcard/substring/suffix inference); TestScopeMismatch rejects a wrong territorial scope; TestRevoked rejects a revoked grant; TestExpired rejects now > valid_until; TestNotYetReviewed rejects now < reviewed_at; TestEvidenceIntegrity rejects non-hash evidence; TestConflictingIdentity rejects a hostname owned by a different party; TestValidateRegistryValid accepts the fixture registry; TestValidateRegistryRejects proves unknown version, duplicate party, Unicode hostname (no IDNA guessing), cross-party hostname conflict, and over-180-day grant are rejected; TestVersionInvalidation proves a superseding registry that revokes/removes a grant invalidates dependent provenance while an unchanged grant does not, without rewriting historical records. Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy/test-canonical/test-party-registry all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/provenance/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new atlas.provenance/v1 package and make test-party-registry target; technical text English; no product decision changed. Gate 6 (trust-boundary): contributor cannot activate own authority (registry is caller-selected, grants require independent evidence hashes); Unicode/IDNA hostname ambiguity rejected; lookalike/suffix/subdomain not trusted; expiry/revocation/scope/evidence all fail closed; fixture anchors only, no unverified production facts, no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0023 intended files staged (internal/provenance/{party.go,party_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0023-*.md).
