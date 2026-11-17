---
id: TQ-0069
title: Verify official source URLs and build the portable preview artifacts
status: done
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/docs
  - area/provenance
  - area/rendering
  - agent/developer
created: 2026-11-17T23:25:18+01:00
updated: 2026-11-17T23:25:25+01:00
---

## Outcome

The run records verified official provenance for its three sources and ships two
portable artifacts: a self-contained HTML view with the summaries and verification
reports embedded, and a single A4 compendium PDF with index, summaries and reports.

## Scope

- Included: official URL and SHA-256 per source with the download check, `web` link
  per party in the view, `Programas-electorales-2026-autocontenido.html`,
  `Programas-29N-2026-compendio.pdf`, and the generators in `tools/`.
- Excluded: embedding the 28 MB of source PDFs inside the HTML (rejected: a 38 MB
  document for no reading benefit), and promoting the sources to the canonical
  catalogue.

## Acceptance criteria

- [ ] `SOURCES.md` records URL, size, hash and the download check for all three sources.
- [ ] Downloading each official URL reproduces the recorded SHA-256.
- [ ] The self-contained view opens with no sibling files and no network, and shows
      the full summary and the verification report for each party.
- [ ] The compendium PDF contains the six text documents, an index and the provenance.
- [ ] `make verify` passes.

## Quality Gates

### Gate 1: Acceptance criteria

`ls runs/29n2026` plus the products table in `runs/29n2026/README.md`.
Pass: every artifact exists and the table describes it.

### Gate 2: Focused behavior

`curl -sSL <url> | shasum -a 256` for the three official URLs.
Pass: the three hashes equal the ones in `SOURCES.md`.

### Gate 3: Regression

`make verify` at the repository root.
Pass: `verify: OK` including docs and relative links.

### Gate 4: Static and structural validation

`node --check` on the extracted script of both HTML files; headless Chrome
`--dump-dom` on the self-contained file; `pdftotext` over the compendium.
Pass: no syntax error, the embedded JSON parses with the six documents and no
base64 payload, and every distinctive word of the six Markdown files is present in
the PDF text.

### Gate 5: Documentation and specification consistency

`runs/29n2026/README.md`, `SOURCES.md` and the root `README.md` pointer.
Pass: the new artifacts are described with their trade-offs and `make verify` reports
`docs OK` and `links OK`.

### Gate 6: Trust-boundary and security behavior

The official URLs are recorded as a point-in-time check, not as certification, and
the documents say so. The embedded JSON escapes `</` so no document can close the
script tag it lives in.
Pass: both statements present in `SOURCES.md` and README; no unescaped `</script>`
inside the payload.

### Gate 7: Artifact and task record

Artifacts and this record are committed and pushed together.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/SPEC.md#priority-for-primary-sources`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-17T23:25:24+01:00 — Gate evidence (all PASS):
  - G1: artifacts in runs/29n2026 -> SOURCES.md, Programas-electorales-2026-autocontenido.html (192 KB), Programas-29N-2026-compendio.pdf (0.9 MB, 45 pages), tools/build_standalone.py, tools/build_compendium.py.
  - G2: downloaded the three official URLs on 2026-11-17; sha256 FO 7b9f423f..., VOX 065abf67..., PP ae0c35d7... all equal to the analysed files and to SOURCES.md. HTTP 200 application/pdf in the three.
  - G3: make verify -> verify: OK (taskqueue, gofmt, docs OK, links OK, all go tests).
  - G4: node --check OK on both HTML scripts; Chrome --dump-dom on the self-contained file renders 'Leer resumen completo' in the toolbar, the three buttons in the sources row and the official link; embedded JSON parses with 6 documents (fo 2, vox 2, pp 2) and no b64 key; pdftotext of the compendium contains 100% of the distinctive words of the six Markdown files.
  - G5: README documents the trade-off (no 28 MB base64), SOURCES.md documents the verification and its limits, root README points at the run.
  - G6: SOURCES.md states it is a point-in-time check, not a third-party certification, and that a replaced file invalidates the run for that document. Payload escapes '</' so no embedded document can close its script tag.

  Decisions: rejected the 38 MB all-in-one HTML (base64 of the source PDFs) after measuring it; the portable artifacts are the 192 KB self-contained view (text embedded) and the 45-page compendium PDF. Source PDFs stay out of Git and are identified by official URL + SHA-256.
