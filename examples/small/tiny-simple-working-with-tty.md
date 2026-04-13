# Requirement: "a small library for working with a TTY"

Lets callers move the cursor, clear regions, and toggle raw mode on a terminal. Bytes are emitted via a thin std writer so output can be captured in tests.

std
  std.io
    std.io.write_stdout
      @ (data: bytes) -> void
      + writes the bytes to standard output
      # io

tty
  tty.clear_screen
    @ () -> void
    + writes the ANSI sequence that clears the entire screen and homes the cursor
    # terminal
    -> std.io.write_stdout
  tty.move_cursor
    @ (row: i32, col: i32) -> void
    + writes the ANSI sequence that moves the cursor to (row, col), 1-indexed
    ? coordinates below 1 are clamped to 1
    # terminal
    -> std.io.write_stdout
  tty.set_raw_mode
    @ (enable: bool) -> result[void, string]
    + disables line buffering and echo when enable is true
    + restores the previous terminal mode when enable is false
    - returns error when the standard input is not a terminal
    # terminal
