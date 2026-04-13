# Requirement: "a library for describing and constructing types at runtime"

A registry for runtime type descriptors (named structs with typed fields) and a factory that produces zero-valued instances.

std: (all units exist)

type_registry
  type_registry.new
    @ () -> registry_state
    + creates an empty type registry
    # construction
  type_registry.define
    @ (state: registry_state, name: string, fields: list[field_spec]) -> result[registry_state, string]
    + registers a new type by name with an ordered field list
    - returns error when name is empty
    - returns error when a type with that name already exists
    - returns error when two fields share a name
    # definition
  type_registry.lookup
    @ (state: registry_state, name: string) -> optional[type_descriptor]
    + returns the descriptor for a registered type
    # lookup
  type_registry.instantiate
    @ (state: registry_state, name: string) -> result[dynamic_value, string]
    + returns a fresh value with all fields at their type's zero value
    - returns error when the type is unknown
    # instantiation
  type_registry.set_field
    @ (value: dynamic_value, field: string, v: dynamic_value) -> result[dynamic_value, string]
    + sets a field on a dynamic value when the field type matches
    - returns error when the field does not exist
    - returns error when the value type does not match the field type
    # mutation
