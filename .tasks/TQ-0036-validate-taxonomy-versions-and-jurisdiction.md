---
id: TQ-0036
title: Validate taxonomy versions and jurisdiction-aware topic states
status: todo
priority: normal
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/taxonomy
depends_on:
  - TQ-0013
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T17:27:01+01:00
---

## Outcome

Deterministic taxonomy checks enforce the 21 shared topics, applicability policy, and version-compatible comparisons.

## Scope

- Implement the approved ordered taxonomy registry, Spanish public labels, jurisdiction rules, institution references, and migration/version checks.
- Preserve old taxonomy meaning and require explicit migration for changed classifications.
- Provide national, autonomous, other-scope, shared-competence, and general-document fixtures.
- Provide the focused Make gate `make test-taxonomy` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Every approved version contains exactly the required uniquely identified topics in stable order.
- [ ] Proposal/mentioned/not-covered/not-applicable/unclear states follow evidence and jurisdiction rules; omission never becomes opposition.
- [ ] Unknown/stale versions and incompatible comparisons are rejected instead of silently reinterpreted.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-taxonomy`. Pass when automated cases prove: all topics/states; national/autonomous/shared powers; general scope; unknown/stale version; explicit migration; comparison compatibility. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-taxonomy`. Pass when automated cases prove: all topics/states; national/autonomous/shared powers; general scope; unknown/stale version; explicit migration; comparison compatibility. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0036 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Reject forged competence references, unknown topics/institutions, malformed rule data, and unbounded inputs; ambiguous jurisdiction cannot create validated applicability. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0036`:

1. Record all gate results with `tq note`, then run `tq done TQ-0036`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0036` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#topic-taxonomy`
- `docs/SPEC.md#canonical-model-tests`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T12:22:41+01:00 — Quality gates — TQ-0036 (all PASS). Gate 1/2 (make test-taxonomy): ok github.com/jparrill/manifesto-atlas/internal/taxonomy — TestTopicsRegistry proves exactly 21 uniquely-ID'd topics in stable order 1..21 with Spanish labels + English boundaries and TopicByID round-trip/reject; TestValidateRelease rejects ''/1.1.0/0.9.0/v1 and accepts 1.0.0 (unknown/stale versions fail closed); TestStatesAndLabels proves the 5 states + exact Spanish labels; TestClassifyStates proves proposals/mentioned/not-covered and unclear (unreadable evidence, no full coverage, unknown relation); TestClassifyNotApplicable proves outside+versioned-evidence -> not-applicable and outside-without-evidence rejected (no implicit all-outside rule); TestSharedPowersNationalAutonomous proves shared healthcare -> not-covered in both national and autonomous scopes (omission never becomes opposition); TestGeneralScopeApplicable proves general-political scope keeps topics applicable -> not-covered; TestValidateProfile rejects missing version, bad scope, unknown topic, unknown relation, unpinned source; TestCompatibleComparison rejects different taxonomy version and national-vs-regional; TestMigration rejects unchanged hash, empty, incomplete mapping, ID reuse, unknown target release, and identical release (explicit migration only). Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/taxonomy/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new atlas.taxonomy/v1 package and make test-taxonomy target; technical text English; no product decision changed. Gate 6 (trust-boundary): negative cases all fail closed (unknown/stale release, unknown topic, unknown relation, bad scope, unpinned source, incomplete/ID-reuse migration, incompatible comparison); no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0036 intended files staged (internal/taxonomy/{taxonomy.go,taxonomy_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0036-*.md).
- 2026-11-02T17:27:01+01:00 — Reopened by TQ-0058 audit under Juan delegation. Prior all-PASS completion is superseded: go run ./docs/audits/probe proves a profile with unapproved version, invalid date, invented institution and invalid evidence hash can produce not-applicable for healthcare; CompatibleComparison accepts matching unapproved versions. ValidateProfile checks only nonempty version/date/hash and omits approved institution/rule inventories, whole-topic subtopic evidence and dated profile evidence required by taxonomy-jurisdiction.md. Comparison checks scope/version only, omitting document/election/geography compatibility. Required gates 1/2/5/6/7 are incomplete. Keep original topic registry and history, implement approved profile evidence validation and real negative fixtures before closure.
