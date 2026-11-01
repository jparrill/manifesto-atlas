# Immutable source identity and registration

Contract: `atlas.source/v1`, TQ-0011; implementation TQ-0027. Consumes
[artifact identities](artifact-contract.md), [official provenance](provenance-trust.md)
and [isolation](submission-isolation.md). Refines
[immutable storage](../SPEC.md#immutable-source-storage),
[source tests](../SPEC.md#source-registry-tests) and
[assessment separation](../SPEC.md#assessment-separation).

## Decisions and record interface

`Register(verifiedProvenanceID, metadata, store) -> sourceRecordID` accepts only
complete metadata, matching independently downloaded bytes, and a current valid
provenance attestation/chain. Registry code does not retrieve new URLs or repair
PDFs. Source bytes never change after registration, whether committed yet or not.

The strict JSON record has these required fields (nullable only as specified):

| Field | Contract |
| --- | --- |
| `version` | Exactly `atlas.source/v1` |
| `party_id` | Existing scoped party-registry identity; no name-based inference |
| `programme_id` | H(canonical identity tuple below); stable across source revisions |
| `document_id` | Portable slug identifying this document within that programme; v1 permits exactly one source document |
| `document_type` | `electoral-programme` or `general-political` |
| `series_id` | General-political series slug; null for electoral programmes |
| `scope` | `european`, `national`, `autonomous`, `provincial`, `municipal`, `general-political` |
| `geography` | `{registry_version, id, country, parent_id}`; ID from a pinned official-geography mapping, not an arbitrary name |
| `jurisdiction` | Explicit institution/competence-set identifier and version; not inferred solely from source language |
| `election` | Electoral: `{id, type, date}` with evidenced ISO date; non-electoral: null |
| `publication_date` | Evidenced ISO date or null; never infer from upload date |
| `archive_year`, `year_basis` | Four-digit year; `publication` when evidenced, otherwise `retrieval` from canonical retrieval UTC year |
| `revision_id` | `r-` plus the full 64-hex canonical source checksum |
| `supersedes` | Previous source-record ID or null; same programme, different checksum, no cycles |
| `source` | Exact PDF hash/size, artifact manifest ID and safe registered relative path |
| `provenance` | Verified provenance manifest ID, party-registry identity, requested/final URLs and evidence references |
| `retrieved_at` | Exact canonical-download UTC timestamp from provenance, not registration wall time |
| `original_filename` | Bounded display-only string; basename metadata, never an output path |
| `publisher_rights` | Publisher identity, original rights statement/evidence or explicit unknown; never MIT by default |

The programme identity tuple is `{party_id, document_type, scope, geography_id,
election_id}`; election_id is null only for general-political documents.
General-political series also include `series_id`, an explicit evidenced or
owner-assigned administrative slug, so unrelated political documents are not
merged. That administrative identifier makes no electoral/date claim. `document_id`
is a stable archive label, not a second identity source. Reject unknown fields,
unknown registry versions, conflicting IDs and missing required evidence.

Use ISO 3166/official statistical IDs through the versioned geography mapping:
`es`, `es-cl`, and municipality IDs such as `es-ine-01059` are illustrative
identifiers; the mapping must verify their official association before use.
Parent/territory conflicts fail. Election IDs are explicit registry entries
matching date/type/scope/geography; identical dates alone do not identify events.

Processing lifecycle (`draft`, validating, validated, published) is derived
from current approvals as in the artifact contract, not mutable PDF metadata.
Viability, historical reliability, fulfilment and political scores cannot be
fields in this source-facts record.

## Year-first storage and worked records

Path construction is deterministic:

`sources/<archive-year>/<scope>/<geography-id>/<event>/<party-id>/<document-id>/<revision-id>/source.pdf`

`event` is the evidenced electoral date `YYYY-MM-DD`, or `none-<series-id>` for
non-electoral records. Inside the revision directory, `source.pdf` and
`metadata.json` are immutable. Store the record's own ID in an external index,
not in its own hashed payload. The programme index and publication selection
are derived/rebuildable; they do not replace historical records.

Every slug is lowercase ASCII `[a-z][a-z0-9-]{1,63}`, not a Windows reserved
device basename; disallow dots, separators, empty segments and trailing spaces.
Year/date and checksum segments use their fixed validated syntax instead of
the slug pattern. Refuse an overall relative path longer than 512 bytes.
Original filenames never enter path construction. Open rooted regular files
without following symlinks; never concatenate unchecked contributor strings.

The following are synthetic examples, not real election/date claims. `R` below
means an actual `r-<64hex>` source checksum segment, not a literal directory.
Each record uses the common fields above plus the values in this table.

| Record | Identity fields and worked path |
| --- | --- |
| National electoral | scope=national, geography=es, election={id: fixture-national, type: national, date: 2030-01-01}, archive_year=2030, year_basis=publication; `sources/2030/national/es/2030-01-01/fixture-party/programme/R/source.pdf` |
| Autonomous electoral | scope=autonomous, geography=es-cl, election={id: fixture-regional, type: autonomous, date: 2030-01-01}; `sources/2030/autonomous/es-cl/2030-01-01/fixture-party/programme/R/source.pdf` |
| Municipal electoral | scope=municipal, geography=es-ine-01059, election={id: fixture-local, type: municipal, date: 2030-01-01}; `sources/2030/municipal/es-ine-01059/2030-01-01/fixture-party/programme/R/source.pdf` |
| European electoral | scope=european, geography=es (Spanish electoral constituency), election={id: fixture-eu, type: european, date: 2030-01-01}; `sources/2030/european/es/2030-01-01/fixture-party/programme/R/source.pdf` |
| Non-electoral | type/scope=general-political, geography=es, election=null, series_id=principles, publication_date=null, canonical retrieval in 2030, year_basis=retrieval; `sources/2030/general-political/es/none-principles/fixture-party/programme/R/source.pdf` |

European scope does not mean every proposal belongs to EU institutions; the
jurisdiction contract distinguishes scope and competence. Provincial documents
use the same path grammar with scope=provincial and a mapped province ID.

### Initial golden sources

The owner-approved Frente Obrero document is a general political document, not
an invented national-election programme: use general-political type/scope,
election=null, and verify its actual territory, date/year and series evidence
on recovery. The VOX Castilla y León source must use the evidenced regional
party/scope/geography; do not invent its election date from the year 2026 or
reuse a different election's PDF. Electoral metadata cannot remain null on an
accepted electoral record. Missing evidence leaves TQ-0028/TQ-0029 blocked.
These rules intentionally provide no fabricated real checksum, URL or filename.

## Duplicates, revisions and atomicity

Compute a global checksum lookup before allocating a new source directory.
Exact same programme/checksum/metadata is an idempotent no-op returning the
existing record. Same bytes under another evidenced programme create a new
association to the existing immutable PDF path, not duplicate PDF bytes; the
association record has a new identity and lives in the artifact store. It must
pass independent provenance/scope checks for that programme. Similar text or a
shared party name never implies a duplicate.

A changed official PDF requires a new verified source and explicit `supersedes`
link. The prior revision must belong to the same programme; no implicit
"latest by timestamp", accidental fork overwrite or replacement of old bytes.
Independent correction branches are retained; current selection must be explicit
and cannot publish an ambiguous revision. Same hash with conflicting metadata
is `metadata-conflict`, not a new PDF revision. Resolve through a new, evidenced
metadata-correction association referencing the prior record; old metadata/path
and source bytes remain historical. An amended election identity is a corrected
programme association, never a rename/rewrite of the archived source.

Registration locks the registry, validates the full provenance closure and
rechecks existing paths/checksums under the lock. Copy from an already opened
verified handle into a private sibling staging directory, fsync files and
directory, then install the complete revision directory atomically with
no-replace semantics. Install the immutable source-record manifest only after
the directory is durable. Derived indexes update last and rebuild on restart.
On collision compare content; never overwrite. Interrupted orphan directories
remain non-publishable until verified recovery/rebuild; no partial registration
is accepted. Source and metadata file paths are opened no-follow.

For every check, compare archived bytes against the stored checksum and the
source manifest. CI compares immutable source paths with the configured trusted
base Git tree: changed/deleted existing PDFs fail even if someone also edited
the checksum. New revisions are allowed, overwrites are not. No privileged or
force option bypasses this. Missing baseline/evidence blocks immutability proof.

## Scenario-to-expected-outcome table

| Scenario | Expected outcome |
| --- | --- |
| First complete verified source | One durable revision and source-record; source bytes immutable |
| Exact same programme, bytes and metadata submitted twice | Same record/path; no duplicate PDF or spurious timestamp revision |
| Same bytes, different evidenced programme | New association, existing bytes reused; scope/provenance still required |
| Corrected official PDF, valid prior revision link | New checksum directory and record; old bytes untouched; descendants invalidated |
| Changed bytes with no explicit revision relationship | Reject ambiguous registration, do not infer supersession |
| Same checksum, conflicting party/election/scope metadata | `metadata-conflict`; append evidenced correction association, never overwrite |
| Missing election date on electoral programme | Block registration; retain quarantined evidence, never substitute retrieval date |
| Non-electoral document with no publication date | election=null; use retrieval archive year with explicit year_basis |
| Original filename includes traversal or shell metacharacters | Display metadata only; generated path unchanged, no execution |
| Unsafe slug, reserved name, symlink or long generated path | `unsafe-path`; no write outside trusted root |
| Committed PDF changed/deleted, even with edited metadata hash | Immutability gate fails against trusted Git base; publication blocked |
| Stale/mismatched source checksum or provenance report | Reject; no registration or downstream approval |
| Concurrent duplicate or path collision | Lock and compare; idempotent reuse or explicit conflict, never replace |
| Interrupted install / lost index update | No apparent publication; durable directory inspected and index rebuilt |
| Supersedes another party/programme or revision cycle | Reject invalid relationship |
| Historical PDF later disappears from party website | Archive bytes retained; current authority/evidence checks still apply |

## Versioning, test strategy and rejected alternatives

Schema/registry/geography changes create new versioned associations and
invalidate dependent evidence, never rewrite old sources. TQ-0027 must test
duplicates, races, interrupted registration, metadata conflicts, unsafe paths,
linked corrections and committed-byte tampering. TQ-0028/TQ-0029 must recover
the exact approved sources and prove real metadata rather than using these
synthetic examples. Third-party source rights remain with the publisher.

Rejected: source filenames as IDs; guessed election dates; content similarity
as checksum deduplication; overwrite-in-place corrections; automatic latest
selection; mixing performance/fulfilment assessments into programme facts.
