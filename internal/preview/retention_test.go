package preview

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jparrill/manifesto-atlas/internal/semantic"
)

type previewTransport func(*http.Request) (*http.Response, error)

func (f previewTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestInvocationAttemptRetention(t *testing.T) {
	t.Run("When transport retries once, it should retain both attempted requests and original outcomes", func(t *testing.T) {
		count := 0
		var provider *semantic.OpenAIProvider
		inputs := map[string]string{"draft": hashBytes([]byte("draft"))}
		client := &http.Client{Transport: previewTransport(func(r *http.Request) (*http.Response, error) {
			count++
			if count == 1 {
				return nil, errors.New("controlled pre-response transport failure")
			}
			reply := semantic.Response{Version: semantic.JobVersion, JobID: "retry-test-retry1", ModelIdentity: provider.Identity(),
				Report: semantic.Report{Version: semantic.ReportVersion, Role: semantic.RoleClaimValidator, Subject: "source", Inputs: inputs,
					Result: semantic.ResultPass, Findings: []semantic.Finding{}, Checked: semantic.Checked{Statements: []string{}, Measures: []string{}, Pages: []int{}, Topics: []string{}}, Details: json.RawMessage(`{}`)}}
			raw, _ := json.Marshal(reply)
			envelope, _ := json.Marshal(map[string]any{"model": "test-model", "choices": []any{map[string]any{"message": map[string]string{"content": string(raw)}}}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(envelope))), Header: http.Header{}}, nil
		})}
		var err error
		provider, err = semantic.NewOpenAIProvider(semantic.ProviderConfig{BaseURL: "http://provider.invalid/v1", Model: "test-model", HTTPClient: client})
		if err != nil {
			t.Fatal(err)
		}
		_, err = semantic.NewRunner(provider).Run(t.Context(), semantic.Request{Version: semantic.JobVersion, JobID: "retry-test", EvaluationSet: "retention-test",
			Role: semantic.RoleClaimValidator, Subject: "source", Inputs: inputs, Model: "test-model"}, evidence{})
		if err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		files, err := persistInvocation(dir, semantic.RoleClaimValidator, provider)
		if err != nil || len(files) != 6 || count != 2 {
			t.Fatalf("attempt retention failed: files=%v calls=%d err=%v", files, count, err)
		}
		for i, want := range []string{"transient", "ok"} {
			prefix := semantic.RoleClaimValidator + "-attempt-" + string(rune('1'+i))
			raw, err := os.ReadFile(filepath.Join(dir, prefix+"-invocation.json"))
			if err != nil {
				t.Fatal(err)
			}
			var record map[string]any
			if err := json.Unmarshal(raw, &record); err != nil || record["status"] != want {
				t.Fatalf("original attempt outcome lost: %s", raw)
			}
			request, err := os.ReadFile(filepath.Join(dir, prefix+"-request.json"))
			if err != nil || len(request) == 0 {
				t.Fatal("original attempted request not retained")
			}
		}
	})
}

func TestRetainedPDFBudget(t *testing.T) {
	t.Run("When a PDF exceeds the JSON artifact ceiling, it should retain the documented PDF budget", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "source.pdf")
		f, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.Truncate(18 << 20); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := regularRead(path, retainedLimit("source.pdf")); err != nil {
			t.Fatalf("supported PDF rejected during retention: %v", err)
		}
		if _, err := regularRead(path, retainedLimit("programme.json")); err == nil {
			t.Fatal("JSON artifact limit was accidentally widened")
		}
	})
}
