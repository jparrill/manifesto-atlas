# Product-scope and premature-platform audit

Task: TQ-0060. Audited implementation baseline: `94c4cad`.
This audit recommends scope decisions; it does not authorize or apply deletions,
change the accepted specification, or relax an active trust boundary.

## Executive judgment

The project has built foundations for an authenticated, hostile-input processing
platform before delivering its first useful programme summary. Many individual
checks follow logically from that platform's contracts. That does not establish
that the platform is the smallest way to deliver the product.

The valuable product is an archive of original political documents, Spanish
topic-based summaries, source-page citations, independent claim/coverage checks,
and static publication. Recording provenance, exact input hashes and model/tool
versions supports that product. A generic artifact DAG, ten report-signing roles,
publisher-key lifecycle, arbitrary-URL transport qualification and mandatory
pixel-exact browser profiles are not prerequisites inherent in that idea.

The previous trust-boundary audit asked whether implementations satisfied their
contracts. This audit asks whether those contracts deserve to be in the first
release. A repair can be locally correct while maintaining the wrong scope.

## What history actually establishes

The earliest available Git commit, `e10e85a`, already contains the current
[README](../../README.md) and [SPEC](../SPEC.md), including 40 user stories,
independent downloading, isolated conversion, separate claim/coverage passes,
provider substitution, report signatures and browser/visual checks. These files
have no diff from that baseline through `94c4cad`:

```bash
git diff --exit-code e10e85a 94c4cad -- README.md docs/SPEC.md
```

There is no earlier committed minimal specification to compare. The audit cannot
attribute every expanded requirement to later agents or reconstruct unrecorded
conversation. The earliest product narrative is simple; the earliest technical
specification already expands it substantially. Subsequent architecture adds
specific mechanisms and makes almost all of them mandatory before the first
useful public result.

Examples of later choices rather than intrinsic product requirements:

- [Artifact contract](../design/artifact-contract.md): a 26-kind content-addressed
  DAG, a separate object/manifest store, schema/tool/configuration leaves at every
  stage, and ten Ed25519 report-role mappings. The SPEC requires traceability,
  invalidation and authenticated reports, but does not prescribe this machinery.
- [Isolation contract](../design/submission-isolation.md): one exact rootless
  Podman/crun/SELinux/systemd/cgroup profile, disposable qualification probes and
  pinned OCI-image prerequisites. Safe bounded conversion is a real requirement;
  this particular platform and qualification apparatus is a design choice.
- [Provenance contract](../design/provenance-trust.md): 180-day domain grants,
  key grants, precise SSRF address policy and a custom bounded transport. This is
  a large trust-root/network subsystem beyond recording an official source.
- [Site contract](../design/site-quality.md): a precise Chromium/driver/font
  profile and a 0.1% changed-pixel threshold selected before an implemented site
  or measured screenshot corpus exists. Readability, accessibility and content
  parity are justified; this release gate needs an actual demonstrated need.

## Measured state

Counts are raw lines, including comments and blank lines, at `94c4cad`. They are
inventory measurements, not estimates of programmer effort or software quality.

| Area | Measurement |
| --- | --- |
| Tracked Go source | 19 files, 7,048 lines |
| Application non-test Go | 3,938 lines, including 182 for CLI/entry point |
| Go tests | 2,897 lines |
| Existing audit probe | 213 lines, counted separately from application code |
| Approved design documents | 9 files, 1,968 lines |
| Tracked task records, excluding generated task guide | 59 files, 5,475 lines before this audit record |
| External Go module dependencies | 0 |
| Working public pipeline stages | 0 of 12 in both Make and CLI |
| Tracked original PDFs / project skills / machine-readable schema JSON | 0 / 0 / 0 |
| Tracked presentation assets / GitHub workflow files | 0 / 0 |

The queue during this audit contains 60 tasks including TQ-0060, 143 dependency
edges, and a longest native path of 29 task nodes. This includes completed design
work; it is not a claim that 29 future tasks remain sequentially executable.

The ordinary executable import graph is only:

```text
cmd/manifesto-atlas -> internal/cli -> Go standard library
```

All six domain packages are absent from that executable's dependency closure.
The artifact parser/hash helpers do have real inter-library reuse by canonical,
semantic, attestation and provenance. That distinction matters: disconnected
product integration does not make every internal helper dead code.

