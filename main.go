package main

import (
	"fmt"
	"os"

	"github.com/stackdiff/stackdiff/internal/parser"
	"github.com/stackdiff/stackdiff/internal/report"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: stackdiff <baseline.tfstate> <current.tfstate>\n")
		os.Exit(1)
	}

	baselinePath := os.Args[1]
	currentPath := os.Args[2]

	baseline, err := parser.ParseStateFile(baselinePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading baseline state: %v\n", err)
		os.Exit(1)
	}

	current, err := parser.ParseStateFile(currentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading current state: %v\n", err)
		os.Exit(1)
	}

	drift := report.Compare(baseline, current)
	report.Print(drift, os.Stdout)
}
