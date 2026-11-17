// Package cli implements the manifesto-atlas command dispatcher.
//
// The dispatcher is the user-facing CLI. It supports safe commands (help,
// version) and declares the 12 pipeline stages from docs/SPEC.md#make-targets.
// Stage commands are registered now so a later implementation can replace their
// not-available behavior without changing the user-facing command names.
//
// No command performs network retrieval, model invocation, or publication.
// Not-available stages create no artifacts: they validate their arguments and
// then report their owning task. Argument values are treated as opaque data,
// never interpolated or executed.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// Version is the manifesto-atlas version reported by the version command.
const Version = "0.1.0"

// Process exit codes.
const (
	exitOK     = 0
	exitFailed = 1 // command ran but reported a not-available or operational failure
	exitUsage  = 2 // usage error: unknown command, invalid flag, wrong argument count
)

// stage is a declared pipeline stage that is not yet available.
type stage struct {
	name    string
	owner   string
	purpose string
}

// stages mirrors the 12 pipeline stages from docs/SPEC.md#make-targets. The
// order and names must stay in sync with the Makefile and the Spec.
var stages = []stage{
	{"submit", "TQ-0021", "inspect a selected local candidate"},
	{"verify-source", "TQ-0027", "verify local selection evidence and matching bytes"},
	{"add", "TQ-0027", "register a verified immutable source"},
	{"validate-extraction", "TQ-0030", "reject deterministic page/warning blockers"},
	{"normalize", "TQ-0038", "produce draft canonical programme JSON"},
	{"validate-site", "TQ-0049", "run basic single-programme site quality checks"},
	{"publish-check", "TQ-0050 TQ-0049", "verify offline Markdown eligibility first"},
	{"serve", "TQ-0065", "preview an explicit private generated root locally"},
}

// stageCommand is an implemented pipeline stage that takes flags.
type stageCommand struct {
	name       string
	positional int
	flags      func(string) (*flag.FlagSet, *stageFlags)
	run        func(pos []string, f *stageFlags, w io.Writer) int
}

// command is a dispatchable CLI command.
type command struct {
	name       string
	positional int // exact number of positional arguments accepted
	run        func(args []string, w io.Writer) int
	// flags, when set, parses stage flags before runWithFlags is called.
	flags        func(string) (*flag.FlagSet, *stageFlags)
	runWithFlags func(pos []string, f *stageFlags, w io.Writer) int
}

