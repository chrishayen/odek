# Requirement: "an HTML parser and DOM library"

Parses HTML into a node tree and offers query and mutation operations. No rendering, no layout, no scripting.

std
  std.strings
    std.strings.to_lower
      @ (s: string) -> string
      + returns an ASCII-lowercased copy
      # strings
    std.strings.trim
      @ (s: string) -> string
      + trims ASCII whitespace from both ends
      # strings

dom
  dom.parse
    @ (html: string) -> result[node, string]
    + parses an HTML document and returns the root node
    - returns error when input cannot be coerced into any node tree
    ? tolerant of unclosed tags like a browser; never panics on malformed input
    # parsing
    -> std.strings.to_lower
  dom.tag_name
    @ (n: node) -> optional[string]
    + returns the element tag name in lowercase
    - returns none for text nodes and the document root
    # introspection
  dom.text_content
    @ (n: node) -> string
    + returns the concatenated text of n and all descendants
    # introspection
    -> std.strings.trim
  dom.get_attribute
    @ (n: node, name: string) -> optional[string]
    + returns the value of the named attribute on an element
    - returns none when attribute or element does not exist
    # attributes
  dom.set_attribute
    @ (n: node, name: string, value: string) -> node
    + returns the node with the attribute added or replaced
    # attributes
  dom.children
    @ (n: node) -> list[node]
    + returns direct child nodes in document order
    # traversal
  dom.find_by_tag
    @ (root: node, tag: string) -> list[node]
    + returns all descendant elements with the given tag
    # queries
    -> std.strings.to_lower
  dom.find_by_id
    @ (root: node, id: string) -> optional[node]
    + returns the first descendant element whose id attribute equals id
    - returns none when not found
    # queries
  dom.find_by_class
    @ (root: node, class_name: string) -> list[node]
    + returns all descendant elements containing class_name in their class attribute
    # queries
  dom.append_child
    @ (parent: node, child: node) -> node
    + returns parent with child appended as its last child
    # mutation
  dom.remove_child
    @ (parent: node, child: node) -> node
    + returns parent with the given child removed
    - returns parent unchanged when child is not present
    # mutation
  dom.serialize
    @ (root: node) -> string
    + renders the node tree back to HTML
    # serialization
