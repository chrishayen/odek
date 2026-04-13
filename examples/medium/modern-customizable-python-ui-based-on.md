# Requirement: "a themed widget toolkit over a base GUI library"

Provides themed widgets and a theme registry; rendering delegates to an underlying canvas abstraction.

std: (all units exist)

ui
  ui.new_theme
    @ (name: string) -> theme
    + creates an empty theme with the given identifier
    # construction
  ui.set_color
    @ (t: theme, token: string, rgba: u32) -> theme
    + sets a named color token on the theme
    # theming
  ui.set_font
    @ (t: theme, token: string, family: string, size: i32) -> theme
    + sets a named font token on the theme
    # theming
  ui.new_window
    @ (title: string, width: i32, height: i32, t: theme) -> window_state
    + creates a window with the given title, size, and theme
    # construction
  ui.add_button
    @ (w: window_state, id: string, label: string, x: i32, y: i32, on_click: fn() -> void) -> window_state
    + places a themed button at (x, y)
    # widgets
  ui.add_label
    @ (w: window_state, id: string, text: string, x: i32, y: i32) -> window_state
    + places a themed text label at (x, y)
    # widgets
  ui.add_entry
    @ (w: window_state, id: string, x: i32, y: i32, width: i32) -> window_state
    + places a themed single-line text entry at (x, y)
    # widgets
  ui.get_entry_text
    @ (w: window_state, id: string) -> result[string, string]
    + returns the current text of the named entry widget
    - returns error when no entry has that id
    # widgets
  ui.dispatch_click
    @ (w: window_state, x: i32, y: i32) -> window_state
    + finds the widget at (x, y) and invokes its click handler if it has one
    # events
