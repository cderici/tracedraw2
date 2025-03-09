package common

// Trace holds all the sections of a trace, structured and typed.
type Trace struct {
}

// TraceRaw contains all the sections as raw strings
type TraceRaw struct {
	Jit_BackendDumps_Raw   []string
	Jit_Backends_Raw       map[string][]string
	Jit_Loops_Raw          []string
	Jit_Bridges_Raw        []string
	Jit_Summary_Raw        string
	Jit_Backend_Counts_Raw string
}
