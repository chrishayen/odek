# Requirement: "an interactive shell line-editing library with modern UI features"

Core is a reusable line editor with history, completion, and simple cursor controls.

std
  std.io
    std.io.read_key
      @ () -> result[i32, string]
      + returns the next key code from stdin in raw mode
      - returns error when stdin is closed
      # input
    std.io.write_string
      @ (s: string) -> void
      + writes bytes to stdout without buffering
      # output
  std.term
    std.term.enable_raw_mode
      @ () -> result[term_state, string]
      + switches the terminal into raw mode and returns a restore handle
      - returns error when stdout is not a tty
      # terminal
    std.term.restore
      @ (state: term_state) -> void
      + restores the terminal to its prior mode
      # terminal
    std.term.get_width
      @ () -> i32
      + returns the current terminal column count
      # terminal

line_editor
  line_editor.new
    @ (prompt: string) -> editor_state
    + creates an editor with the given prompt and empty buffer
    # construction
  line_editor.read_line
    @ (state: editor_state) -> result[string, string]
    + returns the committed line when the user presses enter
    + supports left/right cursor movement and backspace
    - returns error on read failure
    # line_editing
    -> std.term.enable_raw_mode
    -> std.term.restore
    -> std.io.read_key
    -> std.io.write_string
  line_editor.set_completer
    @ (state: editor_state, fn: completer_fn) -> editor_state
    + registers a completion callback invoked on tab
    # completion
  line_editor.history_push
    @ (state: editor_state, line: string) -> editor_state
    + appends a line to the in-memory history ring
    # history
  line_editor.history_recall
    @ (state: editor_state, direction: i32) -> optional[string]
    + returns the previous or next history entry relative to the cursor
    - returns none when the history is empty
    # history
  line_editor.redraw
    @ (state: editor_state) -> void
    + repaints the prompt and buffer on the current line
    # rendering
    -> std.term.get_width
    -> std.io.write_string
