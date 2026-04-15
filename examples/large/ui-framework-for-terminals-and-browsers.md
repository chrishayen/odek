# Requirement: "a framework for building interactive user interfaces that run in terminals and browsers"

A retained-mode UI framework. Widgets declare a tree, an event loop dispatches input, and a renderer walks the tree producing abstract draw commands that a backend (terminal, browser canvas, test harness) turns into pixels or cells.

std
  std.collections
    std.collections.list_index_of
      fn (items: list[string], target: string) -> optional[i32]
      + returns the first index where items[i] equals target
      - returns none when target is absent
      # collections
  std.strings
    std.strings.measure_width
      fn (s: string) -> i32
      + returns the display width in cells (1 per ASCII, 2 per wide char)
      # text_metrics

ui
  ui.new_app
    fn (root_widget_id: string) -> app_state
    + creates an app with the given root widget and an empty event queue
    # construction
  ui.mount
    fn (state: app_state, parent_id: string, child_id: string, props: map[string,string]) -> app_state
    + attaches a child widget to a parent and stores its props
    - returns unchanged state when parent_id is not in the tree
    # tree_construction
  ui.unmount
    fn (state: app_state, widget_id: string) -> app_state
    + removes a widget and all its descendants
    # tree_construction
  ui.set_prop
    fn (state: app_state, widget_id: string, key: string, value: string) -> app_state
    + updates a single prop and marks the widget dirty
    # state_update
  ui.push_event
    fn (state: app_state, event_kind: string, target_id: string, payload: string) -> app_state
    + appends an event to the queue for the next tick
    # event_queue
  ui.tick
    fn (state: app_state) -> app_state
    + drains the event queue, invoking each widget's handler in insertion order
    + bubbles unhandled events from child to parent
    # event_loop
    -> std.collections.list_index_of
  ui.layout
    fn (state: app_state, viewport_cols: i32, viewport_rows: i32) -> layout_tree
    + computes x/y/width/height for every mounted widget given the viewport
    ? uses a simple box model: each widget fills its parent unless a width/height prop overrides
    # layout
    -> std.strings.measure_width
  ui.render
    fn (state: app_state, layout: layout_tree) -> list[draw_command]
    + walks the layout tree producing draw commands (move_to, write_text, set_fg, set_bg, clear)
    + emits nothing for clean (non-dirty) subtrees when possible
    # rendering
  ui.clear_dirty
    fn (state: app_state) -> app_state
    + marks every widget as clean after a successful render
    # state_update
  ui.focus
    fn (state: app_state, widget_id: string) -> app_state
    + sets the focused widget so keyboard events target it first
    - returns unchanged state when widget_id is not mounted
    # focus_management
  ui.inject_key
    fn (state: app_state, key: string, modifiers: i32) -> app_state
    + translates a key press into a "key" event targeted at the focused widget
    # input
    -> std.collections.list_index_of
