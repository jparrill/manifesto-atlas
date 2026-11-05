package preview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/semantic"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

type Config struct {
	Source   string
	Output   string
	Endpoint string
	Model    string
	APIKey   string
}

type Outcome struct {
	Result   string             `json:"result"`
	Error    string             `json:"error,omitempty"`
	Response *semantic.Response `json:"response,omitempty"`
}

type Record struct {
	Version         string             `json:"version"`
	PrivatePreview  bool               `json:"private_preview"`
	PublicEligible  bool               `json:"public_eligible"`
	Inputs          map[string]string  `json:"inputs"`
	Reports         map[string]Outcome `json:"reports"`
	ReferenceSHA256 string             `json:"reference_sha256"`
}

func requiredRunInputs() []string {
	names := []string{"source.pdf", "metadata.json", "inspection.txt", "extraction.txt", "extraction.json", "evidence.json", "programme.json", "taxonomy.json", "model.json", "normalizer.json"}
	for _, role := range []string{semantic.RoleNormalizer, semantic.RoleClaimValidator, semantic.RoleCoverageValidator} {
		for _, suffix := range []string{"-request.json", "-reply.txt", "-invocation.json"} {
			names = append(names, role+suffix)
		}
	}
	return append(names, semantic.RoleClaimValidator+".json", semantic.RoleCoverageValidator+".json")
}

type evidence map[string][]byte

func (e evidence) Read(id string, page *int) ([]byte, error) {
	value, ok := e[id]
	if !ok || page != nil {
		return nil, errors.New("evidence not allowlisted")
	}
	return value, nil
}

const normalizationPrompt = `Use Spanish plain text. This is a concise PRIVATE preview of the principal proposals, not a claim of exhaustive coverage.
In report.details return exactly {"topics":[{"id":"TOPIC_ID","state":"proposals|mentioned|not-covered|unclear","summary":{"text":"short Spanish factual summary","evidence_id":"p5-e1"},"measures":[{"text":"one concrete source-supported action, preserving relevant conditions","evidence_id":"p5-e1"}]}]}.
Return ALL 21 supplied topic IDs in the exact supplied order. For proposals give one cited summary and exactly ONE principal cited measure. For mentioned give a cited summary and no measures. For not-covered or unclear give summary:null and measures:[]. Never use not-applicable or infer opposition from omission. Prefer unclear to a guessed classification. Do not invent facts, institutions, quantities or legal powers.
Each evidence_id MUST be selected from the supplied original-evidence catalogue. These IDs identify exact original page slices; never invent IDs, quote text, page fields or extra citation fields. Select a slice that actually supports that entire statement. The host resolves the literal quotation and original PDF page. Keep every summary/action under 250 characters and the complete reply concise. Return a valid details object even if normalization is insufficient. Do not repeat the source. Nothing outside the common JSON envelope.`

func sourceRecord(metadata []byte, checksum string) (canonical.SourceRecord, error) {
	var m struct {
		ProgrammeID     string              `json:"programme_id"`
		Party           canonical.Party     `json:"party"`
		DocumentType    string              `json:"document_type"`
		Scope           string              `json:"scope"`
		Geography       canonical.Geography `json:"geography"`
		Election        *canonical.Election `json:"election"`
		PublicationDate *string             `json:"publication_date"`
		Rights          canonical.Rights    `json:"rights"`
	}
	if err := json.Unmarshal(metadata, &m); err != nil {
		return canonical.SourceRecord{}, err
	}
	return canonical.SourceRecord{ProgrammeID: m.ProgrammeID, SourceChecksum: checksum, SourceMetadataChecksum: hashBytes(metadata),
		Party: m.Party, DocumentType: m.DocumentType, Scope: m.Scope, Geography: m.Geography, Election: m.Election,
		PublicationDate: m.PublicationDate, Rights: m.Rights}, nil
}

