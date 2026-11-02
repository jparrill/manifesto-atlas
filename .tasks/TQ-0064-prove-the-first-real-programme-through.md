---
id: TQ-0064
title: Prove the first real programme through validated Spanish Markdown
status: todo
priority: high
labels:
  - type/feature
  - area/validation
  - agent/developer
depends_on:
  - TQ-0028
  - TQ-0051
  - TQ-0053
created: 2026-11-02T18:05:04+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

The approved Frente Obrero source produces a genuinely validated readable Spanish Markdown summary with offline replay and deliberate failure evidence.

## Scope

- Run the real selected source through bounded inspection/registration/extraction, canonical generation, independent claim/coverage jobs, Markdown and local publish-check.
- Retain the actual original PDF/evidence, schema/configuration/prompts, bounded provider replies/reports, canonical draft, run record and rendered Markdown. A fake full path is regression evidence only.
- Exclude the second source, HTML assets/browser, comparison, PKI and automatic network admission from this first acceptance milestone.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] The real approved source completes the public documented one-programme path with source-page citations and both independent validators passing on identical current inputs.
- [ ] Offline replay verifies committed inputs/reports/output without calling a model or official website, and unchanged deterministic output matches.
- [ ] Missing citation/page/coverage, changed input, unsupported claim, disagreement and failed/cancelled attempts cannot become eligible output.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-first-programme`. Pass when automated cases prove stored first-source replay, deterministic Markdown, each required stage failure, stale input and immutable-source mutation; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails. Also run `make evaluate-first-programme TARGET=fixtures/frente/source.json` and retain actual source/model/tool/input identities and original outputs; unavailable required resources keep the task open, never a skipped pass.

### Gate 2: Focused behavior

Run `make test-first-programme`. Pass with exit zero and named positive/negative cases for stored first-source replay, deterministic Markdown, each required stage failure, stale input and immutable-source mutation. Do not substitute a fake for required actual-tool/model evidence.

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

- docs/SPEC.md#end-to-end-tests
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:05:04+01:00 — Created directly in todo under Juan's explicit TQ-0061 request to revise existing tasks and create the missing MVP outcomes. Outcome, scope, criteria, seven gates, priority, dependencies and execution role are declared. No authority to promote unrelated inbox proposals is granted.
