# Manifesto Atlas — Make orchestration interface.
#
# Baseline targets are implemented now (help, check, verify). The 12 pipeline
# stages from docs/SPEC.md#make-targets are declared but fail until their owning
# backlog task lands. Target-to-owner mapping and integration notes: docs/MAKEFILE.md
#
# Orchestration lives here. Deterministic work moves to Go once the module exists
# (owned by TQ-0018); these baseline checks only use Git and Taskqueue.

.DEFAULT_GOAL := help
SHELL := /bin/sh

# Baseline targets (implemented).
.PHONY: help check verify test-cli test-artifacts test-taxonomy test-canonical test-semantic-runner test-fixture-frente test-runs

# Pipeline stages (planned): fail until their owning task lands.
.PHONY: submit verify-source add extract validate-extraction normalize validate-content build validate-site publish-check process serve

help:
	@printf '%s\n' \
	  "Manifesto Atlas — Make targets" \
	  "" \
	  "Baseline (implemented):" \
	  "  make help                  Show this help." \
	  "  make check                 Run baseline checks (Git + Taskqueue + Go)." \
	  "  make verify                Run full verification (check + docs + links + tests)." \
	  "  make test-cli              Run the focused Go CLI test suite." \
	  "  make test-artifacts        Run the focused strict canonical-JSON/hash helper suite." \
	  "  make test-taxonomy         Run the focused taxonomy/document-scope validation suite." \
	  "  make test-canonical        Run the focused canonical programme/citation validation suite." \
	  "  make test-semantic-runner  Run the focused semantic job runner (isolation/provider) suite." \
	  "  make test-fixture-frente   Run the focused Frente Obrero candidate/source fixture suite." \
	  "  make test-runs             Run the focused selection/run record validation suite." \
	  "" \
	  "Pipeline stages (planned — fail until their owning task lands):" \
	  "  make submit                TQ-0021  Inspect a selected local candidate." \
	  "  make verify-source         TQ-0027  Verify local selection evidence and bytes." \
	  "  make add                   TQ-0027  Register a verified immutable source." \
	  "  make extract               TQ-0030  Page-aware extraction." \
	  "  make validate-extraction   TQ-0030  Deterministic page/warning blockers." \
	  "  make normalize             TQ-0038  Draft canonical programme JSON." \
	  "  make validate-content      TQ-0041  Claim and coverage validation." \
	  "  make build                 TQ-0043 TQ-0044 TQ-0049  Markdown, previews, HTML." \
	  "  make validate-site         TQ-0049  Site validation gate." \
	  "  make publish-check         TQ-0050 TQ-0049  Offline Markdown eligibility first." \
	  "  make process               TQ-0051  Orchestrate programme processing." \
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
	@echo "verify: go test (test-runs)"
	@$(MAKE) test-runs
	@echo "verify: OK"

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

# Focused gate for the atlas.selection/v1 owner selection record and the
# atlas.run/v1 per-programme run record: strict parsing, deterministic
# round-trip, and offline bounded-file hash verification. No live models or
# official websites.
test-runs:
	@go test ./internal/run/...

# --- Planned pipeline stages (fail until their owning task lands) -------------

submit:
	$(call unavailable,submit,TQ-0021,inspect a selected local candidate)
verify-source:
	$(call unavailable,verify-source,TQ-0027,verify local selection evidence and matching bytes)
add:
	$(call unavailable,add,TQ-0027,register a verified immutable source)
extract:
	$(call unavailable,extract,TQ-0030,produce structured page-aware extraction)
validate-extraction:
	$(call unavailable,validate-extraction,TQ-0030,reject deterministic page/warning blockers)
normalize:
	$(call unavailable,normalize,TQ-0038,produce draft canonical programme JSON)
validate-content:
	$(call unavailable,validate-content,TQ-0041,run claim and coverage validators)
build:
	$(call unavailable,build,TQ-0043 TQ-0044 TQ-0049,render Spanish Markdown, previews, and HTML)
validate-site:
	$(call unavailable,validate-site,TQ-0049,run basic single-programme site quality checks)
publish-check:
	$(call unavailable,publish-check,TQ-0050 TQ-0049,verify offline Markdown eligibility first)
process:
	$(call unavailable,process,TQ-0051,orchestrate the applicable processing stages)
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
