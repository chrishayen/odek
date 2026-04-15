# Requirement: "a file-watching process supervisor that restarts a command on source changes"

Language-agnostic. Watches a set of paths, debounces bursts of events, and drives a child process through a simple state machine.

std
  std.fs
    std.fs.watch
      fn (paths: list[string]) -> result[watch_handle, string]
      + returns a handle that yields file-change events for the given paths recursively
      - returns error when any path does not exist
      # filesystem
    std.fs.next_event
      fn (handle: watch_handle) -> optional[fs_event]
      + returns the next pending event or none if the handle is closed
      # filesystem
  std.proc
    std.proc.spawn
      fn (argv: list[string]) -> result[proc_handle, string]
      + starts a child process with the given argument vector
      - returns error when the executable cannot be found
      # process
    std.proc.kill
      fn (handle: proc_handle) -> result[void, string]
      + sends a termination signal to the child
      # process
    std.proc.wait
      fn (handle: proc_handle) -> i32
      + blocks until the child exits and returns its exit code
      # process
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current time in milliseconds since the epoch
      # time

watcher
  watcher.new
    fn (paths: list[string], argv: list[string], debounce_ms: i32) -> watcher_state
    + creates a watcher for the given paths and command with the supplied debounce window
    # construction
  watcher.start
    fn (state: watcher_state) -> result[watcher_state, string]
    + spawns the child and begins watching the paths
    - returns error when no path can be watched
    # lifecycle
    -> std.fs.watch
    -> std.proc.spawn
  watcher.tick
    fn (state: watcher_state) -> watcher_state
    + consumes any pending events and, if the debounce window has elapsed since the last change, restarts the child
    # supervision
    -> std.fs.next_event
    -> std.time.now_millis
    -> std.proc.kill
    -> std.proc.spawn
  watcher.stop
    fn (state: watcher_state) -> void
    + kills the child and closes the watch handle
    # lifecycle
    -> std.proc.kill
    -> std.proc.wait
