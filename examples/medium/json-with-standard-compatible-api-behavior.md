# Requirement: "a JSON parsing and encoding library with an API compatible with a standard library"

Parse produces a tagged value tree; encode serializes it back. A handful of accessor functions let callers inspect the tree.

std: (all units exist)

jzon
  jzon.parse
    @ (raw: string) -> result[json_value, string]
    + parses arbitrary JSON into a tagged value tree
    - returns error on unterminated strings
    - returns error on trailing garbage
    # parsing
  jzon.encode
    @ (value: json_value) -> string
    + returns the canonical JSON encoding of the value
    + escapes control characters in strings
    # encoding
  jzon.type_of
    @ (value: json_value) -> string
    + returns one of "null", "bool", "number", "string", "array", "object"
    # inspection
  jzon.get_field
    @ (value: json_value, name: string) -> optional[json_value]
    + returns the named field of an object
    - returns none when the value is not an object
    # inspection
  jzon.index
    @ (value: json_value, i: i32) -> optional[json_value]
    + returns the i-th element of an array
    - returns none when the value is not an array or i is out of range
    # inspection
  jzon.as_string
    @ (value: json_value) -> optional[string]
    + returns the underlying string when value is a JSON string
    # inspection
  jzon.as_number
    @ (value: json_value) -> optional[f64]
    + returns the underlying number when value is a JSON number
    # inspection
