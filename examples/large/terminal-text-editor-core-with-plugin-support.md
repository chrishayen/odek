# Requirement: "a terminal text editor core with plugin support"

Editor core: buffer, cursor, edits, undo, file I/O, and a plugin hook registry. Terminal rendering is the caller's concern.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the entire file as text
      - returns error when the file does not exist
      # fs
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to the path replacing any existing file
      - returns error on I/O failure
      # fs
  std.strings
    std.strings.split_lines
      fn (input: string) -> list[string]
      + splits input into lines preserving order
      + returns empty list for empty string
      # strings
    std.strings.join_lines
      fn (lines: list[string]) -> string
      + joins lines with newline separators
      # strings

editor
  editor.new
    fn () -> editor_state
    + creates an empty buffer with cursor at origin
    # construction
  editor.open_file
    fn (state: editor_state, path: string) -> result[editor_state, string]
    + loads file contents into the buffer and records the path
    - returns error when the file cannot be read
    # io
    -> std.fs.read_all
    -> std.strings.split_lines
  editor.save_file
    fn (state: editor_state) -> result[void, string]
    + writes the current buffer back to the recorded path
    - returns error when no path is associated with the buffer
    # io
    -> std.strings.join_lines
    -> std.fs.write_all
  editor.insert_text
    fn (state: editor_state, text: string) -> editor_state
    + inserts text at the cursor and advances the cursor
    # editing
  editor.delete_range
    fn (state: editor_state, start: cursor_pos, end: cursor_pos) -> editor_state
    + removes the text between start and end
    - leaves buffer unchanged when start >= end
    # editing
  editor.move_cursor
    fn (state: editor_state, pos: cursor_pos) -> editor_state
    + relocates the cursor, clamping to buffer bounds
    # navigation
  editor.undo
    fn (state: editor_state) -> editor_state
    + reverts the most recent edit
    - returns unchanged state when history is empty
    # history
  editor.redo
    fn (state: editor_state) -> editor_state
    + reapplies the most recently undone edit
    - returns unchanged state when redo stack is empty
    # history
  editor.get_line
    fn (state: editor_state, index: i32) -> result[string, string]
    + returns the line at the given zero-based index
    - returns error when index is out of range
    # inspection
  editor.line_count
    fn (state: editor_state) -> i32
    + returns the number of lines in the buffer
    # inspection
  editor.register_hook
    fn (state: editor_state, event: string, handler: fn(ctx: hook_context) -> editor_state) -> result[editor_state, string]
    + attaches a plugin callback to an editor event
    - returns error when event name is not recognized
    # plugins
  editor.fire_hook
    fn (state: editor_state, event: string, ctx: hook_context) -> editor_state
    + invokes every registered handler for the event in registration order
    # plugins
