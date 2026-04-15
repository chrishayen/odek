# Requirement: "an http proxy library that can inject configurable failure scenarios and network conditions"

The project exposes a proxy handler that, given a request, consults a chain of interceptors before forwarding. Interceptors can inject latency, drop bodies, or return synthetic errors.

std
  std.net
    std.net.forward_request
      fn (target: string, request: http_request) -> result[http_response, string]
      + sends a request upstream and returns the response
      - returns error when the upstream is unreachable
      # networking
  std.time
    std.time.sleep_ms
      fn (millis: i64) -> void
      + blocks the current task for the given number of milliseconds
      # time
  std.random
    std.random.bernoulli
      fn (probability: f64) -> bool
      + returns true with the given probability in [0, 1]
      # randomness

fault_proxy
  fault_proxy.new
    fn (upstream: string) -> fault_proxy_state
    + creates a proxy that forwards to the given upstream by default
    # construction
  fault_proxy.add_latency
    fn (state: fault_proxy_state, millis: i64, match_path: string) -> fault_proxy_state
    + registers a latency injector for requests matching the path
    # interception
  fault_proxy.add_error
    fn (state: fault_proxy_state, status: i32, probability: f64, match_path: string) -> fault_proxy_state
    + registers a synthetic error responder applied to the fraction of matching requests
    # interception
  fault_proxy.add_drop
    fn (state: fault_proxy_state, probability: f64, match_path: string) -> fault_proxy_state
    + registers a connection-drop injector applied to the fraction of matching requests
    # interception
  fault_proxy.handle
    fn (state: fault_proxy_state, request: http_request) -> result[http_response, string]
    + runs all matching interceptors then forwards to upstream if none short-circuit
    + applies configured latency before forwarding
    - returns the configured synthetic error when an error interceptor fires
    - returns error when a drop interceptor fires
    # dispatch
    -> std.time.sleep_ms
    -> std.random.bernoulli
    -> std.net.forward_request
