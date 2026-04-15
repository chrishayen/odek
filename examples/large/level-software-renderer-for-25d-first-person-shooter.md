# Requirement: "a software renderer for a classic 2.5D first-person shooter level format"

Loads a level archive, walks the BSP to determine visible sectors from a viewpoint, and rasterizes walls, flats, and sprites into a framebuffer.

std
  std.fs
    std.fs.read_all_bytes
      fn (path: string) -> result[bytes, string]
      + reads the entire file as bytes
      - returns error when the file cannot be read
      # filesystem
  std.math
    std.math.sin
      fn (radians: f64) -> f64
      + returns the sine of the angle
      # math
    std.math.cos
      fn (radians: f64) -> f64
      + returns the cosine of the angle
      # math
    std.math.atan2
      fn (y: f64, x: f64) -> f64
      + returns the angle of the vector (x, y) in radians
      # math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root
      # math

level
  level.load_archive
    fn (path: string) -> result[level_archive, string]
    + parses a packed level archive into its named lumps
    - returns error when the magic header is missing
    # loading
    -> std.fs.read_all_bytes
  level.load_map
    fn (archive: level_archive, name: string) -> result[map_data, string]
    + assembles vertices, lines, sides, sectors, and the BSP tree for the named map
    - returns error when any required lump is missing
    # loading

renderer
  renderer.new
    fn (width: i32, height: i32, fov_degrees: f64) -> renderer_state
    + creates a renderer targeting a width by height framebuffer with the given horizontal FOV
    # construction
  renderer.set_viewpoint
    fn (state: renderer_state, x: f64, y: f64, z: f64, angle: f64) -> renderer_state
    + positions the camera at (x, y, z) facing angle radians
    # camera
  renderer.walk_bsp
    fn (state: renderer_state, map: map_data) -> list[sub_sector_ref]
    + returns the visible sub-sectors in front-to-back order from the current viewpoint
    # visibility
  renderer.clip_wall
    fn (state: renderer_state, wall: wall_segment) -> optional[clipped_wall]
    + returns the wall clipped to the view frustum, or none when fully off-screen
    # clipping
    -> std.math.sin
    -> std.math.cos
  renderer.project_column
    fn (state: renderer_state, wall: clipped_wall, screen_x: i32) -> column_span
    + returns the top and bottom screen rows for one wall column
    # projection
    -> std.math.atan2
    -> std.math.sqrt
  renderer.draw_wall_column
    fn (fb: framebuffer, span: column_span, texture: texture, light: f64) -> framebuffer
    + writes the sampled texture column into the framebuffer between top and bottom
    # rasterization
  renderer.draw_flat
    fn (fb: framebuffer, state: renderer_state, sub: sub_sector_ref, texture: texture) -> framebuffer
    + fills floor and ceiling pixels for the sub-sector
    # rasterization
  renderer.draw_sprite
    fn (fb: framebuffer, state: renderer_state, sprite: sprite_ref) -> framebuffer
    + draws the sprite with depth clipping against previously drawn walls
    # rasterization
  renderer.render_frame
    fn (state: renderer_state, map: map_data, fb: framebuffer) -> framebuffer
    + renders one full frame into the framebuffer
    # orchestration
