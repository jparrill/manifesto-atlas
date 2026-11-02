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
.PHONY: help check verify test-cli test-artifacts test-attestations test-taxonomy test-canonical test-party-registry test-fetch

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
	  "  make test-artifacts        Run the focused artifact envelope/graph test suite." \
	  "  make test-attestations     Run the focused attestation signing/verification suite." \
	  "  make test-taxonomy         Run the focused taxonomy/jurisdiction validation suite." \
	  "  make test-canonical        Run the focused canonical programme/citation validation suite." \
	  "  make test-party-registry   Run the focused party trust-anchor/registry validation suite." \
	  "  make test-fetch            Run the focused bounded safe HTTP transport (fetch) suite." \
	  "" \
	  "Pipeline stages (planned — fail until their owning task lands):" \
	  "  make submit                TQ-0021  Quarantine candidate bytes and metadata." \
	  "  make verify-source         TQ-0022 TQ-0026  Provenance and canonical bytes." \
	  "  make add                   TQ-0027  Register a verified immutable source." \
	  "  make extract               TQ-0030  Docling page-aware extraction." \
	  "  make validate-extraction   TQ-0035  Extraction fidelity approval." \
	  "  make normalize             TQ-0038  Draft canonical programme JSON." \
	  "  make validate-content      TQ-0041  Claim and coverage validation." \
	  "  make build                 TQ-0043 TQ-0044  Deterministic Markdown and HTML." \
	  "  make validate-site         TQ-0049  Site validation gate." \
	  "  make publish-check         TQ-0050  Publishable artifact graph." \
	  "  make process               TQ-0051  Orchestrate programme processing." \
	  "  make serve                 TQ-0051  Local site preview." \
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
	@echo "verify: go test (test-attestations)"
	@$(MAKE) test-attestations
	@echo "verify: go test (test-taxonomy)"
	@$(MAKE) test-taxonomy
	@echo "verify: go test (test-canonical)"
	@$(MAKE) test-canonical
	@echo "verify: go test (test-party-registry)"
	@$(MAKE) test-party-registry
	@echo "verify: go test (test-fetch)"
	@$(MAKE) test-fetch
	@echo "verify: OK"

# Focused gate for the Go CLI test suite. Runs the real tests without invoking
# live models or requiring live official websites.
test-cli:
	@go test ./internal/cli/...

# Focused gate for the atlas.artifact/v1 envelope and graph validator. Exercises
# canonicalization, hashing, manifest/graph validation, and invalidation fixtures
# offline; no live models or official websites.
test-artifacts:
	@go test ./internal/artifact/...

# Focused gate for the atlas.attestation/v1 and atlas.trust/v1 contracts.
# Exercises Ed25519 signing, offline verification, role/rotation/revocation, and
# trust-bundle validation with test-only keys; no live models or official websites.
test-attestations:
	@go test ./internal/attestation/...

# Focused gate for the atlas.taxonomy/v1 registry, topic states, jurisdiction
# profiles, classification, comparison compatibility, and migration validation.
# Deterministic and offline; no live models or official websites.
test-taxonomy:
	@go test ./internal/taxonomy/...

# Focused gate for the atlas.programme/v1 canonical schema, topic/state rules,
# citation/page/excerpt evidence, source binding, and deterministic bytes.
# Deterministic and offline; no live models or official websites.
test-canonical:
	@go test ./internal/canonical/...

# Focused gate for the atlas.provenance/v1 party/territory/domain evidence
# registry, fail-closed anchor authority checks, and version invalidation.
# Deterministic and offline; fixture anchors only, no live official websites.
test-party-registry:
	@go test ./internal/provenance/...

# Focused gate for the atlas.provenance/v1 bounded safe HTTP transport: URL/
# DNS/IP (SSRF) policy, redirect revalidation, bounded downloads, and content
# inspection. Deterministic and offline; synthetic local TLS fixtures with an
# injected resolver/dialer, no live official websites.
test-fetch:
	@go test ./internal/provenance/...

# --- Planned pipeline stages (fail until their owning task lands) -------------

submit:
	$(call unavailable,submit,TQ-0021,quarantine candidate bytes and submission metadata)
verify-source:
	$(call unavailable,verify-source,TQ-0022 TQ-0026,establish provenance and retrieve canonical bytes)
add:
	$(call unavailable,add,TQ-0027,register a verified immutable source)
extract:
	$(call unavailable,extract,TQ-0030,produce structured page-aware extraction)
validate-extraction:
	$(call unavailable,validate-extraction,TQ-0035,audit extraction fidelity)
normalize:
	$(call unavailable,normalize,TQ-0038,produce draft canonical programme JSON)
validate-content:
	$(call unavailable,validate-content,TQ-0041,run claim and coverage validators)
build:
	$(call unavailable,build,TQ-0043 TQ-0044,render Spanish Markdown and HTML deterministically)
validate-site:
	$(call unavailable,validate-site,TQ-0049,run content parity, HTML, accessibility, responsive, and visual checks)
publish-check:
	$(call unavailable,publish-check,TQ-0050,verify the complete artifact graph is publishable)
process:
	$(call unavailable,process,TQ-0051,orchestrate the applicable processing stages)
serve:
	$(call unavailable,serve,TQ-0051,preview the generated static website locally)

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
