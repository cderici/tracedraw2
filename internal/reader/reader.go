package reader

import (
	"bufio"
	"errors"
	"strings"

	"github.com/cderici/tracedrawer/internal/common"
)

const JITSUMMARY_SECTION = "jit-summary"

func IngestTraceFile(filename string) common.Trace {
	return common.Trace{}
}

func IngestRaw(scanner bufio.Scanner) (common.TraceRaw, error) {
	sections := make(map[string]*strings.Builder)

	var currentSection *strings.Builder
	var currentSectionName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// check if line starts a new section
		if strings.HasPrefix(line, "{") {
			currentSectionName = strings.TrimPrefix(line, "{")
			if currentSectionName != "" {
				sections[currentSectionName] = &strings.Builder{}
				currentSection = sections[currentSectionName]
			}
			continue
		}

		// check if line ends the current section
		if line == currentSectionName+"}" {
			currentSection = nil
			currentSectionName = ""
			continue
		}

		// if inside a secton, append to the corresponding builder
		if currentSection != nil {
			currentSection.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return common.TraceRaw{}, errors.New("read trace")
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
	if _, exists := sections[sectionName]; !exists {
		return "", errors.New("section not there")
	}
	return sectionContents, nil
}
