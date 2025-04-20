package reader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

var SECTION_STARTER = regexp.MustCompile(`(?i)^\[([0-9a-f]+)\]\s*\{([a-z-]+)$`)
var SECTION_FINISHER = regexp.MustCompile(`(?i)^\[([0-9a-f]+)\]\s*([a-z-]+)\}$`)

func (f *fileReader) Ingest(scanner *bufio.Scanner) (common.Trace, error) {
	// We can use IngestRaw to get the TraceRaw
	rawTrace, err := f.IngestRaw(scanner)
	if err != nil {
		return common.Trace{}, err
	}

	// BackendCounts need to be processed before creating Loop objects.
	var loopBackendCounts map[common.LoopID]int
	loopBackendCounts = getBackendCounts(rawTrace.JitBackendCountsRaw)

	return common.Trace{}, nil
}

func getBackendCounts(rawCounts string) map[common.LoopID]int {
	counts := make(map[common.LoopID]int)

	re := regexp.MustCompile(`TargetToken\((\d+)\):(\d+)`)
	matches := re.FindAllStringSubmatch(rawCounts, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		loopID := common.LoopID(match[1])
		count, err := strconv.Atoi(match[2])
		if err != nil {
			fmt.Fprint(os.Stderr, "Warning: weird backend count detected %q", count)
			continue
		}
		counts[loopID] = count
	}

	return counts
}

func (f *fileReader) IngestRaw(scanner *bufio.Scanner) (common.TraceRaw, error) {

	var SECTIONS = map[string]*strings.Builder{
		JIT_LOOP_SECTION:           {},
		JIT_BRIDGE_SECTION:         {},
		JIT_SUMMARY_SECTION:        {},
		JIT_BACKEND_COUNTS_SECTION: {},
	}

	raw_trace := common.TraceRaw{}

	var currentSectionBuilder *strings.Builder
	var currentSectionName string

	// ignoreSection is  a flag we activate whenever we see a section we want to ignore
	ignoreSectionName := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if ignoreSectionName == "" && SECTION_STARTER.MatchString(line) {
			// A section is starting
			sub := SECTION_STARTER.FindStringSubmatch(line)
			sectionName := sub[2]
			// Is this a section that I care about?
			if builder, found := SECTIONS[sectionName]; found {
				currentSectionName = sectionName
				currentSectionBuilder = builder
				// Also write the current line
				currentSectionBuilder.WriteString(line + "\n")
			} else {
				ignoreSectionName = sectionName
			}
		} else if SECTION_FINISHER.MatchString(line) {
			// A section is ending
			sub := SECTION_FINISHER.FindStringSubmatch(line)
			sectionName := sub[2]
			if _, found := SECTIONS[sectionName]; found {
				// Write the final line as well
				currentSectionBuilder.WriteString(line + "\n")

				// Build the section string
				str := currentSectionBuilder.String()

				switch currentSectionName {
				case JIT_SUMMARY_SECTION:
					raw_trace.JitSummaryRaw = str
				case JIT_BACKEND_COUNTS_SECTION:
					raw_trace.JitBackendCountsRaw = str
				case JIT_BRIDGE_SECTION:
					raw_trace.JitBridgesRaw = append(raw_trace.JitBridgesRaw, str)
				case JIT_LOOP_SECTION:
					raw_trace.JitLoopsRaw = append(raw_trace.JitLoopsRaw, str)
				}

				currentSectionBuilder = nil
				currentSectionName = ""

			} else {
				// We're ending a section that we were ignoring
				ignoreSectionName = ""
			}
		} else {
			if ignoreSectionName == "" {
				// Assert that we have a builder for the current section we're scanning
				if currentSectionBuilder == nil {
					return common.TraceRaw{}, je.Annotatef(common.ErrScanningTrace, "section: %q", currentSectionName)
				}
				// Write line to the builder for the current section
				currentSectionBuilder.WriteString(line + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return common.TraceRaw{}, je.Annotate(err, "scanning trace file")
	}

	return raw_trace, nil
}
