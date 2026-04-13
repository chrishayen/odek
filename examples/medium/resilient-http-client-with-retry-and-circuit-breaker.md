# Requirement: "an HTTP client with retry and circuit-breaker capabilities"

Wraps a transport with configurable retry policy and a circuit breaker that opens after repeated failures.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      @ (duration: i64) -> void
      + blocks the current task for the given duration
      # time

resilient_http
  resilient_http.new_client
    @ (max_retries: i32, base_backoff_ms: i64, breaker_threshold: i32, breaker_cooldown_ms: i64) -> client_state
    + creates a client with the given retry and breaker settings
    # construction
  resilient_http.execute
    @ (state: client_state, req: http_request, transport: fn(http_request) -> result[http_response, string]) -> tuple[client_state, result[http_response, string]]
    + returns the response, retrying on transient errors with exponential backoff
    + returns fast with a circuit_open error when the breaker is open
    - returns the last error after max_retries attempts fail
    # execution
    -> std.time.now_millis
    -> std.time.sleep_millis
  resilient_http.should_retry
    @ (res: result[http_response, string], attempt: i32, max_retries: i32) -> bool
    + returns true for network errors and 5xx responses when attempts remain
    - returns false for 4xx responses other than 408 and 429
    # retry_policy
  resilient_http.next_backoff
    @ (base_ms: i64, attempt: i32) -> i64
    + returns base * 2^attempt capped at a sane maximum
    # retry_policy
  resilient_http.breaker_record_success
    @ (state: client_state) -> client_state
    + resets the failure counter and closes the breaker
    # circuit_breaker
  resilient_http.breaker_record_failure
    @ (state: client_state, now_ms: i64) -> client_state
    + increments the failure counter and opens the breaker when the threshold is reached
    # circuit_breaker
  resilient_http.breaker_allows
    @ (state: client_state, now_ms: i64) -> bool
    + returns true when the breaker is closed, or when cooldown has elapsed in half-open
    - returns false while the breaker is open and cooldown has not elapsed
    # circuit_breaker
