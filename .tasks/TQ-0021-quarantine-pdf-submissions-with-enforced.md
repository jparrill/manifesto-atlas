---
id: TQ-0021
title: Quarantine PDF submissions with enforced resource and access limits
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/provenance
  - area/cli
depends_on:
  - TQ-0018
  - TQ-0019
  - TQ-0009
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T17:27:01+01:00
---

## Outcome

make submit accepts a valid candidate into quarantine and safely rejects hostile or over-budget inputs.

## Scope

- Implement submission metadata, byte checksum, MIME/magic inspection, page inspection, quarantine lifecycle, and the approved sandbox runtime.
- Enforce byte/page/time/memory/storage limits, repository-write and network restrictions, subprocess cleanup, and failure records.
- A submission is not a registered official source and cannot appear in the public catalogue.
- Provide the focused Make gate `make test-submission` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] Valid PDFs produce a quarantined record with exact submitted checksum and bounded inspection metadata.
- [ ] Invalid MIME, malformed/encrypted PDFs under the chosen policy, excessive resources, unsupported isolation, and unsafe paths fail with retained bounded evidence.
- [ ] Cancellation cleans up the complete process tree; hostile input cannot write committed repository paths or access network/secrets.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-submission`. Pass when automated cases prove: valid PDF; MIME spoof; corrupt/encrypted PDF; byte/page limits; timeout/memory cap; cancellation; sandbox denial. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-submission`. Pass when automated cases prove: valid PDF; MIME spoof; corrupt/encrypted PDF; byte/page limits; timeout/memory cap; cancellation; sandbox denial. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0021 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Use adversarial fixtures and sandbox probes to prove write/network denial and enforce real time/memory limits, not only configuration flags. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0021`:

1. Record all gate results with `tq note`, then run `tq done TQ-0021`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0021` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#untrusted-document-handling`
- `docs/SPEC.md#submission-and-provenance-tests`
- `docs/SPEC.md#make-targets`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T15:39:02+01:00 — BLOCKER (2026-11-02): cannot pass the required sandbox gate in this environment. The approved isolation contract (submission-isolation.md) mandates a rootless Podman sandbox with a pinned OCI image (SHA-256, --pull=never), crun, seccomp, SELinux enforcing, delegated cgroup v2 limits, AND a pinned qpdf adapter for sandbox page-tree inspection, plus a disposable qualification job that proves REAL OS enforcement (memory.max, pids.max, write/network denial, timeout/fork/memory probes). Host check: kernel 7.1.13 (>=6.1 OK), podman present, crun present, SELinux Enforcing, user namespaces enabled, cgroup v2 present — but qpdf is MISSING and no pinned OCI image is present locally. The design states: 'Unavailability of required sandbox enforcement is a failed required gate, not a skipped pass.' Implementing the Go Inspect/RunIsolated code is feasible, but make test-submission cannot honestly pass without qpdf + the pinned image + a conforming host to run the qualification probes. Left in todo, not claimed. Unblocked by: installing/building the pinned qpdf adapter and OCI image, then running the sandbox qualification on a conforming host.
- 2026-11-02T17:27:01+01:00 — Correction under TQ-0058: qpdf is now installed from Fedora (qpdf 12.3.2 verified by qpdf --version), so missing host qpdf is not a current blocker. A missing host binary never proved a container image could not be built. Juan requested avoiding unnecessary rootless setup; review the isolation design against that direction and retain demonstrable resource, process cleanup and repository/network protections. Do not equate installing qpdf with satisfying those gates or execute untrusted PDFs using an undocumented fallback. Earlier notes are preserved as history, not current verified blockers.
