# Requirement: "a library that builds and restarts a project when it crashes or watched files change"

A supervised build-and-run loop that reacts to file changes and process exits. The project layer orchestrates build, spawn, and watch primitives from std.

std
  std.fs
    std.fs.walk
      @ (root: string) -> list[string]
      + returns all file paths under root recursively
      # filesystem
    std.fs.mtime
      @ (path: string) -> result[i64, string]
      + returns the last-modified unix time in seconds
      - returns error when path does not exist
      # filesystem
  std.process
    std.process.spawn
      @ (cmd: string, args: list[string]) -> result[proc_handle, string]
      + starts a child process
      - returns error when the executable cannot be found
      # process
    std.process.wait
      @ (proc: proc_handle) -> i32
      + blocks until the process exits and returns the exit code
      # process
    std.process.kill
      @ (proc: proc_handle) -> void
      + sends a terminate signal
      # process
  std.time
    std.time.sleep_millis
      @ (ms: i32) -> void
      + pauses the current thread for ms milliseconds
      # time

supervisor
  supervisor.new
    @ (watch_root: string, build_cmd: list[string], run_cmd: list[string]) -> supervisor_state
    + creates a supervisor with empty file snapshot and no running process
    # construction
  supervisor.snapshot
    @ (root: string) -> map[string, i64]
    + builds a path-to-mtime map for all files under root
    # file_tracking
    -> std.fs.walk
    -> std.fs.mtime
  supervisor.changed_files
    @ (before: map[string,i64], after: map[string,i64]) -> list[string]
    + returns paths whose mtime differs, were added, or were removed
    # file_tracking
  supervisor.build_and_run
    @ (state: supervisor_state) -> result[supervisor_state, string]
    + runs the build command then spawns the run command
    - returns error when the build exits non-zero
    # lifecycle
    -> std.process.spawn
    -> std.process.wait
  supervisor.restart
    @ (state: supervisor_state) -> result[supervisor_state, string]
    + kills the running process and invokes build_and_run
    # lifecycle
    -> std.process.kill
  supervisor.tick
    @ (state: supervisor_state) -> supervisor_state
    + takes a fresh snapshot; restarts on file change or process exit
    ? one poll step; the caller loops with sleep_millis
    # supervision_loop
    -> std.time.sleep_millis
