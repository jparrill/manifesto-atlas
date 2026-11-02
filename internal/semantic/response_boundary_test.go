package semantic_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/semantic"
)

func TestStrictResponseBoundary(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "strict-job", "strict-set")
	valid := validReply(req, testIdentity, semantic.ResultPass)
	cases := map[string][]byte{
		"trailing object":   append(append([]byte{}, valid...), []byte(`{}`)...),
		"duplicate member":  []byte(strings.Replace(string(valid), `"result":"pass"`, `"result":"fail","result":"pass"`, 1)),
		"case alias":        []byte(strings.Replace(string(valid), `"job_id":`, `"JOB_ID":`, 1)),
		"unknown member":    []byte(strings.Replace(string(valid), `"job_id":`, `"extra":1,"job_id":`, 1)),
		"fractional number": []byte(strings.Replace(string(valid), `"details":{}`, `"details":{"number":1.5}`, 1)),
		"null details":      []byte(strings.Replace(string(valid), `"details":{}`, `"details":null`, 1)),
		"array details":     []byte(strings.Replace(string(valid), `"details":{}`, `"details":[]`, 1)),
		"null finding":      []byte(strings.Replace(string(valid), `"findings":[]`, `"findings":[null]`, 1)),
		"missing finding evidence": []byte(strings.Replace(string(valid), `"findings":[]`,
			`"findings":[{"id":"f-1","code":"fixture","severity":"info","target":"document","message":"Synthetic diagnostic"}]`, 1)),
		"alias finding evidence": []byte(strings.Replace(string(valid), `"findings":[]`,
			`"findings":[{"id":"f-1","code":"fixture","severity":"info","target":"document","EVIDENCE":"fixture","message":"Synthetic diagnostic"}]`, 1)),
	}
	// Every common object member must be present, non-null, and exactly named.
	objects := map[string][]string{
		"":               {"version", "job_id", "model_identity", "report"},
		"model_identity": {"provider", "model", "revision"},
		"report":         {"version", "role", "subject", "inputs", "result", "findings", "checked", "details"},
		"report.checked": {"statements", "measures", "pages", "topics"},
	}
	for path, keys := range objects {
		for _, key := range keys {
			for _, mutation := range []string{"missing", "null", "alias"} {
				var object map[string]any
				if err := json.Unmarshal(valid, &object); err != nil {
					t.Fatal(err)
				}
				nested := object
				if path != "" {
					for _, part := range strings.Split(path, ".") {
						nested = nested[part].(map[string]any)
					}
				}
				switch mutation {
				case "missing":
					delete(nested, key)
				case "null":
					nested[key] = nil
				case "alias":
					nested[strings.ToUpper(key)] = nested[key]
					delete(nested, key)
				}
				raw, err := json.Marshal(object)
				if err != nil {
					t.Fatal(err)
				}
				cases[path+"/"+key+"/"+mutation] = raw
			}
		}
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			a := &fakeAdapter{identity: testIdentity, reply: raw}
			resp, err := semantic.NewRunner(a).Run(context.Background(), req, stubReader{})
			if resp != nil || code(err) != semantic.CodeBadReport || len(a.calls) != 1 {
				t.Fatalf("malformed reply: response=%v error=%v calls=%d", resp, err, len(a.calls))
			}
		})
	}
}

func TestReportDecisionConsistency(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "decision-job", "decision-set")
	for _, result := range []string{semantic.ResultPass, semantic.ResultFail, semantic.ResultInsufficient, semantic.ResultError} {
		for _, severity := range []string{semantic.SeverityInfo, semantic.SeverityWarning, semantic.SeverityBlocker} {
			t.Run(result+"/"+severity, func(t *testing.T) {
				var reply semantic.Response
				if err := json.Unmarshal(validReply(req, testIdentity, result), &reply); err != nil {
					t.Fatal(err)
				}
				reply.Report.Findings = []semantic.Finding{{ID: "f-1", Code: "fixture", Severity: severity,
					Target: "document", Evidence: "fixture-evidence", Message: "Synthetic diagnostic"}}
				raw, err := json.Marshal(reply)
				if err != nil {
					t.Fatal(err)
				}
				a := &fakeAdapter{identity: testIdentity, reply: raw}
				resp, err := semantic.NewRunner(a).Run(context.Background(), req, stubReader{})
				if result == semantic.ResultPass && severity != semantic.SeverityInfo {
					if resp != nil || code(err) != semantic.CodeBadReport {
						t.Fatalf("contradictory pass accepted: %v, %v", resp, err)
					}
				} else if err != nil || resp == nil || resp.Report.Result != result {
					t.Fatalf("valid decision not retained: %v, %v", resp, err)
				}
				if len(a.calls) != 1 {
					t.Fatal("semantic finding caused an automatic retry")
				}
			})
		}
	}
}

