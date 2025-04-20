package finder

import "github.com/cderici/tracedraw2/internal/common"

// LoopPicker represents a loop picker.
// e.g. what's the loop that has the largest inner loop use count?
type LoopPicker interface {
	Pick([]common.Loop) common.Loop
}

// BridgePicker represents a bridge picker.
// e.g. get me the bridge that this guard jumps to.
type BridgePicker interface {
	Pick([]common.Bridge) common.Bridge
}
