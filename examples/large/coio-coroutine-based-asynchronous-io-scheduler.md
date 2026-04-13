# Requirement: "a coroutine-based asynchronous I/O scheduler"

An event loop multiplexes ready I/O sources onto a pool of lightweight tasks. The project runes manage task lifecycle, registration, and the run loop.

std
  std.io
    std.io.nonblocking_read
      @ (fd: i32, max: i32) -> result[bytes, string]
      + reads up to max bytes without blocking
      - returns a would-block error when no data is ready
      # io
    std.io.nonblocking_write
      @ (fd: i32, data: bytes) -> result[i32, string]
      + writes without blocking and returns the number of bytes accepted
      - returns a would-block error when the kernel buffer is full
      # io
  std.poll
    std.poll.new
      @ () -> poller
      + creates an empty I/O readiness poller
      # poll
    std.poll.register
      @ (p: poller, fd: i32, interest: i32) -> result[void, string]
      + adds a file descriptor with a read/write interest mask
      # poll
    std.poll.deregister
      @ (p: poller, fd: i32) -> result[void, string]
      + removes a file descriptor
      # poll
    std.poll.wait
      @ (p: poller, timeout_ms: i32) -> list[tuple[i32,i32]]
      + blocks until at least one registered fd is ready or the timeout elapses
      + returns pairs of (fd, ready_mask)
      # poll
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

coio
  coio.new
    @ () -> coio_state
    + creates a scheduler with an empty task table and a fresh poller
    # construction
    -> std.poll.new
  coio.spawn
    @ (state: coio_state, body: coroutine) -> tuple[task_id, coio_state]
    + registers a new coroutine and returns its id
    # tasks
  coio.register_io
    @ (state: coio_state, tid: task_id, fd: i32, interest: i32) -> result[coio_state, string]
    + associates a file descriptor with a waiting task
    # tasks
    -> std.poll.register
  coio.read_async
    @ (state: coio_state, tid: task_id, fd: i32, max: i32) -> result[bytes, string]
    + reads when ready, otherwise parks the task until read-readiness
    - returns error on closed descriptor or underlying io error
    # io
    -> std.io.nonblocking_read
  coio.write_async
    @ (state: coio_state, tid: task_id, fd: i32, data: bytes) -> result[i32, string]
    + writes what fits, otherwise parks the task until write-readiness
    - returns error on closed descriptor or underlying io error
    # io
    -> std.io.nonblocking_write
  coio.sleep_async
    @ (state: coio_state, tid: task_id, ms: i64) -> coio_state
    + parks the task in a timer wheel until ms elapse
    # timers
    -> std.time.now_millis
  coio.run_ready
    @ (state: coio_state) -> coio_state
    + advances every runnable task by one step until all are parked
    # loop
  coio.dispatch_ready_io
    @ (state: coio_state, events: list[tuple[i32,i32]]) -> coio_state
    + wakes tasks waiting on the reported file descriptors
    # loop
  coio.tick
    @ (state: coio_state, max_wait_ms: i32) -> coio_state
    + runs all ready tasks, polls, and wakes fired timers
    # loop
    -> std.poll.wait
    -> std.time.now_millis
  coio.run_until_idle
    @ (state: coio_state) -> coio_state
    + ticks until no tasks remain and no I/O is pending
    # loop
  coio.join
    @ (state: coio_state, tid: task_id) -> result[void, string]
    + drives the loop until the given task completes
    - returns error when the task id is unknown
    # tasks
