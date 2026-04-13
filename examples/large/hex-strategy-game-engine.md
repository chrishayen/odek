# Requirement: "a turn-based hexagonal strategy game engine"

Pure state transitions over a hex grid. Rendering and input are the caller's concern; this library owns rules, movement, and combat.

std: (all units exist)

hex_strategy
  hex_strategy.new_game
    @ (width: i32, height: i32, player_count: i32) -> game_state
    + creates a fresh game with empty tiles and the given player count
    ? uses axial coordinates internally
    # construction
  hex_strategy.place_unit
    @ (state: game_state, player: i32, kind: string, q: i32, r: i32) -> result[game_state, string]
    + places a unit of the given kind for player at the hex
    - returns error when the hex is off the map
    - returns error when the hex is already occupied
    # setup
  hex_strategy.neighbors
    @ (q: i32, r: i32) -> list[tuple[i32, i32]]
    + returns the six adjacent hex coordinates
    # grid
  hex_strategy.distance
    @ (a_q: i32, a_r: i32, b_q: i32, b_r: i32) -> i32
    + returns the hex-grid distance between two coordinates
    # grid
  hex_strategy.reachable
    @ (state: game_state, unit_id: i32) -> list[tuple[i32, i32]]
    + returns hexes the unit can move to this turn given remaining movement
    # pathfinding
  hex_strategy.move_unit
    @ (state: game_state, unit_id: i32, to_q: i32, to_r: i32) -> result[game_state, string]
    + moves the unit, consuming movement points
    - returns error when the destination is unreachable
    - returns error when the unit belongs to another player
    # movement
  hex_strategy.attack
    @ (state: game_state, attacker_id: i32, defender_id: i32) -> result[game_state, string]
    + resolves an attack, subtracting damage from the defender
    - returns error when defender is not adjacent
    - returns error when attacker has already attacked this turn
    # combat
  hex_strategy.end_turn
    @ (state: game_state) -> game_state
    + advances to the next player and refreshes unit movement
    # turn_flow
  hex_strategy.current_player
    @ (state: game_state) -> i32
    + returns the index of the player whose turn it is
    # turn_flow
  hex_strategy.winner
    @ (state: game_state) -> optional[i32]
    + returns the surviving player when only one remains
    # victory
  hex_strategy.units_at
    @ (state: game_state, q: i32, r: i32) -> list[unit]
    + returns the units on the given hex
    # query
