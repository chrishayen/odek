# Requirement: "an event-driven networking framework that accepts connections, runs an event loop per worker, and delivers read, write, and close events to user handlers"

Manages listeners, schedules connections across workers, drains read buffers, buffers writes, and notifies user handlers via registered callbacks.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener_handle, string]
      + binds a TCP listener on host and port
      - returns error when the address cannot be bound
      # networking
    std.net.accept_nonblocking
      fn (listener: listener_handle) -> result[optional[conn_handle], string]
      + returns a new connection when one is pending, none otherwise
      - returns error on listener failure
      # networking
    std.net.read_nonblocking
      fn (conn: conn_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes without blocking
      + returns empty bytes when no data is currently available
      - returns error on read failure
      # networking
    std.net.write_nonblocking
      fn (conn: conn_handle, data: bytes) -> result[i32, string]
      + writes bytes and returns the count written
      - returns error on write failure
      # networking
    std.net.close
      fn (conn: conn_handle) -> void
      + closes the connection
      # networking
  std.io
    std.io.poll
      fn (fds: list[i32], timeout_ms: i32) -> result[list[i32], string]
      + returns the indices of fds that are ready for i/o
      - returns error on poll failure
      # io

event_loop_server
  event_loop_server.new_server
    fn (num_workers: i32) -> server_state
    + creates a server with the given worker count and empty listeners
    # construction
  event_loop_server.bind
    fn (state: server_state, host: string, port: i32) -> result[server_state, string]
    + binds a new listener and adds it to the server
    - returns error when the address cannot be bound
    # setup
    -> std.net.listen_tcp
  event_loop_server.assign_worker
    fn (state: server_state, conn_id: i64) -> i32
    + returns the worker index that owns the connection id using stable hashing
    # scheduling
  event_loop_server.new_worker
    fn (index: i32) -> worker_state
    + creates an empty worker with no connections
    # construction
  event_loop_server.register_conn
    fn (worker: worker_state, conn: conn_handle) -> worker_state
    + adds a connection with empty read and write buffers
    # worker
  event_loop_server.unregister_conn
    fn (worker: worker_state, conn_id: i64) -> worker_state
    + removes the connection and discards its buffers
    # worker
  event_loop_server.enqueue_write
    fn (worker: worker_state, conn_id: i64, data: bytes) -> worker_state
    + appends bytes to the connection's pending-write buffer
    # buffering
  event_loop_server.drain_reads
    fn (worker: worker_state, conn_id: i64) -> result[tuple[worker_state, bytes], string]
    + reads all currently available bytes and returns them with the updated state
    - returns error on read failure
    # io
    -> std.net.read_nonblocking
  event_loop_server.flush_writes
    fn (worker: worker_state, conn_id: i64) -> result[worker_state, string]
    + writes as many pending bytes as possible and removes them from the buffer
    - returns error on write failure
    # io
    -> std.net.write_nonblocking
  event_loop_server.close_conn
    fn (worker: worker_state, conn_id: i64) -> worker_state
    + flushes remaining writes, closes the connection, and unregisters it
    # lifecycle
    -> std.net.close
  event_loop_server.new_handler_set
    fn () -> handler_set
    + creates an empty handler set
    # handlers
  event_loop_server.set_on_open
    fn (set: handler_set, handler_id: string) -> handler_set
    + registers the open handler
    # handlers
  event_loop_server.set_on_read
    fn (set: handler_set, handler_id: string) -> handler_set
    + registers the read handler
    # handlers
  event_loop_server.set_on_close
    fn (set: handler_set, handler_id: string) -> handler_set
    + registers the close handler
    # handlers
  event_loop_server.dispatch_events
    fn (worker: worker_state, handlers: handler_set, invoker: handler_invoker) -> result[worker_state, string]
    + polls connections and invokes handlers for ready reads, writes, and closes
    - returns error when polling fails
    # dispatch
    -> std.io.poll
  event_loop_server.accept_pending
    fn (state: server_state) -> result[server_state, string]
    + polls every listener once for a new connection and registers it on a worker
    - returns error on accept failure
    # acceptance
    -> std.net.accept_nonblocking
