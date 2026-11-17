---
id: TQ-0070
title: Make the programme-run process deterministic and gated
status: done
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/feature
  - area/validation
  - area/rendering
  - agent/developer
created: 2026-11-17T23:43:42+01:00
updated: 2026-11-17T23:43:55+01:00
---

## Outcome

The refinement process used for the 29-N run is part of the project: a documented
run contract, a deterministic gate implemented in Go, curated claim and exception
files, and a `make check-run` target that `make verify` always runs. The gate also
proves that the view embeds every document and needs no sibling file.

## Scope

- Included: `docs/design/programme-run.md`, `internal/runcheck` with tests,
  `cmd/runcheck`, `runs/29n2026/run.json`, `checks/claims.tsv`, `checks/allow.tsv`,
  the `check-run` target wired into `verify`, and the pointer in `AGENTS.md` and
  `docs/MAKEFILE.md`. The reference Python checkers are removed; Go is the authority.
- Excluded: promoting the run to the canonical catalogue, model calls, network use,
  and embedding the source or summary PDFs.

## Acceptance criteria

- [ ] `make check-run` verifies a run offline and exits 1 with one line per issue.
- [ ] It checks literal quotes, page citations, proposal coverage, curated claims,
      section codes and self-containment of the view.
- [ ] Coverage supports documents that number proposals once and documents that
      restart numbering per section.
- [ ] Exceptions require a reason and never silence a different issue.
- [ ] `make verify` runs the gate, and the 29-N run passes it.
- [ ] The run contract documents the phases, the failure modes seen in practice and
      how to add a run.

## Quality Gates

### Gate 1: Acceptance criteria

`make check-run` on the real run and `ls runs/29n2026`.
Pass: exit 0 with every check reporting "ok"; the files named in the contract exist.

### Gate 2: Focused behavior

`go test ./internal/runcheck/...`
Pass: one subtest per check and per failure mode, all passing, including the
exception-without-reason rejection.

### Gate 3: Regression

`make verify`
Pass: `verify: OK`, with `verify: programme run gate` and `go test (test-runcheck)`
in the output.

### Gate 4: Static and structural validation

`gofmt -l`, `go vet ./...` and a negative check: removing a document from the view
payload must make the gate fail.
Pass: no formatting or vet findings; the gate reports `runcheck.embed` with exit 1.

### Gate 5: Documentation and specification consistency

`docs/design/programme-run.md`, `docs/MAKEFILE.md#baseline`, `AGENTS.md` and the run
`README.md`.
Pass: `make verify` reports `docs OK` and `links OK`, and the four files describe the
same target and contract.

### Gate 6: Trust-boundary and security behavior

The gate never contacts the network, never reads the source PDFs and never writes to
the run. Provenance stays a point-in-time verification, not a certification.
Pass: the checks only open files inside the run directory; the contract states the
limit; `SOURCES.md` keeps the disclaimer.

### Gate 7: Artifact and task record

The gate output is recorded in the task note and the run README, and code, documents
and this record are committed and pushed together.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/design/programme-run.md`
- `docs/QUALITY_GATES.md`

---

## Notes

- 2026-11-17T23:43:55+01:00 — Gate evidence (all PASS):
  - G1: make check-run on runs/29n2026 -> quotes 222 ok, page refs 379 ok, proposals 497 ok, claims 69 ok, sections 28 ok, embedding 6 documents ok, 4 reviewed exceptions. Exit 0.
  - G2: go test ./internal/runcheck/... -> ok, 20 subtests covering each check and each failure mode (literal quote, paraphrase, digit-vs-word, wrong page, coverage by numbers and by items, invented section code, claim on the wrong page, malformed claim row, exception applied, exception without reason rejected, exception not hiding another issue, view without payload, local ref in markup, local ref inside an embedded document, external stylesheet, programme without embedded document, missing viewer).
  - G3: make verify -> verify: OK with 'verify: programme run gate' and 'go test (test-runcheck)' in the output.
  - G4: gofmt -l empty, go vet ./... clean. Negative test: emptying the vox payload makes the gate report 'programme vox embeds no readable document' and exit 1; the view is restored afterwards and the gate passes again.
  - G5: docs/design/programme-run.md (phases, artifacts, checks, failure taxonomy of 14 real cases, how to add a run, gate mapping), docs/MAKEFILE.md row, AGENTS.md pointer, run README with the retained output. make verify reports docs OK and links OK.
  - G6: the gate reads only files inside the run directory, never the network and never the source PDFs; the contract keeps provenance as a point-in-time check and keeps the empty-proposal rule.
  - G7: this note, the run README and the commit carry the evidence.

  Decisions: the reference Python checkers (verify.py, datos.py) are removed and their 69 curated claims moved to checks/allow+claims TSV, so there is one authority. The view is now a single self-contained file: tools/template.html plus tools/build_standalone.py, and the previous two-file arrangement (template plus autocontenido) is gone because it invited opening the file without the payload. The summary PDFs are not embedded: they are generated from the same text the gate already covers.
