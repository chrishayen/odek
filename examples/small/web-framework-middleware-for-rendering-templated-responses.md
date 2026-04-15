# Requirement: "a web framework middleware for rendering templated responses"

The original entry paired a template engine with a web framework; the generic capability is middleware that compiles templates and attaches a render function to responses.

std
  std.template
    std.template.compile
      fn (source: string) -> result[template_value, string]
      + compiles a template source into an executable template
      - returns error on unclosed tags or unknown directives
      # templating
    std.template.render
      fn (tpl: template_value, context: map[string,string]) -> result[string, string]
      + renders the template with the given context
      - returns error when a referenced variable is missing
      # templating

template_middleware
  template_middleware.new
    fn (template_dir: string) -> middleware_state
    + creates middleware that will load templates from template_dir on demand
    # construction
  template_middleware.preload
    fn (state: middleware_state, names: list[string]) -> result[middleware_state, string]
    + compiles the named templates eagerly and caches them
    - returns error when any template fails to compile
    # caching
    -> std.template.compile
  template_middleware.render_response
    fn (state: middleware_state, template_name: string, context: map[string,string]) -> result[http_response, string]
    + produces a 200 response whose body is the rendered template
    - returns error when template_name is not found
    # rendering
    -> std.template.render
