# Requirement: "a middleware-based HTTP server framework"

A server framework centered on composable middleware. A request flows through an ordered chain; each middleware may short-circuit, mutate, or pass through.

std
  std.http
    std.http.listen
      fn (port: i32, handler: http_handler_fn) -> result[void, string]
      + starts an HTTP listener bound to the given port
      - returns error when the port is in use
      # http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses a raw HTTP request
      - returns error on malformed header
      # http
    std.http.build_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + builds an HTTP response bytes
      # http

iron
  iron.new_chain
    fn () -> chain_state
    + creates an empty middleware chain
    # construction
  iron.link
    fn (chain: chain_state, middleware: middleware_fn) -> chain_state
    + appends a middleware to the chain
    ? each middleware may return a response directly, forward to next, or return an error
    # composition
  iron.handle
    fn (chain: chain_state, req: http_request) -> result[http_response, string]
    + runs the request through the chain in registration order
    + returns the first response produced by any middleware
    - returns error when no middleware produces a response
    # dispatch
  iron.serve
    fn (chain: chain_state, port: i32) -> result[void, string]
    + starts an HTTP listener dispatching each request through the chain
    # transport
    -> std.http.listen
    -> std.http.parse_request
    -> std.http.build_response
    -> iron.handle
  iron.before_after
    fn (before: middleware_fn, inner: middleware_fn, after: middleware_fn) -> middleware_fn
    + composes a before-hook, an inner handler, and an after-hook into one middleware
    # composition
