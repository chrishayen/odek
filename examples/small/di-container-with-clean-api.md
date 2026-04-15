# Requirement: "a dependency injection container focused on a clean api"

Container stores constructors keyed by type name and resolves them lazily, caching singletons.

std: (all units exist)

di_container
  di_container.new
    fn () -> container_state
    + returns an empty container
    # construction
  di_container.register
    fn (state: container_state, type_name: string, factory: factory_fn) -> container_state
    + associates a factory with the given type name
    ? later registrations replace earlier ones for the same name
    # registration
  di_container.resolve
    fn (state: container_state, type_name: string) -> result[any_value, string]
    + invokes the factory for type_name and returns the value
    + caches the result so subsequent resolves return the same instance
    - returns error when type_name has no registered factory
    - returns error when the factory call returns an error
    # resolution
  di_container.resolve_fresh
    fn (state: container_state, type_name: string) -> result[any_value, string]
    + calls the factory without consulting or updating the singleton cache
    - returns error when type_name has no registered factory
    # resolution
