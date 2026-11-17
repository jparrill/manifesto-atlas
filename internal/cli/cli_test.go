package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type result struct {
	code   int
	stdout string
	stderr string
}

func doRun(args ...string) result {
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return result{code: code, stdout: out.String(), stderr: errb.String()}
}

func TestHelp(t *testing.T) {
	r := doRun("help")
	if r.code != exitOK {
		t.Fatalf("help exit = %d, want %d", r.code, exitOK)
	}
	if !strings.Contains(r.stdout, "Usage:") {
		t.Errorf("help output missing Usage:\n%s", r.stdout)
	}
	for _, s := range stages {
		if !strings.Contains(r.stdout, s.name) {
			t.Errorf("help output missing stage %q", s.name)
		}
	}
}

func TestNoArgsShowsHelp(t *testing.T) {
	r := doRun()
	if r.code != exitOK {
		t.Fatalf("no-args exit = %d, want %d", r.code, exitOK)
	}
	if !strings.Contains(r.stdout, "Usage:") {
		t.Errorf("no-args output missing Usage:\n%s", r.stdout)
	}
}

func TestVersion(t *testing.T) {
	r := doRun("version")
	if r.code != exitOK {
		t.Fatalf("version exit = %d, want %d", r.code, exitOK)
	}
	if !strings.Contains(r.stdout, "manifesto-atlas "+Version) {
		t.Errorf("version output missing version:\n%s", r.stdout)
	}
	if !strings.Contains(r.stdout, "go go") {
		t.Errorf("version output missing go runtime:\n%s", r.stdout)
	}
}

func TestUnknownCommand(t *testing.T) {
	r := doRun("bogus")
	if r.code != exitUsage {
		t.Fatalf("unknown command exit = %d, want %d", r.code, exitUsage)
	}
	if !strings.Contains(r.stderr, "unknown command") {
		t.Errorf("unknown command stderr missing message:\n%s", r.stderr)
	}
}

func TestInvalidFlag(t *testing.T) {
	r := doRun("version", "--bogus")
	if r.code != exitUsage {
		t.Fatalf("invalid flag exit = %d, want %d", r.code, exitUsage)
	}
	if !strings.Contains(r.stderr, "flag provided but not defined") {
		t.Errorf("invalid flag stderr missing message:\n%s", r.stderr)
	}
}

func TestHelpFlagShowsUsage(t *testing.T) {
	r := doRun("submit", "--help")
	if r.code != exitOK {
		t.Fatalf("--help exit = %d, want %d", r.code, exitOK)
	}
	if !strings.Contains(r.stdout, "Usage:") {
		t.Errorf("--help output missing Usage:\n%s", r.stdout)
	}
}

func TestMissingRequiredArgument(t *testing.T) {
	r := doRun("submit")
	if r.code != exitUsage {
		t.Fatalf("missing arg exit = %d, want %d", r.code, exitUsage)
	}
	if !strings.Contains(r.stderr, "target-path") {
		t.Errorf("missing arg stderr missing message:\n%s", r.stderr)
	}
}

func TestStageNotAvailable(t *testing.T) {
	// build, extract and validate-content are implemented; the publication stages
	// still declare their owning task.
	r := doRun("normalize", "programme.json")
	if r.code != exitFailed {
		t.Fatalf("stage exit = %d, want %d", r.code, exitFailed)
	}
	if !strings.Contains(r.stdout, "not available yet") {
		t.Errorf("stage output missing not-available:\n%s", r.stdout)
	}
	if !strings.Contains(r.stdout, "TQ-0038") {
		t.Errorf("stage output missing owning task:\n%s", r.stdout)
	}
}

// TestPathWithSpaces proves a target path containing spaces is treated as a
// single argument (no word splitting) and echoed intact.
func TestPathWithSpaces(t *testing.T) {
	const path = "/path/with spaces/file.pdf"
	r := doRun("submit", path)
	if r.code != exitFailed {
		t.Fatalf("path-with-spaces exit = %d, want %d", r.code, exitFailed)
	}
	if !strings.Contains(r.stdout, path) {
		t.Errorf("path with spaces not preserved intact; output:\n%s", r.stdout)
	}
}

