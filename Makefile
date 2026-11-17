# Manifesto Atlas — Make orchestration interface.
#
# Baseline targets are implemented now (help, check, verify). The 12 pipeline
# stages from docs/SPEC.md#make-targets are declared but fail until their owning
# backlog task lands. Target-to-owner mapping and integration notes: docs/MAKEFILE.md
#
# Orchestration lives here. Deterministic work moves to Go once the module exists
# (owned by TQ-0018); these baseline checks only use Git and Taskqueue.

.DEFAULT_GOAL := help

# Programme run verified by `make check-run`. A run is a directory holding
# run.json; its contract is docs/design/programme-run.md.
RUN ?= runs/29n2026
SHELL := /bin/sh

# Baseline targets (implemented).
.PHONY: help check verify check-run view-parity test-cli test-artifacts test-taxonomy test-canonical test-semantic-runner test-fixture-frente test-fixture-vox test-runs smoke-provider test-preview e2e-preview evaluate-preview render-preview

# Pipeline stages (planned): fail until their owning task lands.
.PHONY: submit verify-source add extract validate-extraction normalize validate-content build validate-site publish-check process serve

help:
	@printf '%s\n' \
	  "Manifesto Atlas — Make targets" \
	  "" \
	  "Baseline (implemented):" \
	  "  make help                  Show this help." \
	  "  make check                 Run baseline checks (Git + Taskqueue + Go)." \
	  "  make verify                Run full verification (check + docs + links + run gate + tests)." \
	  "  make check-run             Run the deterministic second gate of a programme run (RUN=runs/29n2026)." \
	  "  make extract PDF= URL= PROG=  Run stages: provenance, extraction, manifest and summary skeleton." \
	  "  make validate-content      Run the gate and write Verificacion.md (RUN=runs/29n2026)." \
	  "  make summarize             Write the summary with the configured model and repair it with the gate (RUN= PROG=|-all)." \
	  "  make build                 Render the single-file view with the documents embedded (RUN=..., OUT=...)." \
	  "  make process PDF= URL= PROG=  Everything in one go: extract, summarize, validate, build." \
	  "  make view-parity CURATED= GENERATED=  Compare the generated view with the hand-made one (Playwright)." \
	  "  make test-cli              Run the focused Go CLI test suite." \
	  "  make test-artifacts        Run the focused strict canonical-JSON/hash helper suite." \
	  "  make test-taxonomy         Run the focused taxonomy/document-scope validation suite." \
	  "  make test-canonical        Run the focused canonical programme/citation validation suite." \
	  "  make test-semantic-runner  Run the focused semantic job runner (isolation/provider) suite." \
	  "  make test-fixture-frente   Run the focused Frente Obrero candidate/source fixture suite." \
	  "  make test-fixture-vox      Run the focused VOX Castilla y León candidate/source fixture suite." \
	  "  make test-runs             Run the focused selection/run record validation suite." \
	  "  make smoke-provider        Run the three clean semantic roles through the configured local provider (TARGET=...)." \
	  "  make e2e-preview           Real macOS PDF/model/private HTML run (OUTPUT=new-directory)." \
	  "  make render-preview        Rerender retained private evidence offline (no models)." \
	  "" \
	  "Pipeline stages (planned — fail until their owning task lands):" \
	  "  make submit                TQ-0021  Inspect a selected local candidate." \
	  "  make verify-source         TQ-0027  Verify local selection evidence and bytes." \
	  "  make add                   TQ-0027  Register a verified immutable source." \
	  "  make validate-extraction   TQ-0030  Deterministic page/warning blockers." \
	  "  make normalize             TQ-0038  Draft canonical programme JSON." \
	  "  make validate-site         TQ-0049  Site validation gate." \
	  "  make publish-check         TQ-0050 TQ-0049  Offline Markdown eligibility first." \
	  "  make serve                 TQ-0065  Local private preview." \
	  "" \
	  "Mapping and integration notes: docs/MAKEFILE.md"

