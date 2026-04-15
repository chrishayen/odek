# Requirement: "a library that fills record fields from a map using field annotations"

Given a schema of named fields with default values and a map of raw inputs, produce a fully populated record.

std: (all units exist)

field_filler
  field_filler.new_schema
    fn () -> fill_schema
    + returns an empty schema
    # construction
  field_filler.add_field
    fn (schema: fill_schema, name: string, default_value: string, required: bool) -> fill_schema
    + returns a schema with the field added
    # registration
  field_filler.fill
    fn (schema: fill_schema, inputs: map[string, string]) -> result[map[string, string], string]
    + returns a record with each field populated from inputs or defaults
    - returns error when a required field is missing from inputs
    # filling
