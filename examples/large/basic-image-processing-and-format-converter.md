# Requirement: "a library for basic image processing and conversion between image formats"

Core decode/encode primitives live in std; the project layer exposes pixel-level operations over a decoded image.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns full file contents as bytes
      - returns error when path cannot be read
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to path, replacing any existing file
      - returns error when path cannot be written
      # filesystem
  std.image
    std.image.decode_png
      @ (data: bytes) -> result[image_data, string]
      + decodes PNG bytes into an RGBA pixel buffer with width and height
      - returns error on truncated or invalid PNG
      # image_decoding
    std.image.decode_jpeg
      @ (data: bytes) -> result[image_data, string]
      + decodes JPEG bytes into an RGBA pixel buffer with width and height
      - returns error on invalid or unsupported JPEG
      # image_decoding
    std.image.encode_png
      @ (img: image_data) -> result[bytes, string]
      + encodes an RGBA pixel buffer as PNG bytes
      # image_encoding
    std.image.encode_jpeg
      @ (img: image_data, quality: i32) -> result[bytes, string]
      + encodes an RGBA pixel buffer as JPEG bytes at the given quality (1..100)
      - returns error when quality is out of range
      # image_encoding

imaging
  imaging.load
    @ (path: string) -> result[image_data, string]
    + loads an image from disk, detecting PNG or JPEG by content
    - returns error when format cannot be detected
    # loading
    -> std.fs.read_all
    -> std.image.decode_png
    -> std.image.decode_jpeg
  imaging.save
    @ (img: image_data, path: string, format: string) -> result[void, string]
    + saves the image to disk in the requested format ("png" or "jpeg")
    - returns error on unknown format
    # saving
    -> std.image.encode_png
    -> std.image.encode_jpeg
    -> std.fs.write_all
  imaging.convert
    @ (src_path: string, dst_path: string, dst_format: string) -> result[void, string]
    + loads src, re-encodes it in dst_format, and writes dst
    - returns error when src cannot be read or dst_format is unknown
    # conversion
  imaging.resize
    @ (img: image_data, width: i32, height: i32) -> image_data
    + returns a new image scaled to the given dimensions using bilinear sampling
    ? zero or negative dimensions are clamped to 1
    # resizing
  imaging.crop
    @ (img: image_data, x: i32, y: i32, width: i32, height: i32) -> result[image_data, string]
    + returns a new image containing the given rectangle
    - returns error when the rectangle extends outside the source
    # cropping
  imaging.grayscale
    @ (img: image_data) -> image_data
    + returns a new image with each pixel converted to luminance-preserving gray
    # color_transform
  imaging.flip_horizontal
    @ (img: image_data) -> image_data
    + returns a new image mirrored along the vertical axis
    # geometry
  imaging.rotate_90
    @ (img: image_data) -> image_data
    + returns a new image rotated 90 degrees clockwise
    # geometry
