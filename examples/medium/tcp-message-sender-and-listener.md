# Requirement: "a friendly library for sending messages and listening on TCP sockets"

Wraps a TCP socket with length-prefixed message framing so callers exchange discrete messages instead of byte streams.

std
  std.net
    std.net.tcp_listen
      fn (addr: string) -> result[listener_handle, string]
      + binds and listens on the given address
      - returns error when the address is already in use
      # network
    std.net.tcp_accept
      fn (listener: listener_handle) -> result[conn_handle, string]
      + accepts the next incoming connection
      # network
    std.net.tcp_dial
      fn (addr: string) -> result[conn_handle, string]
      + opens a TCP connection to the address
      - returns error when the peer is unreachable
      # network
    std.net.tcp_read
      fn (conn: conn_handle, n: i32) -> result[bytes, string]
      + reads exactly n bytes from the connection
      - returns error on EOF before n bytes are received
      # network
    std.net.tcp_write
      fn (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # network
    std.net.tcp_close
      fn (conn: conn_handle) -> void
      + closes the connection
      # network

tcp_msg
  tcp_msg.serve
    fn (addr: string, handler: fn(conn_handle, bytes) -> bytes) -> result[void, string]
    + listens on the address and calls handler for each framed message, writing the reply
    - returns error when binding fails
    # server
    -> std.net.tcp_listen
    -> std.net.tcp_accept
  tcp_msg.connect
    fn (addr: string) -> result[conn_handle, string]
    + dials the address and returns an active connection
    # client
    -> std.net.tcp_dial
  tcp_msg.send
    fn (conn: conn_handle, msg: bytes) -> result[void, string]
    + writes a 4-byte big-endian length prefix followed by the payload
    - returns error when the payload exceeds 2^31 - 1 bytes
    # framing
    -> std.net.tcp_write
  tcp_msg.receive
    fn (conn: conn_handle) -> result[bytes, string]
    + reads a length prefix then the payload of that size
    - returns error on a short read
    # framing
    -> std.net.tcp_read
  tcp_msg.close
    fn (conn: conn_handle) -> void
    + closes the connection
    # teardown
    -> std.net.tcp_close
