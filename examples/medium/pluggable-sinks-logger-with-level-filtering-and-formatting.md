# Requirement: "a logging library with pluggable sinks, level filtering, and message formatting"

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

logger
  logger.new
    @ (min_level: log_level) -> logger_state
    + creates a logger that drops records below min_level and has no sinks
    # construction
  logger.add_sink
    @ (state: logger_state, name: string, accepts: log_level) -> logger_state
    + registers a named sink that receives records at or above the given level
    ? sinks identify themselves by name so the caller can attach transports externally
    # registration
  logger.set_format
    @ (state: logger_state, pattern: string) -> logger_state
    + sets a format pattern with placeholders for time, level, and message
    ? supported placeholders are {time}, {level}, and {msg}
    # configuration
  logger.log
    @ (state: logger_state, level: log_level, message: string) -> list[tuple[string, string]]
    + returns (sink_name, formatted_line) pairs for every sink that accepts this level
    - returns an empty list when the record level is below min_level
    # dispatch
    -> std.time.now_millis
  logger.format_record
    @ (pattern: string, level: log_level, message: string, timestamp_ms: i64) -> string
    + substitutes the pattern placeholders with the record fields
    # formatting
  logger.level_at_least
    @ (level: log_level, min: log_level) -> bool
    + returns true when level is numerically at least min
    # filtering
