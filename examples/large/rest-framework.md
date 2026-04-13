# Requirement: "a REST web service framework"

A routing and middleware framework for HTTP-style request handling with content negotiation and error mapping.

std
  std.http
    std.http.parse_request_line
      @ (line: string) -> result[tuple[string, string, string], string]
      + returns (method, path, version) from a request line
      - returns error when the line does not have three whitespace-separated parts
      # http
    std.http.parse_headers
      @ (raw: string) -> result[map[string, string], string]
      + parses CRLF-separated headers into a lowercase-keyed map
      - returns error on malformed header lines
      # http
  std.url
    std.url.decode
      @ (encoded: string) -> result[string, string]
      + decodes percent-encoded sequences
      - returns error on invalid %XX sequences
      # url
    std.url.parse_query
      @ (query: string) -> map[string, string]
      + parses a "?a=1&b=2" string into key-value pairs
      # url
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

rest_framework
  rest_framework.new_router
    @ () -> router_state
    + creates a router with no routes and no middleware
    # construction
  rest_framework.route
    @ (router: router_state, method: string, pattern: string, handler: fn(request) -> response) -> router_state
    + registers a handler for the given method and path pattern
    ? pattern supports ":name" segments that become path parameters
    # routing
  rest_framework.use
    @ (router: router_state, mw: fn(request, fn(request) -> response) -> response) -> router_state
    + appends a middleware applied to every matching route
    # middleware
  rest_framework.match_route
    @ (router: router_state, method: string, path: string) -> result[tuple[handler, map[string, string]], string]
    + returns the handler and extracted path parameters for the request
    - returns error with status 405 when path matches but method does not
    - returns error with status 404 when no pattern matches
    # routing
    -> std.url.decode
  rest_framework.dispatch
    @ (router: router_state, req: request) -> response
    + runs middleware chain and handler, returning the final response
    + returns a 500 response when a handler throws
    # dispatch
  rest_framework.parse_request
    @ (raw: bytes) -> result[request, string]
    + parses a raw HTTP request into method, path, headers, query, and body
    - returns error on malformed request line or headers
    # parsing
    -> std.http.parse_request_line
    -> std.http.parse_headers
    -> std.url.parse_query
  rest_framework.json_response
    @ (status: i32, body: map[string, string]) -> response
    + builds a response with Content-Type application/json and encoded body
    # responses
    -> std.json.encode_object
  rest_framework.read_json_body
    @ (req: request) -> result[map[string, string], string]
    + parses the request body as a JSON object when Content-Type is application/json
    - returns error when the body is not JSON or content type does not match
    # parsing
    -> std.json.parse_object
  rest_framework.error_response
    @ (status: i32, message: string) -> response
    + builds a standard error response with status and message fields
    # responses
    -> std.json.encode_object
  rest_framework.negotiate_content
    @ (accept_header: string, offered: list[string]) -> optional[string]
    + returns the best matching content type from the offered list
    - returns none when nothing matches
    # content_negotiation
