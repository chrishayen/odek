# Requirement: "a collection of utility functions for cancellation and deadline contexts"

Helpers that derive new contexts from existing ones and inspect their lifecycle.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

ctxutil
  ctxutil.with_timeout
    @ (parent: bus_state, timeout_ms: i64) -> bus_state
    + returns a child context that cancels after the timeout elapses
    + child inherits parent cancellation
    # derivation
    -> std.time.now_millis
  ctxutil.with_deadline_at
    @ (parent: bus_state, deadline_unix_ms: i64) -> bus_state
    + returns a child context that cancels at the absolute deadline
    # derivation
    -> std.time.now_millis
  ctxutil.merge
    @ (a: bus_state, b: bus_state) -> bus_state
    + returns a context that cancels when either parent cancels
    # derivation
  ctxutil.remaining_ms
    @ (ctx: bus_state) -> optional[i64]
    + returns milliseconds until deadline, or absent for contexts with no deadline
    # inspection
    -> std.time.now_millis
