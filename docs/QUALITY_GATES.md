# Task Completion Quality Gates

Every Taskqueue task must prove its outcome before it moves to `done`. Quality gates are part of the task contract, not a checklist added after implementation.

This document defines the project-wide policy. Each task body defines the concrete gates that apply to its outcome.

After the gates pass, every task must execute its final delivery phase under [the automatic closure and continuation procedure](../AGENTS.md#closing-tasks-and-continuing-work). A local `done` status alone does not prove committed and pushed delivery.

## Readiness and ownership

Inbox tasks are unapproved proposals in Juan's review queue. They may be incomplete and agents must not execute or promote them.

Before Juan moves a task to `todo`, it must have complete acceptance criteria, quality gates, priority, dependencies, and exactly one `agent/developer` or `agent/architect` execution-role label. Only tasks in `todo` are available for autonomous agent selection.

## Gate contract

A task in `todo` or `in-progress` must contain a numbered `## Quality Gates` section. Each gate must state:

1. **What is checked.** An observable behavior, invariant, or artifact.
2. **How it is checked.** An exact command, test, comparison, or generated report.
3. **What passes.** An exit code, expected value, absence of findings, or artifact condition.

A gate such as "review the output" is invalid because it has no reproducible pass condition.

A task may use `N/A` only when the gate category genuinely cannot apply, and the task note must record why. A failed, flaky, skipped, or unavailable required gate blocks completion.

## Mandatory gates

Every task must cover these gate categories, combining them only when one command proves several categories.

### Gate 1: Acceptance criteria

Prove every acceptance criterion in the task body. Written code or documentation is not itself proof that the requested behavior exists.

### Gate 2: Focused behavior

Run the smallest tests or checks that directly exercise the changed behavior. New non-trivial logic requires a focused automated test or an assert-based deterministic check.

### Gate 3: Regression

Run the relevant existing test suites and confirm the change does not break connected behavior. Once the project Makefile provides `make verify`, it is mandatory for every task.

Before `make verify` exists, the task must list the exact available regression commands instead of using a placeholder or no-op gate.

### Gate 4: Static and structural validation

Run the applicable compiler, formatter, linter, schema validator, diagnostics, and generated-artifact consistency checks. Warning-or-higher diagnostics in changed project files block completion unless the task explicitly exists to remove them in a later required step.

### Gate 5: Documentation and specification consistency

Confirm that affected public interfaces, schemas, prompts, Make targets, task instructions, and technical documentation agree with the implementation. Update `docs/SPEC.md` only when the product contract changes, not for implementation details already discoverable from code.

### Gate 6: Trust-boundary and security behavior

For changes handling PDFs, URLs, metadata, model output, generated HTML, filesystem paths, or external commands, test invalid and hostile input as well as the success path. Verify timeouts, limits, escaping, and failure behavior at the changed trust boundary.

Use `N/A` with a recorded rationale only when the task does not touch a trust boundary.

### Gate 7: Artifact and task record

Confirm required generated artifacts are current and deterministic. Record all gate results with `tq note` before moving the task to `done`.

Use this note format:

```text
Quality gates:
- Gate 1 PASS — acceptance criteria verified by <evidence>.
- Gate 2 PASS — `<command>`.
- Gate 3 PASS — `<command>`.
- Gate 4 PASS — <diagnostic or artifact report>.
- Gate 5 PASS — <documentation check or N/A rationale>.
- Gate 6 PASS — <security check or N/A rationale>.
- Gate 7 PASS — artifacts current and task record complete.
```

## Conditional gates

Add these gates when the task touches the corresponding area.

### Go implementation

- Format changed Go files.
- Run focused package tests.
- Run the repository-wide Go test and static-analysis targets.
- Check errors, cancellation, and resource cleanup at external seams.

### Schemas, canonical data, and taxonomy

- Validate positive and negative fixtures.
- Verify version changes invalidate stale downstream artifacts.
- Confirm deterministic serialization.
- Confirm public Markdown and HTML remain derivable from canonical JSON.

### Agent skills, prompts, and validators

- Validate output against its schema.
- Run pass, fail, disagreement, and insufficient-evidence fixtures where applicable.
- Test prompt-injection content as untrusted data.
- Record model, prompt, schema, taxonomy, and input hashes.
- Confirm one model invocation cannot approve its own generated result.

### Source provenance and PDF extraction

- Verify checksums and immutable-source behavior.
- Run the applicable golden PDF fixtures.
- Check page-level provenance and extraction warnings.
- Test a malformed, mismatched, or unverified source.
- Confirm blocking fidelity findings override aggregate scores.

### Rendering and website changes

- Build the static site.
- Validate canonical JSON-to-DOM content parity.
- Validate links, page citations, HTML structure, and escaping.
- Run accessibility and keyboard-navigation checks.
- Capture desktop and mobile screenshots.
- Run visual regression and the schema-constrained visual QA skill when available.

### CI and build changes

- Run the changed target locally.
- Verify failure propagation and non-zero exit codes.
- Confirm CI uses the same project target rather than duplicating command logic.
- Confirm cached or stored validation reports are rejected when their input hashes are stale.

### Documentation-only changes

Documentation tasks still run structural diagnostics, link checks when available, Taskqueue integrity checks, and any commands needed to prove documented examples are correct. Documentation-only is not an exemption from quality gates.

## Closing workflow

Follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work). It defines the ordered gate-recording, `done`, task-scoped commit, confirmed push, failure recovery, and automatic next-task selection procedure. Include a `## Final phase: delivery and continuation` reference in every task body; delivery runs after the seven gates rather than creating a circular gate that requires its own future commit hash.

A follow-up task may track optional improvement. It cannot be used to defer a failed gate required for the current outcome.
