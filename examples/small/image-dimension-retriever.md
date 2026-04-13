# Requirement: "a library that returns the pixel dimensions of an image"

Reads only the header bytes needed to determine width and height for common raster formats.

std: (all units exist)

image_dimensions
  image_dimensions.detect_format
    @ (data: bytes) -> result[string, string]
    + returns "png", "jpeg", "gif", "webp", or "bmp" based on the magic bytes
    - returns error when no known format is detected
    # format_detection
  image_dimensions.png_size
    @ (data: bytes) -> result[tuple[i32, i32], string]
    + returns (width, height) from the IHDR chunk
    - returns error when the IHDR chunk is missing or truncated
    # png
  image_dimensions.jpeg_size
    @ (data: bytes) -> result[tuple[i32, i32], string]
    + returns (width, height) from the first SOFn marker
    - returns error when no SOFn marker is present
    # jpeg
  image_dimensions.gif_size
    @ (data: bytes) -> result[tuple[i32, i32], string]
    + returns (width, height) from the logical screen descriptor
    - returns error when data is shorter than the GIF header
    # gif
  image_dimensions.size
    @ (data: bytes) -> result[tuple[i32, i32], string]
    + dispatches on detected format and returns (width, height)
    - returns error for unsupported or malformed input
    # entry
