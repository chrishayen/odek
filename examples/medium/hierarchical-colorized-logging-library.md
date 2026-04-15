# Requirement: "a hierarchical structured logging library with colorized output"

Logs form nested stories; each line inherits from its parent story and can be rendered with ANSI colors.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.ansi
    std.ansi.colorize
      fn (text: string, color: string) -> string
      + wraps text in the ANSI escape for the named color and a reset
      + returns text unchanged when color is empty
      # terminal
    std.ansi.strip
      fn (text: string) -> string
      + returns text with all ANSI escape sequences removed
      # terminal

story_log
  story_log.new_root
    fn (name: string) -> story_state
    + creates a top-level story with no parent
    # construction
    -> std.time.now_millis
  story_log.begin_child
    fn (parent: story_state, name: string) -> tuple[story_state, story_id]
    + returns an updated parent and the id of the new child story
    ? child ids are formed by appending to the parent's path
    # hierarchy
    -> std.time.now_millis
  story_log.end
    fn (state: story_state, id: story_id) -> story_state
    + marks the story as finished and records its duration
    - no-op when id is unknown
    # hierarchy
    -> std.time.now_millis
  story_log.record
    fn (state: story_state, id: story_id, level: string, message: string) -> story_state
    + appends a log entry to the story with timestamp and level
    # logging
    -> std.time.now_millis
  story_log.render_line
    fn (state: story_state, id: story_id, entry_index: i32, use_color: bool) -> result[string, string]
    + returns a rendered line with indentation, timestamp, level, and message
    - returns error when id or entry_index is out of range
    # rendering
    -> std.ansi.colorize
  story_log.flatten
    fn (state: story_state, use_color: bool) -> list[string]
    + returns every entry across every story in the order it was recorded
    # rendering
    -> std.ansi.colorize
