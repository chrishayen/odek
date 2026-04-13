# Requirement: "a request property validation library"

Declarative rules validate the fields of an incoming request body map and return a list of errors.

std: (all units exist)

prop_validator
  prop_validator.required
    @ (field: string) -> rule
    + creates a rule that the field must be present and non-empty
    # rule
  prop_validator.is_int
    @ (field: string) -> rule
    + creates a rule that the field must parse as an integer
    # rule
  prop_validator.length_between
    @ (field: string, min: i32, max: i32) -> rule
    + creates a rule that the field length must be in [min, max]
    # rule
  prop_validator.matches
    @ (field: string, pattern: string) -> rule
    + creates a rule that the field must match the given regex pattern
    # rule
  prop_validator.validate
    @ (rules: list[rule], body: map[string, string]) -> list[validation_error]
    + returns an empty list when every rule passes
    - returns a validation_error per failing rule with the field name and message
    # validation
