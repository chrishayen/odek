# Requirement: "an async I/O library with an event loop and worker thread pool"

An event loop that drives timers, I/O readiness, and deferred callbacks, plus a worker pool for offloading blocking work. std provides the OS polling primitive, monotonic time, and a simple FIFO queue.

std
  std.time
    std.time.now_monotonic_ms
      @ () -> i64
      + returns a monotonically non-decreasing millisecond counter
      # time
  std.os
    std.os.poll_wait
      @ (fds: list[i32], timeout_ms: i32) -> result[list[i32], string]
      + blocks until one or more file descriptors are ready or timeout elapses
      + returns the subset of fds that are ready
      - returns error when a descriptor is invalid
      # io
  std.sync
    std.sync.queue_new
      @ () -> queue_state
      + creates an empty thread-safe FIFO queue
      # concurrency
    std.sync.queue_push
      @ (q: queue_state, item: bytes) -> void
      + appends an item to the queue and wakes any waiter
      # concurrency
    std.sync.queue_pop_blocking
      @ (q: queue_state) -> bytes
      + removes and returns the head item, blocking until one is available
      # concurrency
  std.thread
    std.thread.spawn
      @ (name: string) -> thread_handle
      + spawns a worker thread bound to the given name
      # concurrency

event_loop
  event_loop.new
    @ () -> loop_state
    + creates an empty loop with no handles and an internal worker pool
    # construction
  event_loop.add_timer
    @ (loop: loop_state, delay_ms: i64, tag: string) -> u64
    + schedules a one-shot timer, returning a handle id
    + timers fire in order of deadline, ties broken by insertion order
    # timers
    -> std.time.now_monotonic_ms
  event_loop.cancel_timer
    @ (loop: loop_state, id: u64) -> bool
    + cancels a pending timer, returning true if it was still pending
    - returns false when the id is unknown or already fired
    # timers
  event_loop.watch_fd
    @ (loop: loop_state, fd: i32, interest: string) -> void
    + registers an fd for "read", "write", or "read_write" readiness events
    # io
  event_loop.unwatch_fd
    @ (loop: loop_state, fd: i32) -> void
    + stops watching the fd; pending events for it are dropped
    # io
  event_loop.submit_work
    @ (loop: loop_state, payload: bytes) -> u64
    + enqueues a blocking work item to the worker pool, returning a work id
    # worker_pool
    -> std.sync.queue_push
  event_loop.run_once
    @ (loop: loop_state) -> list[loop_event]
    + runs one iteration: expires due timers, polls I/O, drains completed work
    + returns the list of events in fire order for the caller to dispatch
    # event_loop
    -> std.os.poll_wait
    -> std.time.now_monotonic_ms
  event_loop.run_forever
    @ (loop: loop_state) -> void
    + repeatedly calls run_once until no handles remain
    - exits immediately when there are no timers, fds, or in-flight work
    # event_loop
  event_loop.worker_start
    @ (loop: loop_state, count: i32) -> void
    + spins up the given number of worker threads that consume the work queue
    # worker_pool
    -> std.thread.spawn
    -> std.sync.queue_pop_blocking
