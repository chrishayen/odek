# Requirement: "a dependency injection container"

A container that registers named components and resolves them with their dependencies wired in. Components are opaque bean handles.

std
  std.sync
    std.sync.mutex_new
      fn () -> mutex_handle
      + creates a new unlocked mutex
      # concurrency
    std.sync.mutex_with
      fn (m: mutex_handle, f: fn() -> void) -> void
      + runs f while holding the mutex and releases on return
      # concurrency

di
  di.container_new
    fn () -> container_state
    + creates an empty container with no registered beans
    # construction
    -> std.sync.mutex_new
  di.register_singleton
    fn (c: container_state, name: string, factory: fn(container_state) -> bean) -> result[void, string]
    + registers a factory that will be invoked at most once per container
    - returns error when a bean with the same name is already registered
    # registration
    -> std.sync.mutex_with
  di.register_prototype
    fn (c: container_state, name: string, factory: fn(container_state) -> bean) -> result[void, string]
    + registers a factory that is invoked on every resolve call
    - returns error when a bean with the same name is already registered
    # registration
    -> std.sync.mutex_with
  di.resolve
    fn (c: container_state, name: string) -> result[bean, string]
    + returns the bean for name, constructing singletons on first access
    + allows factories to call resolve recursively for their dependencies
    - returns error when the name has no registration
    - returns error on a cyclic dependency between singletons
    # resolution
    -> std.sync.mutex_with
  di.resolve_all
    fn (c: container_state, tag: string) -> result[list[bean], string]
    + returns every bean registered under the given tag in registration order
    ? tags are attached at registration via a separate call, omitted here for brevity
    # resolution
