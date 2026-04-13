# Requirement: "a computer vision library for loading, transforming, and analyzing images"

Core operations: load, convert color spaces, resize, filter, detect edges, and find contours. Image encoding/decoding is delegated to std.

std
  std.image_codec
    std.image_codec.decode
      @ (data: bytes) -> result[raw_image, string]
      + returns a raw image with width, height, channels, and pixel buffer
      - returns error on unsupported or malformed data
      # image_codec
    std.image_codec.encode_png
      @ (img: raw_image) -> result[bytes, string]
      + encodes a raw image as PNG
      # image_codec
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      - returns error when the file is missing
      # filesystem

vision
  vision.load
    @ (path: string) -> result[raw_image, string]
    + loads and decodes an image from disk
    # loading
    -> std.fs.read_all
    -> std.image_codec.decode
  vision.to_grayscale
    @ (img: raw_image) -> raw_image
    + returns a single-channel image using luma weights
    # color_conversion
  vision.resize
    @ (img: raw_image, width: i32, height: i32) -> raw_image
    + returns a bilinearly resampled image at the target size
    # resampling
  vision.gaussian_blur
    @ (img: raw_image, sigma: f32) -> raw_image
    + returns the image convolved with a separable Gaussian kernel
    # filtering
  vision.threshold
    @ (img: raw_image, value: u8) -> raw_image
    + returns a binary image where pixels >= value become 255 and others 0
    ? input must be single-channel
    # segmentation
  vision.detect_edges
    @ (img: raw_image, low: f32, high: f32) -> raw_image
    + returns an edge map using hysteresis thresholding
    # edge_detection
  vision.find_contours
    @ (binary: raw_image) -> list[contour]
    + returns outer contours of connected components in the binary image
    # analysis
  vision.save
    @ (img: raw_image, path: string) -> result[void, string]
    + encodes as PNG and writes to disk
    # saving
    -> std.image_codec.encode_png
