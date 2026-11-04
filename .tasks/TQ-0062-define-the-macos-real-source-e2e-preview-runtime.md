---
id: TQ-0062
title: Define the macOS real-source E2E preview runtime
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - area/extraction
  - agent/architect
depends_on:
  - TQ-0061
created: 2026-11-02T18:05:04+01:00
updated: 2026-11-04T23:02:08+01:00
---

## Outcome

A concrete macOS local-preview contract lets the requested first real PDF reach a truthful HTML preview without waiting for Fedora.

## Scope

- Juan explicitly selected macOS and supplied the existing Programas-electorales-2026.html as the visual reference in this session.
- Define the local Poppler/sandbox-exec path, bounded subprocess execution, exact page/citation mapping, current canonical validator reuse, independent configured-provider requests, and immutable supplied-asset data binding in docs/design/macos-preview.md.
- This is a private preview, not public publication or a claim that all production pipeline tasks are complete. Semantic failure/insufficient evidence remains visible and never becomes public eligibility.
- Keep the supplied HTML/CSS/JavaScript byte-identical as an input asset; Go may replace only its programme data payload with escaped canonical data. Do not author presentation code.
- Actual runtime negative tests, real-source/provider execution and browser proof belong to dependent TQ-0067. This design records observed macOS capabilities without claiming unrun enforcement.

## Acceptance criteria

- [ ] The contract names available macOS tools, file/network restrictions, deadlines/process-group cleanup, output budgets and an RSS watchdog with its sampled-limit limitation.
- [ ] Explicit commands cover importing the owner asset, running the real source, rerendering stored results and opening the resulting local HTML.
- [ ] Canonical input bindings, original-page citations, independent reports and private/draft status are mandatory; no report is invented or silently upgraded.
- [ ] Fedora is not a prerequisite for the requested local preview; later production tasks retain their separate gates.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect docs/design/macos-preview.md. Pass when every acceptance criterion has an explicit contract section and named owning command/task.

### Gate 2: Focused behavior

Run command -v pdfinfo, command -v pdftotext, command -v sandbox-exec and a synthetic sandbox-exec /usr/bin/true probe. Pass when the selected primitives exist and execute on Darwin; record the unavailable RLIMIT_AS probe and the chosen sampled-RSS alternative without claiming a kernel hard limit.

### Gate 3: Regression

Run make verify. Pass with all offline checks and no queue warnings.

### Gate 4: Static and structural validation

Run make check and git diff --check. Pass without findings and with resolvable document links.

### Gate 5: Documentation and specification consistency

Compare SPEC, MVP, handoff and the preview contract. Pass when the new owner-requested macOS preview is explicit and does not masquerade as completed production publication.

### Gate 6: Trust-boundary and security behavior

Inspect negative-case ownership for sandbox writes/network, output bounds, cancellation, source mutation, model injection, stale reports and HTML escaping. Pass when TQ-0067 must prove these through actual interfaces and the supplied presentation asset remains unedited.

### Gate 7: Artifact and task record

Record observations and all gate evidence with tq note before done. No unavailable/unrun implementation probe is called passing; delivery follows the final phase.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work for TQ-0062: all gates PASS, done, scoped commit, confirmed push, then the next approved in-scope task. Delivery failure reopens the task and blocks continuation.

## Specification

- docs/SPEC.md#untrusted-document-handling
- docs/design/MVP.md
- docs/QUALITY_GATES.md

---

## Notes

- 2026-11-02T18:05:04+01:00 — Created directly in todo under Juan's explicit TQ-0061 request to revise existing tasks and create the missing MVP outcomes. Outcome, scope, criteria, seven gates, priority, dependencies and execution role are declared. No authority to promote unrelated inbox proposals is granted.
- 2026-11-04T22:59:53+01:00 — Juan explicitly replaced the immediate Fedora prerequisite with macOS and requested a real E2E plus HTML based on the supplied existing document. The former Fedora contract remains historical; this task now defines the private macOS preview seam. Poppler and sandbox-exec exist, a synthetic sandbox command succeeded, and RLIMIT_AS rejects a reduced address-space limit on this host; actual sampled-RSS enforcement is an implementation gate, not a fabricated hard-kernel-limit claim. Auriga was reached through the existing SSH alias and has the configured Qwen model available.
- 2026-11-04T23:02:08+01:00 — Quality gates: Gate 1 PASS — docs/design/macos-preview.md defines the owner-requested macOS seam, exact data-only presentation binding, commands, retained evidence and separate production scope. Gate 2 PASS — Darwin/arm64 Poppler pdfinfo/pdftotext 26.10.0 and sandbox-exec are available; synthetic sandboxed true succeeds. RLIMIT_AS lowering was observed to fail, so the contract explicitly specifies sampled RSS with an overshoot limitation; real enforcement belongs to TQ-0067. Gate 3 PASS — make verify passed all offline suites without queue warnings. Gate 4 PASS — make check and git diff --check passed; new relative links resolve. Gate 5 PASS — SPEC/MVP/HANDOFF explicitly reference the private macOS exception and preserve production gates. Gate 6 PASS — actual denial, bounds, cleanup, stale-report/citation and escaping checks are required in TQ-0067; no unrun implementation probe is claimed qualified and no presentation code was authored. Gate 7 PASS — contract/observations and dependent TQ-0067 are current; all evidence recorded before done.
