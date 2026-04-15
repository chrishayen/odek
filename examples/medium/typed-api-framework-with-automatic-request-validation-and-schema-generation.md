# Requirement: "a typed API framework with automatic request validation and schema generation"

Registers typed endpoints, validates incoming JSON against their schemas, and emits an OpenAPI-style description.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[value, string]
      + parses JSON into a dynamic value
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      fn (v: value) -> string
      + encodes a dynamic value as JSON
      # serialization

api
  api.new_app
    fn () -> app_state
    + creates an empty app with no endpoints
    # construction
  api.endpoint
    fn (app: app_state, method: string, path: string, input_schema: schema, output_schema: schema, handler: fn(value) -> result[value, api_error]) -> app_state
    + registers a typed endpoint with request and response schemas
    # registration
  api.validate
    fn (sch: schema, v: value) -> result[void, string]
    + checks that v conforms to sch (type, required fields, ranges)
    - returns error describing the first mismatch
    # validation
  api.handle
    fn (app: app_state, method: string, path: string, body: string) -> tuple[i32, string]
    + parses body, validates, runs the handler, validates the response, and returns (status, body)
    - returns (400, error body) on input validation failure
    - returns (404, error body) on unknown route
    - returns handler error status on handler failure
    # dispatch
    -> std.json.parse_value
    -> std.json.encode_value
  api.describe
    fn (app: app_state) -> string
    + emits an OpenAPI 3-compatible JSON description of all registered endpoints
    # schema_export
    -> std.json.encode_value
