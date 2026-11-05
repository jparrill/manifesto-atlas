# macOS real-source preview

Owner-requested scope, TQ-0062/TQ-0067: demonstrate the real local process and
its HTML on macOS now. Fedora is not a prerequisite for this preview. The
production registration/publication tasks remain separately tracked; this seam
does not retrospectively complete them or waive independent content validation.

## Runtime

Use installed Poppler `pdfinfo` and `pdftotext` for the initial native-text path,
not an uninstalled Docling adapter. Record actual versions and preserve the raw
page-separated text. Verify the retained candidate/evidence hashes first. Normalize
whitespace within each original page, keep contiguous one-based page indices, and
bind citations to exact Unicode slices of that retained normalized evidence.
The host labels overlapping 40-word source slices with stable page/evidence IDs;
the normalizer selects those IDs instead of retyping source quotations. Unknown
IDs reject, and the existing canonical validator checks the resolved literal slices.
Empty/native-unreadable pages are explicit warnings, never invented OCR text.

Workers run through `/usr/bin/sandbox-exec`: deny by default, allow executable/
library/font reads, the one selected PDF and process primitives, deny general
network and repository/home writes. Input is read-only. The host owns bounded
stdout/stderr capture, a deadline and process-group kill on failure/cancellation.
Workers receive a minimal environment without host credentials/SSH agent sockets;
host process inspection and forking are denied in the default native PDF policy.
Monitor process-group RSS and terminate the group above its configured budget.
This is a **sampled RSS ceiling**, not a kernel hard memory limit: a short-lived
overshoot can occur between samples. Do not advertise it as Fedora/cgroup parity.
Bounds: candidate 100 MiB, 500 pages, extraction output 8 MiB, each worker 120 s,
process-group RSS 1 GiB sampled every 50 ms. Tests tighten these bounds.

Observed on this host: Darwin/arm64, Poppler and sandbox-exec available; a synthetic
sandboxed `/usr/bin/true` succeeded. Python's attempt to lower `RLIMIT_AS` failed
with `ValueError: current limit exceeds maximum limit`; do not silently rely on
that mechanism. Actual denial, cleanup, RSS and output-limit probes are required
implementation evidence in TQ-0067, not claims made by this design document.

## Semantic path

Use Auriga through the existing trusted SSH configuration and a loopback tunnel.
Endpoint/model are operator configuration, never source-controlled. Credentials
stay out of artifacts. Three stateless requests use the configured semantic
provider: normalizer, claim-validator, coverage-validator. The validators receive
the complete original page evidence and canonical draft, not generator reasoning.
Record prompts, input hashes and original bounded replies. Host checks canonical
structure, exact quotes/page references and complete checked inventories. Retain
fail/insufficient/error outcomes; no retry-until-pass or fabricated approval.

This preview summarizes principal measures; coverage may legitimately reject an
incomplete concise draft. Its outcome remains visible. Private rendering is
allowed for structurally valid data even when either semantic check fails. It
never creates a public catalogue or asserts production publication eligibility.

## Supplied presentation, not LLM-authored markup

Juan supplied `Programas-electorales-2026.html` from the earlier programme folder.
Import it byte-for-byte as an immutable reference, retaining its checksum. Keep
all HTML/CSS/JavaScript unchanged as input. Go binds only the existing `parties`
data literal using canonical data; it does not author layout, script or templates.
HTML-escape every interpolated string before JSON serialization because the
supplied script uses `innerHTML`. Preserve safe, generated source/Markdown links.
Unknown reference structure fails rather than falling back to generated markup.
Unprocessed party tabs contain explicit pending state, not the old summaries.

## Commands and retained evidence

```bash
go run ./cmd/preview -import-reference /path/to/Programas-electorales-2026.html
ssh -fN -L 127.0.0.1:18790:127.0.0.1:8090 auriga
ATLAS_BASE_URL=http://127.0.0.1:18790/v1 \
ATLAS_MODEL=qwen3.8-27b-dflash2-q4 make e2e-preview OUTPUT=runs/preview/new-run
make render-preview
go run ./cmd/preview -check-browser -output runs/preview/frente-e2e
open runs/preview/frente-e2e/index.html
```

The reference is stored under `assets/reference/`. Each explicit run retains raw
extraction, normalized pages, canonical JSON, role requests/replies/reports,
model/tool identities, source checksum, Markdown, HTML and a run record under
`runs/preview/`. New evaluations require a new output directory; previous attempts
are never overwritten. `make render-preview` defaults to the retained real
`frente-e2e` run. Model invocations are explicit and excluded from `make verify`.
Rerender consumes stored canonical data and reports without model calls; input
hashes must match or the stored reports are stale. The output always identifies
itself as a private preview and exposes actual report outcomes.

Browser evidence must exercise the actual supplied script and generated data at
desktop and mobile sizes. Production HTML acceptance remains its own task; the
existing human/non-LLM presentation rule is preserved, not bypassed.

## Retained first execution

`runs/preview/frente-e2e/` contains the genuine 64-page Frente Obrero execution,
21 ordered topics, 17 principal measures, exact original-page citations, and
actual separate Auriga normalization/claim/coverage requests and replies.
Both host outcomes are insufficient: native extraction has no text for pages
1, 2 and 64, and original visual evidence is required. Coverage also reports
potentially material omitted measures; retained findings are evidence of the
validator's actual assessment, not infallible judgments. The preview clearly says
draft/private; it is not publicly eligible. Earlier failed sandbox, cancelled
proxy and unsupported-quotation attempts remain separately recorded, as does a
first draft whose validator replies omitted mandatory finding fields.

The browser check uses Chrome's debugging pipe and CDP device emulation, not
generated JavaScript. It checks 21 actually rendered cards and horizontal layout
at 1440×1000 and 390×844, retaining screenshots and their hashes. Ordinary offline
tests replay the stored run, check exact data parity/determinism and reject stale
or missing source/model/draft/report bytes without calling Auriga.

To repeat actual native extraction independently of the model, run
`ATLAS_EVALUATE_REAL_PDF=1 go test ./internal/preview -run TestActualMacSource -count=1 -v`.
It asserts exact equality of all original pages, source identity, warnings and
actual tools against the retained run. The explicit real-tool gate is opt-in;
ordinary regression replays stored evidence offline.
