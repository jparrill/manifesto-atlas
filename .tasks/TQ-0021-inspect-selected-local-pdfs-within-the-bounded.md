---
id: TQ-0021
title: Inspect selected local PDFs within the bounded macOS deployment
status: todo
priority: high
labels:
  - type/feature
  - agent/developer
  - area/provenance
  - area/cli
depends_on:
  - TQ-0018
  - TQ-0062
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-05T00:20:21+01:00
---

## Outcome

make submit produces a bounded inspected local candidate without claiming official-origin certification or publication approval.

## Scope

- Implement rooted/no-follow local input, bounded immutable candidate copy/hash, MIME/magic and Poppler metadata/page/encryption inspection using the one deployment approved in TQ-0062.
- Retain bounded failure/resource/cleanup evidence; use installed macOS tool identities rather than requiring the former rootless/OCI platform.
- No arbitrary-URL fetching, signing or catalogue admission belongs to this slice.
- Follow docs/design/MVP.md, which supersedes conflicting historical designs. Do not reintroduce deferred PKI/DAG/network admission, calibrated scores or comparison prerequisites.
- Regression checks are offline; explicit real converter/model evaluations are separate owning-task gates, not calls hidden inside make verify.

## Acceptance criteria

- [ ] A valid synthetic PDF records exact bytes/hash/page inspection and the actual runtime/tool configuration.
- [ ] Spoofed MIME, header-only/corrupt/encrypted PDF, excessive bytes/pages/resources and unsafe input paths reject safely.
- [ ] Real cancellation/timeout cleans the complete job and the worker cannot write the checkout or access home/secrets/general network.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-submission`. Pass when automated cases prove valid PDF, corrupt/encrypted/spoofed input, byte/page limits, actual timeout/memory enforcement, write/network denial, cancellation and teardown; every acceptance criterion must map to a named test or retained artifact. Unexpected acceptance/rejection fails.

### Gate 2: Focused behavior

Run `make test-submission`. Pass with exit zero and named positive/negative cases for valid PDF, corrupt/encrypted/spoofed input, byte/page limits, actual timeout/memory enforcement, write/network denial, cancellation and teardown. Do not substitute a fake for required actual-tool/model evidence.

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

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T15:39:02+01:00 — BLOCKER (2026-11-02): cannot pass the required sandbox gate in this environment. The approved isolation contract (submission-isolation.md) mandates a rootless Podman sandbox with a pinned OCI image (SHA-256, --pull=never), crun, seccomp, SELinux enforcing, delegated cgroup v2 limits, AND a pinned qpdf adapter for sandbox page-tree inspection, plus a disposable qualification job that proves REAL OS enforcement (memory.max, pids.max, write/network denial, timeout/fork/memory probes). Host check: kernel 7.1.13 (>=6.1 OK), podman present, crun present, SELinux Enforcing, user namespaces enabled, cgroup v2 present — but qpdf is MISSING and no pinned OCI image is present locally. The design states: 'Unavailability of required sandbox enforcement is a failed required gate, not a skipped pass.' Implementing the Go Inspect/RunIsolated code is feasible, but make test-submission cannot honestly pass without qpdf + the pinned image + a conforming host to run the qualification probes. Left in todo, not claimed. Unblocked by: installing/building the pinned qpdf adapter and OCI image, then running the sandbox qualification on a conforming host.
- 2026-11-02T17:27:01+01:00 — Correction under TQ-0058: qpdf is now installed from Fedora (qpdf 12.3.2 verified by qpdf --version), so missing host qpdf is not a current blocker. A missing host binary never proved a container image could not be built. Juan requested avoiding unnecessary rootless setup; review the isolation design against that direction and retain demonstrable resource, process cleanup and repository/network protections. Do not equate installing qpdf with satisfying those gates or execute untrusted PDFs using an undocumented fallback. Earlier notes are preserved as history, not current verified blockers.
- 2026-11-02T18:05:04+01:00 — Rescoped under Juan's explicit TQ-0061 MVP replanning request after TQ-0060. The new outcome, boundaries, criteria, seven gates, priority and native dependencies supersede the former platform contract. Previous approvals/failures/notes remain history, not claims that the new outcome already passes. Execution role is unchanged; docs/SPEC.md and docs/design/MVP.md now govern.
- 2026-11-05T00:20:21+01:00 — Owner explicitly chose macOS for the immediate E2E. TQ-0062 now defines the native Poppler/sandbox-exec runtime and TQ-0067 demonstrates it. The intake task remains todo because the production make submit/CLI interface still needs integration; reuse the demonstrated worker rather than restoring a Fedora prerequisite.
