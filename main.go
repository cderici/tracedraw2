package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cderici/tracedraw2/internal/reader"
)

func main() {
	// Open trace file
	file, err := os.Open("demo.trace")
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

	fmt.Println(traceRaw.Jit_Summary_Raw)

	fmt.Println(traceRaw.Jit_Backend_Counts_Raw)

	fmt.Printf("\nLoop count: %d\nBridge count: %d\n", len(traceRaw.Jit_Loops_Raw), len(traceRaw.Jit_Bridges_Raw))
}
