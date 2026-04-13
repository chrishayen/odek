# Requirement: "a library for parsing and validating HTTP request arguments from query, form, and JSON bodies"

Declares a schema of expected fields with types and required flags, then extracts them from any supported source into a typed map.

std
  std.http
    std.http.parse_query
      @ (query: string) -> map[string, string]
      + parses a urlencoded query string into a string-to-string map
      # http
  std.json
    std.json.parse_value
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

request_args
  request_args.new_schema
    @ () -> schema_state
    + returns an empty schema with no fields
    # construction
  request_args.add_field
    @ (schema: schema_state, name: string, kind: string, required: bool) -> schema_state
    + registers a field with its expected type ("string", "int", "bool") and required flag
    # schema
  request_args.parse_query_args
    @ (schema: schema_state, query: string) -> result[map[string, string], list[string]]
    + returns parsed values when every required field is present and well-typed
    - returns a list of error messages naming missing or malformed fields
    # parsing
    -> std.http.parse_query
  request_args.parse_json_args
    @ (schema: schema_state, body: string) -> result[map[string, string], list[string]]
    + returns parsed values from a JSON body validated against the schema
    - returns a list of error messages naming missing or malformed fields
    # parsing
    -> std.json.parse_value
