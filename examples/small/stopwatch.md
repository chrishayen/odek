# Requirement: "a stopwatch"

Three operations on an opaque state value. Uses std monotonic-clock reads so tests aren't flaky under wall-clock changes.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns monotonic time in nanoseconds, not wall-clock time
      ? monotonic clock only moves forward and is immune to ntp adjustments
      # time

stopwatch
  stopwatch.start
    fn () -> stopwatch_state
    + returns a running stopwatch with the current monotonic time captured
    # lifecycle
    -> std.time.now_nanos
  stopwatch.stop
    fn (s: stopwatch_state) -> stopwatch_state
    + freezes elapsed time and marks the stopwatch stopped
    + stopping an already-stopped stopwatch is a no-op
    # lifecycle
    -> std.time.now_nanos
  stopwatch.elapsed_seconds
    fn (s: stopwatch_state) -> f64
    + returns elapsed seconds as f64
    + works on both running and stopped stopwatches
    # measurement
