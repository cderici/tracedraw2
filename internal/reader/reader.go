package reader

import (
	"bufio"
	"strings"

	"github.com/cderici/tracedraw2/internal/common"
	je "github.com/juju/errors"
)

/*
		Trace File

jit-log-opt-loop
...
...
jit-log-opt-bridge
...
...
jit-summary
jit-backend-counts

*/

// fileReader implements TraceReader
type fileReader struct {
}

// NewFileReader returns a fileReader that implements TraceReader
func NewFileReader() TraceReader {
	return &fileReader{}
}

const JIT_LOOP_SECTION = "jit-log-opt-loop"
const JIT_BRIDGE_SECTION = "jit-log-opt-bridge"
const JIT_SUMMARY_SECTION = "jit-summary"
const JIT_BACKEND_COUNTS_SECTION = "jit-backend-counts"

var SECTIONS = map[string]*strings.Builder{
	JIT_LOOP_SECTION:           {},
	JIT_BRIDGE_SECTION:         {},
	JIT_SUMMARY_SECTION:        {},
	JIT_BACKEND_COUNTS_SECTION: {},
}

func IngestTraceFile(filename string) common.Trace {
	return common.Trace{}
}

func (f *fileReader) IngestRaw(scanner *bufio.Scanner) (common.TraceRaw, error) {
	var jitSummaryRaw, jitBackendCounts string
	var jitBridges, jitLoops []string

	raw_trace := common.TraceRaw{}

	var currentSectionBuilder *strings.Builder
	var currentSectionName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// check if line ends the current section
		if currentSectionBuilder != nil {
			if strings.Contains(line, currentSectionName+"}") {
				str := currentSectionBuilder.String()
				switch currentSectionName {
				case JIT_SUMMARY_SECTION:
					jitSummaryRaw = str
				case JIT_BACKEND_COUNTS_SECTION:
					jitBackendCounts = str
				case JIT_BRIDGE_SECTION:
					jitBridges = append(jitBridges, str)
				case JIT_LOOP_SECTION:
					jitLoops = append(jitLoops, str)
				}

				currentSectionBuilder = nil
				currentSectionName = ""

				continue
			}
		}

		// check if line starts a new section
		if currentSectionBuilder == nil {
			for sectionName, builder := range SECTIONS {
				if strings.Contains(line, "{"+sectionName) {
					currentSectionName = sectionName
					currentSectionBuilder = builder
				}
				break
			}
			if currentSectionBuilder == nil {
				return raw_trace, je.Annotatef(common.ErrDataWithNoSectionInTraceFile, "%s", line)
			}
		}

		// if inside a secton, append to the corresponding builder
		if currentSectionBuilder != nil {
			currentSectionBuilder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return common.TraceRaw{}, je.Annotate(err, "scanning trace file")
	}

	raw_trace.Jit_Summary_Raw = jitSummaryRaw
	raw_trace.Jit_Backend_Counts_Raw = jitBackendCounts
	raw_trace.Jit_Bridges_Raw = jitBridges
	raw_trace.Jit_Loops_Raw = jitLoops

	return raw_trace, nil
}
