# Requirement: "a library for arbitrary transformation of JSON documents driven by a declarative spec"

The spec describes shift, default, delete, and concat operations. A single apply function walks the spec and produces a new document.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses any JSON document into a tagged value tree
      - returns error on invalid JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a json_value back to a JSON string
      # serialization
    std.json.get_path
      @ (root: json_value, path: string) -> optional[json_value]
      + returns the value at a dotted path, or none if missing
      ? path segments are split on "."
      # traversal
    std.json.set_path
      @ (root: json_value, path: string, value: json_value) -> json_value
      + returns a new document with the value written at the dotted path
      + creates intermediate objects as needed
      # traversal
    std.json.delete_path
      @ (root: json_value, path: string) -> json_value
      + returns a new document with the key at the dotted path removed
      # traversal

jsontransform
  jsontransform.compile
    @ (spec_source: string) -> result[transform_spec, string]
    + parses a spec JSON document describing one or more operations
    - returns error when an operation kind is unknown
    # compile
    -> std.json.parse
  jsontransform.apply
    @ (spec: transform_spec, input: string) -> result[string, string]
    + applies each operation in order and returns the rewritten document
    - returns error when the input is not valid JSON
    # execution
    -> std.json.parse
    -> std.json.encode
  jsontransform.op_shift
    @ (spec: transform_spec, doc: json_value) -> json_value
    + moves values from source paths to destination paths
    # operation
    -> std.json.get_path
    -> std.json.set_path
  jsontransform.op_default
    @ (spec: transform_spec, doc: json_value) -> json_value
    + writes values at paths that are currently missing
    # operation
    -> std.json.get_path
    -> std.json.set_path
  jsontransform.op_delete
    @ (spec: transform_spec, doc: json_value) -> json_value
    + removes keys listed in the spec
    # operation
    -> std.json.delete_path
  jsontransform.op_concat
    @ (spec: transform_spec, doc: json_value) -> json_value
    + joins string values from source paths into a single destination string
    # operation
    -> std.json.get_path
    -> std.json.set_path
