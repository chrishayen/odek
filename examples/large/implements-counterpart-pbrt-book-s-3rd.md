# Requirement: "a physically based ray tracing renderer"

A renderer library: scene description in, image bytes out. Real substance lives in geometry, shading, and integration primitives.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root
      - returns NaN for negative input
      # math
    std.math.clamp
      @ (x: f64, lo: f64, hi: f64) -> f64
      + returns x clamped to [lo, hi]
      # math
  std.random
    std.random.uniform_f64
      @ (seed_state: rng_state) -> tuple[f64, rng_state]
      + returns a uniform sample in [0, 1) and the advanced rng state
      # random

renderer
  renderer.make_scene
    @ () -> scene
    + creates an empty scene with no shapes and no lights
    # construction
  renderer.add_sphere
    @ (s: scene, center: vec3, radius: f64, material_id: i32) -> scene
    + appends a sphere primitive to the scene
    - radius must be positive
    # geometry
  renderer.add_triangle_mesh
    @ (s: scene, vertices: list[vec3], indices: list[i32], material_id: i32) -> result[scene, string]
    + appends a triangle mesh
    - returns error when indices length is not a multiple of three
    # geometry
  renderer.add_point_light
    @ (s: scene, position: vec3, intensity: vec3) -> scene
    + appends a point light at the given position with rgb intensity
    # lighting
  renderer.define_lambert_material
    @ (s: scene, albedo: vec3) -> tuple[scene, i32]
    + registers a diffuse material and returns its id
    # materials
  renderer.define_mirror_material
    @ (s: scene, reflectance: vec3) -> tuple[scene, i32]
    + registers a perfectly specular material and returns its id
    # materials
  renderer.intersect_ray
    @ (s: scene, origin: vec3, direction: vec3) -> optional[hit_record]
    + returns the nearest intersection along the ray
    - returns none when the ray misses everything
    # intersection
    -> std.math.sqrt
  renderer.shade_hit
    @ (s: scene, hit: hit_record, view_dir: vec3) -> vec3
    + returns outgoing radiance at the hit point under all lights
    # shading
    -> std.math.clamp
  renderer.trace_path
    @ (s: scene, origin: vec3, direction: vec3, depth: i32, rng: rng_state) -> tuple[vec3, rng_state]
    + returns estimated radiance along the ray via recursive path tracing
    ? terminates when depth reaches zero or by russian roulette
    # integration
    -> std.random.uniform_f64
  renderer.render_tile
    @ (s: scene, cam: camera, x0: i32, y0: i32, x1: i32, y1: i32, samples: i32) -> list[vec3]
    + returns radiance per pixel for the tile in row-major order
    # rendering
  renderer.tonemap_to_srgb
    @ (pixels: list[vec3]) -> bytes
    + applies gamma correction and returns 8-bit rgb bytes
    # output
    -> std.math.clamp
