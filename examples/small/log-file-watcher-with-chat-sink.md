# Requirement: "a library that watches log files and pushes new lines to a chat message sink"

Tails a file and forwards appended content line-by-line to a pluggable sink.

std
  std.fs
    std.fs.stat_size
      fn (path: string) -> result[i64, string]
      + returns the current byte length of the file
      - returns error when the path does not exist
      # filesystem
    std.fs.read_range
      fn (path: string, offset: i64, length: i64) -> result[bytes, string]
      + returns the requested byte range
      # filesystem

log_pusher
  log_pusher.new
    fn (path: string) -> result[watcher_state, string]
    + returns a watcher anchored at the current end of the file
    - returns error when the file does not exist
    # construction
    -> std.fs.stat_size
  log_pusher.poll
    fn (w: watcher_state) -> result[tuple[list[string], watcher_state], string]
    + returns newly appended lines since the last poll and the advanced state
    + handles file truncation by resetting the offset to zero
    - returns error when the file disappears
    # tailing
    -> std.fs.stat_size
    -> std.fs.read_range
  log_pusher.push_lines
    fn (lines: list[string], sink: fn(message: string) -> result[void, string]) -> result[i32, string]
    + returns the number of lines successfully delivered
    - returns error at the first sink failure
    ? caller supplies the sink, so the library is transport-agnostic
    # delivery
