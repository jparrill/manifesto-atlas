# Trust-boundary completion audit

Audit task: TQ-0058. Baseline commit: `37b7968`.
Scope: implementation and completion claims for TQ-0022, TQ-0023, TQ-0024,
TQ-0036 and TQ-0037. This is not a complete audit of the artifact/attestation
packages or a claim that the rest of the system is qualified.

## Decision

Reopen all five tasks as `todo`, preserving their implementation, history,
acceptance criteria, execution roles and native dependencies. Their existing
tests pass, but required outcomes and trust-boundary gates are incomplete.
The earlier all-PASS notes are superseded by corrective notes, not erased.
TQ-0025 must no longer be ready through incomplete prerequisites.

## Reproduction

Run from the repository root:

```bash
go run ./docs/audits/probe
```

The probe uses only synthetic in-memory records and an injected dialer that
records its argument and returns an error without connecting. No model, remote
host, political PDF, credential or sandbox is used. Its canonical baseline is
accepted by the current implementation; this does not endorse that baseline as
having complete verified jurisdiction/metadata evidence.

The command asserts the required *safe* behavior. At the audited baseline it
prints fourteen `FAIL` lines and exits nonzero (`exit status 1`). These are
reproduced implementation defects, not passing security tests. As repairs land,
the command should report fewer failures; current runtime gates must acquire
focused regression tests for the repaired behavior. The historical count below
is evidence about the baseline, not an instruction to retain vulnerabilities.

| Task | Required rejection invariant | Observed at baseline |
| --- | --- | --- |
| TQ-0022 | Reject HTTPS port 8443 | Accepted by `ValidateRequestURL` |
| TQ-0022 | Reject IPv6 documentation range `2001:db8::/32` | `IsPublicIP` returns true |
| TQ-0022 | Hand the approved IP to the actual dialer | Dialer receives the hostname instead; validated IP is discarded |
| TQ-0023 | Reject malformed publisher-key grants | Invalid key, state and dates pass `ValidateRegistry` |
| TQ-0036 | Reject forged competence evidence | Unapproved version, invalid date/hash and invented institution produce `not-applicable` |
| TQ-0036 | Reject unknown versions in comparison | Identical unapproved versions pass |
| TQ-0037 | Resolve a not-applicable rule against the profile | An invented rule is accepted |
| TQ-0037 | Bind extraction checksum to source | Wrong extraction checksum is accepted |
| TQ-0037 | Require the summary field | Missing field is accepted |
| TQ-0037 | Reject incomplete taxonomy without panic | Empty supplied topic inventory panics |
| TQ-0024 | Reject trailing JSON | A valid reply followed by another object is accepted |
| TQ-0024 | Reject pass with a blocker | Contradictory pass report is accepted |
| TQ-0024 | Enforce a tightened response-byte limit | A one-byte request limit is ignored |
| TQ-0024 | Discard success after the deadline | Adapter waits for cancellation, then returns an accepted pass |

## Missing deliverables found by inspection

These findings are distinct from the executable counterexamples above.

### TQ-0022: bounded transport

