# Requirement: "a circuit breaker with gradual recovery via probabilistic throttling"

Three states (closed, open, half-open) with the half-open phase admitting a gradually increasing fraction of requests.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.rand
    std.rand.uniform
      @ () -> f64
      + returns a uniformly distributed value in [0, 1)
      # random

breaker
  breaker.new
    @ (failure_threshold: i32, open_duration_ms: i64, recovery_window_ms: i64) -> breaker_state
    + constructs a closed breaker with the given thresholds
    # construction
  breaker.allow
    @ (state: breaker_state) -> tuple[bool, breaker_state]
    + returns (true, state) when closed
    + in half-open, admits requests with a probability rising linearly through the recovery window
    - returns (false, state) when open and the open window has not elapsed
    # admission
    -> std.time.now_millis
    -> std.rand.uniform
  breaker.record_success
    @ (state: breaker_state) -> breaker_state
    + resets failure count and closes the breaker when fully recovered
    # feedback
  breaker.record_failure
    @ (state: breaker_state) -> breaker_state
    + increments failures and transitions to open when the threshold is reached
    # feedback
    -> std.time.now_millis
