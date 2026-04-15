# Requirement: "a microservice framework with pluggable transports, observability, and lifecycle management"

Wires a service together from pluggable transports, exposes a lifecycle, and emits telemetry through thin primitives.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.log
    std.log.emit
      fn (level: string, message: string, fields: map[string, string]) -> void
      + emits a structured log record to the configured backend
      # logging
  std.net
    std.net.listen_tcp
      fn (address: string) -> result[listener, string]
      + binds a TCP listener on the given host:port
      - returns error when the address is already in use
      # network
    std.net.accept
      fn (l: listener) -> result[connection, string]
      + returns the next accepted connection
      - returns error when the listener is closed
      # network

micro
  micro.new_service
    fn (name: string, version: string) -> service_state
    + creates a service state with no transports yet registered
    # construction
  micro.add_transport
    fn (state: service_state, t: transport) -> service_state
    + attaches a transport (HTTP server, queue consumer, etc.) to the service
    ? each transport exposes its own start and stop hooks
    # registration
  micro.add_health_check
    fn (state: service_state, name: string, check: fn() -> result[void, string]) -> service_state
    + registers a named health check invoked by readiness probes
    # health
  micro.start
    fn (state: service_state) -> result[service_state, string]
    + starts every registered transport in order
    - returns error on the first transport that fails to start and stops the rest
    # lifecycle
    -> std.log.emit
  micro.stop
    fn (state: service_state, timeout_ms: i32) -> result[void, string]
    + stops every transport, bounded by the shutdown timeout
    - returns error when any transport fails to stop within the timeout
    # lifecycle
    -> std.log.emit
  micro.record_request
    fn (state: service_state, transport: string, route: string, status: i32, duration_ms: i64) -> void
    + records a request against the service's per-route metrics
    # metrics
    -> std.time.now_millis
  micro.snapshot_metrics
    fn (state: service_state) -> map[string, f64]
    + returns a flat map of current metric values
    # metrics
  micro.ready
    fn (state: service_state) -> result[void, list[string]]
    + returns ok when every health check passes
    - returns the list of failing check names when any fail
    # health
  micro.new_http_transport
    fn (address: string, handler: fn(http_request) -> http_response) -> transport
    + returns a transport that serves HTTP on the given address
    # transport
    -> std.net.listen_tcp
    -> std.net.accept
  micro.new_queue_transport
    fn (queue_name: string, handler: fn(bytes) -> result[void, string]) -> transport
    + returns a transport that consumes messages from a named queue
    # transport
  micro.propagate_trace
    fn (ctx: request_context) -> map[string, string]
    + returns a header map carrying the current trace and span ids
    # tracing
