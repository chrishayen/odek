# Requirement: "a TCP library that uses a worker pool to service connections"

Accepts connections on a listener and dispatches them to a fixed pool of worker threads. Protects the server from spawning one thread per connection.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on host:port
      - returns error when the port is in use
      # network
    std.net.accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks until a client connects
      - returns error when the listener is closed
      # network
    std.net.conn_close
      fn (conn: conn_state) -> void
      + closes the connection
      # network
  std.thread
    std.thread.spawn
      fn (work: thread_fn) -> result[thread_handle, string]
      + starts a worker thread
      - returns error when the OS refuses the thread
      # threading
    std.thread.join
      fn (handle: thread_handle) -> void
      + waits for the thread to exit
      # threading

poolserver
  poolserver.new
    fn (host: string, port: u16, workers: i32, queue_depth: i32) -> result[server_state, string]
    + creates a pooled server bound to host:port with the given pool size and queue
    - returns error when the bind fails
    - returns error when workers or queue_depth are not positive
    # construction
    -> std.net.listen_tcp
  poolserver.set_handler
    fn (server: server_state, handler: handler_fn) -> server_state
    + installs the per-connection handler
    # handler
  poolserver.start
    fn (server: server_state) -> result[server_state, string]
    + spawns the worker pool and begins accepting
    # lifecycle
    -> std.thread.spawn
    -> std.net.accept
  poolserver.submit_conn
    fn (server: server_state, conn: conn_state) -> result[server_state, string]
    + enqueues a connection for a worker to handle
    - returns error when the queue is full; the connection is closed
    # dispatch
    -> std.net.conn_close
  poolserver.stats
    fn (server: server_state) -> pool_stats
    + returns active workers, queued connections, and total handled
    # introspection
  poolserver.shutdown
    fn (server: server_state) -> result[void, string]
    + stops accepting, drains the queue, and joins workers
    # lifecycle
    -> std.thread.join
