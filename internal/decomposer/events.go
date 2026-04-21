package decomposer

// DecompositionEvent is the sum type of progress events emitted by the
// two-pass decompose pipeline. Callers type-switch and handle variants
// they care about; unknown variants are safe to ignore.
type DecompositionEvent interface {
	isDecompositionEvent()
}

// Phase names for EventPhaseStarted.
const (
	PhaseContract   = "contract"
	PhaseExtraction = "extraction"
)

// EventPhaseStarted is emitted at the start of each pipeline phase.
// Phase is one of PhaseContract or PhaseExtraction.
type EventPhaseStarted struct {
	Phase string
}

// EventContractChunk is emitted for every assistant-content delta during
// the contract (pass-1) phase. Concatenating Text across chunks in arrival
// order yields the full contract document.
type EventContractChunk struct {
	Text string
}

// EventContractComplete is emitted once pass 1 has returned. Full is the
// complete contract text (same as the concatenation of EventContractChunk
// deltas).
type EventContractComplete struct {
	Full      string
	ElapsedMs int64
}

// EventExtractionProgress is emitted periodically during pass 2 as the
// decompose tool-call arguments stream in. Bytes is the running total of
// argument bytes received so far.
type EventExtractionProgress struct {
	Bytes int
}

// EventRunesComplete is emitted once pass 2 has returned and the rune
// tree has been parsed and normalized.
type EventRunesComplete struct {
	Response  *DecompositionResponse
	ElapsedMs int64
}

// EventReadExample is emitted when the model called read_example during
// the extraction phase.
type EventReadExample struct {
	Paths []string
	Found []string
}

// EventError is emitted when a phase fails. Non-terminal — EventDone
// still follows and the channel still closes.
type EventError struct {
	Phase string
	Err   string
}

// EventCancelled is emitted when ctx is cancelled mid-run. EventDone
// still follows and then the channel closes.
type EventCancelled struct{}

// EventDone is the terminal event. The channel closes immediately after
// it is received.
type EventDone struct {
	ElapsedMs int64
}

func (EventPhaseStarted) isDecompositionEvent()       {}
func (EventContractChunk) isDecompositionEvent()      {}
func (EventContractComplete) isDecompositionEvent()   {}
func (EventExtractionProgress) isDecompositionEvent() {}
func (EventRunesComplete) isDecompositionEvent()      {}
func (EventReadExample) isDecompositionEvent()        {}
func (EventError) isDecompositionEvent()              {}
func (EventCancelled) isDecompositionEvent()          {}
func (EventDone) isDecompositionEvent()               {}
