# Requirement: "a cross-platform asynchronous I/O event loop"

An event loop that multiplexes timers, I/O readiness, and worker-thread completions.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.os
    std.os.poll
      fn (fds: list[i32], timeout_ms: i32) -> list[i32]
      + returns the fds that are ready within the timeout
      - returns an empty list when no fd becomes ready
      # os
    std.os.spawn_thread
      fn (task: task_handle) -> thread_handle
      + runs a task on a worker thread
      # os

evloop
  evloop.new
    fn () -> loop_state
    + creates an empty event loop
    # construction
  evloop.add_timer
    fn (loop: loop_state, delay_ms: i64, cb: callback_handle) -> loop_state
    + schedules a callback to fire after delay_ms
    # timers
    -> std.time.now_millis
  evloop.add_io
    fn (loop: loop_state, fd: i32, event: string, cb: callback_handle) -> loop_state
    + registers a callback for "read" or "write" readiness on fd
    # io
  evloop.submit_work
    fn (loop: loop_state, task: task_handle, done_cb: callback_handle) -> loop_state
    + offloads a task to the thread pool and fires done_cb on completion
    # workers
    -> std.os.spawn_thread
  evloop.run_once
    fn (loop: loop_state) -> loop_state
    + performs one iteration: due timers, then ready I/O, then completed work
    # loop
    -> std.os.poll
    -> std.time.now_millis
  evloop.run
    fn (loop: loop_state) -> loop_state
    + runs until no timers, handles, or pending work remain
    # loop
  evloop.stop
    fn (loop: loop_state) -> loop_state
    + marks the loop to exit after the current iteration
    # control
