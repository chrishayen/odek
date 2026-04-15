# Requirement: "a terminal progress bar renderer"

Tracks progress state and renders a fixed-width bar string. The caller is responsible for writing the string to the terminal.

std: (all units exist)

progress
  progress.new
    fn (total: i64, width: i32) -> bar_state
    + creates a bar with zero progress and the given total and rendered width
    # construction
  progress.advance
    fn (state: bar_state, delta: i64) -> bar_state
    + adds delta to the current value, clamped to total
    # update
  progress.set
    fn (state: bar_state, value: i64) -> bar_state
    + sets the current value, clamped to [0, total]
    # update
  progress.render
    fn (state: bar_state) -> string
    + returns a line like "[====>     ] 40% (40/100)"
    + shows "[==========] 100%" when complete
    - shows "[>         ] 0%" for a fresh bar with total > 0
    # rendering
