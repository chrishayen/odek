# Requirement: "a game server engine with entities, spatial rooms, and hot-reloadable behavior scripts"

Entities live in spaces (rooms), interact with nearby entities, and run behaviors that can be swapped at runtime without dropping state.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.rand
    std.rand.u64
      fn () -> u64
      + returns a pseudo-random unsigned 64-bit integer
      # random
  std.net
    std.net.send_message
      fn (peer: string, payload: bytes) -> result[void, string]
      + sends a message to a remote peer
      - returns error on connection failure
      # network

game_server
  game_server.entity_new
    fn (kind: string, x: f32, y: f32) -> entity
    + creates an entity with a fresh id and a position
    # construction
    -> std.rand.u64
  game_server.space_new
    fn (name: string, width: f32, height: f32) -> space
    + creates an empty space with the given bounds
    # construction
  game_server.world_new
    fn () -> world_state
    + creates a world containing no spaces and no behaviors
    # construction
  game_server.space_add
    fn (world: world_state, s: space) -> world_state
    + registers a space in the world
    # space_mgmt
  game_server.entity_spawn
    fn (world: world_state, space_name: string, e: entity) -> result[world_state, string]
    + places an entity into the named space
    - returns error when the space does not exist
    # space_mgmt
  game_server.entity_move
    fn (world: world_state, id: entity_id, dx: f32, dy: f32) -> result[world_state, string]
    + translates an entity's position, clamping to space bounds
    - returns error when the entity does not exist
    # movement
  game_server.entity_migrate
    fn (world: world_state, id: entity_id, target_space: string) -> result[world_state, string]
    + moves an entity from its current space to another, preserving state
    - returns error when the target space does not exist
    # space_mgmt
  game_server.aoi_query
    fn (world: world_state, space_name: string, x: f32, y: f32, radius: f32) -> list[entity_id]
    + returns entity ids within a radius of a point (area-of-interest query)
    # queries
  game_server.behavior_register
    fn (world: world_state, name: string, fn: behavior_fn) -> world_state
    + registers a named behavior function
    # behaviors
  game_server.behavior_swap
    fn (world: world_state, name: string, fn: behavior_fn) -> result[world_state, string]
    + replaces a registered behavior without dropping entity state
    - returns error when no behavior with that name is registered
    # hot_swap
  game_server.entity_bind_behavior
    fn (world: world_state, id: entity_id, behavior_name: string) -> result[world_state, string]
    + attaches a registered behavior to an entity
    - returns error when the behavior is unknown
    # behaviors
  game_server.tick
    fn (world: world_state, dt_ms: i64) -> world_state
    + advances one simulation step, running each entity's behavior once
    # simulation
    -> std.time.now_millis
  game_server.broadcast
    fn (world: world_state, space_name: string, payload: bytes) -> result[i32, string]
    + sends payload to every peer currently subscribed to the space and returns the count delivered
    - returns error when the space does not exist
    # networking
    -> std.net.send_message
  game_server.snapshot_space
    fn (world: world_state, space_name: string) -> result[list[entity], string]
    + returns a copy of every entity in the space
    - returns error when the space does not exist
    # queries
