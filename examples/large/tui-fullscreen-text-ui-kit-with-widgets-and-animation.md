# Requirement: "a full-screen text UI toolkit with widgets and animation"

Provides a screen buffer, input handling, and simple widgets. Rendering and input come through thin std terminal primitives.

std
  std.term
    std.term.size
      fn () -> tuple[i32, i32]
      + returns (rows, cols) of the current terminal
      # terminal
    std.term.write
      fn (data: string) -> void
      + writes raw bytes to the terminal output
      # terminal
    std.term.read_key
      fn () -> optional[key_event]
      + returns the next key event or none if none is pending
      # terminal
    std.term.enter_alt_screen
      fn () -> void
      + switches to the alternate screen buffer
      # terminal
    std.term.leave_alt_screen
      fn () -> void
      + restores the primary screen buffer
      # terminal
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (duration: i64) -> void
      + blocks for the given duration
      # time

tui
  tui.new_screen
    fn () -> screen_state
    + creates an empty screen buffer sized to the terminal
    # construction
    -> std.term.size
  tui.clear
    fn (screen: screen_state) -> screen_state
    + fills the buffer with spaces and the default attributes
    # drawing
  tui.put_text
    fn (screen: screen_state, row: i32, col: i32, text: string) -> screen_state
    + writes text starting at the given cell, clipping to the buffer
    # drawing
  tui.draw_box
    fn (screen: screen_state, row: i32, col: i32, height: i32, width: i32) -> screen_state
    + draws a single-line border rectangle
    # drawing
  tui.add_form_field
    fn (screen: screen_state, label: string, row: i32, col: i32, width: i32) -> tuple[screen_state, string]
    + places a labeled input field and returns its generated field id
    # widgets
  tui.set_field_value
    fn (screen: screen_state, field_id: string, value: string) -> result[screen_state, string]
    + updates the text held in the named field
    - returns error if the field id is unknown
    # widgets
  tui.add_animation
    fn (screen: screen_state, frames: list[string], row: i32, col: i32, period_ms: i64) -> tuple[screen_state, string]
    + registers a cycling frame animation and returns its id
    # animation
  tui.tick
    fn (screen: screen_state) -> screen_state
    + advances animations based on elapsed time and redraws affected regions
    # animation
    -> std.time.now_millis
  tui.handle_key
    fn (screen: screen_state, key: key_event) -> screen_state
    + routes a key event to the focused widget
    # input
  tui.flush
    fn (screen: screen_state) -> void
    + emits the buffer to the terminal
    # rendering
    -> std.term.write
  tui.run
    fn (screen: screen_state) -> screen_state
    + enters the alt screen, polls input and ticks animations until a quit key is received
    # event_loop
    -> std.term.enter_alt_screen
    -> std.term.leave_alt_screen
    -> std.term.read_key
    -> std.time.sleep_millis
