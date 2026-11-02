# PDF submission and conversion isolation

> Historical contract. [MVP.md](MVP.md) and [SPEC.md](../SPEC.md) govern active
> work. TQ-0062 chooses one minimal bounded Fedora deployment; rootless Podman,
> custom OCI sealing and the full qualification apparatus below are not MVP prerequisites.

Contract: `atlas.isolation/v1`, established by TQ-0009 for TQ-0021 and the
extraction adapters. This is an implementation contract, not a claim that the
runtime has already been implemented or qualified on this host.

## Decisions and interfaces

Use rootless Podman with an OCI image pinned by SHA-256, crun, Linux namespaces,
seccomp, SELinux and delegated cgroup v2 limits. The first supported host profile
is Linux x86-64, kernel >= 6.1, Podman >= 5, crun >= 1.14, systemd user service,
SELinux enforcing, and working unprivileged user namespaces. macOS/Windows,
remote Podman connections, privileged containers, and hosts without any required
enforcement mechanism are unsupported. Never silently run a host converter.
The image must be present locally: `--pull=never`; no automatic image download.
Its actual digest and converter/model locks are recorded by the implementation,
not invented by this design.

`Inspect(candidateReader, submissionMetadata, approvedProfile) -> Inspection`
streams candidate bytes into private quarantine without parsing the PDF on the
host. `RunIsolated(registeredSource, adapter, approvedProfile) -> JobResult`
runs a fixed executable and argv inside a fresh sandbox. Both accept context
cancellation. Candidate bytes, paths, document text and metadata are data, never
shell snippets, image names, mount options, executable names, or limits.

`Inspection` records source hash/size, detected format, page count, encryption,
parser findings, profile digest, tool/image/runtime/kernel identity, enforcement
evidence, and a decision/error code. `JobResult` records bounded outputs,
resource/exit evidence and teardown status. Both use the
[artifact contract](artifact-contract.md). A successful inspection establishes
safe intake, not official provenance or extraction fidelity.

## Limits and override policy

All sizes are binary units. Host intake counts streamed bytes, stopping at
limit + 1; never trust Content-Length or PDF metadata. Page inspection must run
under the same isolation/resource controls as conversion.

| Resource | Default | Hard ceiling / enforcement |
| --- | --- | --- |
| Candidate PDF | 100 MiB | 256 MiB; counted bounded copy before parser |
| Pages | 512 | 2,048; parser page tree count before conversion |
| Inspection wall time | 60 seconds | 120 seconds; host monotonic deadline |
| Conversion wall time | 900 seconds | 1,800 seconds; host monotonic deadline |
| Per-job memory | 12 GiB | 16 GiB; `memory.max`, `memory.swap.max=0` |
| CPU bandwidth | 2 CPUs | 4 CPUs; `cpu.max`, never host-wide affinity assumptions |
| Process/thread count | 64 | 128; `pids.max`, includes descendants and threads |
| Open descriptors | 256 | 512; hard and soft `RLIMIT_NOFILE` |
| Writable scratch (`/work` + `/tmp`) | 2 GiB total | 4 GiB total; separate size-limited tmpfs allocations summing to the budget |
| Single output file | 64 MiB | 64 MiB; `RLIMIT_FSIZE`, plus post-run validation |
| Structured JSON file | 16 MiB | Artifact contract; smaller than single-file budget |
| Retained stdout/stderr | 4 MiB each | Drain/discard beyond bound; mark `output-limit`, never block a pipe |
| Rendered page | 150 DPI, 24 million pixels | Fixed profile; reject greater dimensions before image allocation |
| Concurrent jobs | 1 | 1 for v1; avoid multiplying memory/scratch limits |
| Teardown grace | 5 seconds TERM + 5 seconds KILL/reap | 10 seconds total; cleanup failure blocks admission |

The default permits the previously measured 172-page, approximately 7.5-GB,
seven-minute Docling trial without presenting that trial as fidelity approval.
The limits do not claim every nationwide document fits this profile.

Overrides are a maintainer-selected, checked-in profile resource, not flags or
fields supplied by the document/contributor/model. Every override must remain
within ceilings, is hashed as an input, and invalidates dependent artifacts.
Larger ceilings require a new reviewed contract version. No unlimited, zero,
negative, NaN, overflowed, or unsupported limit value is allowed.

## Intake, MIME and parser behavior

