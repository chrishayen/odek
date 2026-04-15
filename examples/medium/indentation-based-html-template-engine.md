# Requirement: "an indentation-based HTML template engine"

Templates use indentation to express nesting. The engine parses source into an AST and renders it with a context.

std: (all units exist)

template
  template.tokenize
    fn (source: string) -> list[token]
    + emits indent, dedent, tag, text, and interpolation tokens
    + collapses blank lines and strips trailing whitespace
    # lexing
  template.parse
    fn (tokens: list[token]) -> result[node, string]
    + builds a tree of element, text, and interpolation nodes from indentation
    - returns error on inconsistent indentation
    # parsing
  template.compile
    fn (source: string) -> result[node, string]
    + tokenizes and parses source into a renderable AST
    # compilation
    -> template.tokenize
    -> template.parse
  template.render
    fn (root: node, context: map[string, string]) -> string
    + walks the AST and emits HTML, substituting interpolations from context
    + escapes HTML special characters in interpolated values
    # rendering
  template.escape_html
    fn (s: string) -> string
    + replaces &, <, >, ", and ' with their HTML entities
    # escaping
