# Requirement: "a parser for indentation-based source code that returns a tree structure"

Indentation defines parent/child; blank lines and comments are skipped. The tree stores the line text and children.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits on LF, tolerating CRLF
      # strings

indent_tree
  indent_tree.parse
    fn (source: string, tab_width: i32) -> result[code_node, string]
    + returns a virtual root whose children are the top-level lines
    + uses leading whitespace to assign parents, expanding tabs to tab_width spaces
    - returns error when a child's indentation does not match any ancestor column
    # parsing
    -> std.strings.split_lines
  indent_tree.node_text
    fn (n: code_node) -> string
    + returns the line text with leading whitespace stripped
    # query
  indent_tree.node_children
    fn (n: code_node) -> list[code_node]
    + returns the immediate children of n
    # query
  indent_tree.walk
    fn (n: code_node, visit: node_visitor) -> void
    + performs a pre-order traversal invoking visit on each node
    # traversal
