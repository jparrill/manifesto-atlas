---
id: TQ-0012
title: Define page-aware extraction and pinned adapter reproducibility
status: done
priority: high
assignee: gpt-6.1-sol
labels:
  - type/docs
  - agent/architect
  - area/extraction
  - area/fidelity
depends_on:
  - TQ-0008
  - TQ-0009
  - TQ-0011
created: 2026-11-01T17:07:33+01:00
updated: 2026-11-02T00:29:04+01:00
---

## Outcome

Converters and auditors share a stable page-aware representation with reproducible tool identity.

## Scope

- Define structured page/block JSON, coordinates, reading order, block types, tables, derived Markdown, page images, text coverage, warnings, and hashes.
- Pin the initial Docling version, model assets, secondary extractor, renderer, and environment after evaluating documented trial results; define artifact acquisition and resource budgets.
- Define native-text default with OCR disabled, page-level escalation signals, OCR provenance, and what exact reproducibility means despite runtime metadata.
- Define secondary evidence as independent, non-canonical output, with mapped PDF page numbering.
- Produce an approved implementation contract; do not implement the runtime feature or presentation assets.

## Acceptance criteria

- [ ] Contract examples include native text, tables, lists, multicolumn order, textless pages, suspicious pages, and OCR fallback.
- [ ] All blocks and warnings map to original PDF pages, and tool/model asset identity changes invalidate downstream artifacts.
- [ ] Unresolved or unrenderable content remains a blocker rather than being silently omitted.
- [ ] docs/design/extraction-contract.md records decisions, interfaces, versioning, rejected alternatives, and an explicit scenario-to-expected-outcome table.

## Quality Gates

### Gate 1: Acceptance criteria

Inspect `docs/design/extraction-contract.md` and its explicit scenario table. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Inspect the scenario table in docs/design/extraction-contract.md. Pass only when each of these cases has one explicit, internally consistent expected result: native PDF; image-only page; wrong page map; missing block; converter warning; model upgrade; identical pinned rerun. Cross-check examples against the declared interface; no implementation-critical decision remains implicit.

### Gate 3: Regression

Run `make verify` whenever a Makefile exists. Before it exists, run `git diff --check`, `tq list --json`, and `tq show TQ-0012 --json`; require exit zero and no malformed/skipped-queue warnings on stderr. Do not treat unavailable required tools or warning-only queue failures as a pass.

### Gate 4: Static and structural validation

Run `git diff --check`; inspect the named design artifact for all declared sections, working relative links, and internally consistent example identities. Pass with zero whitespace/structural/link findings.

### Gate 5: Documentation and specification consistency

Compare the design with every referenced Spec section and affected existing design. Record contradictions explicitly; do not silently expand product scope.

### Gate 6: Trust-boundary and security behavior

Inspect the contract's trust-boundary section and negative scenarios. Pass when untrusted inputs cannot acquire authority, unavailable evidence fails closed, and no private keys or secrets are embedded. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0012`:

1. Record all gate results with `tq note`, then run `tq done TQ-0012`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0012` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#extraction-artifacts`
- `docs/SPEC.md#extraction-fidelity`
- `docs/SPEC.md#extraction-tests`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T18:56:21+01:00 — Approved for todo under Juan's explicit delegation in this session (TQ-0057). Reviewed outcome, scope, acceptance criteria, seven quality gates, priority, native dependencies, and execution role against the accepted specification. Implementation remains gated by approved design tasks and all existing dependencies; role restrictions, required real-fixture/model checks, and fail-closed behavior remain unchanged. This approval grants no standing authority to promote future inbox proposals.
- 2026-11-02T00:29:04+01:00 — Design evidence: confirmed PyPI Docling 2.145.0 and selected library versions exist; inspected tagged upstream layout/table options; hf models info/list resolved immutable Heron and TableFormer revisions and published weight LFS SHA-256 values; hashed actual Poppler 26.10.0/poppler-data 0.4.12 source archive bytes; git ls-remote confirmed Tesseract/tessdata release refs. These establish design pins, not runtime compatibility or extraction approval. Full transitive package/component inventory, image digest, runtime qualification and golden rerun tests remain required in TQ-0030/0031. The missing original trial dependency lock is explicitly acknowledged, not reconstructed by guesswork.
- 2026-11-02T00:29:04+01:00 — Quality gates:
  - Gate 1 PASS — extraction-contract.md specifies stable page/block/table/list/column JSON, original-page coordinate/order/warning maps, native/textless/suspicious/OCR examples, exact initial converter/model/secondary/renderer/environment choices and resource/acquisition rules. Missing/unrenderable content blocks approval; every changed identity invalidates dependents.
  - Gate 2 PASS — scenario table reviewed for native PDF, image-only/blank page, wrong page map, missing block, warning despite converter success, model upgrade and identical pinned rerun. Each has one consistent outcome; OCR lineage and independent noncanonical secondary evidence are explicit.
  - Gate 3 PASS — make verify passed existing baseline/Go checks and tests.
  - Gate 4 PASS — git diff --check; python3 /tmp/opencode/check-design.py docs/design/extraction-contract.md passed UTF-8, whitespace, relative-link/anchor, section and scenario/table checks. Version/hash identities match the upstream lookups recorded in the preceding note.
  - Gate 5 PASS — checked SPEC extraction-artifacts/fidelity/tests and artifact/source/isolation contracts: original PDFs immutable, native OCR-off default, original-page evidence, retained warnings, JSON versus extracted Markdown distinction and empirically calibrated future weights preserved. No public/SPEC/presentation-code changes.
  - Gate 6 PASS — design review covers bounded offline conversion, locked pre-acquisition, no automatic model/backend substitution, hostile source as data, pixel/output limits, warning/mapping failures and independent audit. No keys/secrets/private paths embedded; no claim of having run missing golden PDFs.
  - Gate 7 PASS — design and upstream identity evidence current, deterministic payload versus run-log boundaries explicit, all gate evidence recorded; only this design/task record changed. Actual implementation and golden extraction tests remain mandatory in dependent tasks, not deferred gates of this design-only task.
