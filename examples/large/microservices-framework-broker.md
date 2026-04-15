# Requirement: "a microservices framework"

A service broker registers local service handlers, routes calls by service name, and forwards calls it does not own to a pluggable transport.

std
  std.json
    std.json.encode_value
      fn (value: json_value) -> string
      + serializes a generic value tree as JSON
      # serialization
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a generic value tree
      - returns error on malformed JSON
      # serialization
  std.time
    std.time.now_millis
      fn () -> i64
      + returns unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      fn () -> string
      + returns a new random UUID
      # identifiers

broker
  broker.new
    fn (node_id: string) -> broker_state
    + creates a broker with no registered actions and the given node identifier
    # construction
  broker.register_action
    fn (state: broker_state, service: string, action: string, handler_id: string) -> broker_state
    + records that this node handles service.action with the given handler identifier
    ? handler bodies are resolved out of band by handler_id
    # registry
  broker.register_remote
    fn (state: broker_state, service: string, action: string, node_id: string) -> broker_state
    + records that a remote node handles service.action
    # registry
  broker.list_actions
    fn (state: broker_state, service: string) -> list[string]
    + returns every action registered for a service, local or remote
    # discovery
  broker.resolve
    fn (state: broker_state, service: string, action: string) -> result[resolved_target, string]
    + returns either the local handler id or the remote node id for the call
    - returns error when no node advertises the action
    # routing
  broker.encode_request
    fn (caller: string, service: string, action: string, params: json_value) -> string
    + builds a request envelope containing id, caller, target, params, and timestamp
    # transport
    -> std.id.new_uuid
    -> std.time.now_millis
    -> std.json.encode_value
  broker.decode_request
    fn (raw: string) -> result[request_envelope, string]
    + parses a request envelope
    - returns error on malformed JSON or missing fields
    # transport
    -> std.json.parse_value
  broker.encode_response
    fn (request_id: string, result: json_value, error: optional[string]) -> string
    + builds a response envelope referencing the request id
    # transport
    -> std.json.encode_value
  broker.subscribe_event
    fn (state: broker_state, event: string, handler_id: string) -> broker_state
    + registers a local handler to receive a broadcast event
    # events
  broker.list_subscribers
    fn (state: broker_state, event: string) -> list[string]
    + returns the handler ids subscribed to an event
    # events
  broker.heartbeat
    fn (state: broker_state, node_id: string) -> broker_state
    + records a liveness timestamp for a peer node
    # health
    -> std.time.now_millis
  broker.prune_dead_nodes
    fn (state: broker_state, timeout_ms: i64) -> broker_state
    + removes peer nodes whose last heartbeat is older than the timeout
    # health
    -> std.time.now_millis
