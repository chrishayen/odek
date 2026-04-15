# Requirement: "a scientific image processing library"

Operates on 2D and 3D image arrays with filtering, morphology, segmentation, and feature extraction.

std
  std.math
    std.math.sqrt_f64
      fn (x: f64) -> f64
      + returns the square root of x
      - returns NaN for negative input
      # math
    std.math.exp_f64
      fn (x: f64) -> f64
      + returns e raised to x
      # math
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file's contents as bytes
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path, overwriting existing files
      - returns error when the parent directory does not exist
      # filesystem

image
  image.new
    fn (width: i32, height: i32, channels: i32) -> image_state
    + creates a zero-filled image with the given shape
    # construction
  image.load_png
    fn (path: string) -> result[image_state, string]
    + loads a PNG from disk
    - returns error when the file is not a valid PNG
    # io
    -> std.fs.read_all
  image.save_png
    fn (img: image_state, path: string) -> result[void, string]
    + writes the image to disk as PNG
    - returns error when encoding fails
    # io
    -> std.fs.write_all
  image.to_grayscale
    fn (img: image_state) -> image_state
    + returns a single-channel image from a multi-channel input
    # conversion
  image.gaussian_blur
    fn (img: image_state, sigma: f64) -> image_state
    + convolves the image with a Gaussian kernel of the given standard deviation
    ? kernel radius is derived from sigma
    # filtering
    -> std.math.exp_f64
  image.sobel_edges
    fn (img: image_state) -> image_state
    + returns the gradient magnitude from Sobel operators
    # filtering
    -> std.math.sqrt_f64
  image.threshold_otsu
    fn (img: image_state) -> image_state
    + returns a binary image thresholded at the Otsu optimum
    # segmentation
  image.morphology_dilate
    fn (img: image_state, radius: i32) -> image_state
    + dilates a binary image by a disk structuring element
    # morphology
  image.morphology_erode
    fn (img: image_state, radius: i32) -> image_state
    + erodes a binary image by a disk structuring element
    # morphology
  image.connected_components
    fn (img: image_state) -> tuple[image_state, i32]
    + labels each connected foreground region and returns the label count
    # segmentation
  image.watershed_segment
    fn (img: image_state, markers: image_state) -> image_state
    + returns a labeled image from marker-controlled watershed segmentation
    # segmentation
  image.hog_features
    fn (img: image_state, cell_size: i32, bins: i32) -> list[f64]
    + returns a histogram-of-oriented-gradients feature vector
    # features
    -> std.math.sqrt_f64
  image.resize_bilinear
    fn (img: image_state, width: i32, height: i32) -> image_state
    + resamples the image to the target shape with bilinear interpolation
    # geometry
