package summarize

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const extraction = `<!-- page 1 -->

VIVIENDA
1. Construir vivienda pública.
2. Prohibir los pisos turísticos.
`

const goodSummary = "# Resumen\n\n## 1. Vivienda — p. 1\n\n" +
	"1. «Construir vivienda pública» [p. 1]\n" +
	"2. Prohibir los pisos turísticos [p. 1]\n"

const badSummary = "# Resumen\n\n## 1. Vivienda — p. 1\n\n" +
	"1. «Levantar casas públicas» [p. 1]\n" +
	"2. Prohibir los pisos turísticos [p. 1]\n"

// run prepares a run directory with one programme and its extraction.
func run(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, "pp.md", extraction)
	write(t, dir, "run.json", `{"name":"prueba","claims":"checks/claims.tsv",
  "sources":[{"programme":"pp","url":"https://example.org/p.pdf","file":"p.pdf","bytes":10,
    "sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","pages":1}],
  "documents":[{"programme":"pp","source":"pp.md","summary":"Resumen-PP.md","proposals":2,"title":"PP"}]}`)
	write(t, dir, "checks/claims.tsv", "# ninguna\n")
	return dir
}

func write(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

// provider returns a server answering with the given replies in order.
func provider(t *testing.T, replies ...string) (*httptest.Server, *[]chatRequest) {
	t.Helper()
	var seen []chatRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" && got != "Bearer secreto" {
			t.Errorf("unexpected authorization header %q", got)
		}
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		seen = append(seen, req)
		reply := replies[min(len(seen)-1, len(replies)-1)]
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model":   req.Model,
			"choices": []map[string]any{{"message": map[string]string{"content": reply}}},
		})
	}))
	t.Cleanup(srv.Close)
	return srv, &seen
}

func TestGenerate(t *testing.T) {
	t.Run("When the model answers with fences they should be stripped", func(t *testing.T) {
		srv, _ := provider(t, "```markdown\n"+goodSummary+"```")
		got, err := Generate(context.Background(), Config{BaseURL: srv.URL, Model: "m", MaxInput: 1 << 20}, "PP", extraction)
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}
		if strings.Contains(got, "```") || !strings.HasPrefix(got, "# Resumen") {
			t.Fatalf("fences not stripped: %q", got)
		}
	})
	t.Run("When the extraction is over the budget it should refuse", func(t *testing.T) {
		srv, _ := provider(t, goodSummary)
		if _, err := Generate(context.Background(), Config{BaseURL: srv.URL, Model: "m", MaxInput: 10}, "PP", extraction); err == nil {
			t.Fatalf("expected a budget error")
		}
	})
	t.Run("When the endpoint or model is missing it should refuse", func(t *testing.T) {
		if _, err := Generate(context.Background(), Config{}, "PP", extraction); err == nil {
			t.Fatalf("expected a configuration error")
		}
	})
	t.Run("When the provider fails it should return the status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		}))
		defer srv.Close()
		if _, err := Generate(context.Background(), Config{BaseURL: srv.URL, Model: "m", MaxInput: 1 << 20}, "PP", extraction); err == nil {
			t.Fatalf("expected a provider error")
		}
	})
}

func TestRunRepairsUntilTheGateAccepts(t *testing.T) {
	dir := run(t)
	srv, seen := provider(t, badSummary, goodSummary)
	res, err := Run(context.Background(), Config{BaseURL: srv.URL, Model: "m", APIKey: "secreto",
		Attempts: 3, MaxInput: 1 << 20, TracePath: filepath.Join(dir, "trace.json")}, dir, "pp")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.Resolved || res.Attempts != 2 {
		t.Fatalf("resolved=%v attempts=%d, want true/2", res.Resolved, res.Attempts)
	}
	if len(*seen) != 2 {
		t.Fatalf("provider calls = %d, want 2", len(*seen))
	}
	if !strings.Contains((*seen)[1].Messages[1].Content, "runcheck.quote") {
		t.Fatalf("the repair request does not carry the gate issue: %s", (*seen)[1].Messages[1].Content)
	}
	final, err := os.ReadFile(filepath.Join(dir, "Resumen-PP.md"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if string(final) != goodSummary {
		t.Fatalf("summary not replaced by the repaired version:\n%s", final)
	}
}

func TestRunGivesUpAfterTheBudget(t *testing.T) {
	dir := run(t)
	srv, _ := provider(t, badSummary)
	res, err := Run(context.Background(), Config{BaseURL: srv.URL, Model: "m", Attempts: 2, MaxInput: 1 << 20}, dir, "pp")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Resolved {
		t.Fatalf("a failing summary was accepted")
	}
	if res.Attempts != 2 || len(res.Issues) == 0 {
		t.Fatalf("attempts=%d issues=%d", res.Attempts, len(res.Issues))
	}
}

func TestTraceHasNoSecrets(t *testing.T) {
	dir := run(t)
	srv, _ := provider(t, goodSummary)
	trace := filepath.Join(dir, "trace.json")
	if _, err := Run(context.Background(), Config{BaseURL: srv.URL, Model: "m", APIKey: "secreto",
		Attempts: 1, MaxInput: 1 << 20, TracePath: trace}, dir, "pp"); err != nil {
		t.Fatalf("Run: %v", err)
	}
	raw, err := os.ReadFile(trace)
	if err != nil {
		t.Fatalf("read trace: %v", err)
	}
	for _, secret := range []string{"secreto", srv.URL, "Construir vivienda"} {
		if strings.Contains(string(raw), secret) {
			t.Fatalf("trace leaks %q:\n%s", secret, raw)
		}
	}
}
