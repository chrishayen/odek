# Requirement: "an object schema description language and validator"

Schemas are built with a small fluent API and applied to dynamic values producing structured errors.

std: (all units exist)

schema
  schema.string
    @ () -> schema_node
    + returns a schema node that accepts string values
    # schema_builder
  schema.number
    @ () -> schema_node
    + returns a schema node that accepts numeric values
    # schema_builder
  schema.object
    @ (fields: map[string, schema_node]) -> schema_node
    + returns a schema node that accepts objects with the given field schemas
    # schema_builder
  schema.array
    @ (item: schema_node) -> schema_node
    + returns a schema node that accepts arrays whose elements match item
    # schema_builder
  schema.required
    @ (node: schema_node) -> schema_node
    + marks the node as required; validation fails when absent
    # schema_builder
  schema.min
    @ (node: schema_node, n: f64) -> schema_node
    + attaches a minimum constraint for numbers and min length for strings/arrays
    # schema_builder
  schema.max
    @ (node: schema_node, n: f64) -> schema_node
    + attaches a maximum constraint for numbers and max length for strings/arrays
    # schema_builder
  schema.pattern
    @ (node: schema_node, regex: string) -> schema_node
    + attaches a regex pattern constraint to a string node
    # schema_builder
  schema.validate
    @ (node: schema_node, value: dynamic_value) -> list[validation_error]
    + returns an empty list when the value satisfies the schema
    - returns errors with path and reason for each violation
    - reports type mismatches before descending into children
    # validation
