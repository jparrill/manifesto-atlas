---
id: TQ-0005
title: Decompose the complete specification into implementation tasks
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-01T16:57:45+01:00
updated: 2026-11-01T17:22:49+01:00
---

## Outcome

A complete, dependency-aware Taskqueue backlog covers the accepted specification and can be reviewed and executed without reconstructing this session.

## Scope

- Inspect the complete specification and current repository empirically.
- Create explicitly requested gitignore and Makefile tasks in todo, with complete execution contracts.
- Create specification-derived proposals in inbox, including CI design only.
- Record requirement coverage, unresolved decisions, execution restrictions, and real blocking dependencies.

## Acceptance criteria

- [ ] Gitignore and Makefile each have a separate task with complete acceptance criteria and seven gates.
- [ ] Every specification story and implementation/testing section maps to one or more tasks.
- [ ] All new tasks use declared labels and valid, acyclic native dependencies.
- [ ] Specification proposals remain in inbox; no CI implementation task exists.
- [ ] Task bodies are English and distinguish architecture from implementation.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect the created task records and the coverage matrix in this task notes. Pass when every requested outcome and all 40 stories are covered.

### Gate 2: Focused behavior

Run a deterministic assertion-based audit of tq list --json. Pass when every new task has all required sections, seven gates, a specification reference, one type, at least one area, and exactly one execution role.

### Gate 3: Regression

Run tq list --json, tq show for every created task, git diff --check, and make verify if a Makefile exists. Pass when commands exit zero, queue stderr contains no malformed-file warnings, and existing records remain intact. No Makefile exists at planning start.

### Gate 4: Static and structural validation

Assert all dependencies exist, graph has no cycle, and configured labels/statuses/priorities are used. Pass with zero findings.

### Gate 5: Documentation and specification consistency

Compare every SPEC heading and story with the recorded coverage matrix. Pass when nothing is omitted, CI is design-only, and presentation implementation respects the LLM prohibition.

### Gate 6: Trust-boundary and security behavior

Inspect tasks for PDF isolation, provenance anchors, report signing, independent validators, hostile-input tests, and fail-closed publication. Pass when each boundary has an executable task and negative-path acceptance criteria; task records contain no secrets or personal absolute paths.

### Gate 7: Artifact and task record

Record audit counts, coverage matrix, and all gate results with tq note before completion. Pass when Taskqueue records are current and the planning task is closed only after successful checks.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0005`:

1. Record all gate results with `tq note`, then run `tq done TQ-0005`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0005` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- docs/SPEC.md#user-stories
- docs/SPEC.md#implementation-decisions
- docs/SPEC.md#testing-decisions
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-01T16:57:59+01:00 — Repository inspection: git tracks documentation and Taskqueue configuration only; no Go module, Makefile, PDFs, schemas, extraction artifacts, project skills, or CI exist. HANDOFF.md and TQ-0004 are untracked pre-existing work. SPEC references two golden PDFs absent from this checkout. Unresolved architectural contracts include trust anchors, sandbox limits, signatures/hash invalidation, taxonomy applicability, validator isolation, and site quality thresholds. Presentation assets require human/non-LLM execution under the explicit project prohibition. CI will receive an inbox design task only. The Makefile task will provide real documentation/queue checks now and explicit non-zero failures for unavailable stages, with Go checks added when Go exists.
- 2026-11-01T17:08:47+01:00 — Specification coverage — all 40 user stories:
  - Story 1: TQ-0021
  - Story 2: TQ-0052
  - Story 3: TQ-0043, TQ-0044
  - Story 4: TQ-0026, TQ-0053
  - Story 5: TQ-0022, TQ-0026
  - Story 6: TQ-0009, TQ-0021
  - Story 7: TQ-0027
  - Story 8: TQ-0027
  - Story 9: TQ-0027
  - Story 10: TQ-0030
  - Story 11: TQ-0030
  - Story 12: TQ-0030
  - Story 13: TQ-0031
  - Story 14: TQ-0032, TQ-0039
  - Story 15: TQ-0034, TQ-0035
  - Story 16: TQ-0035
  - Story 17: TQ-0036, TQ-0038
  - Story 18: TQ-0036, TQ-0038
  - Story 19: TQ-0043, TQ-0044
  - Story 20: TQ-0013, TQ-0036
  - Story 21: TQ-0037, TQ-0038
  - Story 22: TQ-0039
  - Story 23: TQ-0040
  - Story 24: TQ-0041
  - Story 25: TQ-0019, TQ-0020, TQ-0024
  - Story 26: TQ-0038, TQ-0053
  - Story 27: TQ-0018, TQ-0014, TQ-0052
  - Story 28: TQ-0043, TQ-0044, TQ-0046
  - Story 29: TQ-0042, TQ-0044
  - Story 30: TQ-0047
  - Story 31: TQ-0042, TQ-0047
  - Story 32: TQ-0046
  - Story 33: TQ-0047, TQ-0048
  - Story 34: TQ-0024
  - Story 35: TQ-0007, TQ-0050, TQ-0017
  - Story 36: TQ-0053
  - Story 37: TQ-0019, TQ-0041, TQ-0035
  - Story 38: TQ-0050
  - Story 39: TQ-0011, TQ-0037, TQ-0045
  - Story 40: TQ-0014, TQ-0053
