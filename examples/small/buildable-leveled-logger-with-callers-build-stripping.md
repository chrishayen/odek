# Requirement: "a leveled logger whose inactive levels can be stripped by the caller's build"

The library exposes per-level log entry points that each consult a build-time minimum level. When a call is below the minimum, it returns immediately with no formatting — allowing a tree-shaker or dead-code step to remove it entirely.

std
  std.io
    std.io.print_string
      fn (s: string) -> void
      + writes s followed by a newline to standard output
      # io

build_logger
  build_logger.new
    fn (min_level: i32) -> logger_state
    + creates a logger that drops messages below min_level (0=debug, 1=info, 2=warn, 3=error)
    # construction
  build_logger.debug
    fn (state: logger_state, message: string) -> void
    + prints message when min_level is 0
    - does nothing when min_level is greater than 0
    # logging
    -> std.io.print_string
  build_logger.info
    fn (state: logger_state, message: string) -> void
    + prints message when min_level is at most 1
    - does nothing otherwise
    # logging
    -> std.io.print_string
  build_logger.warn
    fn (state: logger_state, message: string) -> void
    + prints message when min_level is at most 2
    - does nothing otherwise
    # logging
    -> std.io.print_string
  build_logger.error
    fn (state: logger_state, message: string) -> void
    + always prints message
    # logging
    -> std.io.print_string
