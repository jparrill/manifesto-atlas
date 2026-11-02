---
id: TQ-0062
title: Choose one minimal bounded Fedora PDF deployment
status: todo
priority: high
labels:
  - type/docs
  - area/extraction
  - agent/architect
depends_on:
  - TQ-0061
created: 2026-11-02T18:05:04+01:00
updated: 2026-11-02T18:19:18+01:00
---

## Outcome

One concrete Fedora deployment can run installed qpdf/conversion tools with bounded resources and denied repository/home/network access, without requiring rootless Podman or a custom image platform.

## Scope

- Choose exactly one existing OS/container mechanism supported on this host; document its fixed executable/argv, privilege/setup requirements, input/output roots and resource/teardown behavior in docs/design/fedora-runtime.md.
- Use Fedora-packaged qpdf and existing tools; record actual versions. Do not build a generic sandbox framework or require binary-sealed OCI images.
- Run small disposable native-tool probes for timeout/descendant cleanup, memory bound, write denial and network denial. This is architecture plus qualification evidence, not the intake implementation.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] The contract names one working deployment and reproducible commands, with no implementation-critical isolation decision left to a developer.
- [ ] A bounded synthetic qpdf job succeeds; actual timeout/memory/write/network negatives prove the selected controls and cleanup.
- [ ] The contract keeps source bytes read-only and secrets/checkout unavailable without the former platform-wide qualification suite.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect docs/design/fedora-runtime.md and its referenced retained probe log. Pass only when the artifact names one concrete deployment, actual versions, fixed argv/roots/limits/privileges/teardown and an observed scenario table for positive qpdf, timeout/descendant cleanup, memory enforcement, write denial and network denial. Each row must include the exact disposable command, expected outcome and actual output/exit evidence; no implementation-critical decision or unrun probe may be called qualified.

### Gate 2: Focused behavior

Run `qpdf --version` and each exact disposable native-tool probe command recorded in docs/design/fedora-runtime.md. Pass only when a synthetic qpdf job succeeds, the deadline and memory probes are actually stopped/bounded, prohibited writes and network access are denied, source/repository bytes remain unchanged and teardown leaves no descendants. Version output alone, configuration flags alone, simulated denial or an unavailable/unrun probe fails this gate.

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

- docs/SPEC.md#untrusted-document-handling
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:05:04+01:00 — Created directly in todo under Juan's explicit TQ-0061 request to revise existing tasks and create the missing MVP outcomes. Outcome, scope, criteria, seven gates, priority, dependencies and execution role are declared. No authority to promote unrelated inbox proposals is granted.