func checkedInventory(p *canonical.Programme, count int) semantic.Checked {
	result := semantic.Checked{Statements: []string{}, Measures: []string{}, Pages: []int{}, Topics: []string{}}
	for _, s := range p.Summary {
		result.Statements = append(result.Statements, s.ID)
	}
	if p.Metadata.DocumentTitle != nil {
		result.Statements = append(result.Statements, p.Metadata.DocumentTitle.ID)
	}
	for _, t := range p.Topics {
		result.Topics = append(result.Topics, t.ID)
		for _, s := range t.Summary {
			result.Statements = append(result.Statements, s.ID)
		}
	}
	for _, m := range p.Measures {
		result.Measures = append(result.Measures, m.ID)
		result.Statements = append(result.Statements, m.Action.ID)
		for _, c := range m.Conditions {
			result.Statements = append(result.Statements, c.ID)
		}
		if m.Institution != nil {
			result.Statements = append(result.Statements, m.Institution.ID)
		}
	}
	for i := 1; i <= count; i++ {
		result.Pages = append(result.Pages, i)
	}
	return result
}

func completeInventory(got, want semantic.Checked) bool {
	stringsEqual := func(a, b []string) bool {
		a, b = append([]string{}, a...), append([]string{}, b...)
		slices.Sort(a)
		slices.Sort(b)
		return reflect.DeepEqual(a, b)
	}
	a, b := append([]int{}, got.Pages...), append([]int{}, want.Pages...)
	slices.Sort(a)
	slices.Sort(b)
	return reflect.DeepEqual(a, b) && stringsEqual(got.Statements, want.Statements) &&
		stringsEqual(got.Measures, want.Measures) && stringsEqual(got.Topics, want.Topics)
}

