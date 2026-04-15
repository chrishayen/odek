# Requirement: "a turn-based hexagonal strategy game engine"

A hex grid, units with movement and combat, turn order, and a damage model. Randomness is behind a std primitive so battles are reproducible in tests.

std
  std.random
    std.random.next_unit_f64
      fn () -> f64
      + returns a pseudo-random number in [0.0, 1.0)
      # randomness

hexwar
  hexwar.new_board
    fn (radius: i32) -> result[board_state, string]
    + creates an empty hex board with the given radius using axial coordinates
    - returns error when radius is negative
    # construction
  hexwar.hex_distance
    fn (a: hex_coord, b: hex_coord) -> i32
    + returns the hex distance using cube coordinates
    + returns 0 when a equals b
    # geometry
  hexwar.neighbors
    fn (c: hex_coord) -> list[hex_coord]
    + returns the six adjacent hex coordinates
    # geometry
  hexwar.place_unit
    fn (state: board_state, unit: unit_spec, at: hex_coord) -> result[board_state, string]
    + places the unit at the hex
    - returns error when the hex is off-board
    - returns error when the hex is already occupied
    # setup
  hexwar.reachable_tiles
    fn (state: board_state, unit_id: string) -> result[list[hex_coord], string]
    + returns every hex the unit can reach within its movement range
    - returns error when unit_id is unknown
    # movement
    -> hexwar.hex_distance
    -> hexwar.neighbors
  hexwar.move_unit
    fn (state: board_state, unit_id: string, to: hex_coord) -> result[board_state, string]
    + moves the unit and subtracts the path cost from its action points
    - returns error when the destination is out of range
    - returns error when the destination is occupied
    # movement
  hexwar.attack
    fn (state: board_state, attacker_id: string, defender_id: string) -> result[battle_outcome, string]
    + resolves damage using attacker strength, defender armor, and one random roll
    - returns error when the defender is not adjacent
    - returns error when the attacker is out of action points
    # combat
    -> std.random.next_unit_f64
  hexwar.end_turn
    fn (state: board_state) -> board_state
    + advances to the next player and refreshes action points for that player's units
    # turn_order
  hexwar.current_player
    fn (state: board_state) -> i32
    + returns the id of the player whose turn it is
    # turn_order
  hexwar.winner
    fn (state: board_state) -> optional[i32]
    + returns the player id when exactly one player has surviving units
    - returns none while multiple players still have units
    # status
  hexwar.unit_at
    fn (state: board_state, c: hex_coord) -> optional[unit_view]
    + returns the unit occupying the given hex
    - returns none for an empty hex
    # inspection
