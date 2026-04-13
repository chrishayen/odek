# Requirement: "a structured logging library with pluggable sinks and levels"

Levels, structured fields, and pluggable output sinks. Time reads go through a thin std utility for deterministic tests.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.io
    std.io.write_all
      @ (data: string) -> void
      + writes data to stdout
      # io

logger
  logger.new
    @ (min_level: string) -> logger_state
    + creates a logger that drops records below min_level
    ? levels are "debug", "info", "warn", "error"
    # construction
  logger.add_sink
    @ (state: logger_state, sink: fn(string) -> void) -> logger_state
    + registers an output function
    # sinks
  logger.with_field
    @ (state: logger_state, key: string, value: string) -> logger_state
    + returns a logger that injects the field into every record
    # context
  logger.log
    @ (state: logger_state, level: string, message: string, fields: map[string, string]) -> void
    + formats a record and dispatches to every sink
    - drops the record when level is below min_level
    # emit
    -> std.time.now_millis
  logger.format_record
    @ (timestamp_ms: i64, level: string, message: string, fields: map[string, string]) -> string
    + returns a single-line JSON log record
    # formatting
  logger.stdout_sink
    @ (line: string) -> void
    + writes line plus newline to stdout
    # sinks
    -> std.io.write_all
