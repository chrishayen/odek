# Requirement: "a dependency injection framework"

A small container that resolves typed dependencies lazily via registered providers.

std: (all units exist)

di
  di.new_container
    fn () -> container_state
    + creates an empty dependency container
    # construction
  di.provide
    fn (c: container_state, key: string, factory: provider_fn) -> container_state
    + registers a factory for the given key
    ? factories are called at most once; the resolved value is cached
    # registration
  di.resolve
    fn (c: container_state, key: string) -> result[any_value, string]
    + returns the value for a key, invoking its factory on first access
    - returns error when no factory is registered for the key
    - returns error when the factory raises a circular dependency
    # resolution
  di.shutdown
    fn (c: container_state) -> void
    + releases cached values in reverse registration order
    # lifecycle
