package reader

import (
	"bufio"
	"regexp"
	"strings"
)

type BackendCounts map[string]int

func processBackendCounts(countsRaw string) BackendCounts {
	strReader := strings.NewReader(countsRaw)
	scanner := bufio.NewScanner(strReader)

	re := regexp.MustCompile(`([a-zA-Z0-9() ]+):([0-9]+)`)
	// WARN: do I want these to be different types?

	for scanner.Scan() {
		line := scanner.Text()

	}
}
