---
id: TQ-0014
title: Define canonical programme schema and citation semantics
status: todo
priority: high
labels:
  - type/docs
  - agent/architect
  - area/normalization
  - area/taxonomy
  - area/validation
depends_on:
  - TQ-0008
  - TQ-0011
  - TQ-0013
  - TQ-0012
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T18:56:22+01:00
---

## Outcome

One versioned canonical JSON contract can drive public summaries, validation, and both renderers.

## Scope

- Define normalized metadata, all topic states, Spanish summaries, atomic measures, responsible institutions, jurisdictions, citations/excerpts, state, versions, and artifact hashes.
- Specify evidence for every public factual statement, including summaries and topic explanations, not only measure text.
- Define original PDF page references, multiple citations, stable measure IDs, duplicates, ordering, length/size limits, unknown fields, and deterministic serialization.
- Keep draft/validating/validated/published transitions separate from editorial opinions, future scores, and source provenance states.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Representative valid and invalid examples cover all five topic states and every public text field.
- [ ] Missing/invalid citations, incompatible taxonomy, unsupported institution, and missing topic sections are rejected.
- [ ] JSON is the sole public-content source; extracted Markdown is never its rendering contract.
- [ ] docs/design/canonical-programme.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/canonical-programme.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/canonical-programme.md. Pass only when each of these cases has one explicit, internally consistent expected result: valid programme; uncited summary/measure; invalid PDF page; duplicate measure ID; missing topic; unknown schema; unsafe embedded markup. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0014 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0014`:

1. Record all gate results with `tq note`, then run `tq done TQ-0014`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0014` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#canonical-programme-model`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#deterministic-rendering`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
