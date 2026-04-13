# Requirement: "a text template engine with live preview support"

Parses a template once, then renders it repeatedly against different data for live preview. Errors include source positions.

std: (all units exist)

template
  template.parse
    @ (source: string) -> result[template_ast, string]
    + returns a parsed AST on success
    - returns error with line and column on unmatched delimiters
    - returns error on unknown directives
    # parsing
  template.render
    @ (ast: template_ast, data: map[string, dynamic_value]) -> result[string, string]
    + returns the rendered string
    - returns error when the template references an undefined key
    # rendering
  template.referenced_keys
    @ (ast: template_ast) -> list[string]
    + returns the names of all top-level keys the template reads
    ? used by preview UIs to show which fields the template needs
    # introspection
