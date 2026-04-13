# Requirement: "middleware providing logging, metrics, auth, and tracing for an RPC framework"

Interceptors for an RPC framework, mirroring HTTP middleware but scoped to unary and streaming calls.

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

rpckit
  rpckit.new_logger
    @ (service_name: string) -> logger_state
    + creates a logger tagged with service_name
    # logging
  rpckit.log_call
    @ (logger: logger_state, method: string, status: string, duration_ms: i64) -> void
    + emits a structured log line describing a completed call
    # logging
    -> std.io.print_string
  rpckit.logger_unary
    @ (logger: logger_state, handler: unary_handler_fn) -> unary_handler_fn
    + wraps a unary handler with start and end logging
    # composition
    -> std.time.now_millis
    -> rpckit.log_call
  rpckit.logger_stream
    @ (logger: logger_state, handler: stream_handler_fn) -> stream_handler_fn
    + wraps a streaming handler with start and end logging
    # composition
    -> std.time.now_millis
    -> rpckit.log_call
  rpckit.new_metrics
    @ () -> metrics_state
    + creates an empty metrics collector
    # metrics
  rpckit.record_call
    @ (metrics: metrics_state, method: string, status: string, duration_ms: i64) -> metrics_state
    + records a single call observation
    # metrics
  rpckit.metrics_unary
    @ (metrics: metrics_state, handler: unary_handler_fn) -> unary_handler_fn
    + wraps a unary handler to record metrics
    # composition
    -> std.time.now_millis
  rpckit.verify_token
    @ (token: string, secret: bytes) -> result[map[string,string], string]
    + returns claims when the token signature is valid
    - returns error when the signature does not match
    # auth
    -> std.crypto.hmac_sha256
  rpckit.auth_unary
    @ (secret: bytes, handler: unary_handler_fn) -> unary_handler_fn
    + wraps a unary handler to require a valid bearer token in metadata
    - short-circuits with Unauthenticated when token is missing or invalid
    # auth
    -> rpckit.verify_token
  rpckit.new_tracer
    @ (service_name: string) -> tracer_state
    + creates a tracer tagged with service_name
    # tracing
  rpckit.start_span
    @ (tracer: tracer_state, parent_id: optional[string], name: string) -> span_state
    + opens a span
    # tracing
    -> std.time.now_millis
    -> std.encoding.hex_encode
  rpckit.finish_span
    @ (span: span_state) -> finished_span
    + closes the span
    # tracing
    -> std.time.now_millis
  rpckit.tracing_unary
    @ (tracer: tracer_state, handler: unary_handler_fn) -> unary_handler_fn
    + wraps a unary handler to open and close a span around the call
    # composition
    -> rpckit.start_span
    -> rpckit.finish_span
