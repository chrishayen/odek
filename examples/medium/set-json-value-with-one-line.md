# Requirement: "set a JSON value at a dotted path in one call"

Given a JSON document and a dotted path, return a new document with the value set at that path. Intermediate objects or arrays are created as needed.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a generic value tree
      - returns error on invalid JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a generic value tree back to JSON text
      # serialization

sjson
  sjson.parse_path
    @ (path: string) -> list[path_segment]
    + splits a dotted path into segments; numeric segments become array indices
    + supports escaping '.' with a backslash inside a segment
    # path_parsing
  sjson.set_string
    @ (doc: string, path: string, value: string) -> result[string, string]
    + returns the JSON document with value set as a string at path
    - returns error when the document is invalid JSON
    # set
    -> std.json.parse
    -> std.json.encode
  sjson.set_number
    @ (doc: string, path: string, value: f64) -> result[string, string]
    + returns the JSON document with value set as a number at path
    - returns error when the document is invalid JSON
    # set
    -> std.json.parse
    -> std.json.encode
  sjson.set_bool
    @ (doc: string, path: string, value: bool) -> result[string, string]
    + returns the JSON document with value set as a bool at path
    # set
    -> std.json.parse
    -> std.json.encode
  sjson.set_raw
    @ (doc: string, path: string, raw_json: string) -> result[string, string]
    + returns the JSON document with raw_json inserted verbatim at path
    - returns error when raw_json is not valid JSON
    # set
    -> std.json.parse
    -> std.json.encode
  sjson.delete
    @ (doc: string, path: string) -> result[string, string]
    + returns the JSON document with the value at path removed
    - returns error when the document is invalid JSON
    # delete
    -> std.json.parse
    -> std.json.encode
