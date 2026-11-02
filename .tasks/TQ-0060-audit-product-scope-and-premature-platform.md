---
id: TQ-0060
title: Audit product scope and premature platform complexity
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/chore
  - area/project
  - agent/architect
created: 2026-11-02T17:44:22+01:00
updated: 2026-11-02T17:51:05+01:00
---

## Outcome

Juan can decide a minimal product direction using a measured comparison of the original repository intent, current implementation, expanded design contracts and backlog.

## Scope

- Recover the earliest available product baseline from Git; distinguish explicit requirements from later architecture choices.
- Inventory code, tests, dependencies, actual production callers, public pipeline behavior, task graph and missing vertical product outcomes.
- Rank concrete simplification candidates, explain the IPv6/SSRF boundary, and recommend a minimal end-to-end scope with preserved evidence/citation/independent-validation guarantees.
- Quantify removable code with non-overlapping file counts and list prerequisites/risks of each recommendation.
- Audit/report only. Do not delete runtime code, weaken active trust boundaries, rewrite the accepted specification, promote inbox work or resume implementation under the old expanded backlog.

## Acceptance criteria

- [ ] The report links the earliest available baseline and current contract/code evidence and identifies where requirements expanded.
- [ ] Every deletion recommendation has a repository-wide caller/reference check and a quantified non-overlapping estimate.
- [ ] The report explains why IPv6 filtering exists and offers a smaller trust boundary rather than removing protection from an arbitrary-URL fetcher.
- [ ] A concrete MVP path prioritizes one real programme end to end and separates retain, defer and remove recommendations.
- [ ] The final Spanish report is actionable and current implementations/requirements remain unchanged pending owner scope decisions.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect docs/audits/2026-11-02-product-scope.md. Pass when all criteria map to history, paths, command output or measured counts.

### Gate 2: Focused behavior

Run documented git history/inventory/caller checks, task dependency audit and CLI pipeline probes. Pass when observations match the report, counts do not double-count and no pipeline stub is described as working.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero; passing tests are not treated as proof of end-to-end product readiness.

### Gate 4: Static and structural validation

Run git diff --check and make check; validate audit document links and tq show. Pass with no structural errors or queue warnings.

### Gate 5: Documentation and specification consistency

Compare earliest available README/SPEC with current README/SPEC, designs and task bodies. Pass when original commitments and later choices are distinguished and proposed scope changes are not presented as approved changes.

### Gate 6: Trust-boundary and security behavior

Use read-only inspection, help and unavailable-stage probes; no live network/model, credential access, untrusted PDF processing or runtime edits. Pass when proposed boundary reduction is distinguished from weakening the existing boundary.

### Gate 7: Artifact and task record

Record measured evidence and all seven gate results using tq note. Include only audit documentation and task record. Pass when no runtime/source/approval artifact or unrelated task is modified.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended audit artifacts/task record, push and confirm delivery. This assignment is audit-only; do not continue implementation until Juan chooses the product scope.

## Specification

- docs/SPEC.md#solution
- docs/SPEC.md#module-seams
- docs/SPEC.md#out-of-scope
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T17:51:05+01:00 — Quality gates:
  - Gate 1 PASS — docs/audits/2026-11-02-product-scope.md compares the earliest committed intent with current implementation/design/backlog, explains IPv6 as a consequence of arbitrary-URL admission, ranks nine simplifications, quantifies disjoint cuts, and proposes a one-document vertical release with preserved core guarantees.
  - Gate 2 PASS — git diff --exit-code e10e85a 94c4cad -- README.md docs/SPEC.md exited zero. Baseline git-show inventory assertions confirmed 19 Go files/7048 lines (3938 application, 2897 tests, 213 audit), 2109 disjoint gross package-removal candidates plus 460 separate artifact production lines. go list proved cmd imports only cli, no domain packages; repository-wide reference searches included tests, docs, contract strings and tasks. All 12 CLI stages and all 12 Make stages were invoked and correctly reported unavailable/nonzero. Native task graph audit found 60 tasks/143 edges, no cycle/missing dependency, longest path 29 including completed work; reported ordering problems match actual dependencies.
  - Gate 3 PASS — make verify and go test ./... exited zero. Report explicitly distinguishes passing library tests from missing end-to-end product behavior.
  - Gate 4 PASS — make check (via make verify), git diff --check and git diff --cached --check passed. Structural link check verified all 19 local audit links. tq show TQ-0060 is valid; no queue warnings.
  - Gate 5 PASS — earliest README/SPEC already contain the expanded commitments; the report does not invent an earlier committed minimal spec or blame all complexity on later agents. Proposed changes to signed reports/independent downloading/site gates remain recommendations, not authorized contract edits. Technical artifact is English; final user report is Spanish.
  - Gate 6 PASS — inspection, local help/stub probes and existing offline tests only; no network/model call, credential access, PDF parsing or runtime edit. Report preserves converter isolation, strict JSON, evidence, independent validation and accessibility basics; reducing the admission boundary is distinguished from deleting IPv6 protection on an exposed fetcher.
  - Gate 7 PASS — only this audit report and TQ-0060 record change. Original implementation, specification, other task statuses/dependencies, source/approval artifacts and inbox proposals remain unchanged. Gross line counts and zero installed external Go dependencies are not presented as measured final net savings.
