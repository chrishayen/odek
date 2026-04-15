# Requirement: "generate a typed record definition from an XML or JSON document"

Parses either source format, infers a nested record schema, and emits a language-agnostic struct definition.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a generic value tree
      - returns error on malformed input
      # serialization
  std.xml
    std.xml.parse_document
      fn (raw: string) -> result[xml_node, string]
      + parses an XML document into an element tree
      - returns error on malformed XML
      # serialization

xj2struct
  xj2struct.infer_from_json
    fn (value: json_value, root_name: string) -> struct_schema
    + infers scalar types from JSON primitives
    + infers nested records from JSON objects
    + infers list element types from JSON arrays, unifying heterogeneous elements
    # inference
    -> std.json.parse_value
  xj2struct.infer_from_xml
    fn (root: xml_node, root_name: string) -> struct_schema
    + treats repeated child elements with the same tag as a list field
    + treats attributes as scalar fields on the enclosing record
    + treats text content as a dedicated "text" field when mixed with children
    # inference
    -> std.xml.parse_document
  xj2struct.render_schema
    fn (schema: struct_schema) -> string
    + emits nested record definitions in declaration order
    + uses TitleCase names for records and fields
    + emits list fields with the element type annotated
    # rendering
