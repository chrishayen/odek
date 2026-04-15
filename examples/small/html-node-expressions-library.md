# Requirement: "a library for html templates as first-class expressions"

Build HTML fragments by composing typed nodes in code rather than parsing template strings.

std: (all units exist)

html_nodes
  html_nodes.element
    fn (tag: string, attrs: map[string,string], children: list[html_node]) -> html_node
    + produces an element node with tag, attrs, and ordered children
    + empty children list yields a self-contained element
    # construction
  html_nodes.text
    fn (content: string) -> html_node
    + produces a text node whose content will be escaped on render
    # construction
  html_nodes.render
    fn (node: html_node) -> string
    + serializes an element node as "<tag attr=\"value\">...children...</tag>"
    + serializes a text node with &, <, >, " escaped
    - returns "" for a node whose tag is empty and content is empty
    # rendering
