package finder

import "github.com/cderici/tracedraw2/internal/common"

/*
Max Inner Count Loop Finder is a LoopPicker that finds the
loop that has the largest inner loop use count.

Basically, it finds (outer) loop that has the
most used optimized (inner) loop.

*/

// maxInnerCountLoop implements LoopPicker
type maxInnerCountLoop struct{}

func NewMaxInnerCountLoopPicker() LoopPicker {
	return maxInnerCountLoop{}
}

func (m maxInnerCountLoop) Pick(loops common.LoopMap) common.Loop {
	maxCount := 0
	var maxCountLoop common.Loop

	for _, loop := range loops {
		if loop.InnerUseCount > maxCount {
			maxCount = loop.InnerUseCount
			maxCountLoop = loop
		}
	}

	return maxCountLoop
}
