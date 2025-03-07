package reader

import (
	"bufio"
	"strings"

	"github.com/cderici/tracedrawer/internal/common"
	je "github.com/juju/errors"
)

// fileReader implements TraceReader
type fileReader struct {
}

// NewFileReader returns a fileReader that implements TraceReader
func NewFileReader() TraceReader {
	return &fileReader{}
}

const JITSUMMARY_SECTION = "jit-summary"

func IngestTraceFile(filename string) common.Trace {
	return common.Trace{}
}

func (f *fileReader) IngestRaw(scanner *bufio.Scanner) (common.TraceRaw, error) {
	sections := make(map[string]*strings.Builder)

	var currentSectionBuilder *strings.Builder
	var currentSectionName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// check if line starts a new section
		if strings.Contains(line, "{"+JITSUMMARY_SECTION) {
			currentSectionName = JITSUMMARY_SECTION
			sections[currentSectionName] = &strings.Builder{}
			currentSectionBuilder = sections[currentSectionName]

			continue
		}

		// check if line ends the current section
		if strings.Contains(line, currentSectionName+"}") {
			currentSectionBuilder = nil
			currentSectionName = ""

			continue
		}

		// if inside a secton, append to the corresponding builder
		if currentSectionBuilder != nil {
			currentSectionBuilder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return common.TraceRaw{}, je.Annotate(err, "scanning trace file")
	}

	jitSummary, err := getSectionContents(sections, JITSUMMARY_SECTION)
	if err != nil {
		return common.TraceRaw{}, err
	}

	return common.TraceRaw{
		JitSummaryRaw: jitSummary,
	}, nil
}

func getSectionContents(sections map[string]*strings.Builder, sectionName string) (string, error) {
	var sectionContents string
	var err error

	if contents, exists := sections[sectionName]; exists {
		sectionContents = contents.String()
	} else {
		err = je.Annotatef(common.ErrIngestNoSuchSection, "section %s", sectionName)
	}

	return sectionContents, err
}
