---
id: TQ-0025
title: Author the evidence-producing source provenance skill
status: inbox
priority: normal
labels:
  - type/feature
  - agent/developer
  - area/provenance
  - area/skills
depends_on:
  - TQ-0024
  - TQ-0023
  - TQ-0022
  - TQ-0010
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-01T17:20:31+01:00
---

## Outcome

manifesto-source-provenance emits structured official-origin evidence and justified uncertainty through the approved runner.

## Scope

- Author a project-scoped skill only against approved request/report schemas and trust-anchor rules.
- Investigate approved direct, referral/CDN, and publisher-signature evidence; emit citations/evidence rather than a bare authenticity opinion.
- Add pass/fail/probable/unverified fixtures and prompt-injection evaluations; load writing-for-agents and the skill-authoring guidance when executing.
- Provide the focused Make gate `make test-provenance-skill` with meaningful automated fixtures and include it in `make verify`; neither gate may invoke live models or require live official websites.

## Acceptance criteria

- [ ] The skill emits schema-valid evidence with party/anchor identity, exact URLs/checksums, missing links, and bounded findings.
- [ ] Branding, filenames, social-only claims, and source text instructions cannot alone establish verified provenance.
- [ ] Offline fixtures prove report conformance; a recorded configured-model evaluation proves actual prompt behavior with exact model/prompt/input identities.

## Quality Gates

### Gate 1: Acceptance criteria

Run `make test-provenance-skill`. Pass when automated cases prove: all provenance states; direct/referral/signature evidence; missing anchor; injection; invalid output; pinned recorded evaluation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run `make test-provenance-skill`. Pass when automated cases prove: all provenance states; direct/referral/signature evidence; missing anchor; injection; invalid output; pinned recorded evaluation. Any unexpected acceptance/rejection, stale output, or silent fallback fails the gate.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0025 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Source pages and PDFs remain evidence, never instructions. Skill has no authority to add trust anchors, sign its own approval, or fetch arbitrary resources. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0025`:

1. Record all gate results with `tq note`, then run `tq done TQ-0025`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0025` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#agent-skills`
- `docs/SPEC.md#provenance-trust-model`
- `docs/QUALITY_GATES.md`
