# Requirement: "a minimal template engine (variable substitution and conditionals)"

Two project functions: compile once, render many. Keeps the expensive parse step separate from the hot render path.

std: (all units exist)

template
  template.compile
    fn (source: string) -> result[compiled_template, string]
    + parses template source with {{var}} and {% if flag %}...{% endif %} syntax
    - returns error on unclosed {{ ... }} tags
    - returns error on unbalanced {% if %} / {% endif %}
    # compilation
  template.render
    fn (tmpl: compiled_template, context: map[string, string]) -> string
    + substitutes {{name}} with context["name"]
    + evaluates {% if flag %}...{% endif %} blocks using truthy context values
    + missing variables render as empty strings
    ? HTML escaping is the caller's responsibility, not the engine's
    # rendering
