# Requirement: "a progress bar with real-time throughput and eta"

The library computes progress state; rendering animation frames is the caller's job.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current wall-clock time in milliseconds
      # time

progress
  progress.new
    fn (total: i64) -> progress_state
    + creates a bar with the given total and zero items done
    ? total of 0 means indeterminate progress
    # construction
    -> std.time.now_millis
  progress.advance
    fn (state: progress_state, delta: i64) -> progress_state
    + adds delta to the completed count
    + clamps completed to total when total is known
    # updating
    -> std.time.now_millis
  progress.snapshot
    fn (state: progress_state) -> progress_snapshot
    + returns fraction, items/sec throughput, and eta in seconds
    + throughput is a moving average over recent advances
    - eta is -1 when total is 0 or throughput is 0
    # reporting
    -> std.time.now_millis
