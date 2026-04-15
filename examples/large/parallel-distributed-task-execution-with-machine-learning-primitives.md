# Requirement: "a library for parallel and distributed task execution with machine learning primitives"

A task graph scheduler that executes work across a local worker pool and tracks futures. Out-of-scope: actually running on remote nodes — the library exposes a pluggable transport so callers can wire that up.

std
  std.thread
    std.thread.spawn
      fn (work: thread_fn) -> result[thread_handle, string]
      + starts a worker thread running the function
      - returns error when the OS refuses the thread
      # threading
    std.thread.join
      fn (handle: thread_handle) -> void
      + waits for the thread to exit
      # threading
  std.sync
    std.sync.mutex_new
      fn () -> mutex_state
      + creates an unlocked mutex
      # sync
    std.sync.mutex_lock
      fn (mutex: mutex_state) -> void
      + blocks until the mutex is acquired
      # sync
    std.sync.mutex_unlock
      fn (mutex: mutex_state) -> void
      + releases the mutex
      # sync
    std.sync.cond_new
      fn () -> cond_state
      + creates a condition variable
      # sync
    std.sync.cond_wait
      fn (cond: cond_state, mutex: mutex_state) -> void
      + atomically releases and waits
      # sync
    std.sync.cond_broadcast
      fn (cond: cond_state) -> void
      + wakes all waiters
      # sync
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

taskflow
  taskflow.new
    fn (worker_count: i32) -> runtime_state
    + creates a runtime with the requested number of workers
    # construction
    -> std.thread.spawn
    -> std.sync.mutex_new
    -> std.sync.cond_new
  taskflow.submit
    fn (runtime: runtime_state, task: task_fn, inputs: list[bytes]) -> future_id
    + queues a task for execution and returns its future id
    # scheduling
    -> std.sync.mutex_lock
    -> std.sync.cond_broadcast
    -> std.sync.mutex_unlock
  taskflow.depend
    fn (runtime: runtime_state, future: future_id, upstream: list[future_id]) -> runtime_state
    + declares that the future must wait for upstream futures to complete
    # graph
  taskflow.get
    fn (runtime: runtime_state, future: future_id) -> result[bytes, string]
    + blocks until the future resolves and returns its value
    - returns error when the task failed
    # retrieval
    -> std.sync.mutex_lock
    -> std.sync.cond_wait
    -> std.sync.mutex_unlock
  taskflow.cancel
    fn (runtime: runtime_state, future: future_id) -> runtime_state
    + marks the future cancelled; pending dependents fail
    # lifecycle
  taskflow.map
    fn (runtime: runtime_state, task: task_fn, items: list[bytes]) -> list[future_id]
    + submits one task per item and returns their futures
    # scheduling
  taskflow.reduce
    fn (runtime: runtime_state, futures: list[future_id], combine: combine_fn) -> future_id
    + schedules a reduction that depends on all of the input futures
    # scheduling
  taskflow.set_transport
    fn (runtime: runtime_state, transport: transport_fn) -> runtime_state
    + installs a pluggable transport used to ship tasks to remote workers
    ? when no transport is set, all tasks run locally
    # transport
  taskflow.stats
    fn (runtime: runtime_state) -> runtime_stats
    + returns counts of pending, running, completed, and failed tasks
    # introspection
    -> std.time.now_millis
  taskflow.shutdown
    fn (runtime: runtime_state) -> void
    + stops accepting new tasks and joins all workers
    # lifecycle
    -> std.sync.cond_broadcast
    -> std.thread.join
