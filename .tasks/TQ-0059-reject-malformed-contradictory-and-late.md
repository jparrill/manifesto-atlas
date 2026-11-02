---
id: TQ-0059
title: Reject malformed contradictory and late semantic replies
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/security
  - area/validation
  - agent/developer
depends_on:
  - TQ-0058
created: 2026-11-02T17:29:50+01:00
updated: 2026-11-02T17:33:59+01:00
---

## Outcome

The semantic runner rejects malformed common response envelopes, contradictory passing findings, over-budget replies, and replies received after cancellation or deadline before returning a report.

## Scope

- Implement the response-boundary repairs authorized by Juan following TQ-0058, using the existing artifact strict parser and runner interface.
- Enforce exact common required JSON fields, known fields, nullability, strict single JSON value, response-byte limits, and pass/finding consistency.
- Apply a bounded default request deadline, validate limit values, and reject late/cancelled results and ineligible retries.
- Add regression cases at Runner.Run and document precisely the guarantees proved.
- This required repair is a prerequisite of reopened TQ-0024, not completion of its configured-provider, evidence/schema inventory, capability or session-isolation work.

## Acceptance criteria

- [ ] Strict valid common replies retain their pass/fail/insufficient/error results; extra/trailing/duplicate/case-aliased/missing/null fields reject.
- [ ] Passing reports with blocker or warning findings reject; informational findings remain eligible.
- [ ] Tightened byte limits are enforced; cancelled or expired invocations never return a successful report.
- [ ] A transient error with partial response data or cancellation does not retry; a genuine pre-response transient error retains at most one clean retry.
- [ ] TQ-0024 remains open with its other required work explicitly recorded.

## Quality Gates

### Gate 1: Acceptance criteria

Run make test-semantic-runner. Pass when each acceptance criterion maps to named positive and negative Runner.Run tests and TQ-0024 remains open.

### Gate 2: Focused behavior

Run go test ./internal/semantic/... -count=1. Pass with zero test failures. Re-run the audit probe; its four semantic rejection checks must print PASS while remaining unrelated failures stay visible.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero.

### Gate 4: Static and structural validation

Run make check and git diff --check. Pass with no formatting/build/vet/queue/whitespace findings.

### Gate 5: Documentation and specification consistency

Update docs/MAKEFILE.md and the audit repair status against docs/design/semantic-validation.md. Pass when common-envelope checks are distinguished from unimplemented stage-specific schemas and configured-provider isolation.

### Gate 6: Trust-boundary and security behavior

Run strict JSON, contradictory report, partial-response retry, pre-cancellation, late-reply and invalid-limit cases through Runner.Run. Pass only when all reject with no successful report and no extra provider call.

### Gate 7: Artifact and task record

Record all seven gates using tq note; preserve original reports/history and leave TQ-0024 incomplete. Pass when only intended repair files and records are included.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, tq done, commit intended changes and task records, push and confirm upstream before selecting more work. Unresolved required gates keep their owning tasks open.

## Specification

- docs/SPEC.md#llm-validation
- docs/SPEC.md#untrusted-document-handling
- docs/design/semantic-validation.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T17:33:59+01:00 — Quality gates:
  - Gate 1 PASS — make test-semantic-runner passed. TestStrictResponseBoundary covers trailing/duplicate/unknown/case-aliased/missing/null JSON members; TestReportDecisionConsistency covers all four results and three severities; TestTightenedResponseLimit proves exact/tightened/default byte ceilings; TestTerminalInvocationBoundary covers pre-cancel, late cancellation/deadline and ineligible retries; TestInvocationLimitValidation proves bounded defaults and invalid-limit preflight; TestRetryDoesNotReuseExistingJob protects retry identity. Existing TestRetryOnce retains the one clean transient retry. TQ-0024 remains todo.
  - Gate 2 PASS — go test ./internal/semantic/... -count=1 exited zero. Tests were added and observed failing before the repair. go run ./docs/audits/probe now prints PASS for all four semantic invariants and FAIL for the ten unrelated known findings (expected aggregate exit 1; not an overall security pass). The synthetic positive semantic baseline succeeds first.
  - Gate 3 PASS — make verify and go test ./... exited zero after final code/test changes.
  - Gate 4 PASS — make check (via make verify), gofmt, build, vet, Taskqueue integrity and git diff --check passed with no findings.
  - Gate 5 PASS — docs/MAKEFILE.md, Makefile comments, package comments and audit repair status distinguish common response checks from missing provider/evidence/stage-schema/isolation guarantees. No approved contract weakened; original history preserved.
  - Gate 6 PASS — Runner.Run rejects hostile envelope variants, contradictory passes, oversized raw replies, late results and ineligible retries without returning a successful report or making an extra provider call. No network, model, credential or PDF used. Non-cooperative adapter termination remains explicitly outside this repair and required in TQ-0024.
  - Gate 7 PASS — only intended semantic code/tests, audit reproduction/status, Make documentation/comments, TQ-0059 and parent dependency/progress record are included. All broader parent gates stay open; no model or isolation qualification is claimed.
