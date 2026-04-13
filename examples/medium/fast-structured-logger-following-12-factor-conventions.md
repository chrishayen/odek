# Requirement: "a fast structured application logger following 12-factor conventions"

The logger writes structured events to a stream sink. Level filtering and formatting are split so a caller can swap either independently.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.io
    std.io.write_line
      @ (line: string) -> result[void, string]
      + writes a line followed by a newline to the process output stream
      - returns error when the stream is closed
      # io
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      + produces keys in sorted order for deterministic output
      # serialization

logger
  logger.new
    @ (min_level: i32) -> logger_state
    + creates a logger filtering events below the given level
    ? levels follow a fixed order: debug=0, info=1, warn=2, error=3
    # construction
  logger.with_field
    @ (state: logger_state, key: string, value: string) -> logger_state
    + returns a new state carrying an additional context field
    ? fields are merged into every event emitted by the returned state
    # context
  logger.log
    @ (state: logger_state, level: i32, message: string) -> result[void, string]
    + formats the event as a single JSON object with timestamp, level, message and context fields
    + writes nothing when level is below the state's minimum
    - returns error when the underlying stream write fails
    # emit
    -> std.time.now_millis
    -> std.json.encode_object
    -> std.io.write_line
  logger.level_from_name
    @ (name: string) -> result[i32, string]
    + maps "debug", "info", "warn", "error" to their numeric level
    - returns error on any other input
    # parsing
