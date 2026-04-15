# Requirement: "a web framework with low-level control"

A routing and middleware framework that maps HTTP requests to handlers with path parameters and a layered middleware chain.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed request line or headers
      # http
    std.http.serialize_response
      fn (resp: http_response) -> bytes
      + returns a full HTTP/1.1 response byte sequence
      # http
    std.http.url_decode
      fn (s: string) -> result[string, string]
      + decodes percent-encoded characters
      - returns error on truncated escape
      # http
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits on the literal separator
      # text

web_framework
  web_framework.new_router
    fn () -> router_state
    + creates an empty router with no routes and no middleware
    # construction
  web_framework.route
    fn (router: router_state, method: string, pattern: string, handler_id: string) -> result[router_state, string]
    + registers a handler for a method+pattern; supports ":param" and "*rest" segments
    - returns error when pattern is empty or contains invalid segments
    # routing
    -> std.text.split
  web_framework.use
    fn (router: router_state, middleware_id: string) -> router_state
    + appends middleware to the global chain
    # middleware
  web_framework.group
    fn (router: router_state, prefix: string, handler_ids: list[tuple[string, string, string]]) -> result[router_state, string]
    + registers multiple (method, pattern, handler_id) under a common prefix
    # routing
  web_framework.match
    fn (router: router_state, method: string, path: string) -> optional[tuple[string, map[string, string]]]
    + returns (handler_id, path_params) for the best matching route
    - returns none when nothing matches
    # routing
    -> std.text.split
    -> std.http.url_decode
  web_framework.dispatch
    fn (router: router_state, req: http_request, handlers: map[string, handler], middleware: map[string, middleware]) -> http_response
    + runs the middleware chain then the matched handler
    + returns a 404 response when no route matches
    + returns a 405 response when path matches but method does not
    # dispatch
  web_framework.not_found_handler
    fn (router: router_state, handler_id: string) -> router_state
    + overrides the default 404 handler
    # routing
  web_framework.parse_query
    fn (query: string) -> map[string, string]
    + returns key/value pairs from a url-encoded query string
    # request
    -> std.http.url_decode
    -> std.text.split
  web_framework.handle_raw
    fn (router: router_state, raw: bytes, handlers: map[string, handler], middleware: map[string, middleware]) -> bytes
    + parses raw request bytes, dispatches, and returns serialized response bytes
    # io
    -> std.http.parse_request
    -> std.http.serialize_response
  web_framework.add_error_recovery
    fn (router: router_state) -> router_state
    + installs a middleware that returns 500 with a safe body on panic
    # middleware