Contract: [provenance URL/network policy](../design/provenance-trust.md#url-and-bounded-network-policy).

- No production pinned-IP dialer is provided; any injected dialer can re-resolve
  the hostname after validation. Validation alone is not rebinding prevention.
- The API accepts arbitrary TLS configuration, including disabled certificate
  verification; current local TLS tests use that override.
- Special-purpose IP coverage, path/query ambiguity rejection, request headers,
  content-encoding policy and stage-specific deadlines/budgets are incomplete.
- MIME acceptance is substring matching, not media-type validation. The PDF
  default is 1 GiB rather than the isolation contract's 100 MiB default.
- Evidence lacks required hop status/Location/approved-IP and response fields.
  A supplied referral-page URL is merely echoed; no fetched page body or exact
  parsed link is established. Atomic quarantine installation is absent.

### TQ-0023: authority registry

Contract: [party registry and authority](../design/provenance-trust.md#party-registry-and-anchor-authority).

- The public mutable registry structure has no trusted update/admission check.
  Passing caller-selected data is a seam, not proof that unauthorized updates
  are rejected. Existing tests do not attempt contributor self-activation.
- Evidence is checked only for hash syntax, not existence, byte integrity or
  independent authority-chain support. An arbitrary hash is not authority.
- Publisher-key validation, territorial/document grant structure, and
  provenance-report bindings to a selected trusted registry are incomplete.
- `AnchorChanged` compares structures; it does not itself implement approval
  graph invalidation or authority verification.

### TQ-0024: semantic runner

Contract: [provider-independent semantic jobs](../design/semantic-validation.md).

- Only a test fake exists; there is no configured local-provider adapter or
  documented executable smoke command.
- `OutputSchema` is unused. Stage-specific schema enforcement is not provided
  merely by decoding common Go structs.
- `EvidenceReader` is forwarded without evidence allowlisting, byte budgets,
  inventory verification or checked-target resolution.
- Prompt/model/input invocation identities are not recorded as required.
  Whole-evaluation limits and local isolation qualification are absent.
- The canary test creates separate runners and separate fakes. It does not
  exercise provider state reuse. The substitution test instantiates the same
  fake type twice with the same identity; it proves no configured adapter.
- A Go interface without tool methods is not OS capability confinement or
  provider-session/KV-cache isolation.

### TQ-0036: jurisdiction-aware taxonomy

Contract: [competence profiles](../design/taxonomy-jurisdiction.md#competence-evidence-and-scope-profiles).

- The ordered registry and labels exist, but dated profile evidence and
  approved institution/rule inventories are not supplied or verified.
- Nonempty strings stand in for profile version, legal date and content hash.
  A whole-topic `outside` flag does not prove all constituent subtopics are
  outside scope.
- Comparisons omit document type, election cycle/type, jurisdiction and
  geography compatibility required by the design.
- Migration tests exercise only rejection: v1 has no approved successor. Do
  not describe those tests as proof of a successful split/merge migration.

### TQ-0037: canonical programme

Contract: [canonical schema and citation contract](../design/canonical-programme.md).

- No machine-readable schema is supplied. Decoding into structs rejects extra
  fields but does not enforce all required members, nullability or exact names.
- Metadata evidence values are only checked for nonemptiness, not resolved to
  the specified trusted source/registry fields.
- Separate membership maps do not express or verify institution/rule/topic/
  relation compatibility or entire-topic applicability.
- Measure IDs are checked for shape rather than anchor/ordinal derivation;
  canonical citation/measure ordering and statement ordinals are incomplete.
- The supplied taxonomy inventory can differ from the immutable registry.

## Corrections to the reported blockers

**TQ-0021:** `qpdf --version` reports 12.3.2, installed from Fedora. Missing host
qpdf is no longer a blocker; it also never proved a container image could not
be built. Juan requested avoiding unnecessary rootless setup. That direction
needs an explicit architecture update of the current
[isolation contract](../design/submission-isolation.md), preserving demonstrable
resource limits, teardown, repository protection and network restrictions.
Installing qpdf alone does not establish these properties. No untrusted PDF
was parsed on the host during this audit.

**TQ-0025:** the task does not require golden political PDFs for prompt
evaluation. A real configured-model invocation using explicitly synthetic
direct/referral/signature/missing-link/injection fixtures is legitimate.
The agent can execute it when a configured, identified provider is available;
Juan need not execute it personally. Provider availability has not been
established either way. The current concrete blockers are the reopened native
prerequisites TQ-0022, TQ-0023 and TQ-0024. The eventual evaluation must record
actual execution, exact inputs/prompt/schema/model identities, and original
outputs; fake-adapter replies cannot substitute for it.

## Recovery order

1. Repair the semantic response rejection boundary immediately: strict JSON,
   tightened response limits, contradictory pass findings and late replies.
   Keep TQ-0024 open until its other required deliverables also pass.
2. Complete pinned transport and trusted registry admission/evidence resolution
   under TQ-0022/TQ-0023; add real controlled-network and self-approval negatives.
3. Complete the configured provider adapter and evidence/context boundaries in
   TQ-0024. Then author/evaluate the provenance skill in TQ-0025 against that
   actual runner and synthetic evidence.
4. Complete dated jurisdiction profiles in TQ-0036, then the dependent strict
   canonical schema/reference validator in TQ-0037.
5. Resolve the Fedora isolation profile explicitly and implement/test TQ-0021.

The probe is an audit aid, not an exhaustive qualification suite. An individual
repair may fix a listed rejection invariant without making the owning task
complete. All seven original gates still apply.

## Response-boundary repair: TQ-0059

The four semantic rejection probes now print `PASS`; the ten unrelated probes
still print `FAIL`, so the aggregate audit command still exits one. Its synthetic
positive semantic baseline must succeed before these rejection checks run.

The runner now uses the existing strict artifact parser for the entire raw
response and checks exact required common fields before typed decoding. It
rejects contradictory passing findings, honors tightened byte limits, supplies
bounded defaults, rejects invalid limits, and checks cancellation/deadline after
adapter returns and before returning a report. Cancelled transient failures and
transient errors accompanied by partial bytes do not retry. No stage-specific
schema, provider capability or coverage completeness is implied by these checks.

Regression evidence lives in `internal/semantic/response_boundary_test.go`, run
by `make test-semantic-runner`. TQ-0024 remains open for the other deficiencies
above. An adapter that never returns must be terminated by its actual transport
or isolated process implementation; late-return rejection alone is not teardown.
