# Requirement: "a typesafe client for an in-memory key-value store speaking a line-oriented request/response protocol"

The project layer is a small command surface over a socket and a RESP-style encoder/decoder.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[tcp_conn, string]
      + opens a TCP connection
      - returns error when the host cannot be reached
      # networking
    std.net.write_all
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes or returns an error
      # networking
    std.net.read_line
      fn (conn: tcp_conn) -> result[string, string]
      + reads up to and including the next "\r\n"
      - returns error on EOF before a line terminator
      # networking
    std.net.read_exact
      fn (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      - returns error on short read
      # networking
    std.net.close
      fn (conn: tcp_conn) -> void
      # networking

kv_client
  kv_client.connect
    fn (host: string, port: i32) -> result[client_state, string]
    + dials the server and returns a client
    # construction
    -> std.net.dial_tcp
  kv_client.encode_command
    fn (parts: list[string]) -> bytes
    + encodes a command as a length-prefixed array frame
    # protocol
  kv_client.read_reply
    fn (state: client_state) -> result[reply_value, string]
    + reads one typed reply (simple string, integer, bulk string, array, or error)
    - returns error when the server sends an error reply
    # protocol
    -> std.net.read_line
    -> std.net.read_exact
  kv_client.get
    fn (state: client_state, key: string) -> result[optional[string], string]
    + returns the value at key or none when missing
    # commands
    -> std.net.write_all
  kv_client.set
    fn (state: client_state, key: string, value: string) -> result[void, string]
    + stores value at key
    # commands
    -> std.net.write_all
  kv_client.del
    fn (state: client_state, keys: list[string]) -> result[i32, string]
    + removes the keys and returns how many were deleted
    # commands
    -> std.net.write_all
  kv_client.close
    fn (state: client_state) -> void
    # teardown
    -> std.net.close
