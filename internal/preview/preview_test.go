package preview

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/canonical"
	"github.com/jparrill/manifesto-atlas/internal/semantic"
	"github.com/jparrill/manifesto-atlas/internal/taxonomy"
)

func TestCanonicalCitations(t *testing.T) {
	tests := []struct {
		name  string
		quote string
		page  int
		valid bool
	}{
		{"When a quote is a literal Unicode source slice, it should retain exact original-page citations", "Reforma democrática en España", 1, true},
		{"When a model changes the quote, it should reject unsupported citation evidence", "Reforma democrática de España", 1, false},
		{"When a model invents a page, it should reject the draft", "Reforma democrática en España", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Extraction{Version: "atlas.preview.extraction/v1", SourceChecksum: "source", PageCount: 1,
				Pages: []Page{{1, "Reforma democrática en España para mejorar sus instituciones."}}}
			normalized := Normalization{Topics: []DraftTopic{}}
			for _, topic := range taxonomy.Topics() {
				normalized.Topics = append(normalized.Topics, DraftTopic{ID: topic.ID, State: "unclear", Measures: []CitedText{}})
			}
			value := CitedText{Text: "Reformar las instituciones democráticas.", Page: tt.page, Quote: tt.quote}
			normalized.Topics[0].State = "proposals"
			normalized.Topics[0].Summary = &value
			normalized.Topics[0].Measures = []CitedText{value}
			raw, err := json.Marshal(normalized)
			if err != nil {
				t.Fatal(err)
			}
			src := canonical.SourceRecord{ProgrammeID: "source", SourceChecksum: "source", SourceMetadataChecksum: "metadata",
				Party: canonical.Party{ID: "party", Name: "Partido"}, Geography: canonical.Geography{ID: "es", Name: "España"}, Rights: canonical.Rights{Publisher: "Partido"}}
			p, err := assemble(raw, src, ex, "extraction")
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v, err=%v", tt.valid, err)
			}
			if tt.valid && (len(p.Citations) != 2 || p.Citations[0].Excerpt != tt.quote || p.Citations[0].Page != tt.page) {
				t.Fatal("source evidence was lost")
			}
		})
	}
}

func TestHostAssignedEvidence(t *testing.T) {
	for _, tt := range []struct {
		name  string
		id    string
		valid bool
	}{
		{"When a model selects retained evidence, it should bind the exact host-owned quote", "p1-e1", true},
		{"When a model invents an evidence ID, it should reject instead of fabricating a quote", "p1-e999", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			text := "Reforma democrática en España para mejorar sus instituciones."
			ex := &Extraction{SourceChecksum: "source", PageCount: 1, Pages: []Page{{1, text}}}
			n := Normalization{Topics: []DraftTopic{}}
			for _, topic := range taxonomy.Topics() {
				n.Topics = append(n.Topics, DraftTopic{ID: topic.ID, State: "unclear", Measures: []CitedText{}})
			}
			n.Topics[0].State = "mentioned"
			n.Topics[0].Summary = &CitedText{Text: "Reforma de las instituciones.", EvidenceID: tt.id}
			raw, _ := json.Marshal(n)
			src := canonical.SourceRecord{ProgrammeID: "source", SourceChecksum: "source", SourceMetadataChecksum: "metadata",
				Party: canonical.Party{ID: "party", Name: "Partido"}, Geography: canonical.Geography{ID: "es", Name: "España"}, Rights: canonical.Rights{Publisher: "Partido"}}
			p, err := assemble(raw, src, ex, "extraction")
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v err=%v", tt.valid, err)
			}
			if tt.valid && (len(p.Citations) != 1 || p.Citations[0].Excerpt != text) {
				t.Fatal("host-owned evidence changed")
			}
		})
	}
}

func TestSuppliedPresentationBinding(t *testing.T) {
	reference, err := os.ReadFile("../../assets/reference/programmes.html")
	if err != nil {
		t.Fatal(err)
	}
	prefix, suffix, err := referenceParts(reference)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("When hostile canonical text is bound, it should stay escaped data and preserve supplied presentation code", func(t *testing.T) {
		parties := map[string]partyView{"fo": {Name: "Texto <script>alert(1)</script>", Note: "literal & quotation", PDF: "source.pdf", MD: "summary.md", Items: [][2]string{{"Tema", "Texto <img src=x onerror=alert(1)>"}}}}
		first, err := renderHTML(reference, parties)
		if err != nil {
			t.Fatal(err)
		}
		second, err := renderHTML(reference, parties)
		if err != nil || !bytes.Equal(first, second) {
			t.Fatal("unchanged rendering is not byte-deterministic", err)
		}
		if !bytes.HasPrefix(first, prefix) || !bytes.HasSuffix(first, suffix) {
			t.Fatal("owner-supplied presentation code changed")
		}
		var decoded map[string]partyView
		if err := json.Unmarshal(bytes.TrimSpace(first[len(prefix):len(first)-len(suffix)]), &decoded); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(decoded["fo"].Name, "<") || strings.Contains(decoded["fo"].Items[0][1], "<") || !strings.Contains(decoded["fo"].Name, "&lt;") {
			t.Fatal("JSON decoding restored executable markup")
		}
	})
	t.Run("When a generated link has an executable scheme, it should reject rendering", func(t *testing.T) {
		_, err := renderHTML(reference, map[string]partyView{"fo": {PDF: "javascript:alert(1)", MD: "summary.md"}})
		if err == nil {
			t.Fatal("unsafe link accepted")
		}
	})
	t.Run("When the supplied reference structure changes, it should fail instead of inventing markup", func(t *testing.T) {
		_, err := renderHTML([]byte("unrecognized owner asset"), nil)
		if err == nil {
			t.Fatal("unknown reference accepted")
		}
	})
}

