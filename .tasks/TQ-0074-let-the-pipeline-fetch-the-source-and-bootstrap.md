---
id: TQ-0074
title: Let the pipeline fetch the source and bootstrap the parity check
status: in-progress
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/feature
  - area/provenance
  - agent/developer
created: 2026-11-18T00:40:02+01:00
updated: 2026-11-18T00:40:02+01:00
---

## Outcome

`make process URL=… PROG=…` needs no local PDF and no manual setup for the parity
check: the pipeline downloads the document from its official URL into an ignored
cache, hashes it while streaming, and Playwright is installed by one target.

## Scope

- Included: `runsource.Fetch` (https only, PDF content type, 120 MB cap, hash-keyed
  cache with the original file name), `auto` as the PDF argument of `extract` and
  `process`, an optional positional for `process`, `make view-parity-setup` and the
  Playwright discovery in `scripts/view-parity.mjs`, `.cache/` ignored.
- Excluded: automatic admission of arbitrary third-party URLs (the owner still
  supplies the URL), and re-signing a document whose URL changed its bytes.

## Acceptance criteria

- [ ] `make extract URL=… PROG=…` with no PDF downloads, hashes and extracts.
- [ ] The cache path is `.cache/sources/<sha256>/<nombre>.pdf` and reused.
- [ ] A non-https origin, a non-PDF answer and an oversized body are refused.
- [ ] `make process` with no positional works, and without URL or model exits 2 with
      an actionable message.
- [ ] `make view-parity-setup` then `make view-parity` works with no environment.

## Quality Gates

### Gate 1: Acceptance criteria

`go test ./internal/runsource ./internal/cli` plus a real `make extract URL=…`.
Pass: the subtests pass; the real extract reports the expected sha256, 66 pages and
115 numbered proposals.

### Gate 2: Focused behavior

`go test ./internal/runsource -run TestFetch`
Pass: cache path `<hash>/<name>.pdf`, scheme refusal, content-type refusal, status
error.

### Gate 3: Regression

`make verify` → `verify: OK`.

### Gate 4: Static and structural validation

`gofmt -l`, `go vet ./...`, and a real parity run after the bootstrap.
Pass: no findings; `make view-parity` reports 0 style differences and every
behaviour check ok.

### Gate 5: Documentation

`docs/design/programme-run.md#8` documents the download and the parity bootstrapping.
Pass: `make verify` reports `docs OK` and `links OK`.

### Gate 6: Trust-boundary

The pipeline still refuses anything that is not https or not a PDF, and the hash is
computed while streaming before any file is accepted.
Pass: the subtests above.

### Gate 7: Artifact and task record

Committed with this record.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/design/programme-run.md`

---

## Notes

- 2026-11-18T00:40:02+01:00 — Gate evidence (all PASS):
  - G1: real 'make extract URL=https://www.voxespana.es/.../Programa-completo-A5-1.pdf PROG=voxfetch RUN=runs/tmp-fetch' downloaded 14536418 bytes, sha256 065abf67..., 'extracted 66 pages', 'proposals 115 numbered items'; the run was removed afterwards. TestProcessWithoutArguments covers process with no positional and without URL/model (exit 2, actionable message).
  - G2: go test ./internal/runsource -run TestFetch -> cache path .cache/sources/<sha256>/programa.pdf, scheme refusal, non-PDF content type refusal, non-200 status.
  - G3: make verify -> verify: OK.
  - G4: gofmt -l empty, go vet clean. make view-parity-setup installed Playwright into .cache/view-parity and 'make view-parity' then reported 'diferencias de estilo: 0' with 7/7 behaviour checks ok, without any environment variable.
  - G5: docs/design/programme-run.md#8 documents the download (cache layout, https-only, 120 MB cap) and the parity bootstrap.
  - G6: refuse non-https, non-PDF and oversized bodies; the hash is computed while streaming and the file is only renamed into the cache once measured.
  - G7: committed with this record.

  Decisions: the cache keeps the original file name inside a hash directory so the extraction files keep readable names ('Programa-completo-A5-1.md') and two runs of the same document share bytes. Playwright lives in .cache/view-parity and .cache is gitignored: the project itself has no Node dependency.
