# Requirement: "a library that transforms a function to return a future for parallel processing"

Wrap any callable so it runs on a worker pool and returns a handle the caller can await.

std: (all units exist)

channelify
  channelify.new_pool
    @ (workers: i32) -> pool_state
    + returns a pool with the given number of workers ready to accept tasks
    ? workers must be >= 1
    # construction
  channelify.wrap
    @ (pool: pool_state, f: fn(any) -> any) -> fn(any) -> future_handle
    + returns a new function that, when called, schedules the underlying call on the pool
    # wrapping
  channelify.await
    @ (handle: future_handle) -> any
    + blocks until the task completes and returns its result
    # awaiting
  channelify.close_pool
    @ (pool: pool_state) -> void
    + waits for in-flight tasks to finish and releases workers
    # lifecycle
