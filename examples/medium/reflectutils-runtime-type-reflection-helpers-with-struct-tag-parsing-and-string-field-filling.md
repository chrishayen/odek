# Requirement: "helpers for runtime type reflection: struct tag parsing, recursive value walking, filling a field from a string"

The "struct" here is an abstract record described by a type_info value the host language supplies.

std: (all units exist)

reflectutils
  reflectutils.parse_tag
    fn (raw: string) -> map[string, string]
    + parses a tag string like `json:"name,omitempty" db:"user_name"` into a key-to-value map
    + values may themselves contain commas and are returned verbatim
    - unclosed quotes yield an empty map
    # tag_parsing
  reflectutils.tag_options
    fn (tag_value: string) -> tuple[string, list[string]]
    + splits a single tag value on commas, returning (primary, options)
    ? example: "name,omitempty" -> ("name", ["omitempty"])
    # tag_parsing
  reflectutils.walk
    fn (root: type_info, visit: fn(path: string, field: field_info) -> void) -> void
    + performs a depth-first traversal of a record type, invoking visit for every field
    + nested records are descended; leaf types halt recursion
    + paths use dot notation (e.g. "user.address.city")
    # traversal
  reflectutils.fill_from_string
    fn (field: field_info, text: string) -> result[field_info, string]
    + parses text according to the field's declared primitive type and stores it
    + supports integer, float, bool, and string kinds
    - returns error when text cannot be parsed as the target kind
    - returns error when the field kind is non-primitive
    # coercion
  reflectutils.fill_record_from_map
    fn (root: type_info, values: map[string, string]) -> result[type_info, list[string]]
    + assigns each entry in values to a field selected by the map key as a dot path
    - returns the list of error messages when any field fails to parse
    # coercion
    -> reflectutils.walk
    -> reflectutils.fill_from_string
