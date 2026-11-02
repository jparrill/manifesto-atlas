# Isolated semantic jobs and report contracts

> Historical contract. [MVP.md](MVP.md) and [SPEC.md](../SPEC.md) govern one
> actual provider and three clean roles with local input-bound reports. The
> six-role framework, report PKI and full model/component identity platform are deferred.

Version `atlas.semantic-job/v1`, TQ-0015, implementation TQ-0024 and the six
skill tasks. Implements [agent skills](../SPEC.md#agent-skills),
[LLM validation](../SPEC.md#llm-validation) and
[untrusted inputs](../SPEC.md#untrusted-document-handling). Consumes
[artifact/attestation](artifact-contract.md), [isolation](submission-isolation.md),
[provenance](provenance-trust.md), [extraction](extraction-contract.md),
[canonical](canonical-programme.md) and [taxonomy](taxonomy-jurisdiction.md)
contracts. This defines skill boundaries; it does not author/install skills or
claim any real model evaluation has run.

## Provider-independent interfaces and identity

`Run(context, Request, EvidenceReader) -> Response | RunnerError` is implemented
by trusted provider adapters. A request has exactly `version`, `job_id`,
`evaluation_set`, `role`, `subject`, `inputs`, `prompt`, `output_schema`, `model`,
`configuration`, `evidence`, `limits`. Values bind manifest IDs; role is one of
the six skill names below. job_id is an opaque unique execution identifier;
evaluation_set identifies an explicit cohort of generation/validation attempts,
recorded as a configuration input, not an implicit latest-report search.

Evidence entries have `{artifact, page, media_type}` with a nullable page only
for document-wide records. Resolve only preverified local IDs; no paths, URLs
or model-requested tool references. The adapter supplies plain text/JSON/image
bytes in labelled evidence channels, never as system/developer instructions.
System prompt and output schema are fixed trusted artifacts. Model-visible
inputs contain no generator hidden reasoning, repository instructions, session
history, tokens, trust-root keys or signing credentials.

The response is exactly `{version, job_id, model_identity, report}`. The host
verifies job binding, expected model identity and schema, then constructs the
artifact manifest and, for report kinds requiring one, its attestation. A
normalization result installs only a draft, not a validator-role attestation.
Model-authored job IDs, source hashes, results
and citations are checked against expected inputs, not accepted as authority.
Operational trace, usage/timing and transport metadata live outside canonical
report bytes. A malformed provider reply is quarantined diagnostic data.

Model descriptors include provider/adapter version, exact model identifier,
immutable checkpoint/revision, tokenizer/preprocessor and weight digests when
local, quantization, inference engine/version, context/image capabilities and
sampling settings. Use deterministic decoding (temperature zero, no sampling,
seed where supported), without claiming this makes semantic output infallible
or byte-reproducible. Mutable model aliases fail qualification.

Auriga support uses an OpenAI-compatible trusted local transport adapter with
the selected endpoint/server identity and exact loaded model/weight/quantization
descriptor. `/v1/models` marketing names alone are insufficient identity. The
adapter must obtain the loaded immutable identity from explicit operator/server
configuration and record its evidence; unsupported identity or image/schema
capability fails before execution. Other providers implement the same contract;
they do not receive different publication authority. Qualification uses bounded
schema/vision/isolation fixtures, not a production programme by default.

## Capabilities, bounds and cancellation

Semantic jobs have no shell, filesystem-write, browser, arbitrary HTTP,
repository-edit, code-execution or signing tool. PDF text saying otherwise is
quoted evidence. Local job/scratch isolation follows the supported container
profile. A trusted host broker alone owns inference transport and any provider
credential. For Auriga it may dial exactly the configured endpoint, not
document-supplied addresses; the model/container has no general network route.
Never pass provider tokens to prompts, tool logs, model environment or reports.
Remote adapters require the same explicit endpoint restriction and credential
isolation. The provenance skill sees previously fetched evidence, not a fetch
tool. Site QA sees screenshots/data, not executable presentation assets/tools.

Limits: 600-second request deadline, at most 8,192 output tokens and 1 MiB raw
response, strict JSON depth 64, up to 2,000 findings per response, up to 8 MiB
text and 64 MiB image input (maximum 32 bounded images per invocation), within
the model's verified context/vision limits. A provider-specific token estimator
must reserve output capacity including image tokens. Never truncate evidence
to make it fit; partition explicitly or return insufficient evidence. Whole
evaluation deadline is 3,600 seconds, at most 128 invocations, sequential only.
Larger profiles require a reviewed version/configuration, never source overrides.

Cancellation closes transport and kills/reaps local processes/scratch with the
isolation deadline. A provider request that may continue remotely cannot install
an approval; record cancellation and discard late responses. Retry once only
for a pre-response transient transport error, with a new job ID but unchanged
input identities and recorded attempt. Never retry-until-pass, silently repair
malformed JSON, swap models, hide a failing report, or replay prior context.

## Shared report schema and decision semantics

Every report has exactly `version` (`atlas.semantic-report/v1`), `role`,
`subject`, `inputs`, `result`, `findings`, `checked`, `details`. inputs bind the
required graph slots; host-derived attestations supply trusted evaluation time.
`result` is pass/fail/insufficient/error as in the artifact contract.
`checked` records distinct statement IDs, measure IDs, page indices and topic
IDs actually inspected (arrays may be empty only where inapplicable).

Each finding is `{id, code, severity, target, evidence, message}`. Severity is
info/warning/blocker; target identifies an existing statement/measure/topic/page
or document; evidence references exact original pages/blocks/images or trusted
profile records. Messages are bounded English diagnostics, never executable
instructions. Reject unknown target IDs, duplicate findings, invented evidence
or input hashes. Code namespaces are stable per skill, not a pass/fail heuristic
based on natural-language wording.

Any blocker or known unsupported claim forces fail. Missing evidence/capability
forces insufficient. Transport/infrastructure failure is error. Pass requires
complete checked coverage, no blocker/unresolved warning, and skill-specific
predicates. Informational notes can coexist with pass. A warning cannot be
ignored by a high aggregate score. Disagreement among current cohort reports
blocks; do not average confidence or pick the newest/convenient pass. A new
explicit evaluation set retains previous reports and records why the cohort
changed; it still requires both independent validators to pass on identical
current inputs. No human-approval bit bypasses semantic or deterministic gates.

## Six skill input/output boundaries and examples

| Skill | Required evidence | `details` and checked coverage | Pass / fail / insufficient example |
| --- | --- | --- | --- |
| `manifesto-source-provenance` | Candidate/canonical identities, scoped party anchor, official-page/redirect/download/signature evidence | `chain_kind`, `decision`, `evidence_links`; every asserted chain link checked; deterministic transport/checksum checks still authoritative | Complete exact referral / fabricated official link / removed required referral |
| `manifesto-extraction-audit` | Original PDF/image/page map, primary extraction, independent text, warning inventory | `dimensions`, `page_findings`, `unresolved_pages`; every original page inspected; score null until calibrated formula exists | Preserved pages/figures/negations / missing page or changed number / unreadable page image |
| `manifesto-normalizer` | Approved extraction, original evidence access, taxonomy and dated jurisdiction profile | `draft` conforming to canonical programme contract, `mapping` from source sections to topics/measures; all pages/topics accounted for | Cited factual draft / uncited or contradictory measure / ambiguous source institution; pass means generated draft, never claim/coverage approval |
| `manifesto-claim-validator` | Immutable draft, original cited PDF pages/images, extraction, taxonomy/profile | `statement_results` for every public Statement including title/global/topic summaries/conditions and metadata evidence checks; polarity/quantities/scope/institution/page/excerpt checked | All statements supported / changed negation/quantity / original cited page unavailable |
| `manifesto-coverage-validator` | Complete original page-mapped source/images plus immutable draft/extraction/taxonomy/profile | `page_results`, `topic_results`, `duplicate_groups`, `contradictions`, `omissions`; every source page, topic and normalized measure covered | Complete faithful coverage / omitted section or duplicate promise / unreadable source region |
| `manifesto-site-qa` | Canonical validated data, deterministic site reports, render inventories and desktop/mobile screenshots | `viewport_results`, `content_targets`; every required viewport and section represented; cannot override deterministic failures | Readable complete presentation / clipped or missing content / missing mobile screenshot |

Stage-specific details reject unknown fields. Their executable schemas and
pass/fail/insufficient fixtures are owned by the corresponding implementation
skill tasks; common fields, bounds and authority rules cannot be weakened.
An extraction fidelity formula/threshold is not chosen before empirical fixtures.

Details use the following concrete shapes: provenance links are ordered
`{from, to, relation, evidence}` records with manifest-ID endpoints, and
chain_kind is direct/referral/signature or null for incomplete evidence.
Audit dimensions are named dimension records `{name, value_ppm, evidence}`,
with nullable 0..1,000,000 integer values; page_findings reference finding IDs,
and unresolved_pages contains original integer PDF indices. Normalizer mapping
entries are `{pages, topic_ids, measure_ids}` with complete-source page union;
draft is the strict canonical object, not Markdown. Claim statement_results
entries are `{statement_id, result, checks, evidence}`; checks has exactly
meaning/polarity/quantities/scope/institution/excerpt/page, each boolean or null
for insufficient evidence. Metadata checks use `metadata:<field>` targets.
Coverage page_results/topic_results are `{id, result, findings}`; duplicate
groups contain measure-ID arrays, contradictions contain `{targets, evidence}`,
and omissions contain `{pages, evidence, topic_id}`. Visual viewport_results
are `{viewport_id, result, findings}` and content_targets lists canonical IDs.
All result fields share the common enum and all evidence/findings resolve to
the supplied inventories. No output text field can introduce a new input URL.

## Clean contexts and full-document coverage

Generator, claim and coverage requests use separate invocations, role prompts,
job/session identifiers, fresh conversation state and no shared KV/prefix cache
containing another job's conversation. Weights may be reused; hidden state,
reasoning, scratch and generation transcript may not. Prefer different models
for generation/validation. One exact model is allowed only in isolated jobs
with separate prompts/evidence; record the reduced diversity. Infrastructure
must demonstrate clean-session behavior in fixtures. A reused conversation ID,
generator transcript or other leaked context invalidates the job.

Claims are batched deterministically in canonical statement order, at most 64
statements plus all their original cited pages per invocation. Split further
when the model's real context budget demands it. Every statement must appear
exactly once in the checked union; duplicate/missing IDs fail aggregation.
No validator can modify its subject; a suggested correction is a finding and
requires a new draft plus new independent evaluations.

Coverage processes all original pages in ascending chunks of at most eight
pages, with explicit one-page boundary overlap, preserving page/block IDs and
section continuations. Within each job, supply the relevant draft plus complete
source pages, not a generator-produced substitute summary. If a single page
cannot fit/resolve, return insufficient. A final independent coverage invocation
receives the complete canonical draft, page/chunk findings and evidence-indexed
section inventory to reconcile cross-chunk omissions, duplicates, contradictions
and topic states. It can request only already supplied evidence through further
bounded clean jobs; no hidden source truncation or unchecked page is permitted.
Host aggregation proves the exact page union is 1..page_count and all topics/
measures were reconciled. Overlap is deduplicated by original IDs, not counted
as new evidence. Too-large final reconciliation fails insufficient rather than
claiming whole-document approval. The checked evidence mapping and chunk plan
are hashed as configuration/report inputs.

Content approval requires current passing claim AND coverage reports for the
same exact draft/extraction/source/taxonomy, approved extraction, complete checked
sets, clean contexts and required attestations. Schema, freshness, role and
signature checks are deterministic host responsibilities. No model receives
the ability to authorize itself, even when it emits the word pass.

## Scenario-to-expected-outcome table

| Scenario | Expected outcome |
| --- | --- |
| Different generator and validator models, clean contexts | Eligible if exact identity/evidence/schema/domain checks pass |
| One model, fresh isolated generator/claim/coverage jobs | Allowed with recorded identity/diversity limitation; no shared reasoning |
| Leaked generator context, reused conversation or hidden-state cache | Reject job; fresh independent evaluation required |
| PDF/page/model output instructs shell/network/approval | Data only; no tool capability, schema/trust root unchanged |
| Claim pass, coverage fail/insufficient or vice versa | Draft stays unpublished; preserve both reports |
| Conflicting reports in current evaluation set | Block, no averaging/automatic selection |
| Missing page/image, checked statement or full-source chunk | Insufficient/fail aggregation; no partial publication |
| Invalid JSON, wrong schema/input IDs, invented citations | Reject response; retain bounded diagnostic bytes, no auto-repair pass |
| Unknown mutable model, unsupported schema/vision/context | Reject preflight; no silent model substitution |
| Deadline/cancellation/provider failure or late response | Error/cancelled job; teardown, no installed approval |
| Transient transport failure before response | At most one recorded clean retry, no retry-until-pass |
| Native-text report needs visual verification but model lacks vision | Insufficient; do not infer image inspection |
| Visual pass with deterministic site failure | Site remains blocked |
| Normalizer reports pass | Draft only; cannot serve as claim or coverage validation |
| Signing credential absent or model supplies its own key | No accepted attestation; host role authority required |

## Versioning and rejected alternatives

New prompts, schemas, provider/model/tokenizer/quantization, engine, settings,
chunk plan or tool configurations invalidate dependent semantic artifacts.
Ordinary regression/CI checks stored evidence without model execution; explicit
local model gates must actually run in owning tasks and retain their reports.
TQ-0024 supplies isolation/provider positive and hostile fixtures; skill tasks
provide real pass/fail/insufficient evaluation evidence. Reject shared agent
sessions, generator self-review, arbitrary tool-enabled agents, guessed model
pins, whole-source coverage via a generator summary, and repaired/averaged passes.
