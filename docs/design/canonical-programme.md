# Canonical programme and citation contract

> Historical contract. The smaller `atlas.programme.local/v1` contract in
> [MVP.md](MVP.md) and [SPEC.md](../SPEC.md) governs active work. Retain cited
> factual content; graph bindings, per-citation hash identities and legal registries are superseded.

Version: `atlas.programme/v1`, TQ-0014; implementation TQ-0037/TQ-0038.
Implements [canonical data](../SPEC.md#canonical-programme-model),
[independent validation](../SPEC.md#llm-validation) and
[deterministic rendering](../SPEC.md#deterministic-rendering). Consumes
[artifact](artifact-contract.md), [source](source-identity.md),
[extraction](extraction-contract.md) and [taxonomy](taxonomy-jurisdiction.md)
contracts. Canonical JSON is the sole public-content source; extracted Markdown
is never parsed as the public rendering or citation contract.

## Decisions and strict schema interface

`ValidateProgramme(bytes, sourceRecord, extraction, taxonomy, jurisdiction)`
returns a validated draft structure or typed failures. This deterministic
structural/evidence check cannot mark the programme semantically validated.
Only independent claim/coverage passes and exact input-bound approvals do that.
All objects reject unknown fields; schema changes create a new version.

| Top-level field | Required value |
| --- | --- |
| `version` | `atlas.programme/v1` |
| `programme_id` | Exact registered programme identity |
| `source_record` | Source-record manifest ID, matching the extraction |
| `source_checksum` | Exact archived PDF SHA-256 |
| `extraction` | Approved extraction manifest ID |
| `taxonomy_version`, `taxonomy` | Exact release `1.0.0` and taxonomy manifest ID |
| `jurisdiction` | Exact profile/configuration manifest ID |
| `metadata` | Source-derived party/election/scope/geography/document fields below |
| `summary` | Ordered array of cited `Statement` objects |
| `topics` | Exactly 21 `Topic` objects in taxonomy order |
| `measures` | Ordered array of unique atomic `Measure` objects |
| `citations` | Unique `Citation` objects referenced by all factual statements |

`metadata` has exactly `party_id`, `party_name`, `document_title`,
`document_type`, `scope`, `geography_id`, `geography_name`, `election`,
`publication_date`, `rights`, `source_evidence`. IDs/dates agree with the source
record; display strings come from the verified registry/source metadata, not
model embellishment. `source_evidence` maps each factual display field to its
source-record/registry field/evidence reference. `document_title` is a cited
original-page Statement or null; when null, use the deterministic document-type
label, not a guessed title. Its evidence is its citations rather than an invented
source-record title field.
Election/publication-date nullability follows source identity. Metadata cannot
override the party registry, provenance state or original PDF link/checksum.

A `Statement` is exactly `{id, text, citations}`: stable local ASCII ID,
nonempty plain Spanish text, and nonempty array of citation IDs. Every global
summary, topic summary, measure text and optional measure condition is a
Statement. No free-form uncited subtitle, introduction, footnote or rationale
may enter a renderer via another field. Original quotations retain source
language; generated factual summaries are Spanish.

A `Topic` is exactly `{id, state, summary, measure_ids, secondary_measure_ids,
basis}`. `summary` is an array of Statements. Primary measures belong to exactly
one topic; secondary references reuse their ID without duplication. `basis`
contains `{page_ranges, applicability_rule}`: sorted disjoint one-based PDF page
intervals covering the complete source for absence claims, and a profile rule
ID or null. A profile rule is required for not-applicable. This is a request
for coverage validation, not a self-issued coverage report.

A `Measure` is exactly `{id, primary_topic, action, conditions, institution_id,
action_relation, jurisdiction_rule, citations}`. `action` is one Statement;
`conditions` is an array of Statements (empty if none); `citations` is the union
of its Statement citation IDs, deduplicated in source order. `institution_id`
must occur in the jurisdiction profile and source-supported evidence; if the
source does not name one, use the profile's explicit `unspecified` identity and
preserve uncertainty, not a guessed institution. `action_relation` is `direct`,
`shared`, `delegated`, `advocacy` or `unspecified`. Relations/targets must agree
with the cited action and competence rule; a model cannot turn advocacy into
direct governmental authority. An unresolved material institution ambiguity
blocks claim approval even if the draft is structurally valid.

## Citation semantics, stable IDs and ordering

A `Citation` is exactly `{id, source_record, page, block_id, start, end,
excerpt}`. `page` is the one-based original PDF index, not a printed label.
The block must belong to the same approved extraction/page/source revision.
`start`/`end` are zero-based Unicode code-point offsets, start inclusive and end
exclusive, inside that block's text. `excerpt` equals that exact nonempty slice.
No normalization of negations/numbers, ellipses inside the asserted slice, or
substitution from secondary text is allowed. For multi-block/multi-page support,
provide multiple citations, each with its own exact excerpt and page.
Validators still compare every claim with the original page/image, since
substring equality alone cannot establish meaning or polarity.

Citation IDs are `c-` plus the full SHA-256 hex of canonical
`{source_record, page, block_id, start, end, excerpt}`. Measure IDs are `m-` plus
the full hex hash of `{source_record, anchors, ordinal}`: anchors are the sorted
unique original page/block pairs, ordinal is a one-based action ordinal within
those anchors. A wording-only edit preserves identity; a changed source/anchor
or split/merge gets new IDs and explicit lineage in the normalization run log.
Do not hash prose alone or generate random IDs. Statement IDs are `s-<ordinal>`
assigned in deterministic document/topic/measure traversal order.

Measures sort by first supporting PDF page, first block order and action
ordinal; ties use ID byte order. Citations sort by page, block order, start,
end and ID. Topic order is fixed by the taxonomy. JSON object-key ordering and
integer encoding follow the artifact contract; no timestamps/approvals inside
the draft's hashed payload. A validation view combines the immutable draft and
its approval graph to expose draft/validating/validated/published state.

## Five-state and public-field examples

The following are synthetic examples with fixture-backed citations, not claims
about real parties. `C1` means the computed citation ID for original page 1,
block p1-b1, excerpt «Construiremos un centro de salud.»; `S1` is a Statement
with text «Propone construir un centro de salud.» and citations=[C1]. `M1` is
the computed measure anchored there, with a distinct action Statement S2 using
the same cited content and an evidenced regional health institution. These
aliases describe fixture construction, not literal
valid IDs or a whole 21-topic document.

| State/example | Valid draft representation and approval requirement |
| --- | --- |
| proposals/healthcare | summary=[S1], measure_ids=[M1]; source page/quote and responsible institution independently supported |
| mentioned/healthcare | Cited Statement «El texto menciona las listas de espera.»; measure_ids=[], no invented action |
| not-covered/housing | summary=[], measures=[], basis includes all original pages; applicable profile and whole-source coverage pass required |
| not-applicable/outside-topic | summary=[], measures=[], complete page basis plus explicit entire-topic scope rule; static Spanish status label, no inferred opposition |
| unclear/unreadable-topic | summary=[], measures=[], basis retains affected original page references; valid only as draft, publication blocked |

Every public factual field has one of two evidence forms:

| Public content | Evidence / invalid counterpart |
| --- | --- |
| Party, territory, election/date, document title, rights | Exact source/registry field references; a changed party name/date without evidence rejects |
| Global summary | Nonempty Statement citations; `text` with citations=[] rejects |
| Topic summary or factual explanation | Cited Statements only; an extra free-form `explanation` field rejects |
| Atomic promise and conditions | Cited action/condition Statements preserving polarity/numbers/scope; unsupported funding/number/institution rejects |
| Omission/non-applicability/unclear explanation | Deterministic taxonomy status label from state, page-coverage basis and legal rule; no generated factual prose shortcut |
| Source link/checksum, validation state/methodology | Derived from exact source and approval/configuration records, never a model-written URL or claim of validation |

For mentioned, summary must be nonempty and primary measures empty. For
proposals, primary measures and summary must be nonempty. Not-covered,
not-applicable and unclear have no asserted proposal or factual summary unless
the topic is reclassified to mentioned/proposals after resolution. Secondary
measure references constitute relevant content and therefore also prevent an
absence state; they may appear in mentioned alongside a cited summary but do
not manufacture another primary promise.

## Limits and trust boundary

Maximum canonical payload 16 MiB, JSON depth 64; exactly 21 topics; at most
5,000 measures, 20,000 citations, 10 global summary Statements and 10 per-topic
summary Statements. Each Statement is 1..2,000 Unicode code points; measure
conditions at most 20, each at most 2,000; excerpts 1..2,000; IDs bounded to 80
ASCII bytes. Metadata display fields at most 500 code points, original rights
statement at most 2,000. Citation count per Statement 1..32. Reject over-limit
input without truncation. Use safe integer bounds and strict JSON parsing.

Strings are plain data: reject NUL/control characters except newline/tab, raw
markup delimiters `<`/`>`, Markdown link/image/autolink syntax, embedded scripts
and executable URL schemes in generated Statement fields. A source excerpt can
contain literal punctuation/markup as a quotation: preserve it as text and
escape it deterministically in both renderers; never execute it or interpret
it as presentation code. No arbitrary URLs/paths in body fields, no model-authored
HTML/CSS/JavaScript/template fields. Link destinations are constructed from
validated source/citation metadata only. Language/faithfulness are semantic
validator responsibilities, not a simplistic ASCII/keyword filter.

## Scenario-to-expected-outcome table

| Scenario | Expected outcome |
| --- | --- |
| Complete 21-topic programme with cited Spanish Statements | Structurally valid draft; independent claim/coverage approval still required |
| Uncited global/topic summary or action/condition | Reject missing support, not merely an advisory finding |
| Page zero/out of range, wrong source revision/block/offset/excerpt | Reject citation mismatch; original-page reference mandatory |
| Duplicate citation/measure/Statement ID or dangling reference | Reject uniqueness/reference failure |
| Duplicate semantic promises under different IDs | Coverage validator blocks; never count both silently |
| Missing/reordered/unknown topic, stale taxonomy or schema | Reject contract mismatch, no inferred sections |
| Unknown or unsupported institution/power relation | Reject reference/competence mismatch or block unresolved draft claim |
| Multi-page promise with two valid excerpts | Multiple citations retained; no invented single-page support |
| Funding condition omitted from concise action | Claim validation fails, even when citation substring exists |
| not-covered without complete page basis, not-applicable without rule | Reject insufficient absence/applicability evidence |
| unclear state or unresolved extraction finding | Retain draft, no validated/public state |
| Source text contains prompt instructions | Evidence only; cannot grant tools, edit schema or approve output |
| Unsafe generated markup/URL field | Reject generated field; quoted original excerpt escaped as data |
| Content approval embedded inside the draft itself | Reject unknown field/cycle; lifecycle belongs to derived validation view |
| Viability, fulfilment, political rating or uncited editorial opinion | Reject outside source-facts contract |

## Versioning and rejected alternatives

TQ-0037 supplies machine-readable schema/positive and negative fixtures for
these fields and all states; TQ-0038 produces conforming drafts. Semantic skill
contracts may tighten evidence requirements, never weaken them. Migrations are
new schema/taxonomy inputs and immutable draft artifacts with revalidation.
Rejected: uncited topic prose, free-form model HTML, extracted-Markdown parsing,
mutable approval flags in the draft, random/text-only measure IDs, secondary
evidence as canonical citations, and treating structurally valid JSON as truth.
