package preview

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Explicit real-tool gate, separate from ordinary offline stored-output replay.
func TestActualMacSource(t *testing.T) {
	if runtime.GOOS != "darwin" || os.Getenv("ATLAS_EVALUATE_REAL_PDF") != "1" {
		t.Skip("explicit gate: ATLAS_EVALUATE_REAL_PDF=1 go test ./internal/preview -run TestActualMacSource -count=1")
	}
	descriptor, err := filepath.Abs("../../fixtures/frente/source.json")
	if err != nil {
		t.Fatal(err)
	}
	c, _, pdf, err := loadCandidate(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := extract(t.Context(), filepath.Join(filepath.Dir(descriptor), c.CandidatePath), t.TempDir(), hashBytes(pdf))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("../../runs/preview/frente-e2e/extraction.json")
	if err != nil {
		t.Fatal(err)
	}
	var stored Extraction
	if err := json.Unmarshal(raw, &stored); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, &stored) {
		t.Fatal("real original pages, source identity, warnings or tools drifted from the retained run")
	}
	t.Logf("Actual bounded native extraction matches %d retained original pages with %d explicit warnings", actual.PageCount, len(actual.Warnings))
}