func TestCheckedInventory(t *testing.T) {
	want := semantic.Checked{Statements: []string{"s-1"}, Measures: []string{"m-1"}, Topics: []string{"topic"}, Pages: []int{1, 2}}
	for _, tt := range []struct {
		name string
		got  semantic.Checked
		want bool
	}{
		{"When all actual targets are reported in another order, it should accept completeness", semantic.Checked{Statements: []string{"s-1"}, Measures: []string{"m-1"}, Topics: []string{"topic"}, Pages: []int{2, 1}}, true},
		{"When a whole-source page is omitted, it should reject the claimed complete pass", semantic.Checked{Statements: []string{"s-1"}, Measures: []string{"m-1"}, Topics: []string{"topic"}, Pages: []int{1}}, false},
		{"When duplicate targets replace real pages, it should reject coverage", semantic.Checked{Statements: []string{"s-1"}, Measures: []string{"m-1"}, Topics: []string{"topic"}, Pages: []int{1, 1}}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := completeInventory(tt.got, want); got != tt.want {
				t.Fatalf("complete=%v, want=%v", got, tt.want)
			}
		})
	}
}

func TestOutputBoundary(t *testing.T) {
	for _, path := range []string{".", "..", "../escape", "/tmp/absolute"} {
		t.Run("When output escapes the checkout, it should reject "+path, func(t *testing.T) {
			if _, err := outputPath(path); err == nil {
				t.Fatal("unsafe output accepted")
			}
		})
	}
}

func TestEndpointCredentials(t *testing.T) {
	t.Run("When a provider URL embeds credentials, it should reject without retaining or printing them", func(t *testing.T) {
		err := Run(t.Context(), Config{Endpoint: "http://user:private-password@example.invalid/v1"})
		if err == nil || !strings.Contains(err.Error(), "without embedded credentials") || strings.Contains(err.Error(), "private-password") {
			t.Fatalf("embedded transport credentials were not rejected safely: %v", err)
		}
	})
}

func TestMarkdownEscaping(t *testing.T) {
	t.Run("When source excerpts contain markup, it should display literal data rather than create executable Markdown", func(t *testing.T) {
		value := markdownText("<script>bad</script> [link](javascript:bad)\n# injected")
		if strings.Contains(value, "<script>") || strings.Contains(value, "[link](") || strings.Contains(value, "\n#") || !strings.Contains(value, "&lt;script&gt;") {
			t.Fatalf("source syntax was not neutralized: %q", value)
		}
	})
}

