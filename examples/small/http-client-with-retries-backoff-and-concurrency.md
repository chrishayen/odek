# Requirement: "an HTTP client with retries, backoff, and concurrency"

Wraps a pluggable transport with retry and backoff policy, plus a concurrency helper for parallel requests.

std
  std.time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + sleeps for the given number of milliseconds
      # time
  std.math
    std.math.pow_f64
      fn (base: f64, exponent: f64) -> f64
      + returns base raised to exponent
      # math

http_client
  http_client.new
    fn (max_attempts: i32, base_backoff_millis: i64) -> client_state
    + creates a client configured with retry count and initial backoff
    # construction
  http_client.do_with_retry
    fn (state: client_state, transport_id: string, request: bytes) -> result[bytes, string]
    + invokes the transport, retrying on transient failures with exponential backoff
    - returns the last error after exhausting max_attempts
    # retry
    -> std.time.sleep_millis
    -> std.math.pow_f64
  http_client.do_concurrent
    fn (state: client_state, transport_id: string, requests: list[bytes], max_parallel: i32) -> list[result[bytes, string]]
    + fans out requests in parallel up to max_parallel and preserves input order in results
    # concurrency
