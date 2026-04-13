# Requirement: "a selector-based HTML parsing and querying library"

Parses HTML and lets the caller query nodes with CSS-like selectors. No mutation, no rendering.

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

query
  query.parse_html
    @ (html: string) -> result[node, string]
    + parses an HTML string into a node tree tolerant of minor malformations
    - returns error only when the input is empty
    # parsing
    -> std.strings.to_lower
  query.parse_selector
    @ (sel: string) -> result[selector, string]
    + parses a selector of the form "tag", "#id", ".class", or combinations like "div.card a"
    - returns error on unsupported syntax
    # selectors
    -> std.strings.trim
  query.find
    @ (root: node, sel: selector) -> list[node]
    + returns all descendants matching sel in document order
    # queries
  query.find_one
    @ (root: node, sel: selector) -> optional[node]
    + returns the first descendant matching sel
    - returns none when nothing matches
    # queries
  query.text
    @ (n: node) -> string
    + returns the concatenated text content of n and its descendants
    # accessors
    -> std.strings.trim
  query.attribute
    @ (n: node, name: string) -> optional[string]
    + returns the value of a named attribute
    - returns none when absent
    # accessors
