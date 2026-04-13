# Requirement: "a panic-safe asynchronous function runner"

Runs a zero-argument function on a background task, capturing any panic as a structured error and exposing a join handle.

std
  std.task
    std.task.spawn
      @ (f: func() -> result[string, string]) -> task_handle
      + starts f on a background task and returns a handle
      # concurrency
    std.task.join
      @ (h: task_handle) -> result[string, string]
      + blocks until the task finishes and returns its result
      # concurrency

async_safe
  async_safe.run
    @ (f: func() -> string) -> task_handle
    + spawns f and wraps it so any panic becomes an error result instead of crashing the task
    # spawn
    -> std.task.spawn
  async_safe.wait
    @ (h: task_handle) -> result[string, string]
    + returns the function's return value on success
    - returns an error carrying the panic message when f panicked
    # join
    -> std.task.join
