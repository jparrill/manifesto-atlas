# Manifesto Atlas Specification

Current scope: **local-first MVP**, approved by Juan through TQ-0061 following
the [product-scope audit](audits/2026-11-02-product-scope.md).
The [MVP contract](design/MVP.md) defines the current interfaces and milestones.
Earlier design contracts are historical where they conflict with these files;
their completed tasks and execution history remain preserved.

Owner-requested local demonstration: [macOS preview](design/macos-preview.md)
defines the current real-source E2E/HTML preview path (TQ-0062/TQ-0067). It uses
installed macOS tools and supplied presentation assets without waiting for Fedora.
Private preview outcomes are not production publication approval.

## Problem Statement

Political programmes are difficult to find, preserve and read. Manifesto Atlas
archives the original official documents and produces concise Spanish summaries
whose factual claims can be checked against the original PDF pages. It must not
invent proposals, confuse territories/elections, infer opposition from omission,
or present probabilistic automated validation as certainty.

## Solution

Deliver **one real document end to end before generalizing**:

1. A trusted maintainer selects an official document and retains its origin
   evidence, exact bytes, checksum and factual metadata.
2. Inspect and register that local source without modifying it.
3. Extract page-aware evidence in one supported bounded deployment.
4. Generate cited, topic-based canonical JSON through one configured provider.
5. Independently validate every public claim and complete-source coverage.
6. Produce readable, deterministic Spanish Markdown and an offline eligibility
   check based on actual current inputs and reports.
7. Repeat with the second golden document. Extend to supplied static HTML assets,
   basic accessibility and local preview independently of regional comparisons.

Source selection is an administrative authority boundary, not a manual override
of content validation. An uploaded file or contributor report cannot grant itself
source admission or publication authority. The MVP does not automatically fetch
arbitrary contributor URLs, certify party domains or accept third-party approvals.

## User Stories

1. A maintainer can preserve a selected official PDF with its origin evidence.
2. A reader can inspect the archived original, checksum and source-page citations.
3. A maintainer can process one document without first building an admission platform.
4. A reader sees comparable topic organization and explicit absence/uncertainty.
5. Every public summary, measure and condition has original-source support.
6. Separate claim and whole-source coverage checks block unsupported publication.
7. Changed inputs invalidate that programme's results; old sources/reports survive.
8. A maintainer can reproduce deterministic output without rerunning models in CI.
9. A reader receives Spanish Markdown and, in the site extension, accessible HTML.
10. Contributors may propose sources, but cannot self-approve bytes or reports.
11. Methodology states model limitations and separates facts from future assessments.
12. A local draft preview provides feedback without placing drafts in a public catalogue.

## Implementation Decisions

### Module seams

Use a small local pipeline: intake/registration, extraction/evidence, generation,
independent validators, rendering and eligibility. Reuse strict JSON/hash helpers.
A simple per-programme run record binds actual file hashes, prompts, schemas,
model/tool configuration and reports. No content-addressed object store, 26-kind
artifact graph, report-signing PKI or universal agent framework is required.

Go and Make own deterministic operations. Schema-constrained skills own semantic
judgment. Start with one actual provider and three roles: normalizer, claim
validator and coverage validator. Recompute the affected programme initially;
precise transitive cache optimization is outside the MVP.

### Make targets

The names remain discoverable; implementation tasks make them real incrementally:

| Target | MVP behavior / owner |
| --- | --- |
| `make submit` | Inspect a selected local candidate; TQ-0021 |
| `make verify-source` | Verify local selection evidence and matching bytes, **no network certification**; TQ-0027 |
| `make add` | Register immutable selected source; TQ-0027 |
| `make extract` | Page-aware extraction; TQ-0030 |
| `make validate-extraction` | Deterministic page/warning blockers, no numeric score; TQ-0030 |
| `make normalize` | Cited Spanish draft; TQ-0038 |
| `make validate-content` | Independent claim and coverage passes; TQ-0041 |
| `make build` | Markdown first (TQ-0043), private supplied-HTML previews (TQ-0044), public HTML integration later (TQ-0049) |
| `make validate-site` | Basic deterministic HTML quality extension; TQ-0049 |
| `make publish-check` | Offline Markdown eligibility first (TQ-0050), HTML extension in TQ-0049 |
| `make process` | One-programme Markdown pipeline; TQ-0051 |
| `make serve` | Explicit private fixture/draft preview, no publication prerequisite; TQ-0065 |

TQ-0063 updates old stub ownership/messages and uses `go test ./...` for global
regression while retaining focused targets. No stage may report fake success.
`make verify` remains offline: no model, source-site retrieval or automatic setup.
Explicit real converter/model evaluations are separate required owning-task gates.

### Provenance trust model

The MVP stores **maintainer-selected source** evidence, not an assertion that an
automated official-origin certification ran. Record the official URL, retrieval
time, exact PDF checksum and supporting page/record evidence; preserve original
publisher rights. Missing or contradictory origin evidence blocks selection.

