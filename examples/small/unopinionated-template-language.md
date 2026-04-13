# Requirement: "a simple unopinionated templating language"

Parse a template string with `{{name}}` placeholders and render it against a variable map.

std: (all units exist)

template
  template.parse
    @ (source: string) -> result[list[template_token], string]
    + splits the source into literal and placeholder tokens
    + a placeholder is `{{` followed by a name followed by `}}`
    - returns error when a `{{` is never closed
    # parsing
  template.render
    @ (tokens: list[template_token], vars: map[string, string]) -> result[string, string]
    + concatenates literal tokens and substitutes placeholders from vars
    - returns error when a placeholder name is not present in vars
    # rendering
  template.render_string
    @ (source: string, vars: map[string, string]) -> result[string, string]
    + convenience that parses then renders in one call
    # rendering
