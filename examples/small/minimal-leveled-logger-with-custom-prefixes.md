# Requirement: "a minimal leveled logger with custom prefixes"

A logger struct carrying a level threshold and a prefix string. Writing goes through a thin std sink.

std
  std.io
    std.io.print_line
      fn (line: string) -> void
      + writes a line to standard output
      # io

logger
  logger.new
    fn (prefix: string, min_level: i32) -> logger_state
    + creates a logger with the given prefix and minimum level
    ? levels: 0=debug, 1=info, 2=warn, 3=error
    # construction
  logger.log
    fn (state: logger_state, level: i32, message: string) -> void
    + writes "<prefix> [<level_name>] <message>" when level >= threshold
    - writes nothing when level is below threshold
    # logging
    -> std.io.print_line
