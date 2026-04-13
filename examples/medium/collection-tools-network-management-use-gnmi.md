# Requirement: "a client library for a network device management protocol based on path-addressed configuration trees"

Clients against a streaming configuration protocol: capabilities, get, set, and subscribe.

std
  std.net
    std.net.dial_tls
      @ (host: string, port: u16) -> result[bus_state, string]
      + opens a mutually authenticated stream
      - returns error on handshake failure
      # transport
    std.net.close
      @ (conn: bus_state) -> void
      + releases the underlying stream
      # transport
  std.encoding
    std.encoding.proto_marshal
      @ (message_id: i32, fields: map[string, bytes]) -> bytes
      + encodes a tagged field map into length-prefixed wire bytes
      # serialization
    std.encoding.proto_unmarshal
      @ (payload: bytes) -> result[map[string, bytes], string]
      + decodes length-prefixed wire bytes into a tagged field map
      - returns error on truncated input
      # serialization

netmgmt
  netmgmt.connect
    @ (host: string, port: u16) -> result[bus_state, string]
    + establishes a secure session to a device
    - returns error when authentication fails
    # session
    -> std.net.dial_tls
  netmgmt.capabilities
    @ (conn: bus_state) -> result[list[string], string]
    + requests the list of supported models
    # discovery
    -> std.encoding.proto_marshal
    -> std.encoding.proto_unmarshal
  netmgmt.get
    @ (conn: bus_state, paths: list[string]) -> result[map[string, string], string]
    + returns the value at each path
    - returns error when any path is malformed
    # read
    -> std.encoding.proto_marshal
    -> std.encoding.proto_unmarshal
  netmgmt.set
    @ (conn: bus_state, updates: map[string, string]) -> result[void, string]
    + applies path-value updates atomically
    - returns error when the device rejects a path
    # write
    -> std.encoding.proto_marshal
    -> std.encoding.proto_unmarshal
  netmgmt.subscribe
    @ (conn: bus_state, paths: list[string]) -> result[bus_state, string]
    + opens a streaming subscription to the given paths
    # streaming
    -> std.encoding.proto_marshal
  netmgmt.next_update
    @ (stream: bus_state) -> result[tuple[string, string], string]
    + returns the next (path, value) event
    - returns error when the stream is closed
    # streaming
    -> std.encoding.proto_unmarshal
  netmgmt.disconnect
    @ (conn: bus_state) -> void
    + closes the session
    # session
    -> std.net.close
