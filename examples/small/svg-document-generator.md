# Requirement: "an SVG document generation library"

Build an SVG by appending shape primitives, then serialize to a string.

std: (all units exist)

svg
  svg.new_canvas
    fn (width: i32, height: i32) -> canvas_state
    + creates a canvas with the given pixel dimensions
    # construction
  svg.add_rect
    fn (state: canvas_state, x: i32, y: i32, w: i32, h: i32, fill: string) -> canvas_state
    + appends a rectangle with the given position, size, and fill color
    # shapes
  svg.add_circle
    fn (state: canvas_state, cx: i32, cy: i32, r: i32, fill: string) -> canvas_state
    + appends a circle
    # shapes
  svg.add_line
    fn (state: canvas_state, x1: i32, y1: i32, x2: i32, y2: i32, stroke: string) -> canvas_state
    + appends a line segment
    # shapes
  svg.add_text
    fn (state: canvas_state, x: i32, y: i32, content: string, fill: string) -> canvas_state
    + appends a text element at the given anchor
    # shapes
  svg.render
    fn (state: canvas_state) -> string
    + serializes the canvas to an SVG document string with proper header
    # serialization
