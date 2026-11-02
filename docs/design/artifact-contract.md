# Artifact identity, invalidation, and attestation contract

> Historical contract. For active MVP work, [MVP.md](MVP.md) and
> [SPEC.md](../SPEC.md) supersede conflicting requirements. TQ-0063 retires the
> inactive DAG/PKI platform while retaining useful strict JSON/hash helpers.

Contract: `atlas.artifact/v1`. Status: implementation contract for TQ-0019 and
TQ-0020, established by TQ-0008 under Juan's architecture assignment.

## Decisions and scope

Use a content-addressed directed acyclic graph (DAG), SHA-256, strict JSON,
and detached Ed25519 attestations. Verification is offline and fail-closed.
The graph records evidence and exact computation inputs; it does not prove that
a political claim is correct. A valid signature never overrides a failed check.

This implements the identity and authority requirements of the
[Make targets](../SPEC.md#make-targets), [LLM validation](../SPEC.md#llm-validation),
and [publication and CI](../SPEC.md#publication-and-ci) specification sections.
It does not change those product decisions. Stage-specific payload schemas,
source storage hierarchy, runner isolation, and presentation contracts remain
owned by their dependent design tasks. Their absence blocks those stages; it
does not block implementation of this common envelope and graph validator.

## Serialization and hash boundaries

`H(bytes)` means `sha256:` followed by 64 lowercase hexadecimal digits. There
is no alternative algorithm, uppercase form, abbreviated hash, or remote ID in
v1. A manifest ID is `H(C(manifest))`, not the hash of its filename or Git blob.

`C` is the following versioned canonical JSON encoding:

1. Accept exactly one UTF-8 JSON value, with no BOM, trailing value, duplicate
   object member, invalid Unicode, or unpaired surrogate. Reject before decoding
   loses duplicate keys or replaces invalid characters.
2. Values are objects, arrays, strings, booleans, null, or integers in
   `[-9007199254740991, 9007199254740991]`. Reject fractional/exponent numeric
   tokens and negative zero. Scores requiring fractions use schema-defined
   integer units or decimal strings, never floating-point JSON numbers.
3. Recursively sort object keys by their UTF-8 byte order; preserve array order
   and string code points, including whitespace. Do not normalize Unicode,
   trim excerpts, reorder measures, or translate text.
4. Encode without indentation, insignificant whitespace, or trailing newline.
   Strings use Go `encoding/json` escaping with `SetEscapeHTML(false)`:
   quotes/backslashes and control characters are escaped; U+2028 and U+2029
   remain escaped as Go emits them. `/`, `<`, `>`, and `&` are not escaped.
   For valid Unicode, UTF-8 key order agrees with Go string order. Remove only
   the encoder's final newline, not any whitespace inside a string.

This restricted profile is deliberately not advertised as RFC 8785/JCS.
Changing these rules requires a new contract version. Implementations must
validate before canonicalization and use integer-preserving decoding, not a
`float64` intermediate. Unknown manifest fields and unknown contract/kind
versions fail; they must not be silently discarded before hashing.

Payloads declared JSON are stored as `C(payload)` and hashed as those exact
stored bytes. PDFs, extracted text, Markdown, images, fonts, HTML, templates,
model weights, and other non-JSON files are hashed byte-for-byte. No newline,
metadata, PDF, image, or encoding rewrite is permitted during verification.
Converters may produce new derived artifacts; they never rewrite source bytes.

### Reference vectors

| Input | Canonical bytes or expected rejection |
| --- | --- |
| `{"b":2,"a":1}` | `{"a":1,"b":2}` |
| `{"a":1, "b":2}` | Same bytes and ID as the preceding row |
| `{"quote":"<salud>&"}` | `{"quote":"<salud>&"}` |
| `{"a":1,"a":2}` | Reject duplicate member |
| `{"score":0.5}` or `{"n":1e0}` or `-0` | Reject numeric token |
| `{"n":9007199254740992}` | Reject integer outside the v1 range |
| Two JSON values, invalid UTF-8, or an unpaired surrogate | Reject malformed input |

## Manifest and storage interface

A manifest has exactly these required fields:

```json
{
  "contract": "atlas.artifact/v1",
  "kind": "source",
  "payload": {
    "media_type": "application/pdf",
    "sha256": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "size": 0
  },
  "inputs": {}
}
```

This is a worked envelope for the empty-byte SHA-256 vector, **not** an accepted
PDF: the submission/domain validator must reject its empty payload. Envelope
validity is necessary but insufficient for stage acceptance.

- `kind` is one of the v1 kinds below. No default or inferred kind is allowed.
- `payload` has exactly `media_type` (nonempty MIME type without parameters),
  `sha256` (a full `H` ID), and `size` (nonnegative integer). The verifier checks
  size and hash against the stored bytes before parsing them.
- `inputs` maps unique, named input roles to full manifest IDs. Slot names are
  lowercase ASCII letters/digits/hyphens, start with a letter, and are at most
  64 bytes. An edge binds a manifest, not merely its payload hash.
- IDs are computed, not self-embedded in manifests. A supplied ID/index must
  match the computed ID. Reject repeated IDs in a submitted graph, even if the
  manifest bytes agree. A store may reuse one existing immutable object.

The Go-facing interface is conceptually `Canonicalize(bytes) -> bytes`,
`ManifestID(manifest) -> ID`, `ValidateGraph(manifests, selectedRoots) -> error`,
and `ValidatePayload(manifest, bytes) -> error`. Errors include a stable code,
the subject ID/slot when available, and no private data. No API implicitly
retrieves files, follows URLs, invokes a model, or approves domain content.
Graph validation requires **all** references in the supplied graph to resolve
and all nodes to be acyclic, including disconnected nodes.

Resolve payloads/manifests only through a caller-supplied local store, under
`objects/sha256/<64hex>` and `manifests/sha256/<64hex>.json`. These are derived
paths, never input path fields. Reject absolute paths, separators in IDs,
traversal, symlinks, non-regular files, external references, and unknown roots;
never fall back to a URL, another checkout, or arbitrary filesystem lookup.
Use rooted/no-follow filesystem operations (Go `os.Root` or equivalent) and
check the opened file, rather than a check-then-open path sequence.

Limits: 1 MiB per manifest, 16 MiB per JSON payload, nesting depth 64 (root is
depth 1), 1,024 edges per node, 10,000 nodes and 64 MiB total manifest input per
graph. Reject excess before allocation/traversal; use iterative traversal or a
bounded stack. Binary size budgets belong to submission/extraction policy.
Limits cannot be disabled by values inside an untrusted manifest.

Write verified bytes to private scratch, then install immutable objects
atomically without replacing an existing different object. Compare an existing
object before reuse. An interrupted write never creates an accepted manifest or
approval. The store is separate from the year-first public source hierarchy;
the source registry can reference these IDs without relocating archived PDFs.

## Artifact kinds and required graph edges

Leaf kinds are `source`, `schema`, `taxonomy`, `prompt`, `model`, `tool`,
`configuration`, `templates`, and `trust-policy`. They have no inputs. A `source`
is candidate/canonical PDF **bytes**, not an assertion of official provenance.
Each other leaf payload is a versioned resource description (or a binary/bundle
descriptor): it includes a stable name/version and exact immutable component
digests. Mutable aliases such as `latest` or an unversioned model name cannot
establish a publishable identity. Bundles enumerate sorted safe relative paths,
sizes, and digests for every component; no symlinks or external components.

For every non-leaf, require `schema`, `tool`, and `configuration` inputs of
those leaf kinds. A tool description binds executable/container digest,
version, dependency lock, platform/runtime, and relevant execution settings.
Configuration binds actual options, seeds, limits, and policy versions, never
credentials. The schema input binds exact payload schema bytes/version; a
missing or unsupported domain schema cannot be treated as a generic pass.

The table lists **additional** required inputs and kinds. No optional or
undeclared input slots are permitted in v1.
A missing required edge fails even if the payload copies the missing hash.
Future optional inputs require an explicit versioned kind/slot rule.

| Kind | Additional input slots (`slot: kind`) |
| --- | --- |
| `provenance` | `submitted-source: source`, `canonical-source: source`, `trust-policy: trust-policy`, `evidence: configuration`, `prompt: prompt`, `model: model` |
| `source-record` | `source: source`, `provenance: provenance` |
| `extraction` | `source-record: source-record`, `extraction-model: model` |
| `secondary-evidence` | `source-record: source-record` |
| `page-images` | `source-record: source-record` |
| `extraction-audit` | `extraction: extraction`, `secondary-evidence: secondary-evidence`, `page-images: page-images`, `prompt: prompt`, `model: model` |
| `extraction-approval` | `audit: extraction-audit`, `extraction: extraction`, `formula: configuration` |
| `programme` | `extraction-approval: extraction-approval`, `extraction: extraction`, `taxonomy: taxonomy`, `jurisdiction: configuration`, `prompt: prompt`, `model: model` |
| `claim-report` | `programme: programme`, `extraction: extraction`, `source-record: source-record`, `page-images: page-images`, `prompt: prompt`, `model: model` |
| `coverage-report` | `programme: programme`, `extraction: extraction`, `source-record: source-record`, `page-images: page-images`, `prompt: prompt`, `model: model` |
| `content-approval` | `programme: programme`, `claim-report: claim-report`, `coverage-report: coverage-report`, `extraction-approval: extraction-approval` |
| `markdown` | `programme: programme`, `content-approval: content-approval` |
| `html` | `programme: programme`, `content-approval: content-approval`, `templates: templates` |
| `site-deterministic-report` | `programme: programme`, `markdown: markdown`, `html: html`, `qa: configuration` |
| `site-visual-report` | `programme: programme`, `html: html`, `deterministic-report: site-deterministic-report`, `prompt: prompt`, `model: model`, `qa: configuration` |
| `site-approval` | `deterministic-report: site-deterministic-report`, `visual-report: site-visual-report` |
| `publication` | `source-record: source-record`, `extraction-approval: extraction-approval`, `content-approval: content-approval`, `markdown: markdown`, `html: html`, `site-approval: site-approval`, `trust-policy: trust-policy` |

`page-images` may have an empty, schema-valid image inventory when the approved
page-handling policy finds no page needing an image; an absent inventory or an
unresolved suspicious page is not equivalent. All content validators receive
original evidence through the source-record and image edges. The secondary
extractor's output is corroborating evidence, never a substitute for extraction
or canonical programme JSON. Semantic provenance evidence is collected by a
bounded transport and passed as data; the model itself receives no network.

In addition to types/edges, stage gates enforce cross-edge equality: claim and
coverage reports must name the same programme, extraction, and source as their
approval root; renderings and site reports must name that programme and exact
renderings; source-record's source must equal provenance's canonical-source;
extraction-approval must bind its audit's extraction. A well-typed graph with
mixed revisions fails `input-mismatch`, not a warning.

## Worked manifests across the pipeline

The following notation defines complete worked manifests without pretending to
know real PDF/model checksums. It is an **executable mathematical recipe**, not
literal JSON with fake IDs: `M(kind, payloadBytes, inputs)` constructs the exact
four-field envelope above, computes payload size/hash, and returns its manifest
ID. To instantiate the example without external files, each named payload is
the UTF-8 bytes of its label (for example `candidatePDF`); use `text/plain` for
these synthetic payloads. They exercise envelope identities and edges, not
domain acceptance. The implementation task replaces these bytes with its own
positive/negative domain fixtures before claiming payload-schema approval.
For a matching A/B example, use the same label bytes for both. `base(K)` expands
to all three required edges:
`schema=SK`, `tool=TK`, `configuration=CK` (distinct leaf manifests for kind K).
Each `SK`, `TK`, and `CK` is produced with `M` and empty inputs, using the
UTF-8 label `schema:K`, `tool:K`, or `configuration:K` as payload. Other leaves
use their named label bytes in the same way. `+` merges disjoint named slots;
it never silently overwrites one. Media type is an explicit argument fixed to
`text/plain` for this worked example; real stages use their actual MIME type.

| Name / manifest | Complete input map in addition to `base(kind)` |
| --- | --- |
| `A = M(source, candidatePDF, {})` | Leaf; no base edges |
| `B = M(source, canonicalPDF, {})` | Leaf; in a matching submission `A = B`, stored once |
| `R = M(trust-policy, officialDomainEvidence, {})` | Leaf; official domains and evidence, not report-signing authority |
| `E = M(configuration, retrievalEvidenceInventory, {})` | Leaf; exact referral/redirect/response evidence digests |
| `Q = M(prompt, provenancePrompt, {})`, `L = M(model, provenanceModelDescriptor, {})` | Leaves; pinned identifiers and bytes/component digests |
| `P = M(provenance, verifiedReport, base(provenance) + {...})` | `submitted-source=A, canonical-source=B, trust-policy=R, evidence=E, prompt=Q, model=L` |
| `S = M(source-record, immutableRegistration, base(source-record) + {...})` | `source=B, provenance=P` |
| `X = M(extraction, pageMappedJSON, base(extraction) + {...})` | `source-record=S, extraction-model=MX` (pinned Docling model bundle leaf) |
| `Y = M(secondary-evidence, secondaryText, base(secondary-evidence) + {...})` | `source-record=S` |
| `I = M(page-images, originalPageImageInventory, base(page-images) + {...})` | `source-record=S` |
| `U = M(extraction-audit, auditReport, base(extraction-audit) + {...})` | `extraction=X, secondary-evidence=Y, page-images=I, prompt=QU, model=MU` |
| `V = M(extraction-approval, extractionDecision, base(extraction-approval) + {...})` | `audit=U, extraction=X, formula=F` (calibrated formula configuration leaf) |
| `D = M(programme, canonicalDraft, base(programme) + {...})` | `extraction-approval=V, extraction=X, taxonomy=N, jurisdiction=J, prompt=QD, model=MD` |
| `C = M(claim-report, claimReport, base(claim-report) + {...})` | `programme=D, extraction=X, source-record=S, page-images=I, prompt=QC, model=MC` |
| `O = M(coverage-report, coverageReport, base(coverage-report) + {...})` | `programme=D, extraction=X, source-record=S, page-images=I, prompt=QO, model=MO` |
| `G = M(content-approval, contentDecision, base(content-approval) + {...})` | `programme=D, claim-report=C, coverage-report=O, extraction-approval=V` |
| `W = M(markdown, spanishMarkdown, base(markdown) + {...})` | `programme=D, content-approval=G` |
| `H = M(html, suppliedTemplateRendering, base(html) + {...})` | `programme=D, content-approval=G, templates=HT` (human-authored asset bundle leaf) |
| `K = M(site-deterministic-report, deterministicQA, base(site-deterministic-report) + {...})` | `programme=D, markdown=W, html=H, qa=Z` |
| `ZV = M(site-visual-report, visualQA, base(site-visual-report) + {...})` | `programme=D, html=H, deterministic-report=K, prompt=QV, model=MV, qa=Z` |
| `GS = M(site-approval, siteDecision, base(site-approval) + {...})` | `deterministic-report=K, visual-report=ZV` |
| `PUB = M(publication, distributionInventory, base(publication) + {...})` | `source-record=S, extraction-approval=V, content-approval=G, markdown=W, html=H, site-approval=GS, trust-policy=R` |

`N`, `J`, `HT`, `Z`, and every `Q*`/`M*` are leaves of the indicated taxonomy,
configuration, templates, prompt, or model kind. Tool/schema/configuration
leaves are separate per stage; reusing identical leaf bytes is allowed. No
attestation is embedded in D, C, or O, avoiding a self-reference. Fixtures must
deduplicate A/B when their IDs agree. Synthetic fixtures cannot serve as real
official provenance or empirical fidelity calibration evidence.

## Time, state, reproducibility, and invalidation

Manifests contain no execution timestamp. Report payloads contain their outcome
and evidence; the attestation supplies the evaluation timestamp required by the
Spec. Non-semantic run bookkeeping (`started_at`, `finished_at`, elapsed time,
temporary paths, request IDs) lives outside the hashed graph in local run logs.
Changing such bookkeeping alone does not invalidate content or manufacture a
fresh evaluation. Retrieval time and response headers **are evidence**, live
inside the retrieval evidence resource, and therefore change its ID.

There are two distinct reproducibility claims:

- Deterministic stages: identical pinned inputs must produce identical payload
  bytes and manifest IDs. A mismatch fails `nondeterministic-output` and must be
  investigated; it must not overwrite the old object.
- Semantic/model stages: identical inputs need not reproduce wording or the
  same judgment. Reuse a stored, valid, authorized result if inputs agree. An
  explicit fresh evaluation creates a new report when its payload differs;
  retain all results. Conflicting current reports cannot be resolved by picking
  the convenient one: disagreement blocks publication until a new recorded
  evaluation set resolves it under the runner contract.

Stage execution states are `pending`, `running`, `succeeded`, `failed`, and
`insufficient`. They are run-log state, not mutable manifest fields. Report
payloads use the required `result` enum `pass`, `fail`, `insufficient`, `error`
and a `findings` array; stage schemas add exact decision/evidence fields.
Only `pass` plus the stage's domain predicates can advance a stage. Preserve
schema-valid failing reports as immutable artifacts. A malformed response is
retained as quarantined diagnostic bytes, never wrapped as a passing report.

`draft -> validating -> validated -> published` is derived programme lifecycle
state, not a rewrite of canonical draft D: validated means valid G for exactly
D; published means admitted PUB. The public view joins D with approvals and
validation metadata. Embedding G inside D would create a cycle and is forbidden.
Failed/insufficient reports and their signatures remain inspectable but are not
listed as approved programmes. No human-review fallback bypasses these gates.

A caller supplies selected roots for the **current** programme/configuration,
not just a historical publication ID. Compare expected current named inputs
against the manifests, then propagate changes along reverse edges. Recompute
only the transitive dependent set, not siblings or unrelated programmes. Old
artifacts stay valid historical bytes, but are `stale` for the current selection.
The verifier distinguishes `invalid` structure, `stale` inputs, `untrusted`
authority, and `blocked` domain result; all are non-publishable.

## Attestations and offline authority

Use Ed25519 from Go's standard library. Git commit signatures establish review
history for the checkout, not report authority. A contributor's own public key,
the PDF's signature, or a `trust-policy` graph leaf cannot grant report roles.

An attestation has exactly `contract`, `subject`, `role`, `key_id`, `issued_at`,
and `signature`. Contract is `atlas.attestation/v1`; subject is a manifest ID;
key ID is `H(raw 32-byte Ed25519 public key)`; signature is canonical padded
standard base64 of 64 bytes. Reject other key/signature sizes/encodings.
`issued_at` is UTC RFC3339 with whole seconds (`YYYY-MM-DDTHH:MM:SSZ`).
Sign the bytes `atlas.attestation/v1\n` followed by `C` of the other five fields
(excluding `signature`). The signature thus binds exact inputs transitively,
the role, signer, and timestamp. It is not a bare payload signature. New
attestation time changes the signed statement, not the artifact ID.

Required report-kind to role mapping:

| Kinds | Required signer role |
| --- | --- |
| `provenance` | `provenance` |
| `extraction-audit` | `extraction-audit` |
| `extraction-approval` | `extraction-gate` |
| `claim-report` | `claim-validator` |
| `coverage-report` | `coverage-validator` |
| `content-approval` | `content-gate` |
| `site-deterministic-report` | `site-deterministic` |
| `site-visual-report` | `site-visual` |
| `site-approval` | `site-gate` |
| `publication` | `publication-gate` |

Manifests of other kinds need no report attestation. Their bytes/inputs are
still verified. Every required report in a publication closure must have at
least one valid current attestation for its mapped role. Signing a failed
report authenticates failure, never approval. Sign only after schema/domain
checks; missing credentials or signing failure must leave no installed
approval envelope. The model process cannot access signing keys or mint roles.

The caller explicitly supplies a trusted local `atlas.trust/v1` bundle from
the owner-reviewed checkout/configuration, outside submitted evidence. It has
exact fields `contract` and `keys`; keys contain `key_id`, `public_key` (base64
raw key), `roles` (distinct allowed roles above), `not_before`, `not_after`, and
`revoked` (boolean). Reject duplicate keys/roles, mismatched computed key IDs,
unknown fields/versions, empty roles, invalid time ranges, and invalid encodings.
No private material is stored in this bundle. No automatic key discovery, key
server, TOFU, graph-supplied trust root, or certificate download is allowed.

Verify against an explicit trusted evaluation time: the issuance time must be
within the key's inclusive validity interval and not later than evaluation
time. The evaluation time must also be within that interval, and `revoked`
must be false. There is no grace period or guessed clock correction. Expired
keys require a fresh authorized attestation, not a model rerun if the evidence
and inputs are still current. Changing the trusted bundle requires rechecking
authority across the selected graph even though content IDs remain unchanged.

Rotation adds the new key with overlapping validity, then removes authority
from the old key after re-attestation. Revocation is retroactive for admission:
all signatures by that key fail even if dated before revocation. Retain them
as history; a new authorized key can re-attest unchanged evidence only after
re-verifying it. Retired/unknown/expired keys fail closed. CI's trusted checkout
revision is the authority boundary: offline verification cannot discover a
revocation it has not received. Never claim offline verification guarantees
fresh global revocation knowledge or prevents rollback of the trusted checkout.

Claim and coverage independence is a runner/job property recorded in their
schemas, not a requirement that two humans hold different signing keys. Keys
may have multiple explicit roles; role checks still prohibit substituting a
generator or claim attestation for a coverage attestation. Later runner
contracts must enforce clean contexts and cannot weaken this role mapping.

Envelope size is at most 16 KiB, trust bundle 1 MiB/256 keys, JSON depth at most
64. Signing APIs accept key bytes or an explicit private-key handle outside
the repository and runner sandbox. If a file handle is used, open a regular
non-symlink owner-only file in a caller-selected private root; never accept a
key path from a report or echo key bytes in logs. Signing credentials are not
artifact inputs. Verification requires public material only, no network.

## Scenario-to-expected-outcome table

In this table, descendants refer to the worked graph above. Include the changed
node itself for recomputation; retain the previous graph as historical evidence.

| Scenario | Expected outcome |
| --- | --- |
| Identical logical JSON, different object key order/indentation | Same canonical payload/manifest bytes and IDs; reuse current results |
| Execution-log timestamp changes only | No artifact ID or descendant changes; no new semantic approval implied |
| Attestation `issued_at` changes only | Artifact ID unchanged; new signature required; no content recomputation |
| Retrieval timestamp/evidence changes | E changes; P, S, X/Y/I and every dependent node become stale |
| Canonical PDF source changes B | P/S and their entire dependent closure become stale; no overwrite of B |
| Submitted candidate changes A only | P and its descendants become stale; unrelated source roots unchanged |
| Extraction bytes/model artifact MX/tool TX/options CX change | X and its descendants stale; P/S/Y/I remain current |
| Secondary extractor/tool changes | Y, U/V/D/C/O/G/W/H/K/ZV/GS/PUB stale; X remains current; Y never replaces X |
| Extraction schema SX changes | X and descendants stale; unrelated schemas and source bytes unchanged |
| Canonical programme schema SD changes | D and descendants stale; X/U/V remain current |
| Taxonomy N or jurisdiction J changes | D and descendants stale; provenance/extraction remain current |
| Normalizer prompt QD, model MD, tool TD, or configuration CD changes | D and descendants stale; sibling programmes not using that leaf unchanged |
| Claim prompt QC/model MC/configuration CC changes | C/G/W/H/K/ZV/GS/PUB stale; O and D remain current |
| Coverage prompt QO/model MO/configuration CO changes | O/G/W/H/K/ZV/GS/PUB stale; C and D remain current |
| Template/asset bundle HT or HTML renderer TH changes | H/K/ZV/GS/PUB stale; W, C/O/G remain current |
| QA configuration Z or deterministic QA tool TK changes | K/ZV/GS/PUB stale; W/H remain current |
| Visual QA prompt QV/model MV changes | ZV/GS/PUB stale; K/W/H remain current |
| Any other schema, tool, prompt, model, or config leaf changes | Exactly its reverse-reachable nodes stale; disconnected roots untouched |
| Required edge absent or unknown slot present | Reject `missing-input` / `unknown-input`; non-publishable even if a hash is copied into prose |
| Dangling ID, duplicate identity, self-edge, or graph cycle | Reject `missing-dependency`, `duplicate-id`, or `cycle`; no partial acceptance |
| Unknown contract/kind/schema/model version or mutable alias | Reject unsupported identity/version; no compatibility guessing |
| Manifest/payload size/hash mismatch or edited signed report | Reject integrity/signature verification; neither can advance a stage |
| Signed report for a different programme/source revision | Reject `input-mismatch`, even if signer is authorized |
| Forged signature, unknown signer, or graph-supplied key | Reject `untrusted`; graph content cannot grant authority |
| Claim signature reused as coverage signature | Reject subject/role mismatch; no role inference |
| Rotated key, both keys valid/authorized during overlap | Either signature can authenticate its exact mapped subject; old artifacts reused |
| Revoked/expired/removed key, even with old issuance timestamp | Reject admission until current authorized re-attestation; retain historical bytes |
| Future timestamp, unknown trust version, missing trust bundle | Reject authority verification; no network recovery or fallback |
| `fail`, `insufficient`, `error`, or contradictory current reports | Retain evidence; block publication irrespective of valid signatures |
| High fidelity score plus blocking finding | Block extraction approval and publication, never average away the blocker |
| Visual pass plus deterministic failure | Site approval blocked; visual report cannot override deterministic findings |
| Missing schema, original page evidence, or extraction approval | Block content/publication; no human-review shortcut |
| Path traversal, symlink, URL ID, oversized JSON, deep JSON | Reject before filesystem/network traversal or unbounded allocation |
| Interrupted write or signing without credential | No installed apparent approval; retain diagnostic failure safely |

## Versioning, alternatives, and hand-off

Contract/kind compatibility is explicit. v1 readers reject unknown versions;
migrations create new objects/manifests and invalidate dependent artifacts.
Never reinterpret old bytes under new rules. A stage schema release is a new
schema leaf even if it retains the same filename. Exact executable/model/image
and component digests, not only a marketing version, identify dependencies.
Unavailable model revision/digest evidence produces insufficient identity, not
a fabricated pin. CI can verify committed evidence without possessing model
weights only when their exact descriptor/digests and required runner evidence
are available under the later runner contract.

Rejected alternatives:

- Hashing Markdown or filenames: loses structure and exact source identity.
- Hashing raw JSON formatting: creates needless cache misses for key ordering.
- Floating-point canonicalization or a general JCS library: unnecessary for the
  bounded v1 envelope; integer/decimal-string payload schemas avoid ambiguity.
- Putting timestamps, state, approvals, or signatures inside the draft hash:
  breaks idempotency or creates cycles.
- Trusting a Git signature or self-supplied report key: lacks per-role authority.
- Online certificate checks in ordinary regression: violates offline checking.
- Re-running expensive models in CI: validates a different execution rather
  than the exact stored evidence and contradicts the Spec.

TQ-0019 implements canonicalization, manifests, bounded graph validation,
payload integrity and mutation fixtures. TQ-0020 implements attestations,
explicit public trust configuration, rotation/revocation and signing failures.
Later domain tasks must supply their payload schemas and cross-edge predicates
before enabling their stages. This design supplies no real source provenance,
approved fidelity weights, model results, presentation assets, or deployment.
