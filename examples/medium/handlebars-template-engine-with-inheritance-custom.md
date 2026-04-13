# Requirement: "a template engine with block inheritance and user-registered helpers"

std: (all units exist)

template_engine
  template_engine.new
    @ () -> engine_state
    + creates an engine with no registered templates or helpers
    # construction
  template_engine.register_template
    @ (state: engine_state, name: string, source: string) -> result[engine_state, string]
    + parses and stores a template under the given name
    - returns error on parse failure
    # registration
  template_engine.register_partial
    @ (state: engine_state, name: string, source: string) -> result[engine_state, string]
    + parses and stores a partial used by {{> name}} references
    - returns error on parse failure
    # registration
  template_engine.register_helper
    @ (state: engine_state, name: string, helper: helper_fn) -> engine_state
    + stores a helper callable that receives arguments and a context
    # registration
  template_engine.render
    @ (state: engine_state, name: string, context: map[string, string]) -> result[string, string]
    + expands the named template, resolving partials, helpers, and variables in context
    - returns error when the template name is not registered
    - returns error when a helper is invoked that has not been registered
    # rendering
  template_engine.render_with_parent
    @ (state: engine_state, child: string, parent: string, context: map[string, string]) -> result[string, string]
    + renders child as an inheriting template whose named blocks override those of parent
    - returns error when parent or child is not registered
    # inheritance
