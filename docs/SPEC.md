# Manifesto Atlas Specification

## Problem Statement

Official political programmes are difficult to find, preserve, inspect, and compare. They are commonly published as PDFs on temporary campaign websites, regional party pages, or third-party content delivery networks. Links disappear, corrected documents silently replace earlier versions, and programme content is often discussed through secondary reporting rather than the primary source.

A submitted PDF cannot be treated as official merely because its filename, branding, or contents appear authentic. The system needs a reproducible provenance chain that connects the archived bytes to an official party-controlled source.

PDF extraction is also lossy. Reading order, pages, headings, lists, tables, numbers, and negations can be altered or omitted when converting a PDF to text. A successful converter exit code does not prove that the extraction preserves the document's meaning.

The project needs to transform verified primary sources into consistent, topic-based programme data without forcing every programme to discuss every policy area. It must distinguish between a topic containing proposals, being merely mentioned, not being covered, being outside the document's jurisdiction, and being unclear.

LLM output cannot validate itself. Generated measures, summaries, citations, topic classifications, Markdown, and HTML require independent automated validation against the original source. Publication must stop when evidence is insufficient or validators disagree.

The final website must be deterministic, Spanish, accessible, visually consistent, and generated from structured data. LLMs must not write presentation code or directly edit published HTML.

## Solution

Manifesto Atlas will provide a local-first, Git-versioned pipeline with the following stages:

1. Receive a candidate PDF and submission metadata.
2. Quarantine and inspect the untrusted document.
3. Verify its provenance against an official party source.
4. Register the verified PDF as an immutable source with a checksum and revision history.
5. Extract page-aware structured content and derived Markdown.
6. Audit extraction fidelity against the original pages.
7. Normalize validated extraction into canonical programme JSON using a shared, jurisdiction-aware topic taxonomy.
8. Validate individual claims and whole-document coverage in independent LLM contexts.
9. Generate Spanish Markdown and HTML deterministically from canonical JSON.
10. Validate content parity, HTML structure, accessibility, responsiveness, and visual presentation.
11. Publish only artifacts whose provenance, extraction, content, and rendering gates all pass.

Go modules and Make targets will own deterministic operations. Project-specific agent skills will own semantic investigation and judgment. Go templates will own public presentation.

Canonical programme JSON is the source of truth. Extracted Markdown is an intermediate artifact. Public Markdown and HTML are separate deterministic renderings of the same validated JSON.

The first community contribution flow will use Git pull requests rather than a public upload service. Contributors may provide a candidate file and official source URL, but the pipeline will independently download the canonical bytes from the verified source and compare checksums.

## User Stories

