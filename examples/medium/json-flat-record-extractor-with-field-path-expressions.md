# Requirement: "a library to extract a flat record from a nested JSON document using field path expressions"

Callers declare a mapping from output field names to JSON paths; the library walks the document and produces a flat string map.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_node, string]
      + parses a JSON document into an in-memory tree
      - returns error on malformed input
      # serialization
    std.json.node_at_path
      fn (root: json_node, segments: list[string]) -> optional[json_node]
      + walks the segments (object keys or numeric indices) through the tree
      - returns none when any segment misses
      # traversal
    std.json.node_as_string
      fn (node: json_node) -> optional[string]
      + returns the string form of a scalar node (string, number, bool, null)
      - returns none for object or array nodes
      # serialization
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on the given separator
      # strings

json_extract
  json_extract.compile_mapping
    fn (mapping: map[string,string]) -> map[string, list[string]]
    + splits each dotted path value into segments keyed by the field name
    # mapping
    -> std.strings.split
  json_extract.extract
    fn (raw: string, mapping: map[string, list[string]]) -> result[map[string,string], string]
    + parses the document and returns a flat field-to-string map
    + fields whose paths miss are omitted from the result
    - returns error when the document itself cannot be parsed
    # extraction
    -> std.json.parse
    -> std.json.node_at_path
    -> std.json.node_as_string
