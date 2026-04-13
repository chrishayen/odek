# Requirement: "a dependency injection and inversion-of-control container with component lifecycle hooks"

A container that registers component factories by type key, resolves their dependencies, and runs start/stop lifecycle phases in topological order.

std: (all units exist)

container
  container.new
    @ () -> container_state
    + creates an empty container
    # construction
  container.register
    @ (state: container_state, key: string, deps: list[string], factory: component_factory) -> result[container_state, string]
    + registers a component under a key with its dependency keys and a factory
    - returns error when the key is already registered
    # registration
  container.resolve
    @ (state: container_state, key: string) -> result[component, string]
    + returns the constructed component for the key, building dependencies first
    - returns error when the key is not registered
    - returns error when a dependency cycle is detected
    # resolution
  container.start_all
    @ (state: container_state) -> result[container_state, string]
    + invokes the start hook on every component in dependency order
    - returns error when any component's start hook fails
    # lifecycle
  container.stop_all
    @ (state: container_state) -> result[container_state, string]
    + invokes the stop hook on every component in reverse dependency order
    + continues stopping remaining components even if one fails and aggregates errors
    # lifecycle
  container.topological_order
    @ (state: container_state) -> result[list[string], string]
    + returns the registered keys in dependency order
    - returns error on cycles
    # ordering
