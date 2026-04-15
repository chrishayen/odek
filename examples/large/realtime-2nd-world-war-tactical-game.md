# Requirement: "a realtime 2nd world war tactical game"

A tactical game engine for historical ground combat: units, orders, line-of-sight, and a tick-based simulation loop. The project package hosts game logic; std provides math and time primitives.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.math
    std.math.distance2d
      fn (ax: f64, ay: f64, bx: f64, by: f64) -> f64
      + returns euclidean distance between two points
      # geometry
    std.math.angle_between
      fn (ax: f64, ay: f64, bx: f64, by: f64) -> f64
      + returns angle in radians from a to b
      # geometry
  std.rand
    std.rand.new_seeded
      fn (seed: u64) -> rng_state
      + returns a deterministic rng seeded with the given value
      # randomness
    std.rand.next_f64
      fn (rng: rng_state) -> tuple[f64, rng_state]
      + returns a uniform random number in [0,1) and the advanced rng
      # randomness

tactical_game
  tactical_game.new_world
    fn (map_width: i32, map_height: i32, seed: u64) -> world_state
    + creates an empty world with the given map dimensions and rng seed
    # construction
    -> std.rand.new_seeded
  tactical_game.spawn_unit
    fn (world: world_state, side: i32, kind: string, x: f64, y: f64) -> tuple[unit_id, world_state]
    + places a unit of the given side and kind at the given tile
    + each unit has hp, morale, ammo, and a facing angle
    - returns invalid unit_id when coordinates are outside the map
    # unit_management
  tactical_game.issue_move_order
    fn (world: world_state, unit: unit_id, tx: f64, ty: f64) -> result[world_state, string]
    + queues a move order for the unit toward the target tile
    - returns error when the unit is already dead
    # commands
  tactical_game.issue_fire_order
    fn (world: world_state, unit: unit_id, target: unit_id) -> result[world_state, string]
    + queues a fire order against the target unit
    - returns error when the target is on the same side
    - returns error when the target is out of range
    # commands
  tactical_game.line_of_sight
    fn (world: world_state, from: unit_id, to: unit_id) -> bool
    + returns true when no terrain obstacle blocks sight between the two units
    # perception
    -> std.math.distance2d
  tactical_game.tick
    fn (world: world_state, dt_millis: i64) -> world_state
    + advances movement, fire resolution, morale, and order completion by dt
    + fire resolution uses rng and distance to compute hit probability
    ? damage is applied immediately; no shell travel time
    # simulation
    -> std.math.distance2d
    -> std.math.angle_between
    -> std.rand.next_f64
  tactical_game.unit_view
    fn (world: world_state, unit: unit_id) -> optional[unit_snapshot]
    + returns a read-only snapshot of the unit (hp, pos, facing, morale)
    - returns none when the unit does not exist
    # query
  tactical_game.side_won
    fn (world: world_state) -> optional[i32]
    + returns the winning side when only one side has living units
    - returns none while both sides still have living units
    # victory
  tactical_game.set_terrain
    fn (world: world_state, x: i32, y: i32, kind: string) -> world_state
    + sets the terrain kind at the given tile (open, forest, building, water)
    # terrain
  tactical_game.run_for
    fn (world: world_state, total_millis: i64, step_millis: i64) -> world_state
    + advances the simulation by total_millis in step_millis increments
    # simulation
    -> std.time.now_millis
