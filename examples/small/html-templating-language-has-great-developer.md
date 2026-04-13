# Requirement: "an html templating language with typed placeholders"

Parse a template string with `{name}` placeholders and render it against a bindings map.

std: (all units exist)

html_template
  html_template.parse
    @ (source: string) -> result[template_ast, string]
    + returns an ast that alternates literal and placeholder segments
    - returns error on an unclosed "{" placeholder
    - returns error on a placeholder name containing whitespace
    # parsing
  html_template.render
    @ (tpl: template_ast, bindings: map[string,string]) -> result[string, string]
    + substitutes each placeholder with its binding, html-escaping the value
    - returns error when a placeholder has no binding
    # rendering