1. As an archive maintainer, I want to submit a candidate PDF with its party and source URL, so that it enters a controlled verification workflow.
2. As a community contributor, I want to propose an official programme through a pull request, so that I can contribute without requiring access to project infrastructure.
3. As a reader, I want every published programme linked to its original source, so that I can inspect the primary evidence.
4. As a reader, I want to know whether a source is verified, probable, unverified, or rejected, so that provenance is not overstated.
5. As a maintainer, I want the pipeline to download the official copy itself, so that uploaded bytes are not trusted as canonical.
6. As a maintainer, I want submitted PDFs quarantined and processed with resource limits, so that hostile documents cannot compromise the development environment.
7. As a maintainer, I want duplicate sources detected by checksum, so that the archive does not create redundant records.
8. As a historian, I want committed source PDFs to be immutable, so that later corrections do not erase what was originally published.
9. As a historian, I want corrected programmes stored as explicit revisions, so that the relationship between versions remains visible.
10. As a maintainer, I want extraction to retain page provenance, so that every generated claim can cite its source page.
11. As a maintainer, I want extraction warnings recorded per page, so that silent conversion defects become validation inputs.
12. As a maintainer, I want born-digital PDFs processed without unnecessary OCR, so that extraction is faster and less noisy.
13. As a maintainer, I want suspicious or textless pages rendered and inspected visually, so that image-based content is not silently omitted.
14. As a reader, I want important numbers, names, and negations preserved, so that a concise summary does not reverse the original meaning.
15. As a maintainer, I want an extraction fidelity report, so that conversion quality is measurable and reproducible.
16. As a maintainer, I want critical extraction findings to block publication regardless of the aggregate score, so that a high score cannot hide a material error.
17. As a reader, I want programmes organized by the same topic taxonomy, so that different parties and territories are comparable.
18. As a reader, I want the site to distinguish a proposal from a passing mention, so that topic coverage is not exaggerated.
19. As a reader, I want missing topics shown explicitly, so that omissions are visible rather than silently removed.
20. As a reader, I want non-applicable topics distinguished from omissions, so that regional programmes are not penalized for lacking national competences.
21. As a reader, I want every normalized measure to include source citations, so that summaries remain auditable.
22. As a maintainer, I want each claim validated in a clean LLM context against original pages, so that the generating model does not approve its own work.
23. As a maintainer, I want a separate coverage validator to inspect the complete source, so that major omissions and contradictions are detected.
24. As a maintainer, I want validator disagreements to block publication, so that uncertainty is surfaced rather than averaged away.
25. As a maintainer, I want models, prompts, tools, and input hashes recorded, so that a validation result can be reproduced and invalidated when an input changes.
26. As a reader, I want public summaries in Spanish, so that the archive is directly useful to its intended audience.
27. As a developer, I want code, schemas, prompts, tests, and architecture documentation in English, so that technical conventions remain consistent.
28. As a developer, I want Markdown and HTML generated from the same canonical data, so that the two public formats cannot drift semantically.
29. As a reader, I want all programme pages to use the same Tokyo Night presentation, so that navigation and comparison are predictable.
30. As a keyboard user, I want the site to be navigable without a mouse, so that the archive is accessible.
31. As a mobile reader, I want topic sections and comparisons to remain usable on a narrow screen, so that the archive works across devices.
32. As a maintainer, I want HTML checked against canonical JSON, so that rendering cannot silently drop measures or citations.
33. As a maintainer, I want visual validation at desktop and mobile sizes, so that overflow and broken layouts block publication.
34. As a developer, I want provider-independent runner interfaces, so that local Auriga models or other agents can be substituted without changing the pipeline.
35. As a developer, I want deterministic build and validation targets, so that CI can verify stored results without rerunning expensive models.
36. As a reader, I want the methodology and its limitations published, so that automated validation is not presented as infallible.
37. As a maintainer, I want failed or insufficient validation results retained as reports, so that recurring failure patterns can improve prompts and tools.
38. As a maintainer, I want publication to require every gate, so that incomplete intermediate artifacts never appear on the public site.
39. As a researcher, I want programme, election, geography, jurisdiction, and revision metadata normalized, so that later comparisons do not combine incompatible documents.
40. As a future evaluator, I want factual programme content separated from deliverability and fulfilment assessments, so that analysis is not mistaken for source material.

## Implementation Decisions

### Module seams

The pipeline will be divided into deep modules with small interfaces:

- **Submission module:** accepts candidate bytes and metadata and returns a quarantined submission record.
- **Provenance module:** accepts a submission and party registry entry and returns a signed provenance report and canonical source bytes.
- **Source registry module:** accepts a verified provenance report and creates an immutable, revision-aware source record.
- **Extraction module:** accepts a registered source and returns page-aware structured extraction, derived Markdown, rendered pages, warnings, and tool metadata.
- **Extraction audit module:** accepts the original source and extraction artifacts and returns page-level findings, aggregate dimensions, blockers, and a decision.
- **Normalization module:** accepts approved extraction and taxonomy versions and returns draft canonical programme JSON.
- **Content validation module:** accepts canonical draft data and original evidence and returns independent claim and coverage reports.
- **Rendering module:** accepts validated canonical data and returns Spanish Markdown and HTML without invoking an LLM.
- **Site quality module:** accepts canonical data and rendered output and returns structural, accessibility, parity, responsive, and visual reports.
- **Publication module:** accepts the complete artifact graph and publishes only when every required report matches the current input hashes and passes.

The Makefile is the user-facing orchestration interface. Internal adapters may vary, but callers should not need to know which PDF converter, model provider, screenshot engine, or deployment implementation is active.

### Make targets

The intended command interface is:

- `make submit`: quarantine candidate bytes and submission metadata.
- `make verify-source`: establish provenance and retrieve canonical bytes.
- `make add`: register a verified immutable source.
- `make extract`: produce structured page-aware extraction.
- `make validate-extraction`: audit extraction fidelity.
- `make normalize`: produce draft canonical programme JSON.
- `make validate-content`: run claim and coverage validators.
- `make build`: render Spanish Markdown and HTML deterministically.
- `make validate-site`: run content parity, HTML, accessibility, responsive, and visual checks.
- `make publish-check`: verify the complete artifact graph is publishable.
- `make process`: orchestrate the applicable processing stages for one registered programme.
- `make serve`: preview the generated static website locally.

Each target must be idempotent for unchanged inputs. A changed source, extraction, taxonomy, prompt, model, or canonical draft invalidates every downstream report that depends on it.