Acquisition of the exact initial documents is explicit maintainer-authorized work
with existing tools, before pipeline registration. Unknown URLs/metadata must not
be invented. Contributor proposals are data only. An owner-controlled source
selection record identifies the accepted metadata/PDF/evidence hashes; that record
is never generated or activated by a model or untrusted payload.

Arbitrary-URL admission, independent automatic canonical downloads, expiring
domain/key grants, publisher-signature verification and semantic provenance
skills are deferred. If brought back, their network/SSRF protections (including
IPv6 and DNS/redirect validation) remain necessary for that expanded boundary.

### Untrusted document handling

Selected PDFs and model outputs remain untrusted data. Enforce safe rooted paths,
PDF MIME/magic/parser checks, byte/page bounds, deadlines, memory limits and
whole-job cleanup. PDF workers cannot write the repository, read home/secrets or
use general network access. TQ-0062 chooses one minimal supported Fedora
deployment using installed tools and actual bounded denial/resource checks;
TQ-0021 implements it. Rootless Podman, custom OCI-image sealing and a generalized
platform qualification laboratory are not prerequisites imposed by this Spec.

Use Fedora qpdf for inspection; record its actual package/tool version. Do not
run a converter unsafely just because qpdf is installed. A missing required
enforcement mechanism is a genuine failure, never a skipped passing gate.

### Immutable source storage

Store selected originals and factual metadata under a stable year/programme/
checksum hierarchy. Year basis is explicit publication or retrieval, never an
invented election date. Exact duplicates are idempotent; corrected bytes create a
new revision referencing the previous one. Safe atomic installation must never
replace committed PDF bytes. Git records selection and revision history.

### Extraction artifacts

Use one actual Docling adapter, recording its version/configuration and source
checksum, with OCR off for native text. Keep page count, contiguous original PDF
page indices, block IDs/text/order and warnings. Preserve native structured
output for audit; do not build a second general-purpose layout representation
beyond what actual citations and page inspection require.

Original-page images and, where useful, existing secondary text tools corroborate
evidence. Required unreadable, scanned or materially ambiguous regions remain
blocked/insufficient if the selected tools/model cannot resolve them. No mandatory
whole-source second-extractor framework or automatic OCR escalation engine is
needed for the initial native-document path.

### Extraction fidelity

Check page completeness, page references, warnings and observed headings/numbers/
negations. Material extraction defects cannot be hidden by a score. Independent
claim/coverage validation examines original evidence and detects remaining
meaning/omission problems. Numeric fidelity formula/calibration and a separate
extraction-audit skill are deferred; deterministic blockers move into TQ-0030
and original-page/full-source checks into TQ-0031/TQ-0039/TQ-0040.

### Canonical programme model

Versioned canonical JSON is the public-content source. It contains source-bound
metadata, cited global/topic summaries, 21 ordered topics, atomic measures and
conditions, source-supported institution text and exact page/excerpt citations.
Use local stable IDs within a source revision, not a new hashed-identity protocol
for each measure/citation. Extracted Markdown is never the public rendering input.
Draft JSON cannot assert approval; eligibility is derived from current reports.

### Topic taxonomy

Retain the existing 21 topic IDs/order and Spanish labels, release `1.0.0`, in
`internal/taxonomy`. States remain `proposals`, `mentioned`, `not-covered`,
`not-applicable`, `unclear`. Omission is not opposition. Entire-topic exclusion
requires explicit evidenced document scope; missing/ambiguous applicability is
unclear. Shared powers or advocacy must not become blanket topic exclusions.
An all-institution legal registry, comparison engine and split/merge migration
framework are deferred; unknown current versions still reject.

### Agent skills

Implement `manifesto-normalizer`, `manifesto-claim-validator` and
`manifesto-coverage-validator` through one configured provider with separate clean
requests. Fixed role prompts, machine-readable output schemas and bounded
evidence inventories are required. Source text is evidence, never instructions.
Skills have no shell, arbitrary network, write or signing tools. Provenance,
dedicated extraction-audit and visual-site skills are deferred.

### LLM validation

Public eligibility requires **both** independent claim and whole-source coverage
passes for the exact source/extraction/draft/prompt/schema/model/tool inputs.
Claims include every public title/summary/action/condition and source-derived
metadata. Coverage accounts for every original page/topic/measure, not only
generated citations; use bounded explicit chunks/reconciliation when needed.

No shared generation transcript or reused conversation is allowed. One model
is acceptable in separate stateless requests; prefer distinct models when
available. Record actual identity/configuration, disclose unknown revision/weight
details rather than inventing pins, and require reruns when configuration changes.
Strict JSON, checked-target completeness, bounded outputs, cancellation and
pass/finding consistency are deterministic host checks. Fail, insufficient,
missing evidence, stale inputs or disagreement block publication. No manual
content-approval fallback or retry-until-pass is allowed.

### Deterministic rendering

