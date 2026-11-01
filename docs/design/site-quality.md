# Rendering, comparisons and site quality

Contract `atlas.site/v1`, TQ-0016. Implements
[rendering](../SPEC.md#deterministic-rendering),
[site validation](../SPEC.md#site-validation),
[site tests](../SPEC.md#site-quality-tests) and
[topic ordering](../SPEC.md#topic-taxonomy). Consumes the
[canonical](canonical-programme.md), [taxonomy](taxonomy-jurisdiction.md),
[artifact](artifact-contract.md) and [semantic](semantic-validation.md) contracts.

## Authorship boundary and interfaces

TQ-0042 requires Juan or a non-LLM author to supply every template/HTML/CSS/
JavaScript/presentation asset. LLM agents may implement Go loading, validation,
data preparation, escaping and orchestration, but cannot author or edit those
assets indirectly through embedded strings, generated fixtures, code generators
or template factories. Missing supplied assets block rendering integration;
an LLM must not replace them with a minimal invented page.

`Render(validatedProgrammeView, assetBundle) -> Markdown + HTML + inventory`
uses Go deterministic rendering, without model calls. `CheckParity` and
`CheckBrowser` produce deterministic reports; `VisualQA` runs in a separate
semantic job only after these pass. None edits assets or canonical data to make
checks pass. No presentation code is supplied by this design.

## Required content, public cases and presentation

Render exact party/document/territory/election-or-non-electoral metadata,
provenance, archived original PDF link/checksum, revision relationship, validation
methodology/version, cited global summary, all 21 topics in taxonomy order,
every Statement/measure/condition/institution and original-page citation.
Source rights remain visible. Display no invented percentages or political
rankings. Omission text never implies opposition.

| Canonical state | Required Spanish presentation |
| --- | --- |
| proposals | «Contiene propuestas», cited summary and each atomic measure/condition |
| mentioned | «Mencionado sin propuestas concretas», cited discussion, no invented measure |
| not-covered | «No se aborda en el documento», visible section, not a negative political stance |
| not-applicable | «Fuera del ámbito del documento», distinct from an omission |
| unclear | «No se puede determinar con fiabilidad», draft/fixture preview only; publication blocked |

No topic disappears because it has no proposals. Use the exact Spanish labels
from the selected taxonomy. Metadata dates use an unambiguous Spanish display
with machine date identity retained. Non-electoral documents have no invented
election label/date. Cite the archived original PDF's one-based page, retain
excerpts, and expose multiple citations individually.

Supplied HTML must have Spanish document language, descriptive title, one main
landmark and one top-level heading, navigation, header/footer landmarks, logical
headings and a keyboard-visible skip-to-main link. All sections have unique
stable IDs. Public canonical objects have a machine-readable canonical-ID
attribute mapping back to the immutable view; parity checks also inspect actual
visible text/links, not just hidden attributes.

Tokyo Night dark appearance: background #1a1b26, principal text #c0caf5,
accent #7aa2f7; the asset author must choose variants satisfying actual computed
WCAG 2.2 AA contrast (4.5:1 normal text, 3:1 large text/non-text controls).
Base body text at least 16 CSS pixels, line height at least 1.5, visible focus,
non-colour-only status cues, readable long quotes and mobile spacing. No remote
font/assets, analytics, dynamic generated summaries, or required JavaScript
for reading/navigation. Static content must remain usable with scripts disabled.

## Content parity, links and trust boundary

Compare decoded DOM text and Markdown semantic content against the canonical
view's exact factual Statement/excerpt/measure text, metadata, IDs, ordering and
citations. Normalize only renderer-introduced whitespace boundaries, never
numbers/polarity/punctuation/source text. Require exact counts and one primary
measure occurrence; secondary cross-references must be labelled references.
Reject hidden/truncated factual content, duplicate/missing content and extra
uncited factual text. Fixed Spanish UI labels are an explicit versioned allowlist,
not unrestricted prose. Markdown escaping must preserve text as well as HTML.

Use Go's contextual HTML escaping; never cast untrusted text to trusted markup,
template definitions or executable content. Validate canonical generated fields
before rendering; literal original quotes remain escaped text. Paths/link
destinations derive only from verified source/citation IDs. PDF page fragments
are constructed integers within source page_count; all archive/local/anchor links
resolve inside the distribution. Original external source URLs are preserved
as evidence links, not fetched in ordinary offline QA. Reject executable URL
schemes, traversal, inline event/code injection and unexpected external requests.
Browser execution uses an isolated local preview origin with no general network,
credentials or archive/repository write access.

## Compatible regional comparisons

`Compare(programmeIDs, comparisonPolicy) -> comparisonView | error` requires
explicitly selected validated immutable revisions, same party, electoral type,
declared cycle, autonomous scope, document type and exact taxonomy version, plus
distinct mapped regional geographies in the same country. Election dates may
differ only within an explicit reviewed election-cycle mapping; do not infer
compatibility from year or party alone. Non-electoral/national/municipal/EU
documents cannot silently enter this regional comparison.

Preserve separate metadata/source/citation/provenance columns. Common measures
require an independently validated mapping of semantic equivalence with source
citations from each programme; literal similarity is not proof. Equivalence
classes carry member programme+measure IDs and mapping-report identity; partial
overlap or regional/funding/competence differences remain distinct. Show common,
territory-specific and topic-state differences without merging sources or
implying a missing proposal is opposition. Missing mapping evidence blocks the
common/different claim, not permission to invent it. Comparison artifacts bind
every programme revision/mapping/configuration so any change invalidates QA.

## Reproducible browser and visual checks

Initial reference tool profile: Linux x86-64 Debian 12 userland, Go CDP adapter
chromedp v0.14.2, Chromium headless 141.0.7390.37 (Playwright distribution
revision 1194), DejaVu Sans/DejaVu Sans Mono 2.37 supplied local font files,
axe-core 4.10.3 pre-authored audit distribution. This selects existing tools,
not LLM-authored JavaScript. The browser/version pairing is recorded in the
[upstream browser manifest](https://github.com/microsoft/playwright/blob/v1.56.1/packages/playwright-core/browsers.json).
TQ-0047 must seal actual browser/driver/font/axe/native-library bytes and their
hash inventory in the tool/profile artifact, qualify supported sandboxing, and
fail on missing assets or version substitution. No disabling the browser sandbox
as an automated fallback. Updates create a new profile and revalidation.

Capture at 1440x900 desktop, 390x844 mobile and 320x800 narrow-screen CSS pixels,
device scale 1, locale es-ES, timezone UTC, fixed colour scheme, and no reduced
motion ambiguity: all animations/transitions disabled by supplied assets or
browser emulation. Wait for document/font loading; all assets local and
successful, no console/resource errors. Timeout 30 seconds per load, 120 seconds
per page/profile; viewport screenshots plus complete-page coverage using
ordered overlapping tiles if a full image exceeds 24 million pixels. Retain
every original capture and tile/page coordinates; never crop away a failure.

Binary gates in addition to content/link/landmark checks:

- Keyboard: Tab/Shift-Tab reach every interactive control in a stable order;
  skip link works, focus is visible, no trap, Enter/Space operate controls as
  appropriate, topic navigation reaches actual sections.
- Accessibility: no axe violations of any severity in the required WCAG2/2.1/
  2.2 A/AA ruleset, no unresolved incomplete findings; manual/semantic screenshot
  inspection supplements rather than substitutes for deterministic checks.
- Reflow: no document-level horizontal overflow beyond 1 CSS pixel at any
  viewport or 200% text zoom. Scrollable comparison/table regions may overflow
  only inside labelled keyboard-accessible containers; factual text cannot be
  clipped or concealed. Long words/URLs/quotes must remain readable.
- Screenshot regression: fixed fixture corpus of all topic states, long text,
  hostile quotation, multiple citations and compatible regional comparisons.
  Same dimensions/profile, channel difference >16/255 marks a changed pixel;
  at most 0.1% changed pixels per image and no missing/extra tile. Any content,
  layout/overflow/a11y failure still blocks regardless of pixel threshold.

Baselines are generated artifacts with exact canonical-fixture, asset, browser,
font, viewport and QA-configuration hashes. Initial baseline establishment or
update is an explicit reviewed asset/QA change; never auto-accept current output
as the baseline. Keep old/new captures and a difference report, invalidate all
dependent site approvals and rerun deterministic/visual checks. This review is
of presentation/QA assets, not a human fallback gate for programme publication.
New programmes use the approved fixture-regression baseline and get their own
fresh complete screenshots; they need no invented per-programme golden image.

Deterministic site reports bind actual render/screenshot/inventory/baseline
digests and every viewport/check result. The visual job receives these plus
canonical data and images and emits the semantic contract's site report.
Approve only if deterministic **and** independent visual reports pass and all
inputs match. Timestamp-only run logs do not change output identity. Changed
template/browser/font/baseline/configuration requires new dependent evidence.

## Scenario-to-expected-outcome table

| Scenario | Expected outcome |
| --- | --- |
| All five topic states | Every section visible with distinct Spanish label; unclear fixture remains unpublishable |
| Long measure/quote/URL on narrow mobile screen | Readable/reflowed, no clipped text/document overflow; otherwise fail |
| Hostile markup in source quotation | Literal escaped text, no execution/network; unsafe generated field rejects |
| Keyboard navigation/skip link/focus works | Eligible deterministic pass; trapped/unreachable/invisible focus fails |
| Citation/condition/section omitted in either renderer | Parity/link gate fails, even with matching screenshot baseline |
| Same party but wrong election/type/scope/taxonomy/revision mapping | Reject regional comparison; no merged source identity |
| Claimed common proposal lacks both-programme evidence | Block mapping, no similarity-based substitution |
| Old screenshot after browser/template/font/canonical change | Stale, reject approval until fresh capture/QA |
| Pixels pass but deterministic content/a11y fails | Site remains blocked |
| Visual model pass but deterministic report fails/missing | Site remains blocked |
| Missing human-authored template/font/baseline | Block relevant rendering/QA task; LLM cannot create presentation assets |
| Browser cannot run with required isolation | Fail required gate; no unsandboxed/browser-skip pass |
| Legitimate reviewed asset/baseline update | New identities and complete deterministic/visual revalidation |

## Versioning and rejected alternatives

TQ-0042 supplies assets; TQ-0043/0044 implement deterministic Go plumbing;
TQ-0045 comparison mapping, TQ-0046 parity/links, TQ-0047 actual browser tests,
TQ-0048 visual skill, TQ-0049 the composite gate. All require real executable
fixtures/evidence, not this design's scenario review. Reject LLM presentation
authorship, screenshot-only semantic parity, automatic baselines, external-font
drift, hidden section omissions and visual overrides of deterministic failures.
