---
id: TQ-0063
title: Replace inactive platform prerequisites with local programme run records
status: done
priority: high
labels:
  - type/feature
  - area/project
  - agent/developer
depends_on:
  - TQ-0061
created: 2026-11-02T18:05:04+01:00
updated: 2026-11-02T22:14:17+01:00
---

## Outcome

One local programme run can bind files, reports and output hashes without report PKI, a content-addressed DAG or a custom source-network platform.

## Scope

- Implement the small atlas.selection/v1 and atlas.run/v1 schemas/readers from the MVP contract; use declared rooted local files and preserve useful strict JSON/hash helpers.
- Remove inactive attestation, arbitrary-URL transport and party-grant code/tests, and unused manifest/kind/graph machinery after checking all callers. Adapt existing audit probes and documentation honestly; preserve historical audit reports/task notes.
- Update CLI/Make stub ownership to active MVP tasks and deferred-stage messages; use go test ./... for global regression and retain meaningful focused targets. Do not remove retained canonical/semantic negative cases as a shortcut.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Valid local run records round-trip deterministically and expose actual current input/report/output references without a trusted approval boolean.
- [ ] Missing/tampered/unsafe references, unknown fields/versions, duplicate selections and stale hashes reject without network lookup or panic.
- [ ] No production/stub/dependency references require retired platform packages; source bytes and historical task/audit evidence remain intact.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-runs`. Pass when automated cases prove selection/run round-trip, stale file mutation, missing reference, symlink/traversal, duplicate source selection, malformed JSON and unavailable-stage ownership; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-runs`. Pass with exit zero and named positive/negative cases for selection/run round-trip, stale file mutation, missing reference, symlink/traversal, duplicate source selection, malformed JSON and unavailable-stage ownership. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#module-seams
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:05:04+01:00 — Created directly in todo under Juan's explicit TQ-0061 request to revise existing tasks and create the missing MVP outcomes. Outcome, scope, criteria, seven gates, priority, dependencies and execution role are declared. No authority to promote unrelated inbox proposals is granted.
- 2026-11-02T22:13:12+01:00 — Gate 1 PASS: make test-runs exit 0. Named cases: TestSelectionRoundTrip + TestRunRoundTrip (deterministic round-trip, null refs preserved); TestParseSelectionRejects (duplicate programme_id->duplicate, malformed->bad-json, v2->bad-version, bad hash->bad-hash, unknown field->bad-field); TestParseRunRejectsUnknownStage (unknown input stage->bad-field, malformed->bad-json, v2->bad-version, bad ref hash->bad-hash, ../escape.json->bad-path); TestVerifyRunStaleMutation->hash-mismatch; TestVerifyRunMissingReference->missing-file; TestVerifyRunSymlink->symlink. Unavailable-stage ownership proven by make test-cli (build stage -> exit 1, 'not available yet', TQ-0043) plus the updated Makefile/CLI ownership map matching docs/SPEC.md#make-targets.
- 2026-11-02T22:13:12+01:00 — Gate 2 PASS: make test-runs exit 0 with named positive/negative cases (round-trip, stale mutation, missing ref, symlink, traversal, duplicate, malformed JSON, bad version/hash/field). No fake model/converter evidence substituted; all offline.
- 2026-11-02T22:13:42+01:00 — Gate 3 PASS: make verify exit 0 (check + docs structure + relative links + all focused suites incl. test-runs); go test ./... all ok (artifact, candidate, canonical, cli, run, semantic, taxonomy; cmd/probe no test files). No Taskqueue warnings, no live model/source-site requests, no skipped required tests.
- 2026-11-02T22:13:42+01:00 — Gate 4 PASS: make check OK (git diff --check clean, tq list clean, gofmt empty, go build ./cmd/manifesto-atlas clean, go vet ./... clean); git diff --check clean; gofmt -l cmd internal empty.
- 2026-11-02T22:13:42+01:00 — Gate 5 PASS: docs/MAKEFILE.md updated (baseline target table, Go-tests semantics, module layout, MVP ownership map, post-TQ-0063 migration wording); Makefile + internal/cli stage ownership now match docs/SPEC.md#make-targets (verify-source/add TQ-0027, validate-extraction TQ-0030, build TQ-0043/0044/0049, publish-check TQ-0050/0049, serve TQ-0065). SPEC.md already declares arbitrary-URL/attestation/domain-key-grants/semantic-provenance as deferred (104-107) and local-only acquisition (243-247); no conflict. Historical audit/design docs preserved. Technical prose English, no LLM presentation assets.
- 2026-11-02T22:14:04+01:00 — Gate 6 PASS: negative cases exercised through the real interface with no network lookup (VerifyRun/Parse* are local file ops + strict JSON only): traversal ../escape.json->bad-path, symlink ref->symlink, stale mutation->hash-mismatch, absent file->missing-file, duplicate programme_id->duplicate, malformed JSON->bad-json, wrong version->bad-version, unknown field/stage->bad-field, malformed hash->bad-hash. No false approval, no unsafe writes, no leaked credentials, no panic.
- 2026-11-02T22:14:04+01:00 — Gate 7 PASS: run/selection Marshal deterministic (repeat marshal identical, tested). Removed only internal/attestation, internal/provenance, and unused artifact manifest/kind/graph machinery after caller check (canonical/semantic use only Canonicalize/CanonicalValue, both retained); retained canonical/semantic/candidate/taxonomy negative cases. fixtures/frente/ source bytes + historical task/audit evidence intact; no unrelated user work removed. All required gates PASS.
