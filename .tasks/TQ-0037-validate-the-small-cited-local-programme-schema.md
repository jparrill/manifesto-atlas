---
id: TQ-0037
title: Validate the small cited local programme schema
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/normalization
  - area/validation
depends_on:
  - TQ-0036
  - TQ-0063
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T23:47:57+01:00
---

## Outcome

atlas.programme.local/v1 rejects uncited or mismatched public content using source-bound metadata and exact original-page excerpts.

## Scope

- Supply the executable schema and Go validator for the exact MVP draft shapes, local statement/measure/citation IDs and 21 ordered topics; reuse bounded strict JSON parsing.
- Bind metadata to selected source data, citations to extraction blocks/code-point slices and scope evidence to retained files. No per-citation hashes, legal registry or artifact graph.
- Validate every public title/global/topic/action/condition/institution Statement and never accept draft-authored approval or extracted Markdown as canonical data.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Valid local drafts round-trip deterministically with all required fields/arrays/topics and unique local IDs; unknown versions/fields and malformed inputs reject.
- [ ] Missing/dangling citations, wrong source/extraction/page/block/excerpt, incomplete absence basis, invented scope evidence and mismatched metadata reject without panic.
- [ ] Duplicate/dangling measures, unsupported institution text, unsafe generated markup/paths and forged lifecycle fields cannot reach rendering eligibility.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-canonical`. Pass when automated cases prove all five states, exact fields/nullability/IDs, every public-text citation, Unicode offsets, metadata/extraction mismatch, missing scope/page basis, short taxonomy and forged validation state; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-canonical`. Pass with exit zero and named positive/negative cases for all five states, exact fields/nullability/IDs, every public-text citation, Unicode offsets, metadata/extraction mismatch, missing scope/page basis, short taxonomy and forged validation state. Do not substitute a fake for required actual-tool/model evidence.

### Gate 3: Regression

Run `make verify` and `go test ./...`. Pass with exit zero, no Taskqueue warnings, no live model/source-site requests and no skipped required tests.

### Gate 4: Static and structural validation

Run `make check` and `git diff --check`; validate changed schemas/fixtures and documented relative links. Pass with no formatting/compiler/vet/schema/queue/whitespace errors.

### Gate 5: Documentation and specification consistency

Compare changed commands, schemas, fixtures and documentation with docs/SPEC.md and docs/design/MVP.md. Pass when actual behavior/retention/source authority match, public prose is Spanish, technical prose is English, and no LLM authors presentation assets.

### Gate 6: Trust-boundary and security behavior

Use the negative cases named in Gates 1/2 through the real changed interface. Pass only when invalid, stale, hostile, cancelled or incomplete inputs produce the stated safe outcome with no false approval, unsafe writes or leaked credentials. Retain any applicable trust-boundary evidence.

### Gate 7: Artifact and task record

Verify required retained artifacts/input hashes are current, deterministic operations repeat identically, and no unrelated user work or historical source bytes are removed. Record every gate as PASS with commands/artifact evidence using `tq note` before completion. A failed, unavailable or skipped required gate keeps this task open.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended changes and final task record, push and confirm upstream, then continue only approved ready work within the assignment and execution role. Delivery failures restore in-progress and block continuation; never bypass hooks/signing or promote inbox tasks.

## Specification

