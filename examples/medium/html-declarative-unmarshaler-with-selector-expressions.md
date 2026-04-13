# Requirement: "declarative unmarshalling of HTML into schema values using selector expressions"

A schema maps field names to CSS-like selectors. Selection primitives live in std; the project layer walks the schema and collects values.

std
  std.html
    std.html.parse
      @ (source: string) -> result[html_doc, string]
      + parses an HTML document into an opaque tree
      - returns error on severely malformed input
      # parsing
    std.html.select_all
      @ (doc: html_doc, selector: string) -> result[list[html_node], string]
      + returns all nodes matching the selector in document order
      - returns error on an invalid selector
      # querying
    std.html.text
      @ (node: html_node) -> string
      + returns the concatenated text content of the node
      # inspection
    std.html.attr
      @ (node: html_node, name: string) -> optional[string]
      + returns the attribute value when present
      - returns none when the attribute is missing
      # inspection

html_unmarshal
  html_unmarshal.new_schema
    @ () -> schema_state
    + creates an empty field schema
    # construction
  html_unmarshal.bind_text
    @ (schema: schema_state, field: string, selector: string) -> schema_state
    + binds a field to the text of the first element matching selector
    # schema
  html_unmarshal.bind_attr
    @ (schema: schema_state, field: string, selector: string, attr: string) -> schema_state
    + binds a field to the attribute of the first element matching selector
    # schema
  html_unmarshal.bind_list
    @ (schema: schema_state, field: string, selector: string) -> schema_state
    + binds a field to the list of text values for every matching element
    # schema
  html_unmarshal.unmarshal
    @ (schema: schema_state, source: string) -> result[map[string, string], string]
    + returns a map of field name to extracted scalar value
    - returns error when a bound selector is invalid
    - returns error when a scalar-bound field has no match
    # extraction
    -> std.html.parse
    -> std.html.select_all
    -> std.html.text
    -> std.html.attr
