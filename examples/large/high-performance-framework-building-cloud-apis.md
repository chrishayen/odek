# Requirement: "a framework for building cloud apis and web backends"

Resource-oriented API framework where resources expose method handlers, requests flow through middleware, and responses are structured values the caller serializes over a transport.

std
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
  std.text
    std.text.split
      @ (s: string, sep: string) -> list[string]
      + splits a string on a separator
      # text
    std.text.trim
      @ (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # text

cloudapi
  cloudapi.api_new
    @ () -> api_state
    + returns an empty API with no resources
    # construction
  cloudapi.add_resource
    @ (api: api_state, path: string, resource_id: string) -> result[api_state, string]
    + mounts a resource handler at a path pattern
    - returns error on duplicate path
    # resources
    -> std.text.split
  cloudapi.add_middleware
    @ (api: api_state, middleware_id: string) -> api_state
    + appends middleware to the global chain
    # middleware
  cloudapi.add_error_handler
    @ (api: api_state, status: i32, handler_id: string) -> api_state
    + registers an error handler for a specific HTTP status
    # error_handling
  cloudapi.parse_request
    @ (method: string, path: string, headers: map[string, string], body: bytes) -> http_request
    + constructs a canonical request value
    # request_parsing
    -> std.text.trim
  cloudapi.dispatch
    @ (api: api_state, request: http_request) -> http_response
    + routes the request through middleware and the matched resource handler
    + returns a 404 response when no resource matches
    + returns a registered error handler's response on status matches
    # dispatch
  cloudapi.make_response
    @ (status: i32, body: map[string, string]) -> http_response
    + constructs a JSON response
    # responses
    -> std.serialization.encode_json
  cloudapi.parse_body_json
    @ (request: http_request) -> result[map[string, string], string]
    + decodes the request body as JSON
    - returns error on invalid JSON
    # request_parsing
    -> std.serialization.decode_json
  cloudapi.validate_schema
    @ (body: map[string, string], required: list[string]) -> result[void, string]
    + checks that every required field is present and non-empty
    - returns error listing missing fields
    # validation
  cloudapi.make_error_response
    @ (status: i32, message: string) -> http_response
    + constructs a JSON error body with "error" and "message" fields
    # error_handling
    -> std.serialization.encode_json
