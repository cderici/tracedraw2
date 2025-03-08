package reader

import (
	"bufio"
	"strings"

	"github.com/cderici/tracedraw2/internal/common"
	je "github.com/juju/errors"
)

// fileReader implements TraceReader
type fileReader struct {
}

// NewFileReader returns a fileReader that implements TraceReader
func NewFileReader() TraceReader {
	return &fileReader{}
}

const JITLOOP_SECTION = "jit-log-opt-loop"
const JITBRIDGE_SECTION = "jit-log-opt-bridge"
const JITBACKEND_COUNTS_SECTION = "jit-backend-counts"
const JITSUMMARY_SECTION = "jit-summary"

var SECTIONS = map[string]*strings.Builder{
	JITLOOP_SECTION:           {},
	JITBRIDGE_SECTION:         {},
	JITBACKEND_COUNTS_SECTION: {},
	JITSUMMARY_SECTION:        {},
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
				case JITSUMMARY_SECTION:
					jitSummaryRaw = str
				case JITBACKEND_COUNTS_SECTION:
					jitBackendCounts = str
				case JITBRIDGE_SECTION:
					jitBridges = append(jitBridges, str)
				case JITLOOP_SECTION:
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

	raw_trace.JitSummaryRaw = jitSummaryRaw
	raw_trace.JitBackendCountsRaw = jitBackendCounts
	raw_trace.BridgesRaw = jitBridges
	raw_trace.LoopsRaw = jitLoops

	return raw_trace, nil
}