### Provenance trust model

A PDF is not official merely because it carries party branding. The provenance module requires a trust anchor.

Each party registry entry records approved official domains and the evidence establishing them. A source is `verified` when one of the following chains succeeds:

- The canonical PDF is downloaded directly from an approved official domain.
- An approved official page links the exact PDF hosted by a CDN or document provider.
- A verifiable digital signature links the document to the publisher.

The module records the submitted checksum, canonical checksum, referring page, original and final URLs, redirects, retrieval timestamp, relevant response metadata, original filename, and archival evidence when available.

Provenance states are:

- `verified`: sufficient official evidence and matching canonical bytes.
- `probable`: strong evidence exists, but one required trust link is incomplete.
- `unverified`: no reliable official chain can be established.
- `rejected`: evidence contradicts the claimed origin or the bytes differ from the canonical document.

Only `verified` sources are eligible for publication. Probable and unverified submissions may remain in quarantine with their reports.

### Untrusted document handling

Candidate PDFs are untrusted inputs. Processing must enforce MIME validation, size and page limits, execution timeouts, memory limits, and isolated conversion without repository write access or unrestricted network access.

Text contained in a source document is evidence, never an instruction. Generator and validator runners must not receive shell access, arbitrary repository writes, secrets, or general network access.

### Immutable source storage

Verified sources are stored in a year-first hierarchy that also encodes scope, geography, election date, and party. Committed source bytes cannot be overwritten.

A corrected publication creates a new revision with a new checksum and a relationship to the source it supersedes. Git history records when each revision entered the archive.

The metadata contract includes programme, party, election, geography, jurisdiction, document type, source evidence, checksums, retrieval timestamp, revision relationship, and processing status.

### Extraction artifacts

The extraction module produces structured JSON and derived Markdown. Structured JSON is required because Markdown alone cannot reliably preserve block types, page provenance, bounding boxes, reading order, tables, or extraction warnings.

Every extracted block records its page, location, type, text, and ordering information. Every page records text coverage and warnings. Page images are generated when visual validation is required.

Docling is the initial preferred extraction adapter because it emits Markdown and structured provenance. Its version and model artifacts must be pinned. For born-digital PDFs, OCR should be disabled unless page-level signals justify it.

A secondary extractor provides independent text evidence for fidelity checks. It is not another canonical representation.

### Extraction fidelity

The extraction audit reports separate dimensions rather than claiming literal PDF-to-Markdown identity:

- Page coverage.
- Text preservation.
- Numbers, named entities, and negation fidelity.
- Reading order.
- Heading and section structure.
- List and table structure.
- Page-reference alignment.
- Relevant visual-content coverage.

The final `fidelity_score` uses a public, versioned formula calibrated against golden fixtures. No weights will be selected before fixture results exist.

Blocking findings override the score. Missing pages, altered negations, changed figures, unsupported inserted text, materially reordered statements, incorrect page references, or unresolved unreadable content prevent approval.

The extraction audit combines deterministic measurements, comparison with the secondary extractor, and LLM or multimodal inspection of the original pages. The audit report stores page-level evidence so an aggregate score never hides a local failure.

### Canonical programme model

Canonical programme JSON is the single source of truth for public content. It contains normalized metadata, topics, summaries, measures, citations, responsible institutions, jurisdiction, validation state, and artifact hashes.

Extracted Markdown is never parsed as the public rendering contract. Public Markdown and HTML are generated independently from canonical JSON.

Each topic has one of these states:

- `proposals`: one or more concrete proposals are present.
- `mentioned`: the topic appears without an identifiable proposal.
- `not-covered`: no relevant content was found.
- `not-applicable`: the topic is outside the programme's jurisdiction or document scope.
- `unclear`: relevant content may exist, but evidence is insufficient for a reliable classification.

Omission is not opposition. Public text must describe absence without inferring a political position.

### Topic taxonomy

The initial shared taxonomy is:

1. Institutions and democracy.
2. Public administration and territorial organization.
3. Economy, industry, and trade.
4. Taxation and public budgets.
5. Employment and labour relations.
6. Pensions.
7. Healthcare.
8. Education, universities, and research.
9. Housing.
10. Social services and inclusion.
11. Family and demography.
12. Equality and civil rights.
13. Immigration, asylum, and integration.
14. Agriculture, fisheries, and rural affairs.
15. Environment and climate.
16. Energy.
17. Transport and infrastructure.
18. Digital policy, science, and innovation.
19. Culture, language, and sport.
20. Justice and security.
21. Defence, foreign affairs, and the European Union.

