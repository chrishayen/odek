# Requirement: "a minimal leveled logging library"

Produces formatted log records at configurable levels; emission is the caller's responsibility.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time

yell
  yell.new
    @ (min_level: log_level) -> logger_state
    + creates a logger that suppresses records below min_level
    # construction
  yell.log
    @ (state: logger_state, level: log_level, message: string) -> optional[string]
    + returns a formatted "timestamp level message" record when level meets threshold
    - returns none when level is below the logger's minimum
    # logging
    -> std.time.now_millis
  yell.set_level
    @ (state: logger_state, level: log_level) -> logger_state
    + returns a logger state with the updated minimum level
    # configuration