Render readable Spanish Markdown first from canonical data, preserving all
topics/states/measures/conditions and original-page references. HTML is a separate
extension using supplied human/non-LLM-authored Tokyo Night assets. LLMs must not
author HTML, CSS, JavaScript or templates, including fixtures/embedded strings.
Keep contextual escaping, safe links and byte-deterministic unchanged builds.

### Site validation

The HTML extension requires content parity, structure/landmarks, working source
links, keyboard navigation, focus and desktop/mobile reflow/accessibility checks.
Use existing browser/audit tooling, recording its version. No mandatory exact
browser/font binary lock, pixel-difference threshold or visual LLM approval is
needed. Missing required basic checks block HTML eligibility, not Markdown.
Regional comparison support is not a prerequisite of single-programme pages.

### Publication and CI

The trusted local processing job creates current reports; contributor-written
reports do not enter the approval set. A small run manifest records input/report/
output hashes. Offline checks revalidate schemas, exact bytes and required
passing results against the selected trusted checkout; they do not prove a
model execution from a bare `pass` field or a checksum alone. Report PKI and
generic artifact graphs are removed by TQ-0063, not replaced by blind trust in JSON.

`make publish-check` emits eligible Markdown output first; HTML requires its
additional site checks. Draft/fixture preview output has a separate root and
never joins the public catalogue. GitHub Pages/CI workflow design stays in
owner-only inbox TQ-0017; this plan does not authorize deployment implementation.

### Language and licensing

Public prose/output is Spanish. Code, schemas, prompts, tests, task content and
technical documentation are English. Original project code/docs use MIT;
archived third-party documents retain publisher rights.

### Assessment separation

No viability, political ranking, historical reliability, fulfilment percentage or
combined confidence score is part of the MVP source-facts layer.

## Testing Decisions

### Submission and provenance tests

Test selected source evidence/checksums, wrong bytes, spoofed MIME, encrypted/
corrupt PDFs, unsafe paths, limits and cleanup. Acquisition fixtures are local;
automatic network provenance is not claimed.

### Source registry tests

Test immutable bytes, duplicate idempotency, new revisions, metadata mismatch,
unsafe paths and interrupted installation. Recovery precedes registration code.

### Extraction tests

First run one real approved document plus synthetic boundary fixtures; retain
original page maps, anchors, warnings and tool/resource observations. Repeated
unchanged deterministic adapters must not silently drift. The second golden
source is a later E2E milestone, not a prerequisite to the first extraction.

### Extraction audit tests

Missing pages, wrong mapping, altered numbers/negation, unreadable regions and
major omissions must be detected by deterministic checks or independent original-
page claim/coverage evaluations. No numeric score is an acceptance shortcut.

### Canonical model tests

Test strict schema/required fields, 21 ordered topics, all five states, unknown
versions, uncited public prose, invalid page/block/code-point excerpts, duplicate
IDs, source metadata mismatch and forged validation fields.

### Validator contract tests

Keep offline pass/fail/insufficient/disagreement and injection fixtures. Each
semantic owning task must actually run a configured model and retain original
bounded outputs with exact inputs and observed identities. Synthetic evaluation
evidence is allowed when labelled; fake replies never prove actual prompt behavior.

### Rendering tests

Test Spanish labels/order, exact measure/condition/citation preservation,
escaping, stale-input rejection, deterministic bytes and explicit draft previews.

### Site quality tests

Test actual supplied assets for parity, safe links, landmarks, keyboard use and
mobile/desktop readability. Do not require comparisons, pixel baselines or a
visual model to validate one programme page.

### End-to-end tests

TQ-0064 proves the first real programme through Markdown and offline eligibility,
including deliberate failures; TQ-0054 repeats with the second source. Stored
evidence permits offline regression without calling a model or source website.

## Out of Scope

- Public upload services, accounts, runtime backend, database/CMS or embeddings.
- Automatic arbitrary-contributor URL/download admission and publisher-key/domain PKI.
- Signed per-role report approvals, content-addressed object stores and generic DAG engines.
- Numeric fidelity scores, a general legal-competence registry and taxonomy migration engines.
- Regional comparisons, visual-model site approval and pixel-exact screenshot baselines.
- Untrusted contributor report self-approval, manual bypass of failed content checks,
  LLM presentation authorship and automatic deployment implementation.
- Political scoring or automatic verification of every statistic a party asserts.

## Further Notes

The exact approved Frente Obrero general-political and VOX Castilla y León
documents are the initial source targets; both are retained as local candidate
fixtures under `fixtures/` with evidenced origin metadata.
Never invent URLs, checksums or election dates, substitute another revision, or
call a candidate an automatically verified official source. If exact recovery
is unavailable, leave its task open with concrete evidence of the blocker.
qpdf 12.3.2 is installed from Fedora; tool availability is not a completed pipeline.
Docling trials are historical observations, not current extraction approval.
Taskqueue is authoritative for current milestones; completed legacy tasks are
historical implementation records, not permission to reintroduce deferred scope.