1. Open a caller-selected regular local file under a declared input root, using
   rooted/no-follow operations; reject symlinks in every path component,
   traversal, devices, sockets, directories and FIFOs. Read from that opened
   handle, not a later path lookup. Use an opaque generated quarantine filename;
   preserve the submitted name only as escaped metadata.
2. Stream into a mode-0700 job directory and mode-0600 immutable candidate copy,
   count size and hash bytes. Cancellation/oversize removes partial bytes.
3. Require `%PDF-1.0` through `%PDF-1.7` or `%PDF-2.0` at offset zero, a valid
   line terminator after the header, and host magic sniffing consistent with
   `application/pdf`. Header-only acceptance is forbidden. Reject conflicting
   declared MIME; absent HTTP MIME is allowed only when magic and parser agree.
4. In the sandbox use a pinned qpdf adapter (`--check` and JSON metadata, with
   argv constructed in Go). Reject encrypted documents even with empty-password
   access. Count the complete page tree: at least 1 page, within the profile.
   Corrupt/repaired/ambiguous documents and parser warnings (including qpdf exit
   3) yield `invalid-pdf`, never accepted conversion evidence. Do not repair or
   decrypt and silently substitute new source bytes.
5. Active content, embedded files, launch actions, external streams and remote
   references are not executed or followed. If the parser cannot inspect them
   within limits, fail closed. Page images never authorize external retrieval.

## Sandbox enforcement and evidence

The adapter constructs argv directly using `os/exec`, never `sh -c` or an
interpolated command. Mount only the private candidate/source copy at
`/input/source.pdf` read-only. Do not mount the checkout, home, SSH/GPG agents,
container socket, host `/proc`, model cache directory or credentials. All model
assets are preinstalled in the immutable image or separately pinned read-only
bundle. No runtime download/OCR-model bootstrap is allowed.

Required Podman settings include `--remote=false`, `--pull=never`,
`--network=none`, `--read-only`, `--cap-drop=ALL`,
`--security-opt=no-new-privileges`, non-root UID 10001, private PID/IPC/UTS/user
namespaces, `--pids-limit`, `--memory`, `--memory-swap` equal to memory,
`--cpus`, the descriptor/file-size ulimits, and bounded tmpfs mounts with
`nosuid,nodev,noexec`. Use the default seccomp allowlist with no unconfined
override and normal SELinux container confinement; private input copies may be
relabelled, registered sources must not be relabelled or changed in place.
Only minimal runtime devices such as `/dev/null`, zero and random are allowed.
Do not expose GPUs/devices in the first profile. Do not pass host environment;
use a fixed minimal PATH, locale and deterministic adapter variables.

At startup and after a runtime/profile change, a disposable qualification job
must prove the actual controls, not merely inspect CLI flags:

- Read back effective `memory.max`, zero `memory.swap.max`, `pids.max` and
  `cpu.max`; reject missing delegation, ignored limits, or Podman warnings.
- Probe UID, capabilities, `NoNewPrivs`, seccomp and enforcing SELinux label.
- Attempt writes to root/input/checkout paths (all denied or absent); create
  only scratch files and verify scratch capacity and per-file limits.
- Attempt DNS, public and loopback host connections and child process escape;
  no reachable host/service route is allowed. Container-internal loopback does
  not authorize an inference server or host networking.
- Run bounded timeout, fork/thread, scratch-fill and memory-pressure probes;
  confirm whole-cgroup enforcement and no surviving descendant.

Store probe results, effective limits, runtime/image/profile digests, process
exit status, OOM events and cleanup result as evidence. Requalify each new
session/profile/runtime/kernel; any failed/unknown probe yields
`isolation-unavailable` before processing a real PDF. Qualification never
requires public internet success, a real political source, or a model call.

The host owns the deadline and container ID independently of the worker process.
On timeout/cancel/error: stop input, send TERM, then KILL/remove the container
and entire cgroup, wait/reap descendants, and verify teardown. Container ID/name
is generated, never derived from a candidate filename. Bounded emergency cleanup
also runs after parent-process restart using recorded job IDs.

Keep the container alive using a fixed, minimal supervisor while the host runs
the adapter with `podman exec` and observes its exit status directly. The
deadline includes export. After the adapter exits cleanly, but **before**
container teardown destroys its tmpfs, copy only whitelisted regular output
files into host-private pending storage. Bound the copy and validate names,
sizes and content hashes; reject symlinks, hardlinks, special files, unexpected
outputs and hostile archive entries before installation. Do not extract a
worker-provided tar archive into the checkout. Only after successful export,
container teardown and payload validation may the host install derived
artifacts. Container outputs remain untrusted. Copying output is not permission
for workers to write the repository. Remove scratch on success and failure;
retain only the original
quarantined bytes (if complete), bounded diagnostics and failure evidence.
Cleanup failure leaves the job blocked and prevents further processing until
the stale resource is removed and qualification succeeds.

