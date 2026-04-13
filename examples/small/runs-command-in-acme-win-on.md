# Requirement: "a library that watches a set of files and re-runs a command when any change is detected"

Small watcher facade: register paths, poll for modification changes, and invoke a callback. The filesystem poller is thin enough to swap out for an event-driven backend.

std
  std.fs
    std.fs.stat_mtime
      @ (path: string) -> result[i64, string]
      + returns the modification time of the path in milliseconds since epoch
      - returns error when the path does not exist
      # filesystem
  std.process
    std.process.run
      @ (cmd: string, args: list[string]) -> result[i32, string]
      + runs a command to completion and returns its exit code
      - returns error when the command fails to start
      # process

rerun
  rerun.new_watcher
    @ (paths: list[string]) -> watcher_state
    + creates a watcher tracking the given paths with no baseline yet
    # construction
  rerun.poll
    @ (state: watcher_state) -> tuple[watcher_state, list[string]]
    + stats each path and returns the updated state and the list of paths whose mtime changed since the last poll
    + the first poll after construction reports every existing path as changed
    # detection
    -> std.fs.stat_mtime
  rerun.run_on_change
    @ (state: watcher_state, cmd: string, args: list[string]) -> result[tuple[watcher_state, bool], string]
    + polls once; when any path changed, runs the command and returns (new_state, true)
    + returns (new_state, false) when nothing changed
    - propagates command execution errors
    # orchestration
    -> std.process.run
