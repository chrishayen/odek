# Requirement: "a readline-like line editor for interactive terminals"

Maintains an editable line buffer with a cursor, history navigation, and a tab-completion hook. Raw terminal I/O stays in std so the editor is testable.

std
  std.term
    std.term.read_key
      fn () -> result[key_event, string]
      + returns the next key event from the terminal
      - returns error when the terminal is closed
      # terminal
    std.term.write
      fn (s: string) -> result[bool, string]
      + writes a string to the terminal
      # terminal
    std.term.set_raw
      fn (enable: bool) -> result[bool, string]
      + enters or leaves raw mode
      # terminal

liner
  liner.new
    fn () -> liner_state
    + creates an editor with empty buffer and empty history
    # construction
  liner.insert
    fn (state: liner_state, ch: string) -> liner_state
    + inserts a character at the cursor and advances it
    # editing
  liner.backspace
    fn (state: liner_state) -> liner_state
    + deletes the character before the cursor
    - leaves state unchanged when the cursor is at 0
    # editing
  liner.move_cursor
    fn (state: liner_state, delta: i32) -> liner_state
    + moves the cursor, clamped to [0, length]
    # editing
  liner.history_prev
    fn (state: liner_state) -> liner_state
    + replaces the buffer with the previous history entry
    - leaves state unchanged when there is no previous entry
    # history
  liner.history_next
    fn (state: liner_state) -> liner_state
    + replaces the buffer with the next history entry or clears it past the end
    # history
  liner.set_completer
    fn (state: liner_state, completer: fn(string) -> list[string]) -> liner_state
    + registers a completion function invoked on tab
    # completion
  liner.read_line
    fn (state: liner_state, prompt: string) -> result[tuple[string, liner_state], string]
    + writes the prompt, processes keys, and returns the finished line on enter
    + appends non-empty lines to history
    - returns error on terminal failure
    # io
    -> std.term.set_raw
    -> std.term.write
    -> std.term.read_key