## Status and scenario-to-expected-outcome table

`accepted` advances intake only. `rejected` is a permanent invalid candidate;
`failed` is a bounded processing failure; `cancelled` is caller cancellation;
`blocked` means unavailable enforcement/cleanup. None except `accepted` can
advance to provenance. A converter exit zero is not an extraction approval.

| Scenario | Expected decision / error and cleanup |
| --- | --- |
| Valid native PDF within all limits | `accepted`; page count/hash recorded; no OCR required merely for intake |
| Branding or filename claims official origin | No provenance authority; independent official download still required |
| Spoofed MIME, wrong magic, HTML disguised as PDF | `rejected/mime-mismatch`; do not invoke converter |
| Missing HTTP MIME, valid magic and parser | Eligible for `accepted`; record absence, never infer official origin |
| Empty or header-only file | `rejected/invalid-pdf`; no accepted artifact |
| Oversize or byte-count overflow | `rejected/size-limit`; bounded copy stops, partial file removed |
| Zero pages or excess page count | `rejected/page-limit`; conversion never starts |
| Encrypted, corrupt, repaired or parser-warning PDF | `rejected/invalid-pdf`; immutable submitted bytes retained, no repair fallback |
| Huge streams, image dimensions or decompression bomb | `failed/resource-limit`; memory/time/pixel/scratch bounds apply before success |
| Inspection/conversion timeout | `failed/timeout`; kill/reap entire job, no partial approval |
| Memory exhaustion / cgroup OOM | `failed/memory-limit`; record OOM evidence, remove scratch |
| Fork/thread explosion, descriptor or output exhaustion | `failed/process-limit`, `descriptor-limit` or `output-limit`; teardown whole job |
| Attempted repository/source write | Denied/absent mount; qualification must prove this; no bytes modified |
| Attempted DNS, remote URL or host network access | Denied by network namespace; no network fallback |
| Candidate path contains spaces or shell metacharacters | One argv/data value; no shell interpolation |
| Symlink, traversal, FIFO, device or changed input path | `rejected/unsafe-input`; open-handle policy prevents replacement races |
| Malicious output symlink, hardlink or extra file | `failed/unsafe-output`; do not install any derived artifact |
| Caller cancellation | `cancelled`; bounded teardown and no running descendants |
| Unsupported host/runtime, missing image or limit ignored | `blocked/isolation-unavailable`; never run host/privileged fallback |
| Profile override within checked-in ceilings | New configuration identity; requalification and downstream invalidation |
| Contributor/document requests relaxed sandbox | Ignore as data; reject unauthorized configuration |
| Teardown fails or stale sandbox survives | `blocked/cleanup-failed`; halt further jobs until recovery |

## Versioning, test strategy and rejected alternatives

TQ-0021 must implement the above qualification and negative fixtures as a real
focused sandbox gate. Unit-test bounded copy, hostile paths, error mapping and
cancellation; integration-test the actual supported sandbox with synthetic
native/encrypted/corrupt/oversized/high-page-count PDFs and resource probes.
Tests must assert absence of surviving processes and repository/source changes.
Unavailability of required sandbox enforcement is a failed required gate, not a
skipped pass. Extraction tasks additionally run the real golden PDFs after their
verified recovery and independently audit conversion fidelity.

Policy/runtime/image/adapter changes create new configuration/tool identities
under the artifact contract. A later platform/GPU profile requires an explicit
versioned design with equivalent enforcement. Reject host-only `ulimit`, timeout
of just the parent process, containers with writable checkout mounts, privileged
Docker, and warning-only isolation checks: each leaves a specified boundary
unenforced. The supported narrow Linux profile is chosen over pretending all
platforms provide equivalent guarantees.

This design is consistent with [untrusted document handling](../SPEC.md#untrusted-document-handling)
and [submission and provenance tests](../SPEC.md#submission-and-provenance-tests).
It adds concrete enforcement choices without changing publication/provenance
requirements. The source registry and extraction contracts must consume these
results without treating intake success as semantic approval.
