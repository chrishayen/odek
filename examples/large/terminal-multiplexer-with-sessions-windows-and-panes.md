# Requirement: "a terminal multiplexer managing sessions, windows, and panes"

A session holds windows; a window holds a tree of panes. Each pane is backed by a pseudoterminal. The library manages layout, input routing, and output buffering.

std
  std.pty
    std.pty.spawn
      @ (command: string, args: list[string]) -> result[pty_handle, string]
      + starts a process attached to a new pseudoterminal
      - returns error when exec fails
      # process
    std.pty.write
      @ (handle: pty_handle, data: bytes) -> result[void, string]
      + writes to the pty input
      # io
    std.pty.read
      @ (handle: pty_handle, max: i32) -> result[bytes, string]
      + reads available output up to max bytes
      # io
    std.pty.resize
      @ (handle: pty_handle, cols: u16, rows: u16) -> result[void, string]
      + sends a window-size update to the pty
      # control
    std.pty.close
      @ (handle: pty_handle) -> void
      + closes the pty and terminates the child
      # lifecycle

multiplexer
  multiplexer.new_session
    @ (name: string) -> session_state
    + creates an empty session with one default window
    # session
  multiplexer.create_window
    @ (session: session_state, name: string) -> tuple[i32, session_state]
    + creates a window with a single pane and returns its id
    # window
  multiplexer.split_pane
    @ (session: session_state, pane_id: i32, direction: i8) -> result[tuple[i32, session_state], string]
    + splits pane_id horizontally (0) or vertically (1) and returns the new pane id
    - returns error when pane_id is unknown
    # pane
  multiplexer.close_pane
    @ (session: session_state, pane_id: i32) -> result[session_state, string]
    + removes a pane and collapses the layout
    - returns error when pane_id is unknown
    # pane
  multiplexer.compute_layout
    @ (session: session_state, window_id: i32, total_cols: u16, total_rows: u16) -> result[list[pane_rect], string]
    + divides the window area between panes based on the split tree
    - returns error when window_id is unknown
    # layout
  multiplexer.attach_process
    @ (session: session_state, pane_id: i32, command: string, args: list[string]) -> result[session_state, string]
    + starts a process under pane_id
    - returns error when pane_id already has a process
    # process_attach
    -> std.pty.spawn
  multiplexer.send_input
    @ (session: session_state, pane_id: i32, data: bytes) -> result[void, string]
    + writes data to the active pane's pty
    - returns error when pane has no attached process
    # input
    -> std.pty.write
  multiplexer.read_output
    @ (session: session_state, pane_id: i32) -> result[bytes, string]
    + reads buffered output from the pane
    # output
    -> std.pty.read
  multiplexer.resize_window
    @ (session: session_state, window_id: i32, cols: u16, rows: u16) -> result[session_state, string]
    + recomputes layout and sends resize to each pane's pty
    # resize
    -> std.pty.resize
  multiplexer.focus_pane
    @ (session: session_state, pane_id: i32) -> result[session_state, string]
    + marks a pane as the active input target
    - returns error when pane_id is unknown
    # focus
  multiplexer.save_layout
    @ (session: session_state) -> string
    + serializes session layout to a restorable string
    # persistence
  multiplexer.restore_layout
    @ (serialized: string) -> result[session_state, string]
    + reconstructs a session from a saved layout
    - returns error on malformed input
    # persistence
