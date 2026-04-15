# Requirement: "a library for JSON payload sanitization and schema validation"

Validates JSON payloads against a declarative schema and sanitizes values (coerce types, trim strings, clamp numbers).

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a tagged value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      fn (v: json_value) -> string
      + encodes a tagged json_value back to a JSON string
      # serialization
  std.strings
    std.strings.trim
      fn (s: string) -> string
      + trims ASCII whitespace from both ends
      # strings

schema
  schema.field_string
    fn (name: string, required: bool, min_len: i32, max_len: i32) -> field_spec
    + builds a string field spec with length bounds
    # spec
  schema.field_integer
    fn (name: string, required: bool, min_val: i64, max_val: i64) -> field_spec
    + builds an integer field spec with inclusive range
    # spec
  schema.field_boolean
    fn (name: string, required: bool) -> field_spec
    + builds a boolean field spec
    # spec
  schema.object
    fn (fields: list[field_spec]) -> object_schema
    + bundles field specs into a top-level schema
    # spec
  schema.validate
    fn (s: object_schema, raw: string) -> result[json_value, list[validation_error]]
    + parses raw and returns it when it conforms
    - returns a list of per-field errors when any constraint fails
    # validation
    -> std.json.parse_value
  schema.sanitize
    fn (s: object_schema, raw: string) -> result[string, list[validation_error]]
    + parses raw, coerces strings (trimmed), clamps numbers to bounds, drops unknown fields, and re-encodes
    - returns errors only for unfixable violations (missing required fields, wrong base type)
    # sanitization
    -> std.json.parse_value
    -> std.json.encode_value
    -> std.strings.trim