check:
	@echo "check: git whitespace"
	@git diff --check
	@echo "check: taskqueue integrity"
	@$(call tq-clean)
	@echo "check: go format"
	@$(call go-fmt-check)
	@echo "check: go build"
	@go build -o /dev/null ./cmd/manifesto-atlas
	@echo "check: go vet"
	@go vet ./...
	@echo "check: OK"

verify: check
	@echo "verify: programme run gate"
	@$(MAKE) check-run
	@echo "verify: documentation structure"
	@$(call docs-check)
	@echo "verify: relative links"
	@$(call link-check)
	@echo "verify: go test (test-cli)"
	@$(MAKE) test-cli
	@echo "verify: go test (test-artifacts)"
	@$(MAKE) test-artifacts
	@echo "verify: go test (test-taxonomy)"
	@$(MAKE) test-taxonomy
	@echo "verify: go test (test-canonical)"
	@$(MAKE) test-canonical
	@echo "verify: go test (test-semantic-runner)"
	@$(MAKE) test-semantic-runner
	@echo "verify: go test (test-fixture-frente)"
	@$(MAKE) test-fixture-frente
	@echo "verify: go test (test-fixture-vox)"
	@$(MAKE) test-fixture-vox
	@echo "verify: go test (test-runs)"
	@$(MAKE) test-runs
	@$(MAKE) test-preview
	@echo "verify: go test (run pipeline)"
	@go test ./internal/runcheck/... ./internal/runsource/... ./internal/runview/... ./internal/mdrender/... ./internal/summarize/...
	@echo "verify: OK"

# Deterministic second gate of a programme run: literal quotes, page citations,
# proposal coverage, curated claims, section codes and a self-contained view.
# Offline, no models, no source PDFs; see docs/design/programme-run.md.
check-run:
	@go run ./cmd/runcheck -run $(RUN)

# Visual and behavioural parity between the generated view and the hand-made one.
# Needs Playwright reachable through PLAYWRIGHT_PATH (dev tool, not a dependency).
view-parity:
	@test -n "$(GENERATED)" -a -n "$(CURATED)" || { \
		echo "make view-parity CURATED=<html a mano> GENERATED=<html generado>"; exit 2; }
	@node scripts/view-parity.mjs "$(CURATED)" "$(GENERATED)"

# Focused gate for the Go CLI test suite. Runs the real tests without invoking
# live models or requiring live official websites.
test-cli:
	@go test ./internal/cli/...

# Focused gate for the strict canonical-JSON serialization and SHA-256 hash
# helpers. Exercises canonicalization, integer/UTF-8/surrogate handling, depth
# and size bounds, and hash validation offline; no live models or websites.
test-artifacts:
	@go test ./internal/artifact/...

# Focused gate for the atlas.taxonomy/v1 registry, topic states, and
# evidence-based document-scope classification. Deterministic and offline; no
# live models or official websites.
test-taxonomy:
	@go test ./internal/taxonomy/...

# Focused gate for the atlas.programme/v1 canonical schema, topic/state rules,
# citation/page/excerpt evidence, source binding, and deterministic bytes.
# Deterministic and offline; no live models or official websites.
test-canonical:
	@go test ./internal/canonical/...

# Common semantic response rejection and fake-adapter regression gate.
# Actual configured-provider, evidence/schema and runtime/session isolation
# qualification remains TQ-0024. No live models or committed credentials.
test-semantic-runner:
	@go test ./internal/semantic/...

# Focused gate for the retained Frente Obrero candidate (atlas.candidate/v1 with
# the embedded atlas.source.local/v1 record): offline hash and metadata
# verification of the recovered general-political source and its origin evidence.
# Deterministic and offline; no live models or official websites.
test-fixture-frente:
	@go test ./internal/candidate/...

# Focused gate for the retained VOX Castilla y León candidate (atlas.candidate/v1
# with the embedded atlas.source.local/v1 record): offline hash and metadata
# verification of the recovered electoral source, its evidenced 15 March 2026
# election, and its origin evidence. Deterministic and offline; no live models
# or official websites.
test-fixture-vox:
	@go test -run TestCheckVoxFixture ./internal/candidate/...

