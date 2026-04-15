# Requirement: "a source-file watcher that recompiles and restarts a target on change"

Same shape as a generic auto-restarter, but the restart step runs a build command before launching the target.

std
  std.fs
    std.fs.watch_dir
      fn (path: string, recursive: bool) -> watcher_state
      + returns a handle that emits filesystem change events
      # filesystem
    std.fs.next_event
      fn (w: watcher_state) -> optional[string]
      + returns the next changed path or none
      # filesystem
  std.proc
    std.proc.run
      fn (command: string, args: list[string]) -> result[i32, string]
      + runs a command to completion and returns its exit code
      # process
    std.proc.spawn
      fn (command: string, args: list[string]) -> result[i32, string]
      + starts a child process and returns its pid
      # process
    std.proc.kill
      fn (pid: i32) -> void
      + terminates the process identified by pid
      # process

dev_loop
  dev_loop.new
    fn (watch_path: string, build: list[string], run: list[string]) -> dev_loop_state
    + creates a supervisor bound to build and run commands
    # construction
    -> std.fs.watch_dir
  dev_loop.build_and_start
    fn (s: dev_loop_state) -> result[dev_loop_state, string]
    + runs the build command then starts the target
    - returns error when build exits non-zero
    # lifecycle
    -> std.proc.run
    -> std.proc.spawn
  dev_loop.tick
    fn (s: dev_loop_state) -> dev_loop_state
    + when events are pending, kills the target, rebuilds, and restarts
    # supervision
    -> std.fs.next_event
    -> std.proc.kill
    -> std.proc.run
    -> std.proc.spawn
