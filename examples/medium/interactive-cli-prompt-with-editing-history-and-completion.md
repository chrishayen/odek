# Requirement: "a library for building interactive command-line prompts with editing, history, and completion"

Provides a reader that consumes raw keystrokes and produces edited lines with history and completion support.

std
  std.term
    std.term.read_key
      @ () -> result[key_event, string]
      + reads one keystroke in raw mode
      - returns error when the terminal cannot be placed in raw mode
      # terminal
    std.term.write
      @ (s: string) -> void
      + writes a string to the terminal with ANSI escapes
      # terminal

interline
  interline.new
    @ (prompt: string) -> line_state
    + creates a line buffer showing the given prompt
    # construction
  interline.insert_char
    @ (state: line_state, ch: string) -> line_state
    + inserts a character at the cursor
    # editing
  interline.delete_back
    @ (state: line_state) -> line_state
    + deletes the character before the cursor
    # editing
  interline.move_cursor
    @ (state: line_state, delta: i32) -> line_state
    + moves the cursor left or right, clamped to the buffer
    # editing
  interline.history_new
    @ (capacity: i32) -> history_state
    + creates a bounded ring-buffer history
    # history
  interline.history_push
    @ (hist: history_state, entry: string) -> history_state
    + appends an entry, evicting the oldest when full
    # history
  interline.history_recall
    @ (hist: history_state, offset: i32) -> optional[string]
    + returns the entry at offset from the most recent
    # history
  interline.set_completer
    @ (state: line_state, completer: callback) -> line_state
    + installs a function that returns candidate completions for the current prefix
    # completion
  interline.complete
    @ (state: line_state) -> line_state
    + replaces the current word with the longest common prefix of candidates
    # completion
  interline.read_line
    @ (state: line_state, hist: history_state) -> result[tuple[string, history_state], string]
    + reads keys in a loop until Enter and returns the final line with updated history
    - returns error on EOF
    # reading
    -> std.term.read_key
    -> std.term.write
