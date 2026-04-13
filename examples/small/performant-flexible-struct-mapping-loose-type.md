# Requirement: "a flexible struct mapping library with loose type conversion between fields"

Copies field values between two structured records by name, converting compatible types along the way.

std: (all units exist)

struct_map
  struct_map.describe
    @ (fields: list[tuple[string, string]]) -> struct_schema
    + returns a schema listing field names and their declared type tags
    # schema
  struct_map.map
    @ (src_schema: struct_schema, src_values: map[string, string], dst_schema: struct_schema) -> result[map[string, string], list[string]]
    + returns destination field values populated from matching source fields
    - returns the list of field names that could not be converted
    ? only fields present in both schemas are considered
    # mapping
  struct_map.coerce
    @ (value: string, from_type: string, to_type: string) -> result[string, string]
    + returns the value converted between compatible scalar types
    - returns error when the conversion is not defined (e.g. bool to float)
    # coercion
  struct_map.rename_rule
    @ (from_field: string, to_field: string) -> field_rule
    + returns a rule that maps one source field name to a different destination name
    # rules
