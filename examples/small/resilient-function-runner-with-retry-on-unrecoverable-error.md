# Requirement: "a library that runs functions resiliently, catching unrecoverable errors and restarting them"

Wraps a user function with crash capture and a restart policy. The policy is explicit so callers control backoff and attempt limits.

std
  std.time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + blocks the current task for ms milliseconds
      # time
  std.panic
    std.panic.try
      fn (body: fn() -> void) -> optional[string]
      + runs body and returns none on clean exit or the captured error message if it crashed
      # recovery

resilient
  resilient.policy
    fn (max_attempts: i32, initial_backoff_ms: i64, backoff_multiplier: f64) -> restart_policy
    + constructs a policy; max_attempts of 0 means retry forever
    # policy
  resilient.run
    fn (body: fn() -> void, policy: restart_policy) -> result[i32, string]
    + runs body, catching crashes and retrying per policy; returns the number of attempts on clean exit
    - returns the last captured error when max_attempts is exhausted
    + sleeps between attempts using the current backoff, multiplied after each failure
    # supervision
    -> std.panic.try
    -> std.time.sleep_millis
  resilient.run_forever
    fn (body: fn() -> void, initial_backoff_ms: i64, backoff_multiplier: f64) -> void
    + convenience wrapper that retries indefinitely with exponential backoff
    # supervision
    -> std.panic.try
    -> std.time.sleep_millis
