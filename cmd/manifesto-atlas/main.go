// Command manifesto-atlas is the repository entry point. It dispatches to the
// cli package and propagates the computed exit code to the process.
package main

import (
	"os"

	"github.com/jparrill/manifesto-atlas/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
