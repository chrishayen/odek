# Requirement: "a retry mechanism that repeats an action until it succeeds"

Functional retry with pluggable backoff. The action is a caller-supplied closure returning a result.

std
  std.time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + blocks for the specified number of milliseconds
      # time

retry
  retry.forever
    fn (action: fn[void, result[bytes, string]], backoff: fn[i32, i64]) -> bytes
    + calls action repeatedly until it returns Ok
    + sleeps backoff(attempt) milliseconds between attempts
    # retry
    -> std.time.sleep_millis
  retry.with_limit
    fn (action: fn[void, result[bytes, string]], backoff: fn[i32, i64], max_attempts: i32) -> result[bytes, string]
    + calls action up to max_attempts times
    - returns the last error when all attempts fail
    # retry
    -> std.time.sleep_millis
  retry.constant_backoff
    fn (delay_ms: i64) -> fn[i32, i64]
    + returns a backoff function that always yields delay_ms
    # backoff
  retry.exponential_backoff
    fn (base_ms: i64, factor: f64, cap_ms: i64) -> fn[i32, i64]
    + returns a backoff function that grows geometrically and clamps at cap_ms
    # backoff
