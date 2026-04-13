# Requirement: "a GraphQL field builder with typed resolvers"

A fluent builder for declaring GraphQL fields with typed argument schemas, resolvers, and composable middleware.

std: (all units exist)

field_builder
  field_builder.new
    @ (name: string, return_type: string) -> field_state
    + creates a field with the given name and return type
    # construction
  field_builder.arg
    @ (state: field_state, name: string, type_name: string, required: bool) -> field_state
    + declares an argument on the field
    # schema
  field_builder.description
    @ (state: field_state, text: string) -> field_state
    + attaches a human-readable description
    # schema
  field_builder.deprecated
    @ (state: field_state, reason: string) -> field_state
    + marks the field deprecated with a reason string
    # schema
  field_builder.resolver
    @ (state: field_state, resolve: closure[map[string, string]]) -> field_state
    + attaches the function that computes the field value from its arguments
    # resolution
  field_builder.middleware
    @ (state: field_state, wrap: closure[map[string, string]]) -> field_state
    + pushes a middleware that wraps the resolver chain
    ? middlewares run in the order they were pushed
    # middleware
  field_builder.build
    @ (state: field_state) -> result[field_definition, string]
    + finalizes the field into an immutable schema entry
    - returns error when the resolver is missing
    - returns error on duplicate argument names
    # finalization
