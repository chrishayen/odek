# Requirement: "a library of simple functions to manipulate structured records"

A struct here is a dynamic string-keyed record. Operations are pure: they return new records.

std: (all units exist)

structs
  structs.to_map
    fn (s: struct_record) -> map[string, string]
    + returns a map of every field name to its string form
    # introspection
  structs.field_names
    fn (s: struct_record) -> list[string]
    + returns the declared field names in declaration order
    # introspection
  structs.has_field
    fn (s: struct_record, name: string) -> bool
    + returns true when the named field exists
    - returns false for unknown names
    # introspection
  structs.get
    fn (s: struct_record, name: string) -> result[string, string]
    + returns the string form of the named field
    - returns error when the field does not exist
    # access
  structs.set
    fn (s: struct_record, name: string, value: string) -> result[struct_record, string]
    + returns a new record with the named field replaced
    - returns error when the field does not exist
    # update
  structs.merge
    fn (base: struct_record, overlay: struct_record) -> struct_record
    + returns a new record with overlay's fields overriding base's
    ? fields present only in overlay are added
    # merge
  structs.equal
    fn (a: struct_record, b: struct_record) -> bool
    + returns true when both records have the same field names and values
    # comparison
