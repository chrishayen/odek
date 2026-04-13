# Requirement: "a structured logger for context-aware HTTP handlers with pluggable dispatching"

Context-scoped loggers carry field bags; records are dispatched to one or more sinks.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

structlog
  structlog.new_logger
    @ () -> logger_state
    + creates a logger with no fields and no sinks
    # construction
  structlog.add_sink
    @ (logger: logger_state, sink_id: string, handler: sink_fn) -> logger_state
    + registers a sink that receives formatted records
    ? sink_fn is an opaque callable taking a record and returning void
    # dispatch
  structlog.with_field
    @ (logger: logger_state, key: string, value: string) -> logger_state
    + returns a child logger carrying an additional field
    # context
  structlog.with_request
    @ (logger: logger_state, method: string, path: string, request_id: string) -> logger_state
    + attaches common HTTP fields to the logger's context
    # context
  structlog.log
    @ (logger: logger_state, level: string, message: string) -> void
    + formats a record with timestamp, level, message, and fields and fans it out to sinks
    + accepts levels "debug","info","warn","error"
    # emission
    -> std.time.now_millis
    -> std.json.encode_object
  structlog.set_min_level
    @ (logger: logger_state, level: string) -> logger_state
    + records below the minimum level are dropped before dispatch
    # filtering
