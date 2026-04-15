# Requirement: "nullable primitive wrappers that can be serialized to and from json"

Optional-string, optional-int, optional-float, optional-bool with null-aware json round-trip.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses a json document into a generic value tree
      - returns error on malformed input
      # serialization
    std.json.encode_value
      fn (v: json_value) -> string
      + encodes a generic value tree as a json document
      # serialization

nullable
  nullable.string_to_json
    fn (v: optional[string]) -> string
    + returns the json string form, or "null" when absent
    # marshalling
    -> std.json.encode_value
  nullable.string_from_json
    fn (raw: string) -> result[optional[string], string]
    + returns none when the document is null, some(s) when it is a string
    - returns error when the document is neither null nor a string
    # unmarshalling
    -> std.json.parse_value
  nullable.i64_to_json
    fn (v: optional[i64]) -> string
    + returns the json number form, or "null" when absent
    # marshalling
    -> std.json.encode_value
  nullable.i64_from_json
    fn (raw: string) -> result[optional[i64], string]
    + returns none when the document is null, some(n) when it is an integer
    - returns error when the document is neither null nor an integer
    # unmarshalling
    -> std.json.parse_value
  nullable.f64_to_json
    fn (v: optional[f64]) -> string
    + returns the json number form, or "null" when absent
    # marshalling
    -> std.json.encode_value
  nullable.f64_from_json
    fn (raw: string) -> result[optional[f64], string]
    + returns none when the document is null, some(x) when it is a number
    - returns error when the document is neither null nor a number
    # unmarshalling
    -> std.json.parse_value
  nullable.bool_to_json
    fn (v: optional[bool]) -> string
    + returns the json boolean form, or "null" when absent
    # marshalling
    -> std.json.encode_value
  nullable.bool_from_json
    fn (raw: string) -> result[optional[bool], string]
    + returns none when the document is null, some(b) when it is a boolean
    - returns error when the document is neither null nor a boolean
    # unmarshalling
    -> std.json.parse_value
