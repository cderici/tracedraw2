package main

import (
	"fmt"
	"os"

	"github.com/cderici/tracedraw2/cmd/tracedraw2-cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tracedraw2-cli print-raw <tracefile>")
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "print-raw":
		commands.DoPrintRaw(args)
	case "trace-summaries":
		commands.DoTraceSummaries(args)
	default:
		fmt.Printf("Unknown command: %q\n", cmd)
		os.Exit(1)
	}
}
