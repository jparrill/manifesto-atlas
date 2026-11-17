// Command runcheck verifies a programme run offline: literal quotes, page
// citations, proposal coverage, curated claims and section codes.
//
//	go run ./cmd/runcheck -run runs/29n2026
//
// It exits 0 when the run passes and 1 with one line per failed check.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jparrill/manifesto-atlas/internal/runcheck"
)

func main() {
	dir := flag.String("run", "runs/29n2026", "run directory holding run.json")
	flag.Parse()

	report, err := runcheck.Verify(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, line := range report.Summary() {
		fmt.Println(line)
	}
	for _, issue := range report.Issues {
		fmt.Println("  FAIL", issue)
	}
	if !report.OK() {
		fmt.Printf("runcheck: %d issues\n", len(report.Issues))
		os.Exit(1)
	}
}
