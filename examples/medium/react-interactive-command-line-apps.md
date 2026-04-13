# Requirement: "a component-based interactive terminal UI library"

Virtual-tree style terminal UI: components produce a node tree, the library diffs it against the previous frame and writes the minimum updates to the terminal.

std
  std.term
    std.term.size
      @ () -> tuple[i32, i32]
      + returns the current terminal columns and rows
      # terminal
    std.term.write
      @ (data: string) -> result[void, string]
      + writes data to the terminal output stream
      # terminal
    std.term.read_key
      @ () -> result[key_event, string]
      + blocks until a key event is available
      - returns error when stdin is closed
      # input

tui
  tui.text
    @ (content: string) -> node
    + returns a leaf node rendering content in one row
    # component
  tui.box
    @ (children: list[node], direction: string) -> node
    + returns a container node laying out children row-wise or column-wise
    ? direction is "row" or "column"
    # component
  tui.render_frame
    @ (tree: node, previous: optional[frame], cols: i32, rows: i32) -> frame
    + returns a frame describing the cells to display for this tree within cols x rows
    # rendering
  tui.diff_and_flush
    @ (current: frame, previous: optional[frame]) -> result[void, string]
    + writes only the cells that changed since previous to the terminal
    # rendering
    -> std.term.write
  tui.run
    @ (initial_state: app_state, view: fn(app_state) -> node, update: fn(app_state, key_event) -> app_state) -> result[void, string]
    + runs an input/update/render loop until update returns a state marked as exited
    # event_loop
    -> std.term.size
    -> std.term.read_key
