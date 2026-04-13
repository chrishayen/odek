# Requirement: "a library for building user interfaces with a declarative syntax"

A declarative UI tree is parsed from source text into a component graph, bound to a reactive state model, and rendered to a backend-agnostic draw list.

std
  std.lexer
    std.lexer.tokenize
      @ (source: string) -> result[list[token], string]
      + splits source into identifier, string, number, and punctuation tokens
      - returns error with line and column on unterminated strings
      # lexing
  std.collections
    std.collections.topo_sort
      @ (nodes: list[string], edges: list[edge]) -> result[list[string], string]
      + returns nodes in dependency order
      - returns error when a cycle is detected
      # graph

declarative_ui
  declarative_ui.parse
    @ (source: string) -> result[component_tree, string]
    + parses a declarative component definition into a tree
    - returns error on syntax errors with location info
    # parsing
    -> std.lexer.tokenize
  declarative_ui.new_state
    @ () -> ui_state
    + creates an empty reactive state store
    # construction
  declarative_ui.set
    @ (state: ui_state, key: string, value: dynamic_value) -> ui_state
    + updates a reactive variable and marks dependents dirty
    # state
  declarative_ui.get
    @ (state: ui_state, key: string) -> optional[dynamic_value]
    + returns the current value of a reactive variable
    # state
  declarative_ui.bind
    @ (tree: component_tree, state: ui_state) -> result[bound_tree, string]
    + resolves bindings from components to reactive variables
    - returns error when a binding references an unknown variable
    # binding
    -> std.collections.topo_sort
  declarative_ui.layout
    @ (tree: bound_tree, viewport_w: i32, viewport_h: i32) -> layout_tree
    + computes positions and sizes given a viewport
    # layout
  declarative_ui.render
    @ (tree: layout_tree) -> list[draw_command]
    + emits a flat list of backend-agnostic draw commands
    # rendering
  declarative_ui.dispatch_event
    @ (tree: bound_tree, event: ui_event) -> result[bound_tree, string]
    + routes a pointer or key event to the component under it
    - returns error when no component handles the event and it is required
    # events
