# Requirement: "a cross-platform tui widget library with interactive widgets"

A small widget tree with focus, events, and a renderer that writes to a cell buffer. Terminal IO is a thin std primitive so tests can substitute a fake backend.

std
  std.term
    std.term.get_size
      @ () -> tuple[u16, u16]
      + returns (columns, rows) of the current terminal
      # terminal
    std.term.write
      @ (data: string) -> void
      + writes data to the terminal output stream
      # terminal
    std.term.read_key
      @ () -> result[key_event, string]
      + blocks until a key event is available and returns it
      - returns error when stdin is closed
      # terminal

tui
  tui.new_screen
    @ (width: u16, height: u16) -> screen_state
    + creates an empty cell buffer of the given dimensions
    # screen
  tui.make_button
    @ (label: string) -> widget
    + creates a button widget that emits a click event when activated
    # widget
  tui.make_text_input
    @ (placeholder: string) -> widget
    + creates a single-line text input widget
    # widget
  tui.make_vbox
    @ (children: list[widget]) -> widget
    + creates a vertical container that stacks children top-to-bottom
    # layout
  tui.focus_next
    @ (screen: screen_state) -> screen_state
    + advances focus to the next focusable widget in tree order
    + wraps to the first focusable widget when at the end
    # focus
  tui.dispatch_key
    @ (screen: screen_state, event: key_event) -> screen_state
    + routes the key event to the focused widget and returns the updated screen
    - ignores events when no widget is focusable
    # events
  tui.render
    @ (screen: screen_state) -> void
    + draws the widget tree into the cell buffer and flushes to the terminal
    # rendering
    -> std.term.write
