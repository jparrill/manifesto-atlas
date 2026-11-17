// Package summarize writes the factual summary of a programme with a configured
// language model, and repairs it until the deterministic gate accepts it.
//
// The model is untrusted: everything it writes is checked afterwards by
// internal/runcheck. The pipeline never trusts the model, never lets it choose
// the endpoint or the model, and never sends credentials inside a body.
package summarize

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jparrill/manifesto-atlas/internal/runcheck"
	"github.com/jparrill/manifesto-atlas/internal/runmanifest"
)

// Config is the operator-fixed transport and budget of one summarisation.
type Config struct {
	BaseURL   string // OpenAI-compatible root, e.g. http://127.0.0.1:8787/v1
	Model     string
	APIKey    string // header only; never logged nor put in a body
	Attempts  int    // total attempts, at least 1
	MaxInput  int    // bytes of extraction sent to the model
	HTTP      *http.Client
	TracePath string // optional credential-free record of the run
}

// DefaultMaxInput bounds the extraction sent in one request.
const DefaultMaxInput = 2_000_000

// Result is the outcome of summarising one programme.
type Result struct {
	Attempts int
	Issues   []runcheck.Issue
	Markdown []byte
	Model    string
	Reply    string
	Resolved bool
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type chatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Rules is the contract the model must follow. It is identical for every run so
// the behaviour does not drift with the prompt of the day.
const Rules = `Eres un redactor documental. Conviertes un programa electoral publicado en un
resumen factual en español. Reglas absolutas:

1. Escribe SOLO lo que el documento dice. No añadas datos externos, opiniones,
   valoraciones de viabilidad ni comparaciones con otros partidos.
2. Una línea por propuesta, terminada en la página del PDF entre corchetes:
   [p. 12] o [pp. 41-43]. La página debe ser la del PDF, no la impresa.
3. Las comillas angulares «» solo para texto LITERAL del documento, copiado
   carácter a carácter. Prohibido meter paráfrasis, negaciones reescritas o
   títulos inventados dentro de «».
4. Respeta el número y el orden de las propuestas del original. Si el documento
   numera 1..N, tu resumen numera 1..N. Si reinicia la numeración en cada
   apartado, reiníciala igual.
5. Usa la estructura del documento: sus apartados, con su título y su página.
6. Si una propuesta aparece en el documento sin contenido («permanezca atento a
   la campaña»), decláralo así; no la rellenes.
7. Nada de tablas de datos externos. Empieza con el título del documento y una
   tabla de metadatos: documento, naturaleza, ámbito, fecha, páginas y número de
   propuestas.
8. Devuelve Markdown y nada más: sin bloques de código, sin comentarios, sin
   explicaciones antes o después.`

var fence = regexp.MustCompile("(?s)^```[a-zA-Z]*\\n(.*)\\n```\\s*$")

// Generate asks the model for the summary of one extraction.
func Generate(ctx context.Context, cfg Config, title, extraction string) (string, error) {
	if cfg.BaseURL == "" || cfg.Model == "" {
		return "", fmt.Errorf("a summary needs an operator-fixed endpoint and model")
	}
	if len(extraction) > cfg.MaxInput {
		return "", fmt.Errorf("extraction has %d bytes, over the %d limit; raise MaxInput explicitly",
			len(extraction), cfg.MaxInput)
	}
	reply, err := call(ctx, cfg, []message{
		{Role: "system", Content: Rules},
		{Role: "user", Content: "Programa a resumir: " + title + "\n\nTexto extraído página a página,\ncada página precedida por un marcador <!-- page N -->:\n\n" + extraction},
	})
	if err != nil {
		return "", err
	}
	return clean(reply), nil
}

// Repair asks the model to fix the issues the gate reported, without touching
// anything the gate accepted.
func Repair(ctx context.Context, cfg Config, title, current string, issues []runcheck.Issue) (string, error) {
	var b strings.Builder
	b.WriteString("Tu resumen anterior ha fallado la comprobación automática. Corrige SOLO lo que indican\n")
	b.WriteString("estos fallos y devuelve el documento completo otra vez, con las mismas reglas:\n\n")
	for _, i := range issues {
		fmt.Fprintf(&b, "- %s en %s: %s\n", i.Code, i.File, i.Detail)
	}
	b.WriteString("\nResumen actual:\n\n")
	b.WriteString(current)
	reply, err := call(ctx, cfg, []message{
		{Role: "system", Content: Rules},
		{Role: "user", Content: b.String()},
	})
	if err != nil {
		return "", err
	}
	return clean(reply), nil
}

func call(ctx context.Context, cfg Config, msgs []message) (string, error) {
	body, err := json.Marshal(chatRequest{Model: cfg.Model, Messages: msgs, Temperature: 0})
	if err != nil {
		return "", err
	}
	url := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}
	client := cfg.HTTP
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("provider request: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("provider returned %d: %s", resp.StatusCode, firstLine(raw))
	}
	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("provider reply is not JSON: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("provider reply has no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}

func firstLine(raw []byte) string {
	s := strings.TrimSpace(string(raw))
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 300 {
		s = s[:300]
	}
	return s
}

// clean removes code fences and surrounding blank lines a model may add.
func clean(reply string) string {
	reply = strings.TrimSpace(reply)
	if m := fence.FindStringSubmatch(reply); m != nil {
		reply = m[1]
	}
	return strings.TrimSpace(reply) + "\n"
}

// Run summarises one programme and repairs it until the gate accepts it or the
// attempts run out. The gate is the only authority: a reply that fails it is
// never written as final.
func Run(ctx context.Context, cfg Config, runDir, programme string) (Result, error) {
	spec, err := runmanifest.Load(runDir)
	if err != nil {
		return Result{}, err
	}
	var doc runmanifest.Document
	found := false
	for _, d := range spec.Documents {
		if d.Programme == programme {
			doc, found = d, true
		}
	}
	if !found {
		return Result{}, fmt.Errorf("programme %q is not in %s", programme, filepath.Join(runDir, "run.json"))
	}
	extraction, err := os.ReadFile(filepath.Join(runDir, doc.Source))
	if err != nil {
		return Result{}, fmt.Errorf("read extraction: %w", err)
	}
	if cfg.Attempts <= 0 {
		cfg.Attempts = 1
	}
	summaryPath := filepath.Join(runDir, doc.Summary)
	title := doc.Title
	if title == "" {
		title = doc.Programme
	}

	var res Result
	res.Model = cfg.Model
	for attempt := 1; attempt <= cfg.Attempts; attempt++ {
		res.Attempts = attempt
		var markdown string
		var err error
		if attempt == 1 {
			markdown, err = Generate(ctx, cfg, title, string(extraction))
		} else {
			markdown, err = Repair(ctx, cfg, title, string(res.Markdown), res.Issues)
		}
		if err != nil {
			return res, err
		}
		res.Markdown = []byte(markdown)
		if err := os.WriteFile(summaryPath, res.Markdown, 0o644); err != nil {
			return res, err
		}
		report, err := runcheck.Verify(runDir)
		if err != nil {
			return res, err
		}
		res.Issues = report.Issues
		if report.OK() {
			res.Resolved = true
			break
		}
	}
	if err := writeTrace(cfg, res, summaryPath); err != nil {
		return res, err
	}
	return res, nil
}

type trace struct {
	Model    string   `json:"model"`
	Attempts int      `json:"attempts"`
	Resolved bool     `json:"resolved"`
	Summary  string   `json:"summary_sha256"`
	Issues   []string `json:"issues,omitempty"`
}

// writeTrace records what happened without the prompt, the reply or the key.
func writeTrace(cfg Config, res Result, summaryPath string) error {
	if cfg.TracePath == "" {
		return nil
	}
	sum := sha256.Sum256(res.Markdown)
	t := trace{Model: res.Model, Attempts: res.Attempts, Resolved: res.Resolved,
		Summary: hex.EncodeToString(sum[:])}
	for _, i := range res.Issues {
		t.Issues = append(t.Issues, i.String())
	}
	raw, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cfg.TracePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(cfg.TracePath, append(raw, '\n'), 0o644)
}
