// Command probe reproduces trust-boundary findings using synthetic local data.
// Exit zero means every listed rejection invariant holds; any failure exits one.
// It is an audit tool, not part of the passing implementation regression gates.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/semantic"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

var failures int

func check(name string, safe bool) {
	if !safe {
		failures++
		fmt.Println("FAIL", name)
		return
	}
	fmt.Println("PASS", name)
}

func hash() string { return "sha256:" + strings.Repeat("a", 64) }

func main() {
	jurisdiction()
	programme()
	jobs()
	fmt.Printf("%d failed rejection invariants\n", failures)
	if failures != 0 {
		os.Exit(1)
	}
}

func jurisdiction() {
	// Out-of-scope without retained scope evidence must reject (no implicit
	// all-outside rule, no guessed exclusion).
	_, err := taxonomy.Classify(taxonomy.ClassifyInput{
		TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true,
		Applicability: taxonomy.AppOutOfScope,
	})
	check("taxonomy rejects out-of-scope without scope evidence", err != nil)
	// Ambiguous applicability is unclear, never a guessed exclusion.
	s, err := taxonomy.Classify(taxonomy.ClassifyInput{
		TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true,
		Applicability: taxonomy.AppUnknown,
	})
	check("taxonomy keeps unknown applicability unclear", err == nil && s == taxonomy.StateUnclear)
	// In-scope omission is not-covered (omission is never opposition).
	s, err = taxonomy.Classify(taxonomy.ClassifyInput{
		TopicID: "healthcare", FullCoverage: true, EvidenceReadable: true,
		Applicability: taxonomy.AppInScope,
	})
	check("taxonomy treats in-scope omission as not-covered", err == nil && s == taxonomy.StateNotCovered)
}

func programme() {
	src := canonical.SourceRecord{
		ProgrammeID: "fixture-prog", SourceChecksum: hash(), SourceMetadataChecksum: hash(),
		Party:        canonical.Party{ID: "fixture-party", Name: "Partido sintético"},
		Scope:        "general-political",
		Geography:    canonical.Geography{ID: "es", Name: "España"},
		DocumentType: "general-political",
		Rights:       canonical.Rights{Publisher: "fixture-publisher"},
	}
	ex := canonical.Extraction{ExtractionChecksum: hash(), SourceChecksum: src.SourceChecksum, PageCount: 1}
	p := canonical.Programme{
		Version: canonical.Version, ProgrammeID: src.ProgrammeID,
		SourceChecksum: src.SourceChecksum, SourceMetadataChecksum: src.SourceMetadataChecksum,
		ExtractionChecksum: ex.ExtractionChecksum, TaxonomyVersion: taxonomy.Release,
		Metadata: canonical.Metadata{
			Party: src.Party, DocumentType: src.DocumentType, Scope: src.Scope,
			Geography: src.Geography, Rights: src.Rights,
		},
		Summary: []canonical.Statement{}, Measures: []canonical.Measure{}, Citations: []canonical.Citation{},
	}
	for _, topic := range taxonomy.Topics() {
		p.Topics = append(p.Topics, canonical.Topic{
			ID: topic.ID, State: "not-covered", Summary: []canonical.Statement{},
			MeasureIDs: []string{}, SecondaryMeasureIDs: []string{},
			Basis: canonical.Basis{PageRanges: []canonical.PageRange{{Start: 1, End: 1}}},
		})
	}
	validate := func() error {
		raw, err := json.Marshal(p)
		if err != nil {
			panic(err)
		}
		_, err = canonical.ValidateProgramme(raw, "", src, ex)
		return err
	}
	if err := validate(); err != nil {
		panic(fmt.Sprintf("canonical baseline rejected: %v", err))
	}
	// Entire-topic exclusion without retained scope evidence must reject.
	p.Topics[0].State = "not-applicable"
	check("canonical rejects not-applicable without scope evidence", validate() != nil)
	p.Topics[0].State = "not-covered"
	// Wrong extraction checksum binding.
	ex.ExtractionChecksum = "sha256:" + strings.Repeat("b", 64)
	check("canonical rejects wrong extraction checksum", validate() != nil)
	ex.ExtractionChecksum = p.ExtractionChecksum
	raw, _ := json.Marshal(p)
	var object map[string]json.RawMessage
	_ = json.Unmarshal(raw, &object)
	delete(object, "summary")
	raw, _ = json.Marshal(object)
	_, err := canonical.ValidateProgramme(raw, "", src, ex)
	check("canonical rejects missing required summary field", err != nil)
	// A draft with fewer than 21 topics rejects without panic.
	p.Topics = p.Topics[:5]
	panicked := false
	func() {
		defer func() { panicked = recover() != nil }()
		err = validate()
	}()
	check("canonical rejects short topic list without panic", !panicked && err != nil)
}

type adapter struct {
	mutate func(*semantic.Response)
	suffix string
	cancel bool
}

func (adapter) Identity() semantic.ModelIdentity {
	return semantic.ModelIdentity{Provider: "offline-fixture", Model: "fixture", Revision: hash()}
}

func (a adapter) Run(ctx context.Context, req semantic.Request, _ semantic.EvidenceReader) ([]byte, error) {
	if a.cancel {
		<-ctx.Done()
	}
	r := semantic.Response{
		Version: semantic.JobVersion, JobID: req.JobID, ModelIdentity: a.Identity(),
		Report: semantic.Report{Version: semantic.ReportVersion, Role: req.Role, Subject: req.Subject,
			Inputs: req.Inputs, Result: semantic.ResultPass, Details: json.RawMessage(`{}`),
			Findings: []semantic.Finding{},
			Checked:  semantic.Checked{Statements: []string{}, Measures: []string{}, Pages: []int{}, Topics: []string{}}},
	}
	if a.mutate != nil {
		a.mutate(&r)
	}
	raw, _ := json.Marshal(r)
	return append(raw, a.suffix...), nil
}

func jobs() {
	req := semantic.Request{
		Version: semantic.JobVersion, JobID: "audit-job", EvaluationSet: "audit-set",
		Role: semantic.RoleClaimValidator, Subject: hash(), Inputs: map[string]string{"draft": hash()},
		Model: "fixture", Prompt: hash(), OutputSchema: hash(),
		Limits: semantic.Limits{Deadline: time.Second},
	}
	run := func(a adapter) error {
		_, err := semantic.NewRunner(a).Run(context.Background(), req, nil)
		return err
	}
	if err := run(adapter{}); err != nil {
		panic(fmt.Sprintf("semantic baseline rejected: %v", err))
	}
	check("runner rejects trailing JSON response", run(adapter{suffix: `{}`}) != nil)
	check("runner rejects pass with blocker", run(adapter{mutate: func(r *semantic.Response) {
		r.Report.Findings = []semantic.Finding{{ID: "f-1", Code: "unsupported", Severity: "blocker",
			Target: "document", Evidence: hash(), Message: "Unsupported claim"}}
	}}) != nil)
	req.Limits.MaxResponseBytes = 1
	check("runner enforces tightened response byte limit", run(adapter{}) != nil)
	req.Limits.MaxResponseBytes = 0
	req.Limits.Deadline = time.Millisecond
	check("runner discards late successful response", run(adapter{cancel: true}) != nil)
}
