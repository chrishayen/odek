# Requirement: "a safe, fast web framework"

Routing, middleware pipeline, typed request/response, and server startup. The core machinery lives in std.

std
  std.http
    std.http.serve
      @ (addr: string, handler: http_handler) -> result[void, string]
      + starts an HTTP server dispatching every request to the handler
      - returns error when the address cannot be bound
      # http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses a raw request into method, path, headers, body
      - returns error on malformed request line
      # http
    std.http.write_response
      @ (resp: http_response, status: i32, body: bytes, headers: map[string, string]) -> void
      + writes a response to the client
      # http
  std.url
    std.url.decode
      @ (s: string) -> string
      + returns the percent-decoded string
      # url
    std.url.parse_query
      @ (query: string) -> map[string, string]
      + parses a URL query string into a map
      # url

web
  web.new_router
    @ () -> router_state
    + creates an empty router
    # construction
  web.route
    @ (router: router_state, method: string, pattern: string, handler: route_handler) -> router_state
    + registers a route, supporting path parameters like "/users/:id"
    # routing
  web.use
    @ (router: router_state, mw: middleware) -> router_state
    + appends middleware applied to every request in registration order
    # middleware
  web.group
    @ (router: router_state, prefix: string, mw: list[middleware]) -> router_group
    + creates a route group sharing a path prefix and middleware
    # routing
  web.group_route
    @ (g: router_group, method: string, pattern: string, handler: route_handler) -> router_group
    + registers a route inside a group
    # routing
  web.match
    @ (router: router_state, method: string, path: string) -> result[route_match, string]
    + returns the handler and extracted path parameters for a request
    - returns error when no route matches
    # routing
  web.dispatch
    @ (router: router_state, req: http_request) -> http_response
    + runs middleware and the matched handler, producing a response
    + returns a 404 response when no route matches
    # execution
    -> std.url.decode
    -> std.url.parse_query
  web.serve
    @ (router: router_state, addr: string) -> result[void, string]
    + starts the HTTP server bound to this router
    - returns error when the address cannot be bound
    # server
    -> std.http.serve
    -> std.http.parse_request
    -> std.http.write_response
