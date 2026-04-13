# Requirement: "a fast image math kernel library"

Pixel-level operations (blend, convolve, gamma) on planar RGBA buffers with pluggable kernel backends.

std
  std.math
    std.math.pow
      @ (base: f64, exp: f64) -> f64
      + returns base raised to exp
      # math
    std.math.clamp_f64
      @ (x: f64, lo: f64, hi: f64) -> f64
      + returns x clamped to the closed range [lo, hi]
      # math

imath
  imath.new_image
    @ (width: i32, height: i32) -> image
    + allocates an RGBA image with the given dimensions, zero-filled
    - returns an empty image when width or height is not positive
    # allocation
  imath.blend
    @ (dst: image, src: image, op: blend_op) -> image
    + blends src over dst using the named porter-duff operator
    + assumes premultiplied alpha
    # compositing
  imath.convolve
    @ (img: image, kernel: list[f32], kernel_w: i32) -> image
    + applies a separable or square kernel and clamps at borders
    - returns the input unchanged when kernel_w is not a positive odd integer
    # filtering
    -> std.math.clamp_f64
  imath.gamma
    @ (img: image, gamma: f64) -> image
    + applies a per-channel gamma curve to RGB, preserving alpha
    # color
    -> std.math.pow
  imath.resize
    @ (img: image, new_w: i32, new_h: i32) -> image
    + resamples to a new size using bilinear interpolation
    # resampling
