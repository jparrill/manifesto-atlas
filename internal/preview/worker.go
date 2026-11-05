// Package preview implements the explicit private macOS demonstration seam.
// It does not implement or authorize production catalogue publication.
package preview

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type workerLimits struct {
	Timeout   time.Duration
	Output    int
	RSS       int64
	AllowFork bool
}

var defaultWorkerLimits = workerLimits{Timeout: 120 * time.Second, Output: 8 << 20, RSS: 1 << 30}

// The worker may read libraries/fonts and its one selected input, but not home
// data, arbitrary repository files, general network or writable repository paths.
func sandboxPolicy(input string, allowFork bool) string {
	fork := ""
	if allowFork {
		fork = "(allow process-fork)"
	}
	return `(version 1)(deny default)
(allow process-exec)(allow process-info* (target self))` + fork + `
(allow sysctl-read)(allow mach-lookup)(allow file-map-executable)
(allow file-read-metadata)
(allow file-read-data (literal "/"))
(allow file-read* (subpath "/System") (subpath "/usr") (subpath "/bin") (subpath "/sbin") (subpath "/Library")
 (subpath "/opt/homebrew") (subpath "/private/var/db")
 (subpath "/private/etc") (subpath "/private/var/run") (subpath "/dev")
 (literal ` + strconv.Quote(input) + `))
(allow file-write* (literal "/dev/null"))`
}

type cappedWriter struct {
	mu    sync.Mutex
	buf   bytes.Buffer
	max   int
	abort func(string)
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf.Len()+len(p) > w.max {
		w.abort("worker output budget exceeded")
		return 0, errors.New("output budget exceeded")
	}
	return w.buf.Write(p)
}

// worker bounds the whole process group; RSS enforcement is sampled, not a
// kernel hard limit. This limitation is part of the recorded runtime contract.
func worker(ctx context.Context, input string, limits workerLimits, tool string, args ...string) ([]byte, []byte, error) {
	if runtime.GOOS != "darwin" {
		return nil, nil, errors.New("private preview worker requires macOS")
	}
	if limits.Timeout <= 0 || limits.Output <= 0 || limits.RSS <= 0 {
		return nil, nil, errors.New("invalid worker limits")
	}
	fi, err := os.Lstat(input)
	if err != nil || !fi.Mode().IsRegular() {
		return nil, nil, errors.New("worker input must be a regular non-symlink file")
	}
	// macOS /var is an OS symlink to /private/var; seatbelt matches resolved paths.
	resolvedInput, err := filepath.EvalSymlinks(input)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, limits.Timeout)
	defer cancel()
	reason := make(chan string, 1)
	abort := func(s string) {
		select {
		case reason <- s:
		default:
		}
		cancel()
	}
	argv := append([]string{"-p", sandboxPolicy(resolvedInput, limits.AllowFork), tool}, args...)
	cmd := exec.CommandContext(ctx, "/usr/bin/sandbox-exec", argv...)
	// PDF workers must not inherit transport credentials, SSH agent sockets, or
	// other host secrets even when a parser is compromised.
	cmd.Env = []string{"PATH=/usr/bin:/bin:/opt/homebrew/bin", "HOME=/nonexistent", "LANG=C", "LC_ALL=C", "TMPDIR=/nonexistent"}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = time.Second
	cmd.Cancel = func() error {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		err := cmd.Process.Kill()
		if errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	out := &cappedWriter{max: limits.Output, abort: abort}
	errout := &cappedWriter{max: 64 << 10, abort: abort}
	cmd.Stdout, cmd.Stderr = out, errout
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	monitorExited := make(chan struct{})
	go func() {
		defer close(monitorExited)
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				pctx, stop := context.WithTimeout(ctx, time.Second)
				raw, err := exec.CommandContext(pctx, "/bin/ps", "-o", "rss=", "-p", strconv.Itoa(cmd.Process.Pid), "-g", strconv.Itoa(cmd.Process.Pid)).Output()
				stop()
				if err == nil {
					var rss int64
					for _, field := range strings.Fields(string(raw)) {
						n, _ := strconv.ParseInt(field, 10, 64)
						rss += n * 1024
					}
					if rss > limits.RSS {
						abort("worker sampled RSS budget exceeded")
						return
					}
				}
			}
		}
	}()
	err = cmd.Wait()
	close(done)
	<-monitorExited
	// Descendants cannot survive a successful parent either.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	select {
	case why := <-reason:
		err = errors.New(why)
	default:
		if ctx.Err() != nil {
			err = ctx.Err()
		}
	}
	if err != nil {
		return out.buf.Bytes(), errout.buf.Bytes(), fmt.Errorf("%s: %w", tool, err)
	}
	return out.buf.Bytes(), errout.buf.Bytes(), nil
}
