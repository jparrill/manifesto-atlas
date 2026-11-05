package preview

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/candidate"
)

type Page struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}

type Extraction struct {
	Version        string            `json:"version"`
	SourceChecksum string            `json:"source_checksum"`
	PageCount      int               `json:"page_count"`
	Pages          []Page            `json:"pages"`
	Warnings       []string          `json:"warnings"`
	Tools          map[string]string `json:"tools"`
}

func regularRead(path string, max int64) ([]byte, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() || fi.Size() > max {
		return nil, fmt.Errorf("not a bounded regular file: %s", filepath.Base(path))
	}
	return os.ReadFile(path)
}

func saveJSON(path string, value any) error {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeOutput(path, append(raw, '\n'))
}

// Replace generated files atomically: never follow an output symlink or mutate
// another retained file through a hard link to an existing output inode.
func writeOutput(path string, raw []byte) error {
	fi, err := os.Lstat(path)
	if err == nil && !fi.Mode().IsRegular() {
		return fmt.Errorf("unsafe output leaf: %s", filepath.Base(path))
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".atlas-output-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := f.Chmod(0o644); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.Write(raw); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func retainedLimit(name string) int64 {
	if name == "source.pdf" {
		return 100 << 20
	}
	return 16 << 20
}

func loadCandidate(descriptor string) (*candidate.Candidate, []byte, []byte, error) {
	raw, err := regularRead(descriptor, 1<<20)
	if err != nil {
		return nil, nil, nil, err
	}
	c, err := candidate.Parse(raw)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := candidate.VerifyFiles(filepath.Dir(descriptor), c); err != nil {
		return nil, nil, nil, err
	}
	var embedded struct {
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(raw, &embedded); err != nil {
		return nil, nil, nil, err
	}
	metadata, err := artifact.Canonicalize(embedded.Metadata)
	if err != nil {
		return nil, nil, nil, err
	}
	pdf, err := regularRead(filepath.Join(filepath.Dir(descriptor), c.CandidatePath), 100<<20)
	if err != nil {
		return nil, nil, nil, err
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		return nil, nil, nil, errors.New("candidate is not a PDF")
	}
	if strings.TrimPrefix(hashBytes(pdf), "sha256:") != c.Metadata.Source.SHA256 || len(pdf) != c.Metadata.Source.Size {
		return nil, nil, nil, errors.New("candidate bytes changed after descriptor verification")
	}
	return c, metadata, pdf, nil
}

func extract(ctx context.Context, pdfPath, output string, checksum string) (*Extraction, error) {
	info, err := exec.LookPath("pdfinfo")
	if err != nil {
		return nil, err
	}
	text, err := exec.LookPath("pdftotext")
	if err != nil {
		return nil, err
	}
	rawInfo, stderr, err := worker(ctx, pdfPath, defaultWorkerLimits, info, pdfPath)
	if err != nil {
		return nil, fmt.Errorf("PDF inspection: %w (%s)", err, stderr)
	}
	pageMatch := regexp.MustCompile(`(?m)^Pages:\s+(\d+)`).FindSubmatch(rawInfo)
	if len(pageMatch) != 2 || regexp.MustCompile(`(?m)^Encrypted:\s+yes`).Match(rawInfo) {
		return nil, errors.New("missing page count or encrypted PDF")
	}
	count, _ := strconv.Atoi(string(pageMatch[1]))
	if count < 1 || count > 500 {
		return nil, errors.New("PDF page budget exceeded")
	}
	rawText, stderr, err := worker(ctx, pdfPath, defaultWorkerLimits, text, "-layout", "-enc", "UTF-8", pdfPath, "-")
	if err != nil {
		return nil, fmt.Errorf("PDF extraction: %w (%s)", err, stderr)
	}
	parts := strings.Split(string(rawText), "\f")
	if len(parts) > 0 && strings.TrimSpace(parts[len(parts)-1]) == "" {
		parts = parts[:len(parts)-1]
	}
	if len(parts) != count {
		return nil, fmt.Errorf("page count mismatch: inspected %d, extracted %d", count, len(parts))
	}
	ex := &Extraction{Version: "atlas.preview.extraction/v1", SourceChecksum: checksum, PageCount: count,
		Pages: []Page{}, Warnings: []string{}, Tools: map[string]string{"runtime": "macOS sandbox-exec; sampled RSS watchdog"}}
	for i, part := range parts {
		page := Page{i + 1, strings.Join(strings.Fields(part), " ")}
		ex.Pages = append(ex.Pages, page)
		if page.Text == "" {
			ex.Warnings = append(ex.Warnings, fmt.Sprintf("page %d has no native text; visual evidence required", page.Page))
		}
	}
	for name, tool := range map[string]string{"pdfinfo": info, "pdftotext": text} {
		out, errout, err := worker(ctx, pdfPath, defaultWorkerLimits, tool, "-v")
		if err != nil {
			return nil, err
		}
		ex.Tools[name] = strings.TrimSpace(string(append(out, errout...)))
	}
	if err := writeOutput(filepath.Join(output, "extraction.txt"), rawText); err != nil {
		return nil, err
	}
	if err := writeOutput(filepath.Join(output, "inspection.txt"), rawInfo); err != nil {
		return nil, err
	}
	if err := saveJSON(filepath.Join(output, "extraction.json"), ex); err != nil {
		return nil, err
	}
	return ex, nil
}
