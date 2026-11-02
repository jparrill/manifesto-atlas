# Manifesto Atlas Handoff

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

The repository contains a Go CLI, tested domain libraries, Make gates and design/
task documentation. The ordinary executable still provides help/version and
unavailable pipeline-stage messages. No source PDF, project skill, configured
provider integration, public output, supplied presentation asset or CI workflow
is committed. Passing library tests are not end-to-end product readiness.

The previous technical audit reopened TQ-0022/TQ-0023/TQ-0024/TQ-0036/TQ-0037;
TQ-0059 fixed four common semantic-response rejection defects. TQ-0061 changes
what should be implemented next: inactive network/registry code is to be retired,
not fully qualified as an MVP dependency. Canonical/semantic rejection protections
needed by the smaller product must survive the simplification.

qpdf 12.3.2 is installed from Fedora. This neither qualifies PDF isolation nor
means the pipeline works. TQ-0062 selects one supported bounded deployment with
actual basic enforcement evidence; TQ-0021 implements intake against it.

## Execution frontier after this planning delivery

```bash
tq ready --label agent/developer --json
tq ready --label agent/architect --json
git status --short --branch
git log --oneline -10
```

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
- Follow task-scoped gate/commit/push/confirmed-delivery rules. A planning-only
  assignment does not automatically expand into implementation.
