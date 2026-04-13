# Requirement: "a type-safe dependency-injection framework that wires providers into handlers by matching parameter types"

Providers declare what types they produce; handlers declare what types they consume. The framework topologically orders providers so every handler gets its inputs.

std: (all units exist)

inject
  inject.new_container
    @ () -> container_state
    + returns an empty container
    # construction
  inject.provide
    @ (state: container_state, produces: list[string], consumes: list[string], factory_id: string) -> result[container_state, string]
    + registers a provider described by the types it produces and consumes
    - returns error when the same type is produced by two providers
    # registration
  inject.resolve_order
    @ (state: container_state, required_types: list[string]) -> result[list[string], string]
    + returns factory ids in a topological order that satisfies the required types
    - returns error when a required type has no provider
    - returns error on cycles
    # planning
  inject.invoke
    @ (state: container_state, handler_consumes: list[string], run_factory: fn(string, map[string,bytes]) -> result[map[string,bytes],string]) -> result[map[string,bytes], string]
    + executes providers in order and returns the type-to-value map requested by the handler
    - returns error when any factory fails
    # execution
    -> inject.resolve_order
  inject.check
    @ (state: container_state, handler_consumes: list[string]) -> result[void, string]
    + validates a handler without executing any factories
    - returns error listing all missing or ambiguous types
    # validation
    -> inject.resolve_order
