# Requirement: "an image processing library with filters, convolutions, and geometric transforms"

Operates on an in-memory raster image. Includes per-pixel transformations, convolution with arbitrary kernels, and common geometric operations.

std
  std.math
    std.math.clamp_f64
      fn (x: f64, lo: f64, hi: f64) -> f64
      + returns x clamped to [lo, hi]
      # math
    std.math.round
      fn (x: f64) -> i64
      + rounds to the nearest integer
      # math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root
      # math

imageproc
  imageproc.new_image
    fn (width: i32, height: i32) -> image
    + creates a zero-initialized RGBA image of the given size
    - errors when width or height is negative
    # construction
  imageproc.get_pixel
    fn (img: image, x: i32, y: i32) -> optional[pixel]
    + returns the pixel at (x, y)
    - returns none when out of bounds
    # access
  imageproc.set_pixel
    fn (img: image, x: i32, y: i32, p: pixel) -> image
    + returns a new image with the pixel replaced
    - returns unchanged when out of bounds
    # access
  imageproc.to_grayscale
    fn (img: image) -> image
    + returns a grayscale copy using luminance weighting
    # filter
  imageproc.invert
    fn (img: image) -> image
    + returns an image with colors inverted
    # filter
  imageproc.adjust_brightness
    fn (img: image, delta: i32) -> image
    + returns an image with each channel shifted, clamped to [0, 255]
    # filter
    -> std.math.clamp_f64
  imageproc.convolve
    fn (img: image, kernel: list[list[f64]]) -> image
    + applies an arbitrary 2D kernel; edges clamp to the nearest in-bounds pixel
    - errors when the kernel is empty or not rectangular
    # convolution
    -> std.math.clamp_f64
  imageproc.box_blur
    fn (img: image, radius: i32) -> image
    + applies a box blur of the given radius
    # filter
    -> std.math.clamp_f64
  imageproc.sobel_edges
    fn (img: image) -> image
    + returns an edge-magnitude image using Sobel operators
    # edge
    -> std.math.sqrt
  imageproc.resize_nearest
    fn (img: image, new_width: i32, new_height: i32) -> image
    + returns a resized image using nearest-neighbor sampling
    # geometry
  imageproc.resize_bilinear
    fn (img: image, new_width: i32, new_height: i32) -> image
    + returns a resized image using bilinear interpolation
    # geometry
    -> std.math.round
  imageproc.rotate_90
    fn (img: image) -> image
    + returns the image rotated 90 degrees clockwise
    # geometry
  imageproc.crop
    fn (img: image, x: i32, y: i32, w: i32, h: i32) -> result[image, string]
    + returns a sub-image of the given region
    - returns error when the region exits the image bounds
    # geometry
  imageproc.threshold
    fn (img: image, cutoff: u8) -> image
    + returns a binarized image based on luminance
    # filter
