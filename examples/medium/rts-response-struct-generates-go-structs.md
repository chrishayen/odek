# Requirement: "a library that generates typed record definitions from JSON server responses"

Inspects a JSON sample and emits a language-agnostic struct schema description. The emitter is pluggable so callers can render to any target syntax.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value (object, array, string, number, bool, null)
      - returns error on malformed input
      # parsing
  std.strings
    std.strings.to_pascal_case
      @ (s: string) -> string
      + converts snake_case or kebab-case identifiers to PascalCase
      # naming

struct_gen
  struct_gen.infer_schema
    @ (root: json_value, root_name: string) -> schema
    + walks the JSON value and produces a schema with one record per distinct object shape
    + merges array element shapes into a single element type
    ? keys across object instances are unioned; missing keys become optional fields
    # inference
    -> std.strings.to_pascal_case
  struct_gen.unify_types
    @ (a: field_type, b: field_type) -> field_type
    + returns the narrowest type that covers both inputs
    + widens number and null combinations into optional numeric types
    # type_unification
  struct_gen.from_sample
    @ (raw: string, root_name: string) -> result[schema, string]
    + convenience entry point that parses the sample and infers its schema
    - propagates parse errors
    # entry
    -> std.json.parse_value
  struct_gen.render
    @ (s: schema, emit: fn(record) -> string) -> string
    + renders each record via the supplied emitter and joins the results in dependency order
    # rendering
