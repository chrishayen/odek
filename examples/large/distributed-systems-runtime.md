# Requirement: "a distributed systems runtime"

A minimal runtime surface for registering services, discovering peers, and dispatching RPC requests. The transport and registry are pluggable via opaque state.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.encoding
    std.encoding.encode_message
      fn (payload: map[string, bytes]) -> bytes
      + encodes a map of fields into a length-prefixed binary frame
      # serialization
    std.encoding.decode_message
      fn (frame: bytes) -> result[map[string, bytes], string]
      + decodes a length-prefixed frame back into a field map
      - returns error on truncated input
      # serialization
  std.net
    std.net.dial
      fn (address: string) -> result[conn_state, string]
      + opens a connection to host:port
      - returns error when address is unreachable
      # networking
    std.net.send_frame
      fn (conn: conn_state, frame: bytes) -> result[void, string]
      + writes a framed message to the connection
      # networking
    std.net.recv_frame
      fn (conn: conn_state) -> result[bytes, string]
      + reads the next framed message from the connection
      - returns error on closed connection
      # networking

runtime
  runtime.new
    fn (node_id: string) -> runtime_state
    + creates a runtime with an empty service table and registry
    # construction
  runtime.register_service
    fn (state: runtime_state, name: string, version: string) -> runtime_state
    + adds a named service endpoint to the local service table
    + replaces an existing entry with the same name
    # service_registration
    -> std.time.now_millis
  runtime.deregister_service
    fn (state: runtime_state, name: string) -> runtime_state
    + removes a service from the local table
    - leaves the table unchanged when the name is not registered
    # service_registration
  runtime.list_services
    fn (state: runtime_state) -> list[string]
    + returns the names of all locally registered services
    # introspection
  runtime.add_peer
    fn (state: runtime_state, node_id: string, address: string) -> runtime_state
    + records a peer node with its contact address
    # registry
  runtime.resolve
    fn (state: runtime_state, service: string) -> result[string, string]
    + returns a peer address known to host the named service
    - returns error when no peer advertises the service
    # discovery
  runtime.call
    fn (state: runtime_state, service: string, method: string, body: bytes) -> result[bytes, string]
    + dials the resolved peer, sends a request frame, and returns the response body
    - returns error when no peer hosts the service
    - returns error when the transport fails
    # rpc
    -> std.net.dial
    -> std.encoding.encode_message
    -> std.net.send_frame
    -> std.net.recv_frame
    -> std.encoding.decode_message
  runtime.serve_once
    fn (state: runtime_state, conn: conn_state, handler: string) -> result[void, string]
    + reads one request, dispatches to the handler tag, and writes the response frame
    - returns error when the request frame is malformed
    # rpc
    -> std.net.recv_frame
    -> std.encoding.decode_message
    -> std.encoding.encode_message
    -> std.net.send_frame
  runtime.heartbeat_tick
    fn (state: runtime_state) -> runtime_state
    + updates last-seen timestamps for active peers
    + drops peers whose last-seen exceeds the stale threshold
    # membership
    -> std.time.now_millis
