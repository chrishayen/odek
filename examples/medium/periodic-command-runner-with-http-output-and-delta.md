# Requirement: "a library that runs a command periodically and exposes the latest output and its delta over HTTP"

Core library manages a polling loop that snapshots command output and exposes a read model. The HTTP surface is a thin handler that reads the snapshot.

std
  std.process
    std.process.run_capture
      fn (cmd: string, args: list[string]) -> result[string, string]
      + runs the command and returns its captured standard output
      - returns error when the command exits non-zero or fails to start
      # process
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.diff
    std.diff.line_diff
      fn (old: string, new_text: string) -> string
      + returns a unified line diff between two multiline strings
      # diff

watch
  watch.new_watcher
    fn (cmd: string, args: list[string], interval_ms: i64) -> watcher_state
    + creates a watcher with no snapshot yet
    # construction
  watch.tick
    fn (state: watcher_state) -> watcher_state
    + when enough time has elapsed, re-runs the command and stores the new snapshot and delta
    ? snapshot includes both the raw latest output and a diff against the previous snapshot
    # polling
    -> std.process.run_capture
    -> std.time.now_millis
    -> std.diff.line_diff
  watch.latest
    fn (state: watcher_state) -> optional[snapshot]
    + returns the most recent snapshot or none if no tick has produced one yet
    # read_model
  watch.latest_delta
    fn (state: watcher_state) -> optional[string]
    + returns the diff between the two most recent snapshots
    # read_model
  watch.handle_http
    fn (state: watcher_state, path: string) -> http_response
    + routes "/" to the latest snapshot and "/delta" to the latest delta
    - returns 404 for any other path
    # http
