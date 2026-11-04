# Local-first MVP contract

Current contract: `atlas.mvp/v1`, established by TQ-0061 under Juan's explicit
scope-reduction and replanning request. Read [SPEC](../SPEC.md) first.
This file supersedes conflicting requirements in the nine older design documents
for active MVP tasks. Historical artifacts/tasks are not retrospectively certified
or erased. Runtime changes occur only in their owning implementation tasks.

For the owner's immediate macOS E2E/HTML demonstration, the narrower
[macOS preview contract](macos-preview.md) governs TQ-0062/TQ-0067. It uses installed
Poppler and sandbox-exec plus the supplied immutable HTML reference. Fedora and
the production Docling path are not prerequisites to this explicitly private
preview; incomplete semantic checks remain visible and non-publishable.

## Milestones and native task map

| Milestone | Tasks |
| --- | --- |
| Recover first source without platform prerequisites | TQ-0028 |
| Choose one bounded Fedora converter deployment | TQ-0062 (architect), then TQ-0021 |
| Retire unused platform code and supply simple run records | TQ-0063 |
| Register/extract original and page evidence | TQ-0027, TQ-0030, TQ-0031 |
| Fixed topics, canonical schema and one real provider | TQ-0036, TQ-0037, TQ-0024 |
| Generate and independently check factual content | TQ-0038, TQ-0039, TQ-0040, TQ-0041 |
| Produce eligible Markdown and one-programme command | TQ-0043, TQ-0050, TQ-0051, TQ-0053 |
| First real end-to-end acceptance | TQ-0064 |
| Recover and process second source | TQ-0029, TQ-0054 |
| Supplied static HTML, basic quality and private preview | TQ-0044, TQ-0046, TQ-0047, TQ-0049, TQ-0065 |
| Document maintainer workflow / proposal-only contributions | TQ-0052 |

TQ-0042 remains the pre-existing owner-only inbox proposal for human/non-LLM
presentation assets; it is not promoted by this plan. TQ-0044 is legitimately
blocked until those assets are supplied/approved. Markdown does not depend on
that boundary. TQ-0017 remains owner-only CI/deployment design, not authorization
to implement a workflow. No active Markdown task depends on deferred proposals.

TQ-0022/TQ-0023/TQ-0025/TQ-0026/TQ-0033/TQ-0034/TQ-0045/TQ-0048 are deferred to
inbox, preserving their old future contracts and notes. TQ-0032/TQ-0035 are
superseded by the smaller extraction/evidence/content checks, not marked done.

## Authority and source selection

The maintainer supplies an owner-controlled selection record in the trusted
checkout. Its exact fields are `version` (`atlas.selection/v1`) and `sources`;
each source entry is `{programme_id, source_checksum, metadata_checksum,
evidence_checksums}`. Hashes are full lowercase `sha256:<64hex>` values. Entries
are unique and bounded. Only the trusted operator can select a new record or
update it; a model, PDF or contributor path cannot activate a selection.

Recovery tasks obtain and retain candidate bytes, factual metadata and actual
official-origin evidence **before** registration exists. Use existing acquisition
tools in explicit maintainer-authorized work; do not implement a general fetcher
or claim automatic network certification. No parser runs against a recovered PDF
until the bounded deployment exists. Candidate recovery can hash/copy bytes and
record externally established evidence without asserting a page count.

The recovery descriptor (`fixtures/frente/source.json` or `fixtures/vox/source.json`)
has exactly `version` (`atlas.candidate/v1`), `programme_id`, `candidate_path` and
`metadata`. Candidate path resolves only within the declared input root. Metadata
is the complete local source record below, including the deterministic eventual
registered path, not a claim that registration already ran. The trusted selection
binds its exact metadata bytes and candidate/evidence hashes. Registration verifies
that the expected destination equals the generated path, copies from the opened
candidate handle, and installs those selected metadata bytes without rewriting
them or manufacturing fresh retrieval timestamps.
Use the retained strict canonical JSON helper to serialize embedded metadata for
both `metadata_checksum` and `metadata.json`; candidate-file indentation does not
change that selected metadata identity. Raw PDF/evidence bytes are never normalized.

## Local source record

Version: `atlas.source.local/v1`. The `version` field and all fields below are required:

- `programme_id`: local lowercase ASCII slug, independent of a display name.
- `party`: `{id, name}`; `document_type`: electoral-programme/general-political.
- `scope`: european/national/autonomous/provincial/municipal/general-political.
- `geography`: `{id, name}`; `election`: null or `{id, type, date}`.
- `publication_date`: ISO date or null, never guessed from retrieval/upload.
- `archive_year`, `year_basis` (publication/retrieval).
- `source`: `{path, sha256, size, original_filename, origin_url, retrieved_at}`.
- `origin`: `{authority, evidence}`; authority is exactly `maintainer-selected`;
  evidence is a nonempty array of `{path, sha256}` retained evidence references.
- `rights`: `{publisher, statement}` with explicit null for unknown statement.
- `supersedes`: previous registered source checksum or null.