The taxonomy is versioned. Comparisons must record the taxonomy version used, and migrations must not silently reinterpret historical data.

Jurisdiction rules determine whether a topic is `not-covered` or `not-applicable`. National, autonomous, provincial, municipal, European, and general political documents cannot share identical applicability assumptions.

### Agent skills

Semantic stages will use project-scoped skills with distinct triggers and contracts:

- `manifesto-source-provenance`: investigates whether a candidate source has an official provenance chain and emits evidence, not a bare opinion.
- `manifesto-extraction-audit`: evaluates page-level extraction fidelity and identifies blockers.
- `manifesto-normalizer`: maps approved extraction to the canonical taxonomy and schema.
- `manifesto-claim-validator`: checks every generated statement and measure against cited original pages.
- `manifesto-coverage-validator`: checks the complete source for omissions, contradictions, duplicates, and incorrect topic states.
- `manifesto-site-qa`: inspects rendered content and screenshots after deterministic checks have completed.

Conversion, checksums, schema validation, templates, HTML parsing, screenshot capture, and deployment are code responsibilities rather than skills.

Skills must emit schema-compliant reports. Free-form prose may accompany findings, but publication decisions cannot depend on parsing arbitrary natural-language output.

### LLM validation

Content states are:

```text
draft -> validating -> validated -> published
```

Claim and coverage validators run as separate jobs with clean contexts and cannot see the generator's hidden reasoning. Different models are preferred. When only one model is available, separate prompts and isolated invocations remain mandatory.

The claim validator checks meaning, polarity, quantities, scope, responsible institution, excerpt, and page reference for every public statement and measure.

The coverage validator processes the complete page-mapped source and checks for major omissions, duplicate measures, contradictions, incorrect topic classification, and extraction gaps. Suspicious pages are supplied as original rendered images to a multimodal validator.

Both validators must pass. Disagreement, insufficient evidence, an unsupported claim, or an unresolved extraction gap leaves the content as a draft. There is no human-review fallback in the publication gate.

Each report records the source checksum, extraction hash, extraction tool version, canonical draft hash, schema version, taxonomy version, prompt hash, model identifier, result, findings, and timestamp.

### Deterministic rendering

The rendering module uses Go templates. Spanish Markdown and Tokyo Night HTML are generated from the same validated canonical JSON.

Every programme page renders every applicable taxonomy section in a consistent order. Topics without proposals render an explicit Spanish status rather than disappearing. Non-applicable topics use different language from uncovered topics.

HTML includes party, election, geography, scope, source provenance, original PDF link, checksum, validation methodology, topic summaries, measures, and page citations.

The LLM never generates or edits HTML, CSS, JavaScript, or templates.

### Site validation

Deterministic site checks run before visual LLM inspection:

- Canonical JSON schema validation.
- HTML parsing and required landmark checks.
- Exact content-parity checks between canonical JSON and rendered output.
- Source and citation-link checks.
- Keyboard-navigation and basic accessibility checks.
- Responsive layout checks at defined viewport sizes.
- Screenshot generation and visual regression comparison.

The visual QA skill receives canonical data, deterministic reports, rendered HTML, and desktop and mobile screenshots. It checks hierarchy, readability, overflow, clipping, missing sections, misplaced content, and obvious inconsistencies. It cannot override a failed deterministic check.

### Publication and CI

LLM generation and validation run locally. Successful reports are committed with their exact input hashes.

CI does not rerun expensive models. It verifies schemas, source immutability, checksums, report signatures and hashes, deterministic output, tests, and the absence of stale downstream artifacts.

Publication requires verified provenance, approved extraction fidelity, validated canonical data, successful deterministic rendering, and successful site validation. Merges to the main branch generate the static distribution and deploy it to GitHub Pages.

### Language and licensing

Public summaries, methodology, status labels, and website content are Spanish. Code, comments, schemas, prompts, tests, Make targets, architecture, and commit messages are English.

Project code and original documentation use the MIT licence. Archived third-party political documents retain their publishers' rights and are not relicensed.

### Assessment separation

Factual programme summaries remain separate from deliverability estimates, historical reliability, observed fulfilment, and delivery confidence.

Any future score must expose its formula, evidence, uncertainty, and methodology version. No unexplained political ranking is permitted.

## Testing Decisions

Tests will exercise module interfaces and observable artifact contracts rather than internal implementation details.

The existing Frente Obrero and VOX Castilla y León documents are the initial golden fixtures.

### Submission and provenance tests

