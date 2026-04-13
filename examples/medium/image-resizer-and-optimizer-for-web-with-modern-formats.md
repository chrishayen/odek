# Requirement: "a library that resizes and optimizes images for the web using modern formats"

Load an image, resize it, and re-encode it as a modern format with quality settings.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when the path cannot be opened
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path
      - returns error when the path cannot be created
      # io
  std.image
    std.image.decode
      @ (data: bytes) -> result[image, string]
      + decodes png, jpeg, gif, or webp based on the magic bytes
      - returns error when the format is unrecognized
      # image
    std.image.encode
      @ (img: image, format: string, quality: i32) -> result[bytes, string]
      + encodes an image as jpeg, webp, or avif with the given quality (0-100)
      - returns error when format is unknown or quality is out of range
      # image

image_optim
  image_optim.resize
    @ (img: image, width: i32, height: i32) -> image
    + returns a resampled image using bicubic filtering
    + preserves aspect ratio when one dimension is 0
    # resizing
  image_optim.fit_within
    @ (img: image, max_width: i32, max_height: i32) -> image
    + returns the image scaled down to fit within the bounds, leaving small images unchanged
    # resizing
  image_optim.strip_metadata
    @ (img: image) -> image
    + returns a copy with EXIF and color profile metadata removed
    # optimization
  image_optim.best_format_for
    @ (accept_header: string) -> string
    + returns "avif", "webp", or "jpeg" based on the HTTP Accept header
    + returns "jpeg" when no modern format is accepted
    # negotiation
  image_optim.process_file
    @ (input_path: string, output_path: string, max_width: i32, max_height: i32, format: string, quality: i32) -> result[void, string]
    + reads, resizes, strips metadata, re-encodes, and writes the image
    - returns error when any step fails
    # pipeline
    -> std.fs.read_all
    -> std.image.decode
    -> std.image.encode
    -> std.fs.write_all