Normal Go imports among the unfinished libraries are:

```text
attestation -> artifact
canonical   -> artifact
provenance  -> artifact
semantic    -> artifact
taxonomy    -> standard library only
```

Tests and the audit probe exercise these libraries; no public command processes
a PDF, calls a configured model or renders a summary. `make verify` passing is
therefore meaningful library regression evidence, not evidence that the product
works. Removing useful tests is not the recommended simplification.

## IPv6: a symptom of the selected boundary

`IsPublicIP` and `ResolvePublicIP` exist because `FetchEvidence` is intended to
download URLs supplied through an untrusted contribution path. A malicious URL,
redirect or changed DNS answer can point to loopback/internal services or cloud
metadata instead of the public PDF. Checking only IPv4 would leave the analogous
IPv6 route available. Those checks are legitimate within that threat model.

The product does not intrinsically need an IP-policy engine. The upstream scope
decision was to build a reusable arbitrary-URL downloader before processing one
maintainer-selected source. That is what introduced DNS pinning, reserved-range
tables, TLS/redirect policy, rebinding tests and the IPv6 discussion.

Recommended first-release boundary: the trusted maintainer selects and acquires
the initial official documents and records their exact source evidence and
checksums; automated content processing accepts those local registered inputs.
Contributor-provided files, URLs and reports do not automatically become trusted
or execute on a privileged CI runner. Automated publication still requires the
independent content checks, not a manual bypass of failed claims.

This deliberately defers automated arbitrary-contributor source admission. It is
a scope change to SPEC's independent-download workflow, not an equivalent
implementation of it. Keep origin labels honest: a selected source is not proof
that an unimplemented automated provenance chain ran. If arbitrary-URL admission
stays in scope, retain and finish IPv4/IPv6/DNS/TLS protections; do not just delete
IPv6 checks or assume that invoking curl alone supplies an SSRF boundary.

## Ranked simplification findings

The ranking favors removing whole prerequisites to a useful product, rather than
minimizing isolated lines. Tags follow the ponytail-audit vocabulary.

1. **yagni:** Stop making the authenticated artifact platform the first vertical
   slice. Replace the 26-kind DAG/object-store contract with per-programme run
   manifests containing exact file/input hashes and a simple stage checklist;
   rerun that programme when its inputs change. Keep strict JSON/hash helpers.
   [Code](../../internal/artifact/artifact.go),
   [contract](../design/artifact-contract.md), TQ-0019/TQ-0050.
2. **delete:** Remove report-role PKI from the MVP, conditional on trusted local
   generation/admission of reports. Use the owner-controlled Git workflow plus
   stored input-bound outputs; do not accept arbitrary contributor `pass` files.
   Git alone does not authenticate model execution, and this is a deliberate
   change to the signed-report requirement. Reintroduce attestations when report
   producers actually cross separate authority boundaries.
   [Package](../../internal/attestation/attestation.go), TQ-0020.
3. **delete:** Defer the automatic source-admission platform: custom network
   transport, expiring party grants, publisher-key grants and semantic provenance
   approval. Start with explicitly selected local official sources and retained
   origin evidence. Add automation after actual source variation demonstrates
   the need. [Transport](../../internal/provenance/transport.go),
   [registry](../../internal/provenance/party.go), TQ-0022–TQ-0026.
4. **shrink:** Reduce converter isolation to one supported, documented bounded
   deployment, using existing OS/container mechanisms. Keep write/network
   restrictions, limits and whole-job cleanup; reconsider the mandatory custom
   platform qualification suite as a prerequisite to every real fixture.
   [Contract](../design/submission-isolation.md), TQ-0021.
5. **yagni:** Implement one actual model adapter and the generator/claim/coverage
   calls before a universal six-role skill execution framework. Keep clean
   conversations, strict output checks, complete evidence coverage and input
   hashes. Defer publisher/visual roles, general adapter qualification and an
   all-component identity graph until a real second integration requires them.
   [Runner](../../internal/semantic/semantic.go), TQ-0024/TQ-0025/TQ-0048.
