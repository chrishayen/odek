# Requirement: "a framework for exposing the same actions over TCP sockets, WebSockets, and HTTP"

Actions are named units of work with typed input and output. The framework registers them once and dispatches requests arriving over any transport through a common path: decode input, run action, encode output.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a flat string map
      - returns error on malformed JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a flat string map as a JSON object
      # serialization
  std.ids
    std.ids.new_id
      fn () -> string
      + returns a unique opaque id
      # identity

actions
  actions.new
    fn () -> server_state
    + creates a server with no registered actions and no open connections
    # construction
  actions.register_action
    fn (state: server_state, name: string, handler_id: string) -> server_state
    + binds an action name to a handler id
    - returns unchanged state when name is already registered
    # registry
  actions.lookup
    fn (state: server_state, name: string) -> result[string, string]
    + returns the handler_id for the given action
    - returns error "unknown action" when name is not registered
    # registry
  actions.decode_request
    fn (raw: string) -> result[request, string]
    + parses a JSON envelope with action, request_id, and params fields
    - returns error when action or request_id is missing
    # wire_format
    -> std.json.parse_object
  actions.encode_response
    fn (request_id: string, status: string, body: map[string,string]) -> string
    + encodes a JSON envelope containing request_id, status, and body
    # wire_format
    -> std.json.encode_object
  actions.handle_tcp
    fn (state: server_state, connection_id: string, data: bytes) -> tuple[server_state, bytes]
    + frames newline-delimited JSON from TCP, dispatches complete requests, and returns reply bytes
    # transport_tcp
  actions.handle_websocket
    fn (state: server_state, connection_id: string, message: string) -> tuple[server_state, string]
    + dispatches one text-frame request and returns the reply text
    # transport_websocket
  actions.handle_http
    fn (state: server_state, method: string, path: string, body: string) -> tuple[server_state, i32, string]
    + maps POST /actions/{name} to a dispatched request and returns (state, status, body)
    - returns (state, 404, "") when path does not match /actions/{name}
    - returns (state, 405, "") when method is not POST
    # transport_http
  actions.open_connection
    fn (state: server_state, transport: string) -> tuple[server_state, string]
    + allocates a connection id tagged with its transport
    # connection_lifecycle
    -> std.ids.new_id
  actions.close_connection
    fn (state: server_state, connection_id: string) -> server_state
    + releases resources held for the given connection
    # connection_lifecycle
  actions.broadcast
    fn (state: server_state, transport: string, message: string) -> list[string]
    + returns the ids of connections on the given transport that should receive the message
    ? actual transmission is the caller's responsibility
    # fan_out
