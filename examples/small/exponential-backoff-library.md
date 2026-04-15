# Requirement: "an exponential backoff library"

Computes the next delay in an exponential backoff sequence and retries a caller-supplied operation.

std
  std.time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + suspends the current task for ms milliseconds
      # time

backoff
  backoff.new
    fn (initial_ms: i64, max_ms: i64, factor: f64) -> backoff_state
    + creates a backoff starting at initial_ms and capped at max_ms
    ? factor must be > 1; caller is responsible for providing a sensible value
    # construction
  backoff.next
    fn (state: backoff_state) -> tuple[i64, backoff_state]
    + returns the current delay and an advanced state with delay * factor, capped at max_ms
    # scheduling
  backoff.reset
    fn (state: backoff_state) -> backoff_state
    + returns state restored to initial_ms
    # scheduling
  backoff.retry
    fn (state: backoff_state, max_attempts: i32, op: fn() -> result[bool, string]) -> result[bool, string]
    + returns success as soon as op returns success, sleeping between attempts
    - returns the last error after max_attempts failures
    # retry
    -> std.time.sleep_millis
