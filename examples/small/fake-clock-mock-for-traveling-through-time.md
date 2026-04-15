# Requirement: "a library for travelling through time by mocking the current-time source"

A controllable clock other code can read from, so tests can freeze or advance time.

std: (all units exist)

fake_clock
  fake_clock.new
    fn (epoch_millis: i64) -> clock_state
    + returns a clock frozen at the given unix millisecond
    # construction
  fake_clock.now_millis
    fn (state: clock_state) -> i64
    + returns the current frozen time
    # read
  fake_clock.freeze_at
    fn (state: clock_state, epoch_millis: i64) -> clock_state
    + returns a clock state jumped to the given unix millisecond
    # control
  fake_clock.tick
    fn (state: clock_state, delta_millis: i64) -> clock_state
    + returns a clock state advanced by delta (negative values rewind)
    # control
