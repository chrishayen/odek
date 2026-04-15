# Requirement: "a library for building terminal user interface applications with widgets, events, and rich colors"

A widget tree rendered to a cell grid, fed by a keyboard and resize event stream, with color and style support.

std
  std.term
    std.term.read_event
      fn () -> result[term_event, string]
      + returns the next key, resize, or mouse event
      - returns error when the terminal stream closes
      # terminal
    std.term.get_size
      fn () -> tuple[i32, i32]
      + returns the current width and height in cells
      # terminal
    std.term.write_frame
      fn (cells: list[cell]) -> result[void, string]
      + writes a full frame of cells to the terminal
      - returns error when stdout is closed
      # terminal
  std.text
    std.text.grapheme_width
      fn (g: string) -> i32
      + returns the display width of a grapheme cluster (0, 1, or 2)
      # text

tui
  tui.new_app
    fn () -> app_state
    + creates an empty application with no root widget
    # construction
  tui.set_root
    fn (state: app_state, root: widget) -> app_state
    + installs the root widget
    # construction
  tui.layout
    fn (state: app_state, width: i32, height: i32) -> layout_tree
    + computes rectangles for every widget in the tree
    # layout
  tui.render
    fn (tree: layout_tree) -> list[cell]
    + emits a full frame of styled cells
    # rendering
    -> std.text.grapheme_width
  tui.dispatch_key
    fn (state: app_state, key: key_event) -> app_state
    + routes a key event to the focused widget
    # events
  tui.focus_next
    fn (state: app_state) -> app_state
    + advances focus to the next focusable widget
    # focus
  tui.make_text
    fn (content: string, style: style_spec) -> widget
    + creates a text widget with the given style
    # widgets
  tui.make_box
    fn (child: widget, border: border_spec) -> widget
    + wraps a child in a bordered box
    # widgets
  tui.make_list
    fn (items: list[string], style: style_spec) -> widget
    + creates a scrollable list widget
    # widgets
  tui.make_input
    fn (placeholder: string) -> widget
    + creates a single-line text input widget
    # widgets
  tui.run
    fn (state: app_state) -> result[void, string]
    + runs the event loop: read event, dispatch, layout, render
    - returns error when terminal I/O fails
    # event_loop
    -> std.term.read_event
    -> std.term.get_size
    -> std.term.write_frame
