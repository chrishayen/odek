# Requirement: "an API framework with OpenAPI schema generation"

A minimal routing layer that lets a caller register typed operations and emit an OpenAPI 3 document describing them. Request parsing and response encoding go through thin std primitives.

std
  std.json
    std.json.encode_value
      @ (value: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
    std.json.parse_value
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from raw request bytes
      - returns error on malformed request line
      # http
    std.http.format_response
      @ (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes a status, headers, and body into wire-format bytes
      # http

api_framework
  api_framework.new_router
    @ () -> router_state
    + returns an empty router with no registered operations
    # construction
  api_framework.register_operation
    @ (state: router_state, method: string, path: string, summary: string) -> router_state
    + adds an operation keyed by (method, path) with its metadata
    - leaves state unchanged if the same method and path are already registered
    # registration
  api_framework.dispatch
    @ (state: router_state, raw: bytes) -> result[bytes, string]
    + parses the request, finds the matching operation, and returns the response bytes
    - returns a 404 response when no operation matches
    # dispatch
    -> std.http.parse_request
    -> std.http.format_response
    -> std.json.parse_value
    -> std.json.encode_value
  api_framework.generate_openapi
    @ (state: router_state) -> string
    + emits an OpenAPI 3 document describing every registered operation
    + returns a document with an empty paths object when no operations are registered
    # documentation
