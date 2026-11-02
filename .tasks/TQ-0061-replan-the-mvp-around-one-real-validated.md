---
id: TQ-0061
title: Replan the MVP around one real validated programme
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-02T17:57:35+01:00
updated: 2026-11-02T18:20:14+01:00
---

## Outcome

The specification and executable Taskqueue backlog prioritize one maintainer-selected PDF through cited Spanish Markdown and independent validation, with a separate static-site extension and no premature platform prerequisites.

## Scope

- Adopt the owner-approved product-scope audit as the MVP direction and define small local source/run/evidence interfaces.
- Rewrite accepted unfinished tasks and native dependencies; defer non-MVP work to inbox without promoting, rejecting, executing or otherwise triaging existing inbox proposals.
- Create only the missing approved runtime-profile, platform-retirement/run-record, first-programme E2E and local-preview slices with complete acceptance criteria and seven gates.
- Preserve completed task history and append supersession notes where needed; preserve previous gate-failure evidence.
- Update SPEC, current design precedence, README/Make reference and HANDOFF so implementation agents cannot accidentally execute obsolete contracts.
- Planning only: no runtime implementation, source acquisition, model execution, deployment or presentation authorship.

## Acceptance criteria

- [ ] One real source can be recovered before generic runners or registration, and the first extraction/Markdown path does not depend on both golden PDFs, report PKI, arbitrary-URL transport, fidelity scores, regional comparisons or visual-model QA.
- [ ] All executable tasks have a complete approved role/scope/acceptance/gate contract and valid acyclic dependencies with no active dependency on deferred work, except the explicit supplied-presentation boundary.
- [ ] New tasks cover removal of unused platform code, simple input-bound run records, a bounded Fedora runtime choice, first real end-to-end proof and draft preview independent of publication.
- [ ] Content citations, immutable sources, full-source coverage, independent claim checks, safe bounded conversion and fail-closed publication remain required.
- [ ] Current documentation clearly supersedes old architecture for the MVP and distinguishes planned behavior from existing unavailable commands.

## Quality Gates

### Gate 1: Acceptance criteria

Audit tq list --json against the MVP milestone/dependency map and docs/SPEC.md. Pass when every criterion is represented and no prohibited platform dependency lies on the first-programme path.

### Gate 2: Focused behavior

Run a deterministic queue audit of roles, labels, sections, seven numbered gates, specification anchors, native dependency existence/cycles, preservation of existing inbox proposals and completed history. Pass with zero findings.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero and no Taskqueue warnings.

### Gate 4: Static and structural validation

Run git diff --check, make check, documentation link/anchor checks and tq ready queries. Pass with zero structural issues and only intended new-scope tasks ready.

### Gate 5: Documentation and specification consistency

Compare SPEC, MVP contract, historical design banners, README, HANDOFF, Make reference and every changed task. Pass when current scope, deferred features, source authority and presentation restrictions agree.

### Gate 6: Trust-boundary and security behavior

Verify each changed outcome retains its relevant negative-input/failure tests and that source selection never bypasses content validation. Planning must not invoke models, fetch PDFs, change runtime code or expose secrets.

### Gate 7: Artifact and task record

Record task counts, scope decisions, dependency audit and all seven gates with tq note. Stage only intended planning/documentation/task records; preserve history and unrelated work.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended changes, push and confirm upstream. This is a planning-only assignment; do not auto-expand it into implementation after delivery.

## Specification

