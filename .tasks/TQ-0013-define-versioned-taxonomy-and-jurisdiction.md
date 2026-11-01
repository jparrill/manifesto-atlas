---
id: TQ-0013
title: Define versioned taxonomy and jurisdiction applicability rules
status: todo
priority: high
labels:
  - type/docs
  - agent/architect
  - area/taxonomy
depends_on:
  - TQ-0011
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T18:56:21+01:00
---

## Outcome

All 21 policy topics have stable identities and defensible applicability rules across document scopes.

## Scope

- Define ordered topic IDs, Spanish public labels, English technical descriptions, cross-topic mapping rules, and taxonomy version/migration policy.
- Define proposal, mention, not-covered, not-applicable, and unclear classifications, with evidence requirements and no inferred opposition.
- Specify national/autonomous/provincial/municipal/European/general-document applicability, shared powers, institutions, and cross-jurisdiction proposals; cite authoritative competence evidence.
- Explain compatible comparison and preservation of historical classifications across taxonomy updates.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Exactly the 21 Spec topics have unique stable IDs and deterministic order.
- [ ] National and autonomous examples distinguish omission from non-applicability without blanket exclusions of shared competences.
- [ ] Unknown/stale taxonomy versions and ambiguous jurisdiction do not silently produce validated data.
- [ ] docs/design/taxonomy-jurisdiction.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/taxonomy-jurisdiction.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/taxonomy-jurisdiction.md. Pass only when each of these cases has one explicit, internally consistent expected result: concrete proposal; passing mention; omission; outside scope; shared competence; unclear evidence; national/autonomous pair; taxonomy migration. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0013 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0013`:

1. Record all gate results with `tq note`, then run `tq done TQ-0013`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0013` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#topic-taxonomy`
- `docs/SPEC.md#canonical-programme-model`
- `docs/SPEC.md#canonical-model-tests`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
