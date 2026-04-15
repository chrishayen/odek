# Requirement: "a templating library extending mustache with helpers and advanced blocks"

Parses a template into an AST of text, variable, block, and helper-call nodes; renders against a context with registered helpers.

std: (all units exist)

htpl
  htpl.parse
    fn (source: string) -> result[template_ast, string]
    + parses text, {{var}}, {{#block}}...{{/block}}, {{^inverse}}...{{/inverse}}, and {{helper arg1 arg2}}
    - returns error on unmatched block tags
    - returns error on malformed tag syntax
    # parsing
  htpl.new_env
    fn () -> env_state
    + creates a rendering environment with an empty helper registry
    # construction
  htpl.register_helper
    fn (env: env_state, name: string, helper: helper_fn) -> env_state
    + installs a helper callable invoked when a tag uses its name
    ? helper_fn takes a list of string args and the current context and returns a string
    # helpers
  htpl.render
    fn (env: env_state, ast: template_ast, context: map[string, string]) -> result[string, string]
    + resolves variables via context, walks blocks, and invokes helpers
    - returns error when a helper name is not registered
    # rendering
  htpl.render_string
    fn (env: env_state, source: string, context: map[string, string]) -> result[string, string]
    + convenience that parses then renders in one step
    - returns error on parse or render failure
    # rendering
    -> htpl.parse
    -> htpl.render
  htpl.escape_html
    fn (text: string) -> string
    + escapes &, <, >, ", and ' for safe HTML interpolation
    # escaping
