package decomposer

// ExpansionEvent is the sum type of progress events emitted by
// ExpandStreaming. Callers should type-switch and handle variants they
// care about; unknown variants are safe to ignore.
type ExpansionEvent interface {
	isExpansionEvent()
}

// EventLevelStarted is emitted once per BFS level, before any rune in
// that level is dispatched.
type EventLevelStarted struct {
	Depth int
	Count int
}

// EventRuneStarted is emitted as a rune expansion is dispatched to the
// underlying request. Status → InFlight.
type EventRuneStarted struct {
	Path  string
	Depth int
}

// EventRuneExpanded is emitted when a single rune expansion returns
// successfully. Response may be non-nil even when ChildCount == 0 (leaf).
// ParentPath is the Path of the AutoDecomposition that produced this rune's
// expansion context — "root" for top-level expansions.
type EventRuneExpanded struct {
	Path       string
	ParentPath string
	Depth      int
	Response   *DecompositionResponse
	ElapsedMs  int64
	ChildCount int
}

// EventRuneError is emitted when a single rune expansion fails. Non-fatal;
// the BFS continues with siblings at the same level.
type EventRuneError struct {
	Path      string
	Depth     int
	Err       string
	ElapsedMs int64
}

// EventLevelCompleted is emitted when every expansion in a BFS level has
// finished (success or error). Carries wall-clock vs request-sum timings
// for the CLI parallelism-factor readout.
type EventLevelCompleted struct {
	Depth        int
	WallClockMs  int64
	SumRequestMs int64
}

// EventReadExample is emitted when the model called read_example during
// any decompose attempt. Used by the CLI for a progress line; TUI ignores.
type EventReadExample struct {
	Paths []string
	Found []string
}

// EventCapReached is emitted when the total rune count crosses RuneCap
// before the BFS was otherwise done. No more levels will be dispatched;
// EventDone follows immediately.
type EventCapReached struct {
	TotalRunes int
	Cap        int
}

// EventCancelled is emitted when ctx is cancelled mid-expansion. EventDone
// still follows, then the channel closes.
type EventCancelled struct{}

// EventDone is the terminal event. The channel closes immediately after it
// is received.
type EventDone struct {
	TotalDecompositions int
	TotalRunes          int
	MaxDepth            int
}

func (EventLevelStarted) isExpansionEvent()   {}
func (EventRuneStarted) isExpansionEvent()    {}
func (EventRuneExpanded) isExpansionEvent()   {}
func (EventRuneError) isExpansionEvent()      {}
func (EventLevelCompleted) isExpansionEvent() {}
func (EventReadExample) isExpansionEvent()    {}
func (EventCapReached) isExpansionEvent()     {}
func (EventCancelled) isExpansionEvent()      {}
func (EventDone) isExpansionEvent()           {}
