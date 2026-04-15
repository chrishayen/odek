# Requirement: "a library for building and rendering HTML5 tag trees"

Constructs attribute-carrying tag nodes and serializes them to well-formed HTML with proper escaping.

std: (all units exist)

html5
  html5.tag
    fn (name: string, attrs: map[string, string], children: list[html_node]) -> html_node
    + returns a tag node with the given name, attributes, and children
    # construction
  html5.text
    fn (content: string) -> html_node
    + returns a text node with its content escaped at render time
    # construction
  html5.void_tag
    fn (name: string, attrs: map[string, string]) -> html_node
    + returns a void (self-closing) tag such as br or img
    # construction
  html5.render
    fn (node: html_node) -> string
    + returns the serialized HTML string for a tree of nodes
    + escapes ampersand, less-than, greater-than, and quotes in attributes and text
    - returns an empty string for an empty tree
    # rendering
