# Requirement: "a library that converts JSON samples to typed record definitions"

Infers a record schema from one or more JSON samples, merging types across samples, and emits a language-agnostic schema description.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a tagged value tree
      - returns error on malformed JSON
      # serialization

json_to_record
  json_to_record.infer_type
    fn (value: json_value) -> inferred_type
    + returns a scalar type for primitives
    + returns a list type whose element is the merged type of all elements
    + returns a record type with inferred field types for objects
    ? null maps to optional[unknown] so merges can refine it
    # inference
  json_to_record.merge_types
    fn (a: inferred_type, b: inferred_type) -> inferred_type
    + returns the unified type when both are compatible
    + widens integer+float to float
    + wraps in optional when one side is null
    - returns "any" when types cannot be unified
    # inference
  json_to_record.from_samples
    fn (samples: list[string]) -> result[inferred_type, string]
    + parses each sample and merges their inferred types
    - returns error when any sample fails to parse
    # pipeline
    -> std.json.parse_value
    -> json_to_record.infer_type
    -> json_to_record.merge_types
  json_to_record.to_schema_text
    fn (schema: inferred_type, root_name: string) -> string
    + returns a textual schema description with nested record definitions
    + generates unique sub-record names from field paths
    # emission
