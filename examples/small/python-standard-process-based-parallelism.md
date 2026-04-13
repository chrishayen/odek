# Requirement: "a process-based parallelism library"

Wraps opaque worker handles; the project layer is a small pool API.

std
  std.process
    std.process.spawn
      @ (command: string, args: list[string]) -> result[process_handle, string]
      + launches a child process and returns its handle
      - returns error when the executable cannot be found
      # process
    std.process.wait
      @ (handle: process_handle) -> result[i32, string]
      + blocks until the process exits and returns its exit code
      # process

parallel
  parallel.new_pool
    @ (size: i32) -> pool_state
    + creates a pool with the requested number of worker slots
    ? size is clamped to at least one
    # construction
  parallel.submit
    @ (pool: pool_state, command: string, args: list[string]) -> result[pool_state, string]
    + spawns a worker for the given command when a slot is free
    - returns error when all slots are busy
    # scheduling
    -> std.process.spawn
  parallel.join_all
    @ (pool: pool_state) -> list[i32]
    + waits for every running worker and returns their exit codes in submit order
    # synchronization
    -> std.process.wait
