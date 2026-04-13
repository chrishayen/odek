# Requirement: "a retry helper that re-runs an operation until it succeeds or a budget is exhausted"

Computes retry decisions from attempt count and elapsed time; the caller drives the loop.

std
  std.time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + suspends the caller for the given number of milliseconds
      # time

retry
  retry.new_policy
    @ (max_attempts: i32, base_delay_ms: i64, max_delay_ms: i64) -> retry_policy
    + returns a policy with exponential backoff between base and max
    # construction
  retry.next_delay
    @ (policy: retry_policy, attempt: i32) -> optional[i64]
    + returns the delay to wait before attempt number N
    - returns none when attempt >= max_attempts
    ? delay grows as base * 2^(attempt-1), capped at max_delay_ms
    # scheduling
  retry.wait_before
    @ (policy: retry_policy, attempt: i32) -> bool
    + sleeps the computed delay and returns true, or returns false when retries are exhausted
    # scheduling
    -> std.time.sleep_millis
