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

The `Makefile` exposes twelve implemented baseline targets and twelve declared
pipeline stages.

### Baseline targets (implemented)

Some currently installed targets exercise historical graph/attestation/network
libraries. They remain truthful library checks until TQ-0063 removes that inactive
scope; their presence is not a requirement to finish the former platform.

| Target | What it does | Tools used |
| --- | --- | --- |
| `make help` | Print this target list. Distinguish implemented checks from planned stages. | none (print only) |
| `make check` | Run baseline checks: Git whitespace, Taskqueue integrity, and Go format/build/vet. | `git`, `tq`, `gofmt`, `go` |
| `make verify` | Run `check`, plus documentation-structure, relative-link, and Go test checks. | `git`, `tq`, `grep`, `dirname`, `test`, `go` |
| `make test-cli` | Run the focused Go CLI test suite (`go test ./internal/cli/...`). | `go` |
| `make test-artifacts` | Run the focused `atlas.artifact/v1` envelope/graph test suite (`go test ./internal/artifact/...`). | `go` |
| `make test-attestations` | Run the focused `atlas.attestation/v1`/`atlas.trust/v1` signing/verification suite (`go test ./internal/attestation/...`). | `go` |
| `make test-taxonomy` | Run the focused `atlas.taxonomy/v1` registry, topic-state, and jurisdiction-validation suite (`go test ./internal/taxonomy/...`). | `go` |
| `make test-canonical` | Run the focused `atlas.programme/v1` canonical schema, citation, and source-binding validation suite (`go test ./internal/canonical/...`). | `go` |
| `make test-party-registry` | Run the focused `atlas.provenance/v1` party trust-anchor and registry-validation suite (`go test ./internal/provenance/...`). | `go` |
| `make test-fetch` | Run the focused `atlas.provenance/v1` bounded safe HTTP transport (URL/SSRF/redirect/size) suite (`go test ./internal/provenance/...`). | `go` |
| `make test-semantic-runner` | Run common semantic-response rejection and fake-adapter regression tests (`go test ./internal/semantic/...`). | `go` |
| `make test-fixture-frente` | Run the focused `atlas.candidate/v1`/`atlas.source.local/v1` Frente Obrero candidate hash/metadata fixture suite (`go test ./internal/candidate/...`). | `go` |

`make` with no arguments runs the default goal (`help`) and performs no
network retrieval, model invocation, or publication.

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
  test-attestations` runs `go test ./internal/attestation/...`, `make
  test-taxonomy` runs `go test ./internal/taxonomy/...`, `make
  test-canonical` runs `go test ./internal/canonical/...`, `make
  test-party-registry` and `make test-fetch` both run
   `go test ./internal/provenance/...` (the party registry and the bounded
   transport share the package), `make test-semantic-runner` runs
   `go test ./internal/semantic/...`, and `make test-fixture-frente` runs
   `go test ./internal/candidate/...`; all suites must pass. `make verify`
   includes them.

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
  - `internal/artifact/` — the `atlas.artifact/v1` canonicalization, hashing,
    manifest/graph validation, and invalidation logic, and its fixtures
    (TQ-0019).
  - `internal/attestation/` — the `atlas.attestation/v1`/`atlas.trust/v1`
    Ed25519 signing and offline verification logic, and its fixtures (TQ-0020).
  - `internal/taxonomy/` — the `atlas.taxonomy/v1` topic registry, topic-state
    classification, jurisdiction profiles, comparison compatibility, and
    migration validation logic, and its fixtures (TQ-0036).
  - `internal/canonical/` — the `atlas.programme/v1` canonical schema, topic/
    state rules, citation/page/excerpt evidence, source binding, and
    deterministic-serialization validation logic, and its fixtures (TQ-0037).
  - `internal/provenance/` — the `atlas.provenance/v1` party/territory/domain
    evidence registry, fail-closed anchor authority checks, version
    invalidation (TQ-0023), and the bounded safe HTTP transport with URL/SSRF/
    redirect/size policy (TQ-0022), and their fixture anchors and local TLS
    fixtures.
  - `internal/semantic/` — common `atlas.semantic-job/v1` response verification,
    rejection limits, retry eligibility and terminal cancellation checks
    (TQ-0059), with offline fake fixtures; the full runner remains TQ-0024.
  - `internal/candidate/` — the `atlas.candidate/v1` recovery descriptor and
    embedded `atlas.source.local/v1` record, with offline hash/metadata
    verification of the retained Frente Obrero general-political source and its
    origin evidence (TQ-0028), and its fixtures.

The CLI supports `help` and `version`, and declares the 12 pipeline stages. A
stage requires exactly one target-path argument and reports not-available (exit
code 1) until its owning task lands; unknown commands, invalid flags, and wrong
argument counts exit with code 2. Argument values are treated as opaque data.
These behaviors are covered by the `test-cli` suite.

## Planned pipeline stages

The following is the **new MVP ownership map**, not a claim that the current
unavailable recipes/CLI strings have already been migrated. Those still name
legacy owners until TQ-0063. No target gains fake success during replanning.

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

Invoking any of these today prints the stage name, a historical owning task,
the purpose and a `tq show <id>` pointer, then exits non-zero. After TQ-0063,
unavailable messages must identify the current owner above. Stages become real
only in their implementation tasks.

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
