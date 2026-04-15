# Requirement: "a high-level async concurrency and networking framework"

A structured-concurrency task system with cancellation, task groups, synchronization primitives, and async networking.

std
  std.async
    std.async.spawn
      fn (op: fn() -> bytes) -> task_handle
      + schedules op on the runtime and returns a handle to its future result
      # async
    std.async.await
      fn (handle: task_handle) -> result[bytes, string]
      + blocks the current async context until the task completes
      - returns error when the task was cancelled
      # async
    std.async.sleep
      fn (millis: i64) -> void
      + suspends the current task for the given duration
      # async
    std.async.cancel
      fn (handle: task_handle) -> void
      + requests cancellation of a running task
      # async
  std.net
    std.net.tcp_connect
      fn (host: string, port: i32) -> result[tcp_conn, string]
      + opens a TCP connection
      # networking
    std.net.tcp_listen
      fn (host: string, port: i32) -> result[tcp_listener, string]
      + binds and listens for incoming TCP connections
      # networking
    std.net.tcp_accept
      fn (listener: tcp_listener) -> result[tcp_conn, string]
      + accepts the next pending connection
      # networking
    std.net.tcp_read
      fn (conn: tcp_conn, max_bytes: i32) -> result[bytes, string]
      + reads up to max_bytes from the connection
      # networking
    std.net.tcp_write
      fn (conn: tcp_conn, data: bytes) -> result[i32, string]
      + writes data and returns the number of bytes written
      # networking

anyio
  anyio.run
    fn (main: fn() -> bytes) -> result[bytes, string]
    + starts the runtime, runs main to completion, and shuts down
    # runtime
    -> std.async.spawn
    -> std.async.await
  anyio.task_group
    fn () -> task_group_state
    + creates a task group that owns spawned children
    # structured_concurrency
  anyio.group_spawn
    fn (group: task_group_state, op: fn() -> bytes) -> task_group_state
    + adds op to the group; its lifetime is bounded by the group
    # structured_concurrency
    -> std.async.spawn
  anyio.group_wait
    fn (group: task_group_state) -> result[list[bytes], string]
    + waits for all tasks in the group to finish
    - if any task fails, the remaining tasks are cancelled and the first error is returned
    # structured_concurrency
    -> std.async.await
    -> std.async.cancel
  anyio.cancel_scope
    fn (deadline_ms: i64) -> cancel_scope_state
    + creates a cancellation scope with an absolute deadline
    # cancellation
  anyio.with_timeout
    fn (millis: i64, op: fn() -> bytes) -> result[bytes, string]
    + runs op subject to a timeout
    - returns a timeout error when the deadline is reached
    # cancellation
    -> anyio.cancel_scope
  anyio.sleep
    fn (millis: i64) -> void
    + suspends the current task cooperatively
    # async
    -> std.async.sleep
  anyio.new_semaphore
    fn (permits: i32) -> semaphore_state
    + creates a counting semaphore with permits available slots
    # synchronization
  anyio.acquire
    fn (sem: semaphore_state) -> semaphore_state
    + blocks until a permit is free, then decrements
    # synchronization
  anyio.release
    fn (sem: semaphore_state) -> semaphore_state
    + increments the permit count and wakes a waiter
    # synchronization
  anyio.open_tcp_stream
    fn (host: string, port: i32) -> result[tcp_conn, string]
    + opens an async TCP stream
    # networking
    -> std.net.tcp_connect
  anyio.serve_tcp
    fn (host: string, port: i32, handler: fn(conn: tcp_conn) -> void) -> result[void, string]
    + accepts connections in a loop, spawning handler per client
    # networking
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.async.spawn
