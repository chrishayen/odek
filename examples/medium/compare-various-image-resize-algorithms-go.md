# Requirement: "an image resize library with multiple algorithms"

Offer nearest-neighbor, bilinear, and bicubic resampling on a simple raster image type. The library owns the sampling kernels; pixel buffers are plain bytes supplied by the caller.

std: (all units exist)

image_resize
  image_resize.new_image
    @ (width: i32, height: i32, pixels: bytes) -> result[raster_image, string]
    + wraps pixels as an RGBA raster image
    - returns error when pixels length does not equal width*height*4
    # construction
  image_resize.nearest
    @ (src: raster_image, new_width: i32, new_height: i32) -> result[raster_image, string]
    + resamples src to (new_width, new_height) using nearest-neighbor
    - returns error when new_width or new_height is not positive
    # resampling
  image_resize.bilinear
    @ (src: raster_image, new_width: i32, new_height: i32) -> result[raster_image, string]
    + resamples src using bilinear interpolation
    + averages the four nearest source pixels weighted by subpixel distance
    - returns error when new_width or new_height is not positive
    # resampling
  image_resize.bicubic
    @ (src: raster_image, new_width: i32, new_height: i32) -> result[raster_image, string]
    + resamples src using a Catmull-Rom bicubic kernel
    - returns error when new_width or new_height is not positive
    # resampling
  image_resize.compare_ssim
    @ (a: raster_image, b: raster_image) -> result[f64, string]
    + returns the structural similarity score between a and b in [-1.0, 1.0]
    - returns error when a and b differ in dimensions
    # quality