All objects reject unknown/missing fields, invalid enums/dates/hashes and unsafe
paths. Electoral records need an evidenced election/date; general political
documents have election=null. Public source text is escaped data. Display names
come from the selected source record, not model embellishment.

Storage: `sources/<archive-year>/<programme-id>/<64hex-checksum>/source.pdf`
and `metadata.json`. Output path segments are validated/generated, never derived
from original filenames. Exact repeated registration is idempotent; changed bytes
require a new revision. Do not overwrite committed source bytes or historical
metadata. Corrections append new records, and current selection is explicit.
No distributed/concurrent store or content-addressed manifest layer is required.

## Minimal extraction and evidence

Version: `atlas.extraction.local/v1`; exact common fields are `version`,
`programme_id`, `source_checksum`, `page_count`, `pages`, `blocks`, `warnings`,
`tool`, `configuration`. `tool` and `configuration` identify actual recorded
version/configuration files by `{path, sha256}`; native raw tool output is retained
separately. `pages` is a contiguous array of `{page, handling, image}` records,
where page is the one-based original PDF index, handling is native/visual-required/
unresolved, and image is null or `{path, sha256}`. `blocks` contains exactly
`{id, page, order, type, text}` per block. IDs are local and unique; order is
positive and deterministic within each page. Native table/layout data remains in
the raw converter evidence rather than a speculative second layout schema.

Warnings are `{id, page, block_id, code, severity, message}` with nullable page/
block for document-wide warnings. Severity is info/warning/blocker. Missing pages,
invalid mapping, material unresolved warnings and unrenderable required regions
block public eligibility. A deterministic extraction check is not proof of semantic
fidelity; original-page claim/coverage checks remain mandatory.

Page images use a bounded existing renderer in the same supported deployment.
Every requested original page must resolve to the same PDF checksum/index.
Validators inspect original page/image evidence, not only generator text. If a
required image/visual capability is unavailable, return insufficient. OCR can be
requested explicitly later; the initial native path must not silently invent
text for scanned/unreadable pages. No numeric fidelity score is computed.

## Canonical draft

Version: `atlas.programme.local/v1`. Preserve the useful shapes from the former
canonical contract, with these deliberately smaller identity/binding rules:

- Top-level exact fields: `version`, `programme_id`, `source_checksum`,
  `source_metadata_checksum`, `extraction_checksum`, `taxonomy_version`, `metadata`,
  `summary`, `topics`, `measures`, `citations`.
- `metadata` is exactly `{party, document_type, scope, geography, election,
  publication_date, rights, document_title}`. All except document_title must equal
  the selected source record. document_title is a cited Statement or null.
- Statement: `{id, text, citations}`. IDs `s-1`, `s-2`, ... are unique in document,
  topic, measure traversal order. Text is plain Spanish; citations are nonempty.
- Topic: `{id, state, summary, measure_ids, secondary_measure_ids, basis}` in the
  existing 21-topic order/release `1.0.0`. Basis is `{page_ranges, scope_evidence}`:
  sorted disjoint inclusive `{start, end}` PDF intervals and null or a retained
  source-evidence reference `{path, sha256}` establishing explicit document scope.
- Measure: `{id, primary_topic, action, conditions, institution, citations}`.
  IDs `m-1`, ... are stable local ordinals within the source revision, primary
  ownership is unique, and secondary references reuse them. Institution is null
  or a cited Statement naming the source-supported responsible/target institution;
  do not infer legal authority. Citations are the ordered deduplicated union of
  the measure's action/conditions/institution Statement references.
- Citation: `{id, page, block_id, start, end, excerpt}` with local IDs `c-1`, ...,
  original PDF indices and exact Unicode code-point block slices. Bind every
  citation to the top-level exact source/extraction revision; do not hash each
  citation/measure into a separate artifact identity.

Proposals require primary measures and a cited summary. Mentioned requires a
cited summary and no primary measures. Absence/unclear states have no factual
summary or measure references. Not-covered/not-applicable require full page
basis; not-applicable additionally needs explicit entire-topic scope evidence.
No universal institution/competence registry or taxonomy migration engine exists
in this version. Ambiguous scope is unclear, not guessed legal fact.

Keep bounded strict JSON (16 MiB, depth 64), 21 topics, at most 5,000 measures /
20,000 citations, 10 global/topic summary Statements, 20 conditions per measure,
1..2,000 code points per Statement/excerpt, 1..32 citation references per Statement
and bounded local IDs/display metadata. Required empty arrays are arrays, not
missing/null. Generated markup/unsafe links reject; original excerpts remain
literal escaped data. Unknown fields, duplicate IDs, invalid references, source
mismatch or draft-authored approval fields reject without panic or truncation.

## One provider and three clean roles

TQ-0024 implements an actual OpenAI-compatible configured provider adapter and
an offline fake. The trusted operator fixes endpoint/model/configuration; source
content cannot change transport or tool capabilities. The three active prompts
are normalizer, claim-validator and coverage-validator, each in a fresh stateless
conversation with only its own labelled evidence. The provider receives no shell,
filesystem-write or signing tool. Credentials stay in host transport only.

