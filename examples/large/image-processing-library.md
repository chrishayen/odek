# Requirement: "an image processing library"

Core decode/encode through std codec primitives; project layer exposes common pixel operations on an in-memory image.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of a file
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating
      - returns error on permission failure
      # filesystem
  std.codec
    std.codec.decode_png
      fn (data: bytes) -> result[image_buffer, string]
      + decodes a PNG byte stream into an RGBA pixel buffer
      - returns error on truncated or corrupt input
      # decoding
    std.codec.encode_png
      fn (img: image_buffer) -> result[bytes, string]
      + encodes an RGBA pixel buffer as PNG
      # encoding
    std.codec.decode_jpeg
      fn (data: bytes) -> result[image_buffer, string]
      + decodes a JPEG byte stream into an RGBA pixel buffer
      - returns error on invalid markers
      # decoding
    std.codec.encode_jpeg
      fn (img: image_buffer, quality: i32) -> result[bytes, string]
      + encodes an RGBA buffer as JPEG at the given quality
      - returns error when quality is outside 1..100
      # encoding

image
  image.load
    fn (path: string) -> result[image_buffer, string]
    + returns an image_buffer decoded from the file
    - returns error when format is unrecognized
    # loading
    -> std.fs.read_all
    -> std.codec.decode_png
    -> std.codec.decode_jpeg
  image.save
    fn (img: image_buffer, path: string) -> result[void, string]
    + writes the image to disk, choosing encoder by file extension
    - returns error when extension is unsupported
    # saving
    -> std.codec.encode_png
    -> std.codec.encode_jpeg
    -> std.fs.write_all
  image.resize
    fn (img: image_buffer, width: i32, height: i32) -> image_buffer
    + returns a new buffer scaled to the target dimensions using bilinear sampling
    ? aspect ratio is not preserved; callers compute target dims themselves
    # resize
  image.crop
    fn (img: image_buffer, x: i32, y: i32, width: i32, height: i32) -> result[image_buffer, string]
    + returns a sub-region as a new buffer
    - returns error when the region exceeds the source bounds
    # crop
  image.rotate
    fn (img: image_buffer, degrees: f64) -> image_buffer
    + returns a new buffer rotated by the given angle with transparent fill
    # rotate
  image.flip_horizontal
    fn (img: image_buffer) -> image_buffer
    + returns a new buffer mirrored left-to-right
    # flip
  image.flip_vertical
    fn (img: image_buffer) -> image_buffer
    + returns a new buffer mirrored top-to-bottom
    # flip
  image.grayscale
    fn (img: image_buffer) -> image_buffer
    + returns a new buffer with each pixel converted to luminance
    # color
  image.blur
    fn (img: image_buffer, radius: f64) -> image_buffer
    + returns a new buffer with a gaussian blur of the given radius
    - returns the input unchanged when radius is 0
    # filter
  image.composite
    fn (base: image_buffer, overlay: image_buffer, x: i32, y: i32) -> image_buffer
    + returns a new buffer with overlay alpha-blended onto base at (x, y)
    # compositing
