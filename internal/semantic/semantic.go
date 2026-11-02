// Package semantic provides common host-side response validation for
// atlas.semantic-job/v1. It checks job/model binding, strict common envelopes,
// decision consistency, response limits and terminal cancellation. Stage-specific
// schemas, evidence inventories, configured providers and runtime/session
// isolation remain required work in TQ-0024; common validation is not approval.
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
	maxRequestTime   = 600 * time.Second
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

// Limits are request bounds that may only tighten the runner ceilings. Zero
// selects the approved default (600 seconds, 8192 tokens, 1 MiB); negative values
// are invalid. A provider adapter must enforce generation token/stream limits.
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

// ProviderAdapter returns untrusted provider reply bytes for host verification.
// No model tool API is exposed here; implementations must separately establish
// the approved broker capabilities and provider/session isolation. Implementations
// must honor context cancellation and bound response allocation while streaming.
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
	if ctx.Err() != nil {
		return nil, mapErr(ctx, ctx.Err())
	}
	if req.Limits.Deadline == 0 {
		req.Limits.Deadline = maxRequestTime
	}
	if req.Limits.MaxResponseBytes == 0 {
		req.Limits.MaxResponseBytes = maxResponseBytes
	}
	if req.Limits.MaxOutputTokens == 0 {
		req.Limits.MaxOutputTokens = maxOutputTokens
	}
	// Reject job reuse within this runner/evaluation set. This bookkeeping does
	// not establish provider conversation or KV-cache isolation.
	set := r.usedJobs[req.EvaluationSet]
	if set == nil {
		set = map[string]bool{}
		r.usedJobs[req.EvaluationSet] = set
	}
	if set[req.JobID] {
		return nil, errf(CodeCleanContext, req.JobID)
	}
	set[req.JobID] = true

	ctx, cancel := context.WithTimeout(ctx, req.Limits.Deadline)
	defer cancel()
	if ctx.Err() != nil {
		return nil, mapErr(ctx, ctx.Err())
	}

	raw, err := r.adapter.Run(ctx, req, ev)
	if ctx.Err() != nil {
		return nil, mapErr(ctx, ctx.Err())
	}
	if err != nil {
		if isTransient(err) && len(raw) == 0 {
			retryReq := req
			retryReq.JobID = req.JobID + "-retry1"
			if set[retryReq.JobID] {
				return nil, errf(CodeCleanContext, retryReq.JobID)
			}
			set[retryReq.JobID] = true
			raw, err = r.adapter.Run(ctx, retryReq, ev)
			if ctx.Err() != nil {
				return nil, mapErr(ctx, ctx.Err())
			}
			if err != nil {
				return nil, mapErr(ctx, err)
			}
			req = retryReq
		} else {
			return nil, mapErr(ctx, err)
		}
	}
	resp, err := r.verify(req, raw)
	if ctx.Err() != nil {
		return nil, mapErr(ctx, ctx.Err())
	}
	return resp, err
}

func (r *Runner) verify(req Request, raw []byte) (*Response, error) {
	if len(raw) > req.Limits.MaxResponseBytes {
		return nil, errf(CodeOverLimit, "response")
	}
	canon, err := artifact.Canonicalize(raw)
	if err != nil {
		return nil, &Error{Code: CodeBadReport, Subject: "strict-json", Err: err}
	}
	if err := validateEnvelope(canon); err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(canon))
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
		if rep.Result == ResultPass && f.Severity != SeverityInfo {
			return nil, errf(CodeBadReport, "pass-with-unresolved-finding")
		}
		if seen[f.ID] {
			return nil, errf(CodeBadReport, "duplicate-finding")
		}
		seen[f.ID] = true
	}
	return &resp, nil
}

// validateEnvelope checks exact common field names/presence before the standard
// decoder (which otherwise accepts case aliases and missing/null fields).
// Details must be an object, but its skill-specific schema remains the caller's
// responsibility until the output-schema resolver is implemented in TQ-0024.
func validateEnvelope(raw []byte) error {
	root, err := objectFields(raw, "version", "job_id", "model_identity", "report")
	if err != nil {
		return err
	}
	if _, err := objectFields(root["model_identity"], "provider", "model", "revision"); err != nil {
		return err
	}
	report, err := objectFields(root["report"], "version", "role", "subject", "inputs", "result", "findings", "checked", "details")
	if err != nil {
		return err
	}
	if _, err := objectFields(report["checked"], "statements", "measures", "pages", "topics"); err != nil {
		return err
	}
	if _, err := objectFields(report["details"]); err != nil {
		return err
	}
	var findings []json.RawMessage
	if err := json.Unmarshal(report["findings"], &findings); err != nil {
		return errf(CodeBadReport, "findings")
	}
	if len(findings) > maxFindings {
		return errf(CodeOverLimit, "findings")
	}
	for _, finding := range findings {
		if _, err := objectFields(finding, "id", "code", "severity", "target", "evidence", "message"); err != nil {
			return err
		}
	}
	return nil
}

func objectFields(raw []byte, fields ...string) (map[string]json.RawMessage, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil || object == nil {
		return nil, errf(CodeBadReport, "expected-object")
	}
	if len(fields) != 0 && len(object) != len(fields) {
		return nil, errf(CodeBadReport, "object-fields")
	}
	for _, field := range fields {
		value, ok := object[field]
		if !ok || bytes.Equal(value, []byte("null")) {
			return nil, errf(CodeBadReport, field)
		}
	}
	return object, nil
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
	if req.Limits.Deadline < 0 || req.Limits.Deadline > maxRequestTime {
		return errf(CodeOverLimit, "deadline")
	}
	if req.Limits.MaxOutputTokens < 0 || req.Limits.MaxOutputTokens > maxOutputTokens {
		return errf(CodeOverLimit, "tokens")
	}
	if req.Limits.MaxResponseBytes < 0 || req.Limits.MaxResponseBytes > maxResponseBytes {
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
