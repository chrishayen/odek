# Requirement: "a library that returns the current terminal window size"

One function that queries the controlling terminal and returns its column and row count.

std: (all units exist)

terminal_size
  terminal_size.get
    @ () -> result[terminal_dimensions, string]
    + returns the terminal's column and row counts when a tty is attached
    - returns error when stdout is not a terminal
    ? dimensions are sampled at call time; callers re-poll on SIGWINCH
    # terminal
