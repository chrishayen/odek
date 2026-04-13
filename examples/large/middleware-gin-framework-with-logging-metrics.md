# Requirement: "middleware providing logging, metrics, auth, and tracing for an HTTP server framework"

Four general-purpose middleware families: structured logging, request metrics, token authentication, and distributed tracing spans.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current time in milliseconds
      # time
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
  std.io
    std.io.print_string
      @ (s: string) -> void
      + writes a line to stdout
      # io

srvkit
  srvkit.new_logger
    @ (service_name: string) -> logger_state
    + creates a logger tagged with service_name
    # logging
  srvkit.log_event
    @ (logger: logger_state, level: string, fields: map[string,string]) -> void
    + emits a structured log line
    # logging
    -> std.io.print_string
  srvkit.logger_middleware
    @ (logger: logger_state, next: handler_fn) -> handler_fn
    + wraps a handler to log request and response lines
    # composition
    -> std.time.now_millis
    -> srvkit.log_event
  srvkit.new_metrics
    @ () -> metrics_state
    + creates an empty metrics collector
    # metrics
  srvkit.record_request
    @ (metrics: metrics_state, route: string, status: i32, duration_ms: i64) -> metrics_state
    + records a request observation
    # metrics
  srvkit.metrics_middleware
    @ (metrics: metrics_state, next: handler_fn) -> handler_fn
    + wraps a handler to record timing and status metrics
    # composition
    -> std.time.now_millis
  srvkit.verify_token
    @ (token: string, secret: bytes) -> result[map[string,string], string]
    + returns claims when the token signature is valid
    - returns error when the signature does not match
    # auth
    -> std.crypto.hmac_sha256
  srvkit.auth_middleware
    @ (secret: bytes, next: handler_fn) -> handler_fn
    + wraps a handler to require a valid bearer token
    - short-circuits with 401 when no token is present
    # auth
    -> srvkit.verify_token
  srvkit.new_tracer
    @ (service_name: string) -> tracer_state
    + creates a tracer tagged with service_name
    # tracing
  srvkit.start_span
    @ (tracer: tracer_state, parent_id: optional[string], name: string) -> span_state
    + opens a span with a unique id and start time
    # tracing
    -> std.time.now_millis
    -> std.encoding.hex_encode
  srvkit.finish_span
    @ (span: span_state) -> finished_span
    + closes the span and records its duration
    # tracing
    -> std.time.now_millis
  srvkit.tracing_middleware
    @ (tracer: tracer_state, next: handler_fn) -> handler_fn
    + wraps a handler to start a span per request and finish it after
    # composition
    -> srvkit.start_span
    -> srvkit.finish_span
