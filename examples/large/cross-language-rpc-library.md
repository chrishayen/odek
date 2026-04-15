# Requirement: "a cross-language RPC library"

A length-prefixed RPC framework: a pluggable serializer, a client that dispatches calls, and a server that routes to registered handlers.

std
  std.io
    std.io.write_bytes
      fn (conn: connection, data: bytes) -> result[i32, string]
      + writes data to the connection and returns bytes written
      - returns error on broken connection
      # io
    std.io.read_bytes
      fn (conn: connection, n: i32) -> result[bytes, string]
      + reads exactly n bytes from the connection
      - returns error on short read or connection close
      # io
  std.encoding
    std.encoding.varint_encode
      fn (n: u64) -> bytes
      + encodes n as an unsigned varint
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i32) -> result[tuple[u64, i32], string]
      + decodes a varint, returning value and new offset
      - returns error when buffer is too short
      # encoding

rpc
  rpc.new_client
    fn (conn: connection, codec: codec) -> rpc_client
    + builds a client bound to a connection with a pluggable codec
    # construction
  rpc.call
    fn (client: rpc_client, method: string, payload: bytes) -> result[bytes, rpc_error]
    + serializes the request, writes it length-prefixed, and reads the response
    - returns error when the remote returns an error frame
    - returns error on connection failure
    # client
    -> std.io.write_bytes
    -> std.io.read_bytes
    -> std.encoding.varint_encode
    -> std.encoding.varint_decode
  rpc.new_server
    fn () -> rpc_server
    + creates an empty server with no registered methods
    # construction
  rpc.register_method
    fn (server: rpc_server, method: string, handler: fn(bytes) -> result[bytes, string]) -> rpc_server
    + binds a handler function to a method name
    # registration
  rpc.serve_connection
    fn (server: rpc_server, conn: connection) -> result[void, string]
    + reads request frames until the connection closes and dispatches each to its handler
    - responds with an error frame when method is not registered
    - returns error when a handler panics or returns an error
    # server
    -> std.io.write_bytes
    -> std.io.read_bytes
  rpc.encode_request
    fn (method: string, payload: bytes) -> bytes
    + produces a length-prefixed frame containing method and payload
    # framing
    -> std.encoding.varint_encode
  rpc.decode_request
    fn (frame: bytes) -> result[tuple[string, bytes], string]
    + extracts method name and payload from a frame
    - returns error on truncated frame
    # framing
    -> std.encoding.varint_decode
  rpc.encode_response
    fn (payload: bytes, err: optional[string]) -> bytes
    + produces a response frame, tagging success or error
    # framing
  rpc.decode_response
    fn (frame: bytes) -> result[bytes, rpc_error]
    + parses a response frame, returning the payload or the error
    - returns error when the frame signals a remote error
    # framing
