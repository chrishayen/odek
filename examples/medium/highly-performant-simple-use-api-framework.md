# Requirement: "an HTTP API framework"

A minimal router, handler registration, and request dispatch. Network I/O is a std primitive.

std
  std.net
    std.net.http_listen
      @ (host: string, port: i32) -> result[listener_handle, string]
      + binds an HTTP listener to the given host and port
      - returns error when the port is already in use
      # networking
    std.net.http_accept
      @ (lis: listener_handle) -> result[http_request, string]
      + blocks until the next HTTP request arrives
      # networking
    std.net.http_respond
      @ (req: http_request, status: i32, headers: map[string, string], body: bytes) -> result[void, string]
      + writes a response to the given request
      # networking

api_framework
  api_framework.new
    @ () -> framework_state
    + returns a framework with an empty route table
    # construction
  api_framework.route
    @ (fw: framework_state, method: string, path: string, handler_id: i64) -> framework_state
    + registers a handler for an exact method and path
    + accepts path segments starting with ":" as captured parameters
    # routing
  api_framework.use
    @ (fw: framework_state, middleware_id: i64) -> framework_state
    + appends a middleware invoked before handlers in registration order
    # middleware
  api_framework.match_route
    @ (fw: framework_state, method: string, path: string) -> optional[tuple[i64, map[string, string]]]
    + returns the matched handler id and captured parameters
    - returns none when no route matches
    # routing
  api_framework.listen
    @ (fw: framework_state, host: string, port: i32) -> result[void, string]
    + accepts requests and dispatches them through middleware and handlers
    - returns error when the listener cannot be bound
    # serving
    -> std.net.http_listen
    -> std.net.http_accept
    -> std.net.http_respond