func TestTightenedResponseLimit(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "bytes-job", "bytes-set")
	raw := validReply(req, testIdentity, semantic.ResultPass)
	for _, limit := range []int{1, len(raw) - 1, len(raw)} {
		req.Limits.MaxResponseBytes = limit
		a := &fakeAdapter{identity: testIdentity, reply: raw}
		resp, err := semantic.NewRunner(a).Run(context.Background(), req, stubReader{})
		if limit < len(raw) {
			if resp != nil || code(err) != semantic.CodeOverLimit {
				t.Fatalf("limit %d not enforced: %v, %v", limit, resp, err)
			}
		} else if err != nil || resp == nil {
			t.Fatalf("exact byte limit rejected: %v", err)
		}
	}
	// Raw output must be bounded before parsing, including malformed output.
	req.Limits.MaxResponseBytes = 0
	a := &fakeAdapter{identity: testIdentity, reply: []byte(strings.Repeat("x", (1<<20)+1))}
	resp, err := semantic.NewRunner(a).Run(context.Background(), req, stubReader{})
	if resp != nil || code(err) != semantic.CodeOverLimit || len(a.calls) != 1 {
		t.Fatalf("default raw byte ceiling not enforced: %v, %v", resp, err)
	}
}

type boundaryAdapter struct {
	run func(context.Context, semantic.Request) ([]byte, error)
}

func (boundaryAdapter) Identity() semantic.ModelIdentity { return testIdentity }
func (a boundaryAdapter) Run(ctx context.Context, req semantic.Request, _ semantic.EvidenceReader) ([]byte, error) {
	return a.run(ctx, req)
}

func TestTerminalInvocationBoundary(t *testing.T) {
	for _, mode := range []string{"pre-cancel", "late-cancel", "late-deadline", "cancelled-transient", "partial-transient"} {
		t.Run(mode, func(t *testing.T) {
			req := baseReq(semantic.RoleClaimValidator, "terminal-job", "terminal-set")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			want, wantCalls := semantic.CodeCancellation, 1
			switch mode {
			case "pre-cancel":
				cancel()
				wantCalls = 0
			case "late-deadline":
				req.Limits.Deadline = time.Millisecond
				want = semantic.CodeTimeout
			case "partial-transient":
				want = semantic.CodeProviderError
			}
			calls := 0
			a := boundaryAdapter{run: func(ctx context.Context, invocation semantic.Request) ([]byte, error) {
				calls++
				if mode == "late-deadline" {
					<-ctx.Done()
				} else if mode != "partial-transient" {
					cancel()
				}
				if mode == "cancelled-transient" {
					return nil, &semantic.TransientError{Err: errors.New("synthetic transient")}
				}
				if mode == "partial-transient" {
					return []byte("partial response"), &semantic.TransientError{Err: errors.New("synthetic transient")}
				}
				return validReply(invocation, testIdentity, semantic.ResultPass), nil
			}}
			resp, err := semantic.NewRunner(a).Run(ctx, req, stubReader{})
			if resp != nil || code(err) != want || calls != wantCalls {
				t.Fatalf("terminal boundary: response=%v error=%v calls=%d; want %s/%d", resp, err, calls, want, wantCalls)
			}
		})
	}
}

func TestInvocationLimitValidation(t *testing.T) {
	for _, mode := range []string{"negative deadline", "negative bytes", "negative tokens", "excess deadline", "default limits"} {
		t.Run(mode, func(t *testing.T) {
			req := baseReq(semantic.RoleClaimValidator, "limits-job", "limits-set")
			switch mode {
			case "negative deadline":
				req.Limits.Deadline = -time.Second
			case "negative bytes":
				req.Limits.MaxResponseBytes = -1
			case "negative tokens":
				req.Limits.MaxOutputTokens = -1
			case "excess deadline":
				req.Limits.Deadline = 601 * time.Second
			case "default limits":
				req.Limits = semantic.Limits{}
			}
			calls := 0
			a := boundaryAdapter{run: func(ctx context.Context, invocation semantic.Request) ([]byte, error) {
				calls++
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) > 600*time.Second || invocation.Limits.MaxResponseBytes != 1<<20 || invocation.Limits.MaxOutputTokens != 8192 {
					t.Error("default invocation was not bounded")
				}
				return validReply(invocation, testIdentity, semantic.ResultPass), nil
			}}
			resp, err := semantic.NewRunner(a).Run(context.Background(), req, stubReader{})
			if mode == "default limits" {
				if resp == nil || err != nil || calls != 1 {
					t.Fatalf("valid defaults rejected: %v", err)
				}
			} else if resp != nil || code(err) != semantic.CodeOverLimit || calls != 0 {
				t.Fatalf("invalid limit reached adapter: %v, %v, calls=%d", resp, err, calls)
			}
		})
	}
}

func TestRetryDoesNotReuseExistingJob(t *testing.T) {
	calls := 0
	a := boundaryAdapter{run: func(_ context.Context, req semantic.Request) ([]byte, error) {
		calls++
		if calls == 1 {
			return validReply(req, testIdentity, semantic.ResultPass), nil
		}
		return nil, &semantic.TransientError{Err: errors.New("synthetic transient")}
	}}
	runner := semantic.NewRunner(a)
	reserved := baseReq(semantic.RoleClaimValidator, "job-retry1", "retry-set")
	if _, err := runner.Run(context.Background(), reserved, stubReader{}); err != nil {
		t.Fatal(err)
	}
	req := baseReq(semantic.RoleClaimValidator, "job", "retry-set")
	resp, err := runner.Run(context.Background(), req, stubReader{})
	if resp != nil || code(err) != semantic.CodeCleanContext || calls != 2 {
		t.Fatalf("retry reused a job: response=%v error=%v calls=%d", resp, err, calls)
	}
}
