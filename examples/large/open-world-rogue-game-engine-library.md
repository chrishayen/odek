# Requirement: "an open world roguelike game engine library"

Tile-based world, player and entity state, turn loop, procedural map generation, and item handling. The project layer is broad; std supplies generic primitives.

std
  std.random
    std.random.new_rng
      fn (seed: u64) -> rng_state
      + creates a deterministic pseudorandom generator from a seed
      # random
    std.random.next_int
      fn (rng: rng_state, max_exclusive: i32) -> tuple[i32, rng_state]
      + returns a uniform integer in [0, max_exclusive) and the advanced rng
      # random
  std.grid
    std.grid.new
      fn (width: i32, height: i32, fill: i32) -> grid_state
      + creates a rectangular grid populated with the fill value
      # grid
    std.grid.get
      fn (grid: grid_state, x: i32, y: i32) -> optional[i32]
      + returns the tile at (x,y), or None when out of bounds
      # grid
    std.grid.set
      fn (grid: grid_state, x: i32, y: i32, value: i32) -> grid_state
      + returns a new grid with the tile at (x,y) replaced
      # grid

rogue
  rogue.new_world
    fn (seed: u64, width: i32, height: i32) -> world_state
    + creates a world with a procedurally generated map from the given seed
    # construction
    -> std.random.new_rng
    -> std.grid.new
  rogue.generate_map
    fn (world: world_state) -> world_state
    + carves rooms and corridors into the world's grid using its rng
    # world_generation
    -> std.random.next_int
    -> std.grid.set
  rogue.spawn_player
    fn (world: world_state, x: i32, y: i32) -> result[world_state, string]
    + places the player at (x,y) with default stats
    - returns error when the tile is not walkable
    # player
    -> std.grid.get
  rogue.spawn_entity
    fn (world: world_state, kind: string, x: i32, y: i32) -> result[tuple[string, world_state], string]
    + adds an entity of the given kind and returns its id
    - returns error when the tile is occupied or out of bounds
    # entities
    -> std.grid.get
  rogue.move_player
    fn (world: world_state, dx: i32, dy: i32) -> result[world_state, string]
    + moves the player by (dx, dy) and consumes one turn
    - returns error when the destination is blocked
    # movement
    -> std.grid.get
  rogue.attack
    fn (world: world_state, target_id: string) -> result[world_state, string]
    + resolves a melee attack against target_id, applying damage via the rng
    - returns error when target_id is unknown or out of range
    # combat
    -> std.random.next_int
  rogue.pickup_item
    fn (world: world_state, item_id: string) -> result[world_state, string]
    + moves the item from the ground into the player's inventory
    - returns error when the item is not on the player's tile
    # inventory
  rogue.use_item
    fn (world: world_state, item_id: string) -> result[world_state, string]
    + applies the item's effect and removes it from inventory if consumable
    - returns error when the item is not in the player's inventory
    # inventory
  rogue.tick
    fn (world: world_state) -> world_state
    + advances all non-player entities by one turn using their simple ai
    # turn_loop
    -> std.random.next_int
  rogue.visible_tiles
    fn (world: world_state, radius: i32) -> list[tuple[i32, i32]]
    + returns the tiles visible to the player within radius using line-of-sight
    # field_of_view
    -> std.grid.get
  rogue.register_tile_kind
    fn (world: world_state, kind: string, walkable: bool, glyph: string) -> world_state
    + defines a new tile kind so maps can include custom terrain
    # extensibility
