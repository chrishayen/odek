# Requirement: "a simple minimalist log system with features for debugging and differentiation of messages"

A leveled logger that formats records. The caller owns the sink; we return formatted strings.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

log
  log.new
    @ (min_level: log_level) -> logger_state
    + creates a logger that drops records below min_level
    # construction
  log.record
    @ (state: logger_state, level: log_level, message: string) -> optional[string]
    + returns a formatted line "<ts> <level> <message>" when level >= min_level
    - returns none when level is below min_level
    # logging
    -> std.time.now_millis
  log.set_level
    @ (state: logger_state, level: log_level) -> logger_state
    + returns a new state with the given minimum level
    # configuration
