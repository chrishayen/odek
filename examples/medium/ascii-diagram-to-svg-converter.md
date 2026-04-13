# Requirement: "an ASCII diagram to SVG converter"

Parse ASCII art into a grid of glyphs, classify shapes (lines, boxes, arrows), and emit SVG elements.

std
  std.text
    std.text.split_lines
      @ (input: string) -> list[string]
      + splits a string into lines on LF and CRLF
      + returns a single-element list when no newline is present
      # text
  std.xml
    std.xml.escape_attr
      @ (value: string) -> string
      + escapes &, <, >, and quotes for safe XML attribute values
      # serialization

svgbob
  svgbob.parse_grid
    @ (source: string) -> grid
    + builds a row-major grid of characters from the source text
    + pads short lines so every row has the same width
    # parsing
    -> std.text.split_lines
  svgbob.classify_cell
    @ (g: grid, row: i32, col: i32) -> cell_kind
    + classifies a cell as line segment, corner, arrow head, text, or empty
    - returns empty for out-of-bounds coordinates
    # classification
  svgbob.trace_shapes
    @ (g: grid) -> list[shape]
    + walks the grid and emits connected line runs, boxes, and arrows
    + preserves text runs as labeled shapes
    # shape_extraction
  svgbob.render_svg
    @ (shapes: list[shape], cell_w: i32, cell_h: i32) -> string
    + emits an SVG document sized to the shape bounds
    + each shape becomes a path, rect, or text element
    # rendering
    -> std.xml.escape_attr
  svgbob.convert
    @ (source: string) -> string
    + end-to-end: parses, traces, and renders ASCII art to SVG
    # pipeline
