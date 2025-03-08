package common

// Trace holds all the sections of a trace, structured and typed.
type Trace struct {
}

// TraceRaw contains all the sections as raw strings
type TraceRaw struct {
	JitSummaryRaw       string
	JitBackendCountsRaw string
	BridgesRaw          []string
	LoopsRaw            []string
}
