# Requirement: "a directory navigation helper that learns from usage frequency"

The library keeps a frecency-scored database of directory entries and answers "best match for a hint" queries.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns the current unix time in seconds
      # time
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes a file, replacing existing contents
      - returns error when the parent directory is missing
      # filesystem

jump
  jump.new
    fn () -> jump_state
    + creates an empty navigation database
    # construction
  jump.visit
    fn (state: jump_state, path: string) -> jump_state
    + records a visit to the given path and updates its frecency score
    ? score combines visit count with recency decay
    # tracking
    -> std.time.now_seconds
  jump.query
    fn (state: jump_state, hint: string) -> optional[string]
    + returns the highest-scoring path whose name contains the hint
    - returns none when no path matches
    # query
  jump.list
    fn (state: jump_state) -> list[string]
    + returns all tracked paths in descending score order
    # query
  jump.forget
    fn (state: jump_state, path: string) -> jump_state
    + removes a path from the database
    # maintenance
  jump.save
    fn (state: jump_state, path: string) -> result[void, string]
    + serializes the database to disk
    - returns error when the file cannot be written
    # persistence
    -> std.fs.write_all
  jump.load
    fn (path: string) -> result[jump_state, string]
    + loads a previously saved database
    - returns error when the file is missing or corrupt
    # persistence
    -> std.fs.read_all
