# Requirement: "a GUI toolkit binding"

Widgets are an opaque tree; the project layer exposes construction, layout, and an event loop. The real drawing primitives live in std.

std
  std.display
    std.display.open_window
      fn (title: string, width: i32, height: i32) -> result[window_handle, string]
      + creates a top-level window
      - returns error when no display is available
      # display
    std.display.close_window
      fn (handle: window_handle) -> void
      + destroys a window and its children
      # display
    std.display.run_loop
      fn (handle: window_handle) -> void
      + enters the platform event loop until the window closes
      # display
    std.display.draw_text
      fn (handle: window_handle, x: i32, y: i32, text: string) -> void
      + renders text at the given pixel coordinates
      # display
    std.display.draw_rect
      fn (handle: window_handle, x: i32, y: i32, w: i32, h: i32) -> void
      + renders a rectangle outline
      # display
  std.input
    std.input.poll_event
      fn (handle: window_handle) -> optional[input_event]
      + returns the next keyboard or mouse event, or none if the queue is empty
      # input

gui
  gui.new_root
    fn (title: string, width: i32, height: i32) -> result[widget_state, string]
    + creates a root window widget with a blank widget tree
    - returns error when the window cannot be opened
    # construction
    -> std.display.open_window
  gui.add_label
    fn (parent: widget_state, text: string) -> widget_state
    + appends a label widget to the parent
    # widgets
  gui.add_button
    fn (parent: widget_state, text: string, on_click: string) -> widget_state
    + appends a button that dispatches the named handler when clicked
    # widgets
  gui.add_entry
    fn (parent: widget_state, name: string) -> widget_state
    + appends a single-line text entry field
    # widgets
  gui.pack
    fn (widget: widget_state, side: string) -> widget_state
    + lays out a widget on the given side ("top", "left", "right", "bottom")
    - returns unchanged when side is unknown
    # layout
  gui.grid
    fn (widget: widget_state, row: i32, column: i32) -> widget_state
    + places a widget at the given row and column
    # layout
  gui.get_entry_text
    fn (widget: widget_state, name: string) -> optional[string]
    + returns the current text of the named entry
    - returns none when no entry by that name exists
    # access
  gui.bind_handler
    fn (widget: widget_state, name: string, handler_id: string) -> widget_state
    + associates a handler id with an event name
    # events
  gui.dispatch
    fn (widget: widget_state, event: input_event) -> widget_state
    + routes an event to the matching bound handler
    # events
  gui.render
    fn (widget: widget_state) -> widget_state
    + draws every visible widget to its root window
    # rendering
    -> std.display.draw_text
    -> std.display.draw_rect
  gui.run
    fn (widget: widget_state) -> void
    + renders and processes events until the window is closed
    # lifecycle
    -> std.display.run_loop
    -> std.input.poll_event
  gui.destroy
    fn (widget: widget_state) -> void
    + tears down the widget tree and closes the window
    # lifecycle
    -> std.display.close_window
