# Requirement: "a serializer and deserializer specialized for XML"

Converts between structured values and XML text via a small schema of field-to-element mappings.

std
  std.xml
    std.xml.parse
      fn (raw: string) -> result[xml_node, string]
      + parses an XML document into a tree of nodes
      - returns error on malformed XML
      # parsing
    std.xml.serialize
      fn (node: xml_node) -> string
      + serializes an XML node tree back to a textual document
      # serialization

xserde
  xserde.serialize
    fn (value: record_value, schema: xml_schema) -> result[string, string]
    + serializes a record to XML using the schema to map fields to elements and attributes
    - returns error when the value does not match the schema
    # serialization
    -> std.xml.serialize
  xserde.deserialize
    fn (raw: string, schema: xml_schema) -> result[record_value, string]
    + parses XML into a record driven by the schema
    - returns error on malformed XML or missing required fields
    # deserialization
    -> std.xml.parse
  xserde.schema_from_fields
    fn (fields: list[field_descriptor]) -> xml_schema
    + builds an xml_schema from a list of field descriptors
    # schema
