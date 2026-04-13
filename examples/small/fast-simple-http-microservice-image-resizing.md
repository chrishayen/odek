# Requirement: "an image resizing library"

Decode, resize, and re-encode image bytes. The std image codec handles format detection.

std
  std.image
    std.image.decode
      @ (data: bytes) -> result[raw_image, string]
      + decodes image bytes into a raw pixel buffer with width, height, and channels
      - returns error on unknown or corrupt format
      # image
    std.image.encode
      @ (img: raw_image, format: string) -> result[bytes, string]
      + encodes a raw image in the named format ("png", "jpeg", "webp")
      - returns error on unknown format
      # image

resize
  resize.bilinear
    @ (img: raw_image, width: i32, height: i32) -> result[raw_image, string]
    + resamples an image to the target dimensions using bilinear interpolation
    - returns error when width or height is non-positive
    # resampling
  resize.fit
    @ (img: raw_image, max_width: i32, max_height: i32) -> raw_image
    + resizes preserving aspect ratio so the result fits inside the box
    + returns the input unchanged when it already fits
    # resampling
  resize.process_bytes
    @ (data: bytes, width: i32, height: i32, format: string) -> result[bytes, string]
    + decodes input bytes, resizes to exact dimensions, and re-encodes
    - returns error on decode, resize, or encode failure
    # pipeline
    -> std.image.decode
    -> std.image.encode
