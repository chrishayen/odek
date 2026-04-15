# Requirement: "a client library for a redis-compatible key-value store"

Speaks the RESP wire protocol over a TCP connection and exposes common commands. Project layer owns the connection, command dispatch, and reply parsing; std provides sockets.

std
  std.socket
    std.socket.dial_tcp
      fn (host: string, port: i32) -> result[connection, string]
      + opens a tcp connection to host:port
      - returns error on connection failure
      # networking
    std.socket.write
      fn (conn: connection, data: bytes) -> result[i64, string]
      + writes data to the connection and returns bytes written
      # networking
    std.socket.read
      fn (conn: connection, max: i64) -> result[bytes, string]
      + reads up to max bytes from the connection
      # networking

redis_client
  redis_client.connect
    fn (host: string, port: i32) -> result[client_state, string]
    + opens a connection to the server and returns a client handle
    # construction
    -> std.socket.dial_tcp
  redis_client.encode_command
    fn (args: list[string]) -> bytes
    + returns the RESP-encoded multi-bulk representation of the command
    # wire_protocol
  redis_client.parse_reply
    fn (buf: bytes) -> result[tuple[reply_value, i64], string]
    + parses one RESP reply and returns it with the number of bytes consumed
    - returns error on malformed RESP
    # wire_protocol
  redis_client.do_command
    fn (client: client_state, args: list[string]) -> result[reply_value, string]
    + sends a command and returns the parsed reply
    - returns error on wire-level failures
    # execution
    -> std.socket.write
    -> std.socket.read
  redis_client.get
    fn (client: client_state, key: string) -> result[optional[string], string]
    + returns the string value for key, or none when the key is absent
    # commands
  redis_client.set
    fn (client: client_state, key: string, value: string) -> result[void, string]
    + stores value under key
    # commands
  redis_client.del
    fn (client: client_state, key: string) -> result[i64, string]
    + deletes the key and returns the number of keys removed
    # commands
  redis_client.pipeline
    fn (client: client_state, commands: list[list[string]]) -> result[list[reply_value], string]
    + sends all commands before reading any reply, then reads replies in order
    # pipelining
  redis_client.close
    fn (client: client_state) -> result[void, string]
    + closes the underlying connection
    # lifecycle
