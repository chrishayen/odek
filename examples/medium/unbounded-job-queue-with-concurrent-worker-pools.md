# Requirement: "an unbounded job queue with concurrent worker pools"

Producers enqueue jobs; a pool of workers drains them concurrently. The queue has no fixed capacity.

std
  std.sync
    std.sync.spawn
      @ (task: fn() -> void) -> void
      + runs the task on a background worker
      # concurrency
    std.sync.mutex_lock
      @ (m: mutex_handle) -> void
      + acquires exclusive access to the protected region
      # concurrency
    std.sync.mutex_unlock
      @ (m: mutex_handle) -> void
      + releases exclusive access
      # concurrency
    std.sync.condvar_wait
      @ (cv: condvar_handle, m: mutex_handle) -> void
      + atomically releases the mutex and waits for notification
      # concurrency
    std.sync.condvar_notify_one
      @ (cv: condvar_handle) -> void
      + wakes one waiting worker
      # concurrency

kyoo
  kyoo.new
    @ (worker_count: i32) -> kyoo_state
    + returns a queue state with the given number of workers not yet started
    ? internal queue grows without a fixed upper bound
    # construction
  kyoo.start
    @ (q: kyoo_state) -> void
    + spawns the configured workers, each looping over pending jobs
    # lifecycle
    -> std.sync.spawn
  kyoo.submit
    @ (q: kyoo_state, job: fn() -> void) -> void
    + appends the job to the back of the queue
    + wakes one worker when at least one was idle
    # enqueue
    -> std.sync.mutex_lock
    -> std.sync.mutex_unlock
    -> std.sync.condvar_notify_one
  kyoo.drain
    @ (q: kyoo_state) -> void
    + blocks until every submitted job has completed
    # synchronization
    -> std.sync.condvar_wait
  kyoo.stop
    @ (q: kyoo_state) -> void
    + signals every worker to exit after the current job
    # lifecycle
  kyoo.pending_count
    @ (q: kyoo_state) -> i64
    + returns the number of jobs waiting to be picked up
    # introspection
