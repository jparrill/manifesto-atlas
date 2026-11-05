---
id: TQ-0067
title: Run a real macOS PDF-to-HTML preview with independent model checks
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/feature
  - area/extraction
  - area/normalization
  - area/rendering
  - agent/developer
depends_on:
  - TQ-0062
created: 2026-11-04T23:00:21+01:00
updated: 2026-11-05T01:05:19+01:00
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

---

## Notes

- 2026-11-05T01:01:06+01:00 — Delivered real macOS private E2E under runs/preview/frente-e2e: 64 retained original pages, 21 ordered topics, 17 principal measures, exact host-resolved source citations, three genuine independent Auriga calls, canonical JSON, Markdown and supplied-reference HTML. Model claim result was pass but host conservatively marks it insufficient because original pages 1/2/64 lack visual evidence; coverage is insufficient and flags possible omissions. Public eligibility remains false. Earlier sandbox failure, cancelled proxy request, nonliteral quote and malformed finding-schema attempts are preserved separately, never relabelled passing. The original owner reference was imported byte-identically; all presentation/script bytes outside the programme data payload remain unchanged. Fedora is not a prerequisite. Follow-up evidence belongs to existing original-page/content tasks; production CLI/registry/publication integration is not claimed complete.
  Quality gates:
  - Gate 1 PASS — actual make evaluate-preview OUTPUT=runs/preview/frente-e2e through the configured Auriga endpoint completed normalization, independent checks and private rendering. Genuine bounded requests/replies, source/model/tool hashes and expected non-publishable outcomes are retained.
  - Gate 2 PASS — make test-preview via make verify passed literal Unicode/evidence-ID binding, source/schema validation, all-topic data parity, deterministic offline rerender, missing/stale model/source/draft/report rejection, all-outcome report bindings, symlink/hard-link-safe output replacement, 100 MiB PDF vs 16 MiB JSON bounds, complete retry-attempt retention, escaping and native sandbox/timeout/RSS cases.
  - Gate 3 PASS — make verify and go test ./... passed all required offline suites without live model calls or queue warnings. The separate actual-tool gate ATLAS_EVALUATE_REAL_PDF=1 go test ./internal/preview -run TestActualMacSource -count=1 -v passed and exactly matched all 64 retained pages/tools/warnings; it was not substituted with a skipped test.
  - Gate 4 PASS — formatting, build, vet, queue, documentation/link checks, git diff --check and staged whitespace validation passed. Retained JSON is parsed by current helpers; actual Chrome desktop/mobile checks passed.
  - Gate 5 PASS — README, MAKEFILE usage, SPEC/MVP exception and handoff match actual private commands and limitations. Production targets remain truthfully unavailable; TQ-0021 now references macOS, and TQ-0042 only records the supplied reference without unauthorized inbox promotion. Standards review: two findings resolved. Spec review: three findings resolved, including one shared finding; both reviewers confirmed the fixes.
  - Gate 6 PASS — actual native positives and write/private-read/network denial, host-credential exclusion, parent-process inspection denial, default fork denial, whole-process-group timeout/cleanup, output ceiling and sampled RSS termination passed. HTML/Markdown escaping, unsafe paths/links, stale and non-passing report substitution, output symlink/hard-link preservation and retry evidence are covered. Provider URL credentials are rejected and no signing key/API credential is retained. RSS is explicitly sampled, not claimed as a kernel hard limit.
  - Gate 7 PASS — latest real run, preserved failures, immutable owner asset and desktop/mobile screenshots are current. CDP proves 21 rendered cards and no horizontal overflow at 1440x1000 and 390x844. Offline replay repeats identically; source bytes are preserved. All evidence is recorded before done; scoped commit/push confirmation follows.
- 2026-11-05T01:04:02+01:00 — Delivery blocker: the combined commit/push command exceeded its 120-second tool deadline after commit succeeded. Remote inspection still shows the prior contract commit and local HEAD is one commit ahead. No task continuation is allowed until normal push and remote confirmation succeed; the task is reopened under the delivery policy. Source/HTML artifacts and all content gates remain intact.
- 2026-11-05T01:05:19+01:00 — Delivery recovery: retrying the normal configured git push succeeded and remote confirmation matches the delivered preview commit. No code or evidence was changed during recovery. Regression make verify and git diff --check passed again; queue/task records remain valid. All seven pre-delivery gates remain PASS. Closing the reopened task and delivering this recovery record as a new commit, without amending history or weakening hooks/signing.
