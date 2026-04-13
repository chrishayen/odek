# Requirement: "a retry-on-error library with flexible backoff strategies and caps"

The project layer composes backoff policies; std provides timing and randomness primitives.

std
  std.time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + blocks the caller for the given number of milliseconds
      # time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.uniform_f64
      @ () -> f64
      + returns a uniformly distributed float in [0.0, 1.0)
      # randomness

retry
  retry.constant_backoff
    @ (delay_ms: i64) -> backoff_policy
    + produces a policy that always returns delay_ms
    # policy_construction
  retry.exponential_backoff
    @ (base_ms: i64, factor: f64, cap_ms: i64) -> backoff_policy
    + produces a policy where each attempt multiplies the delay by factor, capped at cap_ms
    ? attempt 0 returns base_ms; attempt n returns min(base_ms * factor^n, cap_ms)
    # policy_construction
  retry.with_jitter
    @ (inner: backoff_policy, jitter_ratio: f64) -> backoff_policy
    + wraps a policy so each delay is multiplied by (1 - jitter_ratio + 2*jitter_ratio*random)
    - jitter_ratio outside [0, 1] is clamped
    # policy_construction
    -> std.random.uniform_f64
  retry.delay_for
    @ (policy: backoff_policy, attempt: i32) -> i64
    + returns the delay in milliseconds for the given attempt index
    # policy_query
  retry.run
    @ (policy: backoff_policy, max_attempts: i32, op: fn() -> result[bytes, string]) -> result[bytes, string]
    + calls op; returns the first Ok result
    + sleeps according to policy.delay_for(attempt) between retries
    - returns the last error after max_attempts exhausted
    # retry_loop
    -> std.time.sleep_millis