# Focused gate for the atlas.selection/v1 owner selection record and the
# atlas.run/v1 per-programme run record: strict parsing, deterministic
# round-trip, and offline bounded-file hash verification. No live models or
# official websites.
test-runs:
	@go test ./internal/run/...

# Explicit real-provider smoke gate (TQ-0024): runs the three clean semantic roles
# (normalizer, claim-validator, coverage-validator) through the operator-configured
# OpenAI-compatible local provider and retains the actual request/model identities
# and original bounded replies. Requires ATLAS_SMOKE_BASE_URL and ATLAS_SMOKE_MODEL
# in the environment; it is an operator gate, never part of make verify (regression
# is offline). TARGET defaults to the committed synthetic smoke fixture.
TARGET ?= fixtures/semantic/smoke.json
smoke-provider:
	@go run ./cmd/smoke-provider -target $(TARGET)

# The macOS preview is an explicit operator evaluation, never an offline gate.
SOURCE ?= fixtures/frente/source.json
OUTPUT ?= runs/preview/frente-e2e
test-preview:
	@go test ./internal/preview/...
e2e-preview evaluate-preview:
	@go run ./cmd/preview -source "$(SOURCE)" -output "$(OUTPUT)"
render-preview:
	@go run ./cmd/preview -render -output "$(OUTPUT)"

# --- Planned pipeline stages (fail until their owning task lands) -------------

submit:
	$(call unavailable,submit,TQ-0021,inspect a selected local candidate)
verify-source:
	$(call unavailable,verify-source,TQ-0027,verify local selection evidence and matching bytes)
add:
	$(call unavailable,add,TQ-0027,register a verified immutable source)
summarize:
	@go run ./cmd/manifesto-atlas summarize "$(or $(RUN),runs/29n2026)" \
		$(if $(PROG),-programme "$(PROG)",-all) -attempts "$(or $(ATTEMPTS),3)" \
		$(if $(ENDPOINT),-endpoint "$(ENDPOINT)",) $(if $(MODEL),-model "$(MODEL)",)

extract:
	@test -n "$(PDF)" || { echo "make extract: PDF=<ruta al pdf> URL=<url oficial> PROG=<id> [RUN=dir]"; exit 2; }
	@test -n "$(URL)" || { echo "make extract: falta URL (la procedencia oficial del PDF)"; exit 2; }
	@test -n "$(PROG)" || { echo "make extract: falta PROG (identificador corto, p. ej. pp)"; exit 2; }
	@go run ./cmd/manifesto-atlas extract "$(PDF)" -url "$(URL)" -programme "$(PROG)" \
		-run "$(or $(RUN),runs/$(PROG))" -title "$(or $(TITLE),$(PROG))" -coverage "$(or $(COVERAGE),numbers)"
validate-extraction:
	$(call unavailable,validate-extraction,TQ-0030,reject deterministic page/warning blockers)
normalize:
	$(call unavailable,normalize,TQ-0038,produce draft canonical programme JSON)
validate-content:
	@go run ./cmd/manifesto-atlas validate-content "$(or $(RUN),runs/29n2026)"
build:
	@go run ./cmd/manifesto-atlas build "$(or $(RUN),runs/29n2026)" $(if $(OUT),-out "$(OUT)",) $(if $(FORCE),-force,)
validate-site:
	$(call unavailable,validate-site,TQ-0049,run basic single-programme site quality checks)
publish-check:
	$(call unavailable,publish-check,TQ-0050 TQ-0049,verify offline Markdown eligibility first)
