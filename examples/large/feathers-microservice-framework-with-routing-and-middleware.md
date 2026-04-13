# Requirement: "a microservice framework with routing, middleware, and service hooks"

A minimal microservice framework: services register typed methods, middleware wraps method calls, and a router dispatches requests by path.

std
  std.json
    std.json.parse_object
      @ (raw: bytes) -> result[map[string,string], string]
      + parses a JSON object
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string,string]) -> bytes
      + encodes an object to JSON bytes
      # serialization
  std.http
    std.http.listen
      @ (port: i32, handler: http_handler_fn) -> result[void, string]
      + starts an HTTP listener on the given port
      - returns error when the port is in use
      # http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses a raw HTTP request
      - returns error on malformed header
      # http
    std.http.build_response
      @ (status: i32, body: bytes) -> bytes
      + builds an HTTP response
      # http

feathers
  feathers.new_app
    @ () -> app_state
    + creates an empty application
    # construction
  feathers.register_service
    @ (app: app_state, path: string, service: service_handle) -> app_state
    + mounts a service at the given path
    # services
  feathers.new_service
    @ (find: handler_fn, get: handler_fn, create: handler_fn, update: handler_fn, remove: handler_fn) -> service_handle
    + builds a CRUD service with the given method handlers
    ? any handler may be null, in which case the method returns 405
    # services
  feathers.use
    @ (app: app_state, middleware: middleware_fn) -> app_state
    + appends a middleware function to the global chain
    # middleware
  feathers.service_use
    @ (app: app_state, path: string, middleware: middleware_fn) -> app_state
    + appends a middleware scoped to one service
    # middleware
  feathers.dispatch
    @ (app: app_state, method: string, path: string, body: bytes) -> result[bytes, string]
    + routes the request to the matching service method through the middleware chain
    - returns error when no service matches the path
    - returns error when the service does not implement the method
    # routing
    -> std.json.parse_object
    -> std.json.encode_object
  feathers.emit
    @ (app: app_state, event: string, payload: bytes) -> app_state
    + invokes all listeners registered for the event
    # events
  feathers.on
    @ (app: app_state, event: string, listener: listener_fn) -> app_state
    + registers a listener for an event name
    # events
  feathers.serve
    @ (app: app_state, port: i32) -> result[void, string]
    + starts an HTTP listener and dispatches incoming requests
    # transport
    -> std.http.listen
    -> std.http.parse_request
    -> std.http.build_response
    -> feathers.dispatch
  feathers.not_found_handler
    @ (path: string) -> bytes
    + returns a 404 JSON body describing the unknown path
    # routing
    -> std.json.encode_object
