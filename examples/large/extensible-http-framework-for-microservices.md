# Requirement: "an extensible HTTP framework for building microservices"

Routing, middleware chains, and request/response helpers. The underlying transport is an opaque listener.

std
  std.net
    std.net.tcp_listen
      @ (addr: string, port: u16) -> result[tcp_listener, string]
      + binds a TCP listener
      - returns error when the port is in use
      # networking
    std.net.tcp_accept
      @ (l: tcp_listener) -> result[tcp_conn, string]
      + returns the next accepted connection
      # networking
  std.io
    std.io.read_all
      @ (conn: tcp_conn) -> result[bytes, string]
      + reads until the connection is closed or an error occurs
      # io
    std.io.write_all
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # io

http_framework
  http_framework.new_router
    @ () -> router_state
    + creates an empty router with no routes or middleware
    # construction
  http_framework.add_route
    @ (r: router_state, method: string, pattern: string, handler: fn(request) -> response) -> router_state
    + registers a handler for the given method and path pattern
    + pattern may contain ":name" segments bound to request params
    # routing
  http_framework.use_middleware
    @ (r: router_state, mw: fn(fn(request) -> response) -> fn(request) -> response) -> router_state
    + appends a middleware that wraps downstream handlers
    # middleware
  http_framework.match_route
    @ (r: router_state, method: string, path: string) -> optional[route_match]
    + returns the matched handler and extracted path params
    - returns none when no route matches
    # routing
  http_framework.parse_request
    @ (raw: bytes) -> result[request, string]
    + parses request line, headers, and body into a request record
    - returns error on malformed start line
    # parsing
  http_framework.encode_response
    @ (resp: response) -> bytes
    + serializes status line, headers, and body
    # serialization
  http_framework.handle_connection
    @ (r: router_state, conn: tcp_conn) -> void
    + reads one request, runs it through middleware and a matched handler, and writes the response
    + returns a 404 response when no route matches
    # request_handling
    -> std.io.read_all
    -> std.io.write_all
    -> http_framework.parse_request
    -> http_framework.match_route
    -> http_framework.encode_response
  http_framework.serve
    @ (r: router_state, addr: string, port: u16) -> result[void, string]
    + accepts connections in a loop and dispatches them to handle_connection
    # server
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> http_framework.handle_connection
  http_framework.json_response
    @ (status: i32, body: string) -> response
    + builds a response with Content-Type: application/json
    # helpers
