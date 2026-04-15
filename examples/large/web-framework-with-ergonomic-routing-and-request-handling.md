# Requirement: "a web framework focused on ergonomic routing and request handling"

Route registration with typed path parameters, request parsing, response building, and a minimal middleware chain.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses a request line, headers, and body
      - returns error on malformed start line
      # http_parsing
    std.http.encode_response
      fn (status: u16, headers: map[string, string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # http_encoding
  std.url
    std.url.decode
      fn (s: string) -> string
      + returns the percent-decoded string
      # url
    std.url.split_query
      fn (query: string) -> map[string, string]
      + returns the parsed key-value pairs
      # url
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

framework
  framework.new_app
    fn () -> app_state
    + creates an app with no routes and an empty middleware chain
    # construction
  framework.route
    fn (app: app_state, method: string, pattern: string, handler_id: string) -> app_state
    + registers a route; pattern may contain ":name" segments
    - returns unchanged app when pattern does not start with "/"
    # routing
  framework.match_route
    fn (app: app_state, method: string, path: string) -> optional[route_match]
    + returns the matched handler id and extracted path parameters
    # routing
  framework.parse_form
    fn (body: string) -> map[string, string]
    + returns form-encoded fields as a string-to-string map
    # extraction
    -> std.url.decode
    -> std.url.split_query
  framework.parse_json_body
    fn (body: string) -> result[map[string, string], string]
    + returns the parsed JSON body
    - returns error on invalid JSON
    # extraction
    -> std.json.parse_object
  framework.respond_text
    fn (status: u16, body: string) -> bytes
    + returns the wire bytes for a plain text response
    # response_building
    -> std.http.encode_response
  framework.respond_json
    fn (status: u16, obj: map[string, string]) -> bytes
    + returns the wire bytes for a JSON response
    # response_building
    -> std.json.encode_object
    -> std.http.encode_response
  framework.use
    fn (app: app_state, middleware_id: string) -> app_state
    + appends a middleware to the chain
    # middleware
  framework.handle
    fn (app: app_state, raw: bytes) -> result[tuple[string, route_match], string]
    + parses the request and returns (handler_id, match) after running the middleware chain
    - returns error when no route matches
    # dispatch
    -> std.http.parse_request
