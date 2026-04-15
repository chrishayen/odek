# Requirement: "a library of ANSI escape codes for manipulating a terminal"

Returns well-known terminal control sequences as strings. No IO — callers choose what to write and where.

std: (all units exist)

ansi_escapes
  ansi_escapes.cursor_to
    fn (row: i32, col: i32) -> string
    + returns the escape sequence to move the cursor to 1-indexed (row, col)
    # cursor
  ansi_escapes.cursor_move
    fn (dx: i32, dy: i32) -> string
    + returns the escape sequence to move the cursor by a relative offset
    + returns an empty string when both offsets are 0
    # cursor
  ansi_escapes.clear_screen
    fn () -> string
    + returns the escape sequence that clears the screen and moves the cursor home
    # screen
  ansi_escapes.clear_line
    fn () -> string
    + returns the escape sequence that erases the current line
    # screen
  ansi_escapes.hide_cursor
    fn () -> string
    + returns the escape sequence that hides the cursor
    # cursor
  ansi_escapes.show_cursor
    fn () -> string
    + returns the escape sequence that shows the cursor
    # cursor
