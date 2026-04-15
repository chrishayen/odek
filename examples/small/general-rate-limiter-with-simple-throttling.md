# Requirement: "a general rate limiter with simple throttling"

Queue-style throttle: callers submit jobs and the limiter hands them back when it is their turn. Jobs run in the caller; the limiter only gates the timing.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

throttle
  throttle.new
    fn (max_concurrent: i32, min_interval_ms: i64) -> throttle_state
    + creates a throttle with a concurrency cap and minimum spacing between job starts
    # construction
  throttle.submit
    fn (state: throttle_state, job_id: string) -> throttle_state
    + enqueues a job by id at the tail of the queue
    # enqueue
  throttle.next_ready
    fn (state: throttle_state) -> tuple[optional[string], throttle_state]
    + returns the next job that is allowed to start given the concurrency cap and interval
    + returns none when no job is ready yet
    # scheduling
    -> std.time.now_millis
  throttle.done
    fn (state: throttle_state, job_id: string) -> throttle_state
    + marks a running job as complete, freeing a concurrency slot
    - no-op when the id is unknown
    # completion
