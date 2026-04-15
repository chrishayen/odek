# Requirement: "a counting semaphore with timeout on acquire and release"

A fixed-capacity semaphore whose acquire blocks up to a timeout. Time reads go through a thin std utility so tests can drive a deterministic clock.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

semaphore
  semaphore.new
    fn (capacity: i32) -> semaphore_state
    + creates a semaphore with the given number of permits available
    ? capacity must be > 0
    # construction
  semaphore.try_acquire
    fn (state: semaphore_state) -> tuple[bool, semaphore_state]
    + returns (true, new_state) and decrements permits when one is available
    - returns (false, unchanged_state) when no permits remain
    # acquire
  semaphore.acquire_with_timeout
    fn (state: semaphore_state, timeout_ms: i64) -> result[semaphore_state, string]
    + returns the updated state once a permit is obtained within timeout_ms
    - returns error "timeout" when no permit becomes available before the deadline
    # acquire_timeout
    -> std.time.now_millis
  semaphore.release
    fn (state: semaphore_state) -> result[semaphore_state, string]
    + returns state with one permit added back
    - returns error when releasing would exceed the original capacity
    # release
  semaphore.available
    fn (state: semaphore_state) -> i32
    + returns the current number of free permits
    # introspection
