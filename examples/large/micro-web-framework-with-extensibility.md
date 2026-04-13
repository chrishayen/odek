# Requirement: "an extensible micro web framework"

Route registration with path parameters, middleware chains, and a request handler dispatcher. The framework operates on in-memory request and response values; the caller integrates with a transport.

std
  std.text
    std.text.split
      @ (s: string, sep: string) -> list[string]
      + splits a string on a separator, returning empty strings for adjacent separators
      # text
    std.text.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # text
  std.serialization
    std.serialization.encode_json
      @ (value: map[string, string]) -> string
      + encodes a string-keyed map as JSON
      # serialization
    std.serialization.decode_json
      @ (raw: string) -> result[map[string, string], string]
      + decodes a JSON object into a string-keyed map
      - returns error on invalid JSON
      # serialization

microweb
  microweb.app_new
    @ () -> app_state
    + returns an app with no routes and no middleware
    # construction
  microweb.use
    @ (app: app_state, middleware_id: string) -> app_state
    + appends middleware to the global chain, run in registration order
    # middleware
  microweb.route
    @ (app: app_state, method: string, path: string, handler_id: string) -> result[app_state, string]
    + registers a handler for a method-path pair, supporting ":param" segments
    - returns error on duplicate method-path registration
    # routing
    -> std.text.split
  microweb.group
    @ (app: app_state, prefix: string) -> route_group
    + returns a route group sharing a path prefix and its own middleware chain
    # routing
  microweb.group_route
    @ (group: route_group, method: string, path: string, handler_id: string) -> result[route_group, string]
    + registers a handler on the group
    - returns error on duplicate registration
    # routing
  microweb.match
    @ (app: app_state, method: string, path: string) -> optional[matched_route]
    + returns the matched handler id and extracted path parameters
    # routing
    -> std.text.split
    -> std.text.starts_with
  microweb.handle
    @ (app: app_state, request: http_request) -> http_response
    + runs middleware then the matched handler, returning a response
    + returns a 404 response when no route matches
    # dispatch
  microweb.json_response
    @ (status: i32, body: map[string, string]) -> http_response
    + constructs a JSON response with the given status and body
    # response_helpers
    -> std.serialization.encode_json
  microweb.parse_json_body
    @ (request: http_request) -> result[map[string, string], string]
    + decodes the request body as a JSON object
    - returns error on invalid JSON
    # request_helpers
    -> std.serialization.decode_json
  microweb.set_header
    @ (response: http_response, name: string, value: string) -> http_response
    + returns a new response with the header added or replaced
    # response_helpers
