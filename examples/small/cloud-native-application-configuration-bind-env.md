# Requirement: "a library for binding environment variables to typed configuration structures"

Given a schema describing fields and types, read environment variables and populate a typed config.

std
  std.env
    std.env.get
      @ (key: string) -> optional[string]
      + returns the value when the variable is set
      - returns none when unset
      # environment

envconfig
  envconfig.define_field
    @ (name: string, env_key: string, kind: string, required: bool, default_value: string) -> field_spec
    + returns a field descriptor
    + kind is one of "string", "int", "bool", "float"
    # schema
  envconfig.load
    @ (schema: list[field_spec]) -> result[map[string, config_value], string]
    + returns a map of field name to typed value
    - returns error when a required field is missing
    - returns error when a value fails to parse as its kind
    # binding
    -> std.env.get
  envconfig.coerce
    @ (raw: string, kind: string) -> result[config_value, string]
    + parses a raw string into the typed value
    - returns error when the string does not match the kind
    # coercion
