# Requirement: "an asynchronous I/O event loop with coroutines and tasks"

A single-threaded event loop that drives ready tasks, schedules timers, and wakes tasks on I/O readiness.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.io
    std.io.poll_ready
      fn (fds: list[i32], timeout_ms: i64) -> list[io_event]
      + returns the set of file descriptors that are ready for read or write
      - returns empty list on timeout
      # io

asyncio
  asyncio.new_loop
    fn () -> loop_state
    + creates an empty event loop
    # construction
  asyncio.spawn
    fn (state: loop_state, coro: coroutine) -> tuple[loop_state, task_id]
    + schedules a coroutine as a new task and returns its id
    # tasks
  asyncio.cancel
    fn (state: loop_state, task: task_id) -> loop_state
    + marks a task for cancellation on its next resume
    - no-op when the task is already finished
    # tasks
  asyncio.call_later
    fn (state: loop_state, delay_ms: i64, callback: callback_fn) -> loop_state
    + schedules a callback to fire after delay_ms
    # timers
    -> std.time.now_millis
  asyncio.wait_readable
    fn (state: loop_state, task: task_id, fd: i32) -> loop_state
    + suspends a task until fd is readable
    # io
  asyncio.wait_writable
    fn (state: loop_state, task: task_id, fd: i32) -> loop_state
    + suspends a task until fd is writable
    # io
  asyncio.run_once
    fn (state: loop_state) -> loop_state
    + advances the loop one step: fires due timers, drains I/O events, resumes ready tasks
    # loop
    -> std.io.poll_ready
    -> std.time.now_millis
  asyncio.run_until_complete
    fn (state: loop_state, task: task_id) -> result[loop_state, string]
    + runs the loop until the given task finishes
    - returns error when the task was cancelled before completion
    # loop
