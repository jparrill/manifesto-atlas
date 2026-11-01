---
id: TQ-0009
title: Define PDF quarantine and sandbox enforcement policy
status: inbox
priority: high
labels:
  - type/docs
  - agent/architect
  - area/provenance
  - area/extraction
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:29+01:00
---

## Outcome

Submission and conversion have enforceable resource and access limits before handling real PDFs.

## Scope

- Choose concrete MIME/magic, byte, page, wall-time, memory, subprocess, and temporary-storage limits, with bounded override policy.
- Specify supported sandbox runtime and host platforms, read-only source mounts, writable scratch boundaries, network denial, teardown, and evidence of limit enforcement.
- Cover encrypted/malformed PDFs, symlinks, traversal, decompression bombs, subprocess trees, cancellation, and unsupported isolation environments.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] A supported local execution profile can enforce all required restrictions; unsupported platforms fail clearly instead of silently running unsandboxed.
- [ ] Limit values, error/status mapping, cleanup behavior, and test strategy are concrete enough for a developer to implement.
- [ ] docs/design/submission-isolation.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/submission-isolation.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/submission-isolation.md. Pass only when each of these cases has one explicit, internally consistent expected result: valid native PDF; spoofed MIME; oversize; excess pages; encrypted/corrupt PDF; timeout; memory exhaustion; attempted repository write/network access. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0009 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0009`:

1. Record all gate results with `tq note`, then run `tq done TQ-0009`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0009` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#untrusted-document-handling`
- `docs/SPEC.md#submission-and-provenance-tests`
- `docs/QUALITY_GATES.md`
