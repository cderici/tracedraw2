package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cderici/tracedraw2/internal/reader"
)

func DoPrintRaw(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: tracedraw2-cli print-raw <tracefile>")
		os.Exit(1)
	}

	tracefilePath := args[0]
	// Open trace file
	file, err := os.Open(tracefilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var traceReader reader.TraceReader //nolint
	traceReader = reader.NewFileReader()

	traceRaw, err := traceReader.IngestRaw(scanner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err processing trace: %v", err)
		os.Exit(1)
	}

	fmt.Println(traceRaw.JitSummaryRaw)

	fmt.Println(traceRaw.JitBackendCountsRaw)

	fmt.Printf("\nLoop count: %d\nBridge count: %d\n", len(traceRaw.JitLoopsRaw), len(traceRaw.JitBridgesRaw))
}
