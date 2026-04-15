# Requirement: "a gpu-accelerated immediate-mode gui framework"

An immediate-mode widget API backed by a thin gpu draw-list primitive. The project layer is widgets and layout; gpu submission is a std seam.

std
  std.gpu
    std.gpu.begin_frame
      fn (width: i32, height: i32) -> draw_list
      + starts a new draw list sized to the target surface
      # gpu
    std.gpu.submit
      fn (dl: draw_list) -> result[void, string]
      + submits the draw list to the gpu and presents
      - returns error when the device context is lost
      # gpu
  std.text
    std.text.measure
      fn (text: string, size_px: i32) -> tuple[i32, i32]
      + returns the pixel width and height of a text run
      # text

gui
  gui.new_context
    fn () -> gui_context
    + creates a gui context with default theme and empty state
    # construction
  gui.begin
    fn (ctx: gui_context, width: i32, height: i32) -> gui_context
    + opens a new frame, resetting per-frame widget state
    # frame
    -> std.gpu.begin_frame
  gui.end
    fn (ctx: gui_context) -> result[gui_context, string]
    + closes the frame and submits the draw list
    - returns error when submission fails
    # frame
    -> std.gpu.submit
  gui.button
    fn (ctx: gui_context, id: string, label: string) -> tuple[bool, gui_context]
    + returns true on the frame the button was clicked
    - returns false while the button is merely hovered or idle
    # widget
    -> std.text.measure
  gui.slider_f32
    fn (ctx: gui_context, id: string, value: f32, min: f32, max: f32) -> tuple[f32, gui_context]
    + returns the updated value after user drag input
    + clamps the returned value to [min, max]
    # widget
  gui.text
    fn (ctx: gui_context, content: string) -> gui_context
    + queues a text run for drawing at the current cursor
    # widget
    -> std.text.measure
  gui.begin_window
    fn (ctx: gui_context, id: string, title: string) -> gui_context
    + opens a window; subsequent widgets are laid out inside it
    # layout
  gui.end_window
    fn (ctx: gui_context) -> gui_context
    + closes the most recently opened window
    # layout
