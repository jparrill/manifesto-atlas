# Manifesto Atlas Handoff

## Repository

- Remote: https://github.com/jparrill/manifesto-atlas
- Branch: `main`
- Initial repository baseline: `e10e85a` (`chore: initialize manifesto atlas repository`).
- Session work is delivered in task-scoped commits. Use `git status --short --branch` and `git log --oneline -10` for current delivery state rather than treating this document as a live Git snapshot.

Run Taskqueue from this checkout's root or any descendant directory. Running it from another repository selects that repository's queue:

```bash
tq list --json
tq ready --label agent/developer --json
```

Alternatively, set `TQ_CONFIG_PATH` to this repository's `.taskqueue.yaml` file. Never run `tq init` in a subdirectory.

## Read first

Do not reconstruct decisions from conversation history. These files are authoritative:

- [`docs/SPEC.md`](docs/SPEC.md): product, pipeline, module seams, data model, taxonomy, validation, testing, and scope.
- [`docs/QUALITY_GATES.md`](docs/QUALITY_GATES.md): mandatory task-completion gates.
- [`AGENTS.md`](AGENTS.md): owner-only inbox triage, architect/developer roles, task declaration, and closing workflow.
- [`README.md`](README.md): Spanish public project description and methodology summary.
- [`.tasks/AGENTS.md`](.tasks/AGENTS.md): generated Taskqueue CLI reference; do not edit it manually.
- [`.tasks/`](.tasks/): complete executable contracts, native dependencies, and task history. Use `tq show <id> --json` rather than reconstructing filenames.

## Current state

The repository contains documentation, Taskqueue configuration and backlog, and a root `.gitignore`. There is still no Go module, Makefile, source-ingestion implementation, schema, presentation implementation, CI workflow, GitHub Pages workflow, or project-specific agent skill.

`.gitignore` excludes `.serena/` local agent state and `.DS_Store` macOS metadata. Source PDFs, canonical JSON, reports, fixtures, task files, and this handoff remain trackable. No speculative build/cache exclusions were added before their retention contracts exist.

Completed preparation:

- `TQ-0001`: initialize Taskqueue and shared agent instructions.
- `TQ-0002`: require reproducible quality gates for executable tasks.
- `TQ-0003`: establish owner-only inbox review and execution roles.
- `TQ-0004`: create the original handoff.
- `TQ-0005`: inspect the complete specification and create the dependency-aware backlog.
- `TQ-0006`: implement and verify `.gitignore`.

`TQ-0055` records the original backlog handoff refresh. `TQ-0057` records the subsequent owner-delegated approval batch and this approval-state update.

## Backlog and approval state

The planning session created **49 backlog tasks, `TQ-0006` through `TQ-0054`**:

- `TQ-0006` is now `done` following the owner's implementation request.
- **46 backlog tasks are approved in `todo`: `TQ-0007` plus 45 tasks promoted under Juan's explicit delegation in `TQ-0057`.** Their native dependencies still control readiness.
- **`TQ-0007` is currently the only ready developer task.** It establishes a truthful Makefile interface and real baseline checks. `TQ-0018` is approved but waits for it.
- **`TQ-0008` and `TQ-0009` are ready architect tasks.** An approved architect can establish the artifact and sandbox contracts while developer work proceeds within its own scope.
- Only **two tasks remain in `inbox`**: `TQ-0017` for CI/deployment design discussion with Juan, and `TQ-0042` for human/non-LLM presentation authorship. Downstream presentation dependencies remain blocked until the latter is supplied and completed.
- The backlog has 11 architect tasks and 38 developer tasks. `TQ-0042` carries an additional human/non-LLM execution restriction.

Every backlog task has acceptance criteria, seven quality gates, specification references, declared labels, and native blocking dependencies. The audit checked **140 dependency edges**, with no missing references or cycles. The approval batch preserved all existing contracts, priorities, assignees, execution roles, and dependencies; every promoted task has an owner-authorization note.

`TQ-0005` notes contain the coverage matrix for all **40 user stories**, every implementation/testing subsection, all six planned skills, and all 12 specified Make stages:

```bash
tq show TQ-0005 --json
```

Backlog map; these ranges group related work, not an instruction to execute numerically:

| Tasks | Outcome |
| --- | --- |
| `TQ-0007` | Make help/check/verify, planned stage ownership, explicit failure for unavailable stages. |
| `TQ-0008`–`TQ-0016` | Artifact/signature trust, sandbox, provenance, source identity, extraction, taxonomy, canonical JSON, runner/validator, and site design contracts. |
| `TQ-0017` | **CI and GitHub Pages design only**, in `inbox`. No CI implementation task was created. |
| `TQ-0018`–`TQ-0027` | Go foundation, artifact validation, attestations, quarantine, safe fetching, trust anchors, semantic runner, provenance skill/decision, and immutable registration. |
| `TQ-0028`–`TQ-0029` | Recover and verify the two approved golden source PDFs. |
| `TQ-0030`–`TQ-0035` | Docling extraction, independent evidence/images, measurements, audit skill, empirical formula calibration, and extraction approval. `TQ-0034` is architect work. |
| `TQ-0036`–`TQ-0041` | Taxonomy/canonical validation, normalization skill, independent claim/coverage skills, and content lifecycle gate. |
| `TQ-0042`–`TQ-0049` | Human-authored presentation assets, deterministic Markdown/HTML, compatible regional comparisons, parity/link/browser checks, visual QA skill, and site gate. |
| `TQ-0050`–`TQ-0054` | Offline publication check, processing/preview, contribution guide, public methodology, and both-golden-source end-to-end proof. |

