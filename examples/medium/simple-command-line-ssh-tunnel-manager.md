# Requirement: "an SSH tunnel manager that starts, stops, and monitors named port forwards"

Tunnels are described in a config, launched as child processes, and tracked by name.

std
  std.proc
    std.proc.spawn
      @ (argv: list[string]) -> result[proc_handle, string]
      + starts a child process and returns a handle
      - returns error when the executable cannot be found
      # process
    std.proc.is_running
      @ (handle: proc_handle) -> bool
      + returns whether the child is still alive
      # process
    std.proc.kill
      @ (handle: proc_handle) -> result[void, string]
      + sends a termination signal to the child
      # process
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem

tunnel_manager
  tunnel_manager.load_config
    @ (path: string) -> result[list[tunnel_spec], string]
    + parses a tunnel configuration into a list of specs
    - returns error when the file is missing required fields
    # config
    -> std.fs.read_all
  tunnel_manager.new
    @ (specs: list[tunnel_spec]) -> manager_state
    + creates a manager seeded with the given specs and no running tunnels
    # construction
  tunnel_manager.start
    @ (state: manager_state, name: string) -> result[manager_state, string]
    + spawns the ssh client for the named tunnel and records its handle
    - returns error when the name is unknown
    - returns error when the tunnel is already running
    # lifecycle
    -> std.proc.spawn
  tunnel_manager.stop
    @ (state: manager_state, name: string) -> result[manager_state, string]
    + terminates the named tunnel and forgets its handle
    - returns error when the name is unknown
    # lifecycle
    -> std.proc.kill
  tunnel_manager.status
    @ (state: manager_state) -> list[tunnel_status]
    + returns the liveness of every tunnel in the manager
    # monitoring
    -> std.proc.is_running
  tunnel_manager.restart_dead
    @ (state: manager_state) -> manager_state
    + respawns any tunnel whose child has exited since last check
    # supervision
    -> std.proc.is_running
    -> std.proc.spawn
