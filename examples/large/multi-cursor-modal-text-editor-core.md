# Requirement: "a multi-cursor modal text editor core"

A buffer with a set of selections, modal commands that transform both buffer and selections, and an undo stack. Rendering and I/O are the caller's problem.

std: (all units exist)

editor
  editor.buffer_new
    fn (initial: string) -> buffer_state
    + creates a buffer holding the given text
    # construction
  editor.buffer_text
    fn (buffer: buffer_state) -> string
    + returns the full buffer contents
    # inspection
  editor.selection_set
    fn (buffer: buffer_state, ranges: list[tuple[i32, i32]]) -> buffer_state
    + replaces the current selection set with the given ranges
    ? ranges are (start, end) byte offsets
    # selections
  editor.selections
    fn (buffer: buffer_state) -> list[tuple[i32, i32]]
    + returns the current selection ranges
    # selections
  editor.insert_at_selections
    fn (buffer: buffer_state, text: string) -> buffer_state
    + inserts text at each selection and shifts subsequent ranges
    # editing
  editor.delete_selections
    fn (buffer: buffer_state) -> buffer_state
    + removes text inside each selection and collapses each range
    # editing
  editor.move_selections
    fn (buffer: buffer_state, delta: i32) -> buffer_state
    + shifts each selection by delta bytes, clamped to buffer bounds
    # motion
  editor.extend_to_word
    fn (buffer: buffer_state) -> buffer_state
    + expands each collapsed selection to cover the word under the cursor
    - leaves selection empty when the cursor is on whitespace
    # motion
  editor.add_cursor_below
    fn (buffer: buffer_state) -> buffer_state
    + adds a new selection on the next line at the same column
    - no-op when already on the last line
    # multi_cursor
  editor.apply_mode_command
    fn (buffer: buffer_state, mode: i32, command: string) -> result[buffer_state, string]
    + dispatches a command under the given mode (0=normal, 1=select, 2=insert)
    - returns error on unknown command
    # modal_dispatch
  editor.undo
    fn (buffer: buffer_state) -> buffer_state
    + pops the last edit and restores prior text and selections
    - no-op when the undo stack is empty
    # history
  editor.redo
    fn (buffer: buffer_state) -> buffer_state
    + reapplies the most recently undone edit
    - no-op when there is nothing to redo
    # history