// TestHostilePathTreatedAsData proves hostile argument values are treated as
// opaque data, never interpolated or executed.
func TestHostilePathTreatedAsData(t *testing.T) {
	for _, hostile := range []string{"; rm -rf /", "$(reboot)", "`id`", "a;b", "*"} {
		r := doRun("submit", hostile)
		if r.code != exitFailed {
			t.Fatalf("hostile %q exit = %d, want %d", hostile, r.code, exitFailed)
		}
		if !strings.Contains(r.stdout, hostile) {
			t.Errorf("hostile value %q not echoed as data; output:\n%s", hostile, r.stdout)
		}
	}
}

// buildBinary compiles the real entry point into a temp dir and returns its
// path. It fails the test if the build cannot run.
func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "manifesto-atlas")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/manifesto-atlas")
	cmd.Dir = "../.." // module root, relative to the package dir go test uses
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v; output:\n%s", err, out)
	}
	return bin
}

// TestProcessExitCodes runs the compiled binary and proves the process exit
// code propagates the dispatcher result (help=0, stage failure=1, usage=2).
func TestProcessExitCodes(t *testing.T) {
	bin := buildBinary(t)
	cases := []struct {
		name string
		args []string
		want int
	}{
		{"version", []string{"version"}, 0},
		{"stage", []string{"submit", "x.pdf"}, 1},
		{"unknown", []string{"bogus"}, 2},
		{"missing-arg", []string{"submit"}, 2},
		{"invalid-flag", []string{"version", "--bogus"}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(bin, tc.args...)
			out, err := cmd.CombinedOutput()
			var ee *exec.ExitError
			switch {
			case err == nil:
				if tc.want != 0 {
					t.Fatalf("exit = 0, want %d; output:\n%s", tc.want, out)
				}
			case errors.As(err, &ee):
				if ee.ExitCode() != tc.want {
					t.Fatalf("exit = %d, want %d; output:\n%s", ee.ExitCode(), tc.want, out)
				}
			default:
				t.Fatalf("failed to run binary: %v; output:\n%s", err, out)
			}
		})
	}
}

// TestProcessPathWithSpaces proves the compiled binary, invoked with a path
// containing spaces, echoes it intact (no shell interpolation in the CLI).
func TestProcessPathWithSpaces(t *testing.T) {
	bin := buildBinary(t)
	const path = "/tmp/opencode/dir with spaces/file.pdf"
	cmd := exec.Command(bin, "submit", path)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), path) {
		t.Fatalf("path with spaces not preserved by binary; output:\n%s", out)
	}
}

// TestNoArtifactsOnFailure proves a not-available stage and a usage error
// create no files in the working directory.
func TestNoArtifactsOnFailure(t *testing.T) {
	before := listFiles(t)
	doRun("submit", "target.pdf")
	doRun("submit")
	doRun("bogus")
	after := listFiles(t)
	if before != after {
		t.Fatalf("file count changed: before=%d after=%d", before, after)
	}
}

func listFiles(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	return len(entries)
}

