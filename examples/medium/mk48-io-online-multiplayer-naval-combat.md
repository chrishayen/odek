# Requirement: "a multiplayer naval combat game engine"

Server-side game state for ships on a 2D plane: movement, projectiles, damage, and scoring.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.math
    std.math.hypot
      @ (dx: f64, dy: f64) -> f64
      + returns sqrt(dx*dx + dy*dy)
      # math

naval
  naval.new_world
    @ (width: f64, height: f64) -> world_state
    + creates an empty world with the given bounds
    # construction
  naval.spawn_ship
    @ (world: world_state, player_id: string, x: f64, y: f64) -> tuple[u64, world_state]
    + adds a ship at (x, y) owned by player_id and returns its id
    # spawn
  naval.set_heading
    @ (world: world_state, ship_id: u64, heading_radians: f64, throttle: f64) -> result[world_state, string]
    + updates the ship's heading and normalized throttle in [0, 1]
    - returns error when ship_id does not exist
    # control
  naval.fire
    @ (world: world_state, ship_id: u64) -> result[tuple[u64, world_state], string]
    + spawns a projectile ahead of the ship and returns its id
    - returns error when the ship is on cooldown
    # combat
  naval.tick
    @ (world: world_state, dt_seconds: f64) -> world_state
    + advances ships and projectiles, applies collisions and damage, removes sunk ships and expired projectiles
    + clamps positions to world bounds
    # simulation
    -> std.math.hypot
  naval.step_to_now
    @ (world: world_state) -> world_state
    + advances the world using wall-clock delta since its last tick
    # scheduling
    -> std.time.now_millis
  naval.scoreboard
    @ (world: world_state) -> list[tuple[string, i32]]
    + returns (player_id, score) pairs sorted by score descending
    # scoring
