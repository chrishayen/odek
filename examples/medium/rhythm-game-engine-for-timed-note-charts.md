# Requirement: "a rhythm game engine for timed-note charts"

The library loads a note chart, tracks player input against the chart timeline, and scores hits. Audio playback and rendering are the caller's concern.

std
  std.text
    std.text.split_lines
      fn (input: string) -> list[string]
      + splits on newlines, dropping trailing empty lines
      # text
    std.text.parse_i32
      fn (s: string) -> optional[i32]
      + parses a decimal integer
      - returns empty on non-numeric input
      # parsing

rhythm
  rhythm.parse_chart
    fn (raw: string) -> result[chart, string]
    + parses a chart of "time_ms lane" lines into a note list
    - returns error on lines that cannot be parsed
    # parsing
    -> std.text.split_lines
    -> std.text.parse_i32
  rhythm.new_session
    fn (chart: chart) -> session_state
    + returns a fresh session at time 0 with score 0
    # construction
  rhythm.register_hit
    fn (state: session_state, time_ms: i32, lane: i32) -> session_state
    + awards points based on timing window and marks the note consumed
    - ignores hits outside any note's timing window
    ? "perfect" within +/- 30 ms, "good" within +/- 80 ms, "miss" otherwise
    # scoring
  rhythm.advance
    fn (state: session_state, time_ms: i32) -> session_state
    + marks unhit notes older than the miss window as missed
    # timing
  rhythm.score
    fn (state: session_state) -> i64
    + returns the current score
    # accessor
  rhythm.is_finished
    fn (state: session_state) -> bool
    + returns true when all notes have been resolved
    # lifecycle