// stageEnv builds fake pdftotext/pdfinfo tools so the deterministic flow can be
// exercised end to end without a real PDF or network access.
func stageEnv(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	tools := map[string]string{
		"pdfinfo":   "#!/bin/sh\necho 'Title: prueba'\necho 'Pages:          2'\n",
		"pdftotext": "#!/bin/sh\nprintf '01 VIVIENDA\\n1. Construir vivienda pública.\\n2. Prohibir los pisos turísticos.\\fINTRODUCCIÓN\\nTexto de la introducción.\\f'\n",
	}
	for name, body := range tools {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(body), 0o755); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	pdf := filepath.Join(t.TempDir(), "programa.pdf")
	if err := os.WriteFile(pdf, []byte("%PDF-1.4 prueba"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	return pdf
}

func TestRunStagesEndToEnd(t *testing.T) {
	pdf := stageEnv(t)
	runDir := t.TempDir()

	t.Run("When extract runs on a PDF it should write provenance, extraction and a skeleton", func(t *testing.T) {
		res := doRun("extract", pdf, "-url", "https://example.org/p.pdf",
			"-programme", "pp", "-run", runDir, "-title", "PP")
		if res.code != 0 {
			t.Fatalf("exit = %d, stdout=%s stderr=%s", res.code, res.stdout, res.stderr)
		}
		for _, want := range []string{"run.json", "programa.md", "text/programa.txt", "Resumen-PP.md"} {
			if _, err := os.Stat(filepath.Join(runDir, want)); err != nil {
				t.Fatalf("missing %s: %v", want, err)
			}
		}
		marked, err := os.ReadFile(filepath.Join(runDir, "programa.md"))
		if err != nil {
			t.Fatalf("read extraction: %v", err)
		}
		if !strings.Contains(string(marked), "<!-- page 2 -->") {
			t.Fatalf("page markers missing: %s", marked)
		}
		manifest, err := os.ReadFile(filepath.Join(runDir, "run.json"))
		if err != nil {
			t.Fatalf("read manifest: %v", err)
		}
		if !strings.Contains(string(manifest), "https://example.org/p.pdf") {
			t.Fatalf("provenance not recorded: %s", manifest)
		}
	})

	t.Run("When the skeleton is unfilled the gate should refuse the run", func(t *testing.T) {
		res := doRun("validate-content", runDir)
		if res.code == 0 {
			t.Fatalf("the gate accepted an unwritten summary: %s", res.stdout)
		}
		if !strings.Contains(res.stdout, "runcheck.scaffold") {
			t.Fatalf("scaffold issue missing: %s", res.stdout)
		}
	})

	t.Run("When the summary is written and cited the gate should pass", func(t *testing.T) {
		summary := "# PP — resumen\n\n## 01. Vivienda — p. 1\n\n" +
			"1. «Construir vivienda pública» y prohibir los pisos turísticos [p. 1]\n" +
			"2. Texto de la introducción sobre el parque público [p. 2]\n"
		if err := os.WriteFile(filepath.Join(runDir, "Resumen-PP.md"), []byte(summary), 0o644); err != nil {
			t.Fatalf("write summary: %v", err)
		}
		res := doRun("validate-content", runDir)
		if res.code != 0 {
			t.Fatalf("gate failed: exit=%d stdout=%s", res.code, res.stdout)
		}
		if !strings.Contains(res.stdout, "quotes") {
			t.Fatalf("unexpected output: %s", res.stdout)
		}
		if _, err := os.Stat(filepath.Join(runDir, "Verificacion.md")); err != nil {
			t.Fatalf("report not written: %v", err)
		}
	})

	t.Run("When build runs it should embed the documents in one file", func(t *testing.T) {
		out := filepath.Join(runDir, "vista.html")
		res := doRun("build", runDir, "-out", out, "-template", "../../templates/view.html")
		if res.code != 0 {
			t.Fatalf("build failed: exit=%d stdout=%s stderr=%s", res.code, res.stdout, res.stderr)
		}
		html, err := os.ReadFile(out)
		if err != nil {
			t.Fatalf("read view: %v", err)
		}
		for _, want := range []string{"<!-- generated by manifesto-atlas build -->", `"programmes"`, "Construir vivienda"} {
			if !strings.Contains(string(html), want) {
				t.Fatalf("view misses %q", want)
			}
		}
		verify := doRun("validate-content", runDir)
		if verify.code != 0 || !strings.Contains(verify.stdout, "embedding") {
			t.Fatalf("gate does not see the view: %s", verify.stdout)
		}
	})

	t.Run("When extract is called without origin it should be a usage error", func(t *testing.T) {
		if res := doRun("extract", pdf, "-programme", "pp", "-run", t.TempDir()); res.code != 2 {
			t.Fatalf("exit = %d, want 2", res.code)
		}
	})
}

func TestProcessRunsEverything(t *testing.T) {
	pdf := stageEnv(t)
	runDir := t.TempDir()
	// A model stub that answers with a citable summary for this fixture.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "stub",
			"choices": []map[string]any{{"message": map[string]string{"content": "# PP\n\n## 01. Vivienda — p. 1\n\n" +
				"1. «Construir vivienda pública» y prohibir los pisos turísticos [p. 1]\n" +
				"2. Texto de la introducción sobre el parque público [p. 2]\n"}}},
		})
	}))
	defer srv.Close()

	res := doRun("process", pdf, "-url", "https://example.org/p.pdf", "-programme", "pp",
		"-run", runDir, "-title", "PP", "-endpoint", srv.URL, "-model", "stub",
		"-template", "../../templates/view.html")
	if res.code != 0 {
		t.Fatalf("process failed: exit=%d\nstdout=%s\nstderr=%s", res.code, res.stdout, res.stderr)
	}
	for _, want := range []string{"process: extract", "process: summarize", "process: validate-content", "process: build", "process: done"} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("missing %q in:\n%s", want, res.stdout)
		}
	}
	for _, file := range []string{"run.json", "SOURCES.md", "Resumen-PP.md", "Verificacion.md", "Programas-electorales.html"} {
		if _, err := os.Stat(filepath.Join(runDir, file)); err != nil {
			t.Fatalf("missing %s: %v", file, err)
		}
	}
}
