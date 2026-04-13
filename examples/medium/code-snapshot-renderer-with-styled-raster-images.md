# Requirement: "renders source code snippets as styled raster images"

Tokenizes code, lays it out, and paints it into a pixel buffer. PNG encoding lives in std.

std
  std.image
    std.image.new_canvas
      @ (width: i32, height: i32, background: rgb_color) -> image_data
      + returns a fresh RGBA pixel buffer filled with background
      # image_buffer
    std.image.draw_text
      @ (img: image_data, x: i32, y: i32, text: string, color: rgb_color, font: font_handle) -> void
      + draws text at (x, y) using the given font and color
      # image_drawing
    std.image.load_font
      @ (path: string, size: i32) -> result[font_handle, string]
      + loads a font file at the given pixel size
      - returns error when the file cannot be read or parsed
      # fonts
    std.image.encode_png
      @ (img: image_data) -> result[bytes, string]
      + encodes an image buffer as PNG bytes
      # image_encoding

code_snapshot
  code_snapshot.tokenize
    @ (source: string, language: string) -> list[code_token]
    + returns tokens tagged as keyword, string, number, comment, or plain
    ? unknown languages return a single plain token per line
    # lexing
  code_snapshot.theme_color
    @ (theme: code_theme, kind: string) -> rgb_color
    + returns the color assigned to a token kind in the theme
    ? unknown kinds fall back to the foreground color
    # theming
  code_snapshot.measure
    @ (source: string, font: font_handle, padding: i32) -> tuple[i32, i32]
    + returns the canvas width and height required to render source
    # layout
  code_snapshot.render
    @ (source: string, language: string, theme: code_theme, font: font_handle) -> image_data
    + returns an image with syntax-highlighted source on the theme background
    # rendering
    -> std.image.new_canvas
    -> std.image.draw_text
  code_snapshot.render_png
    @ (source: string, language: string, theme: code_theme, font_path: string) -> result[bytes, string]
    + convenience that loads the font, renders, and encodes as PNG
    - returns error when the font cannot be loaded
    # export
    -> std.image.load_font
    -> std.image.encode_png
