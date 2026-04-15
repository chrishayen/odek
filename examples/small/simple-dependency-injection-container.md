# Requirement: "a simple dependency injection container"

A container that stores named singletons and lazy factories and resolves them on demand.

std: (all units exist)

container
  container.new
    fn () -> container_state
    + returns an empty container
    # construction
  container.register_value
    fn (state: container_state, name: string, value: bytes) -> container_state
    + stores a ready-made dependency under the given name
    + replaces any existing binding with the same name
    # binding
  container.register_factory
    fn (state: container_state, name: string, factory: fn(container_state) -> result[bytes, string]) -> container_state
    + stores a lazy factory that builds the dependency on first resolve
    ? factories are invoked at most once; the result is cached
    # binding
  container.resolve
    fn (state: container_state, name: string) -> result[tuple[bytes, container_state], string]
    + returns the stored value for a registered name
    + invokes and caches the factory on first resolve
    - returns error when the name is not registered
    # resolution
