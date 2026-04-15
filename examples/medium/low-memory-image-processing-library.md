# Requirement: "an image processing library with low memory usage"

A streaming pipeline that decodes, transforms, and encodes images in tiles rather than full-frame buffers.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the path does not exist or is unreadable
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or truncating
      # filesystem

image
  image.decode
    fn (data: bytes) -> result[image_source, string]
    + returns a lazy source that yields tiles without decoding the full image
    + auto-detects jpeg, png, and webp from the magic bytes
    - returns error on unknown or corrupt format
    # decoding
  image.resize
    fn (src: image_source, width: i32, height: i32) -> image_source
    + returns a source that resamples tiles to the target dimensions using lanczos
    # transform
  image.crop
    fn (src: image_source, x: i32, y: i32, width: i32, height: i32) -> image_source
    + returns a source restricted to the given rectangle
    - returns a zero-size source when the rectangle is empty
    # transform
  image.convert_colorspace
    fn (src: image_source, target: i32) -> image_source
    + returns a source in the target colorspace (srgb, linear, grayscale)
    # transform
  image.encode
    fn (src: image_source, format: i32, quality: i32) -> result[bytes, string]
    + pulls tiles from src and encodes to jpeg, png, or webp
    + quality is ignored for lossless formats
    - returns error when quality is outside 1..100 for lossy formats
    # encoding
  image.load_file
    fn (path: string) -> result[image_source, string]
    + reads a file and decodes it
    # convenience
    -> std.fs.read_all
  image.save_file
    fn (src: image_source, path: string, format: i32, quality: i32) -> result[void, string]
    + encodes and writes the result to path
    # convenience
    -> std.fs.write_all
