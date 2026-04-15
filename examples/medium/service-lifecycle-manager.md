# Requirement: "a library to manage component lifecycles in a service"

Registers named components with start, stop, and health hooks, then starts them in dependency order and stops them in reverse.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

lifecycle
  lifecycle.new
    fn () -> lifecycle_state
    + creates an empty component registry
    # construction
  lifecycle.register
    fn (state: lifecycle_state, name: string, deps: list[string], start_fn: string, stop_fn: string) -> result[lifecycle_state, string]
    + registers a component with its dependencies and lifecycle hooks
    - returns error when the name already exists
    # registration
  lifecycle.start_order
    fn (state: lifecycle_state) -> result[list[string], string]
    + returns the topological start order
    - returns error when a dependency cycle is detected
    - returns error when a dependency is not registered
    # planning
  lifecycle.start_all
    fn (state: lifecycle_state) -> result[lifecycle_state, string]
    + invokes each component's start hook in dependency order and records the start time
    - returns error on the first component that fails, leaving earlier components running
    # startup
    -> std.time.now_millis
  lifecycle.stop_all
    fn (state: lifecycle_state) -> lifecycle_state
    + invokes each running component's stop hook in reverse start order
    + continues on errors so all components get a chance to shut down
    # shutdown
  lifecycle.status
    fn (state: lifecycle_state, name: string) -> result[string, string]
    + returns "registered", "running", "stopped", or "failed" for the named component
    - returns error when the name is unknown
    # queries
