# Requirement: "a library that watches a source tree and restarts a build/run command when files change"

Project layer owns the watcher loop and process supervisor; std provides filesystem watching, process spawning, and time.

std
  std.fs
    std.fs.watch_dir
      @ (path: string) -> result[watcher, string]
      + starts a recursive directory watcher
      - returns error when the path does not exist
      # filesystem
    std.fs.next_event
      @ (w: watcher) -> result[fs_event, string]
      + blocks until the next filesystem change event
      # filesystem
  std.process
    std.process.spawn
      @ (cmd: string, args: list[string]) -> result[process_handle, string]
      + spawns a child process with the given arguments
      - returns error when the executable is not found
      # process
    std.process.kill
      @ (p: process_handle) -> result[void, string]
      + terminates the child process
      # process

rerunner
  rerunner.new
    @ (watch_dir: string, build_cmd: string, run_cmd: string) -> rerunner_state
    + creates a rerunner for the given directory and commands
    # construction
  rerunner.trigger_rebuild
    @ (r: rerunner_state) -> result[rerunner_state, string]
    + kills the current run, invokes the build command, then starts the run command
    - returns error when the build command exits non-zero
    # lifecycle
    -> std.process.spawn
    -> std.process.kill
  rerunner.watch_loop
    @ (r: rerunner_state) -> result[void, string]
    + watches the directory and calls trigger_rebuild on each relevant change
    ? debounces multiple events inside a short window into one rebuild
    # control
    -> std.fs.watch_dir
    -> std.fs.next_event
  rerunner.should_rebuild
    @ (event: fs_event) -> bool
    + returns true for changes to source files (not temp/build artifacts)
    # filtering