## Next actionable work

1. Read and claim `TQ-0007` using the actual executing model identifier:

   ```bash
   tq show TQ-0007 --json
   tq update TQ-0007 --assignee "<actual-model-id>"
   tq move TQ-0007 in-progress
   ```

2. Provide meaningful checks for the current documentation/queue checkout. All 12 planned stage names must be discoverable; unavailable stages fail non-zero rather than pretending to work. No Go module exists yet; `TQ-0018` owns that bootstrap and subsequent Go gates.
3. Approved architects may claim `TQ-0008` or `TQ-0009` now. Developers may automatically continue to the already-approved `TQ-0018` after `TQ-0007` is delivered. Later work becomes ready as its native dependencies complete; no new promotion request is needed for the approved batch.
4. Execute the [automatic closure and continuation loop](AGENTS.md#closing-tasks-and-continuing-work): all gates PASS, `done`, task-scoped commit, confirmed push, then the next approved ready task within the assignment and execution role. A delivery failure reopens the task and blocks continuation. Every task has an explicit final-phase section.

`TQ-0056` establishes this automatic delivery policy. This session's scope is backlog planning, owner-delegated approval, and the explicitly requested ignore/handoff/policy updates. Implementation agents can work through the approved backlog within their assigned scope and role.

The `TQ-0057` approval batch used Juan's explicit one-time delegation; it does not grant standing permission to promote future `inbox` proposals. The general owner-only inbox policy in `AGENTS.md` remains unchanged. `agent/architect` requires the owner-approved `gpt-6.1-sol` model or an explicitly approved equivalent. Approved architects may also execute `agent/developer` work.

## Important decisions

The complete decisions live in `docs/SPEC.md`. In particular:

- Verified source provenance is required before publication.
- Source PDFs are immutable and revision-aware.
- Canonical programme JSON, not Markdown, is the source of truth.
- Extracted Markdown is an intermediate artifact.
- Public Markdown and HTML are deterministic renderings of canonical JSON.
- HTML uses Go templates and Tokyo Night styling; LLMs never edit presentation code.
- Validation is automatic: independent claim and coverage LLM passes plus deterministic checks. Human review is not a publication gate.
- Every executable task declares and passes quality gates before `done`.
- Fresh semantic evaluations are explicit local work; normal regression and future CI verify stored evidence without rerunning expensive models.

Only the approved 2026 Frente Obrero and VOX Castilla y León documents are initial golden sources. **Their PDFs are absent from this checkout**; `TQ-0028` and `TQ-0029` recover and verify them. Do not invent URLs/checksums/metadata, substitute other revisions, or reintroduce removed PP/2023 materials. Failed recovery leaves the relevant task blocked.

Fidelity weights must follow empirical extraction/audit results, not precede them; `TQ-0034` depends on those measurements. Project skills likewise follow their approved schema/runner contracts.

**Presentation restriction:** `TQ-0042` requires Juan or a non-LLM author to supply templates/HTML/CSS/JavaScript. LLM agents cannot author these assets, including indirectly through embedded strings or generators. `TQ-0044` integrates supplied assets through Go plumbing.

**CI boundary:** review/design `TQ-0017` first. CI workflows and GitHub Pages deployment implementation require separately approved follow-up work after that design.

## Planned skills

The specification currently anticipates these project skills:

- `manifesto-source-provenance`
- `manifesto-extraction-audit`
- `manifesto-normalizer`
- `manifesto-claim-validator`
- `manifesto-coverage-validator`
- `manifesto-site-qa`

These are plans, not installed skills. Their owning tasks appear in the `TQ-0005` coverage notes; do not implement all six before their contracts exist.

## Suggested skills

A fresh agent should load the following skills when relevant:

- `implement`: execute an approved task against its specification and gates.
- `ponytail`: keep implementation minimal while satisfying the complete task contract.
- `codebase-design`: resolve approved architecture/module-interface work.
- `domain-modeling`: establish canonical vocabulary, taxonomy, states, and architecture decisions.
- `open-knowledge-write-skill` and `writing-for-agents`: author project skills after schemas exist; maintain agent instructions.
- `behavior-driven-testing`: design meaningful Go tests around observable behavior.
- `to-tickets`: refine future backlog proposals through Taskqueue; the initial full decomposition is already complete.
- `handoff`: refresh continuity information while referencing authoritative artifacts.

## Known risks

- Official provenance needs a maintained trust anchor for party domains and linked third-party hosts.
- Sandbox enforcement and signer authorization need approved contracts, not configuration-only assurances.
- LLM validators may share correlated errors even when run independently.
- Docling completed both initial PDF trials but was slow, memory-heavy, and emitted OCR warnings; converter success is not extraction approval.
- Missing golden PDFs and human-authored presentation assets are explicit future dependencies.
- A large immutable PDF archive may eventually require Git LFS, releases, or object storage while preserving checksums in Git.
- The installed `tq` binary reports a development version rather than a tagged release.
