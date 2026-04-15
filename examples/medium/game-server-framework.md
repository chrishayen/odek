# Requirement: "a game server framework"

Manages connected players, groups them into rooms, and dispatches messages between them.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

game_server
  game_server.new
    fn () -> server_state
    + returns an empty server with no players or rooms
    # construction
  game_server.connect_player
    fn (state: server_state, player_id: string) -> result[server_state, string]
    + registers a new connected player
    - returns error when the player id is already connected
    # sessions
    -> std.time.now_millis
  game_server.disconnect_player
    fn (state: server_state, player_id: string) -> server_state
    + removes the player and any room memberships
    # sessions
  game_server.create_room
    fn (state: server_state, room_id: string, capacity: i32) -> result[server_state, string]
    + creates a new room with the given capacity
    - returns error on duplicate room id
    # rooms
  game_server.join_room
    fn (state: server_state, room_id: string, player_id: string) -> result[server_state, string]
    + adds the player to the room
    - returns error when the room is full or missing
    # rooms
  game_server.broadcast
    fn (state: server_state, room_id: string, sender: string, payload: bytes) -> result[list[string], string]
    + returns the list of player ids the message was delivered to
    - returns error when the room does not exist
    # messaging
  game_server.tick
    fn (state: server_state, elapsed_ms: i64) -> server_state
    + advances game state by the given elapsed time
    # simulation
