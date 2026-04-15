# Requirement: "a stream processing engine that applies filter, map, and windowed aggregation operators to continuous event streams"

A dataflow builder lets callers compose operators into a pipeline; the runtime feeds events in and emits outputs. std handles JSON decoding for event payloads, monotonic time for windowing, and a tiny expression evaluator for filters.

std
  std.time
    std.time.now_monotonic_ms
      fn () -> i64
      + returns a monotonically non-decreasing millisecond counter
      # time
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on malformed input
      # parsing
  std.expr
    std.expr.compile
      fn (source: string) -> result[expr_handle, string]
      + compiles a boolean predicate expression like "temp > 30 and site == 'A'"
      - returns error on syntax or unknown operator
      # parsing
    std.expr.eval_bool
      fn (handle: expr_handle, bindings: map[string, json_value]) -> result[bool, string]
      + evaluates the compiled expression against the given bindings
      - returns error when a referenced name is missing
      # evaluation

stream
  stream.new
    fn (name: string) -> stream_state
    + creates a new pipeline builder with the given name and no operators
    # construction
  stream.filter
    fn (s: stream_state, predicate: string) -> result[stream_state, string]
    + appends a filter operator that drops events failing the predicate
    - returns error when the predicate fails to compile
    # operators
    -> std.expr.compile
  stream.map
    fn (s: stream_state, assignments: map[string, string]) -> result[stream_state, string]
    + appends a projection that rewrites named fields from compiled expressions
    - returns error when any assignment fails to compile
    # operators
    -> std.expr.compile
  stream.tumbling_window
    fn (s: stream_state, size_ms: i64, agg: string, over: string) -> result[stream_state, string]
    + appends a tumbling window with the given size and aggregator ("sum","avg","count","max","min")
    - returns error when agg is unknown
    - returns error when size_ms <= 0
    # operators
  stream.push
    fn (s: stream_state, event_json: string) -> result[list[output_event], string]
    + ingests a JSON event, runs it through the pipeline, and returns any emitted outputs
    - returns error when the event is not a JSON object
    # execution
    -> std.json.parse_value
    -> std.expr.eval_bool
    -> std.time.now_monotonic_ms
  stream.flush
    fn (s: stream_state) -> list[output_event]
    + closes all open windows and returns any pending aggregates
    # execution
