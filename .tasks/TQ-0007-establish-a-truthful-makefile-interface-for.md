---
id: TQ-0007
title: Establish a truthful Makefile interface for current and planned stages
status: done
priority: high
assignee: auriga-prime/qwen3.8-27b-dflash2-q4
labels:
  - type/chore
  - agent/developer
  - area/project
  - area/cli
created: 2026-11-01T17:07:32+01:00
updated: 2026-11-01T23:07:20+01:00
---

## Outcome

A fresh checkout exposes discoverable Make commands and real baseline checks, with a clear expansion path for every specified pipeline stage.

## Scope

- Inspect the complete Spec and newly created backlog; inventory submit, verify-source, add, extract, validate-extraction, normalize, validate-content, build, validate-site, publish-check, process, and serve.
- Provide help, check, and verify now using actual available documentation/Git/Taskqueue checks; reject malformed-queue warnings even if tq exits zero.
- List future stage targets and their owning task IDs. Invocation of an unavailable stage must fail non-zero with an actionable explanation, never return a fake success.
- Keep deterministic implementation in Go once introduced and orchestration in Make. Later implementation tasks replace the unavailable-stage failure and extend real gates.
- No CI workflow or pipeline implementation is included.

## Acceptance criteria

- [ ] make help lists all 12 Spec pipeline targets and distinguishes implemented checks from planned stages.
- [ ] make check and make verify perform meaningful current checks and exit zero on a valid checkout without requiring a nonexistent Go module.
- [ ] Whitespace errors and a malformed task fixture cause the applicable baseline gate to fail non-zero, including tq warning-only failures.
- [ ] Every unavailable pipeline target exits non-zero; invoking make without arguments is safe and does not publish or fetch anything.
- [ ] Technical usage documentation maps planned targets to backlog owners and explains how Go/model/browser checks will be integrated.

## Quality Gates

### Gate 1: Acceptance criteria

Run make help, make check, and make verify. In a disposable repository copy inject a whitespace error and a malformed task, then assert failure. Invoke each unavailable stage and assert non-zero with its owning task reference. Verify each acceptance criterion against its corresponding command output, fixture, or named artifact. Pass only when all stated behaviors and failure outcomes hold; record the evidence mapping.

### Gate 2: Focused behavior

Run make help, make check, and make verify. In a disposable repository copy inject a whitespace error and a malformed task, then assert failure. Invoke each unavailable stage and assert non-zero with its owning task reference.

### Gate 3: Regression

Run `make verify`, `git diff --check`, `tq list --json`, and `tq show TQ-0007 --json`. Require exit zero and no malformed/skipped-queue warnings on stderr; repeat regression after any gate fix.

### Gate 4: Static and structural validation

Run `git diff --check` and `make check` when the Makefile exists; for Go changes include formatting, compilation, vet/static checks, and schema/fixture validation in that target. Pass with zero errors or warning-or-higher findings. Before a Makefile exists, verify changed file structure and relative links directly and record the exact checks.

### Gate 5: Documentation and specification consistency

Compare changed interfaces, examples, Make help/targets, schemas, fixtures, and public/technical documentation with the referenced Spec sections and approved design contracts. Pass when commands and artifact states match, public text is Spanish, technical text is English, and no product decision is silently changed.

### Gate 6: Trust-boundary and security behavior

Check shell quoting with a checkout path containing spaces, propagate command failures, and verify default/help/check targets perform no network retrieval, model invocation, or publication. Pass only when all negative cases produce the specified safe outcome and no secrets/private absolute paths are committed.

### Gate 7: Artifact and task record

Verify required artifacts and fixtures are current, repeated identical-input checks satisfy the declared determinism policy, and no unrelated user work is removed. Record every gate as PASS with commands/artifact evidence using `tq note` before `tq done`. A failed, flaky, skipped, or unavailable required gate keeps the task open.


## Final phase: delivery and continuation

