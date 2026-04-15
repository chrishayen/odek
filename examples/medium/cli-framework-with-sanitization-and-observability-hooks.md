# Requirement: "a command-line framework with input sanitization, sandboxed plugin storage, and structured observability hooks"

A command framework that also wraps common hardening and observability concerns: safe argument handling, a scoped plugin data store, and a structured event emitter.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the whole file into memory
      - returns error when the path does not exist
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data, creating or truncating the file
      - returns error when the parent directory does not exist
      # io
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

cli
  cli.command_new
    fn (name: string, summary: string) -> command_spec
    + constructs a command with an empty flag set
    # definition
  cli.command_flag
    fn (spec: command_spec, name: string, kind: flag_kind, required: bool) -> command_spec
    + appends a flag declaration
    # definition
  cli.parse_safe
    fn (spec: command_spec, args: list[string]) -> result[parsed_invocation, string]
    + parses args while rejecting control characters, null bytes, and oversized values
    - returns error when an argument exceeds the configured size limit
    - returns error when an argument contains disallowed bytes
    # sanitization
  cli.plugin_store_open
    fn (root: string, plugin_id: string) -> result[plugin_store, string]
    + opens a storage scope rooted at root/plugin_id, refusing paths that escape via "..""
    - returns error when plugin_id contains path separators
    # storage
    -> std.fs.read_all
  cli.plugin_store_get
    fn (store: plugin_store, key: string) -> result[optional[bytes], string]
    + returns the value stored under key, or none
    - returns error when the key contains path separators
    # storage
    -> std.fs.read_all
  cli.plugin_store_put
    fn (store: plugin_store, key: string, value: bytes) -> result[void, string]
    + writes a value for the key within the plugin's scope
    - returns error when the key contains path separators
    # storage
    -> std.fs.write_all
  cli.observer_new
    fn (sink: event_sink) -> observer_state
    + constructs an observer that publishes events to a pluggable sink
    # observability
  cli.observer_emit
    fn (state: observer_state, level: string, message: string, fields: map[string, string]) -> void
    + publishes a structured event with a timestamp to the sink
    # observability
    -> std.time.now_millis
  cli.observer_span_start
    fn (state: observer_state, name: string) -> span
    + starts a timed span for the named operation
    # observability
    -> std.time.now_millis
  cli.observer_span_end
    fn (state: observer_state, span: span) -> void
    + emits a completion event carrying the span's duration
    # observability
    -> std.time.now_millis
  cli.dispatch
    fn (spec: command_spec, handler: command_fn, observer: observer_state, args: list[string]) -> result[i32, string]
    + parses args safely, emits a span around the handler, and returns its exit code
    - returns error when parsing fails
    # execution
    -> cli.parse_safe
    -> cli.observer_span_start
    -> cli.observer_span_end
