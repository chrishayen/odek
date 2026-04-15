# Requirement: "a physically based renderer for glTF 2.0 scenes using GPU ray tracing"

Loads a glTF scene, builds acceleration structures, and traces rays with a physically based shading model. The project exposes scene load, camera setup, and a render call.

std
  std.fs
    std.fs.read_bytes
      fn (path: string) -> result[bytes, string]
      + returns the full file contents
      - returns error when the file does not exist
      # filesystem
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a generic value tree
      - returns error on malformed input
      # serialization
  std.math
    std.math.vec3_normalize
      fn (v: tuple[f64, f64, f64]) -> tuple[f64, f64, f64]
      + returns v scaled to unit length
      ? returns (0,0,0) when v has zero length
      # math
    std.math.vec3_dot
      fn (a: tuple[f64, f64, f64], b: tuple[f64, f64, f64]) -> f64
      + returns the dot product
      # math
    std.math.vec3_cross
      fn (a: tuple[f64, f64, f64], b: tuple[f64, f64, f64]) -> tuple[f64, f64, f64]
      + returns the cross product
      # math

pbr_rt
  pbr_rt.load_gltf
    fn (path: string) -> result[scene_state, string]
    + parses a glTF 2.0 file and its referenced buffers into a scene
    - returns error when the root is not a valid glTF document
    - returns error when a required buffer is missing
    # loading
    -> std.fs.read_bytes
    -> std.json.parse
  pbr_rt.build_bvh
    fn (scene: scene_state) -> scene_state
    + constructs a bounding volume hierarchy over all triangle meshes in the scene
    ? empty scenes produce a trivial BVH
    # acceleration
  pbr_rt.make_camera
    fn (position: tuple[f64, f64, f64], target: tuple[f64, f64, f64], fov_deg: f64, aspect: f64) -> camera
    + creates a pinhole camera looking from position at target
    # camera
    -> std.math.vec3_normalize
    -> std.math.vec3_cross
  pbr_rt.generate_primary_ray
    fn (cam: camera, u: f64, v: f64) -> ray
    + returns a primary ray for the normalized film coordinates (u, v)
    # ray_generation
  pbr_rt.intersect
    fn (scene: scene_state, r: ray) -> optional[hit]
    + returns the closest hit of r with the scene, or none
    # intersection
  pbr_rt.sample_bsdf
    fn (surface: hit, wi: tuple[f64, f64, f64], rand_u: f64, rand_v: f64) -> tuple[tuple[f64, f64, f64], f64]
    + returns a sampled outgoing direction and its pdf for the surface BSDF
    ? metallic and roughness come from the material's PBR metallic-roughness workflow
    # shading
    -> std.math.vec3_dot
  pbr_rt.evaluate_bsdf
    fn (surface: hit, wi: tuple[f64, f64, f64], wo: tuple[f64, f64, f64]) -> tuple[f64, f64, f64]
    + returns the RGB BSDF value for the given incoming and outgoing directions
    # shading
  pbr_rt.trace_path
    fn (scene: scene_state, r: ray, max_bounces: i32) -> tuple[f64, f64, f64]
    + returns the radiance along the ray by recursively sampling the BSDF up to max_bounces
    # path_tracing
  pbr_rt.render
    fn (scene: scene_state, cam: camera, width: i32, height: i32, samples: i32) -> list[tuple[f64, f64, f64]]
    + returns a flat row-major framebuffer of linear RGB radiance values
    ? dimensions must be positive; samples must be at least 1
    # rendering
  pbr_rt.tonemap
    fn (pixels: list[tuple[f64, f64, f64]]) -> list[tuple[u8, u8, u8]]
    + applies an ACES-like tone map and returns 8-bit sRGB pixels
    # output
