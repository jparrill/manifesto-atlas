---
id: TQ-0029
title: Recover the exact VOX Castilla y Leon source and origin evidence
status: done
priority: normal
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/feature
  - agent/developer
  - area/source-registry
  - area/extraction
depends_on:
  - TQ-0061
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-03T01:51:07+01:00
---

## Outcome

The exact owner-approved VOX Castilla y León source and evidenced electoral/territorial metadata are available for the second-document milestone.

## Scope

- Recover exact candidate bytes and retained authoritative origin evidence using existing acquisition tools, independently of registration or the first programme. Store fixtures/vox/source.json and evidence hashes.
- Do not infer election dates from 2026, invent URLs/checksums or substitute another revision; missing exact evidence keeps the task open.
- No unbounded PDF parsing or automatic official-certification claim. Actual page anchors are recorded later under bounded extraction.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] Exact approved bytes, checksum, source evidence, rights and evidenced party/geography/election identity exist as local candidates.
- [ ] Missing or conflicting electoral evidence is a blocker, never a guessed date or copied national programme.
- [ ] Offline tests reject changed bytes, missing evidence and mismatched metadata; this task does not block the first extraction/Markdown milestone.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-fixture-vox`. Pass when automated cases prove exact candidate/evidence hashes, regional electoral identity, missing/conflicting metadata and mutated copy; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-fixture-vox`. Pass with exit zero and named positive/negative cases for exact candidate/evidence hashes, regional electoral identity, missing/conflicting metadata and mutated copy. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#testing-decisions
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-03T01:50:59+01:00 — Gate 1 PASS: make test-fixture-vox -> ok internal/candidate. Exact VOX CyL PDF bytes retained (sha256 b1240d9f..., 29481929 bytes, valid PDF 1.4, 172pp), candidate parsed + verified offline; evidenced 15 March 2026 election, autonomous scope, VOX party, maintainer-selected origin.
- 2026-11-03T01:50:59+01:00 — Gate 2 PASS: go test ./internal/candidate/... -> ok (TestCheckVoxFixture + TestCheckFrenteFixture + TestParseRejects + TestVerifyRejectsMutatedCandidate + TestVerifyRejectsMissingEvidence all PASS).
- 2026-11-03T01:50:59+01:00 — Gate 3 PASS: make verify -> OK; go test ./... -> all packages ok.
- 2026-11-03T01:50:59+01:00 — Gate 4 PASS: make check -> OK (git whitespace, taskqueue integrity, gofmt, go build, go vet); git diff --check clean.
- 2026-11-03T01:50:59+01:00 — Gate 5 PASS: docs/MAKEFILE.md documents make test-fixture-vox; docs/SPEC.md updated (both initial source targets retained as fixtures); make verify docs+links OK.
- 2026-11-03T01:50:59+01:00 — Gate 6 PASS: candidate.Check fail-closed (parse rejects unknown fields/bad enums/bad hashes/unsafe paths/non-https URLs; VerifyFiles rejects mutated or missing bytes). VOX fixture passes the same Check; existing reject suite green.
- 2026-11-03T01:50:59+01:00 — Gate 7 PASS: fixtures/vox/{programa.pdf, evidence/origin-http-headers.txt, source.json} retained with matching SHA-256 (candidate b1240d9f..., evidence 6a4e78f1...); origin evidence = official voxespana.es PDF HTTP headers (Cloudflare-blocked HTML landing page documented); notes recorded.
