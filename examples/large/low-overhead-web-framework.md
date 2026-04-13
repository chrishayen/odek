# Requirement: "a low-overhead web framework"

A minimal router plus request/response types layered on top of an http server primitive. Handlers are plain functions; middlewares compose by wrapping handlers.

std
  std.http
    std.http.listen
      @ (addr: string, handler: http_server_handler) -> result[void, string]
      + binds a listener on addr and dispatches parsed requests to handler
      - returns error when the address is already in use
      # network
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses a raw http/1.1 request
      - returns error on protocol violation
      # protocol
    std.http.format_response
      @ (resp: http_response) -> bytes
      + serializes an http response to wire bytes
      # protocol
  std.json
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a generic json value
      # serialization
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses a json document
      - returns error on invalid input
      # serialization

web
  web.new_app
    @ () -> web_app
    + creates an app with an empty route table
    # construction
  web.route
    @ (app: web_app, method: string, path: string, handler: route_handler) -> web_app
    + registers a handler for method+path, supporting :param segments
    # routing
  web.use
    @ (app: web_app, middleware: middleware_fn) -> web_app
    + pushes a middleware onto the chain wrapping every handler
    # middleware
  web.match_route
    @ (app: web_app, method: string, path: string) -> optional[tuple[route_handler, map[string, string]]]
    + returns the matched handler and extracted path params
    - returns none when no route matches
    # routing
  web.dispatch
    @ (app: web_app, req: http_request) -> http_response
    + runs the matched handler through the middleware chain and returns the response
    + returns 404 when no route matches
    + returns 405 when the path matches but the method does not
    # request_dispatch
  web.json_response
    @ (status: i32, body: json_value) -> http_response
    + builds a response with Content-Type: application/json
    # helpers
    -> std.json.encode_value
  web.read_json_body
    @ (req: http_request) -> result[json_value, string]
    + parses the request body as json
    - returns error on malformed json
    # helpers
    -> std.json.parse_value
  web.serve
    @ (app: web_app, addr: string) -> result[void, string]
    + starts the http listener and blocks
    # server
    -> std.http.listen
    -> std.http.parse_request
    -> std.http.format_response
