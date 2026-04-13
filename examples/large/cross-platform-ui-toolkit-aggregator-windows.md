# Requirement: "a cross-platform UI toolkit abstraction over multiple native backends"

A widget model compiled to a draw list and event stream that any native backend can implement. The library exposes the model; backends implement a small driver interface.

std
  std.math
    std.math.clamp_i32
      @ (v: i32, lo: i32, hi: i32) -> i32
      + returns v constrained to [lo, hi]
      # math
  std.event
    std.event.new_queue
      @ () -> event_queue
      + creates an empty event queue
      # events
    std.event.push
      @ (queue: event_queue, event: ui_event) -> event_queue
      + enqueues an event
      # events
    std.event.pop
      @ (queue: event_queue) -> tuple[optional[ui_event], event_queue]
      + removes and returns the next event
      # events

ui_toolkit
  ui_toolkit.new_window
    @ (title: string, width: i32, height: i32) -> window_state
    + creates a window with initial title and size
    # construction
  ui_toolkit.set_root
    @ (window: window_state, root: widget) -> window_state
    + installs the root widget
    # construction
  ui_toolkit.make_button
    @ (label: string) -> widget
    + creates a clickable button widget
    # widgets
  ui_toolkit.make_label
    @ (text: string) -> widget
    + creates a static text widget
    # widgets
  ui_toolkit.make_text_field
    @ (placeholder: string) -> widget
    + creates a single-line editable text widget
    # widgets
  ui_toolkit.make_stack
    @ (direction: stack_direction, children: list[widget]) -> widget
    + creates a horizontal or vertical stack container
    # widgets
  ui_toolkit.layout
    @ (window: window_state) -> layout_tree
    + computes rectangles for every widget
    # layout
    -> std.math.clamp_i32
  ui_toolkit.render
    @ (tree: layout_tree) -> list[draw_command]
    + emits a backend-agnostic draw list
    # rendering
  ui_toolkit.register_backend
    @ (driver: backend_driver) -> backend_registry
    + registers a native backend driver that can present a window and deliver events
    # backend
  ui_toolkit.present
    @ (registry: backend_registry, window: window_state, commands: list[draw_command]) -> result[void, string]
    + delegates presentation of the draw list to the active backend
    - returns error when no backend is registered
    # backend
  ui_toolkit.dispatch_event
    @ (window: window_state, event: ui_event) -> window_state
    + routes an event to the target widget through hit-testing
    # events
    -> std.event.pop
  ui_toolkit.run
    @ (registry: backend_registry, window: window_state) -> result[void, string]
    + runs the event loop: pull events, dispatch, layout, render, present
    - returns error when the backend or any widget handler fails fatally
    # event_loop
    -> std.event.new_queue
    -> std.event.push
