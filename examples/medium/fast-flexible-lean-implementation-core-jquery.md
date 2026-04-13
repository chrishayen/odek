# Requirement: "a server-side HTML parser with a query selector API"

Parses HTML into a DOM tree, resolves CSS selectors against it, and exposes basic traversal and attribute reads.

std: (all units exist)

html_query
  html_query.parse
    @ (source: string) -> result[dom_node, string]
    + parses an HTML document and returns the root node
    - returns error only on unrecoverable input; tolerates unclosed tags
    # parsing
  html_query.select
    @ (root: dom_node, selector: string) -> result[list[dom_node], string]
    + returns every node matching the selector in document order
    - returns error on an invalid selector
    ? supports tag, .class, #id, [attr=value], and descendant combinators
    # selection
  html_query.select_first
    @ (root: dom_node, selector: string) -> result[optional[dom_node], string]
    + returns the first match or none
    # selection
  html_query.text
    @ (node: dom_node) -> string
    + returns the concatenated text content of node and its descendants
    # extraction
  html_query.attr
    @ (node: dom_node, name: string) -> optional[string]
    + returns the attribute value when present
    # extraction
  html_query.children
    @ (node: dom_node) -> list[dom_node]
    + returns the direct element children of node
    # traversal
  html_query.parent
    @ (node: dom_node) -> optional[dom_node]
    + returns the parent, or none for the root
    # traversal
