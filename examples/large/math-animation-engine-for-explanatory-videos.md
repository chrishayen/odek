# Requirement: "an animation engine for explanatory math videos"

Builds scenes composed of parametric shapes whose attributes interpolate over time, then renders each frame to pixels.

std
  std.math
    std.math.lerp_f32
      fn (a: f32, b: f32, t: f32) -> f32
      + returns a + (b - a) * t
      # math
    std.math.ease_in_out
      fn (t: f32) -> f32
      + returns a cubic ease-in-out of t in [0,1]
      # math
  std.image
    std.image.encode_png
      fn (pixels: bytes, width: i32, height: i32) -> bytes
      + encodes an RGBA pixel buffer as PNG
      # image

math_anim
  math_anim.new_scene
    fn (width: i32, height: i32, fps: i32) -> scene
    + creates an empty scene with the given canvas and frame rate
    # construction
  math_anim.add_line
    fn (s: scene, id: string, x0: f32, y0: f32, x1: f32, y1: f32, color: u32) -> scene
    + adds a line segment shape with an id that later animations can reference
    # shape
  math_anim.add_circle
    fn (s: scene, id: string, cx: f32, cy: f32, radius: f32, color: u32) -> scene
    + adds a circle shape with an id
    # shape
  math_anim.add_text
    fn (s: scene, id: string, x: f32, y: f32, content: string, color: u32) -> scene
    + adds a text shape anchored at (x,y)
    # shape
  math_anim.animate
    fn (s: scene, id: string, attribute: string, from: f32, to: f32, start_s: f32, duration_s: f32) -> result[scene, string]
    + schedules a tween on an attribute of a shape between two values over a time window
    - returns error when the shape id is unknown
    - returns error when the attribute does not exist on the shape
    # animation
  math_anim.sample_frame
    fn (s: scene, time_s: f32) -> scene
    + returns the scene with all tween attributes resolved for the given time
    # animation
    -> std.math.lerp_f32
    -> std.math.ease_in_out
  math_anim.render_frame
    fn (s: scene) -> bytes
    + rasterizes the current scene to an RGBA pixel buffer
    # rendering
  math_anim.render_png
    fn (s: scene, time_s: f32) -> bytes
    + samples at time_s, rasterizes, and encodes the frame as PNG
    # rendering
    -> std.image.encode_png
  math_anim.total_duration
    fn (s: scene) -> f32
    + returns the end time of the latest scheduled animation
    + returns 0.0 when no animations are scheduled
    # animation
