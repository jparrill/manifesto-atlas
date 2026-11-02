---
id: TQ-0058
title: Audit recent trust-boundary closures and restore truthful readiness
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/security
  - area/project
  - area/provenance
  - area/validation
  - agent/architect
created: 2026-11-02T17:24:47+01:00
updated: 2026-11-02T17:28:11+01:00
---

## Outcome

Recent implementation completion claims match demonstrable behavior, and downstream work is blocked by any incomplete prerequisite.

## Scope

- Audit TQ-0022, TQ-0023, TQ-0024, TQ-0036, and TQ-0037 against their approved designs and existing tests.
- Reproduce concrete failures with isolated offline probes; distinguish reproduced defects from missing deliverables.
- Reopen original tasks with precise corrective evidence rather than create optional follow-ups for required gates.
- Correct the obsolete TQ-0021 qpdf and TQ-0025 golden-source blocker assumptions; preserve historical notes.
- Record a technical audit and prioritized recovery order. No runtime implementation or isolation-policy replacement is included.

## Acceptance criteria

- [ ] The audit maps each reviewed task to its contract, evidence, missing guarantees, and corrective status.
- [ ] Offline probes reproduce concrete false acceptance or failure behavior without external network or real source PDFs.
- [ ] Incomplete tasks no longer satisfy native dependencies; historical commits and gate notes remain intact.
- [ ] TQ-0025 does not claim golden PDFs are mandatory for a synthetic configured-model evaluation; qpdf availability is checked.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect docs/audits/2026-11-02-trust-boundaries.md and tq show for all reviewed tasks. Pass when every criterion maps to explicit evidence and lifecycle state.

### Gate 2: Focused behavior

Run the documented offline Go probes against current public APIs. Pass for this audit when every reported reproduction produces its stated observed behavior; this is evidence of defects, not a runtime security pass.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero; existing passing suites are reported separately from the audit findings.

### Gate 4: Static and structural validation

Run git diff --check and make check; inspect tq ready --label agent/developer --json. Pass with zero structural errors and no incomplete dependency reported as ready.

### Gate 5: Documentation and specification consistency

Compare findings with docs/design/provenance-trust.md, docs/design/semantic-validation.md, docs/design/taxonomy-jurisdiction.md, docs/design/canonical-programme.md, and docs/QUALITY_GATES.md. Pass when technical text is English and no implementation restriction is silently waived.

### Gate 6: Trust-boundary and security behavior

Run only synthetic local probes with no production credentials, remote calls, unsafe PDF execution, or committed source mutations. Pass when findings distinguish behavior proved from behavior merely asserted.

### Gate 7: Artifact and task record

Preserve prior notes and commits, record all gate evidence with tq note, and stage only the audit and intended task records. Pass when audit evidence and queue corrections are current.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work: record gates, mark done, commit intended changes and task records, push and confirm upstream. Continue only with eligible work within the current assignment; unresolved prerequisite or delivery failures block continuation.

## Specification

- docs/SPEC.md#provenance-trust-model
- docs/SPEC.md#llm-validation
- docs/SPEC.md#canonical-model-tests
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T17:28:11+01:00 — Quality gates:
  - Gate 1 PASS — docs/audits/2026-11-02-trust-boundaries.md maps all five reviewed tasks to contracts, observed counterexamples, missing deliverables and reopened todo state; qpdf 12.3.2 confirmed; TQ-0021/TQ-0025 corrective notes preserve prior history.
  - Gate 2 PASS (audit evidence only) — go run ./docs/audits/probe printed the documented fourteen FAIL invariants and exited 1 at baseline 37b7968. This reproduces defects; it does not qualify the runtime. Probe source is retained for deterministic replay.
  - Gate 3 PASS — make verify and go test ./... exited zero. Existing regression success is explicitly distinguished from unmet contract gates.
  - Gate 4 PASS — git diff --check and make check (included by make verify) exited zero. tq ready excludes TQ-0025 and TQ-0037 after original prerequisites were reopened; queue has no malformed-file warnings.
  - Gate 5 PASS — findings compared with provenance-trust, semantic-validation, taxonomy-jurisdiction, canonical-programme and isolation contracts plus QUALITY_GATES.md. Technical documentation is English; no product or sandbox guarantee silently weakened.
  - Gate 6 PASS — only in-memory synthetic records, an offline dial recorder and a bounded late-reply adapter were used; no remote connection, model invocation, PDF parsing, secrets or source writes.
  - Gate 7 PASS — all original commits, gate notes and native dependencies preserved; only audit artifacts and intentional task-record corrections are included. No optional follow-up defers an original required gate.
