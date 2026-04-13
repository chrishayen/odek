# Requirement: "an image processing library"

Same shape as any image library: std owns codecs and io, project layer exposes pixel-level transforms.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full file contents
      - returns error when path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # filesystem
  std.codec
    std.codec.decode_png
      @ (data: bytes) -> result[image_buffer, string]
      + decodes PNG bytes to an RGBA buffer
      - returns error on corrupt input
      # decoding
    std.codec.encode_png
      @ (img: image_buffer) -> result[bytes, string]
      + encodes an RGBA buffer as PNG
      # encoding
    std.codec.decode_bmp
      @ (data: bytes) -> result[image_buffer, string]
      + decodes BMP bytes to an RGBA buffer
      - returns error when header is malformed
      # decoding
    std.codec.encode_bmp
      @ (img: image_buffer) -> result[bytes, string]
      + encodes an RGBA buffer as BMP
      # encoding

image
  image.new
    @ (width: i32, height: i32, fill: u32) -> image_buffer
    + returns a blank buffer with every pixel set to fill (packed RGBA)
    # construction
  image.load
    @ (path: string) -> result[image_buffer, string]
    + returns an image decoded from disk, choosing format by extension
    - returns error when extension is unsupported
    # loading
    -> std.fs.read_all
    -> std.codec.decode_png
    -> std.codec.decode_bmp
  image.save
    @ (img: image_buffer, path: string) -> result[void, string]
    + writes the image to disk, choosing encoder by extension
    # saving
    -> std.codec.encode_png
    -> std.codec.encode_bmp
    -> std.fs.write_all
  image.get_pixel
    @ (img: image_buffer, x: i32, y: i32) -> result[u32, string]
    + returns the packed RGBA value at (x, y)
    - returns error when coordinates are out of bounds
    # pixel_access
  image.set_pixel
    @ (img: image_buffer, x: i32, y: i32, color: u32) -> result[image_buffer, string]
    + returns a new buffer with pixel at (x, y) replaced
    - returns error when coordinates are out of bounds
    # pixel_access
  image.resize
    @ (img: image_buffer, width: i32, height: i32) -> image_buffer
    + returns a new buffer bilinearly scaled to the target size
    # resize
  image.crop
    @ (img: image_buffer, x: i32, y: i32, w: i32, h: i32) -> result[image_buffer, string]
    + returns a sub-region as a new buffer
    - returns error when region exceeds source bounds
    # crop
  image.invert
    @ (img: image_buffer) -> image_buffer
    + returns a new buffer with each channel replaced by 255 - value
    # color
  image.grayscale
    @ (img: image_buffer) -> image_buffer
    + returns a new buffer with each pixel replaced by its luminance
    # color
  image.composite
    @ (base: image_buffer, overlay: image_buffer, x: i32, y: i32) -> image_buffer
    + returns a new buffer with overlay alpha-blended onto base at (x, y)
    # compositing
