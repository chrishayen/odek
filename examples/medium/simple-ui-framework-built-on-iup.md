# Requirement: "a retained-mode ui framework over a native widget toolkit"

A thin declarative layer over a native widget backend: build a tree of widgets, show a window, and dispatch events to handlers.

std
  std.widget
    std.widget.init
      @ () -> result[void, string]
      + initializes the native widget toolkit
      - returns error when the native backend cannot be started
      # backend
    std.widget.create_window
      @ (title: string, width: i32, height: i32) -> result[native_handle, string]
      + creates a top-level window handle
      # backend
    std.widget.create_label
      @ (text: string) -> native_handle
      + creates a label widget
      # backend
    std.widget.create_button
      @ (text: string) -> native_handle
      + creates a button widget
      # backend
    std.widget.create_vbox
      @ (children: list[native_handle]) -> native_handle
      + creates a vertical container around the children
      # backend
    std.widget.set_child
      @ (parent: native_handle, child: native_handle) -> void
      + assigns the child as the content of the parent container
      # backend
    std.widget.on_click
      @ (handle: native_handle, handler_id: string) -> void
      + registers a click handler identifier with the widget
      # backend
    std.widget.show
      @ (window: native_handle) -> void
      + makes the window visible
      # backend
    std.widget.run_loop
      @ () -> void
      + enters the native event loop until the last window is closed
      # backend
    std.widget.next_event
      @ () -> optional[event]
      + returns the next pending event or none
      # backend

ui
  ui.new_app
    @ (title: string, width: i32, height: i32) -> result[app_state, string]
    + initializes the toolkit and creates the root window
    - returns error when the backend cannot start
    # construction
    -> std.widget.init
    -> std.widget.create_window
  ui.label
    @ (text: string) -> ui_node
    + returns a label node with the given text
    # widgets
  ui.button
    @ (text: string, on_click: string) -> ui_node
    + returns a button node bound to the click handler id
    # widgets
  ui.column
    @ (children: list[ui_node]) -> ui_node
    + returns a vertical container node
    # layout
  ui.set_content
    @ (state: app_state, root: ui_node) -> app_state
    + materializes the node tree into native widgets and installs it in the window
    # mounting
    -> std.widget.create_label
    -> std.widget.create_button
    -> std.widget.create_vbox
    -> std.widget.set_child
    -> std.widget.on_click
  ui.run
    @ (state: app_state, dispatch: map[string, string]) -> void
    + shows the window and runs the event loop, calling dispatch entries by handler id
    # runtime
    -> std.widget.show
    -> std.widget.run_loop
    -> std.widget.next_event
