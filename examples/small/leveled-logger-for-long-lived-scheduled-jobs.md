# Requirement: "a leveled logger suitable for long-lived scheduled jobs"

Filters by level and formats each line with a timestamp. Writing the line is the caller's problem via an injected sink.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns the current unix timestamp in seconds
      # time

leveled_logger
  leveled_logger.new
    fn (min_level: string, sink: log_sink) -> logger
    + constructs a logger that drops entries below min_level
    ? levels are "debug", "info", "warn", "error" from low to high
    # construction
  leveled_logger.should_log
    fn (l: logger, level: string) -> bool
    + returns true when level >= min_level
    - returns false for an unknown level
    # filter
  leveled_logger.format_entry
    fn (level: string, message: string, timestamp: i64) -> string
    + returns "<iso-time> <LEVEL> message"
    # format
  leveled_logger.log
    fn (l: logger, level: string, message: string) -> result[void, string]
    + formats and writes the entry when should_log returns true
    - returns error when the sink fails
    ? silently drops the entry when the level is filtered out
    # logging
    -> std.time.now_seconds
