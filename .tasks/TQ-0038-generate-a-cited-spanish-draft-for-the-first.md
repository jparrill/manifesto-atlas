---
id: TQ-0038
title: Generate a cited Spanish draft for the first real programme
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/normalization
  - area/skills
  - area/cli
depends_on:
  - TQ-0037
  - TQ-0024
  - TQ-0030
  - TQ-0031
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

make normalize runs the normalizer skill through the actual configured provider and stores a structurally valid cited Spanish draft for the first source.

## Scope

- Author the normalizer prompt/skill against the local schema and fixed topics; use skill-authoring guidance and writing-for-agents when available, without treating unavailable optional authoring advice as model evidence.
- Use complete bounded extraction and original evidence; preserve atomic actions, conditions, quantities, polarity, scope and institutions without editorial/viability conclusions.
- Retain actual prompt/schema/source/extraction/model/tool identities and raw reply; only a draft is installed, never a validator approval. One real source is sufficient for this slice.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] The first actual source produces a schema-valid Spanish draft with all topics and every factual public Statement cited to exact original evidence.
- [ ] Missing evidence, unknown schema/topic, stale source/extraction, malformed output, injected tool/approval instructions and uncited prose cannot install a usable draft.
- [ ] Real configured-model positive/negative evaluation outputs are retained; offline tests check those records without invoking the model.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-normalizer`. Pass when automated cases prove stored real draft, five states, atomic conditions, missing/uncited evidence, stale input, malformed output and prompt injection; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails. Also run `make evaluate-normalizer TARGET=fixtures/frente/source.json` and retain actual source/model/tool/input identities and original outputs; unavailable required resources keep the task open, never a skipped pass.

### Gate 2: Focused behavior

Run `make test-normalizer`. Pass with exit zero and named positive/negative cases for stored real draft, five states, atomic conditions, missing/uncited evidence, stale input, malformed output and prompt injection. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#agent-skills
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
