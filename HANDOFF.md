# Manifesto Atlas Handoff

## Execution mission: continue autonomously until the approved queue is exhausted

**Juan has explicitly authorized the resuming implementation agent to execute
the approved local-first MVP backlog. Start implementing after reading the
current contracts; do not stop at a summary, an audit or another proposed plan.**

**Do not ask questions or request permission to start, choose the next approved
task, run its checks, commit, push or continue. Do not ask whether to continue
after completing a task. Continue automatically with all eligible approved tasks
until none remain for your execution role.** Routine delivery and continuation
are already authorized in [AGENTS.md](AGENTS.md#closing-tasks-and-continuing-work).

Use the specification, task criteria and existing patterns to resolve ordinary
implementation choices. Choose the smallest conforming solution, record relevant
assumptions with `tq note`, and proceed. Do not invent a new approval checkpoint,
expand the MVP, or restart planning because a historical contract was larger.

## Read first

- [SPEC](docs/SPEC.md) and [current MVP contract](docs/design/MVP.md) define the
  owner-approved reduced scope under TQ-0061.
- [AGENTS.md](AGENTS.md), [quality gates](docs/QUALITY_GATES.md) and the generated
  [Taskqueue guide](.tasks/AGENTS.md) define execution/role/delivery rules.
- Taskqueue is the current work source; use `tq show <id> --json` before claiming.
- The [scope audit](docs/audits/2026-11-02-product-scope.md) explains the change.

The nine older design documents are explicitly historical for conflicting MVP
requirements. Do not rebuild their DAG, report PKI, arbitrary-URL admission,
universal legal registry, calibrated score or visual-model QA because an old
completed task mentions it. Completed tasks retain their original history.

## Current implementation, not planned capability

Current owner assignment: implement the real-source macOS E2E/HTML preview under
[the macOS contract](docs/design/macos-preview.md), TQ-0062/TQ-0067. Use the supplied
existing HTML as immutable presentation input; do not wait for Fedora. The live
queue and retained run evidence take precedence over the older snapshot below.

The repository contains a Go CLI, tested domain libraries, Make gates and design/
task documentation. The ordinary executable still provides help/version and
unavailable production pipeline-stage messages. Both candidate PDFs and origin
evidence are now committed under `fixtures/`; the configured semantic provider
and its real synthetic smoke evidence are implemented. There is still no completed
production public output or CI workflow. Passing library tests are not end-to-end
product readiness.

The previous technical audit reopened TQ-0022/TQ-0023/TQ-0024/TQ-0036/TQ-0037;
TQ-0059 fixed four common semantic-response rejection defects. TQ-0061 changes
what should be implemented next: inactive network/registry code is to be retired,
not fully qualified as an MVP dependency. Canonical/semantic rejection protections
needed by the smaller product must survive the simplification.

qpdf 12.3.2 is installed from Fedora. This neither qualifies PDF isolation nor
means the pipeline works. TQ-0062 selects one supported bounded deployment with
actual basic enforcement evidence; TQ-0021 implements intake against it.

## Start here and keep going

```bash
git status --short --branch
git log --oneline -10
tq list --status in-progress --json
tq ready --label agent/developer --json
tq ready --label agent/architect --json
```

The live queue wins over this snapshot. Inspect any in-progress task's ownership
and notes before resuming it; preserve other agents' work. A developer selects
only developer tasks. An explicitly approved architect can also execute eligible
architect tasks, including TQ-0062, without requesting routine owner confirmation.

Ready first slices: TQ-0028 (first exact source recovery), TQ-0062 (minimal Fedora
runtime choice), TQ-0063 (retire platform prerequisites/local run records), and
TQ-0036 (fixed topics/document-scope states). TQ-0029 can recover the second exact
source independently at normal priority; it does not block the first programme.
Read each full task and use the actual executing model identifier when claiming.

No exact approved PDF/URL/checksum/election date may be invented. The first
source is the approved Frente Obrero **general political** document; the second
is the approved VOX Castilla y León document with evidenced electoral/territorial
identity. If the exact source is unavailable, record a concrete blocker. Candidate
recovery precedes registration, but does not grant content publication authority.

### Mandatory execution loop

1. Select approved ready work in priority/dependency order for your role. Read
   the **complete** record with `tq show <id> --json`; record the actual executing
   model with `tq update <id> --assignee "<actual-model-id>"`, then claim it with
   `tq move <id> in-progress` before editing.
2. Implement the task, including its real callers, fixtures, documentation and
   failure behavior. Record decisions and concrete evidence using `tq note`.
   A plan, stub, fake evaluation or passing unrelated test suite is not delivery.
3. Run **every required gate**. Fix failures and rerun affected checks plus
   regression. Keep expensive real-source/model evaluations explicit and retain
   their genuine outputs; never replace them with invented records or skipped passes.
4. Record all gate results with `tq note`, then `tq done <id>`. Inspect Git status,
   diffs and recent history; commit only intended task changes and its record.
   Push through the configured upstream and confirm the committed changes are
   present remotely, following the exact delivery procedure in AGENTS.md.
5. **Immediately re-query the live ready queue after confirmed delivery.** Newly
   unblocked tasks are part of the assignment. Read, assign and claim the next
   eligible task without asking Juan. Repeat the loop; do not stop after one
   task, the first source, the first summary or an arbitrary batch size.

### Blockers are records, not routine questions

- Verify a suspected blocker with the actual relevant command or evidence.
  Do not assume a package, provider or exact source is unavailable without
  checking, and do not assume every task requires both golden PDFs.
- For a task-local blocker that cannot be resolved within approved scope, leave
  the task open with a precise `tq note`: what was attempted, observed failure,
  affected gate and concrete unblock condition. Use native dependencies for
  known prerequisite tasks. Do not mark it done or repeatedly claim the same
  blocked task without new evidence.
- **Continue with eligible independent work instead of stopping to ask a
  question.** For example, inability to recover the exact Frente Obrero PDF
  leaves TQ-0028 open; it does not block TQ-0036 or TQ-0063. A developer waiting
  for architect TQ-0062 still executes other eligible developer work.
- Record an unresolved new architectural requirement as an inbox proposal under
  the role rules, rather than guessing its implementation or promoting it.
  Missing owner-supplied assets likewise remain a recorded dependency. Continue
  other approved work; no reply or routine approval prompt is required.
- Delivery failures are different: a failed commit, hook, signature, push or
  remote confirmation blocks further task claims. Restore the affected task to
  in-progress, record the failure, and resolve delivery using the existing
  configuration before continuing. Never disable signing/hooks, force-push,
  weaken a gate or modify somebody else's work to keep the loop moving.

### When to finish the execution session

Finish only when the live queue has **no eligible approved work remaining for
your role**, Juan explicitly stops execution, or a genuine environment/delivery
blocker prevents further safe progress. A blocker on one task is not grounds to
stop while independent eligible tasks remain.

At that point provide a concise factual report of delivered tasks, open blocked
tasks and their unblock conditions, and deferred/owner-only work. Do not claim
the entire project is complete merely because remaining tasks are blocked or in
inbox, and do not end with a question requesting permission to continue.

## Milestones

1. Recover one real candidate and origin evidence (TQ-0028).
2. Use one bounded runtime and simple local source/run records
   (TQ-0062/TQ-0021/TQ-0063/TQ-0027).
3. Extract original page evidence; generate and independently validate a cited
   Spanish draft (TQ-0030/TQ-0031/TQ-0036/TQ-0037/TQ-0024/TQ-0038–TQ-0041).
4. Render Markdown, verify current local eligibility, document limitations and
   prove one real end-to-end run (TQ-0043/TQ-0050/TQ-0051/TQ-0053/TQ-0064).
5. Repeat with the second source (TQ-0029/TQ-0054).
6. Integrate supplied static HTML, basic parity/a11y/mobile checks and separate
   local draft preview (TQ-0044/TQ-0046/TQ-0047/TQ-0049/TQ-0065).

Full native dependencies and seven gates are in the individual task records.
Model/converter evaluations must actually run in their owning tasks and retain
original evidence; ordinary `make verify` replays stored evidence offline.

## Deferred and superseded work

Under Juan's explicit replan, former accepted tasks TQ-0022/TQ-0023/TQ-0025/
TQ-0026/TQ-0033/TQ-0034/TQ-0045/TQ-0048 returned to inbox for future owner review.
TQ-0032/TQ-0035 were rejected as separate superseded tasks: necessary extraction
blockers/original-evidence checks remain in TQ-0030/TQ-0031/TQ-0039/TQ-0040.
No implementation is falsely marked done by this cancellation.

Pre-existing inbox proposals TQ-0017 (CI/deployment design) and TQ-0042
(human/non-LLM presentation assets) are untouched. Agents must not promote,
reject or execute them. HTML integration still needs supplied approved assets;
neither Markdown nor first-programme acceptance depends on HTML. No LLM may
author HTML/CSS/JavaScript/templates, fixtures or embedded presentation code.

## Authority and preservation

- The maintainer selects source bytes/evidence through a separate trusted record.
  Do not call this automatic official-origin network certification.
- Contributor source/report proposals cannot activate source selection or supply
  an accepted `pass`; the trusted local job creates actual content reports.
- Both current independent claim and whole-source coverage passes are required.
  Disagreement, insufficient evidence or stale inputs block public eligibility.
- Original source bytes and historical reports are preserved; changed input
  invalidates the affected programme, without a generic artifact graph.
- Private fixture/draft preview is explicitly separate from the public catalogue.
- Follow task-scoped gate/commit/push/confirmed-delivery rules. The previous
  planning task is complete; the resuming implementation agent now has Juan's
  explicit authorization to work through the approved MVP queue autonomously.
