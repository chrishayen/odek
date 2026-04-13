# Requirement: "a library that triggers a callback in response to file modifications"

The library polls watched paths and emits change events. Actually running a command is not in scope — the caller decides what to do with each event.

std
  std.fs
    std.fs.stat_mtime
      @ (path: string) -> result[i64, string]
      + returns the modification time as unix seconds
      - returns error when path does not exist
      # filesystem

watcher
  watcher.new
    @ (paths: list[string]) -> watcher_state
    + creates a watcher tracking the given paths with unknown initial mtimes
    # construction
  watcher.poll
    @ (state: watcher_state) -> tuple[list[string], watcher_state]
    + returns the subset of paths whose mtime changed since the previous poll
    + on first call, returns all paths that currently exist
    - omits paths that no longer exist from the returned list
    # change_detection
    -> std.fs.stat_mtime