func TestMacWorkerBoundaries(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS-specific native enforcement; other platforms have no preview runtime")
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "selected input with spaces.txt")
	if err := os.WriteFile(input, []byte("selected evidence"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Run("When the exact selected input is read, it should succeed in the real sandbox", func(t *testing.T) {
		out, stderr, err := worker(t.Context(), input, defaultWorkerLimits, "/bin/cat", input)
		if err != nil || string(out) != "selected evidence" {
			t.Fatalf("out=%q stderr=%q err=%v", out, stderr, err)
		}
	})
	t.Run("When the worker writes the source, it should deny the write and preserve bytes", func(t *testing.T) {
		_, _, err := worker(t.Context(), input, defaultWorkerLimits, "/bin/sh", "-c", "printf overwrite > \"$1\"", "probe", input)
		current, readErr := os.ReadFile(input)
		if err == nil || readErr != nil || string(current) != "selected evidence" {
			t.Fatalf("write denial failed: err=%v bytes=%q", err, current)
		}
	})
	t.Run("When another private file is read, it should deny access", func(t *testing.T) {
		other := filepath.Join(dir, "private.txt")
		if err := os.WriteFile(other, []byte("unavailable evidence"), 0o600); err != nil {
			t.Fatal(err)
		}
		out, _, err := worker(t.Context(), input, defaultWorkerLimits, "/bin/cat", other)
		if err == nil || strings.Contains(string(out), "unavailable evidence") {
			t.Fatal("worker read an unselected private file")
		}
	})
	t.Run("When transport credentials exist on the host, it should not inherit them", func(t *testing.T) {
		t.Setenv("ATLAS_API_KEY", "host-only-test-secret")
		out, stderr, err := worker(t.Context(), input, defaultWorkerLimits, "/usr/bin/env")
		if err != nil || bytes.Contains(out, []byte("ATLAS_API_KEY")) || bytes.Contains(out, []byte("host-only-test-secret")) {
			t.Fatalf("host credential isolation failed: err=%v stderr=%q", err, stderr)
		}
	})
	t.Run("When the worker inspects the host parent process, it should deny access", func(t *testing.T) {
		out, _, err := worker(t.Context(), input, defaultWorkerLimits, "/bin/ps", "-p", strconv.Itoa(os.Getpid()), "-o", "command=")
		if err == nil || strings.Contains(string(out), "preview.test") {
			t.Fatal("worker could inspect its host parent process")
		}
	})
	t.Run("When network access is attempted, it should deny the connection to a real local listener", func(t *testing.T) {
		var hits atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits.Add(1) }))
		defer server.Close()
		_, _, err := worker(t.Context(), input, defaultWorkerLimits, "/usr/bin/curl", "--noproxy", "*", "--max-time", "2", server.URL)
		if err == nil || hits.Load() != 0 {
			t.Fatal("sandbox worker connected to the network")
		}
	})
	t.Run("When output exceeds its budget, it should terminate the worker group", func(t *testing.T) {
		_, _, err := worker(t.Context(), input, workerLimits{Timeout: 3 * time.Second, Output: 1024, RSS: 1 << 30}, "/usr/bin/yes")
		if err == nil || !strings.Contains(err.Error(), "output budget") {
			t.Fatalf("output limit failed: %v", err)
		}
	})
	t.Run("When the deadline expires, it should cancel the actual process group", func(t *testing.T) {
		start := time.Now()
		out, _, err := worker(context.Background(), input, workerLimits{Timeout: 300 * time.Millisecond, Output: 1024, RSS: 1 << 30, AllowFork: true}, "/bin/sh", "-c", "sleep 20 & echo $!; wait")
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > 3*time.Second {
			t.Fatalf("timeout/cleanup failed: %v", err)
		}
		pid, parseErr := strconv.Atoi(strings.TrimSpace(string(out)))
		if parseErr != nil || pid < 1 {
			t.Fatalf("child process did not actually start: %q", out)
		}
		state, _ := exec.Command("/bin/ps", "-o", "stat=", "-p", strconv.Itoa(pid)).Output()
		if current := strings.TrimSpace(string(state)); current != "" && !strings.HasPrefix(current, "Z") {
			t.Fatalf("descendant survived cancellation: %d %s", pid, current)
		}
	})
	t.Run("When sampled RSS exceeds the configured ceiling, it should terminate actual memory growth", func(t *testing.T) {
		_, _, err := worker(t.Context(), input, workerLimits{Timeout: 4 * time.Second, Output: 1024, RSS: 32 << 20}, "/usr/bin/perl", "-e", "$x = 'x' x (96*1024*1024); sleep 10;")
		if err == nil || !strings.Contains(err.Error(), "RSS budget") {
			t.Fatalf("actual RSS enforcement failed: %v", err)
		}
	})
	t.Run("When a default PDF worker tries to fork, it should deny creating extra processes", func(t *testing.T) {
		_, stderr, err := worker(t.Context(), input, defaultWorkerLimits, "/usr/bin/perl", "-e", "defined(fork()) or die 'fork denied';")
		if err == nil || !strings.Contains(string(stderr), "fork denied") {
			t.Fatalf("default worker fork restriction failed: err=%v stderr=%q", err, stderr)
		}
	})
}

func TestOfflineRunTampering(t *testing.T) {
	t.Run("When a retained input changes, it should reject offline rendering before using stale reports", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if err := os.Mkdir("run", 0o755); err != nil {
			t.Fatal(err)
		}
		record := Record{Version: "atlas.preview.run/v1", PrivatePreview: true, Inputs: map[string]string{}, Reports: map[string]Outcome{}}
		for _, name := range requiredRunInputs() {
			if err := os.WriteFile(filepath.Join("run", name), []byte("retained evidence"), 0o600); err != nil {
				t.Fatal(err)
			}
			record.Inputs[name] = hashBytes([]byte("retained evidence"))
		}
		if err := saveJSON("run/run.json", record); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("run/programme.json", []byte("changed evidence"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := Render("run"); err == nil || !strings.Contains(err.Error(), "stale/missing run input: programme.json") {
			t.Fatalf("stale input not detected at its trust boundary: %v", err)
		}
	})
	t.Run("When output uses a symlink, it should reject the escaped directory", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if err := os.Symlink(t.TempDir(), "linked-output"); err != nil {
			t.Fatal(err)
		}
		if _, err := outputPath("linked-output/preview"); err == nil {
			t.Fatal("symlink output accepted")
		}
	})
}
