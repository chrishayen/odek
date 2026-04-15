# Requirement: "display images in the terminal"

Decodes common image formats, resizes to terminal cell dimensions, and emits ansi half-block or sixel output depending on capability.

std
  std.image
    std.image.decode
      fn (data: bytes) -> result[image, string]
      + decodes png/jpeg/gif into an rgba pixel buffer
      - returns error on unknown or corrupt format
      # image
    std.image.resize
      fn (img: image, width: i32, height: i32) -> image
      + returns a bilinearly resampled copy at the target size
      # image
  std.term
    std.term.size
      fn () -> tuple[i32, i32]
      + returns terminal (cols, rows)
      # terminal
    std.term.supports_sixel
      fn () -> bool
      + returns true when the host terminal advertises sixel graphics
      # terminal

terminal_image
  terminal_image.render
    fn (data: bytes, max_cols: i32, max_rows: i32) -> result[string, string]
    + decodes, scales to fit within the given cell box, and returns renderable ansi text
    + uses sixel when the terminal supports it, otherwise unicode half blocks with truecolor
    - returns error on unknown image format
    # rendering
    -> std.image.decode
    -> std.image.resize
    -> std.term.supports_sixel
  terminal_image.render_for_current_terminal
    fn (data: bytes) -> result[string, string]
    + convenience wrapper that queries terminal size and renders to fit
    - returns error when the terminal size cannot be determined
    # rendering
    -> std.term.size
  terminal_image.render_half_blocks
    fn (img: image) -> string
    + returns ansi text where each terminal cell encodes two vertical pixels using foreground and background truecolor
    # rendering
  terminal_image.render_sixel
    fn (img: image) -> string
    + returns a sixel-encoded escape sequence for the image
    # rendering
