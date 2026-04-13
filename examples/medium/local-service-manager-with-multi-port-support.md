# Requirement: "a library to manage multiple local services running on different ports"

Registers named services with their ports and start commands, starts and stops them, reports health, and exposes a lookup from service name to port.

std
  std.net
    std.net.tcp_probe
      @ (host: string, port: u16, timeout_ms: i32) -> bool
      + returns true when a TCP connection to host:port succeeds within the timeout
      # networking
  std.process
    std.process.spawn
      @ (command: string, args: list[string], env: map[string, string]) -> result[process_handle, string]
      + launches a process in the background and returns a handle
      - returns error when the command cannot be launched
      # process
    std.process.kill
      @ (handle: process_handle) -> result[void, string]
      + terminates a running process
      - returns error when the process has already exited
      # process

service_manager
  service_manager.new
    @ () -> manager_state
    + creates an empty service registry
    # construction
  service_manager.register
    @ (state: manager_state, name: string, port: u16, command: string, args: list[string]) -> result[manager_state, string]
    + registers a service with its listening port and start command
    - returns error when the name already exists
    - returns error when another service is already bound to that port
    # registration
  service_manager.start
    @ (state: manager_state, name: string) -> result[manager_state, string]
    + starts the named service and records its process handle
    - returns error when the service is not registered
    - returns error when the service is already running
    # lifecycle
    -> std.process.spawn
  service_manager.stop
    @ (state: manager_state, name: string) -> result[manager_state, string]
    + stops the named service
    - returns error when the service is not running
    # lifecycle
    -> std.process.kill
  service_manager.health
    @ (state: manager_state, name: string) -> result[bool, string]
    + returns true when the service port accepts TCP connections
    - returns error when the service is not registered
    # health
    -> std.net.tcp_probe
  service_manager.lookup_port
    @ (state: manager_state, name: string) -> result[u16, string]
    + returns the port registered for the given service name
    - returns error when the name is unknown
    # queries
  service_manager.list
    @ (state: manager_state) -> list[tuple[string, u16, bool]]
    + returns (name, port, running) for every registered service
    # queries