- docs/SPEC.md#solution
- docs/SPEC.md#module-seams
- docs/SPEC.md#out-of-scope
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:18:10+01:00 — Scope decision: Juan explicitly requested modifying existing tasks and creating missing tasks to follow the TQ-0060 simplification line. Adopted a trusted local source boundary, one configured provider/three independent roles, small local schemas/run records, first real Markdown acceptance, second-document repetition and a separate supplied static-site/preview extension. This authorizes planning changes, not claims that the new runtime exists. Eight accepted tasks returned to inbox, two separate platform tasks were superseded/rejected with their critical checks reassigned, and four missing outcomes were created directly in todo under this explicit assignment. No pre-existing inbox proposal was triaged.
- 2026-11-02T18:18:10+01:00 — Quality gates:
  - Gate 1 PASS — SPEC and docs/design/MVP.md map the current local-first interfaces and milestones to 23 rewritten accepted tasks plus new TQ-0062/TQ-0063/TQ-0064/TQ-0065. First source recovery depends only on this plan; extraction/first Markdown/E2E have no second-source, PKI/DAG, arbitrary-URL admission, calibrated-score, comparison or HTML dependency. Draft preview has no content/publication prerequisite. Core citations, source immutability, original evidence, both independent passes, bounds/isolation and failure blocking remain required.
  - Gate 2 PASS — deterministic CLI/index audit (python3 /tmp/opencode/validate-mvp-plan.py using tq list --json and git baseline c77fa9c) verified 65 tasks, 124 edges, valid declared labels/one type/one role, required sections/seven gates, all dependencies existent and acyclic, all historical notes retained and 25 pre-existing completed records preserving metadata/body history. Existing inbox TQ-0017/TQ-0042 are byte-identical to baseline. Eight deferred tasks are inbox and TQ-0032/TQ-0035 are rejected, not falsely completed. No active dependency enters deferred work except the explicitly supplied TQ-0042 asset boundary.
  - Gate 3 PASS — make verify and go test ./... exited zero after the final documentation changes, with no queue warnings or live model/source-site calls.
  - Gate 4 PASS — make check (via make verify), git diff --check and git diff --cached --check passed. The planning audit resolved 149 local Markdown links/anchors and all active task SPEC anchors. Ready queries were empty while this planning prerequisite was in progress; the planned post-delivery frontier is TQ-0028/TQ-0036/TQ-0063 plus normal-priority TQ-0029 and architect TQ-0062.
  - Gate 5 PASS — current SPEC/MVP, nine historical-contract banners, README, HANDOFF, Make reference and changed task outcomes/gates agree. Existing unavailable code/Make recipes remain described as legacy until TQ-0063 migrates ownership; no function is claimed implemented by a plan. Public text is Spanish, technical/task text English. Presentation and CI owner-only inbox restrictions remain intact.
  - Gate 6 PASS — every revised task declares its actual interface negatives: source/evidence hash and authority mismatch, unsafe paths/partial writes, real resource/cleanup denial, strict/cited canonical data, clean bounded model requests, full-source checked inventories, stale/swapped/contributor approvals, output escaping and preview exclusion. Planning performed no model call, PDF fetch/parse, runtime edit or credential access. Model evaluations must genuinely run in their owning tasks and remain separate from offline regression.
  - Gate 7 PASS — only intended task records and planning/documentation are staged. CLI-managed titles/filenames/frontmatter/lifecycle and preserved notes were used throughout. No original source, generated approval, runtime/presentation/workflow or unrelated user file was changed. Current legacy code removal is assigned to TQ-0063, not falsely reported as already done.
- 2026-11-02T18:19:18+01:00 — Final manual gate review tightened TQ-0062: qpdf --version alone cannot establish runtime controls. Its acceptance gate now requires an observed five-scenario contract/log, and its focused gate must actually execute the recorded positive/timeout/memory/write/network/teardown commands. No configuration-only or simulated enforcement can qualify the chosen deployment.
- 2026-11-02T18:20:14+01:00 — Final gate rerun after tightening the runtime task and clarifying build ownership: Gates 1/2/4/5/6/7 PASS via the complete 65-task/124-edge CLI/index audit, preserved history/inbox checks, all 149 links/anchors, current scope/negative-case contracts, and zero runtime/workflow edits; Gate 3 PASS via make verify and go test ./..., both exit zero. The actual runtime remains unqualified until TQ-0062 executes its required probe commands.
