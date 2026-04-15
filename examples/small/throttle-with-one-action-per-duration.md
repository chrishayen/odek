# Requirement: "a throttle that allows at most one action per configurable duration"

Two project functions: construction and a gated call. Time reads go through a thin std primitive so tests can substitute a deterministic clock.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

throttle
  throttle.new
    fn (interval_ms: i64) -> throttle_state
    + creates a throttle with the given minimum interval and no prior call
    ? first call is always allowed
    # construction
  throttle.try_fire
    fn (state: throttle_state) -> tuple[bool, throttle_state]
    + returns (true, new_state) when enough time has passed since the last fire, updating the last-fire timestamp
    - returns (false, unchanged_state) when the throttle is still cooling down
    # gating
    -> std.time.now_millis
