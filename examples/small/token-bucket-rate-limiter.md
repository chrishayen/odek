# Requirement: "a token bucket rate limiter"

Two project functions. Time reads go through a thin std utility so tests can substitute a deterministic clock without re-implementing the bucket.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

rate_limiter
  rate_limiter.new
    fn (rate_per_sec: f64, burst: i32) -> rate_limiter_state
    + creates a limiter with the given refill rate and burst capacity
    ? tokens accumulate as f64 between calls so fractional refills compound correctly
    # construction
  rate_limiter.try_acquire
    fn (state: rate_limiter_state) -> tuple[bool, rate_limiter_state]
    + returns (true, new_state) when a token is available and consumes one
    + refills tokens based on elapsed time since the last call
    - returns (false, unchanged_state) when the bucket is empty
    # rate_limiting
    -> std.time.now_millis
