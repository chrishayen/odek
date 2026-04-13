# Requirement: "a cross-platform GUI toolkit"

A retained-mode widget tree with layout, input dispatch, and a pluggable native backend for drawing.

std
  std.math
    std.math.clamp_f32
      @ (value: f32, lo: f32, hi: f32) -> f32
      + returns value constrained to [lo, hi]
      # math
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

gui
  gui.new_app
    @ (title: string) -> app_state
    + creates an application with an empty window
    # construction
  gui.add_window
    @ (state: app_state, title: string, width: i32, height: i32) -> tuple[app_state, window_id]
    + attaches a window and returns its id
    # window
  gui.label
    @ (text: string) -> widget_node
    + returns a label widget with the given text
    # widget
  gui.button
    @ (text: string, on_click: event_handler) -> widget_node
    + returns a button widget that invokes on_click when pressed
    # widget
  gui.text_input
    @ (initial: string, on_change: change_handler) -> widget_node
    + returns an editable text field
    # widget
  gui.box
    @ (direction: layout_dir, children: list[widget_node]) -> widget_node
    + returns a container that lays children out horizontally or vertically
    # layout
  gui.set_root
    @ (state: app_state, window: window_id, root: widget_node) -> app_state
    + replaces the root widget for the given window
    # window
  gui.layout_pass
    @ (state: app_state, window: window_id, width: f32, height: f32) -> app_state
    + computes positions and sizes for every widget in the window
    # layout
    -> std.math.clamp_f32
  gui.dispatch_event
    @ (state: app_state, window: window_id, event: input_event) -> app_state
    + routes an input event to the widget under the pointer and fires handlers
    - no-op when no widget accepts the event
    # input
    -> std.time.now_millis
  gui.render
    @ (state: app_state, window: window_id, backend: draw_backend) -> draw_backend
    + walks the widget tree and issues draw commands to the backend
    # rendering
  gui.run_frame
    @ (state: app_state, backend: draw_backend) -> tuple[app_state, draw_backend]
    + advances one frame: layout, event drain, render
    # loop
  gui.request_redraw
    @ (state: app_state, window: window_id) -> app_state
    + marks the window as needing a redraw on the next frame
    # rendering