After every quality gate passes, follow [AGENTS.md: Closing tasks and continuing work](../AGENTS.md#closing-tasks-and-continuing-work) for `TQ-0007`:

1. Record all gate results with `tq note`, then run `tq done TQ-0007`.
2. Commit only this task's intended changes and final task record, push to the configured upstream, and confirm the remote contains the commit.
3. Only after confirmed delivery, automatically read, assign, and claim the next approved ready `todo` task within the current assignment, execution role, and additional execution restrictions.

A commit/push/confirmation failure restores `TQ-0007` to `in-progress` with a blocker note and stops continuation. The linked procedure defines recovery, Git checks, owner-only inbox triage, and scope restrictions. Historical completed execution is not retroactively repeated.

## Specification

- `docs/SPEC.md#make-targets`
- `docs/SPEC.md#module-seams`
- `docs/SPEC.md#publication-and-ci`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-01T19:32:18+01:00 — Quality gates:
  - Gate 1 PASS — acceptance criteria verified. 'make help' lists all 12 Spec pipeline targets plus the 3 implemented baseline targets, distinguishing implemented from planned. 'make check'/'make verify' exit 0 on the valid checkout using only git+tq (no Go module). Disposable-copy tests: injected trailing whitespace into README.md -> 'make check' failed non-zero (git diff --check: README.md:159 trailing whitespace); injected malformed task file (invalid YAML frontmatter 'id: [unclosed') -> 'make check' failed non-zero because 'tq list --json' exits 0 but emits 'warning: skipped TQ-9998-bad.md' on stderr (warning-only failure rejected); all 12 unavailable stages invoked -> each exited non-zero printing its owning TQ id; 'make' with no args ran default goal (help), exit 0, no fetch/publish. docs/MAKEFILE.md maps every planned target to its owning task and explains Go/model/browser integration.
  - Gate 2 PASS — same focused commands as Gate 1: make help / make check / make verify exit 0 on valid checkout; whitespace fixture and malformed-task fixture each forced non-zero; all 12 unavailable stages non-zero with owning-task reference (12/12).
  - Gate 3 PASS — 'make verify' exit 0; 'git diff --check' exit 0; 'tq list --json' exit 0 with empty stderr; 'tq show TQ-0007 --json' exit 0 with empty stderr. No malformed/skipped-queue warnings.
  - Gate 4 PASS — 'git diff --check' exit 0; 'make check' exit 0; 'make -n help check verify submit serve' dry-run clean (no make-level warnings on stderr). No Go changes, so no gofmt/vet/test. checkmake not installed; make -n dry-run used as the structural static check.
  - Gate 5 PASS — the 12 Makefile pipeline targets are byte-identical (as a set) to the 12 'docs/SPEC.md#make-targets' entries (diff of sorted lists empty). Changed files are technical (Makefile, docs/MAKEFILE.md) and written in English; no public/Spanish file changed; docs/SPEC.md, README.md, HANDOFF.md untouched, so no product decision changed.
  - Gate 6 PASS — 'make check' and 'make verify' both exit 0 when the checkout path contains a space ('disposable with space'); failure propagation proven by the whitespace and malformed-task fixtures; baseline recipes invoke only git diff --check, tq list, grep, dirname, test, mktemp, cat, printf, echo, rm (no curl/wget/git fetch/push/model calls — the only http strings are the link-check 'case' that skips URLs); no secrets or private absolute paths in Makefile or docs/MAKEFILE.md.
  - Gate 7 PASS — artifacts current (Makefile + docs/MAKEFILE.md); 'make verify' is deterministic (two identical-input runs produced byte-identical output, exit 0 both); task record complete (this note).
- 2026-11-01T19:38:19+01:00 — BLOCKER (delivery): 'git commit' failed at GPG signing — 'gpg: signing failed: Timeout'. The configured signing key user.signingkey=DE30627E53AC3969 requires a passphrase; the curses pinentry launched and timed out in the automated context, and the key is not cached in the gpg-agent (loopback test without passphrase returned 'Bad passphrase'). All 7 gates already PASS and the 3 intended files (Makefile, docs/MAKEFILE.md, .tasks/TQ-0007-*.md) are staged. Per AGENTS.md I did NOT weaken the signature (no --no-gpg-sign, no commit.gpgsign change, hooks intact). Task restored to in-progress. Resolution needs owner input: unlock/cache the signing key passphrase in the gpg-agent (e.g. run an interactive 'gpg' sign or 'gpg-connect-agent' preset-passphrase), then re-run the commit+push. After a successful signed commit and confirmed push, mark done and continue.
- 2026-11-01T23:07:20+01:00 — DELIVERED: signing key DE30627E53AC3969 was cached in the gpg-agent after owner unlock; signed commit cf157f8 'chore: establish truthful Makefile interface (TQ-0007)' (3 files: Makefile, docs/MAKEFILE.md, task record) pushed to origin/main and confirmed (local main == origin/main == cf157f8). No unpushed commits. Task marked done.
