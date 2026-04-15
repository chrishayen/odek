# Requirement: "a QR code and barcode generator"

Encodes text into a matrix (QR) or line pattern (barcode) and renders to common output formats.

std
  std.image
    std.image.new_bitmap
      fn (width: i32, height: i32) -> bitmap
      + returns a bitmap filled with white pixels
      # image
    std.image.set_pixel
      fn (bm: bitmap, x: i32, y: i32, black: bool) -> bitmap
      + sets a pixel to black or white
      # image
    std.image.encode_png
      fn (bm: bitmap) -> bytes
      + encodes the bitmap as PNG
      # image

code
  code.encode_qr
    fn (text: string, error_level: string) -> result[matrix, string]
    + returns a square matrix of booleans representing QR modules
    + error_level is one of "L", "M", "Q", "H"
    - returns error when text exceeds the capacity for the given level
    # qr_encoding
  code.encode_barcode
    fn (symbology: string, value: string) -> result[list[bool], string]
    + returns a 1-D bar pattern for the given symbology
    + supports "code128", "code39", and "ean13"
    - returns error when value has invalid characters for the symbology
    # barcode_encoding
  code.render_qr_png
    fn (matrix: matrix, module_pixels: i32) -> bytes
    + renders a QR matrix as a PNG with the given module size in pixels
    # rendering
    -> std.image.new_bitmap
    -> std.image.set_pixel
    -> std.image.encode_png
  code.render_barcode_png
    fn (pattern: list[bool], height: i32, module_pixels: i32) -> bytes
    + renders a 1-D bar pattern as a PNG of the given height
    # rendering
    -> std.image.new_bitmap
    -> std.image.set_pixel
    -> std.image.encode_png
  code.render_qr_svg
    fn (matrix: matrix, module_pixels: i32) -> string
    + renders a QR matrix as an SVG document
    # rendering
