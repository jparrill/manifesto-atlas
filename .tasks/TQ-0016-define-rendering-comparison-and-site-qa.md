---
id: TQ-0016
title: Define rendering, comparison, and site QA contracts
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - agent/architect
  - area/rendering
  - area/site
depends_on:
  - TQ-0008
  - TQ-0014
  - TQ-0013
  - TQ-0015
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T18:05:04+01:00
---

## Outcome

Renderers and QA share an exact Spanish content, accessibility, responsive, comparison, and visual approval contract.

## Scope

- Specify required content/landmarks, every-topic ordering and status language, citations/PDF links, Tokyo Night design requirements, escaping, and JSON-to-Markdown/DOM parity.
- Define compatible regional comparisons without merging incompatible elections, geography, document scope, taxonomy versions, or revisions.
- Choose concrete desktop/mobile viewports, keyboard/a11y checks, browser/font/tool pinning, screenshot retention, visual-baseline updates, thresholds, and visual report hashes.
- Specify human/non-LLM authorship of HTML/CSS/JavaScript/templates; a developer LLM may implement Go renderer plumbing but cannot author presentation assets.
- Keep template/design changes and QA approval independent; visual approval cannot override deterministic failures.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Worked public cases cover proposals, mention, absence, non-applicability, and unclear without inferring political positions.
- [ ] QA decisions have reproducible pass conditions and input-bound reports; browser/baseline updates invalidate dependent approvals.
- [ ] Comparison behavior and presentation execution restrictions are explicit enough to prevent incompatible comparisons or prohibited agent edits.
- [ ] docs/design/site-quality.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/site-quality.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/site-quality.md. Pass only when each of these cases has one explicit, internally consistent expected result: all topic states; long text; hostile markup; keyboard navigation; narrow screen; missing citation; incompatible comparison; stale screenshot; visual pass with deterministic failure. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0016 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0016`:

1. Record all gate results with `tq note`, then run `tq done TQ-0016`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0016` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#deterministic-rendering`
- `docs/SPEC.md#site-validation`
- `docs/SPEC.md#site-quality-tests`
- `docs/SPEC.md#topic-taxonomy`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:22+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T00:37:41+01:00 — Quality gates:
  - Gate 1 PASS — site-quality.md defines all-topic Spanish content/state examples, deterministic Markdown/DOM parity, source/citation/rights/methodology content, compatible regional comparison evidence, strict non-LLM presentation authorship, and binary input-bound keyboard/a11y/reflow/screenshot/visual gates.
  - Gate 2 PASS — scenarios reviewed: all topic states, long text, hostile quotes, keyboard navigation, narrow screens, missing citations, incompatible comparisons, stale screenshots and visual pass with deterministic failure have explicit consistent results. Unclear content is preview-only, never published.
  - Gate 3 PASS — make verify passed all existing checks/tests.
  - Gate 4 PASS — git diff --check; design diagnostic passed UTF-8/whitespace, links/anchors and scenario/alternative/table structure. Verified exact Chromium 141.0.7390.37/revision 1194 against upstream Playwright v1.56.1 browser manifest; actual tool-byte sealing is required in TQ-0047.
  - Gate 5 PASS — reviewed SPEC deterministic-rendering/site-validation/site-quality-tests/taxonomy and canonical/semantic/artifact designs. No political inference, omitted topics, merged source identities, model-generated assets or human programme-approval fallback introduced. Design prose only, no HTML/CSS/JS/templates authored.
  - Gate 6 PASS — design trust review covers contextual escaping, canonical-derived links, no arbitrary paths/network/credentials, isolated browser, hostile literal text, screenshot stale/hash checks, and independent composite approval. Missing supplied assets/isolation/baselines fail relevant runtime gates. No secrets/private paths.
  - Gate 7 PASS — current contract, concrete viewport/tool/threshold/baseline rules and gate notes recorded; deterministic inputs/version invalidation explicit; only design/task record changed. Real browser screenshots/a11y/pixel fixtures remain mandatory implementation gates, not claimed executed by this design task.
- 2026-11-02T18:05:04+01:00 — Historical completion retained. Juan's TQ-0061 scope replan establishes docs/SPEC.md and docs/design/MVP.md as the current MVP contract; conflicting graph/PKI/network admission/legal-registry/calibrated-score/browser-pinning requirements in this task's original design are superseded for MVP work. This note does not recertify any historical implementation or erase prior execution evidence. TQ-0063 owns removal/adaptation of inactive platform code.
