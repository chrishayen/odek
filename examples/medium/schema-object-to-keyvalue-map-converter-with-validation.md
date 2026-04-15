# Requirement: "a library for converting structured objects to and from primitive key-value maps with field-level validation"

A schema-driven serializer: the caller declares fields with types, the library loads from or dumps to a primitive map and reports field errors.

std: (all units exist)

schema
  schema.new
    fn () -> schema_def
    + creates an empty schema with no fields
    # construction
  schema.field
    fn (s: schema_def, name: string, kind: string, required: bool) -> schema_def
    + declares a field of the given primitive kind ("string","int","float","bool")
    ? returns the updated schema so callers can chain declarations
    # declaration
  schema.load
    fn (s: schema_def, raw: map[string, string]) -> result[map[string, primitive], map[string, string]]
    + returns the typed object when all required fields are present and coerce cleanly
    - returns a map of field_name -> error_message when any field fails validation
    # deserialization
  schema.dump
    fn (s: schema_def, obj: map[string, primitive]) -> result[map[string, string], string]
    + returns the stringified form of the object for each declared field
    - returns error when a required field is absent from obj
    # serialization
  schema.validate
    fn (s: schema_def, obj: map[string, primitive]) -> map[string, string]
    + returns a map of field_name -> error_message for every field that fails its rules
    + returns an empty map when the object is valid
    # validation
