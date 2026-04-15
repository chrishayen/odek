# Requirement: "a library that lets you work with XML as if it were JSON-like maps"

Converts XML documents to nested maps and back.

std
  std.xml
    std.xml.parse
      fn (source: bytes) -> result[xml_node, string]
      + parses XML 1.0 into a tree of nodes
      - returns error on malformed markup
      # xml
    std.xml.serialize
      fn (root: xml_node) -> bytes
      + serializes a node tree back to XML
      # xml

xmap
  xmap.to_map
    fn (doc: xml_node) -> json_value
    + converts an XML tree into a nested JSON-like value
    + attributes become keys prefixed with "@"
    + text content becomes a "#text" key when mixed with attributes
    + repeated child elements become arrays
    # conversion
  xmap.from_map
    fn (value: json_value) -> result[xml_node, string]
    + rebuilds an XML tree from the nested representation
    - returns error when the root has multiple top-level keys
    - returns error when an attribute value is not a scalar
    # conversion
  xmap.parse
    fn (source: bytes) -> result[json_value, string]
    + parses XML bytes directly into the map form
    - returns error on malformed markup
    # pipeline
    -> std.xml.parse
  xmap.serialize
    fn (value: json_value) -> result[bytes, string]
    + serializes the map form directly to XML bytes
    - returns error when the structure is invalid
    # pipeline
    -> std.xml.serialize
