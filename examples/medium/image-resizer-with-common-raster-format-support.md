# Requirement: "an image resizing library supporting common raster formats"

Decodes an image from one of several formats, resamples it to the target dimensions, and re-encodes it. Format detection comes from the magic bytes in the input.

std
  std.image
    std.image.decode_jpeg
      @ (data: bytes) -> result[raster, string]
      + decodes a JPEG image to an RGBA raster
      - returns error on corrupt or truncated data
      # decoding
    std.image.decode_png
      @ (data: bytes) -> result[raster, string]
      + decodes a PNG image to an RGBA raster
      - returns error on corrupt or truncated data
      # decoding
    std.image.encode_jpeg
      @ (r: raster, quality: i32) -> bytes
      + encodes a raster as JPEG at the given quality (1-100)
      # encoding
    std.image.encode_png
      @ (r: raster) -> bytes
      + encodes a raster as PNG
      # encoding

image_resize
  image_resize.detect_format
    @ (data: bytes) -> result[image_format, string]
    + identifies JPEG, PNG, WebP, or TIFF from magic bytes
    - returns error when the header is not recognized
    # detection
  image_resize.decode
    @ (data: bytes) -> result[raster, string]
    + decodes an image into an RGBA raster regardless of format
    - returns error when decoding fails
    # decoding
    -> std.image.decode_jpeg
    -> std.image.decode_png
  image_resize.resample
    @ (src: raster, width: i32, height: i32) -> raster
    + returns a new raster resized with a high-quality Lanczos filter
    ? aspect ratio is not preserved automatically
    # resampling
  image_resize.fit
    @ (src: raster, max_w: i32, max_h: i32) -> raster
    + returns a resized raster that fits within the box, preserving aspect ratio
    # fit
  image_resize.encode
    @ (r: raster, format: image_format, quality: i32) -> result[bytes, string]
    + encodes a raster in the requested format
    - returns error when the format does not support the given options
    # encoding
    -> std.image.encode_jpeg
    -> std.image.encode_png
