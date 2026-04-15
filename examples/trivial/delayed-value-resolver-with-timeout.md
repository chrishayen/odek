# Requirement: "a deferred value that resolves after a specified amount of time"

Pure construction of a deferred value; the caller's runtime is responsible for actually waiting.

std: (all units exist)

delay
  delay.after_millis
    fn (ms: i64) -> deferred_state
    + returns a deferred value that becomes ready ms milliseconds after creation
    ? negative values are clamped to zero
    # construction
  delay.is_ready
    fn (d: deferred_state, now_millis: i64) -> bool
    + returns true when now is at or after the deferred's ready time
    - returns false when now is before the ready time
    # inspection