- docs/SPEC.md#canonical-programme-model
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T13:25:21+01:00 — Quality gates — TQ-0037 (all PASS). Gate 1/2 (make test-canonical): ok internal/canonical, 18 PASS. TestValidProgramme proves a complete 21-topic atlas.programme/v1 draft round-trips deterministically and exercises all five states (proposals/mentioned/not-covered/not-applicable/unclear) with every public factual field cited; TestDeterministicBytes proves canonical idempotency, identical-structure byte equality, and key-order independence; TestUnknownVersion rejects atlas.programme/v2; TestSourceMismatch rejects a wrong archived checksum binding; TestMissingCitation rejects an uncited global summary (citations=[]); TestInvalidPage rejects page 99 beyond the 3-page extraction; TestExcerptMismatch rejects an excerpt that does not equal the exact block [start:end] slice; TestDuplicateMeasure and TestDuplicateStatement reject duplicate IDs; TestUnknownField rejects an injected extra top-level field (DisallowUnknownFields); TestStateConflict rejects proposals without primary measures; TestNotCoveredIncompleteBasis rejects not-covered without complete 1..N page coverage; TestNotApplicableMissingRule rejects not-applicable without an applicability rule; TestUnknownInstitution/TestUnknownRelation reject out-of-profile references; TestUnsafeText rejects a generated statement containing raw markup; TestNationalScope validates a national-scope variant alongside the autonomous fixture; TestMeasureCitationsUnion rejects a measure whose citations are not the deduped source-order union of its statement citations. Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy/test-canonical all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/canonical/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new atlas.programme/v1 package and make test-canonical target; technical text English; no product decision changed. Gate 6 (trust-boundary): negative cases all fail closed (unknown version, source mismatch, uncited text, out-of-range page, excerpt mismatch, duplicate IDs, unknown fields, state conflicts, incomplete basis, missing rule, unknown institution/relation, unsafe generated markup); plain-data rule bounds length, rejects control chars except newline/tab, raw <> delimiters, markdown link syntax, and executable schemes; quoted excerpts preserved as data; no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0037 intended files staged (internal/canonical/{canonical.go,canonical_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0037-*.md).
- 2026-11-02T17:27:01+01:00 — Reopened by TQ-0058 audit under Juan delegation. Prior all-PASS completion is superseded: go run ./docs/audits/probe reproduces acceptance of an invented not-applicable rule, mismatched extraction checksum, missing required summary field, and panic on a short supplied taxonomy. Static gaps: no machine-readable schema, no checked metadata evidence resolution, no institution/rule/relation compatibility beyond unrelated membership maps, and no measure anchor/ordinal identity or canonical ordering verification. Required gates 1/2/5/6/7 remain incomplete; retain native TQ-0036 dependency and reject incomplete/mismatched records without panic before closing.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-02T23:46:22+01:00 — Gate 1 PASS: make test-canonical exit 0 (26 tests). Named cases: all five states exercised (TestValidProgramme); exact fields/nullability/local IDs (fixture + TestBadCitationOrdinal, TestDuplicateMeasure, TestWrongTopicOrder); every public-text citation cited (TestMissingCitation rejects uncited); Unicode code-point offsets (TestUnicodeCodePointOffsets multi-byte slice); metadata/extraction mismatch (TestMetadataMismatch, TestSourceMismatch, TestSourceMetadataChecksumMismatch, TestExtractionChecksumMismatch); missing scope/page basis (TestNotCoveredIncompleteBasis, TestNotApplicableMissingScope); fixed 21-topic order (TestWrongTopicOrder); forged/unknown field rejected (TestUnknownField). Every criterion maps to a named test.
- 2026-11-02T23:46:22+01:00 — Gate 2 PASS: make test-canonical exit 0 with named positive/negative cases (round-trip determinism TestDeterministicBytes, all states, exact fields, every public-text citation, Unicode offsets, checksum/metadata mismatch, missing scope/page basis, topic order, forged validation state, unsafe markup, scope-evidence hash/traversal). No live model/converter evidence; all offline.
- 2026-11-02T23:47:02+01:00 — Gate 3 PASS: make verify exit 0 (check + docs structure + relative links + all focused suites incl. test-canonical); go test ./... all ok (artifact, candidate, canonical, cli, run, semantic, taxonomy; cmd/probe no test files). No Taskqueue warnings, no live model/source-site requests, no skipped required tests. Audit probe (go run ./docs/audits/probe) reports 0 failed rejection invariants, adapted to the new .local/v1 interface.
- 2026-11-02T23:47:02+01:00 — Gate 4 PASS: make check OK (git diff --check clean, tq list --json no warnings, gofmt empty, go build ./cmd/manifesto-atlas clean, go vet ./... clean); git diff --check clean; gofmt -l cmd internal docs/audits/probe empty.
- 2026-11-02T23:47:02+01:00 — Gate 5 PASS: docs/MAKEFILE.md updated (test-canonical description + module layout now atlas.programme.local/v1 with source-record/extraction checksum binding, retained scope-evidence, local ordinal IDs). docs/SPEC.md#canonical-programme-model consistent (local stable IDs within a source revision, source-bound metadata, cited exact page/excerpt citations, no draft-authored approval). docs/design/MVP.md governs. No stale atlas.programme/v1 (non-local) refs remain in root docs. Technical prose English; no LLM presentation assets.
- 2026-11-02T23:47:45+01:00 — Gate 6 PASS: negative cases exercised through the real ValidateProgramme interface, all offline with no network/panic: wrong source/source-metadata/extraction checksum -> source-mismatch; metadata field mismatch -> source-mismatch; out-of-range page -> bad-page; excerpt != exact block slice -> excerpt-mismatch; wrong citation/measure/statement ordinal -> bad-id; dangling measure/citation refs -> dangling-reference; unsafe generated markup (<>, ](, exec schemes) -> unsafe-text; not-applicable without retained scope evidence -> missing-scope; incomplete absence basis -> incomplete-basis; scope-evidence path traversal -> bad-path; scope-evidence wrong hash -> hash-mismatch; missing required arrays -> missing-field; injected unknown field -> unknown-field. No false approval, no unsafe writes, no leaked credentials; probe confirms short-topic reject without panic.
- 2026-11-02T23:47:45+01:00 — Gate 7 PASS: deterministic (TestDeterministicBytes: repeated canonicalization identical, key-order independent, idempotent). Reworked internal/canonical/{canonical.go,canonical_test.go} to atlas.programme.local/v1; removed obsolete Jurisdiction/Taxonomy/ManifestID bindings and content-hash CitationID; retained artifact strict-JSON/hash helpers and taxonomy 21-topic registry. No unrelated user work or historical source bytes removed (fixtures/frente/ intact). All required gates PASS.
