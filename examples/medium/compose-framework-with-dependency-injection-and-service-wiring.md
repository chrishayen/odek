# Requirement: "a framework for wiring composable application services with dependency injection"

A registry of named service factories with lazy resolution and cycle detection.

std: (all units exist)

compose
  compose.new
    fn () -> container_state
    + creates an empty container with no registered services
    # construction
  compose.register
    fn (c: container_state, name: string, deps: list[string], factory: factory_ref) -> result[container_state, string]
    + registers a service with its dependency names and a factory
    - returns error when name is already registered
    - returns error when name is empty
    # registration
  compose.resolve
    fn (c: container_state, name: string) -> result[tuple[service_instance, container_state], string]
    + returns the service, instantiating it and its dependencies lazily
    + caches the instance so subsequent resolves return the same value
    - returns error when the service is not registered
    - returns error when a dependency is not registered
    - returns error when a cycle is detected in the dependency graph
    # resolution
  compose.startup_order
    fn (c: container_state) -> result[list[string], string]
    + returns registered service names in dependency-first order
    - returns error when any service references an unknown dependency
    - returns error when the graph contains a cycle
    # lifecycle
  compose.names
    fn (c: container_state) -> list[string]
    + returns all registered service names
    # inspection
  compose.shutdown
    fn (c: container_state) -> container_state
    + drops all cached instances while keeping registrations
    # lifecycle
