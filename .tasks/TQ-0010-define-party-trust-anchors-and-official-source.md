---
id: TQ-0010
title: Define party trust anchors and official-source evidence rules
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - agent/architect
  - area/provenance
depends_on:
  - TQ-0008
  - TQ-0009
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T00:23:48+01:00
---

## Outcome

Official status follows an auditable trust chain rather than branding, filenames, or model opinion.

## Scope

- Define party registry identity, evidence authority, approved domains, regional/campaign domains, expiry/revalidation, and anchor updates.
- Specify direct-domain, exact official-page referral/CDN, and verifiable publisher-signature chains, including URL normalization, redirects, archived evidence, and third-party hosts.
- Define verified/probable/unverified/rejected decision tables and checksum mismatch treatment; resolve removed pages, new parties, and social-only publication without weakening verified status.
- Define safe network policy: schemes, DNS/IP checks, redirects, credential-bearing URLs, response limits, and canonical byte retrieval.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] All three permitted trust chains have required evidence and explicit failure outcomes.
- [ ] A contributor cannot approve their own arbitrary domain or document through self-asserted metadata.
- [ ] Candidate-versus-canonical checksum mismatch is rejected and incomplete chains never become publishable.
- [ ] docs/design/provenance-trust.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/provenance-trust.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/provenance-trust.md. Pass only when each of these cases has one explicit, internally consistent expected result: direct official PDF; exact official referral; valid/invalid publisher signature; lookalike host; hostile redirect; absent trust anchor; removed referral; changed canonical bytes. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0010 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0010`:

1. Record all gate results with `tq note`, then run `tq done TQ-0010`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0010` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#provenance-trust-model`
- `docs/SPEC.md#further-notes`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T00:23:48+01:00 — Quality gates:
  - Gate 1 PASS — docs/design/provenance-trust.md defines all three admissible chains, independent canonical-byte retrieval and mandatory candidate checksum equality, independently reviewed scoped/expiring party anchors, deterministic decision rules, interfaces, versions and rejected alternatives.
  - Gate 2 PASS — inspected scenario table and chain/network sections: direct PDF, exact referral, authorized/invalid publisher signatures, lookalikes, hostile redirects, missing anchors, removed referrals and silently changed canonical bytes each have explicit internally consistent outcomes. Positive incomplete evidence is probable, absent evidence unverified, contradictions rejected; none publishes.
  - Gate 3 PASS — make verify completed successfully, retaining all baseline and Go gates.
  - Gate 4 PASS — git diff --check and python3 /tmp/opencode/check-design.py docs/design/provenance-trust.md passed UTF-8, whitespace, section/table presence, relative links and heading-anchor checks.
  - Gate 5 PASS — compared with SPEC provenance-trust-model/further-notes plus artifact/isolation designs: three permitted chains preserved, uploaded bytes never canonical by assertion, no fabricated golden URLs/anchors, no manual publication fallback. Signature format support explicitly versioned; unsupported formats fail closed. Technical English, no public/SPEC changes.
  - Gate 6 PASS — trust-boundary design reviewed: contributor cannot self-approve grants, exact scoped hosts/keys, independent evidence, current expiry/revocation, no model network, safe TLS/DNS/dial pinning/redirects and bounded response handling. Archived evidence cannot bootstrap new official status. No secrets/private paths embedded.
  - Gate 7 PASS — design and scenario artifacts current; immutable text and canonical identity rules deterministic; only design/task record changed. Actual transport/anchor/signature runtime checks are mandatory dependent implementation gates, not claimed executed by this design task.
