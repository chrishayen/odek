# Requirement: "a distilled levelled logging library"

Minimal logger: severity threshold, structured output to a writer, no scopes or hierarchy.

std
  std.io
    std.io.write_string
      fn (sink: writer, line: string) -> result[void, string]
      + appends line to sink
      # io
  std.time
    std.time.now_iso8601
      fn () -> string
      + returns the current time as an iso-8601 timestamp
      # time

level_log
  level_log.new
    fn (sink: writer, min_level: log_level) -> logger_state
    + returns a logger that writes to sink and drops records below min_level
    # construction
  level_log.set_level
    fn (state: logger_state, min_level: log_level) -> logger_state
    + returns a logger with an updated threshold
    # configuration
  level_log.log
    fn (state: logger_state, level: log_level, message: string) -> result[void, string]
    + writes a formatted record as "<timestamp> <level> <message>" when level >= threshold
    - writes nothing when level is below threshold
    # logging
    -> std.time.now_iso8601
    -> std.io.write_string
