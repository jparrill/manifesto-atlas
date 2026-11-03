package semantic_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/semantic"
)

// openaiReply wraps a model reply (the atlas.semantic-job/v1 envelope) in an
// OpenAI chat.completion transport envelope.
func openaiReply(model string, reply []byte) []byte {
	content, _ := json.Marshal(string(reply))
	return []byte(fmt.Sprintf(
		`{"id":"chatcmpl-test","model":%q,"choices":[{"index":0,"finish_reason":"stop","message":{"role":"assistant","content":%s}}]}`,
		model, content))
}

// newOpenAIProvider starts a controlled httptest OpenAI endpoint and a configured
// adapter. The handler receives the provider (assigned before any request), so it
// can build replies that echo the adapter's trusted identity.
func newOpenAIProvider(t *testing.T, handler func(http.ResponseWriter, *http.Request, *semantic.OpenAIProvider), model, apiKey string) *semantic.OpenAIProvider {
	t.Helper()
	var p *semantic.OpenAIProvider
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, p)
	}))
	t.Cleanup(srv.Close)
	var err error
	p, err = semantic.NewOpenAIProvider(semantic.ProviderConfig{
		BaseURL:  srv.URL + "/v1",
		Model:    model,
		APIKey:   apiKey,
		Revision: "rev-test",
	})
	if err != nil {
		t.Fatalf("NewOpenAIProvider: %v", err)
	}
	return p
}

type mapReader map[string][]byte

func (m mapReader) Read(id string, page *int) ([]byte, error) {
	b, ok := m[id]
	if !ok {
		return nil, errors.New("missing evidence: " + id)
	}
	return b, nil
}

func TestRealAdapterTransport(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	var gotAuth, gotModel string
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		gotAuth = r.Header.Get("Authorization")
		var body struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotModel = body.Model
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(openaiReply(prov.Identity().Model, validReply(req, prov.Identity(), semantic.ResultPass)))
	}, "real-model", "secret-test-key")

	resp, err := semantic.NewRunner(p).Run(context.Background(), req, stubReader{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp == nil || resp.Report.Result != semantic.ResultPass {
		t.Fatalf("response = %+v", resp)
	}
	if gotAuth != "Bearer secret-test-key" {
		t.Errorf("auth header = %q", gotAuth)
	}
	if gotModel != "real-model" {
		t.Errorf("requested model = %q", gotModel)
	}
	recs := p.Records()
	if len(recs) != 1 || recs[0].Status != "ok" || recs[0].ReportedModel != "real-model" {
		t.Fatalf("records = %+v", recs)
	}
	if len(recs[0].Reply) == 0 || recs[0].ReplyHash == "" || recs[0].RequestHash == "" {
		t.Error("reply/request identity not retained")
	}
}

func TestRealAdapterCleanRequestCanary(t *testing.T) {
	const canary = "CANYARY-PRIVATE-GENERATOR-REASONING"
	reqGen := baseReq(semantic.RoleNormalizer, "gen-1", "es-1")
	reqGen.Evidence = []semantic.Evidence{{ArtifactID: "page-1", Page: nil, MediaType: "text"}}
	reqVal := baseReq(semantic.RoleClaimValidator, "val-1", "es-1")
	reqVal.Evidence = []semantic.Evidence{{ArtifactID: "page-2", Page: nil, MediaType: "text"}}

	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		var body struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		var reply []byte
		if len(body.Messages) > 0 && strings.Contains(body.Messages[0].Content, semantic.RoleNormalizer) {
			reply = validReply(reqGen, prov.Identity(), semantic.ResultPass)
		} else {
			reply = validReply(reqVal, prov.Identity(), semantic.ResultPass)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(openaiReply(prov.Identity().Model, reply))
	}, "real-model", "secret-test-key")

	reader := mapReader{"page-1": []byte(canary), "page-2": []byte("clean labelled evidence")}
	runner := semantic.NewRunner(p)
	if _, err := runner.Run(context.Background(), reqGen, reader); err != nil {
		t.Fatalf("generator: %v", err)
	}
	if _, err := runner.Run(context.Background(), reqVal, reader); err != nil {
		t.Fatalf("validator: %v", err)
	}

	var genReq, valReq []byte
	for _, rec := range p.Records() {
		switch rec.Role {
		case semantic.RoleNormalizer:
			genReq = rec.CapturedRequest
		case semantic.RoleClaimValidator:
			valReq = rec.CapturedRequest
		}
	}
	if !bytes.Contains(genReq, []byte(canary)) {
		t.Error("generator request is missing its own full allowlisted evidence (silent truncation)")
	}
	if bytes.Contains(valReq, []byte(canary)) {
		t.Fatal("generator private canary leaked into the validator's actual request")
	}
}

