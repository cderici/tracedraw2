package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cderici/tracedraw2/internal/reader"
)

// DoTraceSummaries takes a folder and scans for .trace files
// For each trace file it grabs the summaries and puts the summary into
// a file in a "summaries" directory within the given directory.
func DoTraceSummaries(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: tracedraw2-cli trace-summaries <directory>")
		os.Exit(1)
	}

	dir := args[0]
	fmt.Printf("Generating summaries for all traces in: %q\n\n", dir)

	// Create summaries/ subdir
	summariesDir := filepath.Join(dir, "summaries")
	if err := os.MkdirAll(summariesDir, 0755); err != nil {
		fmt.Printf("Error creating summaries subdir: %v", err)
		os.Exit(1)
	}

	// Iterate over <name>.trace files and write their summaries
	// into <name>.summary files
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory %s : %v", dir, err)
		os.Exit(1)
	}

	traceReader := reader.NewFileReader()

	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".trace" {
			tracePath := filepath.Join(dir, f.Name())

			// Open the trace file
			inFile, err := os.Open(tracePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err opening trace file %q: %v", f.Name(), err)
				os.Exit(1)
			}
			defer inFile.Close()

			// Ingest the trace log contents
			traceRaw, err := traceReader.IngestRaw(bufio.NewScanner(inFile))
			if err != nil {
				fmt.Fprintf(os.Stderr, "err processing trace %q: %v", f.Name(), err)
				os.Exit(1)
			}

			// Extract the summary
			summary := traceRaw.JitSummaryRaw

			// Remove the ".trace" extension from the original filename
			baseName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))

			// write summary into the "summaries/"+f.Name()+".summary"
			outPath := filepath.Join(summariesDir, baseName+".summary")
			outFile, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err opening output file %q: %v", outPath, err)
				os.Exit(1)
			}
			defer outFile.Close()

			// Write it as json maybe?
			_, err = fmt.Fprintln(outFile, summary)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err writing summary for %q: %v", tracePath, err)
				os.Exit(1)
			}

			fmt.Printf("Wrote summary for %s into %s\n", tracePath, outPath)
		}
	}
}
