# Requirement: "a program execution tracer that records step-by-step state for visualization"

Runs a program (via an injected executor) and produces an ordered trace of frames, variables, and heap state suitable for rendering as an execution diagram.

std
  std.json
    std.json.encode_object
      fn (obj: map[string, dynamic_value]) -> string
      + encodes a dynamic map as JSON
      # serialization

tracer
  tracer.new_session
    fn (source: string) -> result[trace_session, string]
    + creates a session ready to execute the given source
    - returns error when the source is empty
    # construction
  tracer.step
    fn (session: trace_session) -> result[trace_frame, string]
    + advances one executable step and returns the captured frame
    - returns error "finished" when there are no more steps
    # stepping
  tracer.record_full
    fn (session: trace_session, max_steps: i32) -> result[list[trace_frame], string]
    + records up to max_steps frames and returns them in order
    - returns error when execution raises an unrecoverable fault
    # recording
  tracer.capture_locals
    fn (frame: trace_frame) -> map[string, dynamic_value]
    + returns the local variables visible at the given frame
    # inspection
  tracer.capture_heap
    fn (frame: trace_frame) -> map[string, dynamic_value]
    + returns heap objects keyed by address referenced from the frame
    # inspection
  tracer.export
    fn (frames: list[trace_frame]) -> string
    + returns a JSON document describing the full trace
    # export
    -> std.json.encode_object
