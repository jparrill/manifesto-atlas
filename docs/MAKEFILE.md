# Makefile Technical Usage

This document is the technical reference for the `Makefile` orchestration
interface. It is English by project policy (technical text). It maps every
planned pipeline target to its owning backlog task and explains how the
deterministic, model, and browser checks will be integrated as those tasks land.

The current reduced product contract is [SPEC](SPEC.md) and
[MVP](design/MVP.md), established by TQ-0061. Existing code and targets are not
rewritten by the planning change; TQ-0063 owns retiring old platform code and
updating obsolete CLI/Make unavailable-stage ownership messages.

The product contract for the pipeline stages lives in
[`docs/SPEC.md#make-targets`](SPEC.md#make-targets) and
[`docs/SPEC.md#module-seams`](SPEC.md#module-seams). This document does not
change that contract; it only documents the current Make surface and the
expansion path.

## Current surface

The `Makefile` exposes implemented local checks/private-preview commands and
twelve declared production pipeline stages.

### Baseline targets (implemented)

| Target | What it does | Tools used |
| --- | --- | --- |
| `make help` | Print this target list. Distinguish implemented checks from planned stages. | none (print only) |
| `make check` | Run baseline checks: Git whitespace, Taskqueue integrity, and Go format/build/vet. | `git`, `tq`, `gofmt`, `go` |
| `make verify` | Run `check`, plus documentation-structure, relative-link, and Go test checks. | `git`, `tq`, `grep`, `dirname`, `test`, `go` |
| `make test-cli` | Run the focused Go CLI test suite (`go test ./internal/cli/...`). | `go` |
| `make test-artifacts` | Run the focused strict canonical-JSON and SHA-256 hash helper suite (`go test ./internal/artifact/...`). | `go` |
| `make test-taxonomy` | Run the focused `atlas.taxonomy/v1` registry, topic-state, and evidence-based document-scope classification suite (`go test ./internal/taxonomy/...`). | `go` |
| `make test-canonical` | Run the focused `atlas.programme.local/v1` local draft schema, citation, and source/extraction binding validation suite (`go test ./internal/canonical/...`). | `go` |
| `make test-semantic-runner` | Run common semantic-response rejection and fake-adapter regression tests (`go test ./internal/semantic/...`). | `go` |
 | `make test-fixture-frente` | Run the focused `atlas.candidate/v1`/`atlas.source.local/v1` Frente Obrero candidate hash/metadata fixture suite (`go test ./internal/candidate/...`). | `go` |
 | `make test-fixture-vox` | Run the focused `atlas.candidate/v1`/`atlas.source.local/v1` VOX Castilla y León electoral candidate hash/metadata fixture suite (`go test -run TestCheckVoxFixture ./internal/candidate/...`). | `go` |
  | `make test-runs` | Run the focused `atlas.selection/v1`/`atlas.run/v1` record validation suite (`go test ./internal/run/...`). | `go` |
  | `make smoke-provider [TARGET=...]` | Operator gate: run the three clean semantic roles (normalizer, claim-validator, coverage-validator) through the operator-configured OpenAI-compatible local provider and retain the actual request/model identities and original bounded replies (`go run ./cmd/smoke-provider`). Requires `ATLAS_SMOKE_BASE_URL`/`ATLAS_SMOKE_MODEL`; not part of `make verify`. | `go`, configured local provider |

`make` with no arguments runs the default goal (`help`) and performs no
network retrieval, model invocation, or publication. `make smoke-provider` and
`make e2e-preview`/`make evaluate-preview` are explicit operator evaluations that
invoke the configured model and are never run by `make verify` or `make check`.

### macOS private preview (implemented)

- `make test-preview`: offline source/citation, immutable supplied-asset binding,
  injection, stale-input and stored-real-run regression. On macOS it also exercises
  actual sandbox write/network denial, timeout/descendant cleanup and RSS/output
  bounds; other systems have no macOS worker and retain portable offline tests.
- `make render-preview [OUTPUT=...]`: check stored source/draft/report hashes and
  regenerate private Markdown/HTML without model calls. The default retained run
  is `runs/preview/frente-e2e`.
- `make e2e-preview OUTPUT=runs/preview/new-run [SOURCE=...]`: inspect/extract a
  real retained PDF, run three fresh configured-model roles and render a private
  preview. Requires `ATLAS_BASE_URL`, `ATLAS_MODEL`, and the imported owner-supplied
  reference; `ATLAS_API_KEY` is optional and never retained. The output directory
  must be new to preserve previous attempts. `make evaluate-preview` is an alias.
- `go run ./cmd/preview -check-browser -output runs/preview/frente-e2e`: actual
  desktop/mobile Chrome with CDP device emulation, screenshots and layout checks.

See the [macOS preview contract](design/macos-preview.md) for source selection,
reference import, provider setup, actual retained outcomes and limitations. This
private seam does not claim the production targets below are implemented.

### Baseline check semantics

- **Git whitespace** — `git diff --check` exits non-zero on a whitespace error
  (trailing whitespace, conflict markers, or a missing trailing newline) in the
  current diff. This is the whitespace gate that the disposable-copy failure test
  exercises.
- **Taskqueue integrity** — `tq list --json` must exit `0`, print non-empty JSON,
  and emit **no warnings on stderr**. `tq` names a malformed or skipped task file
  on stderr while still exiting `0`; a warning-only failure is therefore rejected
  as non-zero here. This is the malformed-task gate.
- **Documentation structure** — the required tracked files must exist:
  `AGENTS.md`, `README.md`, `HANDOFF.md`, `LICENSE`, `.taskqueue.yaml`,
  `docs/SPEC.md`, `docs/QUALITY_GATES.md`, and this file.
- **Relative links** — every non-URL, non-anchor Markdown link in the tracked
  docs must resolve to an existing file or directory relative to the doc that
  holds it.
- **Go format** — `gofmt -l cmd internal` must list no files; any unformatted Go
  source fails the check.
- **Go build** — `go build -o /dev/null ./cmd/manifesto-atlas` must compile the
  entry point (and its dependencies) without leaving a build artifact.
- **Go vet** — `go vet ./...` must report no static-analysis findings.
- **Go tests** — `make test-cli` runs `go test ./internal/cli/...`, `make
  test-artifacts` runs `go test ./internal/artifact/...`, `make
  test-taxonomy` runs `go test ./internal/taxonomy/...`, `make
  test-canonical` runs `go test ./internal/canonical/...`, `make
  test-semantic-runner` runs `go test ./internal/semantic/...`, `make
   test-fixture-frente` runs `go test ./internal/candidate/...`, `make
   test-fixture-vox` runs `go test -run TestCheckVoxFixture ./internal/candidate/...`,
    and `make test-runs` runs `go test ./internal/run/...`; all suites must pass.
    `make verify` includes them. `make smoke-provider` is deliberately excluded:
    it is an explicit operator gate that invokes the configured local provider
    (TQ-0024) and retains genuine request/model evidence, so it never runs under
    `make verify` or `make check`.

Passing these suites does not qualify every planned contract. The
[trust-boundary audit](audits/2026-11-02-trust-boundaries.md) records reopened
implementation tasks. TQ-0059 repairs strict common semantic envelopes, decision
consistency, tightened response-byte limits and late-response rejection. Empty
arrays must be explicit arrays, not null. Zero request limits select bounded
defaults (600 seconds, 8,192 tokens and 1 MiB raw reply); negative limits reject.
Stage-specific output schemas, evidence/coverage resolution, configured-provider
transport, streaming allocation limits and actual runtime/session isolation
remain required in TQ-0024. Adapters must honor cancellation; the common runner
rejects a late return but cannot terminate an arbitrary non-cooperative Go adapter.

## Go module

The repository is a Go module rooted at the checkout:

- Module path: `github.com/jparrill/manifesto-atlas`.
- Supported Go version: `go >= 1.26` (pinned in `go.mod`). The `version`
  command reports the exact runtime used to build.
- No external dependencies: the CLI uses only the standard library.
- Layout:
  - `cmd/manifesto-atlas/` — the entry point; dispatches to `internal/cli` and
    propagates its exit code.
  - `internal/cli/` — the command dispatcher and its tests.
  - `internal/artifact/` — the strict canonical-JSON serialization and SHA-256
    hash helpers (deterministic key order, integer/UTF-8/surrogate rules, depth
    and size bounds), reused by the selection, canonical, and semantic packages.
  - `internal/taxonomy/` — the `atlas.taxonomy/v1` topic registry, topic-state
    classification, and evidence-based document-scope validation logic, and its
    fixtures (TQ-0036).
  - `internal/canonical/` — the `atlas.programme.local/v1` local draft schema,
    source-record/extraction checksum binding, cited page/block/code-point
    excerpt evidence, retained scope-evidence verification, and local ordinal
    statement/measure/citation IDs, and its fixtures (TQ-0037).
  - `internal/semantic/` — common `atlas.semantic-job/v1` response verification,
    rejection limits, retry eligibility and terminal cancellation checks
    (TQ-0059), with offline fake fixtures; the full runner remains TQ-0024.
  - `internal/candidate/` — the `atlas.candidate/v1` recovery descriptor and
    embedded `atlas.source.local/v1` record, with offline hash/metadata
    verification of the retained Frente Obrero general-political source and its
    origin evidence (TQ-0028), and its fixtures.
  - `internal/run/` — the `atlas.selection/v1` owner selection record and the
    `atlas.run/v1` per-programme run record, with strict parsing, deterministic
    round-trip, and offline bounded-file hash verification (TQ-0063), and its
    fixtures.

The CLI supports `help` and `version`, and declares the 12 pipeline stages. A
stage requires exactly one target-path argument and reports not-available (exit
code 1) until its owning task lands; unknown commands, invalid flags, and wrong
argument counts exit with code 2. Argument values are treated as opaque data.
These behaviors are covered by the `test-cli` suite.

## Planned pipeline stages

The following is the **MVP ownership map**; the unavailable recipes and CLI
strings already name these owners (TQ-0063 removed the old platform code and
migrated the obsolete ownership messages). No target gains fake success during
replanning.

Each stage is declared now so a fresh checkout can discover the full interface,
but it fails non-zero with an actionable message that names its owning task.
When the owning task lands, its implementation replaces the not-available recipe
and the target becomes real. The mapping below is the source of truth for which
task owns which target.

| Target | Owning task(s) | Stage purpose (from `docs/SPEC.md#make-targets`) |
| --- | --- | --- |
| `make submit` | `TQ-0021` | Quarantine candidate bytes and submission metadata. |
| `make verify-source` | `TQ-0027` | Check owner-selected local source/evidence hashes; no automatic network certification. |
| `make add` | `TQ-0027` | Register a verified immutable source. |
| `make extract` | `TQ-0030` | Produce structured page-aware extraction. |
| `make validate-extraction` | `TQ-0030` | Reject deterministic page/warning blockers without a numeric score. |
| `make normalize` | `TQ-0038` | Produce draft canonical programme JSON. |
| `make validate-content` | `TQ-0041` | Run claim and coverage validators. |
| `make build` | `TQ-0043`, `TQ-0044`, `TQ-0049` | Markdown first; supplied-asset private HTML previews; public HTML only after basic checks. |
| `make validate-site` | `TQ-0049` | Combine basic single-programme parity, keyboard/accessibility and mobile checks. |
| `make publish-check` | `TQ-0050`, `TQ-0049` | Check local current Markdown evidence first, then independently add HTML eligibility. |
| `make process` | `TQ-0051` | Orchestrate the applicable processing stages for one programme. |
| `make serve` | `TQ-0065` | Preview an explicit private generated root on loopback without publication approval. |

Invoking any of these today prints the stage name, its current owning task
(from the map above), the purpose and a `tq show <id>` pointer, then exits
non-zero. Stages become real only in their implementation tasks.

## Integration path

The Go module and its compiler/test gates are in place. As the remaining
backlog lands, the same targets gain real behavior without changing their names
or the caller's expectations:

- **Go checks.** `make check` runs `gofmt`, `go build`, and `go vet`; `make
  verify` runs the real `go test` suite via `make test-cli`. As stage
  implementations land, their targets replace the not-available recipe with the
  owning Go command, and `make check`/`make verify` extend to the new packages.
  Deterministic work (conversion, checksums, schema validation, templates, HTML
  parsing, screenshot capture, deployment) stays in Go, per
  `docs/SPEC.md#module-seams`.
- **Model checks.** One actual provider serves normalizer, claim and complete-
  source coverage roles in separate clean requests. Explicit real evaluations
  retain genuine bounded replies; `make verify` replays stored evidence. No
  semantic provenance/extraction-score/visual-site approval role is an MVP gate.
- **Browser checks.** The supplied static-site extension checks parity, links,
  keyboard/accessibility and mobile reflow with existing tooling. Screenshots
  may be diagnostics, not a pixel-identical or model-approved publication gate.
  `serve` is separate read-only local preview and does not generate/approve data.

Until those tasks land, the not-available recipes are the truthful state: the
interface is discoverable and the absence is explicit rather than silent.
