# Requirement: "an async concurrency runtime with structured task nurseries and cancellation"

A nursery spawns child tasks whose lifetimes are bounded by the nursery; cancelling the nursery cancels every child. Task execution goes through a std scheduler primitive.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.scheduler
    std.scheduler.yield_now
      fn () -> void
      + yields control back to the scheduler so other tasks can run
      # scheduling

async_runtime
  async_runtime.new_nursery
    fn () -> nursery_state
    + returns an open nursery with no children and no cancellation
    # construction
  async_runtime.spawn
    fn (nursery: nursery_state, task_id: string) -> result[nursery_state, string]
    + registers a child task id with the nursery
    - returns error when the nursery has already been cancelled
    # tasks
  async_runtime.cancel
    fn (nursery: nursery_state) -> nursery_state
    + marks the nursery and all its children as cancelled
    # cancellation
  async_runtime.is_cancelled
    fn (nursery: nursery_state, task_id: string) -> bool
    + returns true when the nursery is cancelled or the named child was cancelled individually
    # cancellation
  async_runtime.wait_until
    fn (deadline_ms: i64) -> result[void, string]
    + yields repeatedly until wall-clock time passes the deadline
    - returns error on a negative deadline
    # timers
    -> std.time.now_millis
    -> std.scheduler.yield_now
  async_runtime.join
    fn (nursery: nursery_state) -> result[void, string]
    + yields until every child of the nursery has finished
    - returns error when the nursery was cancelled before joining
    # tasks
    -> std.scheduler.yield_now
