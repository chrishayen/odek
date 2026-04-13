# Requirement: "an XML to structured-object converter"

Parses an XML document into a generic nested map structure and reverses the operation.

std
  std.xml
    std.xml.parse
      @ (raw: string) -> result[xml_node, string]
      + parses an XML document into a tree of nodes
      - returns error on malformed XML
      # parsing
    std.xml.serialize
      @ (node: xml_node) -> string
      + serializes an XML node tree back to a textual document
      # serialization

xml2obj
  xml2obj.to_object
    @ (raw: string) -> result[object_value, string]
    + converts an XML document to a nested object where element names map to children
    + repeated sibling elements become lists
    - returns error on invalid XML
    # conversion
    -> std.xml.parse
  xml2obj.to_xml
    @ (obj: object_value, root_name: string) -> result[string, string]
    + serializes a nested object back to XML using root_name as the document element
    - returns error when object structure cannot be represented
    # conversion
    -> std.xml.serialize
