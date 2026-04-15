# Requirement: "a terminal progress bar for long-running iterations"

Tracks progress, estimates rate and ETA, and renders a bar line. Rendering is separated from output so callers can redirect it.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

progress
  progress.new
    fn (total: i64, width: i32) -> progress_state
    + creates a bar with the given total and render width
    # construction
    -> std.time.now_millis
  progress.update
    fn (state: progress_state, current: i64) -> progress_state
    + advances the bar to the given count and refreshes rate and ETA
    ? rate is smoothed with an exponential moving average
    # update
    -> std.time.now_millis
  progress.render
    fn (state: progress_state) -> string
    + returns a line like "[####    ] 40% 120/300 eta 0:05"
    # render
  progress.finish
    fn (state: progress_state) -> string
    + returns a final line showing total elapsed time
    # finalize
    -> std.time.now_millis
