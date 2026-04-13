# Requirement: "a library to handle and query JSON without defining typed schemas"

Parses JSON into a dynamic value tree and queries it with a dotted path, returning typed accessors at the leaves.

std
  std.json
    std.json.parse
      @ (text: string) -> result[json_value, string]
      + parses arbitrary JSON into a tagged value tree
      - returns error on malformed input
      # parsing
    std.json.serialize
      @ (value: json_value) -> string
      + serializes a value tree back to JSON text
      # serialization

jsonic
  jsonic.load
    @ (text: string) -> result[json_value, string]
    + parses JSON text for subsequent querying
    - returns error on malformed input
    # loading
    -> std.json.parse
  jsonic.get
    @ (value: json_value, path: string) -> optional[json_value]
    + walks a dotted path like "a.b[2].c" and returns the node or none
    ? bracket notation selects array elements
    # querying
  jsonic.get_string
    @ (value: json_value, path: string) -> result[string, string]
    + returns the string at path
    - returns error when the path is missing or the node is not a string
    # querying
  jsonic.get_i64
    @ (value: json_value, path: string) -> result[i64, string]
    + returns the integer at path
    - returns error when the node is not an integer
    # querying
  jsonic.get_f64
    @ (value: json_value, path: string) -> result[f64, string]
    + returns the float at path
    - returns error when the node is not numeric
    # querying
  jsonic.get_bool
    @ (value: json_value, path: string) -> result[bool, string]
    + returns the boolean at path
    - returns error when the node is not a boolean
    # querying
  jsonic.keys
    @ (value: json_value, path: string) -> result[list[string], string]
    + returns the keys of the object at path
    - returns error when the node is not an object
    # inspection
  jsonic.set
    @ (value: json_value, path: string, replacement: json_value) -> result[json_value, string]
    + returns a new tree with the node at path replaced
    - returns error when the path does not resolve
    # mutation
  jsonic.dump
    @ (value: json_value) -> string
    + serializes the tree back to JSON text
    # serialization
    -> std.json.serialize
