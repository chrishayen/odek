# Requirement: "a JSON schema expression matching library for requests and responses"

Compiles a small schema expression language and matches JSON payloads against it. Meant for HTTP request/response shape checks.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses JSON text into a json_value tree
      - returns error on malformed input
      # parsing

schema_match
  schema_match.compile
    @ (expression: string) -> result[schema_program, string]
    + returns a compiled schema program supporting type tags, optional fields, and nested objects
    - returns error on unknown type tag
    # compilation
  schema_match.match_value
    @ (program: schema_program, value: json_value) -> result[void, list[string]]
    + returns ok when value conforms to the schema
    - returns a list of path-qualified mismatch messages otherwise
    # matching
  schema_match.match_json
    @ (expression: string, raw_json: string) -> result[void, list[string]]
    + convenience that compiles the schema and matches a raw JSON string
    - returns schema compile errors as a single-element list
    # convenience
    -> std.json.parse
