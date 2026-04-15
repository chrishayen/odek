# Requirement: "a simple image processing library"

Core pixel-buffer operations: load, resize, rotate, crop, save. Decoding and encoding are thin std primitives.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path, overwriting
      - returns error when the parent directory is missing
      # io
  std.image_codec
    std.image_codec.decode
      fn (data: bytes) -> result[image_buffer, string]
      + decodes common raster formats (PNG/JPEG) into an rgba image buffer
      - returns error on unknown format signatures
      # codec
    std.image_codec.encode
      fn (img: image_buffer, format: string) -> result[bytes, string]
      + encodes an image buffer into the requested format
      - returns error when format is unsupported
      # codec

imaging
  imaging.load
    fn (path: string) -> result[image_buffer, string]
    + loads an image from disk into an rgba buffer
    - returns error when decoding fails
    # loading
    -> std.fs.read_all
    -> std.image_codec.decode
  imaging.save
    fn (img: image_buffer, path: string, format: string) -> result[void, string]
    + encodes and writes an image to disk
    - returns error when encoding fails
    # saving
    -> std.image_codec.encode
    -> std.fs.write_all
  imaging.resize
    fn (img: image_buffer, width: i32, height: i32) -> image_buffer
    + returns a new buffer resampled to the target dimensions using bilinear filtering
    ? target dimensions must be positive; zero values fall back to preserving aspect from the other axis
    # resizing
  imaging.crop
    fn (img: image_buffer, x: i32, y: i32, width: i32, height: i32) -> result[image_buffer, string]
    + returns a new buffer covering the given rectangle
    - returns error when the rectangle extends outside the source
    # cropping
  imaging.rotate
    fn (img: image_buffer, degrees: f64) -> image_buffer
    + rotates by an arbitrary angle, expanding the canvas to fit
    + rotations of 0/90/180/270 degrees use exact integer paths
    # rotation
  imaging.flip_horizontal
    fn (img: image_buffer) -> image_buffer
    + returns a new buffer mirrored along the vertical axis
    # flipping