func outputPath(path string) (string, error) {
	clean := filepath.Clean(path)
	if filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("preview output must be a named directory within this checkout")
	}
	current := "."
	for _, part := range strings.Split(clean, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		fi, err := os.Lstat(current)
		if err == nil && (!fi.IsDir() || fi.Mode()&os.ModeSymlink != 0) {
			return "", errors.New("unsafe preview output path")
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	return clean, nil
}

func writeInvocation(dir, prefix string, record semantic.RecordedInvocation) error {
	if err := writeOutput(filepath.Join(dir, prefix+"-request.json"), record.CapturedRequest); err != nil {
		return err
	}
	if err := writeOutput(filepath.Join(dir, prefix+"-reply.txt"), record.Reply); err != nil {
		return err
	}
	return saveJSON(filepath.Join(dir, prefix+"-invocation.json"), map[string]any{
		"role": record.Role, "job_id": record.JobID, "model_identity": record.ModelIdentity,
		"reported_model_basename": filepath.Base(record.ReportedModel), "status": record.Status,
		"request_hash": record.RequestHash, "prompt_hash": record.PromptHash, "inputs_hash": record.InputsHash,
		"evidence_hash": record.EvidenceHash, "reply_hash": record.ReplyHash, "reply_bytes": record.ReplyBytes,
		"request_bytes": record.RequestBytes,
	})
}

func persistInvocation(dir, role string, provider *semantic.OpenAIProvider) ([]string, error) {
	files := []string{}
	var records []semantic.RecordedInvocation
	for _, record := range provider.Records() {
		if record.Role == role {
			records = append(records, record)
		}
	}
	if len(records) == 0 {
		return nil, errors.New("provider invocation evidence missing")
	}
	for i, record := range records {
		prefix := fmt.Sprintf("%s-attempt-%d", role, i+1)
		if err := writeInvocation(dir, prefix, record); err != nil {
			return nil, err
		}
		for _, suffix := range []string{"-request.json", "-reply.txt", "-invocation.json"} {
			files = append(files, prefix+suffix)
		}
	}
	// Preserve the simple last-response paths used by existing retained runs.
	if err := writeInvocation(dir, role, records[len(records)-1]); err != nil {
		return nil, err
	}
	return files, nil
}

// Run executes a real, explicit model-backed preview. It never puts output in a
// public catalogue; failed semantic checks are retained and prominently shown.
func Run(ctx context.Context, cfg Config) error {
	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return errors.New("preview endpoint must be configured without embedded credentials, query or fragment")
	}
	dir, err := outputPath(cfg.Output)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(dir); !errors.Is(err, os.ErrNotExist) {
		return errors.New("output already exists: choose a new run directory or use -render")
	}
	reference, err := regularRead(referencePath, 256<<10)
	if err != nil {
		return fmt.Errorf("import the supplied reference first: %w", err)
	}
	c, metadata, pdf, err := loadCandidate(cfg.Source)
	if err != nil {
		return err
	}
	provider, err := semantic.NewOpenAIProvider(semantic.ProviderConfig{BaseURL: cfg.Endpoint, Model: cfg.Model, APIKey: cfg.APIKey, Revision: "unknown"})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, raw := range map[string][]byte{"source.pdf": pdf, "metadata.json": metadata} {
		if err := writeOutput(filepath.Join(dir, name), raw); err != nil {
			return err
		}
	}
	if err := saveJSON(filepath.Join(dir, "model.json"), map[string]string{"endpoint": cfg.Endpoint, "model": cfg.Model, "revision": "unknown"}); err != nil {
		return err
	}
	originalPath := filepath.Join(filepath.Dir(cfg.Source), c.CandidatePath)
	pdfPath, err := filepath.Abs(filepath.Join(dir, "source.pdf"))
	if err != nil {
		return err
	}
	fmt.Println("1/5 macOS inspection and original-page extraction")
	ex, err := extract(ctx, pdfPath, dir, hashBytes(pdf))
	if err != nil {
		return err
	}
	pages, _ := json.Marshal(ex)
	catalogue, err := json.Marshal(map[string]any{"page_count": ex.PageCount, "warnings": ex.Warnings, "original_evidence": evidenceSpans(ex)})
	if err != nil {
		return err
	}
	if err := writeOutput(filepath.Join(dir, "evidence.json"), catalogue); err != nil {
		return err
	}
	registry, _ := json.Marshal(taxonomy.Topics())
	if err := writeOutput(filepath.Join(dir, "taxonomy.json"), registry); err != nil {
		return err
	}
	exRaw, err := regularRead(filepath.Join(dir, "extraction.json"), 8<<20)
	if err != nil {
		return err
	}
	src, err := sourceRecord(metadata, hashBytes(pdf))
	if err != nil {
		return err
	}
	runner := semantic.NewRunner(provider)
	inputs := map[string]string{"source": hashBytes(pdf), "metadata": hashBytes(metadata), "extraction": hashBytes(exRaw), "taxonomy": hashBytes(registry)}
	invocationFiles := []string{}
	invoke := func(role, prompt string, ev evidence) (*semantic.Response, error) {
		maxTokens := 4096
		if role == semantic.RoleNormalizer {
			maxTokens = 8192
		}
		refs := []semantic.Evidence{}
		keys := make([]string, 0, len(ev))
		for key := range ev {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		for _, key := range keys {
			refs = append(refs, semantic.Evidence{ArtifactID: key, MediaType: "application/json"})
		}
		response, callErr := runner.Run(ctx, semantic.Request{Version: semantic.JobVersion, JobID: "preview-" + role,
			EvaluationSet: c.ProgrammeID, Role: role, Subject: c.ProgrammeID, Inputs: inputs, Prompt: prompt,
			OutputSchema: "atlas.preview/v1", Model: cfg.Model, Evidence: refs,
			Configuration: map[string]string{"enable_thinking": "false", "json_object": "true"},
			Limits:        semantic.Limits{Deadline: 600 * time.Second, MaxOutputTokens: maxTokens, MaxResponseBytes: 1 << 20}}, ev)
		files, err := persistInvocation(dir, role, provider)
		if err != nil {
			return nil, err
		}
		invocationFiles = append(invocationFiles, files...)
		return response, callErr
	}
	fmt.Printf("2/5 normalizer: %d original pages through %s\n", ex.PageCount, cfg.Model)
	normalized, err := invoke(semantic.RoleNormalizer, normalizationPrompt, evidence{"original-pages": catalogue, "ordered-topics": registry, "source-metadata": metadata})
	if err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(dir, "normalizer.json"), normalized); err != nil {
		return err
	}
	p, err := assemble(normalized.Report.Details, src, ex, hashBytes(exRaw))
	if err != nil {
		return fmt.Errorf("canonical source/citation check: %w", err)
	}
	if err := saveJSON(filepath.Join(dir, "programme.json"), p); err != nil {
		return err
	}
	draft, err := regularRead(filepath.Join(dir, "programme.json"), 16<<20)
	if err != nil {
		return err
	}
	inputs["draft"] = hashBytes(draft)
	inventory := checkedInventory(p, ex.PageCount)
	inventoryJSON, _ := json.Marshal(inventory)
	record := Record{Version: "atlas.preview.run/v1", PrivatePreview: true, PublicEligible: false,
		Inputs: map[string]string{}, Reports: map[string]Outcome{}, ReferenceSHA256: hashBytes(reference)}
	for i, role := range []string{semantic.RoleClaimValidator, semantic.RoleCoverageValidator} {
		fmt.Printf("%d/5 independent %s\n", i+3, role)
		instruction := "Check every public statement, exact quantity, polarity, relevant conditions, institution and original page/citation."
		if role == semantic.RoleCoverageValidator {
			instruction = "Read ALL original pages, not only citations. Find material omissions, duplicates, contradictions and wrong topic states. A concise preview can fail coverage: never certify completeness merely because a summary is readable."
		}
		prompt := instruction + " Your own clean request contains source evidence and the canonical draft only. No generator reasoning is supplied. " +
			"Return details:{} and account for exactly this checked inventory, without inventing that unavailable evidence was inspected: " + string(inventoryJSON) +
			". Use result=fail or insufficient for missing evidence or any unresolved material issue. Report at most ten material findings, with concise messages. Finding target must be a statement/measure/topic ID or source/draft; evidence must identify an original page or citation. Do not rewrite the draft. Textless pages without original images remain insufficient."
		response, callErr := invoke(role, prompt, evidence{"original-pages": pages, "canonical-draft": draft, "source-metadata": metadata})
		outcome := Outcome{Result: semantic.ResultError, Response: response}
		if callErr != nil {
			outcome.Error = callErr.Error()
		} else {
			outcome.Result = response.Report.Result
			if !completeInventory(response.Report.Checked, inventory) {
				outcome.Result, outcome.Error = semantic.ResultInsufficient, "checked inventory does not match the complete required evidence set"
			}
			if len(ex.Warnings) > 0 && outcome.Result == semantic.ResultPass {
				outcome.Result, outcome.Error = semantic.ResultInsufficient, "unresolved native-page warnings require original visual evidence"
			}
		}
		record.Reports[role] = outcome
		if err := saveJSON(filepath.Join(dir, role+".json"), outcome); err != nil {
			return err
		}
		fmt.Println("   result:", outcome.Result, outcome.Error)
	}
	current, err := regularRead(originalPath, 100<<20)
	if err != nil || !reflect.DeepEqual(current, pdf) {
		return errors.New("original source changed during the run")
	}
	names := requiredRunInputs()
	names = append(names, invocationFiles...)
	for _, name := range names {
		raw, err := regularRead(filepath.Join(dir, name), retainedLimit(name))
		if err != nil {
			return err
		}
		record.Inputs[name] = hashBytes(raw)
	}
	if err := saveJSON(filepath.Join(dir, "run.json"), record); err != nil {
		return err
	}
	fmt.Println("5/5 deterministic private Markdown and supplied-HTML rendering")
	return Render(dir)
}

