# Requirement: "a brute-force and request-flood protection library"

Tracks requests per key and decides whether to allow, throttle, or block based on recent activity.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

protect
  protect.new
    fn (max_attempts: i32, window_seconds: i64, block_seconds: i64) -> protect_state
    + creates a protector with a sliding window and block duration
    # construction
  protect.record_attempt
    fn (state: protect_state, key: string) -> protect_state
    + adds a timestamped attempt for the given key
    # tracking
    -> std.time.now_seconds
  protect.is_blocked
    fn (state: protect_state, key: string) -> bool
    + returns true when the key has exceeded max_attempts in the window or is still within block duration
    - returns false when attempts are below the threshold
    # decision
    -> std.time.now_seconds
  protect.reset
    fn (state: protect_state, key: string) -> protect_state
    + clears all attempts and blocks for the key
    # administration
  protect.prune_expired
    fn (state: protect_state) -> protect_state
    + removes entries whose last activity falls outside both windows
    # maintenance
    -> std.time.now_seconds
