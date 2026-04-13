# Requirement: "a multi-level logging library"

Holds a level threshold and a list of sinks. Formatting is line-oriented and synchronous; the caller plugs in sinks.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

logdump
  logdump.new_logger
    @ (level: string) -> logger_state
    + creates a logger at the given minimum level
    ? level is one of "debug", "info", "warn", "error"
    # construction
  logdump.add_sink
    @ (state: logger_state, sink_name: string) -> logger_state
    + registers a named sink that will receive formatted records
    # sinks
  logdump.set_level
    @ (state: logger_state, level: string) -> result[logger_state, string]
    + changes the active threshold
    - returns error on an unknown level name
    # configuration
  logdump.log
    @ (state: logger_state, level: string, message: string) -> list[log_record]
    + returns the records that would be emitted to each sink, or an empty list when below threshold
    + includes a millisecond timestamp on every record
    # logging
    -> std.time.now_millis
