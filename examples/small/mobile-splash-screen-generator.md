# Requirement: "a mobile app splash screen generator"

Given a source image and a list of target sizes, produces resized splash images padded to aspect.

std
  std.image
    std.image.decode
      fn (data: bytes) -> result[image, string]
      + decodes PNG or JPEG bytes into an image
      - returns error on unsupported or malformed data
      # image_decoding
    std.image.encode_png
      fn (img: image) -> bytes
      + encodes an image as PNG
      # image_encoding
    std.image.resize
      fn (img: image, width: i32, height: i32) -> image
      + returns a resized image using bilinear sampling
      # image_resize

splash
  splash.render
    fn (source: bytes, width: i32, height: i32, background: u32) -> result[bytes, string]
    + decodes the source, fits it centered into a (width, height) canvas filled with background, and encodes PNG
    - returns error when the source image cannot be decoded
    ? aspect is preserved; letterbox bars use background
    # rendering
    -> std.image.decode
    -> std.image.resize
    -> std.image.encode_png
  splash.render_many
    fn (source: bytes, sizes: list[tuple[i32, i32]], background: u32) -> result[list[bytes], string]
    + renders one PNG per requested size
    - returns error when the source image cannot be decoded
    # batch_rendering
