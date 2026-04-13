# Requirement: "a bidirectional RPC library supporting request-response, notifications, and broadcast"

A symmetric message framing and dispatch layer. Transport read/write are delegated so the same core works over TCP, pipes, or in-memory channels.

std
  std.io
    std.io.read_exact
      @ (conn: conn_handle, n: i32) -> result[bytes, string]
      + reads exactly n bytes from the connection
      - returns error on short read or disconnect
      # io
    std.io.write_all
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error on disconnect
      # io
  std.encoding
    std.encoding.varint_encode
      @ (n: u64) -> bytes
      + returns the unsigned LEB128 encoding
      # encoding
    std.encoding.varint_decode
      @ (data: bytes) -> result[tuple[u64, i32], string]
      + returns the decoded value and the number of bytes consumed
      # encoding

arpc
  arpc.new_peer
    @ () -> peer_state
    + creates a peer with empty request registry and handler table
    # construction
  arpc.register_handler
    @ (p: peer_state, method: string, handler: fn[bytes, bytes]) -> peer_state
    + associates a handler with a method name
    # registration
  arpc.encode_frame
    @ (kind: i32, id: u64, method: string, payload: bytes) -> bytes
    + serializes a frame with kind, correlation id, method, and payload
    # framing
    -> std.encoding.varint_encode
  arpc.decode_frame
    @ (data: bytes) -> result[frame, string]
    + parses a frame from the wire buffer
    - returns error on truncated or malformed data
    # framing
    -> std.encoding.varint_decode
  arpc.call
    @ (p: peer_state, conn: conn_handle, method: string, payload: bytes) -> result[bytes, string]
    + sends a request frame and waits for the matching response
    - returns error when the remote returns an error frame
    # rpc
    -> arpc.encode_frame
    -> arpc.decode_frame
    -> std.io.write_all
    -> std.io.read_exact
  arpc.notify
    @ (p: peer_state, conn: conn_handle, method: string, payload: bytes) -> result[void, string]
    + sends a one-way notification frame
    # rpc
    -> arpc.encode_frame
    -> std.io.write_all
  arpc.broadcast
    @ (p: peer_state, conns: list[conn_handle], method: string, payload: bytes) -> i32
    + sends the same notification to every connection, returns count of successes
    # rpc
    -> arpc.notify
  arpc.handle_incoming
    @ (p: peer_state, conn: conn_handle) -> result[peer_state, string]
    + reads one frame and dispatches it to the registered handler or a pending call
    - returns error on transport failure
    # dispatch
    -> arpc.decode_frame
    -> std.io.read_exact
    -> std.io.write_all
