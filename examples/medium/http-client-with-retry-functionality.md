# Requirement: "enrich an HTTP client with retry functionality"

Wraps a request function with a retry policy that backs off exponentially and honors Retry-After on 429 and 503 responses.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (duration: i64) -> void
      + blocks the current task for the given duration
      # time

http_retry
  http_retry.new_policy
    fn (max_attempts: i32, base_ms: i64, max_ms: i64) -> retry_policy
    + creates a policy with the given attempts and backoff bounds
    # configuration
  http_retry.is_retriable
    fn (status: i32, network_error: bool) -> bool
    + returns true for network errors, 5xx, 408, and 429
    - returns false for 2xx, 3xx, and most 4xx responses
    # policy
  http_retry.backoff_ms
    fn (policy: retry_policy, attempt: i32) -> i64
    + returns base * 2^(attempt-1) capped at max_ms
    # policy
  http_retry.parse_retry_after
    fn (header: string) -> optional[i64]
    + returns a duration in milliseconds for numeric Retry-After values
    - returns none when the header is empty or malformed
    # policy
  http_retry.execute
    fn (policy: retry_policy, send: fn() -> result[http_response, string]) -> result[http_response, string]
    + returns the first non-retriable response or the last attempt result
    + sleeps for Retry-After when the server provides it, otherwise uses computed backoff
    - returns the last error after max_attempts failed attempts
    # execution
    -> std.time.sleep_millis
    -> std.time.now_millis
