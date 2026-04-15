# Requirement: "a simple dependency injection container"

Registers factories keyed by type name and resolves instances on demand. Singleton and transient scopes are supported.

std: (all units exist)

di_container
  di_container.new
    fn () -> container_state
    + returns an empty container
    # construction
  di_container.register_singleton
    fn (state: container_state, type_name: string, factory_id: string) -> container_state
    + registers a factory whose result is cached after first resolve
    # registration
  di_container.register_transient
    fn (state: container_state, type_name: string, factory_id: string) -> container_state
    + registers a factory whose result is produced fresh every resolve
    # registration
  di_container.resolve
    fn (state: container_state, type_name: string) -> result[tuple[string, container_state], string]
    + returns (instance_handle, new_state) using the registered factory
    - returns error when the type name is not registered
    # resolution
