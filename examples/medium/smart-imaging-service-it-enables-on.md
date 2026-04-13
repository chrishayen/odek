# Requirement: "an on-demand image crop, resize, and flip library"

Decodes pixel buffers, applies geometry transforms, and re-encodes. Format codecs live in std.

std
  std.image
    std.image.decode
      @ (data: bytes) -> result[image, string]
      + decodes a PNG or JPEG byte blob into a pixel buffer
      - returns error on unknown format or truncated data
      # image_decode
    std.image.encode_png
      @ (img: image) -> bytes
      + encodes a pixel buffer as PNG
      # image_encode
    std.image.encode_jpeg
      @ (img: image, quality: i32) -> result[bytes, string]
      + encodes a pixel buffer as JPEG at the given quality (1-100)
      - returns error when quality is out of range
      # image_encode
  std.image.geom
    std.image.geom.width
      @ (img: image) -> i32
      + returns the image width in pixels
      # image_inspect
    std.image.geom.height
      @ (img: image) -> i32
      + returns the image height in pixels
      # image_inspect

imaging
  imaging.crop
    @ (img: image, x: i32, y: i32, w: i32, h: i32) -> result[image, string]
    + returns a sub-region of the image
    - returns error when the rectangle is outside the source bounds
    - returns error when w or h is non-positive
    # transform_crop
    -> std.image.geom.width
    -> std.image.geom.height
  imaging.resize
    @ (img: image, target_w: i32, target_h: i32) -> result[image, string]
    + scales the image to the exact target dimensions with bilinear filtering
    - returns error when target dimensions are non-positive
    # transform_resize
  imaging.fit
    @ (img: image, max_w: i32, max_h: i32) -> image
    + scales down proportionally so the image fits inside the box, preserving aspect ratio
    ? upscales are not performed; images smaller than the box are returned unchanged
    # transform_resize
    -> std.image.geom.width
    -> std.image.geom.height
  imaging.flip_horizontal
    @ (img: image) -> image
    + returns a mirror image along the vertical axis
    # transform_flip
  imaging.flip_vertical
    @ (img: image) -> image
    + returns a mirror image along the horizontal axis
    # transform_flip
  imaging.process
    @ (source: bytes, ops: list[imaging_op], out_format: string, quality: i32) -> result[bytes, string]
    + decodes, applies each op in order, and re-encodes as "png" or "jpeg"
    - returns error on any stage failure
    # pipeline
    -> std.image.decode
    -> std.image.encode_png
    -> std.image.encode_jpeg
