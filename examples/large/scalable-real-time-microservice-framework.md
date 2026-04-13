# Requirement: "a scalable real-time microservice framework"

A framework for building microservices that communicate via rpc and pub/sub with client presence tracking.

std
  std.net
    std.net.listen_tcp
      @ (host: string, port: u16) -> result[listener_handle, string]
      + opens a tcp listener
      # networking
    std.net.accept
      @ (listener: listener_handle) -> result[conn_handle, string]
      + accepts an incoming connection
      # networking
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a json document into a dynamic value
      - returns error on malformed input
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + encodes a dynamic json value as a string
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

rtmf
  rtmf.new_server
    @ (host: string, port: u16) -> result[server_state, string]
    + creates a new microservice server bound to the given address
    - returns error when the port cannot be acquired
    # construction
    -> std.net.listen_tcp
  rtmf.register_rpc
    @ (state: server_state, name: string, handler: rpc_handler) -> server_state
    + registers a request/response handler under a public name
    # rpc
  rtmf.register_event
    @ (state: server_state, topic: string) -> server_state
    + declares a topic that clients can subscribe to
    # pub_sub
  rtmf.accept_client
    @ (state: server_state) -> result[tuple[client_handle, server_state], string]
    + accepts a client and assigns it a unique client id and presence entry
    # presence
    -> std.net.accept
    -> std.time.now_millis
  rtmf.dispatch_message
    @ (state: server_state, client: client_handle, raw: string) -> result[server_state, string]
    + parses an incoming message and routes it to the matching rpc or subscription
    - returns error when the message type is unknown
    # routing
    -> std.json.parse
  rtmf.call_rpc
    @ (state: server_state, name: string, payload: json_value) -> result[json_value, string]
    + invokes a registered rpc handler and returns its response
    - returns error when no handler is registered for name
    # rpc
  rtmf.subscribe
    @ (state: server_state, client: client_handle, topic: string) -> server_state
    + adds a client to the subscriber set of a topic
    # pub_sub
  rtmf.publish
    @ (state: server_state, topic: string, payload: json_value) -> result[server_state, string]
    + delivers a payload to every subscriber of a topic
    # pub_sub
    -> std.json.encode
  rtmf.presence_list
    @ (state: server_state, topic: string) -> list[string]
    + returns the ids of clients currently subscribed to a topic
    # presence
  rtmf.disconnect
    @ (state: server_state, client: client_handle) -> server_state
    + removes a client from all subscriptions and presence tracking
    # lifecycle
  rtmf.heartbeat_sweep
    @ (state: server_state, stale_after_ms: i64) -> server_state
    + disconnects clients whose last activity is older than the threshold
    # liveness
    -> std.time.now_millis
