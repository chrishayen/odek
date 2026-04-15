# Requirement: "a library for drawing on the terminal using unicode braille characters"

Each braille glyph holds a 2-wide by 4-tall pixel grid, giving 8 sub-pixels per cell.

std: (all units exist)

braille_canvas
  braille_canvas.new
    fn (width_px: i32, height_px: i32) -> canvas_state
    + creates a blank canvas sized in pixels, rounded up to cell boundaries
    # construction
  braille_canvas.set
    fn (state: canvas_state, x: i32, y: i32) -> canvas_state
    + turns on the sub-pixel at (x, y)
    - leaves state unchanged when (x, y) is outside the canvas
    # drawing
  braille_canvas.unset
    fn (state: canvas_state, x: i32, y: i32) -> canvas_state
    + turns off the sub-pixel at (x, y)
    # drawing
  braille_canvas.line
    fn (state: canvas_state, x0: i32, y0: i32, x1: i32, y1: i32) -> canvas_state
    + rasterizes a straight line using Bresenham's algorithm
    # drawing
  braille_canvas.render
    fn (state: canvas_state) -> string
    + returns the canvas as a multi-line braille string
    + each 2x4 pixel block becomes one U+28xx code point
    # rendering
