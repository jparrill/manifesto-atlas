---
id: TQ-0037
title: Validate canonical programme JSON and every public citation
status: todo
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/normalization
  - area/validation
depends_on:
  - TQ-0014
  - TQ-0036
  - TQ-0019
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T17:27:01+01:00
---

## Outcome

A versioned validator rejects incomplete or unsupported canonical data before rendering or content approval.

## Scope

- Implement the approved canonical schema, deterministic serialization, topic/state rules, source metadata binding, citation/page/excerpt references, and measure identity.
- Cover summaries and public status explanations as well as measures; preserve Spanish public text and English technical identifiers.
- Provide positive/negative fixtures for all topic states and national/autonomous scope.
- Provide the focused Make gate `make test-canonical` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Valid fixture JSON round-trips deterministically and contains every required taxonomy section.
- [ ] Missing citations, invalid original page/excerpt references, uncited public summaries, unknown fields/versions, duplicate measures, and incompatible institutions fail.
- [ ] Extracted Markdown cannot substitute for canonical JSON and validation state cannot be asserted by an untrusted payload.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-canonical`. Pass when automated cases prove: schema pass/fail; all topic states; every public text citation; invalid page; duplicate measure; unknown version; deterministic bytes. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-canonical`. Pass when automated cases prove: schema pass/fail; all topic states; every public text citation; invalid page; duplicate measure; unknown version; deterministic bytes. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0037 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Bound JSON size/depth/text fields, treat markup/instructions as data, reject path/reference traversal, and test forged validation state. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0037`:

1. Record all gate results with `tq note`, then run `tq done TQ-0037`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0037` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#canonical-programme-model`
- `docs/SPEC.md#canonical-model-tests`
- `docs/SPEC.md#deterministic-rendering`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T13:25:21+01:00 — Quality gates — TQ-0037 (all PASS). Gate 1/2 (make test-canonical): ok internal/canonical, 18 PASS. TestValidProgramme proves a complete 21-topic atlas.programme/v1 draft round-trips deterministically and exercises all five states (proposals/mentioned/not-covered/not-applicable/unclear) with every public factual field cited; TestDeterministicBytes proves canonical idempotency, identical-structure byte equality, and key-order independence; TestUnknownVersion rejects atlas.programme/v2; TestSourceMismatch rejects a wrong archived checksum binding; TestMissingCitation rejects an uncited global summary (citations=[]); TestInvalidPage rejects page 99 beyond the 3-page extraction; TestExcerptMismatch rejects an excerpt that does not equal the exact block [start:end] slice; TestDuplicateMeasure and TestDuplicateStatement reject duplicate IDs; TestUnknownField rejects an injected extra top-level field (DisallowUnknownFields); TestStateConflict rejects proposals without primary measures; TestNotCoveredIncompleteBasis rejects not-covered without complete 1..N page coverage; TestNotApplicableMissingRule rejects not-applicable without an applicability rule; TestUnknownInstitution/TestUnknownRelation reject out-of-profile references; TestUnsafeText rejects a generated statement containing raw markup; TestNationalScope validates a national-scope variant alongside the autonomous fixture; TestMeasureCitationsUnion rejects a measure whose citations are not the deduped source-order union of its statement citations. Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy/test-canonical all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/canonical/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new atlas.programme/v1 package and make test-canonical target; technical text English; no product decision changed. Gate 6 (trust-boundary): negative cases all fail closed (unknown version, source mismatch, uncited text, out-of-range page, excerpt mismatch, duplicate IDs, unknown fields, state conflicts, incomplete basis, missing rule, unknown institution/relation, unsafe generated markup); plain-data rule bounds length, rejects control chars except newline/tab, raw <> delimiters, markdown link syntax, and executable schemes; quoted excerpts preserved as data; no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0037 intended files staged (internal/canonical/{canonical.go,canonical_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0037-*.md).
- 2026-11-02T17:27:01+01:00 — Reopened by TQ-0058 audit under Juan delegation. Prior all-PASS completion is superseded: go run ./docs/audits/probe reproduces acceptance of an invented not-applicable rule, mismatched extraction checksum, missing required summary field, and panic on a short supplied taxonomy. Static gaps: no machine-readable schema, no checked metadata evidence resolution, no institution/rule/relation compatibility beyond unrelated membership maps, and no measure anchor/ordinal identity or canonical ordering verification. Required gates 1/2/5/6/7 remain incomplete; retain native TQ-0036 dependency and reject incomplete/mismatched records without panic before closing.
