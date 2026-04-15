# Requirement: "a searchable shell history index"

Ingests command history entries and returns ranked suggestions for a query. Ranking is frequency and recency.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file to bytes
      - returns error when the file is missing
      # filesystem

history
  history.new
    fn () -> history_index
    + creates an empty history index
    # construction
  history.record
    fn (idx: history_index, command: string, cwd: string, exit_code: i32) -> history_index
    + adds a command execution with current timestamp
    # ingestion
    -> std.time.now_seconds
  history.load_file
    fn (idx: history_index, path: string) -> result[history_index, string]
    + imports newline-separated commands from a shell history file
    - returns error when the path cannot be read
    # ingestion
    -> std.fs.read_all
  history.search
    fn (idx: history_index, query: string, limit: i32) -> list[history_suggestion]
    + returns up to limit matches ranked by a combination of frequency and recency
    + returns an empty list when nothing matches
    # search
    -> std.time.now_seconds
  history.score
    fn (entry: history_entry, query: string, now: i64) -> f64
    + computes the rank score for an entry given the current time
    + returns 0.0 when the query does not appear in the command
    # ranking
  history.forget
    fn (idx: history_index, command: string) -> history_index
    + removes all occurrences of a command from the index
    # management
