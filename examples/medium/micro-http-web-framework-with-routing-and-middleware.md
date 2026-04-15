# Requirement: "a minimal HTTP web micro-framework"

A small routing core: register handlers by method and path, match an incoming request, and produce a response. Request parsing and socket I/O live in std.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from a raw HTTP/1.1 request
      - returns error on malformed request line
      # http_parsing
    std.http.encode_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes status, headers, and body to an HTTP/1.1 response
      # http_encoding
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds and listens on the given address
      - returns error when the port is already bound
      # networking
    std.net.accept
      fn (l: listener) -> result[connection, string]
      + blocks until a connection arrives
      # networking

micro_framework
  micro_framework.new
    fn () -> router_state
    + returns an empty router with no routes registered
    # construction
  micro_framework.route
    fn (state: router_state, method: string, path: string, handler: handler_fn) -> router_state
    + registers a handler for the given method and exact path
    + overwrites a previous handler registered for the same method and path
    # registration
  micro_framework.dispatch
    fn (state: router_state, req: http_request) -> http_response
    + invokes the registered handler and returns its response
    - returns a 404 response when no handler matches
    - returns a 405 response when the path exists but the method does not
    # routing
    -> std.http.parse_request
  micro_framework.serve
    fn (state: router_state, host: string, port: i32) -> result[void, string]
    + accepts connections and serves each request through dispatch
    - returns error when binding fails
    # serving
    -> std.net.listen_tcp
    -> std.net.accept
    -> std.http.encode_response
