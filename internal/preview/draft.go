package preview

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

type CitedText struct {
	Text       string `json:"text"`
	EvidenceID string `json:"evidence_id,omitempty"`
	Page       int    `json:"page,omitempty"`
	Quote      string `json:"quote,omitempty"`
}

type EvidenceSpan struct {
	ID   string `json:"id"`
	Page int    `json:"page"`
	Text string `json:"text"`
}

// evidenceSpans labels all original text, with overlapping bounded source slices.
// Models select existing IDs rather than attempting to retype literal quotations.
func evidenceSpans(ex *Extraction) []EvidenceSpan {
	spans := []EvidenceSpan{}
	for _, page := range ex.Pages {
		words := strings.Fields(page.Text)
		for start, ordinal := 0, 1; start < len(words); start, ordinal = start+32, ordinal+1 {
			end := min(start+40, len(words))
			spans = append(spans, EvidenceSpan{ID: fmt.Sprintf("p%d-e%d", page.Page, ordinal), Page: page.Page, Text: strings.Join(words[start:end], " ")})
			if end == len(words) {
				break
			}
		}
	}
	return spans
}

type DraftTopic struct {
	ID       string      `json:"id"`
	State    string      `json:"state"`
	Summary  *CitedText  `json:"summary"`
	Measures []CitedText `json:"measures"`
}

type Normalization struct {
	Topics []DraftTopic `json:"topics"`
}

func strictDecode(raw []byte, value any) error {
	canon, err := artifact.Canonicalize(raw)
	if err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(canon))
	d.DisallowUnknownFields()
	return d.Decode(value)
}

func canonicalEvidence(ex *Extraction, checksum string) canonical.Extraction {
	result := canonical.Extraction{ExtractionChecksum: checksum, SourceChecksum: ex.SourceChecksum,
		PageCount: ex.PageCount, BlockText: map[string]string{}, BlockPage: map[string]int{}}
	for _, p := range ex.Pages {
		id := fmt.Sprintf("page-%d", p.Page)
		result.BlockText[id], result.BlockPage[id] = p.Text, p.Page
	}
	return result
}

// assemble does not repair model quotations. Only literal original-page slices
// can become citations; changed words, invented pages and ellipses fail closed.
func assemble(raw []byte, src canonical.SourceRecord, ex *Extraction, extractionHash string) (*canonical.Programme, error) {
	var normalized Normalization
	if err := strictDecode(raw, &normalized); err != nil {
		return nil, err
	}
	registry := taxonomy.Topics()
	if len(normalized.Topics) != len(registry) {
		return nil, errors.New("normalizer must return all 21 topics")
	}
	p := &canonical.Programme{Version: canonical.Version, ProgrammeID: src.ProgrammeID,
		SourceChecksum: src.SourceChecksum, SourceMetadataChecksum: src.SourceMetadataChecksum,
		ExtractionChecksum: extractionHash, TaxonomyVersion: taxonomy.Release,
		Metadata: canonical.Metadata{Party: src.Party, DocumentType: src.DocumentType, Scope: src.Scope,
			Geography: src.Geography, Election: src.Election, PublicationDate: src.PublicationDate, Rights: src.Rights},
		Summary: []canonical.Statement{}, Topics: []canonical.Topic{}, Measures: []canonical.Measure{}, Citations: []canonical.Citation{}}
	statement := 0
	spans := map[string]EvidenceSpan{}
	for _, span := range evidenceSpans(ex) {
		spans[span.ID] = span
	}
	cited := func(text CitedText) (canonical.Statement, error) {
		if text.EvidenceID != "" {
			span, ok := spans[text.EvidenceID]
			if !ok || text.Page != 0 || text.Quote != "" {
				return canonical.Statement{}, fmt.Errorf("unknown or ambiguous original evidence ID: %s", text.EvidenceID)
			}
			text.Page, text.Quote = span.Page, span.Text
		}
		if text.Page < 1 || text.Page > len(ex.Pages) || text.Quote == "" {
			return canonical.Statement{}, errors.New("missing or invalid source quotation")
		}
		original := ex.Pages[text.Page-1].Text
		index := strings.Index(original, text.Quote)
		if index < 0 {
			return canonical.Statement{}, fmt.Errorf("quote not literal on page %d: %.100s", text.Page, text.Quote)
		}
		start := utf8.RuneCountInString(original[:index])
		cid := fmt.Sprintf("c-%d", len(p.Citations)+1)
		p.Citations = append(p.Citations, canonical.Citation{ID: cid, Page: text.Page,
			BlockID: fmt.Sprintf("page-%d", text.Page), Start: start, End: start + utf8.RuneCountInString(text.Quote), Excerpt: text.Quote})
		statement++
		return canonical.Statement{ID: fmt.Sprintf("s-%d", statement), Text: text.Text, Citations: []string{cid}}, nil
	}
	for i, t := range normalized.Topics {
		if t.ID != registry[i].ID || t.State == "not-applicable" {
			return nil, fmt.Errorf("unknown/order/scope topic: %s", t.ID)
		}
		topic := canonical.Topic{ID: t.ID, State: t.State, Summary: []canonical.Statement{}, MeasureIDs: []string{},
			SecondaryMeasureIDs: []string{}, Basis: canonical.Basis{PageRanges: []canonical.PageRange{{Start: 1, End: ex.PageCount}}}}
		if t.Summary != nil {
			s, err := cited(*t.Summary)
			if err != nil {
				return nil, err
			}
			topic.Summary = append(topic.Summary, s)
		}
		p.Topics = append(p.Topics, topic)
	}
	for i, t := range normalized.Topics {
		for _, m := range t.Measures {
			action, err := cited(m)
			if err != nil {
				return nil, err
			}
			mid := fmt.Sprintf("m-%d", len(p.Measures)+1)
			p.Measures = append(p.Measures, canonical.Measure{ID: mid, PrimaryTopic: t.ID, Action: action,
				Conditions: []canonical.Statement{}, Citations: append([]string{}, action.Citations...)})
			p.Topics[i].MeasureIDs = append(p.Topics[i].MeasureIDs, mid)
		}
	}
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return canonical.ValidateProgramme(data, "", src, canonicalEvidence(ex, extractionHash))
}
