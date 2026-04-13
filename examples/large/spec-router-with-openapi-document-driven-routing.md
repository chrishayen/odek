# Requirement: "a spec-first request routing framework driven by an OpenAPI-style document"

Loads a spec, matches incoming requests to operations, validates parameters and bodies against the declared schema, and dispatches to handlers.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a dynamic value tree
      - returns error on malformed input
      # serialization
    std.json.get_field
      @ (v: json_value, key: string) -> optional[json_value]
      + returns the named field of an object value
      # serialization
    std.json.as_string
      @ (v: json_value) -> optional[string]
      + returns the string content when the value is a string
      # serialization
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[json_value, string]
      + parses a YAML document into the same value tree shape as JSON
      - returns error on malformed input
      # serialization
  std.url
    std.url.parse_query
      @ (query: string) -> map[string, string]
      + parses a "?a=1&b=2" query string into a map
      # url
    std.url.split_path
      @ (path: string) -> list[string]
      + splits a URL path into its slash-separated segments
      # url

spec_router
  spec_router.load_spec
    @ (source: string) -> result[api_spec, string]
    + parses a spec document in JSON or YAML into an internal model
    - returns error on parse failure or missing required sections
    - returns error when the spec declares no operations
    # spec_loading
    -> std.json.parse
    -> std.yaml.parse
  spec_router.compile_routes
    @ (spec: api_spec) -> route_table
    + builds a fast lookup table from (method, path_template) to operation id
    ? path templates with "{param}" segments are compiled into matchers
    # routing
  spec_router.match
    @ (table: route_table, method: string, path: string) -> optional[route_match]
    + returns the operation id and extracted path parameters for a request
    # routing
    -> std.url.split_path
  spec_router.validate_parameters
    @ (spec: api_spec, op_id: string, path_params: map[string, string], query: string) -> result[map[string, string], string]
    + checks required params, coerces declared types, and returns the merged parameter map
    - returns error when a required parameter is missing
    - returns error when a value does not match its declared type
    # validation
    -> std.url.parse_query
  spec_router.validate_body
    @ (spec: api_spec, op_id: string, body: string) -> result[json_value, string]
    + validates the request body against the operation's schema
    - returns error when required fields are missing
    - returns error when a field has the wrong type
    # validation
    -> std.json.parse
    -> std.json.get_field
    -> std.json.as_string
  spec_router.register_handler
    @ (r: router_state, op_id: string, handler_id: string) -> result[void, string]
    + binds a handler id to an operation defined in the spec
    - returns error when the operation id is not in the spec
    # dispatch
  spec_router.dispatch
    @ (r: router_state, method: string, path: string, query: string, body: string) -> result[dispatch_result, string]
    + matches, validates, and returns the handler id plus parameter map to invoke
    - returns dispatch_result with status 404 when no operation matches
    - returns dispatch_result with status 400 on validation errors
    - returns dispatch_result with status 501 when no handler is registered
    # dispatch
  spec_router.list_operations
    @ (spec: api_spec) -> list[operation_summary]
    + returns method, path, operation id, and summary for every declared operation
    # introspection
