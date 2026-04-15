# Requirement: "a file system event notification library"

Registers callbacks for events on a path and delivers created, modified, and removed events.

std
  std.fs
    std.fs.watch
      fn (path: string) -> result[watch_handle, string]
      + opens an OS file-watch descriptor for the path
      - returns error when path does not exist
      # io
    std.fs.read_events
      fn (handle: watch_handle) -> result[list[raw_event], string]
      + returns the next batch of raw events from the descriptor
      # io

notify
  notify.new
    fn () -> watcher_state
    + creates an empty watcher with no watches
    # construction
  notify.watch
    fn (state: watcher_state, path: string) -> result[watcher_state, string]
    + starts watching the given path, recursively for directories
    - returns error when path does not exist
    # subscription
    -> std.fs.watch
  notify.unwatch
    fn (state: watcher_state, path: string) -> watcher_state
    + stops watching the path, if present
    # subscription
  notify.next_events
    fn (state: watcher_state) -> result[list[fs_event], string]
    + returns the next batch of decoded events across all watches
    + event kinds are create, write, remove, rename
    # delivery
    -> std.fs.read_events
