package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cderici/tracedraw2/internal/finder"
	"github.com/cderici/tracedraw2/internal/reader"
)

const MOST_USED_USAGE = `
most-used command is used to find the most used loop
in a given trace log using backend-counts.

It has two options:
  - preamble: finds the loop that has the largest
backend count for its preamble
  - inner (default): finds the loop that has the
largest backend count for its inner optimized loop.

The command prints out the found loop on the stdin.
`

func DoMostUsed(args []string) {
	if len(args) < 1 {
		fmt.Print(MOST_USED_USAGE)
		os.Exit(1)
	}
	chosenOption := "inner"

	var traceFile string
	if len(args) == 2 {
		chosenOption = args[0]
		traceFile = args[1]
	} else {
		traceFile = args[0]
	}

	inFile, err := os.Open(traceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening file %q: %v", traceFile, err)
		os.Exit(1)
	}

	traceReader := reader.NewFileReader()

	trace, err := traceReader.Ingest(bufio.NewScanner(inFile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "err processing trace %q: %v", traceFile, err)
		os.Exit(1)
	}

	var loopPicker finder.LoopPicker

	switch chosenOption {
	case "preamble":
		loopPicker = finder.NewMaxPreambleCountLoopPicker()
	case "inner":
		loopPicker = finder.NewMaxInnerCountLoopPicker()
	default:
		fmt.Fprintf(os.Stderr, "unrecognized option for most-used %q", chosenOption)
		os.Exit(1)
	}

	pickedLoop := loopPicker.Pick(trace.Loops)

	// FIXME: make a proper outputter here
	fmt.Println(pickedLoop)
}
