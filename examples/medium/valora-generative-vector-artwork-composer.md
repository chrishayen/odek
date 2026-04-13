# Requirement: "a library for composing generative vector artwork"

Builds a scene from primitive shapes and affine transforms, then rasterizes it into a pixel buffer. No file IO: the caller decides how to save the buffer.

std
  std.math
    std.math.sin
      @ (x: f64) -> f64
      + returns the sine of x in radians
      # math
    std.math.cos
      @ (x: f64) -> f64
      + returns the cosine of x in radians
      # math
  std.random
    std.random.new
      @ (seed: u64) -> rng_state
      + creates a deterministic RNG state from a seed
      # random
    std.random.uniform
      @ (rng: rng_state) -> tuple[f64, rng_state]
      + returns a uniform f64 in [0, 1) and the advanced RNG state
      # random

valora
  valora.new_canvas
    @ (width: i32, height: i32) -> canvas_state
    + creates a pixel buffer filled with transparent pixels
    # construction
  valora.fill_background
    @ (c: canvas_state, r: u8, g: u8, b: u8) -> canvas_state
    + fills every pixel with the given color
    # shading
  valora.draw_circle
    @ (c: canvas_state, cx: f64, cy: f64, radius: f64, r: u8, g: u8, b: u8) -> canvas_state
    + draws a filled circle, clipped to the canvas bounds
    # drawing
  valora.draw_polygon
    @ (c: canvas_state, points: list[tuple[f64, f64]], r: u8, g: u8, b: u8) -> canvas_state
    + draws a filled polygon defined by ordered vertices
    # drawing
  valora.rotate_points
    @ (points: list[tuple[f64, f64]], center: tuple[f64, f64], radians: f64) -> list[tuple[f64, f64]]
    + returns the points rotated around the center by the given angle
    # transform
    -> std.math.sin
    -> std.math.cos
  valora.jitter_color
    @ (rng: rng_state, r: u8, g: u8, b: u8, amount: f64) -> tuple[tuple[u8, u8, u8], rng_state]
    + perturbs each channel by a uniform offset scaled by amount
    # shading
    -> std.random.uniform
  valora.pixels
    @ (c: canvas_state) -> bytes
    + returns the canvas as a row-major RGBA byte buffer
    # export
