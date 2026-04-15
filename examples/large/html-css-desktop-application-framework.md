# Requirement: "a desktop application framework driven by HTML and CSS"

A library that lets a caller mount HTML/CSS views, wire event handlers, and run a native window event loop. The rendering backend is pluggable.

std
  std.html
    std.html.parse
      fn (source: string) -> result[dom_node, string]
      + parses a well-formed HTML fragment into a DOM tree
      - returns error on unbalanced tags
      # parsing
    std.html.serialize
      fn (root: dom_node) -> string
      + serializes a DOM tree back to HTML
      # serialization
  std.css
    std.css.parse_stylesheet
      fn (source: string) -> result[stylesheet, string]
      + parses CSS rules into a stylesheet structure
      - returns error on unterminated blocks
      # parsing
    std.css.match_rules
      fn (sheet: stylesheet, node: dom_node) -> list[css_rule]
      + returns rules whose selector matches the given node
      # styling
  std.event
    std.event.new_queue
      fn () -> event_queue
      + returns an empty event queue
      # events
    std.event.push
      fn (queue: event_queue, evt: event) -> void
      + appends an event to the queue
      # events
    std.event.pop
      fn (queue: event_queue) -> optional[event]
      + returns the next event or none when empty
      # events

apps
  apps.new_window
    fn (title: string, width: i32, height: i32) -> window_handle
    + creates a window with the given title and dimensions
    # window_management
  apps.mount
    fn (win: window_handle, html: string, css: string) -> result[view_state, string]
    + parses html and css and attaches them as the window's root view
    - returns error when html or css is malformed
    # mounting
    -> std.html.parse
    -> std.css.parse_stylesheet
  apps.on
    fn (view: view_state, selector: string, event_name: string, handler: string) -> result[view_state, string]
    + registers a named handler for an event on nodes matching selector
    - returns error when the selector is invalid
    # event_binding
  apps.dispatch
    fn (view: view_state, event_name: string, target_id: string) -> view_state
    + fires the named event at the target node and runs matching handlers
    # event_dispatch
    -> std.event.push
  apps.set_text
    fn (view: view_state, node_id: string, text: string) -> result[view_state, string]
    + replaces the text content of the identified node
    - returns error when the id is not present
    # dom_mutation
  apps.render
    fn (view: view_state) -> string
    + returns the current rendered HTML for the view
    # rendering
    -> std.html.serialize
    -> std.css.match_rules
  apps.run_loop
    fn (win: window_handle, view: view_state) -> view_state
    + drains the window event queue, dispatching each event, and returns when the window closes
    # event_loop
    -> std.event.pop
  apps.close
    fn (win: window_handle) -> void
    + releases the window and its backing resources
    # lifecycle
