// Command smoke-provider exercises the actual OpenAI-compatible configured
// provider adapter against a local model and retains genuine request/model
// identities and the original bounded replies. It is an explicit operator gate
// (make smoke-provider), never part of make verify, and uses no committed
// credentials: the endpoint, model and key are read from the operator's
// environment at run time.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/semantic"
)

type smokeConfig struct {
	Version         string   `json:"version"`
	Endpoint        endpoint `json:"endpoint"`
	MaxOutputTokens int      `json:"max_output_tokens"`
	Targets         []target `json:"targets"`
}

type endpoint struct {
	BaseURLEnv string `json:"base_url_env"`
	ModelEnv   string `json:"model_env"`
	APIKeyEnv  string `json:"api_key_env"`
	Revision   string `json:"revision"`
}

type target struct {
	Role     string            `json:"role"`
	JobID    string            `json:"job_id"`
	Subject  string            `json:"subject"`
	Inputs   map[string]string `json:"inputs"`
	Evidence []evidence        `json:"evidence"`
}

type evidence struct {
	ArtifactID string `json:"artifact_id"`
	Content    string `json:"content"`
}

type smokeReader struct{ data map[string][]byte }

func (s smokeReader) Read(id string, page *int) ([]byte, error) {
	b, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("missing smoke evidence: %s", id)
	}
	return b, nil
}

// retainedInvocation is the credential-free retained record of one actual run.
type retainedInvocation struct {
	Role           string `json:"role"`
	JobID          string `json:"job_id"`
	Status         string `json:"status"`
	Endpoint       string `json:"endpoint"`
	RequestedModel string `json:"requested_model"`
	ReportedModel  string `json:"reported_model"`
	ModelIdentity  string `json:"model_identity"`
	PromptHash     string `json:"prompt_hash"`
	EvidenceHash   string `json:"evidence_hash"`
	InputsHash     string `json:"inputs_hash"`
	RequestHash    string `json:"request_hash"`
	ReplyBytes     int    `json:"reply_bytes"`
	ReplyHash      string `json:"reply_hash"`
	Reply          string `json:"reply"`
	Verified       bool   `json:"verified"`
	Result         string `json:"result,omitempty"`
	Error          string `json:"error,omitempty"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "smoke-provider:", err)
		os.Exit(1)
	}
}

func run() error {
	targetPath := flag.String("target", "fixtures/semantic/smoke.json", "path to the smoke target fixture")
	outPath := flag.String("out", "fixtures/semantic/smoke-actual.json", "path to write the retained actual evidence")
	flag.Parse()

	cfg, err := loadConfig(*targetPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	baseURL := os.Getenv(cfg.Endpoint.BaseURLEnv)
	model := os.Getenv(cfg.Endpoint.ModelEnv)
	apiKey := os.Getenv(cfg.Endpoint.APIKeyEnv)
	if baseURL == "" || model == "" {
		return fmt.Errorf("provider endpoint not configured: set %s and %s (API key in %s)",
			cfg.Endpoint.BaseURLEnv, cfg.Endpoint.ModelEnv, cfg.Endpoint.APIKeyEnv)
	}

	provider, err := semantic.NewOpenAIProvider(semantic.ProviderConfig{
		BaseURL:  baseURL,
		Model:    model,
		APIKey:   apiKey,
		Revision: cfg.Endpoint.Revision,
	})
	if err != nil {
		return fmt.Errorf("build provider: %w", err)
	}

	reader := smokeReader{data: map[string][]byte{}}
	for _, t := range cfg.Targets {
		for _, e := range t.Evidence {
			reader.data[e.ArtifactID] = []byte(e.Content)
		}
	}

	runner := semantic.NewRunner(provider)
	verified := map[string]struct {
		ok     bool
		result string
		err    string
	}{}

	for _, t := range cfg.Targets {
		req := semantic.Request{
			Version:       semantic.JobVersion,
			JobID:         t.JobID,
			EvaluationSet: "smoke",
			Role:          t.Role,
			Subject:       t.Subject,
			Inputs:        t.Inputs,
			Model:         model,
			Evidence:      buildEvidence(t),
			Limits: semantic.Limits{
				Deadline:        5 * time.Minute,
				MaxOutputTokens: cfg.MaxOutputTokens,
			},
		}
		resp, runErr := runner.Run(context.Background(), req, reader)
		v := struct {
			ok     bool
			result string
			err    string
		}{}
		if runErr == nil && resp != nil {
			v.ok = true
			v.result = resp.Report.Result
		} else if runErr != nil {
			v.err = runErr.Error()
		}
		verified[t.JobID] = v
	}

	recs := provider.Records()
	if len(recs) != len(cfg.Targets) {
		return fmt.Errorf("expected %d invocations, got %d", len(cfg.Targets), len(recs))
	}
	retained := make([]retainedInvocation, 0, len(recs))
	allReached := true
	for _, r := range recs {
		v := verified[r.JobID]
		if r.Status != "ok" {
			allReached = false
		}
		ri := retainedInvocation{
			Role:           r.Role,
			JobID:          r.JobID,
			Status:         r.Status,
			Endpoint:       r.Endpoint,
			RequestedModel: r.RequestedModel,
			ReportedModel:  r.ReportedModel,
			ModelIdentity:  identityString(r.ModelIdentity),
			PromptHash:     r.PromptHash,
			EvidenceHash:   r.EvidenceHash,
			InputsHash:     r.InputsHash,
			RequestHash:    r.RequestHash,
			ReplyBytes:     r.ReplyBytes,
			ReplyHash:      r.ReplyHash,
			Reply:          string(r.Reply),
			Verified:       v.ok,
			Result:         v.result,
			Error:          v.err,
		}
		retained = append(retained, ri)
	}

	out := struct {
		Version     string               `json:"version"`
		Endpoint    string               `json:"endpoint"`
		Model       string               `json:"model"`
		Revision    string               `json:"revision"`
		Identity    string               `json:"identity"`
		Invocations []retainedInvocation `json:"invocations"`
	}{
		Version:     "atlas.semantic-smoke-actual/v1",
		Endpoint:    baseURL,
		Model:       model,
		Revision:    cfg.Endpoint.Revision,
		Identity:    identityString(provider.Identity()),
		Invocations: retained,
	}
	raw, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal retained evidence: %w", err)
	}
	if err := os.WriteFile(*outPath, append(raw, '\n'), 0o644); err != nil {
		return fmt.Errorf("write retained evidence: %w", err)
	}
	fmt.Printf("smoke-provider: retained %d invocation(s) to %s\n", len(retained), *outPath)

	if !allReached {
		return fmt.Errorf("provider unreachable or a role produced no reply (see %s); the task must stay open", *outPath)
	}
	return nil
}

func buildEvidence(t target) []semantic.Evidence {
	out := make([]semantic.Evidence, 0, len(t.Evidence))
	for _, e := range t.Evidence {
		out = append(out, semantic.Evidence{ArtifactID: e.ArtifactID, Page: nil, MediaType: "text"})
	}
	return out
}

func identityString(id semantic.ModelIdentity) string {
	return id.Provider + "/" + id.Model + "@" + id.Revision
}

func loadConfig(path string) (smokeConfig, error) {
	var cfg smokeConfig
	raw, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
