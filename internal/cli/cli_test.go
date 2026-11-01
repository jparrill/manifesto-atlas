package cli

import (
	"bytes"
	"errors"
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
	r := doRun("build", "programme.json")
	if r.code != exitFailed {
		t.Fatalf("stage exit = %d, want %d", r.code, exitFailed)
	}
	if !strings.Contains(r.stdout, "not available yet") {
		t.Errorf("stage output missing not-available:\n%s", r.stdout)
	}
	if !strings.Contains(r.stdout, "TQ-0043") {
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
