# Requirement: "a library for processing OpenAPI spec files"

Loads an OpenAPI document, resolves internal $ref pointers, and exposes a traversable model.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file into memory
      - returns error when the file is missing
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged tree
      - returns error on malformed JSON
      # serialization
    std.json.get_path
      @ (value: json_value, pointer: string) -> optional[json_value]
      + resolves a JSON Pointer (RFC 6901) against the document
      # serialization
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[json_value, string]
      + parses YAML into the same tagged tree shape as JSON
      - returns error on malformed YAML
      # serialization

openapi
  openapi.load
    @ (path: string) -> result[spec, string]
    + reads a file and parses either YAML or JSON based on extension
    - returns error when the file is absent or the format is unrecognized
    # loading
    -> std.fs.read_all
    -> std.json.parse
    -> std.yaml.parse
  openapi.version
    @ (spec: spec) -> string
    + returns the "openapi" field value
    # introspection
  openapi.resolve_refs
    @ (spec: spec) -> result[spec, string]
    + substitutes each internal $ref with the referenced subtree
    - returns error on dangling or cyclic references
    # reference_resolution
    -> std.json.get_path
  openapi.list_paths
    @ (spec: spec) -> list[string]
    + returns the route templates under the "paths" object
    # traversal
  openapi.operations_for_path
    @ (spec: spec, path: string) -> list[tuple[string, operation]]
    + returns (method, operation) pairs for the given route
    # traversal
  openapi.parameters_for_operation
    @ (op: operation) -> list[parameter]
    + returns the declared parameters for an operation
    # traversal
