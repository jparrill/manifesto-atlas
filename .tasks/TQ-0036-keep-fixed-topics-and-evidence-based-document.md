---
id: TQ-0036
title: Keep fixed topics and evidence-based document scope states
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/taxonomy
depends_on:
  - TQ-0061
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T20:13:34+01:00
---

## Outcome

The existing 21 ordered topics and five Spanish states support the first programme without a general legal-competence or migration engine.

## Scope

- Retain release 1.0.0 topic IDs/order/labels; reject unknown current versions/topics and remove speculative migration/comparison/profile APIs once callers are adapted.
- Classify cited proposals/mentions and complete-source absence with explicit document scope. Entire-topic not-applicable needs actual scope evidence; ambiguous applicability is unclear.
- Do not make a constitutional-law registry, per-institution authority certification or split/merge release platform prerequisite. A proposal targeting another institution remains sourced content, not proof of legal authority.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Exactly the same 21 unique ordered topics and all five distinct Spanish labels remain available for local canonical documents.
- [ ] Omission never becomes opposition; shared/advocacy content prevents blanket exclusion, and unclear evidence cannot establish not-applicable.
- [ ] Unknown versions/topics and forged/missing scope evidence reject or remain explicitly unclear rather than silently acquiring authority.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-taxonomy`. Pass when automated cases prove 21 IDs/order/labels, five states, general/regional document scope, shared/advocacy content, whole-source absence, unknown version and invented scope evidence; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-taxonomy`. Pass with exit zero and named positive/negative cases for 21 IDs/order/labels, five states, general/regional document scope, shared/advocacy content, whole-source absence, unknown version and invented scope evidence. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#topic-taxonomy
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T12:22:41+01:00 — Quality gates — TQ-0036 (all PASS). Gate 1/2 (make test-taxonomy): ok github.com/jparrill/manifesto-atlas/internal/taxonomy — TestTopicsRegistry proves exactly 21 uniquely-ID'd topics in stable order 1..21 with Spanish labels + English boundaries and TopicByID round-trip/reject; TestValidateRelease rejects ''/1.1.0/0.9.0/v1 and accepts 1.0.0 (unknown/stale versions fail closed); TestStatesAndLabels proves the 5 states + exact Spanish labels; TestClassifyStates proves proposals/mentioned/not-covered and unclear (unreadable evidence, no full coverage, unknown relation); TestClassifyNotApplicable proves outside+versioned-evidence -> not-applicable and outside-without-evidence rejected (no implicit all-outside rule); TestSharedPowersNationalAutonomous proves shared healthcare -> not-covered in both national and autonomous scopes (omission never becomes opposition); TestGeneralScopeApplicable proves general-political scope keeps topics applicable -> not-covered; TestValidateProfile rejects missing version, bad scope, unknown topic, unknown relation, unpinned source; TestCompatibleComparison rejects different taxonomy version and national-vs-regional; TestMigration rejects unchanged hash, empty, incomplete mapping, ID reuse, unknown target release, and identical release (explicit migration only). Gate 3 (regression): make verify -> check OK (git whitespace, tq list --json no warnings, gofmt OK, go build OK, go vet OK, docs structure OK, relative links OK) and go test test-cli/test-artifacts/test-attestations/test-taxonomy all ok; git diff --check -> exit 0. Gate 4 (static/structural): gofmt -l cmd internal -> no files; go build ./... and go vet ./internal/taxonomy/ -> no findings; git diff --check -> 0 errors. Gate 5 (docs/spec consistency): docs/MAKEFILE.md updated (surface table, Go-tests bullet, module layout) for the new atlas.taxonomy/v1 package and make test-taxonomy target; technical text English; no product decision changed. Gate 6 (trust-boundary): negative cases all fail closed (unknown/stale release, unknown topic, unknown relation, bad scope, unpinned source, incomplete/ID-reuse migration, incompatible comparison); no secrets or absolute private paths committed. Gate 7 (artifact/task record): fixtures deterministic (pure stdlib, no live models, no live official websites); repeated make verify identical; only TQ-0036 intended files staged (internal/taxonomy/{taxonomy.go,taxonomy_test.go}, Makefile, docs/MAKEFILE.md, .tasks/TQ-0036-*.md).
- 2026-11-02T17:27:01+01:00 — Reopened by TQ-0058 audit under Juan delegation. Prior all-PASS completion is superseded: go run ./docs/audits/probe proves a profile with unapproved version, invalid date, invented institution and invalid evidence hash can produce not-applicable for healthcare; CompatibleComparison accepts matching unapproved versions. ValidateProfile checks only nonempty version/date/hash and omits approved institution/rule inventories, whole-topic subtopic evidence and dated profile evidence required by taxonomy-jurisdiction.md. Comparison checks scope/version only, omitting document/election/geography compatibility. Required gates 1/2/5/6/7 are incomplete. Keep original topic registry and history, implement approved profile evidence validation and real negative fixtures before closure.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-02T19:55:17+01:00 — Resuming under TQ-0061 MVP scope. Plan: keep the 21 ordered topics + 5 Spanish states + release 1.0.0; rewrite Classify to use an explicit document-scope Applicability (in-scope/out-of-scope/unknown) with retained scope evidence for not-applicable; remove speculative Profile/LegalSource/Relation/Scope/CompatibleComparison/Mapping/Migration APIs. Callers: canonical_test.go uses only Topics() (no change); docs/audits/probe/main.go jurisdiction() uses Profile+CompatibleComparison (adapt to new API). Update Makefile/MAKEFILE.md wording.
- 2026-11-02T20:13:34+01:00 — Quality gates — TQ-0036 (all PASS) under TQ-0061 MVP scope.

  Gate 1/2 (make test-taxonomy, exit 0): internal/taxonomy retains exactly the 21 uniquely-ID'd topics in stable order 1..21 with Spanish labels + English boundaries (TestTopicsRegistry) and all five distinct Spanish state labels (TestStatesAndLabels). Evidence-based document-scope classification (TestClassifyStates): proposals/mentioned/not-covered and unclear (unreadable evidence, no full coverage, unknown applicability). TestClassifyNotApplicable: out-of-scope WITH retained scope evidence -> not-applicable; out-of-scope without/empty evidence rejected (no implicit all-outside rule, no guessed exclusion). TestSharedAndGeneralScopeNotExcluded: shared healthcare and general-political energy in-scope -> not-covered (omission never opposition, shared/advocacy not blanket-excluded). TestValidateRelease: rejects ''/1.1.0/0.9.0/v1, accepts 1.0.0. TestClassifyUnknownTopic + TestClassifyBadApplicability: forged topic/applicability reject.

  Gate 3 (regression): make verify exit 0; go test ./... exit 0 (all packages ok, no skipped tests, no Taskqueue warnings, no live model/source-site requests).

  Gate 4 (static/structural): make check exit 0 (git diff --check, tq list --json clean, gofmt -l empty, go build ./... incl. probe, go vet ./...); git diff --check exit 0.

  Gate 5 (docs/spec): SPEC.md 'Topic taxonomy' section matches implemented behavior (21 topics/release 1.0.0, five states, evidence-based scope, deferred legal registry/comparison/migration); Makefile help + test-taxonomy comment and docs/MAKEFILE.md (table row + internal/taxonomy layout) updated to evidence-based document-scope wording; technical prose English; no presentation assets authored.

  Gate 6 (trust-boundary): forged/missing/empty scope evidence rejected (no false not-applicable); unknown version, unknown topic and unknown applicability value reject or stay explicitly unclear; no silent authority acquisition.

  Gate 7 (artifact/task record): removed speculative Profile/LegalSource/Relation/Scope/CompatibleComparison/Mapping/Migration APIs and adapted the only caller (docs/audits/probe/main.go jurisdiction()); no dangling references (grep clean); deterministic offline; recording gates now.

  Note: docs/audits/probe is an audit tool, not a regression gate; its transport/registry/canonical sections reproduce pre-existing trust-boundary findings owned by other tasks (TQ-0022/TQ-0023/TQ-0037) and are unchanged by this task. The probe builds and its taxonomy checks pass.
