# Requirement: "middleware providing logging and metrics for an HTTP server framework"

Two general-purpose middleware builders: one logs every request, the other records latency and status counts.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current time in milliseconds
      # time
  std.io
    std.io.print_string
      @ (s: string) -> void
      + writes a line to stdout
      # io

echomw
  echomw.new_logger
    @ (format: string) -> logger_state
    + creates a logger that formats lines using the given template
    ? template supports placeholders for method, path, status, and duration_ms
    # logging
  echomw.log_request
    @ (logger: logger_state, method: string, path: string, status: i32, start_ms: i64) -> void
    + emits a log line describing a completed request
    # logging
    -> std.time.now_millis
    -> std.io.print_string
  echomw.new_metrics
    @ () -> metrics_state
    + creates an empty metrics collector
    # metrics
  echomw.record
    @ (metrics: metrics_state, path: string, status: i32, duration_ms: i64) -> metrics_state
    + increments the request counter and updates the latency histogram for (path, status)
    # metrics
  echomw.snapshot
    @ (metrics: metrics_state) -> map[string,i64]
    + returns a flat map of metric_name to value
    # metrics
  echomw.logger_middleware
    @ (logger: logger_state, metrics: metrics_state, next: handler_fn) -> handler_fn
    + returns a handler that logs and records metrics around the inner handler
    # composition
    -> std.time.now_millis
    -> echomw.log_request
    -> echomw.record
