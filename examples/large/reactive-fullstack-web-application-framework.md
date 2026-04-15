# Requirement: "a framework for building reactive full-stack web applications"

A reactive full-stack framework: a server holds authoritative component state, client events mutate it, and a diff is pushed back to subscribed clients. Components are identified by id; their state is a key-value map.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a flat string map
      - returns error on non-object or malformed JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.collections
    std.collections.map_diff
      fn (before: map[string,string], after: map[string,string]) -> map[string,string]
      + returns keys whose values differ or are newly added
      + removed keys map to the empty string
      # collections
  std.ids
    std.ids.new_id
      fn () -> string
      + returns a unique opaque id string
      # identity

webapp
  webapp.new
    fn () -> app_state
    + creates an app with no components and no connected clients
    # construction
  webapp.register_component
    fn (state: app_state, component_type: string, initial: map[string,string]) -> tuple[app_state, string]
    + stores a new component instance with the given initial state and returns its id
    # component_lifecycle
    -> std.ids.new_id
  webapp.get_component_state
    fn (state: app_state, component_id: string) -> optional[map[string,string]]
    + returns the current state map or none
    # state_access
  webapp.dispatch_event
    fn (state: app_state, component_id: string, event_name: string, payload: string) -> result[app_state, string]
    + routes the event to the component and applies its state transition
    - returns error "unknown component" when component_id is not registered
    - returns error "unknown event" when the component has no handler for event_name
    # event_dispatch
  webapp.connect_client
    fn (state: app_state) -> tuple[app_state, string]
    + allocates a client id and an empty subscription set
    # client_lifecycle
    -> std.ids.new_id
  webapp.subscribe
    fn (state: app_state, client_id: string, component_id: string) -> app_state
    + adds component_id to the client's subscription set
    # subscription
  webapp.unsubscribe
    fn (state: app_state, client_id: string, component_id: string) -> app_state
    + removes component_id from the client's subscription set
    # subscription
  webapp.diff_for_client
    fn (state: app_state, client_id: string, previous: map[string,map[string,string]]) -> map[string,map[string,string]]
    + returns, per subscribed component, the keys whose values changed since previous
    ? the caller tracks previous snapshots per client; the framework does not persist them
    # reactivity
    -> std.collections.map_diff
  webapp.serialize_diff
    fn (diff: map[string,map[string,string]]) -> string
    + encodes a per-component diff as JSON for transport
    # wire_format
    -> std.json.encode_object
  webapp.apply_client_message
    fn (state: app_state, client_id: string, raw: string) -> result[app_state, string]
    + parses an inbound JSON message and dispatches the enclosed event
    - returns error on malformed JSON or missing component_id/event_name fields
    # wire_format
    -> std.json.parse_object
  webapp.disconnect_client
    fn (state: app_state, client_id: string) -> app_state
    + drops the client and its subscriptions
    # client_lifecycle
