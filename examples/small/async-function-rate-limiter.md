# Requirement: "a library to limit the execution rate of an async function"

Wraps a function so concurrent/queued calls respect a maximum rate. Time reads go through a thin std primitive.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

rate_limited
  rate_limited.new
    fn (max_calls: i32, per_millis: i64) -> limiter_state
    + creates a limiter allowing max_calls within any per_millis window
    ? window is a sliding timestamp queue, not a fixed bucket
    # construction
  rate_limited.acquire
    fn (state: limiter_state) -> tuple[bool, i64, limiter_state]
    + returns (true, 0, new_state) when a slot is free and records the call
    + returns (false, wait_millis, unchanged_state) when saturated, with millis until the next slot frees
    # rate_limiting
    -> std.time.now_millis
  rate_limited.release
    fn (state: limiter_state) -> limiter_state
    + removes expired timestamps from the window
    # cleanup
    -> std.time.now_millis
