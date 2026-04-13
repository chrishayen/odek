# Requirement: "an image processing library"

Core transformations on decoded image buffers. Decoding and encoding are delegated to std primitives.

std
  std.image
    std.image.decode
      @ (data: bytes) -> result[image_buffer, string]
      + decodes common image formats into a pixel buffer
      - returns error on unknown or corrupt input
      # image
    std.image.encode
      @ (img: image_buffer, format: string) -> result[bytes, string]
      + encodes a pixel buffer to the requested format
      - returns error when the format is not supported
      # image

image_ops
  image_ops.resize
    @ (img: image_buffer, width: i32, height: i32) -> image_buffer
    + returns a new buffer scaled to the target dimensions
    ? resizing uses bilinear filtering
    # resize
  image_ops.crop
    @ (img: image_buffer, x: i32, y: i32, width: i32, height: i32) -> result[image_buffer, string]
    + returns a subregion of the image
    - returns error when the rectangle is outside the source bounds
    # crop
  image_ops.rotate
    @ (img: image_buffer, degrees: i32) -> image_buffer
    + returns the image rotated by 90, 180, or 270 degrees
    ? other angles are snapped to the nearest quadrant
    # rotate
  image_ops.grayscale
    @ (img: image_buffer) -> image_buffer
    + returns a single-channel version of the image
    # color
  image_ops.process
    @ (input: bytes, width: i32, height: i32, format: string) -> result[bytes, string]
    + decodes, resizes, and re-encodes in one call
    - returns error when any step fails
    # pipeline
    -> std.image.decode
    -> std.image.encode
