# Requirement: "a terminal line-rewriter for updating previous output in place"

Produces the control-character stream needed to erase and rewrite the previous render. Actual writing to a terminal is the caller's job.

std: (all units exist)

log_update
  log_update.new
    fn () -> update_state
    + creates an empty state with no previous render
    # construction
  log_update.render
    fn (state: update_state, frame: string) -> tuple[string, update_state]
    + returns a control sequence that clears the previous frame and prints the new one
    + returns the new frame unchanged when there was no previous render
    ? the clear sequence moves the cursor up by the previous frame's line count and erases each line
    # rendering
  log_update.clear
    fn (state: update_state) -> tuple[string, update_state]
    + returns a control sequence that erases the previous frame and resets state
    - returns an empty string when there is no previous frame
    # rendering
  log_update.done
    fn (state: update_state) -> update_state
    + resets state so the next render starts fresh without clearing the current frame
    # finalize
