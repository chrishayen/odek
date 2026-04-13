# Requirement: "a traffic orchestration system with observability and extensibility"

A library that routes inbound requests across backend pools with health checks, load balancing, and pluggable filters. Metrics and tracing hooks are exposed for external observability backends.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.net
    std.net.resolve_host
      @ (host: string) -> result[list[string], string]
      + returns IP addresses for a hostname
      - returns error on DNS failure
      # networking
  std.hash
    std.hash.fnv32
      @ (data: bytes) -> u32
      + returns a 32-bit FNV-1a hash
      # hashing

traffic
  traffic.new_router
    @ () -> router_state
    + creates an empty router with no routes
    # construction
  traffic.add_backend_pool
    @ (state: router_state, pool_name: string, endpoints: list[string]) -> router_state
    + registers a named pool with a list of upstream endpoints
    - returns unchanged state when pool_name is empty
    # pool_registration
  traffic.add_route
    @ (state: router_state, path_prefix: string, pool_name: string) -> result[router_state, string]
    + maps a URL prefix to a named pool
    - returns error when the pool is unknown
    # routing
  traffic.register_filter
    @ (state: router_state, phase: string, filter_name: string) -> router_state
    + attaches a filter to a pipeline phase (pre, post)
    ? filters are referenced by name; the caller registers implementations separately
    # extensibility
  traffic.pick_endpoint
    @ (state: router_state, path: string, client_key: string) -> result[string, string]
    + returns a healthy endpoint chosen by consistent hash on client_key
    - returns error when no route matches the path
    - returns error when all endpoints in the matched pool are unhealthy
    # load_balancing
    -> std.hash.fnv32
  traffic.mark_endpoint_down
    @ (state: router_state, endpoint: string) -> router_state
    + flags an endpoint as unhealthy so pick_endpoint skips it
    # health
    -> std.time.now_millis
  traffic.mark_endpoint_up
    @ (state: router_state, endpoint: string) -> router_state
    + clears the unhealthy flag on an endpoint
    # health
  traffic.run_health_checks
    @ (state: router_state, check: string) -> router_state
    + re-probes known endpoints using the named check type and updates health
    ? probe execution itself is the caller's responsibility; this updates bookkeeping
    # health
    -> std.net.resolve_host
    -> std.time.now_millis
  traffic.emit_metric
    @ (state: router_state, name: string, value: f64) -> router_state
    + appends a metric sample to the internal buffer for later observation
    # observability
    -> std.time.now_millis
  traffic.drain_metrics
    @ (state: router_state) -> tuple[list[metric_sample], router_state]
    + returns and clears the buffered metrics so they can be forwarded externally
    # observability
