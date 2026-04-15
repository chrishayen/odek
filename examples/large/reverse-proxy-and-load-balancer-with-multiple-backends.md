# Requirement: "a reverse proxy and load balancer supporting multiple backends"

Stateful router: register backends, track health, pick one per request with a chosen strategy, and forward.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses a request line, headers, and body
      - returns error on malformed start line
      # http_parsing
    std.http.encode_request
      fn (req: http_request) -> bytes
      + serializes a request into wire bytes
      # http_encoding
    std.http.parse_response
      fn (raw: bytes) -> result[http_response, string]
      + parses a status line, headers, and body
      - returns error on malformed input
      # http_parsing
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

balancer
  balancer.new
    fn (strategy: u8) -> balancer_state
    + creates an empty balancer with strategy (0=round-robin, 1=least-connections, 2=random)
    # construction
  balancer.add_backend
    fn (state: balancer_state, id: string, address: string) -> balancer_state
    + appends a healthy backend with zero active connections
    # backend_registry
  balancer.remove_backend
    fn (state: balancer_state, id: string) -> balancer_state
    + drops the named backend from the pool
    - returns unchanged state when id is unknown
    # backend_registry
  balancer.mark_unhealthy
    fn (state: balancer_state, id: string, cool_off_ms: i64) -> balancer_state
    + marks the backend unhealthy until the cool-off elapses
    # health
    -> std.time.now_millis
  balancer.refresh_health
    fn (state: balancer_state) -> balancer_state
    + restores backends whose cool-off has expired
    # health
    -> std.time.now_millis
  balancer.pick
    fn (state: balancer_state) -> result[tuple[string, balancer_state], string]
    + returns (backend_id, updated_state) according to the strategy
    - returns error when no healthy backends exist
    # selection
  balancer.release
    fn (state: balancer_state, id: string) -> balancer_state
    + decrements the active-connection counter for a backend
    # selection
  balancer.forward
    fn (state: balancer_state, raw: bytes) -> result[tuple[bytes, balancer_state], string]
    + picks a backend, rewrites the request, and returns the upstream response bytes
    - returns error when the pick step fails
    # forwarding
    -> std.http.parse_request
    -> std.http.encode_request
    -> std.http.parse_response
