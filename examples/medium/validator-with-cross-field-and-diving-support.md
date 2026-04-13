# Requirement: "a struct and field validation library with cross-field, cross-struct, map, slice, and array diving"

A validator over a schema of named rules. Callers declare rules, then validate a record and read back a list of violations.

std: (all units exist)

validator
  validator.new
    @ () -> validator_state
    + creates an empty validator with no registered rules
    # construction
  validator.register_field_rule
    @ (v: validator_state, field: string, rule: fn(string) -> bool, message: string) -> validator_state
    + adds a rule that checks a single field's string value
    # rule_registration
  validator.register_cross_field_rule
    @ (v: validator_state, name: string, rule: fn(map[string, string]) -> bool, message: string) -> validator_state
    + adds a rule that depends on two or more fields of the same record
    # rule_registration
  validator.register_collection_rule
    @ (v: validator_state, field: string, element_rule: fn(string) -> bool, message: string) -> validator_state
    + adds a rule that dives into a list or map field and checks each element
    ? maps are validated over their values; keys are ignored
    # rule_registration
  validator.validate
    @ (v: validator_state, record: map[string, string], collections: map[string, list[string]]) -> list[string]
    + returns an empty list when every rule passes
    - returns a list of messages, one per violated rule, in registration order
    # validation
  validator.validate_nested
    @ (v: validator_state, parent: map[string, string], children: list[map[string, string]]) -> list[string]
    + runs cross-struct rules by validating the parent and each child record together
    # validation