func (c *command) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(c.name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

// usage writes the top-level help text to w.
func usage(w io.Writer) {
	fmt.Fprintln(w, "manifesto-atlas "+Version)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  manifesto-atlas [command] [flags] [arguments]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  help       Show this help.")
	fmt.Fprintln(w, "  version    Show the version and the Go runtime version.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Run stages (implemented, deterministic, offline):")
	fmt.Fprintln(w, "  extract <pdf> -url U -programme ID [-run DIR] [-title NAME] [-coverage numbers|items]")
	fmt.Fprintln(w, "             provenance + page-marked extraction + manifest + summary skeleton.")
	fmt.Fprintln(w, "  summarize <run> -programme ID|-all [-endpoint U] [-model M] [-attempts N]")
	fmt.Fprintln(w, "             write the summary with the configured model and repair it with the gate.")
	fmt.Fprintln(w, "  validate-content <run>")
	fmt.Fprintln(w, "             run the gate and write Verificacion.md next to the run.")
	fmt.Fprintln(w, "  process [<pdf>|auto] -url U -programme ID [flags]")
	fmt.Fprintln(w, "             extract, summarise, verify and render in one go; no human step.")
	fmt.Fprintln(w, "  build <run> [-out FILE] [-template FILE] [-force]")
	fmt.Fprintln(w, "             render the single-file view with the documents embedded.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Pipeline stages (declared; not available until their owning task lands):")
	for _, s := range stages {
		fmt.Fprintf(w, "  %-20s %s  %s\n", s.name, s.owner, s.purpose)
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "A stage requires exactly one target-path argument and reports")
	fmt.Fprintln(w, "not-available until its owning task lands. No command fetches,")
	fmt.Fprintln(w, "invokes a model, or publishes. Supported Go version: go >= 1.26.")
}

// Run dispatches args and returns the process exit code. It writes results to
// stdout and diagnostics to stderr.
func Run(args []string) int {
	return run(args, os.Stdout, os.Stderr)
}

func run(args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		usage(out)
		return exitOK
	}

	name, rest := args[0], args[1:]
	cmd := lookup(name)
	if cmd == nil {
		fmt.Fprintf(errOut, "manifesto-atlas: unknown command %q\n", name)
		fmt.Fprintln(errOut, "Run 'manifesto-atlas help' for usage.")
		return exitUsage
	}

	if cmd.flags != nil {
		fs, f := cmd.flags(cmd.name)
		// Flags may follow the target path ("extract programa.pdf -url …"): take
		// the leading positional first so the standard flag parser sees flags only.
		var pos []string
		if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
			pos = append(pos, rest[0])
			rest = rest[1:]
		}
		if err := fs.Parse(rest); err != nil {
			if err == flag.ErrHelp {
				usage(out)
				return exitOK
			}
			fmt.Fprintf(errOut, "manifesto-atlas %s: %v\n", cmd.name, err)
			fmt.Fprintln(errOut, "Run 'manifesto-atlas help' for usage.")
			return exitUsage
		}
		pos = append(pos, fs.Args()...)
		if cmd.positional == 0 && len(pos) == 1 {
			// A command that can fetch its input still accepts an explicit path.
		} else if len(pos) != cmd.positional {
			fmt.Fprintf(errOut, "manifesto-atlas %s: expected %s\n", cmd.name, wantArgs(cmd.positional))
			fmt.Fprintln(errOut, "Run 'manifesto-atlas help' for usage.")
			return exitUsage
		}
		return cmd.runWithFlags(pos, f, out)
	}

	fs := cmd.flagSet()
	if err := fs.Parse(rest); err != nil {
		if err == flag.ErrHelp {
			usage(out)
			return exitOK
		}
		fmt.Fprintf(errOut, "manifesto-atlas %s: %v\n", cmd.name, err)
		fmt.Fprintln(errOut, "Run 'manifesto-atlas help' for usage.")
		return exitUsage
	}

	pos := fs.Args()
	if len(pos) != cmd.positional {
		fmt.Fprintf(errOut, "manifesto-atlas %s: expected %s\n", cmd.name, wantArgs(cmd.positional))
		fmt.Fprintln(errOut, "Run 'manifesto-atlas help' for usage.")
		return exitUsage
	}
	return cmd.run(pos, out)
}

func wantArgs(n int) string {
	switch n {
	case 0:
		return "no arguments"
	case 1:
		return "exactly one target-path argument"
	default:
		return fmt.Sprintf("exactly %d arguments", n)
	}
}

// processStage adapts the whole-run orchestrator.
func processStage() *command { return adapt(processCommand()) }

// summarizeStage adapts the summarisation stage.
func summarizeStage() *command { return adapt(summarizeCommand()) }

// extractStage, validateStage and buildStage adapt the implemented stages to the
// plain command shape used by the dispatcher.
func extractStage() *command  { return adapt(extractCommand()) }
func validateStage() *command { return adapt(validateContent()) }
func buildStage() *command    { return adapt(buildCommand()) }

func adapt(s stageCommand) *command {
	return &command{name: s.name, positional: s.positional, flags: s.flags, runWithFlags: s.run}
}

// lookup returns the command for name, or nil if it is not registered.
func lookup(name string) *command {
	switch name {
	case "help":
		return &command{name: "help", positional: 0, run: func(_ []string, w io.Writer) int {
			usage(w)
			return exitOK
		}}
	case "version":
		return &command{name: "version", positional: 0, run: func(_ []string, w io.Writer) int {
			fmt.Fprintf(w, "manifesto-atlas %s\n", Version)
			fmt.Fprintf(w, "go %s\n", runtime.Version())
			return exitOK
		}}
	}
	switch name {
	case "extract":
		return extractStage()
	case "summarize":
		return summarizeStage()
	case "process":
		return processStage()
	case "validate-content":
		return validateStage()
	case "build":
		return buildStage()
	}
	for i := range stages {
		if stages[i].name == name {
			s := stages[i]
			return &command{name: s.name, positional: 1, run: func(target []string, w io.Writer) int {
				fmt.Fprintf(w, "manifesto-atlas %s: not available yet.\n", s.name)
				fmt.Fprintf(w, "  Owning task: %s.\n", s.owner)
				fmt.Fprintf(w, "  Target: %s\n", target[0])
				fmt.Fprintf(w, "  This stage fails until %s lands. Run: tq show %s\n", s.owner, s.owner)
				return exitFailed
			}}
		}
	}
	return nil
}
