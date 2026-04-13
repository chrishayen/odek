# Requirement: "a REST API routing library backed by a relational store"

Registers HTTP handlers by method and path pattern and dispatches incoming requests. The store is abstracted behind a repository interface.

std
  std.http
    std.http.parse_request
      @ (raw: string) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed request lines
      # http
    std.http.render_response
      @ (status: i32, body: string) -> string
      + returns a wire-format response with appropriate headers
      # http

router
  router.new
    @ () -> router_state
    + returns an empty router
    # construction
  router.register
    @ (state: router_state, method: string, pattern: string, handler: handler_fn) -> router_state
    + adds a handler for the given method and path pattern
    + supports placeholders like "/items/{id}"
    # registration
  router.match
    @ (state: router_state, method: string, path: string) -> optional[route_match]
    + returns the handler and extracted path parameters when a route matches
    - returns none when no pattern matches
    # dispatch
  router.handle
    @ (state: router_state, raw_request: string) -> string
    + parses the request, dispatches to the matching handler, and renders the response
    - returns a 404 response when no route matches
    - returns a 400 response when the request is malformed
    # dispatch
    -> std.http.parse_request
    -> std.http.render_response