func TestEvidenceBudgetNotTruncated(t *testing.T) {
	big := bytes.Repeat([]byte("x"), (8<<20)+1)
	req := baseReq(semantic.RoleNormalizer, "job-1", "es-1")
	req.Evidence = []semantic.Evidence{{ArtifactID: "big", Page: nil, MediaType: "text"}}
	calls := 0
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		calls++
		w.Write(openaiReply(prov.Identity().Model, validReply(req, prov.Identity(), semantic.ResultPass)))
	}, "real-model", "secret-test-key")
	_, err := semantic.NewRunner(p).Run(context.Background(), req, mapReader{"big": big})
	if code(err) != semantic.CodeOverLimit {
		t.Fatalf("evidence over budget not rejected as over-limit: %v", err)
	}
	if calls != 0 {
		t.Fatalf("evidence over budget still sent %d HTTP request(s)", calls)
	}
}

func TestMissingEvidenceFails(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	req.Evidence = []semantic.Evidence{{ArtifactID: "absent", Page: nil, MediaType: "text"}}
	calls := 0
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) { calls++ }, "real-model", "k")
	_, err := semantic.NewRunner(p).Run(context.Background(), req, mapReader{})
	if code(err) != semantic.CodeMissingField {
		t.Fatalf("missing evidence not rejected: %v", err)
	}
	if calls != 0 {
		t.Fatalf("missing evidence still sent %d HTTP request(s)", calls)
	}
}

func TestRealAdapterStrictSchema(t *testing.T) {
	for name, raw := range map[string][]byte{
		"invalid json":   []byte(`{not json`),
		"unknown member": []byte(`{"version":"atlas.semantic-job/v1","job_id":"job-1","model_identity":{"provider":"h","model":"real-model","revision":"rev-test"},"report":{"version":"atlas.semantic-report/v1","role":"manifesto-claim-validator","subject":"prog-1","inputs":{"source":"sha256:a","draft":"sha256:b"},"result":"pass","checked":{"statements":[],"measures":[],"pages":[],"topics":[]},"details":{},"bogus":1}`),
		"invalid result": []byte(`{"version":"atlas.semantic-job/v1","job_id":"job-1","model_identity":{"provider":"h","model":"real-model","revision":"rev-test"},"report":{"version":"atlas.semantic-report/v1","role":"manifesto-claim-validator","subject":"prog-1","inputs":{"source":"sha256:a","draft":"sha256:b"},"result":"maybe","checked":{"statements":[],"measures":[],"pages":[],"topics":[]},"details":{}}`),
		"contradictory":  []byte(`{"version":"atlas.semantic-job/v1","job_id":"job-1","model_identity":{"provider":"h","model":"real-model","revision":"rev-test"},"report":{"version":"atlas.semantic-report/v1","role":"manifesto-claim-validator","subject":"prog-1","inputs":{"source":"sha256:a","draft":"sha256:b"},"result":"pass","findings":[{"id":"f-1","code":"c","severity":"blocker","target":"t","evidence":"e","message":"m"}],"checked":{"statements":[],"measures":[],"pages":[],"topics":[]},"details":{}}`),
	} {
		t.Run(name, func(t *testing.T) {
			req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
			p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
				w.Write(openaiReply("real-model", raw))
			}, "real-model", "k")
			resp, err := semantic.NewRunner(p).Run(context.Background(), req, stubReader{})
			if resp != nil || code(err) != semantic.CodeBadReport {
				t.Fatalf("malformed reply: resp=%v err=%v", resp, err)
			}
		})
	}
}

