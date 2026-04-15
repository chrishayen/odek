# Requirement: "an image processing library"

Decodes, resizes, crops, and re-encodes raster images without relying on an external command-line tool.

std
  std.fs
    std.fs.read_all_bytes
      fn (path: string) -> result[bytes, string]
      + returns the full contents of a file as bytes
      - returns error when the file is missing
      # filesystem
    std.fs.write_all_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating it
      - returns error when the parent directory does not exist
      # filesystem

image_proc
  image_proc.decode
    fn (data: bytes) -> result[image, string]
    + decodes png or jpeg byte streams into an image
    - returns error on unrecognized or truncated input
    # decoding
  image_proc.encode_png
    fn (img: image) -> bytes
    + encodes an image as a png byte stream
    # encoding
  image_proc.encode_jpeg
    fn (img: image, quality: i32) -> bytes
    + encodes an image as a jpeg byte stream at the given quality
    ? quality is clamped to [1, 100]
    # encoding
  image_proc.resize
    fn (img: image, width: i32, height: i32) -> image
    + returns a new image scaled to the target dimensions
    ? uses bilinear sampling
    # resize
  image_proc.crop
    fn (img: image, x: i32, y: i32, width: i32, height: i32) -> result[image, string]
    + returns a new image containing the specified rectangle
    - returns error when the rectangle falls outside the source bounds
    # crop
  image_proc.load_file
    fn (path: string) -> result[image, string]
    + reads and decodes an image from disk
    - returns error on missing or malformed file
    # io
    -> std.fs.read_all_bytes
  image_proc.save_png
    fn (img: image, path: string) -> result[void, string]
    + encodes an image as png and writes it to disk
    # io
    -> std.fs.write_all_bytes
