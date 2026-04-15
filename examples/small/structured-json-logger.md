# Requirement: "a structured JSON logger"

The library emits one JSON log line per event with fields for level, timestamp, message, and attached key-value pairs.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

jlog
  jlog.new_logger
    fn (level: string) -> logger_state
    + returns a logger that emits events at or above the given level
    ? levels are ordered debug < info < warn < error
    # construction
  jlog.with_field
    fn (logger: logger_state, key: string, value: string) -> logger_state
    + returns a logger with an additional persistent field
    # context
  jlog.format_event
    fn (logger: logger_state, level: string, message: string, fields: map[string, string]) -> optional[string]
    + returns a single-line JSON record with ts, level, msg, and merged fields
    - returns none when level is below the logger's threshold
    # formatting
    -> std.time.now_millis
