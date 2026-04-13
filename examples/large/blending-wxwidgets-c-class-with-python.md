# Requirement: "a cross-platform native-widget GUI toolkit"

A retained-mode widget tree with layout, event dispatch, and a backend-agnostic draw interface.

std
  std.math
    std.math.max_i32
      @ (a: i32, b: i32) -> i32
      + returns the larger of two signed integers
      # arithmetic
    std.math.min_i32
      @ (a: i32, b: i32) -> i32
      + returns the smaller of two signed integers
      # arithmetic

gui
  gui.new_window
    @ (title: string, width: i32, height: i32) -> widget_state
    + creates a top-level window widget with the given size
    # construction
  gui.add_child
    @ (parent: widget_state, child: widget_state) -> widget_state
    + appends child to parent's child list
    # tree
  gui.label
    @ (text: string) -> widget_state
    + creates a text label widget
    # widgets
  gui.button
    @ (text: string, on_click: fn() -> void) -> widget_state
    + creates a button widget with a click handler
    # widgets
  gui.text_input
    @ (placeholder: string) -> widget_state
    + creates a single-line text input widget
    # widgets
  gui.vbox
    @ (spacing: i32) -> widget_state
    + creates a vertical box container with the given pixel spacing between children
    # layout
  gui.hbox
    @ (spacing: i32) -> widget_state
    + creates a horizontal box container
    # layout
  gui.measure
    @ (widget: widget_state) -> size
    + recursively computes each widget's preferred size based on content and layout rules
    # layout
    -> std.math.max_i32
  gui.arrange
    @ (widget: widget_state, bounds: rect) -> widget_state
    + assigns concrete positions to every widget in the tree given the available bounds
    # layout
    -> std.math.min_i32
  gui.dispatch_event
    @ (root: widget_state, event: input_event) -> widget_state
    + routes mouse or key events to the topmost widget under the cursor
    + propagates bubbling events from target to root until one handles it
    - ignores events whose coordinates fall outside the root bounds
    # events
  gui.draw
    @ (root: widget_state, canvas: canvas_sink) -> void
    + walks the arranged tree and issues draw commands to the backend canvas
    ? the canvas_sink abstracts the platform drawing API
    # rendering
  gui.set_text
    @ (widget: widget_state, text: string) -> widget_state
    + updates the text of a label or input and marks layout dirty
    # widgets
