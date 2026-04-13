# Requirement: "a grid-based 2D graphics library"

Maintain a fixed-size grid of colored cells and render it to a raster image. Encoding the image bytes is delegated to a thin std primitive.

std
  std.image
    std.image.encode_png
      @ (width: i32, height: i32, pixels: bytes) -> result[bytes, string]
      + encodes RGBA pixel bytes as a PNG
      - returns error when pixels length does not match width * height * 4
      # encoding

grid_draw
  grid_draw.new
    @ (cols: i32, rows: i32, cell_size_px: i32) -> grid_state
    + creates a grid with every cell transparent
    # construction
  grid_draw.set_cell
    @ (state: grid_state, col: i32, row: i32, rgba: u32) -> result[grid_state, string]
    + paints a single cell with the given RGBA color
    - returns error when coordinates are out of bounds
    # mutation
  grid_draw.draw_line
    @ (state: grid_state, from_col: i32, from_row: i32, to_col: i32, to_row: i32, rgba: u32) -> grid_state
    + paints every cell along the line between two grid coordinates
    ? uses Bresenham's line algorithm at cell resolution
    # drawing
  grid_draw.render_png
    @ (state: grid_state) -> result[bytes, string]
    + rasterizes the grid into a PNG image
    - returns error when encoding fails
    # rendering
    -> std.image.encode_png
