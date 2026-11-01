---
id: TQ-0011
title: Define immutable source metadata, storage, and revision identity
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - agent/architect
  - area/source-registry
depends_on:
  - TQ-0008
  - TQ-0010
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T00:25:33+01:00
---

## Outcome

Programme/election/geography identities and immutable revisions have one unambiguous storage contract.

## Scope

- Define metadata for party, programme, election date/type, year, scope, geography, jurisdiction, document type, provenance, retrieval, filenames, checksums, revisions, and processing state.
- Specify the year-first path hierarchy, portable safe IDs, duplicate checksum behavior, corrected publications, collision handling, and atomic registration.
- Resolve how the general Frente Obrero political document and regional VOX document are identified without inventing election dates or scope.
- Keep third-party publisher rights and future fulfilment/viability assessments separate from source facts.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Worked paths and records cover national, autonomous, municipal, European, and non-electoral documents.
- [ ] Committed bytes cannot be overwritten, while a corrected source produces a linked new revision.
- [ ] Duplicate content, conflicting metadata, missing election data, and stale checksums have explicit outcomes.
- [ ] docs/design/source-identity.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/source-identity.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/source-identity.md. Pass only when each of these cases has one explicit, internally consistent expected result: first source; exact duplicate; corrected revision; metadata collision; missing election; unsafe filename; modified committed bytes. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0011 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0011`:

1. Record all gate results with `tq note`, then run `tq done TQ-0011`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0011` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#immutable-source-storage`
- `docs/SPEC.md#assessment-separation`
- `docs/SPEC.md#source-registry-tests`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T00:25:33+01:00 — Quality gates:
  - Gate 1 PASS — source-identity.md defines complete versioned metadata/identity tuples, safe year-first storage, worked national/autonomous/municipal/European/non-electoral records, immutable linked revisions and atomic registration. Golden-source identification explicitly forbids invented election dates, URLs and hashes.
  - Gate 2 PASS — scenario table reviewed for first source, exact/global duplicate, corrected revision, conflicting metadata, missing election, unsafe filenames and changed committed bytes; each has one explicit result consistent with the record/path/transaction rules.
  - Gate 3 PASS — make verify passed all existing baseline and Go gates.
  - Gate 4 PASS — git diff --check; python3 /tmp/opencode/check-design.py docs/design/source-identity.md passed whitespace/UTF-8, links/anchors and structure/scenario checks; example paths follow the declared grammar with clearly labelled synthetic revision segments.
  - Gate 5 PASS — compared with SPEC immutable-source-storage, assessment-separation, source-registry-tests and existing artifact/provenance/isolation contracts. Source facts remain separate from future assessments, third-party rights retained, evidence required, no speculative real golden metadata. English technical document, SPEC/public files unchanged.
  - Gate 6 PASS — reviewed current provenance authority, no-follow paths, portable slugs, lock/no-replace installation, interrupted writes, checksums and trusted-base committed-byte checks. Contributor filenames cannot control paths, stale reports cannot authorize registration. No private paths or secrets embedded.
  - Gate 7 PASS — current design and task evidence recorded; same inputs yield the same identity/path and exact duplicates are no-ops. Only the design and this task record changed; implementation/race/real-source recovery tests remain explicit dependent-task gates.
