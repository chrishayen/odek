# Requirement: "a realtime api gateway that exposes resources over rest, a realtime push channel, and rpc, keeping all clients synchronized"

A gateway that lets clients fetch resources over request/response, subscribe to change events, and invoke remote procedures. The project layer owns resource registration, subscription routing, and change broadcasting; std provides json, hashing, and socket primitives.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses any JSON value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      fn (value: json_value) -> string
      + encodes a JSON value as a string
      # serialization
  std.hash
    std.hash.sha1_hex
      fn (data: bytes) -> string
      + returns the sha1 of data as a lowercase hex string
      # hashing
  std.socket
    std.socket.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds a tcp listener on the given address
      - returns error when the port is already in use
      # networking
    std.socket.accept
      fn (listener: listener) -> result[connection, string]
      + accepts the next incoming connection
      # networking

gateway
  gateway.new
    fn () -> gateway_state
    + creates an empty gateway with no registered resources
    # construction
  gateway.register_resource
    fn (g: gateway_state, pattern: string, handler: resource_handler) -> gateway_state
    + associates a resource pattern (e.g. "user.$id") with a fetch handler
    # routing
  gateway.register_call
    fn (g: gateway_state, pattern: string, handler: call_handler) -> gateway_state
    + associates an rpc method name with a handler
    # routing
  gateway.fetch
    fn (g: gateway_state, resource_id: string) -> result[json_value, string]
    + invokes the matching resource handler and returns its current value
    - returns error when no pattern matches the resource id
    # query
    -> std.json.encode_value
  gateway.call
    fn (g: gateway_state, method: string, params: json_value) -> result[json_value, string]
    + invokes the matching rpc handler with params
    - returns error when no method is registered
    # rpc
  gateway.subscribe
    fn (g: gateway_state, client: client_id, resource_id: string) -> tuple[gateway_state, json_value]
    + records a subscription and returns the current value for initial state
    # subscription
    -> std.hash.sha1_hex
  gateway.unsubscribe
    fn (g: gateway_state, client: client_id, resource_id: string) -> gateway_state
    + removes the subscription
    # subscription
  gateway.publish_change
    fn (g: gateway_state, resource_id: string, new_value: json_value) -> tuple[gateway_state, list[client_change_event]]
    + updates the cached value and returns the list of change events to deliver
    + produces one event per subscribed client
    # broadcast
  gateway.handle_client_message
    fn (g: gateway_state, client: client_id, raw: string) -> tuple[gateway_state, string]
    + parses a json message from the client and produces a json reply
    - returns an error reply when the message is malformed
    # protocol
    -> std.json.parse_value
    -> std.json.encode_value
  gateway.run_listener
    fn (g: gateway_state, host: string, port: i32) -> result[void, string]
    + binds and serves the gateway on the given address until shutdown
    # networking
    -> std.socket.listen_tcp
    -> std.socket.accept
  gateway.disconnect_client
    fn (g: gateway_state, client: client_id) -> gateway_state
    + drops all subscriptions for the client
    # subscription
