# Requirement: "an image processing filter library"

Pure functions over an RGBA pixel buffer. Each filter returns a new image; the input is never mutated.

std: (all units exist)

image_filters
  image_filters.new_image
    fn (width: i32, height: i32, pixels: bytes) -> result[image, string]
    + creates an image from a raw RGBA byte buffer
    - returns error when buffer size does not equal width*height*4
    # construction
  image_filters.grayscale
    fn (img: image) -> image
    + converts to luminance using the standard Rec. 601 weights
    # color
  image_filters.invert
    fn (img: image) -> image
    + replaces each channel with 255 minus itself, preserving alpha
    # color
  image_filters.brightness
    fn (img: image, delta: i32) -> image
    + adds delta to each RGB channel, clamping to [0, 255]
    # adjustment
  image_filters.contrast
    fn (img: image, factor: f32) -> image
    + scales each channel around the midpoint 128 by factor
    # adjustment
  image_filters.gaussian_blur
    fn (img: image, radius: i32) -> image
    + separable Gaussian blur with the given radius
    + returns a copy of the input when radius is 0
    # convolution
  image_filters.sharpen
    fn (img: image, amount: f32) -> image
    + unsharp mask style sharpening with the given strength
    # convolution
  image_filters.resize
    fn (img: image, new_width: i32, new_height: i32) -> result[image, string]
    + bilinear resize to the target dimensions
    - returns error when either dimension is <= 0
    # geometry
  image_filters.crop
    fn (img: image, x: i32, y: i32, w: i32, h: i32) -> result[image, string]
    + returns the sub-image at the given rectangle
    - returns error when the rectangle falls outside the source
    # geometry