6. **shrink:** Keep topic IDs, factual metadata, page/excerpt citations and the
   five states, but defer a general legal-competence registry and split/merge
   migration engine. Use only evidenced document-specific scope; ambiguous
   applicability stays unclear, never a guessed exclusion. Prefer local stable
   IDs scoped to a source revision over hashing every citation/measure into a
   new identity protocol unless actual revision matching needs it.
   [Taxonomy](../../internal/taxonomy/taxonomy.go),
   [canonical data](../../internal/canonical/canonical.go), TQ-0036/TQ-0037.
7. **yagni:** Defer a calibrated numeric extraction-fidelity score and mandatory
   pixel-regression/visual-model approval from the first publishable slice.
   Keep original-page evidence, warnings, critical extraction blockers,
   independent semantic validation, content parity and accessibility basics.
   Start by measuring actual documents/pages and testing real output.
   [Extraction contract](../design/extraction-contract.md),
   [site contract](../design/site-quality.md), TQ-0034/TQ-0047/TQ-0048.
8. **shrink:** Replan around product milestones, not all horizontal layers.
   Recover a real source before its registration platform is finished; do not
   make basic site checks depend on regional comparison support or local preview
   depend on final publication approval. A draft preview must be explicit and
   excluded from the public catalogue. Preserve historical task records and
   owner-only inbox control. [Task graph](../../.tasks/AGENTS.md),
   TQ-0028/TQ-0030/TQ-0045/TQ-0046/TQ-0051.
9. **reuse:** Replace duplicate broad Make regression invocations with one
   repository-wide `go test ./...` while retaining focused targets for debugging.
   `test-party-registry` and `test-fetch` currently run the same package suite;
   Go caching reduces execution cost but not the duplicated gate bookkeeping.
   [Makefile](../../Makefile).

## Quantified cuts and reference checks

Repository-wide searches included Go, tests, task bodies, documentation and
string references, not only import statements. No runtime/dynamic dispatch of
the proposed platform packages exists in the CLI. Documentation references are
contracts to revise, not hidden production consumers.

| Non-overlapping candidate | Existing production lines | Existing test lines | Gross removable Go |
| --- | ---: | ---: | ---: |
| Report-attestation package, if finding 2 is adopted | 479 | 363 | 842 |
| Provenance transport and its tests, if finding 3 is adopted | 370 | 269 | 639 |
| Party-grant registry and its tests, if finding 3 is adopted | 396 | 232 | 628 |
| **Direct package-removal subtotal** | **1,245** | **864** | **2,109** |

Additional disjoint candidate: 460 existing production lines in
`artifact.go:457–558` and `artifact.go:573–930` implement manifest/kind/report/graph
machinery rather than the shared strict parser/hash helper. Removing them needs
test/helper separation and removal of related manifest-ID callers; no test-line
or net savings estimate is claimed for this partial-file refactor. The entire
1,587-line artifact package must not also be counted as removable.

The 976 existing semantic implementation/test lines and 1,757 combined canonical/
taxonomy lines are **review surfaces**, not deletion savings: useful validation
must survive in the smaller product. The report proposes no new dependency and
claims zero removable external Go dependencies because none are installed.

The 2,109-line subtotal is gross existing code/tests retired under explicit scope
changes, approximately 30% of the current 7,048 Go lines. New MVP glue and
metadata handling, audit-probe adaptation, Make/docs/task changes are not
estimated. Do not present this as a measured final net patch or add planned
browser/container dependencies to today's installed dependency count.

Reference findings:

- Attestation symbols have definitions/tests and design/Make/task references;
  no application package or audit probe imports attestation.
- Provenance entry points are called by their package tests and the synthetic
  audit probe, not the product executable.
- Graph construction/invalidation/report-kind helpers are used by artifact's
  own fixtures. Other libraries use only shared JSON/hash functions, with
  `ManifestID` additionally exercised by artifact tests.
- The semantic runner is used by its fakes and audit probe; no configured
  production provider exists. Retain the repaired rejection logic when wiring
  one real adapter rather than treating fixes as sunk-cost justification for
  the entire framework.

## Backlog dependencies that delay product feedback

- TQ-0028's longest path is TQ-0008 → TQ-0010 → TQ-0011 → TQ-0013 → TQ-0014 →
  TQ-0015 → TQ-0024 → TQ-0025 → TQ-0026 → TQ-0027 → TQ-0028. Thus recovering a
  first accepted golden source currently waits on taxonomy/canonical designs
  and the generic semantic runner. Obtaining immutable candidate bytes and
  origin evidence should precede that platform, without prematurely declaring
  source registration or provenance approval complete.
