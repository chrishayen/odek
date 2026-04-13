# Requirement: "a typing test game engine"

A library for a single-player typing test. It tracks typed characters against a target phrase and reports accuracy and words-per-minute. Rendering is the caller's job.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

typing_test
  typing_test.start
    @ (target: string, start_millis: i64) -> result[game_state, string]
    + creates a game with the target phrase and a start time
    - returns error when target is empty
    # construction
  typing_test.type_char
    @ (state: game_state, ch: string) -> game_state
    + advances the cursor when the typed character matches the next target character
    + records a mistake when the typed character does not match
    # input
  typing_test.backspace
    @ (state: game_state) -> game_state
    + moves the cursor back one position when not already at the start
    + does nothing when at position zero
    # input
  typing_test.is_finished
    @ (state: game_state) -> bool
    + returns true when the cursor has reached the end of the target
    - returns false otherwise
    # status
  typing_test.accuracy
    @ (state: game_state) -> f64
    + returns the fraction of correct keystrokes over total keystrokes
    + returns 1.0 when no keystrokes have been made
    # metrics
  typing_test.words_per_minute
    @ (state: game_state, now_millis: i64) -> f64
    + returns characters-typed divided by five, scaled to per-minute using elapsed time
    + returns 0.0 when elapsed time is zero
    # metrics
    -> std.time.now_millis
