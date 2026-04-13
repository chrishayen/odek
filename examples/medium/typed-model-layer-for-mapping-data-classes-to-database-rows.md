# Requirement: "a typed model layer that maps data classes to database rows"

A lightweight ORM-style layer. Models are described once and used for both serialization and row mapping.

std: (all units exist)

typed_models
  typed_models.define_model
    @ (name: string, fields: list[field_spec]) -> model_def
    + creates a model definition with the given field name, type and nullability
    - returns a model_def marked invalid when two fields share a name
    # schema
  typed_models.to_create_table
    @ (m: model_def) -> string
    + returns a CREATE TABLE statement with columns and nullability matching the field specs
    # ddl
  typed_models.row_to_instance
    @ (m: model_def, row: map[string, sql_value]) -> result[instance, string]
    + returns an instance populated from row values
    - returns error when a non-nullable field is missing from row
    - returns error when a value's type does not match the field spec
    # mapping
  typed_models.instance_to_row
    @ (m: model_def, inst: instance) -> map[string, sql_value]
    + returns a column map suitable for parameterized insert
    + absent optional fields become sql NULL
    # mapping
  typed_models.validate
    @ (m: model_def, inst: instance) -> result[void, list[string]]
    + returns ok when every required field is set and types match
    - returns a list of field-specific error messages on any violation
    # validation
