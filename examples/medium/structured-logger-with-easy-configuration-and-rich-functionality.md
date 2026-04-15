# Requirement: "a structured logger with easy configuration and rich functionality"

Structured key-value logging with levels, named loggers, and pluggable sinks.

std
  std.time
    std.time.now_iso8601
      fn () -> string
      + returns current UTC time formatted as ISO 8601
      # time
  std.io
    std.io.write_line
      fn (fd: i32, line: string) -> result[void, string]
      + writes line followed by a newline to file descriptor fd
      # io

zl
  zl.new
    fn (name: string, min_level: log_level) -> logger_state
    + constructs a named logger that filters messages below min_level
    ? levels are ordered: debug < info < warn < error < fatal
    # construction
  zl.with_field
    fn (logger: logger_state, key: string, value: string) -> logger_state
    + returns a child logger that attaches the field to every subsequent message
    # context
  zl.with_fields
    fn (logger: logger_state, fields: map[string, string]) -> logger_state
    + returns a child logger with multiple fields attached
    # context
  zl.format_entry
    fn (logger: logger_state, level: log_level, message: string, fields: map[string, string]) -> string
    + renders a log entry as a JSON line including timestamp, level, name, message, and fields
    # formatting
    -> std.time.now_iso8601
  zl.log
    fn (logger: logger_state, level: log_level, message: string, fields: map[string, string]) -> void
    + emits a log entry when level passes the logger's minimum filter
    - entries below min_level are suppressed
    # emission
    -> zl.format_entry
    -> std.io.write_line
  zl.set_sink
    fn (logger: logger_state, sink: fn(line: string) -> void) -> logger_state
    + replaces the default stdout sink with a caller-supplied function
    # sinks
