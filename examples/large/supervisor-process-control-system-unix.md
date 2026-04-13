# Requirement: "a process supervisor and control system"

A library that launches child processes, monitors them, restarts on failure, and exposes a control API. Process spawning and signal delivery go through std primitives.

std
  std.process
    std.process.spawn
      @ (command: string, args: list[string]) -> result[i32, string]
      + starts a child process and returns its pid
      - returns error when the executable is not found
      # process
    std.process.wait
      @ (pid: i32) -> result[i32, string]
      + blocks until the process exits and returns its exit code
      - returns error when the pid does not exist
      # process
    std.process.signal
      @ (pid: i32, signal: i32) -> result[void, string]
      + sends the given signal to the process
      - returns error when the process has already exited
      # process
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.sleep_millis
      @ (millis: i64) -> void
      + blocks the current thread for the given duration
      # time

supervisor
  supervisor.new
    @ () -> supervisor_state
    + creates an empty supervisor with no managed programs
    # construction
  supervisor.add_program
    @ (state: supervisor_state, name: string, command: string, args: list[string], autorestart: bool) -> supervisor_state
    + registers a program definition under the given name
    - returns unchanged state when a program with that name already exists
    # configuration
  supervisor.start
    @ (state: supervisor_state, name: string) -> result[supervisor_state, string]
    + spawns the named program and records its pid and start time
    - returns error when the program name is not registered
    # lifecycle
    -> std.process.spawn
    -> std.time.now_seconds
  supervisor.stop
    @ (state: supervisor_state, name: string) -> result[supervisor_state, string]
    + sends the stop signal to the program and waits for exit
    - returns error when the program is not running
    # lifecycle
    -> std.process.signal
    -> std.process.wait
  supervisor.restart
    @ (state: supervisor_state, name: string) -> result[supervisor_state, string]
    + stops and then starts the program
    - returns error when the program is not registered
    # lifecycle
  supervisor.status
    @ (state: supervisor_state, name: string) -> result[program_status, string]
    + returns running/stopped/fatal along with pid and uptime
    - returns error when the program is not registered
    # query
    -> std.time.now_seconds
  supervisor.reap_exited
    @ (state: supervisor_state) -> supervisor_state
    + checks each running program, marks exited ones, and respawns those with autorestart
    ? caller drives this periodically; library does not spawn its own thread
    # reaping
    -> std.process.wait
    -> std.process.spawn
    -> std.time.now_seconds
  supervisor.list_programs
    @ (state: supervisor_state) -> list[string]
    + returns the names of all registered programs
    # query
  supervisor.shutdown_all
    @ (state: supervisor_state) -> supervisor_state
    + stops every running program and clears pids
    # lifecycle
    -> std.process.signal
    -> std.process.wait