// Render rechecks all retained bytes and canonical bindings before offline output.
func Render(output string) error {
	dir, err := outputPath(output)
	if err != nil {
		return err
	}
	raw, err := regularRead(filepath.Join(dir, "run.json"), 1<<20)
	if err != nil {
		return err
	}
	var record Record
	if err := strictDecode(raw, &record); err != nil {
		return err
	}
	if record.Version != "atlas.preview.run/v1" || !record.PrivatePreview || record.PublicEligible || len(record.Inputs) < 8 {
		return errors.New("not a private preview run record")
	}
	for _, name := range requiredRunInputs() {
		if _, ok := record.Inputs[name]; !ok {
			return fmt.Errorf("required run evidence missing: %s", name)
		}
	}
	files := map[string][]byte{}
	for name, expected := range record.Inputs {
		if filepath.Base(name) != name || name == "." || strings.ContainsAny(name, "/\\") {
			return errors.New("unsafe run input path")
		}
		raw, err := regularRead(filepath.Join(dir, name), retainedLimit(name))
		if err != nil || hashBytes(raw) != expected {
			return fmt.Errorf("stale/missing run input: %s", name)
		}
		files[name] = raw
	}
	reference, err := regularRead(referencePath, 256<<10)
	if err != nil || hashBytes(reference) != record.ReferenceSHA256 {
		return errors.New("stale or missing supplied reference")
	}
	var ex Extraction
	if err := strictDecode(files["extraction.json"], &ex); err != nil {
		return err
	}
	src, err := sourceRecord(files["metadata.json"], hashBytes(files["source.pdf"]))
	if err != nil {
		return err
	}
	p, err := canonical.ValidateProgramme(files["programme.json"], "", src, canonicalEvidence(&ex, hashBytes(files["extraction.json"])))
	if err != nil {
		return err
	}
	statuses := []string{"BORRADOR · Vista previa privada; no es una publicación aprobada."}
	for _, pair := range [][2]string{{semantic.RoleClaimValidator, "Afirmaciones"}, {semantic.RoleCoverageValidator, "Cobertura"}} {
		stored, ok := files[pair[0]+".json"]
		if !ok {
			return errors.New("independent report missing")
		}
		var outcome Outcome
		if err := strictDecode(stored, &outcome); err != nil {
			return err
		}
		if !reflect.DeepEqual(outcome, record.Reports[pair[0]]) {
			return errors.New("run report binding mismatch")
		}
		expected := map[string]string{}
		for key, file := range map[string]string{"source": "source.pdf", "metadata": "metadata.json", "extraction": "extraction.json", "draft": "programme.json", "taxonomy": "taxonomy.json"} {
			expected[key] = hashBytes(files[file])
		}
		binding, err := artifact.CanonicalValue(expected)
		if err != nil {
			return err
		}
		var invocation struct {
			InputsHash string `json:"inputs_hash"`
		}
		if err := json.Unmarshal(files[pair[0]+"-invocation.json"], &invocation); err != nil || invocation.InputsHash != hashBytes(binding) {
			return errors.New("report invocation references different inputs")
		}
		if outcome.Response != nil && (outcome.Response.Report.Role != pair[0] || outcome.Response.Report.Subject != p.ProgrammeID || !reflect.DeepEqual(outcome.Response.Report.Inputs, expected)) {
			return errors.New("report references different inputs")
		}
		if outcome.Result == semantic.ResultPass {
			if outcome.Response == nil || outcome.Response.Report.Result != semantic.ResultPass ||
				outcome.Response.Report.Role != pair[0] || outcome.Response.Report.Subject != p.ProgrammeID ||
				!completeInventory(outcome.Response.Report.Checked, checkedInventory(p, ex.PageCount)) {
				return errors.New("incomplete or mismatched passing report")
			}
		}
		labels := map[string]string{semantic.ResultPass: "comprobación superada", semantic.ResultFail: "comprobación fallida", semantic.ResultInsufficient: "evidencia insuficiente", semantic.ResultError: "error de ejecución o formato"}
		label, ok := labels[outcome.Result]
		if !ok {
			return errors.New("unknown semantic outcome")
		}
		statuses = append(statuses, pair[1]+": "+label+".")
		if pair[0] == semantic.RoleCoverageValidator && outcome.Response != nil && len(outcome.Response.Report.Findings) > 0 {
			statuses = append(statuses, fmt.Sprintf("El informe de cobertura señala %d posibles incidencias pendientes de resolución.", len(outcome.Response.Report.Findings)))
		}
	}
	missingVisual := []string{}
	for _, page := range ex.Pages {
		if page.Text == "" {
			missingVisual = append(missingVisual, fmt.Sprint(page.Page))
		}
	}
	if len(missingVisual) > 0 {
		statuses = append(statuses, "Pendiente de evidencia visual en las páginas "+strings.Join(missingVisual, ", ")+".")
	}
	status := strings.Join(statuses, " ")
	otherSource := "fixtures/vox/programa.pdf"
	if p.Metadata.Party.ID == "vox" {
		otherSource = "fixtures/frente/programa.pdf"
	}
	pendingPDF, err := filepath.Rel(dir, otherSource)
	if err != nil {
		return err
	}
	htmlBytes, err := renderHTML(reference, views(p, ex.PageCount, status, filepath.ToSlash(pendingPDF)))
	if err != nil {
		return err
	}
	for name, contents := range map[string][]byte{"index.html": htmlBytes, "summary.md": renderMarkdown(p, ex.PageCount, status),
		"pending.md": []byte("# Documento pendiente\n\nEste documento todavía no se ha procesado en esta ejecución.\n")} {
		if err := writeOutput(filepath.Join(dir, name), contents); err != nil {
			return err
		}
	}
	fmt.Println("HTML:", filepath.Join(dir, "index.html"))
	return nil
}
