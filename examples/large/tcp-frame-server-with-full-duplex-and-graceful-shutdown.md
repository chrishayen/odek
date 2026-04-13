# Requirement: "a TCP server framework with full-duplex connections, graceful shutdown, and pluggable framing"

Handles accept loops, per-connection read/write tasks, a framing hook for custom protocols, and graceful shutdown.

std
  std.net
    std.net.listen_tcp
      @ (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on host:port
      - returns error when the port is in use
      # network
    std.net.accept
      @ (listener: listener_state) -> result[conn_state, string]
      + blocks until a client connects
      - returns error when the listener has been closed
      # network
    std.net.conn_read
      @ (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes
      - returns error when the connection is closed
      # network
    std.net.conn_write
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes all of data
      - returns error when the connection is closed
      # network
    std.net.conn_close
      @ (conn: conn_state) -> void
      + closes the connection
      # network
    std.net.listener_close
      @ (listener: listener_state) -> void
      + stops accepting new connections
      # network
  std.thread
    std.thread.spawn
      @ (work: thread_fn) -> result[thread_handle, string]
      + starts a worker thread
      - returns error when the OS refuses the thread
      # threading
    std.thread.join
      @ (handle: thread_handle) -> void
      + waits for the thread to exit
      # threading

tcpframe
  tcpframe.server_new
    @ (host: string, port: u16) -> result[server_state, string]
    + creates a server bound to host:port
    - returns error when the bind fails
    # construction
    -> std.net.listen_tcp
  tcpframe.set_framer
    @ (server: server_state, framer: framer_fn) -> server_state
    + installs a function that splits a byte stream into message frames
    ? framer returns (message, remaining) or "need more" when incomplete
    # framing
  tcpframe.set_handler
    @ (server: server_state, handler: handler_fn) -> server_state
    + installs the per-message handler that returns response bytes
    # handler
  tcpframe.serve
    @ (server: server_state) -> result[void, string]
    + runs the accept loop, spawning a handler thread per connection
    - returns error when the listener fails
    # lifecycle
    -> std.net.accept
    -> std.thread.spawn
  tcpframe.handle_connection
    @ (conn: conn_state, framer: framer_fn, handler: handler_fn) -> void
    + reads, frames, dispatches, and writes until the connection closes
    # per_connection
    -> std.net.conn_read
    -> std.net.conn_write
    -> std.net.conn_close
  tcpframe.broadcast
    @ (server: server_state, data: bytes) -> i32
    + writes data to every active connection; returns how many received it
    # messaging
    -> std.net.conn_write
  tcpframe.connections
    @ (server: server_state) -> i32
    + returns the count of active connections
    # introspection
  tcpframe.shutdown
    @ (server: server_state, drain_timeout_ms: i64) -> result[void, string]
    + stops accepting, waits for handlers to finish up to the timeout, then closes all connections
    - returns error when handlers do not finish within the timeout
    # lifecycle
    -> std.net.listener_close
    -> std.thread.join
    -> std.net.conn_close
  tcpframe.framer_length_prefix
    @ (header_size: i32) -> framer_fn
    + returns a framer that reads a big-endian length header and then the payload
    ? header_size must be 1, 2, 4, or 8
    # framing
  tcpframe.framer_delimited
    @ (delimiter: bytes) -> framer_fn
    + returns a framer that splits on the given delimiter
    # framing
