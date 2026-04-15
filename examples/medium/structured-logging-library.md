# Requirement: "a structured logging library"

Level-filtered logger that emits key-value records to a pluggable writer.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns current unix time in nanoseconds
      # time
  std.io
    std.io.write_line
      fn (w: writer_handle, line: string) -> result[void, string]
      + writes the line followed by a newline
      # io
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

logger
  logger.new
    fn (min_level: log_level, writer: writer_handle) -> logger_state
    + creates a logger that drops records below min_level
    # construction
  logger.with_field
    fn (l: logger_state, key: string, value: string) -> logger_state
    + returns a child logger that always attaches the given key/value
    # context
  logger.log
    fn (l: logger_state, level: log_level, message: string, fields: map[string, string]) -> void
    + merges fields with the logger's context and emits a record
    - does nothing when level is below the logger's min_level
    # logging
    -> std.time.now_nanos
    -> std.json.encode_object
    -> std.io.write_line
  logger.format_record
    fn (ts_nanos: i64, level: log_level, message: string, fields: map[string, string]) -> string
    + produces a single-line JSON record with ts, level, msg, and fields
    # formatting
    -> std.json.encode_object
