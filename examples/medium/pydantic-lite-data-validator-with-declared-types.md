# Requirement: "a data validation library driven by declared types"

Define a model as a set of typed fields, then validate and coerce raw input against it.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses JSON text into a dynamic value
      - returns error on malformed input
      # serialization

pydantic_lite
  pydantic_lite.field
    @ (name: string, dtype: string, required: bool) -> field_spec
    + creates a field spec with name, type tag, and required flag
    # schema
  pydantic_lite.with_default
    @ (spec: field_spec, default: json_value) -> field_spec
    + attaches a default value used when the field is absent
    # schema
  pydantic_lite.with_min_length
    @ (spec: field_spec, min: u32) -> field_spec
    + attaches a minimum length constraint for strings and lists
    # schema
  pydantic_lite.with_range
    @ (spec: field_spec, min: f64, max: f64) -> field_spec
    + attaches an inclusive numeric range constraint
    # schema
  pydantic_lite.model
    @ (name: string, fields: list[field_spec]) -> model_spec
    + assembles fields into a named model
    # schema
  pydantic_lite.validate
    @ (model: model_spec, raw: map[string, json_value]) -> result[map[string, json_value], list[validation_error]]
    + coerces and validates a raw map, returning a typed value map
    - returns errors when a required field is missing
    - returns errors when a field fails its type or constraint check
    # validation
  pydantic_lite.parse_json
    @ (model: model_spec, source: string) -> result[map[string, json_value], list[validation_error]]
    + parses JSON and validates it against the model
    - returns a parse error when the JSON is malformed or not an object
    # parsing
    -> std.json.parse