process:
	@test -n "$(PDF)" -a -n "$(URL)" -a -n "$(PROG)" || { \
		echo "make process PDF=<ruta> URL=<url oficial> PROG=<id> [TITLE= COVERAGE= RUN= ENDPOINT= MODEL= ATTEMPTS=]"; exit 2; }
	@test -n "$(or $(ENDPOINT),$(ATLAS_ENDPOINT))" -a -n "$(or $(MODEL),$(ATLAS_MODEL))" || { \
		echo "make process: falta ENDPOINT/MODEL (o ATLAS_ENDPOINT y ATLAS_MODEL)"; exit 2; }
	@go run ./cmd/manifesto-atlas process "$(PDF)" -url "$(URL)" -programme "$(PROG)" \
		-run "$(or $(RUN),runs/$(PROG))" -title "$(or $(TITLE),$(PROG))" -coverage "$(or $(COVERAGE),numbers)" \
		-attempts "$(or $(ATTEMPTS),3)" $(if $(ENDPOINT),-endpoint "$(ENDPOINT)",) $(if $(MODEL),-model "$(MODEL)",)
serve:
	$(call unavailable,serve,TQ-0065,preview an explicit private generated root locally)

# --- Shared recipes (single-line shells; no multi-line expansion) -------------

# Print a stage-not-available message and exit non-zero.
# $(1)=target  $(2)=owning task id(s)  $(3)=purpose
define unavailable
@printf '%s\n' "make $(1): not available yet." "  Owning task: $(2)." "  Purpose: $(3)." "  This target fails until $(2) lands. Run: tq show $(2)"; exit 1
endef

# Taskqueue integrity: tq list --json must exit 0, print non-empty JSON, and emit
# no warnings on stderr (tq names malformed or skipped task files on stderr while
# still exiting 0, so a warning-only failure must be rejected).
define tq-clean
@errfile=$$(mktemp); tq list --json >/dev/null 2>"$$errfile"; code=$$?; if [ $$code -ne 0 ]; then echo "tq list exited $$code:" >&2; cat "$$errfile" >&2; rm -f "$$errfile"; exit 1; fi; if [ -s "$$errfile" ]; then echo "tq list emitted warnings (malformed or skipped task files):" >&2; cat "$$errfile" >&2; rm -f "$$errfile"; exit 1; fi; rm -f "$$errfile"; echo "  taskqueue OK"
endef

# Go formatting: gofmt -l lists unformatted files; an empty result means all
# Go sources under cmd/ and internal/ are correctly formatted.
define go-fmt-check
@unformatted=$$(gofmt -l cmd internal); if [ -n "$$unformatted" ]; then echo "gofmt: unformatted files:"; echo "$$unformatted"; exit 1; fi; echo "  gofmt OK"
endef

# Documentation structure: the required tracked files must exist.
define docs-check
@missing=0; for f in AGENTS.md README.md HANDOFF.md LICENSE .taskqueue.yaml docs/SPEC.md docs/QUALITY_GATES.md docs/MAKEFILE.md; do if [ ! -f "$$f" ]; then echo "missing required file: $$f" >&2; missing=1; fi; done; if [ $$missing -ne 0 ]; then echo "docs-check: FAIL" >&2; exit 1; fi; echo "  docs OK"
endef

# Relative links: every non-URL, non-anchor markdown link in the tracked docs must
# resolve to an existing file or directory relative to the doc that holds it.
# The ]( prefix and ) suffix are stripped with shell parameter expansion, and URL
# and pure-anchor links are skipped, so no regex sed step is needed.
define link-check
@broken=0; for f in AGENTS.md README.md HANDOFF.md docs/SPEC.md docs/QUALITY_GATES.md docs/MAKEFILE.md; do [ -f "$$f" ] || continue; dir=$$(dirname "$$f"); for match in $$(grep -oE '\]\([^)]+\)' "$$f"); do link=$${match#']('}; link=$${link%')'}; [ -z "$$link" ] && continue; case "$$link" in http://*|https://*) continue ;; esac; target=$${link%%#*}; [ -z "$$target" ] && continue; if [ ! -e "$$dir/$$target" ]; then echo "broken link in $$f: $$link" >&2; broken=1; fi; done; done; if [ $$broken -ne 0 ]; then echo "link-check: FAIL" >&2; exit 1; fi; echo "  links OK"
endef
