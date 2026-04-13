# Requirement: "a data validation library with a schema-based rule system and extensible custom rules"

A declarative validator: declare a schema, register optional custom rules, validate objects and collect errors per field.

std: (all units exist)

validator
  validator.new
    @ () -> schema_state
    + creates an empty schema with no fields and no custom rules
    # construction
  validator.field
    @ (s: schema_state, name: string, rules: list[string]) -> schema_state
    + declares a field with the given named rules (e.g. "required", "min:3", "max:20", "regex:...")
    # declaration
  validator.register_rule
    @ (s: schema_state, name: string, rule_id: string) -> void
    + registers a custom rule handler so "name" can be used in field declarations
    # extensibility
  validator.validate
    @ (s: schema_state, obj: map[string, string]) -> map[string, list[string]]
    + returns a map of field_name -> error messages for every failing rule
    + returns an empty map when the object satisfies every rule
    # validation
  validator.validate_strict
    @ (s: schema_state, obj: map[string, string]) -> result[void, map[string, list[string]]]
    + returns ok when valid; otherwise returns the same error map as validate
    - rejects unknown keys present in obj but absent from the schema
    # validation
