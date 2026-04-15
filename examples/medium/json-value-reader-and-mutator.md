# Requirement: "a library for reading and mutating unstructured JSON documents"

Exposes a single json_value type with typed accessors and path-based get/set, so callers do not need a schema.

std: (all units exist)

json_value
  json_value.parse
    fn (raw: string) -> result[json_value, string]
    + parses any valid JSON document into a json_value
    - returns error with position on malformed input
    # parsing
  json_value.encode
    fn (value: json_value) -> string
    + returns the canonical JSON text for value
    # encoding
  json_value.kind
    fn (value: json_value) -> string
    + returns one of "null", "bool", "number", "string", "array", "object"
    # introspection
  json_value.get
    fn (value: json_value, path: string) -> optional[json_value]
    + returns the nested value at a dotted path like "a.b.0.c"
    - returns none when any segment is missing
    # navigation
  json_value.set
    fn (value: json_value, path: string, new_value: json_value) -> result[json_value, string]
    + returns a new document with new_value set at path
    - returns error when a path segment type conflicts with new_value
    # mutation
  json_value.as_string
    fn (value: json_value) -> optional[string]
    + returns the string payload when value is a string
    - returns none for any other kind
    # accessors
  json_value.as_i64
    fn (value: json_value) -> optional[i64]
    + returns the integer payload when value is a whole number
    # accessors
  json_value.as_array
    fn (value: json_value) -> optional[list[json_value]]
    + returns the array payload when value is an array
    # accessors
