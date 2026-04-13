# Requirement: "a logging library with severity levels and pluggable handlers"

A structured logger with level thresholds and an ordered chain of handlers (console, file, or user-provided sinks).

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.fs
    std.fs.append_line
      @ (path: string, line: string) -> result[void, string]
      + appends a line to path, creating the file if missing
      - returns error when the directory does not exist
      # filesystem
  std.io
    std.io.print_line
      @ (s: string) -> void
      + writes s followed by a newline to standard output
      # io

logger
  logger.new
    @ (name: string, min_level: level_t) -> logger_state
    + creates a logger with the given channel name and minimum level
    # construction
  logger.add_handler
    @ (state: logger_state, handler: handler_t) -> logger_state
    + appends a handler that receives records meeting the level threshold
    # handlers
  logger.log
    @ (state: logger_state, level: level_t, message: string, context: map[string,string]) -> void
    + dispatches the record to every handler when level >= min_level
    + skips dispatch when level is below min_level
    # dispatch
    -> std.time.now_millis
  logger.console_handler
    @ (min_level: level_t) -> handler_t
    + returns a handler that prints formatted records to stdout
    # handler_console
    -> std.io.print_line
  logger.file_handler
    @ (path: string, min_level: level_t) -> handler_t
    + returns a handler that appends formatted records to a file
    # handler_file
    -> std.fs.append_line
  logger.format_record
    @ (channel: string, level: level_t, ts_millis: i64, message: string, context: map[string,string]) -> string
    + returns "[timestamp] channel.LEVEL: message {k=v ...}"
    # formatting
