---
id: TQ-0067
title: Run a real macOS PDF-to-HTML preview with independent model checks
status: todo
priority: high
labels:
  - type/feature
  - area/extraction
  - area/normalization
  - area/rendering
  - agent/developer
depends_on:
  - TQ-0062
created: 2026-11-04T23:00:21+01:00
updated: 2026-11-04T23:00:21+01:00
---

## Outcome

One command on macOS processes a retained real PDF into cited canonical JSON, actual independent semantic reports, Markdown and a locally viewable HTML derived from the owner-supplied reference.

## Scope

- Implement docs/design/macos-preview.md as a narrow private-preview seam using current candidate/canonical/taxonomy/semantic helpers.
- Bound macOS PDF subprocesses with sandbox-exec, deadlines/process-group cleanup, output limits and sampled RSS enforcement. Retain source hashes, raw page evidence and actual tool identities.
- Run the configured model in fresh normalizer, claim and coverage conversations; retain genuine replies and require exact evidence/citation bindings. Missing/failed/incomplete reports must remain non-eligible, while a structurally valid private preview can still be inspected.
- Import the owner-supplied HTML byte-for-byte and replace only its programme JSON data payload during deterministic Go rendering; author no HTML, CSS, JavaScript or template markup. Escape all inserted strings.
- Provide make e2e-preview, make render-preview and focused tests. Document exact commands, resulting files, limitations and private-preview status. Do not mark unrelated production tasks complete.

## Acceptance criteria

- [ ] A real Frente Obrero run on macOS retains original source identity, all original pages, structurally valid canonical JSON, page citations, independent report outcomes and model/tool/input identities.
- [ ] Local HTML and Markdown reflect the same canonical data and show genuine draft/validation status; unprocessed party tabs are explicitly labelled pending without old unsupported summaries.
- [ ] Supplied asset bytes remain unchanged; deterministic rerender uses stored inputs without model calls and hostile inserted text cannot execute markup or script.
- [ ] Source mutation, malformed model output, invalid citation, missing page/check inventory, stale input, cancellation, excess resources and denied worker writes/network fail safely.
- [ ] A real desktop/mobile browser renders the produced local preview; no public eligibility is claimed without both current independent passes and the separate production gates.

## Quality Gates

### Gate 1: Acceptance criteria

Run make evaluate-preview with the real retained source, configured Auriga provider and supplied HTML. Pass when original bounded replies, source/tool/configuration hashes, canonical validation, independent outcomes, Markdown, HTML and run record are retained and consistent. A semantic fail/insufficient result is visible draft status, never an invented pass.

### Gate 2: Focused behavior

Run make test-preview. Pass with named behavioral cases for literal source citations, schema/state rules, stale/missing reports, data-only template binding, injection escaping, determinism, subprocess cancellation/bounds and actual macOS sandbox denial probes.

### Gate 3: Regression

Run make verify and go test ./.... Pass with exit zero and no live model calls in offline regression or Taskqueue warnings.

### Gate 4: Static and structural validation

Run make check and git diff --check. Pass with zero formatting/build/vet/queue/whitespace findings and valid retained JSON.

### Gate 5: Documentation and specification consistency

Compare CLI/Make help, SPEC/MVP exception, preview contract and handoff with actual commands/artifacts. Pass when macOS/private-preview scope, Poppler use, supplied presentation ownership and pending production capabilities are explicit.

### Gate 6: Trust-boundary and security behavior

Exercise real sandbox write/network denial, timeout/descendant cleanup, sampled RSS termination and output bounds; exercise hostile text/quotes, unsafe paths, malformed/stale reports and unsupported citations. Pass only with safe failures and no false public approval or secret retention.

### Gate 7: Artifact and task record

Retained real run evidence and deterministic rerender are current. Record all results with tq note before done; unavailable required gates keep this task open.

## Final phase: delivery and continuation

Follow AGENTS.md#closing-tasks-and-continuing-work for TQ-0067: all gates PASS, done, scoped commit, push and remote confirmation. Continue only approved in-scope work; delivery failures reopen this task.

## Specification

- docs/SPEC.md#deterministic-rendering
- docs/design/MVP.md
- docs/design/macos-preview.md
- docs/QUALITY_GATES.md
