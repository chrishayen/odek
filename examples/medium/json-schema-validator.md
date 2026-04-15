# Requirement: "a JSON Schema validator"

Compiles a JSON Schema document into an internal form once and then validates instance documents against it, returning the full list of violations with paths.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + returns a tagged union for null, bool, number, string, array, object
      - returns error on invalid JSON
      # serialization
  std.regex
    std.regex.match
      fn (pattern: string, input: string) -> bool
      + returns true when input matches the pattern
      # regex

json_schema
  json_schema.compile
    fn (schema: json_value) -> result[compiled_schema, string]
    + returns a compiled schema with keyword handlers resolved
    - returns error when a keyword is unknown or has the wrong type
    # compilation
  json_schema.validate
    fn (schema: compiled_schema, instance: json_value) -> list[violation]
    + returns an empty list when the instance satisfies the schema
    - returns one violation per failing keyword with a JSON pointer path
    # validation
    -> std.regex.match
  json_schema.validate_string
    fn (schema: compiled_schema, raw: string) -> result[list[violation], string]
    + convenience that parses raw before validating
    - returns error when raw is not valid JSON
    # validation
    -> std.json.parse
  json_schema.format_violation
    fn (v: violation) -> string
    + returns a one-line human-readable description of a violation
    # reporting
