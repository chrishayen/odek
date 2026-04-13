# Requirement: "a rate limiting library with token bucket and AIMD strategies"

Two strategy constructors producing interchangeable limiter state.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

rate_limiter
  rate_limiter.new_token_bucket
    @ (rate_per_sec: f64, burst: i32) -> limiter_state
    + creates a token-bucket limiter with the given refill rate and burst
    # construction
  rate_limiter.new_aimd
    @ (initial: f64, inc: f64, dec_factor: f64, max: f64) -> limiter_state
    + creates an AIMD limiter with additive increase and multiplicative decrease
    ? caller invokes on_success and on_failure to adjust the window
    # construction
  rate_limiter.try_acquire
    @ (state: limiter_state) -> tuple[bool, limiter_state]
    + returns (true, new_state) when the request is allowed
    - returns (false, unchanged_state) when the budget is exhausted
    # rate_limiting
    -> std.time.now_millis
  rate_limiter.on_success
    @ (state: limiter_state) -> limiter_state
    + AIMD: additive increase of the allowed rate
    ? no-op for token bucket
    # feedback
  rate_limiter.on_failure
    @ (state: limiter_state) -> limiter_state
    + AIMD: multiplicative decrease of the allowed rate
    ? no-op for token bucket
    # feedback
