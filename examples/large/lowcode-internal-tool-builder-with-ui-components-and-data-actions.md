# Requirement: "a low-code internal-tool builder that composes UI components bound to data sources and actions"

The library models applications as a graph of components, data queries, and action handlers. Rendering and network transport are external.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random UUID v4 string
      # identifier

toolbuilder
  toolbuilder.new_app
    @ (name: string) -> app_state
    + creates an empty app with a generated id and the given name
    # construction
    -> std.uuid.new_v4
  toolbuilder.add_component
    @ (app: app_state, kind: string, props: map[string, string]) -> tuple[string, app_state]
    + adds a component of the given kind with its props and returns its generated id
    ? kinds include "table", "form", "button", "text"
    # composition
    -> std.uuid.new_v4
  toolbuilder.add_data_source
    @ (app: app_state, name: string, endpoint: string) -> app_state
    + registers a named data source pointing at an endpoint
    # data_binding
  toolbuilder.add_query
    @ (app: app_state, name: string, source: string, template: string) -> result[app_state, string]
    + registers a query that can reference the data source at runtime
    - returns error when the named source is not registered
    # data_binding
  toolbuilder.bind_component
    @ (app: app_state, component_id: string, prop: string, query_name: string) -> result[app_state, string]
    + binds a component prop to the result of a named query
    - returns error when the component or query does not exist
    # data_binding
  toolbuilder.add_action
    @ (app: app_state, name: string, target_query: string) -> result[app_state, string]
    + registers an action that runs the named query when triggered
    - returns error when the named query does not exist
    # actions
  toolbuilder.attach_event
    @ (app: app_state, component_id: string, event: string, action: string) -> result[app_state, string]
    + wires a component event (e.g. "click", "submit") to an action
    - returns error when the component or action does not exist
    # actions
  toolbuilder.resolve_component
    @ (app: app_state, component_id: string) -> result[resolved_component, string]
    + returns the component with its bound query values substituted into its props
    - returns error when the component id is unknown
    # rendering_input
  toolbuilder.dispatch_event
    @ (app: app_state, component_id: string, event: string) -> result[list[string], string]
    + returns the sequence of query names that should run in response to the event
    - returns error when no action is bound to the event
    # action_dispatch
  toolbuilder.export
    @ (app: app_state) -> string
    + serializes the entire app graph to a JSON document
    # export
    -> std.json.encode_object
  toolbuilder.import_app
    @ (raw: string) -> result[app_state, string]
    + reconstructs an app from a JSON document
    - returns error on malformed input or missing required fields
    # import
    -> std.json.parse_object
