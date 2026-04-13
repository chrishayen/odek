# Requirement: "a simple logging framework"

Structured, leveled logs emitted to a pluggable sink.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

log
  log.new
    @ (min_level: i32, sink: log_sink) -> logger
    + builds a logger that drops entries below min_level
    ? levels: 0=debug 1=info 2=warn 3=error
    # construction
  log.with_field
    @ (l: logger, key: string, value: string) -> logger
    + returns a child logger carrying an additional structured field
    # context
  log.emit
    @ (l: logger, level: i32, message: string) -> void
    + writes an entry to the sink when level >= min_level
    - is a no-op when level < min_level
    # logging
    -> std.time.now_millis
  log.stdout_sink
    @ () -> log_sink
    + returns a sink that formats entries as line-delimited JSON to standard output
    # sink
