# Requirement: "an immediate-mode 2D GUI library"

Immediate-mode UI where widgets are declared every frame. The library tracks layout, input, and draw commands.

std
  std.math
    std.math.clamp_f32
      fn (value: f32, low: f32, high: f32) -> f32
      + clamps value into the inclusive range
      # math

gui
  gui.new_frame
    fn (viewport_w: f32, viewport_h: f32, input: input_state) -> frame_state
    + begins a new frame with viewport size and current input
    # frame_lifecycle
  gui.end_frame
    fn (frame: frame_state) -> draw_list
    + finalizes the frame and returns the accumulated draw commands
    # frame_lifecycle
  gui.begin_window
    fn (frame: frame_state, id: string, title: string, rect: rect_f32) -> frame_state
    + pushes a window container onto the layout stack
    # windowing
  gui.end_window
    fn (frame: frame_state) -> frame_state
    + pops the current window container
    - returns unchanged frame if no window is open
    # windowing
  gui.button
    fn (frame: frame_state, id: string, label: string) -> tuple[bool, frame_state]
    + returns true during the frame when the button was clicked
    + lays out the button within the current container
    # widgets
    -> std.math.clamp_f32
  gui.label
    fn (frame: frame_state, text: string) -> frame_state
    + adds a static text element to the current container
    # widgets
  gui.text_input
    fn (frame: frame_state, id: string, value: string) -> tuple[string, frame_state]
    + returns the edited value after applying keyboard input
    # widgets
  gui.slider_f32
    fn (frame: frame_state, id: string, value: f32, lo: f32, hi: f32) -> tuple[f32, frame_state]
    + returns the new slider value after applying drag input
    # widgets
    -> std.math.clamp_f32
  gui.checkbox
    fn (frame: frame_state, id: string, label: string, value: bool) -> tuple[bool, frame_state]
    + toggles value when clicked this frame
    # widgets
  gui.layout_row
    fn (frame: frame_state, widths: list[f32]) -> frame_state
    + configures the next widgets to lay out horizontally with the given widths
    # layout
  gui.hit_test
    fn (frame: frame_state, rect: rect_f32) -> hover_state
    + returns hover and click information for the rectangle given current input
    # input_routing
  gui.push_id
    fn (frame: frame_state, id: string) -> frame_state
    + scopes subsequent widget ids to avoid collisions
    # id_scope
