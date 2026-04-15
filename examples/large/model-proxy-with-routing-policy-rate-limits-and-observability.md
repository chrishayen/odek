# Requirement: "a proxy library that routes requests between clients and model inference backends with policy, rate limits, and observability"

Sits in front of one or more inference backends, applying routing, rate limits, and tracing without being an executable.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns 32 bytes of SHA-256 digest
      # cryptography
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

model_proxy
  model_proxy.new
    fn () -> proxy_state
    + returns an empty proxy with no backends and no policies
    # construction
  model_proxy.register_backend
    fn (state: proxy_state, name: string, endpoint: string, weight: i32) -> proxy_state
    + adds a backend with a positive routing weight
    - returns error-marker state when weight <= 0
    # backends
  model_proxy.remove_backend
    fn (state: proxy_state, name: string) -> proxy_state
    + removes the named backend from the pool
    # backends
  model_proxy.add_route
    fn (state: proxy_state, model: string, backend_names: list[string]) -> proxy_state
    + maps a logical model name to its set of candidate backends
    # routing
  model_proxy.pick_backend
    fn (state: proxy_state, model: string, key: string) -> result[string, string]
    + returns a chosen backend name using weighted hashing of key
    - returns error when the model is unknown or has no healthy backends
    # routing
    -> std.crypto.sha256
  model_proxy.set_rate_limit
    fn (state: proxy_state, tenant: string, rps: f64) -> proxy_state
    + sets a per-tenant requests-per-second budget
    # policy
  model_proxy.check_rate_limit
    fn (state: proxy_state, tenant: string) -> tuple[bool, proxy_state]
    + returns (true, new_state) when the tenant has budget
    - returns (false, unchanged_state) when the tenant has exceeded rps
    # policy
    -> std.time.now_millis
  model_proxy.mark_healthy
    fn (state: proxy_state, backend: string) -> proxy_state
    + flags the backend as healthy
    # health
  model_proxy.mark_unhealthy
    fn (state: proxy_state, backend: string) -> proxy_state
    + flags the backend as unhealthy and excludes it from routing
    # health
  model_proxy.build_request
    fn (model: string, prompt: string, tenant: string) -> string
    + returns a canonical JSON request body
    # transport
    -> std.json.encode_object
  model_proxy.parse_response
    fn (raw: string) -> result[map[string, string], string]
    + returns a parsed response body with normalized fields
    - returns error on invalid JSON
    # transport
    -> std.json.parse_object
  model_proxy.record_trace
    fn (state: proxy_state, tenant: string, backend: string, latency_ms: i64) -> proxy_state
    + appends a trace entry with tenant, backend, and latency
    # observability
    -> std.time.now_millis
  model_proxy.trace_snapshot
    fn (state: proxy_state) -> list[trace_entry]
    + returns the in-memory trace buffer in chronological order
    # observability
