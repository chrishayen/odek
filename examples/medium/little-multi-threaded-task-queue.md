# Requirement: "a multi-threaded task queue"

Producers enqueue tasks; a worker pool dequeues and runs them. Scheduling and retries live in the project layer; concurrency and time are std primitives.

std
  std.sync
    std.sync.spawn
      @ (fn_id: i32, arg: bytes) -> result[i64, string]
      + spawns a worker thread executing the registered function, returns a handle id
      # concurrency
    std.sync.join
      @ (handle: i64) -> result[void, string]
      + waits for a worker to finish
      # concurrency
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      @ (millis: i64) -> void
      + blocks the caller for the given duration
      # time

task_queue
  task_queue.new
    @ (worker_count: i32) -> queue_state
    + creates a queue with the given worker count
    - returns an empty queue when worker_count is 0
    # construction
  task_queue.enqueue
    @ (state: queue_state, fn_id: i32, payload: bytes, max_retries: i32) -> queue_state
    + appends a task to the pending list
    # producer
  task_queue.start
    @ (state: queue_state) -> queue_state
    + spawns workers that loop dequeue-and-run until stop is requested
    # lifecycle
    -> std.sync.spawn
  task_queue.dequeue
    @ (state: queue_state) -> tuple[optional[task], queue_state]
    + removes and returns the oldest pending task, or none when empty
    # consumer
  task_queue.mark_done
    @ (state: queue_state, task_id: i64) -> queue_state
    + records a task as completed
    # bookkeeping
    -> std.time.now_millis
  task_queue.mark_failed
    @ (state: queue_state, task_id: i64, reason: string) -> queue_state
    + increments the task's retry count and re-enqueues it unless max_retries reached
    # retry
  task_queue.stop
    @ (state: queue_state) -> result[void, string]
    + signals workers to exit and joins them
    # lifecycle
    -> std.sync.join
