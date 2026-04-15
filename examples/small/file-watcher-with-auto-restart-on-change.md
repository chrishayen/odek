# Requirement: "a file-watcher that restarts a target process when source files change"

The library exposes a watch loop that consumers drive; process start/stop is delegated to std primitives.

std
  std.fs
    std.fs.watch_dir
      fn (path: string, recursive: bool) -> watcher_state
      + returns a handle that emits filesystem change events
      # filesystem
    std.fs.next_event
      fn (w: watcher_state) -> optional[string]
      + returns the next changed path, or none when the buffer is empty
      # filesystem
  std.proc
    std.proc.spawn
      fn (command: string, args: list[string]) -> result[i32, string]
      + starts a child process and returns its pid
      - returns error when the command cannot be executed
      # process
    std.proc.kill
      fn (pid: i32) -> void
      + terminates the process identified by pid
      # process

auto_restart
  auto_restart.new
    fn (watch_path: string, command: string, args: list[string]) -> auto_restart_state
    + creates a supervisor bound to the watch path and restart command
    # construction
    -> std.fs.watch_dir
  auto_restart.start
    fn (s: auto_restart_state) -> result[auto_restart_state, string]
    + launches the child for the first time
    - returns error when spawn fails
    # lifecycle
    -> std.proc.spawn
  auto_restart.tick
    fn (s: auto_restart_state) -> auto_restart_state
    + consumes pending filesystem events and restarts the child when any arrive
    # supervision
    -> std.fs.next_event
    -> std.proc.kill
    -> std.proc.spawn