func TestRealAdapterOverBudgetNoApproval(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	req.Limits.MaxResponseBytes = 64
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		w.Write(openaiReply("real-model", validReply(req, prov.Identity(), semantic.ResultPass)))
	}, "real-model", "k")
	resp, err := semantic.NewRunner(p).Run(context.Background(), req, stubReader{})
	if resp != nil || err == nil {
		t.Fatalf("over-budget response must not return approval: resp=%v err=%v", resp, err)
	}
}

func TestSecretRedaction(t *testing.T) {
	const secret = "super-secret-provider-key"
	req := baseReq(semantic.RoleNormalizer, "job-1", "es-1")
	req.Evidence = []semantic.Evidence{{ArtifactID: "p", Page: nil, MediaType: "text"}}
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		w.Write(openaiReply(prov.Identity().Model, validReply(req, prov.Identity(), semantic.ResultPass)))
	}, "real-model", secret)
	if _, err := semantic.NewRunner(p).Run(context.Background(), req, mapReader{"p": []byte("evidence")}); err != nil {
		t.Fatalf("run: %v", err)
	}
	for _, rec := range p.Records() {
		if strings.Contains(string(rec.CapturedRequest), secret) {
			t.Fatal("API key leaked into the captured request body")
		}
	}
}

// cancelTransport simulates a slow model: it blocks until the request context is
// done, signalling that the transport observed the cancellation and did not read a
// late reply.
type cancelTransport struct {
	started  chan struct{}
	observed chan struct{}
}

func (c *cancelTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	select {
	case <-c.started:
	default:
		close(c.started)
	}
	<-req.Context().Done()
	select {
	case <-c.observed:
	default:
		close(c.observed)
	}
	return nil, req.Context().Err()
}

func TestCancellationClosesTransport(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	rt := &cancelTransport{started: make(chan struct{}), observed: make(chan struct{})}
	p, err := semantic.NewOpenAIProvider(semantic.ProviderConfig{
		BaseURL:    "http://127.0.0.1:9/v1",
		Model:      "real-model",
		HTTPClient: &http.Client{Transport: rt},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := semantic.NewRunner(p).Run(ctx, req, stubReader{})
		done <- err
	}()
	<-rt.started
	cancel()
	select {
	case err := <-done:
		if code(err) != semantic.CodeCancellation {
			t.Fatalf("cancellation not enforced: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runner did not return after cancellation")
	}
	select {
	case <-rt.observed:
	case <-time.After(5 * time.Second):
		t.Fatal("transport did not observe the cancellation (a late reply would be read)")
	}
}

type failTransport struct{ calls int }

func (f *failTransport) RoundTrip(*http.Request) (*http.Response, error) {
	f.calls++
	return nil, errors.New("connection reset by peer")
}

func TestRealAdapterTransientRetry(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	ft := &failTransport{}
	p, err := semantic.NewOpenAIProvider(semantic.ProviderConfig{
		BaseURL:    "http://127.0.0.1:9/v1",
		Model:      "real-model",
		HTTPClient: &http.Client{Transport: ft},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = semantic.NewRunner(p).Run(context.Background(), req, stubReader{})
	if code(err) != semantic.CodeProviderError {
		t.Fatalf("transient-then-failure not surfaced as provider error: %v", err)
	}
	if ft.calls != 2 {
		t.Fatalf("calls = %d, want 2 (exactly one clean retry)", ft.calls)
	}
}

func TestRealAdapterNon2xxNotRetried(t *testing.T) {
	req := baseReq(semantic.RoleClaimValidator, "job-1", "es-1")
	calls := 0
	p := newOpenAIProvider(t, func(w http.ResponseWriter, r *http.Request, prov *semantic.OpenAIProvider) {
		calls++
		http.Error(w, "server error", http.StatusInternalServerError)
	}, "real-model", "k")
	_, err := semantic.NewRunner(p).Run(context.Background(), req, stubReader{})
	if code(err) != semantic.CodeProviderError {
		t.Fatalf("non-2xx not a provider error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("non-2xx retried %d times, want 0 (no retry on a received response)", calls)
	}
}
