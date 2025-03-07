package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cderici/tracedrawer/internal/reader"
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
	traceRaw, err := reader.IngestRaw(*scanner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err processing trace: %v", err)
		os.Exit(1)
	}

	fmt.Println(traceRaw.JitSummaryRaw)
}