- Accept a valid PDF candidate and reject invalid MIME content.
- Enforce size, page, timeout, and memory limits.
- Detect duplicate checksums.
- Verify a direct official-domain source.
- Verify an official page linking to an approved CDN URL.
- Reject a candidate whose bytes differ from the independently downloaded canonical source.
- Keep probable and unverified sources out of publication.
- Preserve redirect and retrieval evidence.

### Source registry tests

- Refuse to overwrite committed source bytes.
- Create explicit revisions for corrected documents.
- Preserve original filenames, URLs, checksums, and revision relationships.
- Detect stale metadata when source hashes change.

### Extraction tests

- Produce page-aware structured output for both golden PDFs.
- Preserve known headings, numbers, names, and negations.
- Record page-level warnings.
- Keep OCR disabled for native-text fixtures unless explicitly requested.
- Trigger visual handling for textless or suspicious pages.
- Reproduce identical extraction with pinned versions and unchanged inputs.

### Extraction audit tests

- Detect an omitted page.
- Detect altered numbers and negations.
- Detect incorrect page mapping.
- Detect inserted unsupported text.
- Demonstrate that blockers override a high aggregate score.
- Validate the versioned score formula against calibrated fixtures.

### Canonical model tests

- Validate programme JSON against its schema.
- Require citations for public measures.
- Distinguish `mentioned`, `not-covered`, `not-applicable`, and `unclear`.
- Apply jurisdiction rules for national and autonomous programmes.
- Detect unknown or stale taxonomy versions.

### Validator contract tests

- Store pass, fail, disagreement, and insufficient-evidence fixtures.
- Reject unsupported claims.
- Reject wrong polarity, scope, quantities, or responsible institutions.
- Detect a major omitted section.
- Detect duplicate or contradictory normalized measures.
- Invalidate reports when any recorded input hash changes.

LLM wording is not tested directly. Tests validate the schemas, decisions, evidence references, and pipeline behavior surrounding model output.

### Rendering tests

- Generate Spanish Markdown and HTML from the same canonical fixture.
- Render every taxonomy topic in the defined order.
- Render distinct public labels for uncovered and non-applicable topics.
- Preserve every measure and citation in both outputs.
- Escape untrusted source text correctly.
- Produce deterministic output for unchanged input.

### Site quality tests

- Parse generated HTML successfully.
- Verify required landmarks and keyboard navigation.
- Verify links and source references.
- Compare canonical content with rendered DOM content.
- Capture desktop and mobile screenshots.
- Detect overflow, clipping, missing sections, and broken comparison layouts.
- Keep visual LLM approval from overriding deterministic failures.

### End-to-end tests

- Process each golden PDF from verified source metadata to validated canonical data and generated site output.
- Keep failed provenance, extraction, validation, or rendering artifacts out of the public catalogue.
- Verify CI can validate committed reports without invoking an LLM.

## Out of Scope

- A public web upload form, user accounts, or runtime backend in the first version.
- A database or CMS.
- Treating contributor-uploaded bytes as official without independent provenance.
- Publication of probable or unverified sources.
- Human review as a required publication gate.
- LLM-generated HTML, CSS, JavaScript, or templates.
- Semantic search or embeddings in the first version.
- Automatic verification of every factual statistic asserted by a political party.
- Multiple supporting documents per programme in the first version.
- A final deliverability or fulfilment formula before a dedicated methodology specification.
- An unexplained composite political score.
- Runtime regeneration of summaries or pages on the public website.

## Further Notes

The hardest unresolved problem is the provenance trust anchor. Official-domain and referral evidence can automate most established-party submissions, but new parties, campaign microsites, removed pages, social-media-only publication, and third-party document hosts require explicit evidence rules.

Docling 2.145.0 successfully converted both initial PDFs during evaluation, but produced OCR warnings and used substantial resources. The 64-page document completed in about one minute with approximately 4.3 GB peak memory. The 172-page document completed in about seven minutes with approximately 7.5 GB peak memory. Converter success must therefore remain separate from extraction approval.

LLM-only validation is probabilistic. Independent contexts, original-page evidence, deterministic checks, and disagreement blocking reduce correlated errors but do not prove infallibility. The public methodology must state this limitation clearly.

A nationwide archive may eventually exceed practical Git storage limits. Source checksums and provenance metadata must remain stable if binary storage later moves to Git LFS, releases, or object storage.

After this specification is accepted, Taskqueue should be initialized and implementation work should be split into small tasks that reference the relevant specification sections. Skills should be authored only after their input and output schemas are defined.
