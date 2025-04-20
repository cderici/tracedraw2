package finder

import "github.com/cderici/tracedraw2/internal/common"

/*
Max Preamble Count Loop Finder is a LoopPicker that finds the
loop that has the largest PreambleUseCount.
*/

// maxPreambleCountLoop implements LoopPicker
type maxPreambleCountLoop struct{}

func NewMaxPreambleCountLoopPicker() LoopPicker {
	return maxPreambleCountLoop{}
}

func (_ maxPreambleCountLoop) Pick(loops common.LoopMap) common.Loop {
	maxCount := 0
	var maxCountLoop common.Loop

	for _, loop := range loops {
		if loop.PreambleUseCount > maxCount {
			maxCount = loop.PreambleUseCount
			maxCountLoop = loop
		}
	}

	return maxCountLoop
}
