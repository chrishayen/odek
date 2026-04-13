# Requirement: "a desktop UI library that renders an HTML-based view and dispatches DOM events to handlers"

The library owns a window backed by an embedded web view. Callers register event handlers and push state updates; the library handles the bridge.

std
  std.webview
    std.webview.create_window
      @ (title: string, width: i32, height: i32) -> result[window_handle, string]
      + creates a native window hosting a web view
      - returns error when the platform cannot create a window
      # window
    std.webview.load_html
      @ (window: window_handle, html: string) -> result[void, string]
      + loads an HTML document into the window
      # window
    std.webview.eval_js
      @ (window: window_handle, script: string) -> result[string, string]
      + evaluates script in the page and returns the stringified result
      - returns error on syntax error
      # bridge
    std.webview.on_message
      @ (window: window_handle, handler: fn(string) -> void) -> void
      + registers a handler for messages posted from the page via the host bridge
      # bridge

desktop_ui
  desktop_ui.open
    @ (title: string, width: i32, height: i32, initial_html: string) -> result[ui_state, string]
    + opens a window with the given HTML and wires the message bridge
    - returns error when window creation fails
    # lifecycle
    -> std.webview.create_window
    -> std.webview.load_html
    -> std.webview.on_message
  desktop_ui.on_event
    @ (state: ui_state, event_name: string, handler: fn(string) -> void) -> ui_state
    + registers a handler for a named DOM event bridged from the page
    # events
  desktop_ui.set_element_text
    @ (state: ui_state, selector: string, text: string) -> result[void, string]
    + updates the textContent of elements matching the selector
    - returns error when the script fails
    # updates
    -> std.webview.eval_js
  desktop_ui.run_script
    @ (state: ui_state, script: string) -> result[string, string]
    + evaluates arbitrary script in the page and returns the result
    # updates
    -> std.webview.eval_js
