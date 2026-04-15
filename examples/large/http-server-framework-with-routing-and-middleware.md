# Requirement: "an HTTP server framework with routing and middleware"

A minimal framework: build a router, register handlers and middleware, then run a request through the pipeline to get a response.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener_state, string]
      + binds a TCP listener
      - returns error when the address is already in use
      # networking
    std.net.accept
      fn (listener: listener_state) -> result[conn_state, string]
      + accepts the next incoming connection
      # networking
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request line, headers, and body
      - returns error on malformed input
      # http
    std.http.write_response
      fn (conn: conn_state, resp: http_response) -> result[void, string]
      + writes an HTTP response to a connection
      # http

server
  server.new_router
    fn () -> router_state
    + creates an empty router with no routes or middleware
    # construction
  server.handle
    fn (router: router_state, method: string, path: string, handler: handler_fn) -> router_state
    + registers a handler for an exact method and path pattern
    ? path patterns support ":param" segments for path variables
    # routing
  server.use
    fn (router: router_state, middleware: middleware_fn) -> router_state
    + appends middleware to the pipeline in registration order
    # middleware
  server.group
    fn (router: router_state, prefix: string) -> router_state
    + returns a router whose subsequent routes are prefixed and inherit current middleware
    # grouping
  server.match_route
    fn (router: router_state, method: string, path: string) -> result[tuple[handler_fn, map[string,string]], string]
    + returns the handler and extracted path params for a request
    - returns error "not found" when no route matches
    - returns error "method not allowed" when path matches but method does not
    # lookup
  server.dispatch
    fn (router: router_state, req: http_request) -> http_response
    + runs the request through middleware and handler, returning the final response
    + returns 404 response when no route matches
    # request_pipeline
  server.serve
    fn (router: router_state, host: string, port: i32) -> result[void, string]
    + binds a listener and serves requests until the listener closes
    # serving
    -> std.net.listen_tcp
    -> std.net.accept
    -> std.http.parse_request
    -> std.http.write_response
  server.json_response
    fn (status: i32, body: string) -> http_response
    + builds a response with Content-Type application/json and the given body
    # response_helpers
  server.text_response
    fn (status: i32, body: string) -> http_response
    + builds a response with Content-Type text/plain and the given body
    # response_helpers
  server.redirect
    fn (status: i32, location: string) -> http_response
    + builds a redirect response with a Location header
    - status must be between 300 and 399
    # response_helpers
