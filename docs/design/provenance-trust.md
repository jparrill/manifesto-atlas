# Official-source provenance trust

> Historical future-platform contract. [MVP.md](MVP.md) and
> [SPEC.md](../SPEC.md) use explicit maintainer-selected local sources instead.
> Automatic network admission/domain/key grants are deferred, not certified by the MVP.

Contract: `atlas.provenance/v1`, TQ-0010. Implementation owners: TQ-0022,
TQ-0023, TQ-0025 and TQ-0026. This contract refines the
[provenance trust model](../SPEC.md#provenance-trust-model) and
[known trust-anchor problem](../SPEC.md#further-notes), and consumes the
[artifact](artifact-contract.md) and [isolation](submission-isolation.md)
contracts. It establishes no real party/domain anchor or official PDF URL.

## Decisions and interfaces

`FetchEvidence(request, networkPolicy) -> immutableEvidence` owns bounded network
access. `CheckAnchor(partyID, registry, evaluationTime) -> anchor | error`
checks independently established publisher authority. The semantic provenance
job receives evidence only. `DecideProvenance(submission, canonicalDownload,
anchor, evidence, semanticReport) -> decision` performs deterministic checksum,
chain, freshness and signature checks. Reports follow the artifact contract and
are signed by an authorized provenance-report signer, not by the model.

The four decisions are `verified`, `probable`, `unverified`, `rejected`.
`verified` also requires report `result=pass`. All other decisions block source
registration/publication. Transport errors preserve a failed/insufficient
attempt; they do not imply contradictory origin. Contradictory positive
evidence (wrong canonical bytes or invalid claimed publisher signature) is
`rejected`. No branding, model confidence score, filename or contributor vote
can complete a missing authority link.

## Party registry and anchor authority

A party entry has version `atlas.party/v1`, immutable lowercase ASCII `party_id`
(pattern `[a-z][a-z0-9-]{1,63}`), legal/display names, optional electoral
registration identifier, and explicit domain and publisher-key grants. A merger,
coalition or renamed legal entity gets an explicit relationship; never join by
display name or logo. Geographic party organizations retain scoped grants.

Each domain grant records an exact lowercase ASCII A-label hostname, party ID,
allowed territorial/document scope, evidence artifact IDs, `reviewed_at`,
`valid_until` and active/revoked state. No wildcard, substring, suffix inference
or automatic `www`/subdomain trust. A campaign/regional hostname requires its own
grant: evidence from an existing approved official party page linking and
identifying that organization/domain, or an authoritative electoral/legal
record explicitly connecting the legal party to that domain. Generic search
results, social profiles, domain WHOIS alone and self-hosted declarations are
insufficient to bootstrap an anchor.

Juan approves registry changes in the trusted checkout after inspecting the
evidence. Contributors/models can propose grants but cannot install/activate
their own authority. This is trust-root maintenance, not a human-review fallback
for publishing an insufficient programme. Initial grants for a new party need
the same independent legal/electoral evidence; if unavailable, keep submissions
unverified. No invented domains are seeded by implementation fixtures.

Grants expire no later than 180 days after review. Check `reviewed_at <= now <=
valid_until`, revocation, scope, and evidence integrity before every new source
admission and publication check. Revalidation requires fresh independent
evidence and a new reviewed registry version; no automatic extension from the
site's continued existence. Compromise, domain sale, or evidence contradiction
revokes the grant and invalidates current provenance-dependent admission.
Preserve historical reports and old source bytes. New grants never silently
reinterpret earlier registry versions. The caller explicitly selects the current
trusted registry, not one supplied by a contributor's artifact graph.

## Permitted trust chains

All chains independently download the canonical bytes, inspect them under the
isolation contract and require their SHA-256 to equal the candidate's checksum.
An upload cannot serve as its own independent download. Retrieve metadata and
headers with the actual bounded GET; HEAD alone never establishes byte identity.

1. **Direct domain:** the PDF's final HTTPS host has an active, in-scope domain
   grant. Every redirect hop must satisfy the network policy. A final
   third-party host is not direct-domain verification even if the initial URL
   was official; use the referral chain instead.
2. **Exact official referral:** GET a page on an active official host; preserve
   the response body, final page URL, source timestamp, headers and redirect
   chain. A deterministic parser resolves an actual PDF link (`a[href]` or
   `link[href]` with declared PDF type) against the final page URL. Ignore
   `<base>` overrides, scripts, meta refresh, embedded JSON claims and prose
   mentions. The requested canonical PDF URL must equal that resolved link
   under the normalization rules below. Fetch that URL independently and retain
   every redirect to the final third-party/CDN object and exact returned bytes.
   The host need not be a party domain, but each hop still passes network checks.
   Model opinion cannot fabricate a link absent from the stored page.
3. **Publisher signature:** the party registry independently grants the
   publisher key and scope. v1 supports detached Ed25519 signatures over the
   exact downloaded PDF bytes with a pinned 32-byte public key; signatures are
   standard base64 64-byte values. The signature, key grant/evidence, byte
   checksum, independently fetched URL and retrieval evidence must all match.
   Keys have reviewed validity intervals and current revocation, checked as in
   the artifact contract. A self-supplied key or cryptographically valid key
   belonging to another publisher is not authority. Other PDF/CMS/OpenPGP
   signature formats are `probable` only if there is positive publisher evidence
   and unsupported verification is the sole missing link; otherwise unverified.
   Never call an unimplemented signature verifier a success.

Publisher signatures authenticate source origin. Report Ed25519 attestations
authenticate processing roles. Their trust roots and signed messages are
distinct; one cannot substitute for the other.

An old archived referral can preserve the provenance of an already registered
immutable revision with its original evidence. It cannot bootstrap a **new**
verified source when the live authority link is unavailable. Record `probable`
when a reputable archive preserves an exact former official referral but one
current required link is missing; archive hosts are not official party anchors.
Social-only publication remains unverified absent an independent approved
domain/key chain. Removed documents are not silently replaced by a different
election/revision or a press reproduction.

## URL and bounded network policy

- HTTPS only, TCP port 443 only; no URL userinfo, fragments on document requests,
  IP-literal hosts, trailing-dot hosts, backslashes, empty hosts, embedded
  whitespace/control characters, encoded NUL/CR/LF, or invalid escapes.
  Reject Unicode host input rather than guessing IDNA transformations; approved
  A-labels are explicit. Lowercase scheme/host and omit explicit default port.
- Preserve path/query bytes and parameter order; do not drop tracking queries,
  decode reserved separators, collapse paths, or assume two queries identify the
  same object. Resolve relative references with RFC 3986 rules. Reject encoded
  separators/traversal ambiguities. Empty path becomes `/`. Fragment-only page
  links are not PDF references. Reject bearer-style query fields (`token`,
  `access_token`, `auth`, `key`, `signature`, `sig`, case-insensitive), not just
  credentials in userinfo; v1 fetches publicly accessible unsigned URLs only.
- Resolve DNS before **each** connection/redirect. Reject the whole answer set
  if any address is not publicly routable (loopback, private, link-local,
  multicast, unspecified, documentation, benchmark, carrier NAT or other
  reserved/special-purpose ranges, including mapped IPv4 and IPv6 variants).
  A pinned, versioned special-purpose range table supplements Go IP predicates;
  `IsGlobalUnicast` alone is insufficient. Empty/mixed public-private answers
  fail. Pin the approved resolved IP to the actual dial while retaining the
  original hostname for TLS SNI and certificate verification. Never re-resolve
  inside the default transport and thereby reopen DNS-rebinding access.
- Disable environment proxies, cookies, ambient credentials, `.netrc`, client
  certificates, connection reuse across origins and automatic redirects. Verify
  normal TLS certificates/hostname; no insecure override. Allow at most 5
  redirects; reject loops, unsupported schemes/ports, unsafe IPs or TLS errors
  at any hop. No credentials/authorization headers are forwarded.
- Use a fixed User-Agent and `Accept-Encoding: identity`; disable transparent
  decompression and reject non-identity Content-Encoding. Response status must
  be 200; reject partial responses, login/challenge bodies, empty PDF downloads
  and Content-Length/count mismatch. Declared wrong PDF MIME rejects; absent
  MIME needs parser+magic agreement under the isolation contract.
- Budgets per attempt: DNS 5 s, connect/TLS 10 s, response headers 10 s, total
  chain/body deadline 120 s, headers 64 KiB, referral page 2 MiB, signature
  response 16 KiB, canonical PDF as the approved isolation size limit. Count
  actual body bytes with limit + 1 detection. At most 10 GET requests per job;
  no implicit retries or browser/JavaScript execution. Cancellation closes the
  response and partial scratch. New attempts are explicit evidence records.

The transport is a trusted bounded adapter, separate from network-denied model
and PDF processes. Hostnames/URLs are parsed as data, never shell commands.
Unsafe requests fail before dialing; remote 404, timeout or unsafe redirect
cannot produce `verified` even if a previous candidate has a familiar checksum.

## Required evidence and decision scenarios

Every report binds submitted/canonical checksums, party/registry/grant versions,
scope/election identity, chain kind, submitted/requested/final URLs, each hop's
status/Location/approved IPs, UTC retrieval time, MIME/length/ETag/Last-Modified
when supplied, original filename, response body/evidence digests, exact referral
link or publisher signature, inspection result, semantic report and tool/model/
prompt/configuration identities. Missing headers are explicit null values, not
invented facts. Private headers/cookies are never logged. Record archival
evidence as optional corroboration, with its own retrieval identity.

| Scenario | Expected outcome |
| --- | --- |
| Direct final official PDF, active scoped grant, matching bytes | `verified`, provided inspection and report pass |
| Exact official page link to CDN, safe complete redirect chain, matching bytes | `verified` referral chain; preserve page and all hops |
| Valid authorized publisher signature over matching independent download | `verified` signature chain |
| Claimed publisher signature cryptographically invalid / wrong bytes | `rejected`; retain contradictory evidence |
| Valid signature under unapproved/self-supplied key | `unverified`; never bootstrap authority from the submission |
| Lookalike/subdomain of approved name without exact grant | `unverified`; substring/suffix/branding is not evidence |
| Official initial URL redirects to third party without an exact referral | `probable` if official-origin evidence exists; not direct verified |
| Hostile redirect to localhost/private IP/HTTP/credential URL | Transport blocked; insufficient non-publishable attempt |
| DNS answer changes or mixes public/private results | Block before connection; never choose just the convenient public IP |
| No trust anchor / new party with only self-asserted domain | `unverified`; contributor cannot activate grant |
| Removed referral preserved only in third-party archive, new submission | `probable` if exact historical link is evidenced; not publishable |
| Existing immutable source with intact archived original verification evidence | Preserve historical provenance; current publication still checks active authority and exact inputs |
| Social-only post, party branding, filename or press copy | `unverified` absent another permitted chain |
| Candidate differs from independently fetched canonical bytes | `rejected`, even if both look official; correction requires a separate submission |
| Canonical URL silently changes bytes after registration | New revision required; old PDF untouched; stale downstream results blocked |
| Expired/revoked/wrong-scope anchor or publisher key | Current admission blocked; retain historical evidence, require revalidation |
| Same path, different query bytes or guessed CDN equivalent | Not an exact referral; no normalization-based authority shortcut |
| Page contains script-only link, meta refresh or misleading base URL | No admissible exact referral; probable only with positive incomplete official evidence |
| Missing page/PDF evidence, response limit, timeout, corrupt PDF | Failed/insufficient or rejected inspection; never `verified` |
| Unsupported publisher signature format with positive independent evidence | `probable`; implementation gap cannot become a pass |
| Model says official but deterministic chain check fails | Non-publishable; deterministic gate prevails |
| All evidence absent vs positively contradictory | `unverified` vs `rejected`; do not conflate uncertainty with contradiction |

## Versioning and rejected alternatives

Version registry/grants, normalization/network policy and signature formats.
Changes invalidate exactly dependent provenance/registration descendants under
the artifact contract. Keep failure reports. The first approved implementation
must test direct/referral/signature chains using synthetic local HTTP/TLS
fixtures with injected resolver/dialer and a fixture trust registry, including
SSRF, rebinding, redirect, MIME and byte mismatch. Fixture-only private endpoints
must not become a production allow-private option. Real golden sources require
actual independently established anchors/evidence, never fixture substitution.

Rejected: model-only origin votes; uploading a PDF as canonical; wildcard domain
trust; archive/social/search bootstrap; following default HTTP redirects;
re-resolving after validation; HEAD-only checks; silent publisher-format fallbacks.
No external service or human programme-approval gate is introduced. Owner review
of trust-root changes remains distinct from automatic document validation.
