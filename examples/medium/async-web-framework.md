# Requirement: "a lightweight asynchronous web framework"

Routes HTTP requests to handlers through a middleware chain and builds responses.

std: (all units exist)

webfw
  webfw.new_app
    @ () -> app
    + creates an app with no routes and no middleware
    # construction
  webfw.route
    @ (a: app, method: string, pattern: string, handler: http_handler) -> result[app, string]
    + registers handler for method and pattern, supporting ":name" path parameters
    - returns error when method is empty or pattern already registered for that method
    # routing
  webfw.use
    @ (a: app, mw: middleware) -> app
    + appends a middleware to the chain run before handlers
    # middleware
  webfw.match_route
    @ (a: app, method: string, path: string) -> result[route_match, string]
    + returns the matching handler and extracted path parameters
    - returns error when no route matches
    # routing
  webfw.parse_request
    @ (raw: string) -> result[request, string]
    + parses an HTTP/1.1 request line, headers, and body
    - returns error on malformed request line or bad headers
    # http_protocol
  webfw.write_response
    @ (resp: response) -> string
    + serializes status line, headers, and body into wire format
    # http_protocol
  webfw.handle
    @ (a: app, req: request) -> response
    + runs the middleware chain then the matched handler and returns a response
    + returns a 404 response when no route matches
    + returns a 500 response when a handler raises
    # dispatch
    -> webfw.match_route
