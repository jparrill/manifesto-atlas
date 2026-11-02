---
id: TQ-0044
title: Render supplied HTML as explicit safe local previews
status: todo
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/rendering
  - area/cli
depends_on:
  - TQ-0042
  - TQ-0037
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-02T18:07:29+01:00
---

## Outcome

Supplied human/non-LLM assets render deterministic Spanish HTML previews from structurally valid local canonical fixtures or drafts before content/publication implementation is complete.

## Scope

- Follow docs/design/MVP.md. Implement Go data binding, contextual escaping, output plumbing and an explicit labelled private preview mode only; no HTML/CSS/JavaScript/template/fixture authorship or embedded presentation strings.
- Consume the supplied TQ-0042 assets and TQ-0037 canonical schema. No dependency on content approval, Markdown implementation, comparisons or publish-check is needed for a private preview.
- Default public HTML mode stays unavailable until TQ-0049 integrates the real current content and basic site checks. A preview must never masquerade as eligible public output.
- Regression uses supplied assets and local fixtures only, without model or source-site requests.

## Acceptance criteria

- [ ] Supplied assets preserve every required Spanish topic, Statement, measure, condition, institution and citation from a structurally valid fixture/draft.
- [ ] Unchanged canonical data/assets yield byte-identical HTML; hostile text and links are contextually escaped without factual alterations.
- [ ] Preview is explicitly labelled and written under a separate private root; public mode refuses until actual eligibility integration, and no LLM presentation-code changes occur.

## Quality Gates

### Gate 1: Acceptance criteria

Run make test-html-rendering. Pass when supplied canonical/HTML fixtures prove all required content/states, safe escaping, byte determinism, explicit draft labelling/root separation and refusal of public mode without real eligibility.

### Gate 2: Focused behavior

Run make test-html-rendering. Pass with zero failures for good fixtures, malformed/stale canonical data, script/attribute/URL injection, unsafe paths, asset changes and preview/public confusion. Use actual supplied assets, never LLM-authored presentation fixtures.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero and no queue warnings, live model calls or source-site requests.

### Gate 4: Static and structural validation

Run make check and git diff --check; validate changed schemas/fixture references. Pass with zero compiler/vet/format/schema/link findings.

### Gate 5: Documentation and specification consistency

Compare docs/SPEC.md and docs/design/MVP.md with actual preview flags/output roots. Pass when preview is not described as publication and all public labels are Spanish. Confirm no presentation-code authorship in the diff.

### Gate 6: Trust-boundary and security behavior

Run the injection/path/preview-public negative cases through the actual renderer. Pass only when untrusted data remains escaped and no preview or payload flag grants public eligibility or repository/source writes.

### Gate 7: Artifact and task record

Confirm deterministic current fixture outputs and preserved input/assets/source bytes. Record all gates with tq note before completion; unavailable required supplied assets or failed gates keep this task open.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, done, task-scoped commit, push and upstream confirmation before approved in-scope continuation. Never promote inbox assets or bypass signing/hooks.

## Specification

- docs/SPEC.md#deterministic-rendering
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T18:56:23+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-02T18:07:29+01:00 — TQ-0061 dependency audit further splits preview from public HTML: removed TQ-0043 so the preview chain has no transitive claim/coverage/publication prerequisite. This task implements supplied-asset private previews; TQ-0049 later connects actual validated public HTML after Markdown eligibility and basic site checks exist.
