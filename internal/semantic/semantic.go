// Package semantic implements the atlas.semantic-job/v1 runner for TQ-0024: a
// provider-independent, host-verified runner for the six schema-constrained
// skill jobs. It owns the bounded invocation, clean-context enforcement,
// response verification, limits, retry, and cancellation. It never grants a
// model job shell, filesystem write, network, repository edit, or signing
// capability; document text that says otherwise is quoted data.
//
// The runner is host authority: it verifies job binding, expected model
// identity, and the report schema against trusted inputs. A provider reply is
// untrusted until it passes every check; a malformed reply is quarantined
// diagnostic data, never an approval.
package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Contract versions.
const (
	JobVersion    = "atlas.semantic-job/v1"
	ReportVersion = "atlas.semantic-report/v1"
)

// The six skill role names.
const (
	RoleSourceProvenance  = "manifesto-source-provenance"
	RoleExtractionAudit   = "manifesto-extraction-audit"
	RoleNormalizer        = "manifesto-normalizer"
	RoleClaimValidator    = "manifesto-claim-validator"
	RoleCoverageValidator = "manifesto-coverage-validator"
	RoleSiteQA            = "manifesto-site-qa"
)

// Result and severity enums (shared with the artifact contract).
const (
	ResultPass         = "pass"
	ResultFail         = "fail"
	ResultInsufficient = "insufficient"
	ResultError        = "error"

	SeverityInfo    = "info"
	SeverityWarning = "warning"
	SeverityBlocker = "blocker"
)

// Runner limits (defaults; a request may only tighten within these ceilings).
const (
	maxResponseBytes = 1 << 20 // 1 MiB raw response
	maxOutputTokens  = 8192
	maxFindings      = 2000
	maxEvidence      = 32
)

// Error codes.
const (
	CodeBadVersion    = "bad-version"
	CodeUnknownRole   = "unknown-role"
	CodeMissingField  = "missing-field"
	CodeOverLimit     = "over-limit"
	CodeCleanContext  = "clean-context"
	CodeJobMismatch   = "job-mismatch"
	CodeModelMismatch = "model-mismatch"
	CodeBadReport     = "bad-report"
	CodeTimeout       = "timeout"
	CodeCancellation  = "cancellation"
	CodeProviderError = "provider-error"
)

// Error is a stable runner error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "semantic: " + e.Code
	if e.Subject != "" {
		s += " (" + e.Subject + ")"
	}
	if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

// Unwrap exposes the underlying error when present.
func (e *Error) Unwrap() error { return e.Err }

func errf(code, subject string) error { return &Error{Code: code, Subject: subject} }

// TransientError marks a pre-response transport error eligible for the single
// clean retry. It is not a semantic result.
type TransientError struct{ Err error }

func (e *TransientError) Error() string {
	if e.Err == nil {
		return "semantic: transient transport error"
	}
	return "semantic: transient: " + e.Err.Error()
}

func (e *TransientError) Unwrap() error { return e.Err }

// Evidence is one bounded evidence reference resolved to a preverified local ID.
type Evidence struct {
	ArtifactID string
	Page       *int
	MediaType  string
}

// Limits are request bounds that may only tighten the runner ceilings.
type Limits struct {
	Deadline         time.Duration
	MaxOutputTokens  int
	MaxResponseBytes int
}

// Request is one approved semantic job invocation.
type Request struct {
	Version       string
	JobID         string
	EvaluationSet string
	Role          string
	Subject       string
	Inputs        map[string]string
	Prompt        string
	OutputSchema  string
	Model         string
	Configuration map[string]string
	Evidence      []Evidence
	Limits        Limits
}

// ModelIdentity is the trusted, immutable model descriptor.
type ModelIdentity struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Revision string `json:"revision"`
}

// Finding is one bounded diagnostic.
type Finding struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Target   string `json:"target"`
	Evidence string `json:"evidence"`
	Message  string `json:"message"`
}

// Checked records the distinct IDs actually inspected.
type Checked struct {
	Statements []string `json:"statements"`
	Measures   []string `json:"measures"`
	Pages      []int    `json:"pages"`
	Topics     []string `json:"topics"`
}

// Report is the shared semantic report body.
type Report struct {
	Version  string            `json:"version"`
	Role     string            `json:"role"`
	Subject  string            `json:"subject"`
	Inputs   map[string]string `json:"inputs"`
	Result   string            `json:"result"`
	Findings []Finding         `json:"findings"`
	Checked  Checked           `json:"checked"`
	Details  json.RawMessage   `json:"details"`
}

// Response is the provider reply the host verifies.
type Response struct {
	Version       string        `json:"version"`
	JobID         string        `json:"job_id"`
	ModelIdentity ModelIdentity `json:"model_identity"`
	Report        Report        `json:"report"`
}

// EvidenceReader supplies bounded evidence bytes for a job. It resolves only
// preverified local artifact IDs; it never accepts paths, URLs, or
// model-requested references.
type EvidenceReader interface {
	Read(artifactID string, page *int) ([]byte, error)
}

// ProviderAdapter is a trusted provider implementation of the job contract. It
// has no shell, filesystem-write, network, repository-edit, or signing
// capability; it returns the raw provider reply bytes for host verification.
type ProviderAdapter interface {
	Identity() ModelIdentity
	Run(ctx context.Context, req Request, ev EvidenceReader) ([]byte, error)
}

// Runner is the host-side bounded, verified runner.
type Runner struct {
	adapter  ProviderAdapter
	usedJobs map[string]map[string]bool
}

