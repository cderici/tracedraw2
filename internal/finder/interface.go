package finder

import "github.com/cderici/tracedraw2/internal/common"

// LoopPredicate is used for picking a loop that satisfies
// a certain condition
type LoopPredicate func(common.Loop) bool

// LoopPicker represents a loop picker.
// e.g. "get me the Loop that satisfies the given predicate
// among these Loops."
type LoopPicker interface {
	Pick(common.LoopMap) common.Loop
}

// BridgePicker represents a bridge picker.
// e.g. get me the bridge that this guard jumps to.
type BridgePicker interface {
	Pick([]common.Bridge) common.Bridge
}
