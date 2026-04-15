# Requirement: "a library that restarts a supervised process when it exits or when watched files change"

A supervisor owns a child process, restarts it on exit, and also restarts it when any watched file's modification time advances.

std
  std.process
    std.process.spawn
      fn (command: list[string]) -> result[process_handle, string]
      + starts the process and returns its handle
      - returns error when the command cannot be launched
      # process
    std.process.wait
      fn (handle: process_handle) -> i32
      + blocks until the process exits and returns its exit code
      # process
  std.fs
    std.fs.mtime
      fn (path: string) -> result[i64, string]
      + returns the file's modification time in seconds
      - returns error when the path does not exist
      # filesystem

supervisor
  supervisor.new
    fn (command: list[string], watched: list[string]) -> supervisor_state
    + creates a supervisor that runs the command and watches the given paths
    # construction
  supervisor.start
    fn (state: supervisor_state) -> result[supervisor_state, string]
    + spawns the child and records its handle on the state
    # lifecycle
    -> std.process.spawn
  supervisor.poll
    fn (state: supervisor_state) -> supervisor_event
    + returns `exited` when the child has terminated or `changed` when any watched file's mtime advanced
    # watching
    -> std.fs.mtime
  supervisor.restart
    fn (state: supervisor_state) -> result[supervisor_state, string]
    + stops the current child and starts a fresh one
    # lifecycle
    -> std.process.wait
    -> std.process.spawn