Record endpoint/server configuration identity, actual reported model identifier,
known revision/quantization/engine information and explicit unknowns. Do not
invent weight digests or require an unavailable weight-lock inventory to claim
known request configuration. Different validators/models are preferred, but one
model in distinct clean requests is allowed and disclosed. Reuse common strict
reply/limit/cancellation checks; neither an interface nor a fake canary proves
actual provider context isolation. Exercise the real adapter using recorded
synthetic smoke evidence before closing TQ-0024.

Default request bounds: 600 seconds, 8,192 output tokens, 1 MiB response, 8 MiB
text and up to 32 images/64 MiB within actual provider context limits. Tightening
is allowed; negative/unbounded limits reject. Budget the whole programme at
3,600 seconds/128 sequential invocations. Never truncate evidence silently.
Retry once only for a pre-response transient transport failure with unchanged
inputs and a new request ID; malformed/semantic failure is not retried to pass.
Cancellation discards late replies and closes the actual transport.

Each role has a machine-readable schema implemented by its owning task. Common
report fields are `version` (`atlas.report.local/v1`), `role`, `programme_id`,
`inputs`, `result`, `findings`, `checked`, `details`. Inputs are actual named
file/configuration hashes, not graph IDs. Result is pass/fail/insufficient/error.
Finding shape remains `{id, code, severity, target, evidence, message}` with
resolvable supplied evidence/targets. Pass requires complete checked inventories
and no unresolved warning/blocker. Details contain only the role's defined fields.
Checked is `{statement_ids, measure_ids, pages, topic_ids}`; host verifies exact
required sets. Claims account for every public Statement/metadata field; coverage
accounts for the full source page union and all topics/measures, including uncited
pages. Chunking must retain boundary context and a final document reconciliation.

Owning tasks must retain genuine model evaluations and original bounded replies.
Regression tests replay stored evidence offline; fake replies are conformance
fixtures, never claims of actual model execution. A model cannot authorize itself.

## Per-programme run record and eligibility

Version: `atlas.run/v1`; exact fields: `version`, `programme_id`, `source_checksum`,
`inputs`, `reports`, `outputs`. Input keys are exactly source-metadata, selection,
extraction, taxonomy, programme-schema, claim-schema, coverage-schema, draft,
normalizer-prompt, claim-prompt, coverage-prompt, model-configuration and
tool-configuration. Each value is null
until available or `{path, sha256}`. All references resolve to bounded regular
files under declared roots; no arbitrary URL, symlink or traversal lookup.

Reports is exactly `{claims, coverage, site}`; each is null or `{path, sha256}`.
Outputs is exactly `{markdown, html}` with the same reference/null shape. Keep
failed/insufficient reports and operational attempts separately; no timestamp,
selected-latest heuristic, model-authored approval flag or signature creates
authority. Current eligibility is derived, not stored as a trusted boolean.

The owner-controlled trusted checkout and local job are the admission boundary.
Proposal-only contributor reports are not copied into the approval set. Offline
verification rechecks all actual bytes, selected origin records, role results,
checked completeness and input equality. A changed source/extraction/draft/prompt/
schema/model/tool configuration makes that programme's reports stale. Reprocess
that programme; preserve old records. Hashes detect mismatch, not report authority
or semantic truth. No report PKI, DAG store or remote approval exchange is implied.

Markdown eligibility needs selected immutable source, acceptable extraction and
original evidence, structurally valid draft, matching complete claim/coverage
passes and deterministic safe output. TQ-0050 initially rejects HTML eligibility
as unavailable. TQ-0049 later extends it to require current parity/basic browser
quality reports bound to source, draft, output, supplied assets and tool config.
Neither stage deploys. Failed/draft/preview output never joins the public catalogue.

## Preview and static-site extension

Markdown is the first public format. Independently, TQ-0044 renders explicitly
flagged private fixture/draft previews from structurally valid data and provided
assets, without waiting for content approval or claiming public validation.
Default public HTML mode stays unavailable until TQ-0049 integrates actual current
content and basic site eligibility checks. TQ-0065 serves only an explicit generated preview root on loopback with no
uploads, tool execution, arbitrary filesystem exposure or automatic generation.
Preview has no dependency on publish-check; publication excludes its root.

TQ-0046/TQ-0047/TQ-0049 require parity, safe links, meaningful structure, keyboard
navigation/focus, desktop/mobile reflow and basic accessibility with existing tools.
Record actual browser/tool versions; do not impose binary-identical screenshots,
0.1% pixel thresholds, regional comparison or a visual model. Assets/fixtures must
be authored by Juan/non-LLM under the existing project prohibition.

## Versioning and rejected alternatives

These `*.local/v1` schemas are new contracts, not reinterpretations of the old
signed-artifact schemas. TQ-0063 removes inactive PKI/network/party-graph code and
updates Make/CLI/documentation/audit references while retaining useful strict JSON
and source/claim negative tests. Original task history and source bytes remain.
Reject skipping evidence/independence to shrink code, silently changing scope,
guessing official metadata, removing SSRF checks while retaining arbitrary-URL
admission, pretending a fake is a model evaluation, or treating source selection
as a manual content-publication override.
