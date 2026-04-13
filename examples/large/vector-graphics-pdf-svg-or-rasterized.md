# Requirement: "a vector graphics library that renders to PDF, SVG, or a raster image"

A small immediate-mode canvas plus three independent backends that consume the same drawing operations.

std
  std.math
    std.math.sin
      @ (radians: f64) -> f64
      + returns the sine of the given angle in radians
      # math
    std.math.cos
      @ (radians: f64) -> f64
      + returns the cosine of the given angle in radians
      # math
  std.io
    std.io.write_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or overwriting it
      - returns error on IO failure
      # io

vector_canvas
  vector_canvas.new
    @ (width: f64, height: f64) -> canvas_state
    + returns an empty canvas of the given dimensions with an empty operation list
    # construction
  vector_canvas.move_to
    @ (canvas: canvas_state, x: f64, y: f64) -> canvas_state
    + appends a move-to operation starting a new subpath at (x, y)
    # path_building
  vector_canvas.line_to
    @ (canvas: canvas_state, x: f64, y: f64) -> canvas_state
    + appends a line-to operation extending the current subpath
    # path_building
  vector_canvas.arc
    @ (canvas: canvas_state, cx: f64, cy: f64, radius: f64, start_rad: f64, end_rad: f64) -> canvas_state
    + appends an arc operation approximated by the backends
    # path_building
    -> std.math.sin
    -> std.math.cos
  vector_canvas.close_path
    @ (canvas: canvas_state) -> canvas_state
    + appends a close-path operation
    # path_building
  vector_canvas.set_stroke
    @ (canvas: canvas_state, r: f64, g: f64, b: f64, width: f64) -> canvas_state
    + sets the stroke color and width for subsequent path operations
    # styling
  vector_canvas.set_fill
    @ (canvas: canvas_state, r: f64, g: f64, b: f64) -> canvas_state
    + sets the fill color for subsequent path operations
    # styling
  vector_canvas.render_svg
    @ (canvas: canvas_state) -> string
    + returns a self-contained SVG document representing all operations
    # svg_backend
  vector_canvas.render_pdf
    @ (canvas: canvas_state) -> bytes
    + returns a single-page PDF document representing all operations
    # pdf_backend
  vector_canvas.render_raster
    @ (canvas: canvas_state, pixel_width: i32, pixel_height: i32) -> bytes
    + returns a PNG image produced by rasterizing the operation list at the given resolution
    # raster_backend
  vector_canvas.save
    @ (canvas: canvas_state, path: string, format: string) -> result[void, string]
    + writes the canvas to disk as "svg", "pdf", or "png" based on the format argument
    - returns error on unknown format
    - returns error on IO failure
    # io
    -> std.io.write_bytes
