# Requirement: "a dependency injection library using type reflection"

Resolves dependencies by inspecting constructor argument types against a registry.

std: (all units exist)

autowire
  autowire.new_registry
    fn () -> registry_state
    + creates an empty type registry
    # construction
  autowire.register
    fn (r: registry_state, type_name: string, constructor: ctor_fn) -> registry_state
    + records a constructor keyed by its produced type name
    # registration
  autowire.build
    fn (r: registry_state, type_name: string) -> result[any_value, string]
    + recursively resolves the constructor's argument types and invokes it
    - returns error when a required argument type is unregistered
    - returns error when a dependency cycle is detected
    # resolution
