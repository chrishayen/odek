# Requirement: "a terminal spinner for long-running command-line operations"

Manages spinner state and produces frames for the caller to print; the library does not own I/O.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time

spinner
  spinner.new
    @ (frames: list[string], interval_ms: i32) -> spinner_state
    + creates a spinner with the given animation frames and frame interval
    ? frames list must be non-empty; interval_ms controls advance rate
    # construction
  spinner.current_frame
    @ (state: spinner_state) -> tuple[string, spinner_state]
    + returns the frame appropriate for the current time and advances state
    + frame index wraps around the frames list
    # animation
    -> std.time.now_millis
  spinner.with_message
    @ (state: spinner_state, message: string) -> spinner_state
    + attaches a suffix message to display alongside each frame
    # decoration
  spinner.render
    @ (state: spinner_state) -> tuple[string, spinner_state]
    + returns the formatted "frame message" string for the current tick
    # rendering
    -> spinner.current_frame
