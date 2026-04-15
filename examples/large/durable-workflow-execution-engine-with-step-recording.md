# Requirement: "a durable execution engine that records workflow steps so they survive restarts"

Workflows are a sequence of named steps. Each step is persisted to an event history; replay reconstructs state without re-running completed steps.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      fn () -> string
      + returns a fresh random identifier
      # identifiers
  std.encoding
    std.encoding.encode_event
      fn (kind: string, payload: bytes) -> bytes
      + encodes an event kind and payload into bytes
      # serialization
    std.encoding.decode_event
      fn (raw: bytes) -> result[tuple[string, bytes], string]
      + decodes bytes back into (kind, payload)
      - returns error on truncated input
      # serialization

durable
  durable.new_engine
    fn () -> engine_state
    + creates an engine with no active workflows
    # construction
  durable.start_workflow
    fn (state: engine_state, name: string, input: bytes) -> tuple[string, engine_state]
    + starts a new workflow and returns its id
    + records a WORKFLOW_STARTED event in the history
    # workflow_lifecycle
    -> std.id.new_uuid
    -> std.time.now_millis
    -> std.encoding.encode_event
  durable.schedule_step
    fn (state: engine_state, workflow_id: string, step_name: string, input: bytes) -> result[engine_state, string]
    + records a STEP_SCHEDULED event for the next step
    - returns error when the workflow id is unknown
    - returns error when the workflow is already complete
    # step_scheduling
    -> std.encoding.encode_event
  durable.complete_step
    fn (state: engine_state, workflow_id: string, step_name: string, output: bytes) -> result[engine_state, string]
    + records a STEP_COMPLETED event and stores the step output
    - returns error when no matching STEP_SCHEDULED event exists
    # step_completion
    -> std.encoding.encode_event
    -> std.time.now_millis
  durable.fail_step
    fn (state: engine_state, workflow_id: string, step_name: string, reason: string) -> result[engine_state, string]
    + records a STEP_FAILED event with the reason
    - returns error when no matching STEP_SCHEDULED event exists
    # step_completion
    -> std.encoding.encode_event
  durable.get_step_output
    fn (state: engine_state, workflow_id: string, step_name: string) -> optional[bytes]
    + returns the stored output for a completed step
    ? enables idempotent re-execution: callers check this before re-running
    # step_query
  durable.complete_workflow
    fn (state: engine_state, workflow_id: string, result: bytes) -> result[engine_state, string]
    + records a WORKFLOW_COMPLETED event and marks the workflow terminal
    - returns error when the workflow id is unknown
    # workflow_lifecycle
    -> std.encoding.encode_event
  durable.history
    fn (state: engine_state, workflow_id: string) -> result[list[bytes], string]
    + returns the ordered event bytes for the workflow
    - returns error when the workflow id is unknown
    # persistence
  durable.replay
    fn (history: list[bytes]) -> result[engine_state, string]
    + rebuilds engine state from a recorded history
    - returns error when any event is corrupt
    # recovery
    -> std.encoding.decode_event
  durable.pending_steps
    fn (state: engine_state, workflow_id: string) -> result[list[string], string]
    + returns the names of scheduled-but-not-completed steps
    - returns error when the workflow id is unknown
    # introspection
