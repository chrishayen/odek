# Requirement: "a real-time strategy game engine with GPU rendering"

Keeps a tick-based simulation of units, orders, and resources, and renders the world each frame to a GPU surface.

std
  std.gpu
    std.gpu.create_surface
      fn (width: i32, height: i32) -> result[gpu_surface, string]
      + creates a rendering surface of the given size
      - returns error when no adapter is available
      # gpu
    std.gpu.upload_mesh
      fn (surface: gpu_surface, vertices: list[f32], indices: list[u32]) -> result[mesh_handle, string]
      + uploads a mesh and returns a handle
      - returns error when the buffers are empty
      # gpu
    std.gpu.draw_frame
      fn (surface: gpu_surface, draws: list[draw_call]) -> result[void, string]
      + submits a frame composed of draw calls
      - returns error on device loss
      # gpu
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

rts_engine
  rts_engine.new_world
    fn (map_width: i32, map_height: i32) -> world_state
    + creates an empty world of the given dimensions
    # construction
  rts_engine.spawn_unit
    fn (world: world_state, kind: string, owner: i32, x: f32, y: f32) -> tuple[world_state, i64]
    + returns a world with a new unit and its assigned id
    # spawning
  rts_engine.issue_order
    fn (world: world_state, unit_id: i64, order: unit_order) -> result[world_state, string]
    + returns a world with the order queued on the unit
    - returns error when the unit id is unknown
    # commands
  rts_engine.tick
    fn (world: world_state, dt_ms: i32) -> world_state
    + advances the simulation by dt_ms, moving units, applying orders, and updating resources
    ? dt_ms is expected to be a fixed simulation step
    # simulation
  rts_engine.find_path
    fn (world: world_state, from_x: i32, from_y: i32, to_x: i32, to_y: i32) -> result[list[tuple[i32, i32]], string]
    + returns an A* path across the walkable tiles of the map
    - returns error when no path exists
    # pathfinding
  rts_engine.resolve_combat
    fn (world: world_state) -> world_state
    + applies damage between units in range for this tick
    # simulation
  rts_engine.gather_resources
    fn (world: world_state) -> world_state
    + credits each player with resources produced by their gatherers this tick
    # simulation
  rts_engine.build_draw_list
    fn (world: world_state, camera: camera_state) -> list[draw_call]
    + returns the list of draw calls for units, terrain, and UI under the camera
    # rendering
  rts_engine.render
    fn (surface: gpu_surface, world: world_state, camera: camera_state) -> result[void, string]
    + submits a frame for the current world state
    - returns error on device loss
    # rendering
    -> std.gpu.draw_frame
  rts_engine.init_surface
    fn (width: i32, height: i32) -> result[gpu_surface, string]
    + creates a rendering surface and uploads default meshes
    - returns error when the GPU surface cannot be created
    # initialization
    -> std.gpu.create_surface
    -> std.gpu.upload_mesh
  rts_engine.run_tick
    fn (world: world_state) -> world_state
    + computes dt from elapsed real time and advances the world one step
    # loop
    -> std.time.now_millis