// NewRunner binds a runner to a trusted provider adapter.
func NewRunner(adapter ProviderAdapter) *Runner {
	return &Runner{adapter: adapter, usedJobs: map[string]map[string]bool{}}
}

// Run performs one bounded, verified semantic job. It enforces clean context,
// the request deadline, the single clean retry, and full response verification.
// A failure yields no approval.
func (r *Runner) Run(ctx context.Context, req Request, ev EvidenceReader) (*Response, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}
	// Clean context: a reused job id within an evaluation set would leak prior
	// conversation state (the canary). Reject it.
	set := r.usedJobs[req.EvaluationSet]
	if set == nil {
		set = map[string]bool{}
		r.usedJobs[req.EvaluationSet] = set
	}
	if set[req.JobID] {
		return nil, errf(CodeCleanContext, req.JobID)
	}
	set[req.JobID] = true

	if req.Limits.Deadline > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.Limits.Deadline)
		defer cancel()
	}

	raw, err := r.adapter.Run(ctx, req, ev)
	if err != nil {
		if isTransient(err) {
			retryReq := req
			retryReq.JobID = req.JobID + "-retry1"
			set[retryReq.JobID] = true
			raw, err = r.adapter.Run(ctx, retryReq, ev)
			if err != nil {
				return nil, mapErr(ctx, err)
			}
			req = retryReq
		} else {
			return nil, mapErr(ctx, err)
		}
	}
	return r.verify(req, raw)
}

func (r *Runner) verify(req Request, raw []byte) (*Response, error) {
	if len(raw) > maxResponseBytes {
		return nil, errf(CodeOverLimit, "response")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var resp Response
	if err := dec.Decode(&resp); err != nil {
		return nil, errf(CodeBadReport, err.Error())
	}
	if resp.Version != JobVersion {
		return nil, errf(CodeBadVersion, resp.Version)
	}
	if resp.JobID != req.JobID {
		return nil, errf(CodeJobMismatch, resp.JobID)
	}
	if resp.ModelIdentity != r.adapter.Identity() {
		return nil, errf(CodeModelMismatch, resp.ModelIdentity.Model)
	}
	rep := resp.Report
	if rep.Version != ReportVersion {
		return nil, errf(CodeBadReport, rep.Version)
	}
	if rep.Role != req.Role {
		return nil, errf(CodeBadReport, "role")
	}
	if rep.Subject != req.Subject {
		return nil, errf(CodeBadReport, "subject")
	}
	if !equalStringMaps(rep.Inputs, req.Inputs) {
		return nil, errf(CodeBadReport, "inputs")
	}
	if !validResult(rep.Result) {
		return nil, errf(CodeBadReport, rep.Result)
	}
	if len(rep.Findings) > maxFindings {
		return nil, errf(CodeOverLimit, "findings")
	}
	seen := map[string]bool{}
	for _, f := range rep.Findings {
		if f.ID == "" || f.Code == "" || f.Severity == "" || f.Target == "" || f.Message == "" {
			return nil, errf(CodeBadReport, "finding")
		}
		if !validSeverity(f.Severity) {
			return nil, errf(CodeBadReport, f.Severity)
		}
		if seen[f.ID] {
			return nil, errf(CodeBadReport, "duplicate-finding")
		}
		seen[f.ID] = true
	}
	if len(rep.Details) > 0 {
		if _, err := artifact.Canonicalize(rep.Details); err != nil {
			return nil, errf(CodeBadReport, "details")
		}
	}
	return &resp, nil
}

func validateRequest(req Request) error {
	if req.Version != JobVersion {
		return errf(CodeBadVersion, req.Version)
	}
	if req.JobID == "" {
		return errf(CodeMissingField, "job_id")
	}
	if req.EvaluationSet == "" {
		return errf(CodeMissingField, "evaluation_set")
	}
	if !validRole(req.Role) {
		return errf(CodeUnknownRole, req.Role)
	}
	if req.Subject == "" {
		return errf(CodeMissingField, "subject")
	}
	if req.Model == "" {
		return errf(CodeMissingField, "model")
	}
	if len(req.Inputs) == 0 {
		return errf(CodeMissingField, "inputs")
	}
	if req.Limits.MaxOutputTokens > maxOutputTokens {
		return errf(CodeOverLimit, "tokens")
	}
	if req.Limits.MaxResponseBytes > maxResponseBytes {
		return errf(CodeOverLimit, "response")
	}
	if len(req.Evidence) > maxEvidence {
		return errf(CodeOverLimit, "evidence")
	}
	for _, e := range req.Evidence {
		if e.ArtifactID == "" {
			return errf(CodeMissingField, "evidence.artifact")
		}
	}
	return nil
}

func mapErr(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errf(CodeTimeout, "deadline")
		}
		return errf(CodeCancellation, "cancelled")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return errf(CodeTimeout, "deadline")
	}
	if errors.Is(err, context.Canceled) {
		return errf(CodeCancellation, "cancelled")
	}
	return errf(CodeProviderError, err.Error())
}

func isTransient(err error) bool {
	var te *TransientError
	return errors.As(err, &te)
}

func validRole(role string) bool {
	switch role {
	case RoleSourceProvenance, RoleExtractionAudit, RoleNormalizer,
		RoleClaimValidator, RoleCoverageValidator, RoleSiteQA:
		return true
	}
	return false
}

func validResult(r string) bool {
	switch r {
	case ResultPass, ResultFail, ResultInsufficient, ResultError:
		return true
	}
	return false
}

func validSeverity(s string) bool {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityBlocker:
		return true
	}
	return false
}

func equalStringMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
