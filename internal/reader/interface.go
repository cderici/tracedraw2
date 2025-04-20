package reader

import (
	"bufio"

	"github.com/cderici/tracedraw2/internal/common"
)

// TraceReader interface has methods for reading a trace file
type TraceReader interface {
	// IngestRaw produces raw sections
	IngestRaw(scanner *bufio.Scanner) (common.TraceRaw, error)

	// Ingest produces Trace out of a given trace log file scanner
	Ingest(scanner *bufio.Scanner) (common.Trace, error)
}