- TQ-0030 waits for **both** golden-source registration tasks before extraction
  implementation completes. One-source development feedback should come first;
  retain both-document regression as a later acceptance milestone.
- TQ-0046 depends on TQ-0045. This makes regional comparison implementation a
  prerequisite of basic site parity/link checks and consequently publication
  of even one programme. Regional comparison is valuable but not foundational.
- TQ-0051 joins process orchestration and preview and waits for TQ-0050. A safe
  private fixture/draft preview can provide early feedback without publishing
  unapproved content; it need not wait for the entire publication graph.

The longest path to golden E2E is 29 nodes. Cyclic-dependency validation passes,
but an acyclic graph can still encode an unhelpful implementation order.

## Recommended first useful release

Use a trusted local-maintainer workflow, one real configured provider and one
document first. Keep contributors and their reports outside the approval boundary.

1. **Original source:** recover one approved PDF, preserve its bytes and recorded
   official URL/evidence, checksum, party, scope and document identity. Do not
   invent the absent golden sources or substitute another election/revision.
2. **Page evidence:** run one actual extractor in a bounded, isolated deployment;
   retain page text, warnings, original-page references and needed page images.
   Block on missing/unreadable material rather than award a reassuring score.
3. **Canonical summary:** generate Spanish topic-based JSON with atomic measures,
   conditions and page/excerpt citations. Validate actual schema/references;
   retain the topic registry and do not infer opposition from absence.
4. **Independent checks:** run separate clean claim and whole-source coverage
   jobs against the original evidence. Persist actual prompt/model/input
   identities and original reports. Failures, missing coverage and disagreement
   block public eligibility. No manual publication override is proposed.
5. **Useful output:** render a readable Spanish Markdown summary first. Integrate
   human/non-LLM-authored HTML assets for the static site, preserving the existing
   authorship restriction, contextual escaping, citation links, parity and basic
   accessibility. Fixture/draft previews remain distinct from publication.
6. **End-to-end acceptance:** demonstrate one full real programme with deliberate
   missing-citation, changed-input and unsupported-claim failures; then repeat
   with the second golden source. Add regional comparisons after two compatible
   programmes exist. Automate broader admission when recurring use justifies it.

A small programme run manifest can record source/extraction/draft/prompt/schema
and tool/model identities, the claim/coverage report hashes and final eligibility.
Recompute the whole affected programme initially instead of optimizing precise
transitive cache invalidation. This preserves traceability while avoiding the
general content-addressed workflow platform.

## What should remain non-negotiable

- Immutable source bytes, correct source/election/scope identity and provenance
  evidence without overstating what was automatically verified.
- Page mapping, original excerpts, quantities/negations/conditions and warnings.
- Structured canonical data as public-content source, Spanish output and
  deterministic Markdown/HTML generation from that data.
- Separate claim and complete-source coverage passes, input-bound reports and
  refusal to publish when evidence/validation is insufficient.
- JSON/input bounds, cancellation, safe file handling, converter isolation,
  credential separation, output escaping and accessibility basics.
- Explicit scope approval before implementing the recommendations. Deleting
  protection while keeping the same exposed hostile-input workflow is not a
  simplification of the product; it is a weaker implementation.

## Evidence commands

```bash
git diff --exit-code e10e85a 94c4cad -- README.md docs/SPEC.md
git ls-files '*.go' | xargs wc -l
git ls-files 'docs/design/*.md' | xargs wc -l
go list -deps ./cmd/manifesto-atlas
go list -json ./...
go list -m all
tq list --json
go run ./cmd/manifesto-atlas help
```

All twelve CLI stages were invoked with a fixture argument; all returned
`not available yet` and nonzero. The twelve corresponding Make stages also
failed nonzero. No source/model/network operation was executed. Reference
searches included `VerifyAttestation`, `Sign`, `FetchEvidence`, `CheckAnchor`,
`NewGraph`, `StaleIDs`, `ReportResult`, `ProviderAdapter`, `NewRunner`, schema
names and contract strings across the repository, including tests and task docs.

The first step after scope selection should be a revised short specification
and a handful of independently verifiable vertical outcomes, not another batch
of infrastructure-only repairs executed automatically from the old queue.
