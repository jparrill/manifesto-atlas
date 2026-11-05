// Command preview runs the explicitly private macOS real-source demonstration.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jparrill/manifesto-atlas/internal/preview"
)

func main() {
	var cfg preview.Config
	var reference string
	var render bool
	var browser bool
	flag.StringVar(&cfg.Source, "source", "fixtures/frente/source.json", "retained candidate descriptor")
	flag.StringVar(&cfg.Output, "output", "runs/preview/frente-e2e", "private run directory; new for evaluation, existing for offline rendering")
	flag.StringVar(&cfg.Endpoint, "endpoint", os.Getenv("ATLAS_BASE_URL"), "operator-configured OpenAI-compatible endpoint")
	flag.StringVar(&cfg.Model, "model", os.Getenv("ATLAS_MODEL"), "operator-configured model")
	flag.StringVar(&reference, "import-reference", "", "import owner-supplied existing HTML unchanged")
	flag.BoolVar(&render, "render", false, "rerender retained evidence offline, without model calls")
	flag.BoolVar(&browser, "check-browser", false, "capture actual desktop/mobile Chrome views of the retained preview")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected positional arguments")
		os.Exit(2)
	}
	cfg.APIKey = os.Getenv("ATLAS_API_KEY")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var err error
	switch {
	case reference != "":
		err = preview.ImportReference(reference)
	case render:
		err = preview.Render(cfg.Output)
	case browser:
		err = preview.CheckBrowser(ctx, cfg.Output)
	default:
		err = preview.Run(ctx, cfg)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "preview:", err)
		os.Exit(1)
	}
}
