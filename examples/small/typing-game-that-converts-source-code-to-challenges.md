# Requirement: "a typing game that turns source code into typing challenges"

Given a source file, slice it into bite-sized challenges and score a typed attempt against the target.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the whole file into memory
      - returns error when the path does not exist
      # io

typing_game
  typing_game.extract_challenges
    fn (source: string, max_lines: i32) -> list[challenge]
    + splits source into challenges of at most max_lines lines each
    + skips blank-only segments
    # extraction
  typing_game.load_challenges
    fn (path: string, max_lines: i32) -> result[list[challenge], string]
    + reads a file and slices it into challenges
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
    -> typing_game.extract_challenges
  typing_game.score_attempt
    fn (target: string, typed: string, elapsed_ms: i64) -> attempt_score
    + returns characters-per-minute, accuracy as a fraction, and error count
    + accuracy is 0.0 when typed is empty
    # scoring
  typing_game.format_diff
    fn (target: string, typed: string) -> list[char_status]
    + returns a per-character status marking correct, wrong, and missing characters
    # feedback
