# Requirement: "a server-rendered component framework"

Components are functions from props to rendered HTML. The framework composes components and wires routes to a root component. Network IO is out of scope.

std
  std.text
    std.text.html_escape
      @ (s: string) -> string
      + escapes &, <, >, ", ' to their HTML entities
      # text
    std.text.render_template
      @ (template: string, vars: map[string, string]) -> result[string, string]
      + substitutes {{key}} placeholders
      - returns error when a referenced key is missing
      # templating

components
  components.define
    @ (name: string, template: string) -> component_def
    + returns a component definition bound to the given template
    # construction
  components.render
    @ (def: component_def, props: map[string, string]) -> result[string, string]
    + returns the rendered HTML with props interpolated
    + escapes prop values to prevent injection
    - returns error when the template references a missing prop
    # rendering
    -> std.text.html_escape
    -> std.text.render_template
  components.new_app
    @ () -> app_state
    + returns an app with no registered components or routes
    # construction
  components.register
    @ (state: app_state, def: component_def) -> app_state
    + registers a component under its name
    # registration
  components.route
    @ (state: app_state, path: string, component_name: string) -> app_state
    + binds a path to a registered component
    # routing
  components.serve_page
    @ (state: app_state, path: string, props: map[string, string]) -> result[string, string]
    + returns the rendered HTML for the route
    - returns error when no route matches
    - returns error when the route's component is unregistered
    # request_handling
