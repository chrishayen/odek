# Requirement: "convert a JSON document into a Protobuf schema definition"

Infers a message schema from a JSON value tree and renders it as a .proto source string.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a generic value tree
      - returns error on malformed input
      # serialization

json2proto
  json2proto.infer_schema
    @ (value: json_value, root_name: string) -> result[proto_schema, string]
    + infers scalar types int32, double, bool, and string from JSON primitives
    + infers nested message types for JSON objects
    + infers repeated fields from JSON arrays, using the element type
    + unifies heterogeneous arrays to the most permissive matching type
    - returns error when the root is not a JSON object
    # inference
    -> std.json.parse_value
  json2proto.render_schema
    @ (schema: proto_schema) -> string
    + emits a proto3 source document with numbered fields starting at 1
    + emits nested messages in declaration order
    + escapes reserved identifiers by appending an underscore
    # rendering
  json2proto.convert
    @ (raw_json: string, root_name: string) -> result[string, string]
    + parses, infers, and renders in one call
    - returns error when parsing or inference fails
    # pipeline
    -> json2proto.infer_schema
    -> json2proto.render_schema
