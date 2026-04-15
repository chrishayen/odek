# Requirement: "a library for building tcp servers that speak a redis-style wire protocol"

Accepts connections, parses inbound RESP commands, dispatches to registered handlers, and encodes replies. Project layer owns the server and handler registry; std provides sockets.

std
  std.socket
    std.socket.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds a tcp listener on the given address
      - returns error when the port is in use
      # networking
    std.socket.accept
      fn (listener: listener) -> result[connection, string]
      + accepts the next incoming connection
      # networking
    std.socket.read
      fn (conn: connection, max: i64) -> result[bytes, string]
      + reads up to max bytes from the connection
      # networking
    std.socket.write
      fn (conn: connection, data: bytes) -> result[i64, string]
      + writes data to the connection and returns bytes written
      # networking

resp_server
  resp_server.new
    fn () -> server_state
    + creates a server with an empty handler registry
    # construction
  resp_server.register
    fn (s: server_state, name: string, handler: command_handler) -> server_state
    + associates an uppercase command name with a handler
    # routing
  resp_server.parse_command
    fn (buf: bytes) -> result[tuple[list[string], i64], string]
    + parses one RESP multi-bulk command and returns it plus bytes consumed
    - returns error on malformed RESP
    # wire_protocol
  resp_server.encode_reply
    fn (value: reply_value) -> bytes
    + encodes a reply as RESP bytes
    # wire_protocol
  resp_server.handle_connection
    fn (s: server_state, conn: connection) -> result[void, string]
    + reads commands from the connection, dispatches, and writes replies until close
    - returns error on wire-level failure
    # execution
    -> std.socket.read
    -> std.socket.write
  resp_server.serve
    fn (s: server_state, host: string, port: i32) -> result[void, string]
    + binds and accepts connections forever, handing each to handle_connection
    # lifecycle
    -> std.socket.listen_tcp
    -> std.socket.accept
