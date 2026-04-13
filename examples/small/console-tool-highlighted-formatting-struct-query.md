# Requirement: "a JSON pretty-printer with colorized output and path-based querying"

Two pieces: colorize a parsed JSON document for display, and extract a value at a dotted path.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a generic value tree
      - returns error on invalid JSON
      # serialization

json_view
  json_view.colorize
    @ (doc: json_value, indent: i32) -> string
    + renders a JSON value with ANSI color codes for keys, strings, numbers, bools, and null
    + respects the given indent width for nested objects and arrays
    # rendering
  json_view.query_path
    @ (doc: json_value, path: string) -> result[json_value, string]
    + returns the subvalue at a dotted path like "user.addresses.0.city"
    - returns error when any path segment is missing
    - returns error when indexing a non-array with a numeric segment
    # query
    -> std.json.parse
