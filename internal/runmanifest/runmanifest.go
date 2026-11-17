// Package runmanifest owns the run contract: run.json, its documents and the
// curated claim list. Both the gate (internal/runcheck) and the run builder
// read the manifest from here so there is a single definition.
package runmanifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Document is one programme inside a run.
type Document struct {
	Programme string `json:"programme"`
	Source    string `json:"source"`  // page-marked extraction, e.g. pp.md
	Summary   string `json:"summary"` // factual summary, e.g. Resumen-pp.md
	Proposals int    `json:"proposals"`
	// Coverage selects how proposals are checked: "numbers" (the document
	// numbers 1..N once) or "items" (numbering restarts in every section).
	Coverage string `json:"coverage"`
	// Title is the display name used by the view.
	Title string `json:"title,omitempty"`
	// PDF is the source file name, used for the provenance link.
	PDF string `json:"pdf,omitempty"`
	// Scope is the human label of the document's territorial nature.
	Scope string `json:"scope,omitempty"`
}

// Source is the provenance of one immutable PDF.
type Source struct {
	Programme string `json:"programme"`
	URL       string `json:"url"`
	File      string `json:"file"`
	Bytes     int64  `json:"bytes"`
	SHA256    string `json:"sha256"`
	Pages     int    `json:"pages"`
}

// Spec is the run manifest, run.json.
type Spec struct {
	Name      string     `json:"name"`
	Claims    string     `json:"claims"`
	Allow     string     `json:"allow"`
	View      string     `json:"view"`
	Documents []Document `json:"documents"`
	Sources   []Source   `json:"sources,omitempty"`
}

// Path returns the manifest path of a run directory.
func Path(dir string) string { return filepath.Join(dir, "run.json") }

// Load reads and validates a manifest.
func Load(dir string) (Spec, error) {
	raw, err := os.ReadFile(Path(dir))
	if err != nil {
		return Spec{}, fmt.Errorf("read run.json: %w", err)
	}
	var spec Spec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return Spec{}, fmt.Errorf("parse run.json: %w", err)
	}
	if spec.Name == "" || len(spec.Documents) == 0 {
		return Spec{}, fmt.Errorf("run.json needs a name and at least one document")
	}
	for _, doc := range spec.Documents {
		if doc.Programme == "" || doc.Source == "" || doc.Summary == "" || doc.Proposals <= 0 {
			return Spec{}, fmt.Errorf("document %+v is incomplete", doc)
		}
	}
	return spec, nil
}

// Save writes the manifest with a stable, indented encoding.
func Save(dir string, spec Spec) error {
	raw, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(Path(dir), append(raw, '\n'), 0o644)
}

// Claim is one curated factual datum of the run.
type Claim struct {
	Programme string
	Pages     string
	Phrase    string
}

// LoadClaims reads a claims file of tab-separated programme, pages and phrase.
func LoadClaims(path string) ([]Claim, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filepath.Base(path), err)
	}
	var claims []Claim
	for i, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 3 {
			return nil, fmt.Errorf("%s:%d: want 3 tab-separated fields", filepath.Base(path), i+1)
		}
		claims = append(claims, Claim{Programme: f[0], Pages: f[1], Phrase: f[2]})
	}
	return claims, nil
}
