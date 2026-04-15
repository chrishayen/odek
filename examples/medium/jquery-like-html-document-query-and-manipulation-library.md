# Requirement: "a jQuery-like HTML document query and manipulation library"

Parse an HTML document into a tree, then run CSS-selector queries and mutations against it. The project layer is a fluent selection API; std handles parsing and selector compilation.

std
  std.html
    std.html.parse
      fn (source: string) -> result[dom_node, string]
      + parses an HTML document into a DOM tree
      - returns error on unrecoverable syntax
      # parsing
    std.html.serialize
      fn (root: dom_node) -> string
      + serializes a DOM tree back to HTML
      # serialization
  std.css
    std.css.compile_selector
      fn (selector: string) -> result[compiled_selector, string]
      + compiles a CSS selector expression for repeated matching
      - returns error on invalid selector syntax
      # selector_compilation
    std.css.matches
      fn (sel: compiled_selector, node: dom_node) -> bool
      + returns true when the node matches the compiled selector
      # selector_matching

doc_query
  doc_query.from_html
    fn (source: string) -> result[selection, string]
    + parses source and returns a selection wrapping the root node
    - returns error on unparseable input
    # construction
    -> std.html.parse
  doc_query.find
    fn (sel: selection, selector: string) -> result[selection, string]
    + returns a new selection of descendants matching the selector
    - returns error on invalid selector
    # query
    -> std.css.compile_selector
    -> std.css.matches
  doc_query.text
    fn (sel: selection) -> string
    + returns the concatenated text content of all nodes in the selection
    # accessor
  doc_query.attr
    fn (sel: selection, name: string) -> optional[string]
    + returns the attribute value of the first matched node
    - returns none when the attribute is missing or selection is empty
    # accessor
  doc_query.set_attr
    fn (sel: selection, name: string, value: string) -> selection
    + sets the attribute on every node in the selection
    # mutation
  doc_query.html
    fn (sel: selection) -> string
    + returns the serialized HTML of the first matched node
    # accessor
    -> std.html.serialize
