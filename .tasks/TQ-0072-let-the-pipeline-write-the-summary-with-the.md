---
id: TQ-0072
title: Let the pipeline write the summary with the gate as authority
status: done
priority: high
assignee: claude-sonnet-4.6[1m]
labels:
  - type/feature
  - area/skills
  - agent/developer
created: 2026-11-18T00:05:50+01:00
updated: 2026-11-18T00:05:57+01:00
---

## Outcome

`make process PDF=… URL=… PROG=…` produces a complete run with no human step:
extraction, a model-written summary repaired against the gate, the verification
report and the self-contained view.

## Scope

- Included: `internal/summarize` (fixed contract, generation, bounded repair loop,
  credential-free trace), the `summarize` and `process` stages, the Make targets,
  and the documentation of what the pipeline does and does not guarantee.
- Excluded: choosing the model or the endpoint for the operator, storing
  credentials, and claiming byte-identical output between runs.

## Acceptance criteria

- [ ] `summarize` writes the summary from the extraction and passes only when the
      gate accepts it; otherwise it exits 1 leaving the issues on disk.
- [ ] A failing reply is repaired with the gate findings, bounded by `-attempts`.
- [ ] The trace records model, attempts and summary hash, and contains no prompt,
      reply, endpoint or key.
- [ ] `process` chains extract, summarize, validate-content and build, stopping at
      the first failure.
- [ ] `make verify` passes.

## Quality Gates

### Gate 1: Acceptance criteria

`go test ./internal/cli/ -run TestProcessRunsEverything` drives the whole pipeline
against a stub provider and asserts every artifact.
Pass: the subtest passes and the five files exist.

### Gate 2: Focused behavior

`go test ./internal/summarize/...`
Pass: fence stripping, input budget, missing configuration, provider failure, repair
until accepted (2 attempts, second request carries `runcheck.quote`), giving up after
the budget, and a trace without secrets.

### Gate 3: Regression

`make verify`
Pass: `verify: OK` including the pipeline tests.

### Gate 4: Static and structural validation

`gofmt -l`, `go vet ./...`.
Pass: no findings.

### Gate 5: Documentation

`docs/design/programme-run.md#8` and `docs/MAKEFILE.md`.
Pass: both describe `summarize` and `process`, the environment variables and the
non-determinism limit; `make verify` reports `docs OK` and `links OK`.

### Gate 6: Trust-boundary

The model cannot choose the endpoint, the model name or the checks; a reply that
fails the gate is never accepted; the trace has no prompt, reply or key; the
instruction to treat source text as data is in the fixed contract.
Pass: the subtests above and the trace assertion.

### Gate 7: Artifact and task record

Committed with the note holding the gate evidence.

## Final phase: delivery and continuation

See `AGENTS.md#closing-tasks-and-continuing-work`.

## Specification

- `docs/design/programme-run.md`

---

## Notes

- 2026-11-18T00:05:57+01:00 — Gate evidence (all PASS):
  - G1: go test ./internal/cli/ -run TestProcessRunsEverything -> process drives extract, summarize, validate-content and build against a stub provider; the run ends with run.json, SOURCES.md, Resumen-PP.md, Verificacion.md and Programas-electorales.html.
  - G2: go test ./internal/summarize/... -> fences stripped; input budget refused; missing endpoint/model refused; provider 500 reported; repair loop accepted on the second attempt with 'runcheck.quote' present in the repair request; budget exhausted leaves the run failing; trace has no key, no endpoint and no document text.
  - G3: make verify -> verify: OK, including internal/summarize in the pipeline tests.
  - G4: gofmt -l empty, go vet ./... clean.
  - G5: docs/design/programme-run.md#8 and docs/MAKEFILE.md document summarize and process, the ATLAS_* variables, the credential-free trace and the non-determinism limit; make verify reports docs OK and links OK.
  - G6: the fixed contract is a constant, the endpoint and model come only from the operator or the environment, a reply failing the gate is never accepted, and the gate remains the only authority.
  - G7: committed with this record.

  Decisions: the summary generator does not reuse semantic.ProviderAdapter because that contract is built for the canonical JSON publication roles; it speaks plain chat-completions with the same discipline (operator-fixed endpoint and model, key in the header only, bounded request and reply). process reuses the other stages instead of duplicating their logic, and stops at the first failure so a partial run is never presented as finished.
