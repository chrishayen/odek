# Requirement: "a data validation library primarily intended for validating data from untrusted sources"

Users build a schema out of composable validators, then apply it to unknown maps. No std primitives needed: the logic is purely structural.

std: (all units exist)

validator
  validator.schema
    @ () -> schema_state
    + creates an empty schema
    # construction
  validator.field_string
    @ (state: schema_state, name: string, required: bool) -> schema_state
    + declares a field that must be a string when present
    # schema
  validator.field_int
    @ (state: schema_state, name: string, required: bool, min: i64, max: i64) -> schema_state
    + declares an integer field constrained to [min, max] inclusive
    # schema
  validator.field_enum
    @ (state: schema_state, name: string, required: bool, choices: list[string]) -> schema_state
    + declares a field whose string value must be one of choices
    # schema
  validator.field_pattern
    @ (state: schema_state, name: string, required: bool, regex: string) -> schema_state
    + declares a string field that must match regex
    # schema
  validator.nested
    @ (state: schema_state, name: string, required: bool, inner: schema_state) -> schema_state
    + declares a nested object field validated by inner
    # schema
  validator.validate
    @ (state: schema_state, input: map[string, string]) -> result[map[string, string], list[string]]
    + returns the coerced map when every declared field passes
    - returns a list of error messages when one or more fields fail or are missing
    - returns error when the input contains an unknown field
    # validation
