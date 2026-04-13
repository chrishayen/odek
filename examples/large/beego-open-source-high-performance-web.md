# Requirement: "a high-performance web framework"

Routing, middleware, request lifecycle, templating, and session state. Networking primitives live in std.

std
  std.http
    std.http.serve
      @ (addr: string, handler: fn(http_request) -> http_response) -> result[void, string]
      + binds addr and dispatches requests to handler
      - returns error when bind fails
      # networking
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses a raw HTTP/1.1 request
      - returns error on malformed request line or headers
      # networking
    std.http.response
      @ (status: i32, headers: map[string, string], body: bytes) -> http_response
      + builds a response with the given status, headers, and body
      # networking
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns full file contents
      - returns error when path cannot be read
      # filesystem
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

web_framework
  web_framework.new_app
    @ () -> app_state
    + returns an empty application with no routes or middleware
    # construction
  web_framework.add_route
    @ (app: app_state, method: string, pattern: string, handler: fn(request_ctx) -> response_value) -> void
    + registers a handler for the given method and path pattern with ":param" segments
    # routing
  web_framework.match_route
    @ (app: app_state, method: string, path: string) -> optional[route_match]
    + returns the matching route and extracted path parameters if any
    # routing
  web_framework.use_middleware
    @ (app: app_state, mw: fn(request_ctx, fn(request_ctx) -> response_value) -> response_value) -> void
    + appends a middleware that wraps every subsequent request
    # middleware
  web_framework.dispatch
    @ (app: app_state, req: http_request) -> http_response
    + runs the middleware chain and the matched handler, returning a full response
    + returns a 404 response when no route matches
    # dispatch
    -> std.http.response
  web_framework.run
    @ (app: app_state, addr: string) -> result[void, string]
    + starts the HTTP server bound to addr
    - returns error when bind fails
    # lifecycle
    -> std.http.serve
    -> std.http.parse_request
  web_framework.parse_query
    @ (query: string) -> map[string, string]
    + parses "a=1&b=2" into a map, url-decoding keys and values
    # request
  web_framework.parse_form
    @ (body: bytes, content_type: string) -> result[map[string, string], string]
    + parses a form-urlencoded request body
    - returns error on unsupported content type
    # request
  web_framework.render_template
    @ (path: string, data: map[string, string]) -> result[string, string]
    + loads the template at path and substitutes {{key}} placeholders
    - returns error when the template cannot be read
    # templating
    -> std.fs.read_all
  web_framework.new_session_store
    @ () -> session_store
    + returns an empty in-memory session store
    # sessions
  web_framework.session_start
    @ (store: session_store, req: http_request) -> session_state
    + returns an existing session keyed by cookie or creates a new one
    # sessions
    -> std.crypto.random_bytes
  web_framework.session_save
    @ (store: session_store, sess: session_state, resp: http_response) -> http_response
    + writes the session cookie into the response
    # sessions
