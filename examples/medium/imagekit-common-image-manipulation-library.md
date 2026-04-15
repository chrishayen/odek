# Requirement: "an image manipulation library with a small set of common operations"

Decodes an image buffer, runs pixel-wise operations on an in-memory RGBA surface, and re-encodes the result.

std
  std.image
    std.image.decode
      fn (data: bytes) -> result[image_surface, string]
      + decodes PNG or JPEG bytes into an RGBA surface
      - returns error on unrecognized format
      # decoding
    std.image.encode_png
      fn (surf: image_surface) -> bytes
      + encodes the surface as PNG
      # encoding
    std.image.encode_jpeg
      fn (surf: image_surface, quality: i32) -> bytes
      + encodes the surface as JPEG at the given quality
      # encoding

imagekit
  imagekit.resize
    fn (surf: image_surface, width: i32, height: i32) -> result[image_surface, string]
    + returns a new surface scaled to the target dimensions with bilinear sampling
    - returns error when either dimension is less than 1
    # transform
  imagekit.crop
    fn (surf: image_surface, x: i32, y: i32, width: i32, height: i32) -> result[image_surface, string]
    + returns the rectangular region as a new surface
    - returns error when the rectangle is outside the source bounds
    # transform
  imagekit.rotate
    fn (surf: image_surface, degrees: i32) -> result[image_surface, string]
    + rotates by 90, 180, or 270 degrees clockwise
    - returns error for any other angle
    # transform
  imagekit.flip_horizontal
    fn (surf: image_surface) -> image_surface
    + mirrors the surface left-to-right
    # transform
  imagekit.grayscale
    fn (surf: image_surface) -> image_surface
    + converts each pixel to its luminance in all three channels
    # filter
  imagekit.blur
    fn (surf: image_surface, radius: i32) -> result[image_surface, string]
    + applies a box blur with the given radius in pixels
    - returns error when radius is negative
    # filter
  imagekit.adjust_brightness
    fn (surf: image_surface, delta: i32) -> image_surface
    + adds delta to each channel clamped to 0..255
    # filter
