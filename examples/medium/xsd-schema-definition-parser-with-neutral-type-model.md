# Requirement: "an XML schema definition parser that produces a neutral type model"

Parses XSD documents into a language-neutral schema model suitable for downstream code generation.

std
  std.xml
    std.xml.parse
      @ (raw: string) -> result[xml_node, string]
      + parses an XML document into a tree of nodes
      - returns error on malformed XML
      # parsing
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as text
      - returns error when the path cannot be opened
      # filesystem

xsd
  xsd.parse
    @ (raw: string) -> result[schema_model, string]
    + parses an XSD document into a neutral schema model of types and elements
    - returns error on malformed XML or unsupported constructs
    # parsing
    -> std.xml.parse
  xsd.parse_file
    @ (path: string) -> result[schema_model, string]
    + loads an XSD file and parses it
    - returns error when the file is missing
    # parsing
    -> std.fs.read_all
    -> xsd.parse
  xsd.resolve_imports
    @ (model: schema_model, base_dir: string) -> result[schema_model, string]
    + recursively loads imported and included schemas relative to base_dir
    - returns error on circular imports or missing files
    # resolution
    -> xsd.parse_file
  xsd.list_types
    @ (model: schema_model) -> list[type_descriptor]
    + returns all complex and simple types defined in the schema
    # inspection
  xsd.list_elements
    @ (model: schema_model) -> list[element_descriptor]
    + returns all top-level element declarations in the schema
    # inspection
  xsd.type_of_element
    @ (model: schema_model, name: string) -> optional[type_descriptor]
    + looks up the type bound to a named element
    - returns none when the element is not defined
    # inspection
