# Requirement: "a development and testing environment management framework"

Declares required services for a workspace, starts them in dependency order, waits for readiness, and tears them down cleanly.

std
  std.process
    std.process.spawn
      @ (program: string, args: list[string], env: map[string, string]) -> result[process_handle, string]
      + starts a subprocess and returns a handle
      - returns error when the program cannot be launched
      # process
    std.process.terminate
      @ (h: process_handle) -> result[void, string]
      + signals the process to exit and waits for it
      # process
  std.net
    std.net.tcp_probe
      @ (host: string, port: u16) -> bool
      + returns true when a TCP connection can be established
      # network
  std.time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + blocks the current task for the given duration
      # time

envite
  envite.new_environment
    @ (name: string) -> environment_spec
    + creates an empty environment specification
    # construction
  envite.add_service
    @ (env: environment_spec, svc: service_spec) -> environment_spec
    + registers a service with its start command, ports, and dependencies
    # registration
  envite.resolve_order
    @ (env: environment_spec) -> result[list[string], string]
    + topologically orders services by their declared dependencies
    - returns error when the dependency graph has a cycle
    # planning
  envite.start_service
    @ (svc: service_spec) -> result[service_handle, string]
    + launches a single service and returns a handle
    - returns error when the service fails to start
    # lifecycle
    -> std.process.spawn
  envite.wait_ready
    @ (svc: service_spec, timeout_ms: i64) -> result[void, string]
    + polls the service's readiness probe until it succeeds or the timeout elapses
    - returns error on timeout
    # readiness
    -> std.net.tcp_probe
    -> std.time.sleep_millis
  envite.start_all
    @ (env: environment_spec) -> result[environment_state, string]
    + starts every service in dependency order, waiting for readiness between each
    - returns error on the first service that fails to start or become ready
    # orchestration
  envite.stop_all
    @ (state: environment_state) -> result[void, string]
    + terminates services in reverse dependency order
    # shutdown
    -> std.process.terminate
