# Requirement: "a high-throughput RPC library"

Length-prefixed binary RPC over a transport. Servers register handlers by method name; clients send requests and receive responses.

std
  std.encoding
    std.encoding.varint_encode
      fn (value: u64) -> bytes
      + encodes an unsigned integer using LEB128
      # encoding
    std.encoding.varint_decode
      fn (data: bytes) -> result[tuple[u64, i32], string]
      + returns the decoded value and the number of bytes consumed
      - returns error on truncated input
      # encoding
  std.io
    std.io.read_exact
      fn (conn: connection, n: i32) -> result[bytes, string]
      + reads exactly n bytes from conn
      - returns error on connection close before n bytes
      # io
    std.io.write_all
      fn (conn: connection, data: bytes) -> result[void, string]
      + writes all of data to conn
      - returns error on write failure
      # io

rpc
  rpc.encode_request
    fn (method: string, payload: bytes) -> bytes
    + returns a framed request containing the method name and payload
    # framing
    -> std.encoding.varint_encode
  rpc.decode_request
    fn (frame: bytes) -> result[tuple[string, bytes], string]
    + returns the method name and payload
    - returns error on malformed framing
    # framing
    -> std.encoding.varint_decode
  rpc.encode_response
    fn (status: i32, payload: bytes) -> bytes
    + returns a framed response with the given status and body
    # framing
    -> std.encoding.varint_encode
  rpc.decode_response
    fn (frame: bytes) -> result[tuple[i32, bytes], string]
    + returns the status code and body
    - returns error on malformed framing
    # framing
    -> std.encoding.varint_decode
  rpc.new_server
    fn () -> server_state
    + returns an empty server with no registered handlers
    # construction
  rpc.register_handler
    fn (server: server_state, method: string, handler: fn(bytes) -> result[bytes, string]) -> server_state
    + associates method with handler
    + replaces any previous handler for the same method
    # registration
  rpc.serve_one
    fn (server: server_state, conn: connection) -> result[void, string]
    + reads one request, dispatches to the registered handler, writes response
    - returns error when no handler is registered for the method
    - returns error on io failure
    # dispatch
    -> std.io.read_exact
    -> rpc.decode_request
    -> rpc.encode_response
    -> std.io.write_all
  rpc.call
    fn (conn: connection, method: string, payload: bytes) -> result[bytes, string]
    + sends a request and returns the response body
    - returns error on non-zero status
    - returns error on io failure
    # client
    -> rpc.encode_request
    -> std.io.write_all
    -> std.io.read_exact
    -> rpc.decode_response
