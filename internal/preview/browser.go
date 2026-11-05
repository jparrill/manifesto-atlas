package preview

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// CheckBrowser uses real CDP device metrics, rather than relying on Chrome's
// minimum window width (which otherwise crops a nominal 390px CLI screenshot).
// It authors no presentation or JavaScript.
func CheckBrowser(ctx context.Context, output string) error {
	if err := Render(output); err != nil {
		return err
	}
	chrome := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	version, err := exec.CommandContext(ctx, chrome, "--version").Output()
	if err != nil {
		return err
	}
	dir, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	target := (&url.URL{Scheme: "file", Path: filepath.Join(dir, "index.html")}).String()
	evidence := map[string]any{"browser": strings.TrimSpace(string(version)), "private_preview": true, "viewport_control": "CDP Emulation.setDeviceMetricsOverride"}
	scripts := regexp.MustCompile(`(?s)<script\b[^>]*>.*?</script>`)
	for _, viewport := range []struct {
		name          string
		width, height int
	}{{"desktop", 1440, 1000}, {"mobile", 390, 844}} {
		profile, err := os.MkdirTemp("", "atlas-preview-browser-")
		if err != nil {
			return err
		}
		dom, png, err := browserView(ctx, chrome, profile, target, viewport.width, viewport.height, viewport.name == "mobile")
		_ = os.RemoveAll(profile)
		if err != nil {
			return fmt.Errorf("%s browser: %w", viewport.name, err)
		}
		dom = scripts.ReplaceAllString(dom, "")
		cards := strings.Count(dom, `class="card"`)
		if cards != 21 || !strings.Contains(dom, "Borrador · Vista previa") {
			return fmt.Errorf("%s browser: expected actual 21-topic draft view, got %d cards", viewport.name, cards)
		}
		if err := writeOutput(filepath.Join(dir, "browser-"+viewport.name+".png"), png); err != nil {
			return err
		}
		evidence[viewport.name] = map[string]any{"width": viewport.width, "height": viewport.height, "rendered_cards": cards, "horizontal_overflow": false, "screenshot_sha256": hashBytes(png)}
	}
	page, err := regularRead(filepath.Join(dir, "index.html"), 1<<20)
	if err != nil {
		return err
	}
	evidence["html_sha256"] = hashBytes(page)
	return saveJSON(filepath.Join(dir, "browser.json"), evidence)
}

func browserView(ctx context.Context, chrome, profile, target string, width, height int, mobile bool) (string, []byte, error) {
	job, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	toBrowser, input, err := os.Pipe()
	if err != nil {
		return "", nil, err
	}
	defer toBrowser.Close()
	defer input.Close()
	output, fromBrowser, err := os.Pipe()
	if err != nil {
		return "", nil, err
	}
	defer output.Close()
	defer fromBrowser.Close()
	deadline, _ := job.Deadline()
	_ = output.SetReadDeadline(deadline)
	cmd := exec.CommandContext(job, chrome, "--headless", "--disable-gpu", "--no-first-run", "--disable-background-networking", "--disable-component-update", "--disable-sync", "--remote-debugging-pipe", "--user-data-dir="+profile)
	cmd.ExtraFiles = []*os.File{toBrowser, fromBrowser}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = time.Second
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	cmd.Stdout = &cappedWriter{max: 1 << 20, abort: func(string) { cancel() }}
	cmd.Stderr = &cappedWriter{max: 64 << 10, abort: func(string) { cancel() }}
	if err := cmd.Start(); err != nil {
		return "", nil, err
	}
	_ = toBrowser.Close()
	_ = fromBrowser.Close()
	defer func() {
		cancel()
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = cmd.Wait()
	}()
	pipe := &browserPipe{in: input, out: bufio.NewReader(output), events: map[string]bool{}}
	dom, image, err := pipe.page(target, width, height, mobile)
	if err != nil {
		return "", nil, err
	}
	png, err := base64.StdEncoding.DecodeString(image)
	return dom, png, err
}
