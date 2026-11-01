---
id: TQ-0053
title: Publish Spanish methodology and automated-validation limitations
status: inbox
priority: normal
labels:
  - type/docs
  - agent/developer
  - area/project
  - area/validation
  - area/site
depends_on:
  - TQ-0034
  - TQ-0050
  - TQ-0045
created: 2026-11-01T17:07:34+01:00
updated: 2026-11-01T17:20:32+01:00
---

## Outcome

Readers can inspect the real provenance, extraction, content, rendering, and publication methodology without mistaking automation for certainty.

## Scope

- Write Spanish public Markdown explaining evidence chains, immutable revisions, topic states/jurisdiction, calibrated fidelity formula, independent validators, report identities, and fail-closed gates.
- Disclose probabilistic/correlated model errors, tool/resource limits, unresolved evidence, and archive storage limitations.
- Separate factual summaries from future viability/historical reliability/observed fulfilment/confidence; do not invent a political score or claim to fact-check all party statistics.
- Integrate methodology content through existing canonical/content/render interfaces without authoring or editing presentation assets.

## Acceptance criteria

- [ ] Public methodology accurately reproduces the approved formula/version and explains blockers overriding aggregate scores.
- [ ] It states both independent content passes are required, disagreement/insufficient evidence prevents publication, and no human-review publication fallback exists.
- [ ] Generated public pages link to the methodology, source evidence, and revision/validation identity; third-party documents retain publisher rights.
- [ ] Spanish text matches actual implementation and contains no unsupported infallibility or fulfilment claims.

## Quality Gates

### Gate 1: Acceptance criteria

Compare public prose and worked examples with the approved calibration artifact and current reports; run make verify and the real site parity/link gates. Pass when formula/examples, statuses, links, and limitation statements match exactly. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Compare public prose and worked examples with the approved calibration artifact and current reports; run make verify and the real site parity/link gates. Pass when formula/examples, statuses, links, and limitation statements match exactly.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0053 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Cross-check README, SPEC, calibration, trust/validation/site contracts, and generated methodology links; all public explanatory prose must be Spanish.

### Gate 6: Trust-boundary and security behavior

Ensure public documentation exposes no signing secrets or local private paths and cannot present unverified content or model opinion as source authority. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0053`:

1. Record all gate results with `tq note`, then run `tq done TQ-0053`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0053` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#language-and-licensing`
- `docs/SPEC.md#assessment-separation`
- `docs/SPEC.md#extraction-fidelity`
- `docs/SPEC.md#llm-validation`
- `docs/SPEC.md#further-notes`
- `docs/QUALITY_GATES.md`
