---
id: TQ-0066
title: Make the handoff explicit about autonomous MVP execution
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - area/project
  - agent/developer
created: 2026-11-02T18:29:06+01:00
updated: 2026-11-02T18:31:46+01:00
---

## Outcome

The next implementation agent resumes the approved MVP queue immediately, without questions or routine approval prompts, and continues until no eligible approved work remains.

## Scope

- Update HANDOFF.md with an explicit implementation mission reflecting Juan repeated authorization to continue the approved MVP tasks.
- Specify the read/claim/implement/gate/record/commit/push/confirm/requery loop and automatic continuation through newly unblocked work.
- Explain task-local blocker recording and independent-task continuation without claiming blocked work passed; preserve delivery-failure, execution-role, inbox and presentation boundaries.
- Confirm the current task frontier rather than invent progress. This session updates the handoff and its task record only.

## Acceptance criteria

- [ ] The handoff says to start implementation rather than ask questions, offer a plan only, request permission to continue, or stop after one task.
- [ ] It requires complete gates and confirmed task-scoped delivery before normal continuation, then re-queries the live queue.
- [ ] Task-local blockers remain open with evidence while eligible independent work continues; no-work reports distinguish exhausted work from blocked work.
- [ ] Owner-only inbox, role restrictions, supplied presentation assets and no bypass of failed gates/signing remain intact.
- [ ] The documented first tasks match the actual ready queue and current reduced MVP.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect HANDOFF.md against the five criteria. Pass when each is expressed as an imperative instruction with no routine approval question or obsolete planning-only mission for the resuming implementation agent.

### Gate 2: Focused behavior

Run deterministic handoff text/link checks and tq ready queries for developer and architect roles. Pass when the mission, delivery loop, blocker/no-work rules and current frontier match the queue and AGENTS.md.

### Gate 3: Regression

Run make verify. Pass with exit zero and no Taskqueue warnings.

### Gate 4: Static and structural validation

Run git diff --check, make check (included by make verify), handoff relative-link/anchor checks and tq show for this task. Pass with no structural findings.

### Gate 5: Documentation and specification consistency

Compare HANDOFF.md with AGENTS.md, docs/QUALITY_GATES.md, docs/SPEC.md and docs/design/MVP.md. Pass when continuation is autonomous within approved scope/roles and no deferred platform scope is revived.

### Gate 6: Trust-boundary and security behavior

Inspect stop/blocker instructions. Pass when no question is used as routine control flow and missing evidence, unavailable required tools, inbox work, presentation restrictions or delivery failures cannot be bypassed or falsely marked done. No runtime boundary is changed.

### Gate 7: Artifact and task record

Record all gate results with tq note; stage only HANDOFF.md and this task record. Pass when no runtime, existing task outcome/dependency or unrelated user work changes.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended handoff/task record, push and confirm upstream. This assignment prepares the handoff only; the resuming implementation agent follows its full approved-MVP execution mission.

## Specification

- docs/SPEC.md#solution
- docs/SPEC.md#make-targets
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:31:46+01:00 — Quality gates:
  - Gate 1 PASS — HANDOFF.md begins with an explicit authorized implementation mission: start work rather than recap/plan; no questions or permission prompts for selection, checks, commit/push or continuation; exhaust all eligible approved MVP work for the execution role.
  - Gate 2 PASS — deterministic staged-text assertions checked the mission, five-step execution/delivery/requery loop, independent-blocker continuation, global delivery failure and accurate stop/report rules. Live ready queries match developer TQ-0028/TQ-0029/TQ-0036/TQ-0063 and architect TQ-0062; no queue stderr warnings. All seven handoff local links, including the AGENTS delivery anchor, resolve.
  - Gate 3 PASS — make verify exited zero; existing package tests, Go formatting/build/vet, documentation and queue checks passed.
  - Gate 4 PASS — make check (included by make verify), git diff --check, git diff --cached --check, handoff link/anchor assertions and tq show TQ-0066 passed with no structural findings.
  - Gate 5 PASS — autonomous execution follows current SPEC/MVP and AGENTS closure rules. Historical planning-only limits are distinguished from Juan explicit resuming implementation assignment. Approved roles, dependency readiness and owner-only inbox/presentation restrictions remain unchanged.
  - Gate 6 PASS — task-local failures stay open with actual evidence and unblock conditions while independent eligible work continues. No false gate pass, invented source/model run, inbox promotion, presentation authorship, signing/hook bypass or next-task claim after delivery failure is authorized. A genuine terminal blocker produces a factual report rather than a routine question. No runtime trust boundary changed.
  - Gate 7 PASS — only HANDOFF.md and this task record are included; existing task outcomes/dependencies, runtime/source/presentation/workflow files and unrelated user work remain unchanged. No new evaluation artifact or source acquisition is claimed.
