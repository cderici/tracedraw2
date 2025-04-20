package main

import (
	"fmt"
	"os"

	"github.com/cderici/tracedraw2/cmd/tracedraw2-cli/commands"
)

const USAGE = `
USAGE: tracedraw2-cli <command> <args>

commands:
	print-raw <traceFile>
	trace-summaries <directory>
	most-used [preamble|inner*] <traceFile>
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(USAGE)
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "print-raw":
		commands.DoPrintRaw(args)
	case "trace-summaries":
		commands.DoTraceSummaries(args)
	case "most-used":
		commands.DoMostUsed(args)
	default:
		fmt.Printf("Unknown command: %q\n", cmd)
		os.Exit(1)
	}
}
