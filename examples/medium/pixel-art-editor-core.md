# Requirement: "a pixel art editor core"

A pixel canvas with modal editing: tools operate on an indexed grid with undoable commands. No UI; callers drive the canvas programmatically.

std
  std.collections
    std.collections.list_set
      fn (items: list[u8], index: i32, value: u8) -> list[u8]
      + returns a new list with index replaced
      # collections

pixel_editor
  pixel_editor.new_canvas
    fn (width: i32, height: i32) -> canvas_state
    + creates a canvas of the given size with all pixels set to index 0
    ? colors are stored as palette indices; a separate palette maps index -> rgba
    # construction
  pixel_editor.set_pixel
    fn (state: canvas_state, x: i32, y: i32, color_index: u8) -> canvas_state
    + sets the pixel at (x, y) to color_index, recording an undo entry
    - is a no-op when (x, y) is out of bounds
    # drawing
    -> std.collections.list_set
  pixel_editor.draw_line
    fn (state: canvas_state, x0: i32, y0: i32, x1: i32, y1: i32, color_index: u8) -> canvas_state
    + draws a line using Bresenham's algorithm as one undoable operation
    # drawing
  pixel_editor.flood_fill
    fn (state: canvas_state, x: i32, y: i32, color_index: u8) -> canvas_state
    + fills the contiguous region of matching color starting at (x, y)
    - is a no-op when target color equals source color
    # drawing
  pixel_editor.undo
    fn (state: canvas_state) -> canvas_state
    + reverts the most recent drawing operation
    - is a no-op when the undo stack is empty
    # history
  pixel_editor.redo
    fn (state: canvas_state) -> canvas_state
    + reapplies the most recently undone operation
    - is a no-op when the redo stack is empty
    # history
  pixel_editor.get_pixel
    fn (state: canvas_state, x: i32, y: i32) -> optional[u8]
    + returns the color index at (x, y)
    - returns none when (x, y) is out of bounds
    # querying
