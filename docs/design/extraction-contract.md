# Page-aware extraction and reproducibility

> Historical contract. [MVP.md](MVP.md) and [SPEC.md](../SPEC.md) govern the
> one-source adapter, original evidence and actual tool/configuration records.
> Full binary/image locking, a second-extractor platform and calibrated scores are deferred.

Contract: `atlas.extraction/v1`, TQ-0012. Implements the design boundary for
TQ-0030 through TQ-0035, following [extraction artifacts](../SPEC.md#extraction-artifacts),
[fidelity](../SPEC.md#extraction-fidelity) and [tests](../SPEC.md#extraction-tests).
Consumes [artifact identities](artifact-contract.md),
[source identity](source-identity.md) and [isolation](submission-isolation.md).

## Decisions and interfaces

`Extract(sourceRecordID, lockedProfile) -> extractionID` runs the pinned Docling
PDF adapter without OCR first. `SecondaryEvidence` and `RenderPages` run separate
Poppler adapters against the same original PDF. `Audit` consumes all three plus
original-page evidence. No adapter can approve its own fidelity. Docling native
JSON is retained as raw evidence; the stable Atlas representation below is the
contract, not an upstream serialization or extracted Markdown.

## Initial locked profile and acquisition

Retain Docling 2.145.0, the version recorded in the Spec trials. Those trials
reported successful conversion, OCR warnings, about 4.3/7.5 GB memory and
one/seven minutes for 64/172 pages. They did not record a complete model/runtime
lock and do not establish fidelity. The following explicit CPU baseline starts
a reproducible lineage; it does not claim to reconstruct undocumented trial
dependencies.

| Component | Exact initial selection |
| --- | --- |
| Platform | Linux x86-64 CPU, CPython 3.12.12, Debian 12 userland; no GPU/automatic accelerator selection |
| Converter | `docling==2.145.0`, `docling-slim==2.145.0` |
| PDF/data libraries | `docling-parse==7.22.2`, `docling-core==2.103.1`, `docling-ibm-models==4.0.2` |
| Inference | `torch==2.9.1` CPU, `torchvision==0.24.1` CPU, `transformers==5.16.0` |
| Layout | `docling-project/docling-layout-heron` revision `8f39ad3c0b4c58e9c2d2c84a38465abf757272d8` |
| TableFormer | `docling-project/docling-models` revision `fc0f2d45e2218ea24bce5045f58a389aed16dc23` (resolved upstream `v2.3.0`); accurate mode, cell matching enabled |
| Secondary text and page renderer | Poppler 26.10.0 (`pdftotext -bbox-layout`, `pdftoppm` PNG), encoding data 0.4.12 |
| OCR escalation | Tesseract 5.5.1 Git ref `0995615022dc368f65c6b88f31784d8bc384f855`; tessdata_best 4.1.0 ref `f0c187c74d3038e6d135fd080e30909af350eae0`, `spa` and `eng` only |

Identity evidence checked during design:

- PyPI Docling wheel SHA-256:
  `167e7a4736498585040d1680fa93fe582d8c033ee0fdc7e0a39d3b47fcfa545e`.
- Heron `model.safetensors`: 171,658,996 bytes, published LFS SHA-256
  `00333a43451945aaf89db8ca9c0a17e75d1537c17db60fdb91aa95f4c7929e0c`.
- Accurate TableFormer safetensors: 212,758,388 bytes, published LFS SHA-256
  `2a7d6c924b3cd12fb99a09280ca9c33a89c5d60b93253617d2e088c1a40374d9`.
- Poppler source archive SHA-256:
  `6792cb7c69205007ad87d2e936cecc5b3a31fac29ab54ffc3175fdb6b2a6ce35`;
  poppler-data archive:
  `c835b640a40ce357e1b83666aabd95edffa24ddddd49b8daff63adb851cdab74`.

These are upstream identities, not a claim that model weights were downloaded,
executed or qualified here. Sources: [PyPI release](https://pypi.org/pypi/docling/2.145.0/json),
[tagged Docling table adapter](https://github.com/docling-project/docling/blob/v2.145.0/docling/models/stages/table_structure/table_structure_model.py),
[Heron revision](https://huggingface.co/docling-project/docling-layout-heron/tree/8f39ad3c0b4c58e9c2d2c84a38465abf757272d8),
[TableFormer revision](https://huggingface.co/docling-project/docling-models/tree/fc0f2d45e2218ea24bce5045f58a389aed16dc23),
and [Poppler releases](https://poppler.freedesktop.org/).

TQ-0030 must acquire only these revisions in a separate explicit setup step,
verify every file's bytes, and commit a sorted inventory/fully resolved package
lock with hashes for all transitive dependencies, native libraries, fonts,
model configurations/tokenizers and OCR traineddata. Resolve annotated Git tags
to commits and record both; Git SHA-1 revision IDs are upstream locators, not
Atlas SHA-256 payload IDs. Build/seal the OCI image and record its actual digest
as a tool input. No worker can download assets or use `main`, `latest`, a shared
mutable cache, remote-code loading or automatic OCR/backend/model substitution.
Missing/incompatible pinned assets block implementation qualification; changing
the selected baseline requires a reviewed versioned profile, not a silent bump.

Use deterministic Torch operations, fixed seeds, one inference thread and one
page/job at a time, CPU-only, fixed locale/timezone, explicit PDF backend and
all options captured. The adapter uses the threaded Docling Parse backend with
one worker, explicit Heron model spec, TableFormer accurate options,
`do_ocr=false`, remote services and optional enrichment disabled. The locked
image includes only required fonts/models. Apply the isolation contract's
12-GiB/900-second/512-page defaults and bounded scratch/output limits.

## Structured JSON and coordinates

Top-level strict fields: `version`, `source_record`, `source_checksum`,
`page_count`, `pages`, `blocks`, `reading_order`, `warnings`, `tool_identity`,
`configuration_identity`, `model_identity`. References are manifest IDs; source
checksum agrees with the registered PDF. All page numbers are **one-based PDF
indices**, including covers/blank pages; printed page labels are nullable
display metadata only. Require exactly page_count contiguous page records.

Each page records `page`, `printed_label`, `width_mpt`, `height_mpt`, `rotation`,
`native_characters`, `extracted_characters`, `text_coverage_ppm`, `image_id`,
`handling` (`native`, `visual-required`, `ocr`, `unresolved`), and warning IDs.
Coverage is an explicitly versioned measurement, not a fidelity score; null is
required when the denominator cannot be established. Never treat null as 100%.

Each block has `id`, `page`, `bbox_mpt`, `type`, `text`, `order`, `parent_id`,
`origin` (`native`, `ocr`, `visual-reference`), and `details`. IDs are
`p<page>-b<ordinal>` assigned after stable ordering; ordinals start at 1.
Text retains original Unicode, numbers, punctuation and negations. `bbox_mpt`
is `[left, top, right, bottom]` in integer thousandths of PDF points, origin
top-left after page rotation and crop-box normalization, within page bounds.
Convert bottom-left coordinates by subtracting from normalized page height;
round to nearest thousandth, ties away from zero. Preserve raw upstream
coordinates separately. Unsupported transforms or out-of-page boxes block audit.

Types: `heading`, `paragraph`, `list`, `list-item`, `table`, `figure`, `caption`,
`header`, `footer`, `page-number`, `unknown`. Details are discriminated:
headings have levels 1..6; lists record ordered/bullet kind, item IDs and nesting;
tables record row/column counts and cells with zero-based row/column start,
positive spans, text, page and bbox. Reject overlapping/inconsistent table cells
and unresolved split tables; never silently flatten a table to prose. Every
cell remains page-mapped. Repeated headers/footers remain evidence even if
excluded from later summaries. Unknown blocks stay visible as audit warnings.

`reading_order` is a complete permutation of block IDs, consistent with per-page
order and list/table parentage. Explicitly represent columns as ordered groups
in block details; never interleave the left/right columns line-by-line. Missing,
duplicate or foreign IDs, cycles, gaps and page mismatches fail validation.
Do not invent a block for unreadable content: record an unresolved region and
original image reference with an audit blocker.

Warnings have `id`, `page`, nullable `block_id`, `code`, `severity`
(`info`, `warning`, `blocker`), `message`, and original tool diagnostic reference.
Document-wide warnings use page=null; page-local warnings must name the actual
page. No warning is discarded because the converter exited zero.

## Page escalation, independent evidence and examples

Native text is the default. Require visual handling when a page is textless,
has fewer than 20 non-whitespace native characters, any replacement characters,
a parser/layout/table warning, unknown/out-of-bounds ordering, or a native versus
secondary non-whitespace character-count difference greater than 10% (denominator
max(1, secondary count)). These are escalation heuristics, **not** calibrated
fidelity weights or automatic acceptance thresholds. Blank/decorative pages
still need an explicit visual disposition, never silent omission.

Render original pages at 150 DPI RGB PNG with fixed crop/rotation and no
annotations, bounded by the isolation pixel/scratch limits. Store image
hash/dimensions/page/renderer ID in a page-image inventory. Unrenderable pages
block. A page-level explicit OCR rerun uses pinned Tesseract on that image;
record language order, engine/traineddata/options/image IDs and raw OCR output.
Retain initial native output; produce a new extraction lineage identifying
replaced page blocks and the escalation reason. OCR text is not independently
approved merely because it is legible or agrees with the initial model.

Secondary evidence is Poppler page/bbox text from the original PDF, with the
same one-based map and tool/profile hashes. It corroborates fidelity only;
normalization and public rendering cannot read it instead of approved extraction
or canonical JSON. Derived extracted Markdown uses stable page markers and
block ordering but is never the citation/schema/public rendering contract.

| Worked representation | Expected handling |
| --- | --- |
| Native page 1, paragraph `p1-b1`, bbox `[1000,2000,9000,4000]`, text `No se reducirá el 10 %.` | `origin=native`, page=1, retain negation/number verbatim |
| Table on page 2, 2 rows/2 columns, four 1x1 cells with individual boxes | Table/details retain coordinates and row/column structure; audit against original |
| Bullet list on page 3 with two list-item children | Parent lists both IDs; each child retains page/order/nesting |
| Two-column page 4, left IDs b1/b2 then right b3/b4 | Reading order `[p4-b1,p4-b2,p4-b3,p4-b4]`, explicit column groups |
| Page 5 contains only a scanned image | `visual-required`; render/OCR escalation, never an empty successful page |
| Page 6 contains broken Unicode or suspicious number count | Warning+original image; unresolved until independent audit |
| Page 5 rerun with OCR | New extraction ID, `origin=ocr`, exact image/engine/traineddata provenance, original page=5 |

## Scenario-to-expected-outcome table

| Scenario | Expected outcome |
| --- | --- |
| Native PDF, all pages mapped | OCR disabled; structured JSON and derived Markdown available for audit |
| Image-only or blank page | Render and explicitly classify; unresolved content blocks approval |
| Incorrect page map / printed numbering mistaken for PDF index | Reject mapping, invalidate citations and downstream reports |
| Missing block/page or invented inserted text | Audit blocker; no score can override it |
| Converter warning with exit zero | Preserve page diagnostic, escalate; success is not fidelity approval |
| Table/list/column content loses structure/order | Audit blocker when material; preserve original evidence |
| Unrenderable or unreadable page | Unresolved blocker; do not skip or publish a partial programme |
| OCR succeeds | New provenance-bearing extraction candidate; independent audit still required |
| Model/tool/schema/options upgrade, even same wording | New input/artifact identity; exactly dependent reports stale |
| Identical pinned deterministic rerun | Byte-identical Atlas JSON/Markdown/image payloads, excluding external run logs; otherwise reproducibility failure |
| Runtime metadata changes only | No payload rewrite or hash change; retain outside canonical content |
| Assets missing, network requested, or incompatible pinned lock | Fail closed; no runtime download, automatic model/backend fallback or fake pin |
| Secondary text differs from Docling | Evidence for audit; never replace the canonical extraction silently |
| Golden PDF absent | No invented fixture or fidelity result; real end-to-end gates remain blocked |

## Versioning and rejected alternatives

New upstream/model/coordinate/ordering rules create versioned profile/schema
identities. TQ-0030/0031 must prove pinned rerun equality and page-aware native,
table/list/column/image/OCR fixtures; real golden-source gates follow recovery.
TQ-0032/0033 measure and audit before TQ-0034 calibrates any fidelity formula.
No formula weights are chosen here. Differences in actual images or content
must not be hidden by broad normalization to manufacture byte reproducibility.

Rejected: Markdown-only provenance, automatic whole-document OCR, upstream
`main` downloads, semantic agreement as byte equality, hidden warning filtering,
secondary output as canonical data, and converter exit status as approval.
