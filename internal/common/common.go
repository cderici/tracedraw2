package common

// Trace holds all relevant the sections of a trace,
// structured and typed.
type Trace struct {
	Loops   LoopMap      // jit-log-opt-loop
	Bridges BridgeMap    // jit-log-opt-bridge
	Summary TraceSummary // jit-summary
}

// Loop holds all relevant sections of a traced loop
type Loop struct {
	ID      LoopID
	OpCount int    // number of operations
	TopCode string // interpreted code

	PreambleLabel    string     // descr=TargetToken(....)
	PreambleArgs     []Argument // TODO: these might be further typed
	PreambleCode     []Line     // TODO: these will be typed further
	PreambleUseCount int        // how many times preamble is used

	InnerLabel    string
	InnerArgs     []Argument
	InnerCode     []Line
	InnerUseCount int // how many times optimized loop is used

	JumpArgs   []Argument // args passed when jumped
	JumpTarget Jumpable   // loop or bridge we're jumping to
}

type Bridge struct {
	ID         BridgeID // FIXME: what's this ID?
	OutOfGuard Guard    // guard that leads us here
	Args       []Argument
	Code       []Line

	JumpArgs   []Argument
	JumpTarget Jumpable
}

// TODO: Generic types that might be elaborated further
type LoopID string
type BridgeID string
type Argument string
type Line string
type Guard string // just the Guard ID for now, e.g. descr=<Guard0x....>
type LoopMap map[LoopID]Loop
type BridgeMap map[BridgeID]Bridge

// Jumpable represents Loops and Bridges
type Jumpable interface {
	isJumpable()
}

func (*Loop) isJumpable()   {}
func (*Bridge) isJumpable() {}

// TraceRaw contains all the sections as raw strings
type TraceRaw struct {
	JitLoopsRaw         []string `json:"jit_loops_raw"`
	JitBridgesRaw       []string `json:"jit_bridges_raw"`
	JitSummaryRaw       string   `json:"jit_summary_raw"`
	JitBackendCountsRaw string   `json:"jit_backend_counts_raw"`
}

// TraceSummary represents the structured data from a JIT summary trace.
type TraceSummary struct {
	Tracing              CountDuration `json:"tracing"`                 // Number of traces recorded and total time spent recording them.
	Backend              CountDuration `json:"backend"`                 // Number of traces compiled and total time spent compiling them.
	TotalTime            float64       `json:"total_time"`              // Total time spent by the JIT (in seconds).
	Ops                  int           `json:"ops"`                     // Total number of operations emitted.
	HeapCachedOps        int           `json:"heap_cached_ops"`         // Operations benefitting from heap caching.
	RecordedOps          int           `json:"recorded_ops"`            // Operations recorded in traces.
	Calls                int           `json:"calls"`                   // Number of function calls recorded in traces.
	Guards               int           `json:"guards"`                  // Total number of guards.
	OptOps               int           `json:"opt_ops"`                 // Operations remaining after optimizations.
	OptGuards            int           `json:"opt_guards"`              // Guards remaining after optimizations (reduced from total guards).
	OptGuardsShared      int           `json:"opt_guards_shared"`       // Guards that were shared or merged.
	Forcings             int           `json:"forcings"`                // Number of virtuals forced (ideally low).
	AbortTraceTooLong    int           `json:"abort_trace_too_long"`    // Number of aborts due to trace being too long.
	AbortCompiling       int           `json:"abort_compiling"`         // Number of aborts due to compilation failure.
	AbortVableEscape     int           `json:"abort_vable_escape"`      // Number of aborts due to an escaped virtualizable.
	AbortBadLoop         int           `json:"abort_bad_loop"`          // Number of aborts due to a "bad loop" (unclear definition).
	AbortForceQuasiImmut int           `json:"abort_force_quasi_immut"` // Number of aborts due to forced quasi-immutables.
	AbortSegmentingTrace int           `json:"abort_segmenting_trace"`  // Number of aborts due to segmenting constraints.
	VirtualizablesForced int           `json:"virtualizables_forced"`   // Number of virtualizables forced (possibly different from "forcings").
	NVirtuals            int           `json:"nvirtuals"`               // Number of virtuals (higher is better?).
	NVholes              int           `json:"nvholes"`                 // Number of virtualization holes (unclear definition).
	NVReused             int           `json:"nvreused"`                // Number of virtuals successfully reused (higher is better).
	VecOptTried          int           `json:"vecopt_tried"`            // Number of vectorization attempts.
	VecOptSuccess        int           `json:"vecopt_success"`          // Number of successful vectorized operations.
	TotalLoops           int           `json:"total_loops"`             // Total number of loops captured.
	TotalBridges         int           `json:"total_bridges"`           // Total number of bridges (high values may indicate too many guards?).
	FreedLoops           int           `json:"freed_loops"`             // Number of loops freed (unclear impact).
	FreedBridges         int           `json:"freed_bridges"`           // Number of bridges freed (unclear impact).
}

// CountDuration represents fields containing both a count and a duration.
// Used to indicate "time spent"
type CountDuration struct {
	Count   int     `json:"count"`   // The number of occurrences.
	Seconds float64 `json:"seconds"` // The total time spent in seconds.
}
