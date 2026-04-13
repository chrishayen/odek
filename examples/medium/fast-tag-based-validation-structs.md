# Requirement: "a tag-driven record validator"

Callers describe a record schema as fields with constraint tags (e.g. "min=3", "max=100", "required", "email"). The library compiles the schema once and validates records against it, returning per-field errors.

std
  std.regex
    std.regex.match
      @ (pattern: string, input: string) -> bool
      + returns true when input matches pattern
      # regex

validator
  validator.parse_tag
    @ (raw: string) -> result[list[constraint], string]
    + parses a comma-separated tag string into constraints
    - returns error on an unknown constraint name
    # tag_parsing
  validator.compile_schema
    @ (fields: list[field_spec]) -> result[schema, string]
    + returns a compiled schema ready for validation
    - returns error when any field's tag is invalid
    # compilation
  validator.check_field
    @ (c: constraint, value: field_value) -> optional[string]
    + returns none when the constraint is satisfied
    - returns a message like "must be at least 3 characters" when violated
    # field_checking
    -> std.regex.match
  validator.validate
    @ (s: schema, record: map[string, field_value]) -> result[void, list[field_error]]
    + returns ok when every field passes every constraint
    - returns the list of failing (field_name, message) pairs
    # validation
