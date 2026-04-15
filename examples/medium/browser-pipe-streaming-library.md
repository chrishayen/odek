# Requirement: "a library that streams piped input to connected web browsers"

Accepts bytes on an input channel and pushes them live to any HTTP client subscribed to the stream endpoint.

std
  std.net
    std.net.http_listen
      fn (port: i32) -> result[listener_state, string]
      + binds an HTTP listener on port
      - returns error when the port is in use
      # networking
    std.net.http_write_chunk
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes a chunk-encoded fragment to an open HTTP response
      - returns error when the connection is closed
      # networking

browser_pipe
  browser_pipe.start
    fn (port: i32) -> result[pipe_state, string]
    + starts a pipe server that accepts browser subscribers on port
    - returns error when the port cannot be bound
    # construction
    -> std.net.http_listen
  browser_pipe.write
    fn (state: pipe_state, data: bytes) -> pipe_state
    + fans data out to every currently subscribed browser
    + buffers recent data so late subscribers see a replay window
    # broadcast
    -> std.net.http_write_chunk
  browser_pipe.subscriber_count
    fn (state: pipe_state) -> i32
    + returns the number of currently connected browsers
    # introspection
  browser_pipe.close
    fn (state: pipe_state) -> void
    + closes all subscribers and the listener
    # teardown