- 2026-11-01T17:08:47+01:00 — Specification coverage — implementation and testing sections:
  - Module seams: TQ-0007, TQ-0018, TQ-0051
  - Make targets: TQ-0007, TQ-0051
  - Provenance trust model: TQ-0010, TQ-0023, TQ-0025, TQ-0026
  - Untrusted document handling: TQ-0009, TQ-0021, TQ-0022, TQ-0024
  - Immutable source storage: TQ-0011, TQ-0027
  - Extraction artifacts: TQ-0012, TQ-0030, TQ-0031
  - Extraction fidelity: TQ-0032, TQ-0033, TQ-0034, TQ-0035
  - Canonical programme model: TQ-0014, TQ-0037, TQ-0038
  - Topic taxonomy: TQ-0013, TQ-0036, TQ-0045
  - Agent skills: TQ-0025, TQ-0033, TQ-0038, TQ-0039, TQ-0040, TQ-0048
  - LLM validation: TQ-0015, TQ-0024, TQ-0039, TQ-0040, TQ-0041
  - Deterministic rendering: TQ-0016, TQ-0042, TQ-0043, TQ-0044
  - Site validation: TQ-0046, TQ-0047, TQ-0048, TQ-0049
  - Publication and CI: TQ-0008, TQ-0020, TQ-0050, TQ-0017
  - Language and licensing: TQ-0052, TQ-0053
  - Assessment separation: TQ-0014, TQ-0053
  - Submission and provenance tests: TQ-0021, TQ-0022, TQ-0023, TQ-0026
  - Source registry tests: TQ-0027, TQ-0028, TQ-0029
  - Extraction tests: TQ-0030, TQ-0031
  - Extraction audit tests: TQ-0032, TQ-0034, TQ-0035
  - Canonical model tests: TQ-0036, TQ-0037
  - Validator contract tests: TQ-0039, TQ-0040, TQ-0041
  - Rendering tests: TQ-0043, TQ-0044, TQ-0045
  - Site quality tests: TQ-0046, TQ-0047, TQ-0048, TQ-0049
  - End-to-end tests: TQ-0054, TQ-0050, TQ-0017
  - Further Notes: TQ-0010, TQ-0028, TQ-0029, TQ-0034, TQ-0053
- 2026-11-01T17:11:41+01:00 — Coverage boundaries: Problem Statement and Solution are covered by the complete stage chain TQ-0021 through TQ-0051 and contribution flow TQ-0052; User Stories are mapped individually above; Implementation Decisions and Testing Decisions are mapped by subsection above. Out of Scope is enforced by task boundaries: no upload service/accounts/database/embeddings, no unverified publication, no human-review fallback, no LLM presentation code, no composite political scoring, and no runtime generation. CI deployment behavior is deliberately covered by design-only TQ-0017 pending owner approval, not an implementation promise. Future binary-storage migration remains a documented risk, not current product scope. The 12 Spec Make stages have owners: submit TQ-0021; verify-source TQ-0026; add TQ-0027; extract TQ-0030; validate-extraction TQ-0035; normalize TQ-0038; validate-content TQ-0041; build TQ-0043/TQ-0044; validate-site TQ-0049; publish-check TQ-0050; process/serve TQ-0051. Dependency audit: 140 native edges, no missing references or cycles. Execution restriction: TQ-0042 requires human/non-LLM presentation authorship. User subsequently requested implementation of TQ-0006 and an updated HANDOFF.md; these are handled as separately verified work.
- 2026-11-01T17:11:41+01:00 — Quality gates:
  - Gate 1 PASS — 49 backlog tasks created: TQ-0006/TQ-0007 in todo and 47 specification proposals in inbox. All 40 stories, 25 implementation/testing subsections, Further Notes, and all root-section scope constraints are mapped above. CI is design-only TQ-0017.
  - Gate 2 PASS — deterministic assertion audit reread all 49 tasks with tq show --json and checked required sections, seven gates, acceptance criteria, specification references, declared labels, priorities, statuses, and exact dependencies. All assertions passed.
  - Gate 3 PASS — tq list --json and every tq show exited zero with empty stderr; git diff --check passed. No Makefile exists, so the task-defined pre-Makefile regression applies.
  - Gate 4 PASS — all 140 dependency references exist and graph is acyclic. Whitespace checks covered all 50 newly created task files, including untracked files, using git diff --no-index --check; zero findings. Taskqueue structural round trips passed. The audit accounts for CLI trailing-newline normalization and git no-index difference exit status.
  - Gate 5 PASS — full SPEC inspection and coverage matrices include every story, module/test subsection, all six planned skills, all 12 Make stages, public language/licensing, regional comparison, and assessment separation. No CI implementation task exists.
  - Gate 6 PASS — explicit design/negative-case tasks cover quarantine/isolation, safe URL transport, party trust anchors, immutable storage, signing authority, independent semantic jobs, stale reports, escaping, and fail-closed publication. Task content audit contains no private absolute paths or secrets.
  - Gate 7 PASS — all task files are current CLI-generated artifacts; repeated creation/audit reused existing IDs without duplicates; coverage and evidence recorded before completion.
- 2026-11-01T17:22:49+01:00 — Delivery preparation under the newly requested automatic policy: final phases were added through tq update, with all task metadata, native dependencies, acceptance criteria, historical notes, and the original coverage matrix preserved. Regression reread all 56 current records with empty stderr, confirmed the graph remains acyclic, and git diff --check passed. TQ-0006 implementation and record are delivered separately; this planning commit contains TQ-0005 and the remaining created backlog proposals/Makefile task only.
