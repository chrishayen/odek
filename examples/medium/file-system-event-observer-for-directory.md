# Requirement: "a library for observing file system events under a directory"

Opens a watch on a directory tree, delivers create/modify/delete events to subscribers, and supports stopping the watch.

std
  std.fs
    std.fs.watch_open
      @ (root: string, recursive: bool) -> result[watch_handle, string]
      + opens an OS-level watch on the given directory
      - returns error when the path does not exist
      # filesystem
    std.fs.watch_poll
      @ (handle: watch_handle, timeout_ms: i32) -> result[list[fs_raw_event], string]
      + returns raw events received since the last poll, blocking up to timeout_ms
      # filesystem
    std.fs.watch_close
      @ (handle: watch_handle) -> result[void, string]
      + releases the watch and its underlying resources
      # filesystem
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

fswatch
  fswatch.start
    @ (root: string, recursive: bool) -> result[fswatch_state, string]
    + begins observing events under root
    - returns error when root is not a directory
    # lifecycle
    -> std.fs.watch_open
  fswatch.stop
    @ (state: fswatch_state) -> result[void, string]
    + cancels the observation and releases the underlying watch
    # lifecycle
    -> std.fs.watch_close
  fswatch.next_events
    @ (state: fswatch_state, timeout_ms: i32) -> result[list[fs_event], string]
    + returns the next batch of normalized events
    + empty list when timeout elapses with no activity
    # polling
    -> std.fs.watch_poll
    -> std.time.now_millis
  fswatch.classify
    @ (raw: fs_raw_event) -> fs_event
    + converts raw OS flags into one of created, modified, deleted, moved
    # classification
  fswatch.debounce
    @ (events: list[fs_event], window_ms: i32) -> list[fs_event]
    + collapses consecutive events for the same path within the window into one event
    ? keeps the last event in each window
    # debouncing
