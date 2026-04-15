# Requirement: "a terminal music player core"

Queue, playback state, and UI snapshot. Audio decode and rendering are the caller's problem; this library manages state.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

player
  player.new
    fn () -> player_state
    + creates a player with an empty queue, stopped, at position 0
    # construction
  player.enqueue
    fn (state: player_state, track: track_info) -> void
    + appends a track to the end of the play queue
    # queue
  player.enqueue_next
    fn (state: player_state, track: track_info) -> void
    + inserts a track immediately after the currently playing item
    # queue
  player.remove_at
    fn (state: player_state, index: i32) -> result[void, string]
    + removes the track at the given queue index
    - returns error when the index is out of range
    # queue
  player.play
    fn (state: player_state) -> result[void, string]
    + starts or resumes playback of the current track
    - returns error when the queue is empty
    # playback_control
    -> std.time.now_millis
  player.pause
    fn (state: player_state) -> void
    + pauses playback and freezes the current position
    # playback_control
    -> std.time.now_millis
  player.next_track
    fn (state: player_state) -> optional[track_info]
    + advances to the next track, respecting the current repeat mode
    + returns none when the queue has been exhausted and repeat is off
    # playback_control
  player.prev_track
    fn (state: player_state) -> optional[track_info]
    + moves to the previous track if one exists
    # playback_control
  player.seek
    fn (state: player_state, position_millis: i64) -> result[void, string]
    + sets playback position within the current track
    - returns error when position is negative or past the track length
    # playback_control
  player.set_repeat_mode
    fn (state: player_state, mode: string) -> result[void, string]
    + sets repeat mode to "off", "one", or "all"
    - returns error on unknown mode
    # mode
  player.toggle_shuffle
    fn (state: player_state) -> bool
    + toggles shuffle and returns the new state
    # mode
  player.snapshot
    fn (state: player_state) -> player_snapshot
    + returns current track, position, queue contents, repeat mode, and shuffle flag
    # inspection
    -> std.time.now_millis
