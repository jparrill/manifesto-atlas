# Manifesto Atlas Agent Instructions

## Task management

Read and follow [`.tasks/AGENTS.md`](.tasks/AGENTS.md) for the authoritative Taskqueue CLI workflow, statuses, priorities, labels, and file rules.

Taskqueue is the source of truth for planned and active work.

### Owner-only inbox triage

- `inbox` is Juan's review queue. Agents must not claim, execute, reject, or promote inbox tasks.
- Agents may create discovered follow-up work only in `inbox`. They must not move it to `todo`.
- Only Juan decides whether an inbox task is accepted and moves it to `todo`.
- Before promotion, Juan confirms the outcome, scope, acceptance criteria, quality gates, priority, dependencies, and execution role.
- Work explicitly assigned by Juan may be created directly as `todo` or `in-progress` with the required fields complete.

### Execution roles

Every task in `todo` or `in-progress` must have exactly one execution-role label:

- `agent/developer`: implementation work against an approved specification or design. A suitable local model or GPT coding agent may claim it.
- `agent/architect`: architecture, system design, cross-module interfaces, canonical schemas, taxonomy, trust models, validator contracts, or changes that establish implementation direction. Only a model explicitly approved by Juan for architecture work may claim it. The canonical approved model identifier is `gpt-6.1-sol`; equivalent quality is determined by Juan, not by the agent itself.

An approved architect may execute developer tasks. A developer agent must never execute an `agent/architect` task.

If a developer task requires an unresolved architectural decision, stop implementation, record the blocker with `tq note`, and create an architect follow-up in `inbox` for Juan to review. If one task mixes substantial design and implementation, split it into an architect task followed by dependent developer tasks.

### Selecting and claiming work

Developer agents select only developer work:

```bash
tq ready --label agent/developer --json
```

Approved architect agents select architecture work with:

```bash
tq ready --label agent/architect --json
```

Before the first project edit:

1. Read the complete task with `tq show <id> --json`.
2. Confirm the execution-role label matches the current agent.
3. Record the actual model or agent identifier with `tq update <id> --assignee "<model-or-agent>"`.
4. Claim it with `tq move <id> in-progress`.

While working:

- Record decisions, surprises, blockers, and deferred follow-up with `tq note`.
- Use dependencies when one task cannot complete before another; do not encode ordering only in prose.
- Create a separate inbox task for discovered work not required for the current outcome.
- Use the Taskqueue CLI for metadata and lifecycle changes; do not edit task frontmatter or filenames manually.
- Write task titles, bodies, notes, assignees, and labels in English.

## Declaring tasks

Each task must describe one independently verifiable outcome, reference the relevant section of [`docs/SPEC.md`](docs/SPEC.md), and follow [`docs/QUALITY_GATES.md`](docs/QUALITY_GATES.md). Prefer a small vertical slice over a broad phase or a single technical layer.

Executable tasks require:

- Exactly one declared `type/*` label.
- At least one declared `area/*` label.
- Exactly one `agent/developer` or `agent/architect` label.
- Complete acceptance criteria and quality gates.

Use this body structure:

```markdown
## Outcome

The observable result this task must produce.

## Scope

- Included work.
- Explicitly excluded work when the boundary is not obvious.

## Acceptance criteria

- [ ] A checkable externally observable behavior.
- [ ] All required callers, schemas, fixtures, and documentation are updated.

## Quality Gates

### Gate 1: Acceptance criteria

- Commands or artifacts that prove every acceptance criterion.
- Explicit pass conditions.

### Gate 2: Focused behavior

- Focused tests or deterministic checks for the changed behavior.
- Explicit pass conditions.

### Gate 3: Regression

- Relevant existing test suites and `make verify` when available.
- Explicit pass conditions.

### Gate 4: Static and structural validation

- Compiler, formatter, linter, schema, diagnostics, or artifact checks.
- Explicit pass conditions.

### Gate 5: Documentation and specification consistency

- Documentation checks or an explicit `N/A` rationale.

### Gate 6: Trust-boundary and security behavior

- Invalid-input and failure-path checks, or an explicit `N/A` rationale.

### Gate 7: Artifact and task record

- Generated artifacts are current and deterministic.
- Gate results are recorded with `tq note` before completion.

## Specification

- `docs/SPEC.md#relevant-section`
```

Acceptance criteria must describe behavior rather than implementation steps. Quality gates must be binary and reproducible: each gate names the command or artifact and its expected result. A task is not complete merely because code was written.

## Closing tasks

Before moving any task to `done`:

1. Run every gate in its `## Quality Gates` section.
2. Fix failures and rerun the affected gate plus the regression gate.
3. Record every result as `PASS` with evidence using `tq note`.
4. Leave the task open when a required gate is failed, flaky, skipped, or unavailable.

A follow-up task may capture optional improvement, but it cannot defer a gate required for the current task's outcome.

## Project rules

- Canonical programme JSON is the source of truth for public Markdown and HTML.
- Deterministic operations belong in Go code and Make targets; semantic judgment belongs in schema-constrained agent skills.
- LLMs must not generate or edit HTML, CSS, JavaScript, or templates.
- Public-facing content is Spanish. Code, comments, schemas, prompts, tests, task content, and technical documentation are English.
- Source PDFs are immutable after registration. Never overwrite committed source bytes.
