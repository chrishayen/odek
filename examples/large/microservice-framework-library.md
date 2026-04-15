# Requirement: "a library for writing microservices"

Services register named actions that handle typed messages. The library routes incoming messages to the correct action, supports pluggable transports, and exposes a client for invoking remote actions.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: u16) -> result[listener, string]
      + returns a listener bound to the address
      - returns error when the port cannot be bound
      # networking
    std.net.tcp_accept
      fn (l: listener) -> result[connection, string]
      + returns the next accepted connection
      - returns error when the listener is closed
      # networking
    std.net.tcp_connect
      fn (host: string, port: u16, timeout_ms: i32) -> result[connection, string]
      + returns a connection to the remote service
      - returns error on timeout or refused
      # networking
    std.net.send
      fn (conn: connection, data: bytes) -> result[void, string]
      + writes data to the connection
      - returns error when the peer closed the connection
      # networking
    std.net.recv
      fn (conn: connection, max_bytes: i32) -> result[bytes, string]
      + reads up to max_bytes from the connection
      - returns error on read failure
      # networking
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

microservice
  microservice.new
    fn (name: string) -> service_state
    + creates a service with the given name and no actions
    # construction
  microservice.add_action
    fn (svc: service_state, pattern: map[string, string], handler: fn(map[string, string]) -> result[map[string, string], string]) -> service_state
    + adds an action whose handler runs when an incoming message matches every key/value pair in the pattern
    # registration
  microservice.match_action
    fn (svc: service_state, message: map[string, string]) -> optional[action_ref]
    + returns the most specific registered action whose pattern is a subset of the message
    ? specificity is measured by the number of pattern keys
    # routing
  microservice.encode_envelope
    fn (correlation_id: string, payload: map[string, string]) -> string
    + returns a serialized request envelope wrapping the payload with the correlation id
    # framing
    -> std.json.encode_object
  microservice.decode_envelope
    fn (raw: string) -> result[envelope, string]
    + returns the parsed envelope with correlation id and payload
    - returns error on invalid framing
    # framing
    -> std.json.parse_object
  microservice.handle_message
    fn (svc: service_state, message: map[string, string]) -> result[map[string, string], string]
    + runs the matching action and returns its response payload
    - returns error when no action matches
    # dispatch
  microservice.listen
    fn (svc: service_state, host: string, port: u16) -> result[void, string]
    + accepts connections and dispatches each message through the service
    - returns error when the listener cannot be created
    # transport
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.net.recv
    -> std.net.send
  microservice.client_new
    fn (host: string, port: u16) -> client_state
    + creates a client bound to the given address
    # client
  microservice.client_act
    fn (client: client_state, payload: map[string, string]) -> result[map[string, string], string]
    + sends the payload to the remote service and returns the response payload
    - returns error on transport failure or a remote error envelope
    # client
    -> std.net.tcp_connect
    -> std.net.send
    -> std.net.recv
    -> std.crypto.random_bytes
