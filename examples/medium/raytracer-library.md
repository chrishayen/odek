# Requirement: "a simple raytracer"

Classic ray-sphere-plane raytracer. Std is a tiny vector-math module; the project owns scene, intersection, and image assembly.

std
  std.math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the non-negative square root
      + returns 0 for x <= 0
      # math
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # filesystem

raytracer
  raytracer.vec3
    fn (x: f64, y: f64, z: f64) -> vec3_value
    + returns a 3-component vector
    # math
  raytracer.vec3_dot
    fn (a: vec3_value, b: vec3_value) -> f64
    + returns the dot product
    # math
  raytracer.vec3_normalize
    fn (v: vec3_value) -> vec3_value
    + returns the unit vector in the same direction
    + returns the zero vector for a zero-length input
    # math
    -> std.math.sqrt
  raytracer.scene_new
    fn () -> scene_state
    + returns an empty scene
    # construction
  raytracer.scene_add_sphere
    fn (s: scene_state, center: vec3_value, radius: f64, color: vec3_value) -> scene_state
    + returns a new scene with a sphere added
    # scene
  raytracer.scene_add_light
    fn (s: scene_state, position: vec3_value, intensity: f64) -> scene_state
    + returns a new scene with a point light added
    # scene
  raytracer.trace_ray
    fn (s: scene_state, origin: vec3_value, direction: vec3_value, depth: i32) -> vec3_value
    + returns the color sampled along the ray up to the given reflection depth
    ? background returned when no intersection
    # shading
    -> raytracer.vec3_dot
    -> raytracer.vec3_normalize
  raytracer.render
    fn (s: scene_state, width: i32, height: i32) -> image_buffer
    + returns an RGB image of the scene at the given resolution using a perspective camera
    # rendering
    -> raytracer.trace_ray
  raytracer.save_ppm
    fn (img: image_buffer, path: string) -> result[void, string]
    + writes the image to disk as a P6 PPM file
    # output
    -> std.fs.write_all
