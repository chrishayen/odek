# Requirement: "load environment variables into a typed configuration with defaults"

A declarative config loader: the caller builds a schema describing each field, the loader reads process environment values, applies defaults, coerces types, and returns a typed map or a list of errors.

std
  std.env
    std.env.lookup
      @ (name: string) -> optional[string]
      + returns the value of the environment variable when set
      - returns absent when the variable is unset
      # environment

envcfg
  envcfg.new_schema
    @ () -> env_schema
    + returns an empty schema
    # construction
  envcfg.field_string
    @ (schema: env_schema, name: string, default_value: optional[string], required: bool) -> env_schema
    + returns a schema with a string field appended
    # schema
  envcfg.field_int
    @ (schema: env_schema, name: string, default_value: optional[i64], required: bool) -> env_schema
    + returns a schema with an integer field appended
    # schema
  envcfg.field_bool
    @ (schema: env_schema, name: string, default_value: optional[bool], required: bool) -> env_schema
    + returns a schema with a boolean field appended
    # schema
  envcfg.field_list
    @ (schema: env_schema, name: string, separator: string, default_value: optional[list[string]], required: bool) -> env_schema
    + returns a schema with a list field parsed from a separator-delimited string
    # schema
  envcfg.load
    @ (schema: env_schema) -> result[env_values, list[string]]
    + returns typed values when every field is present or defaulted
    - returns the list of per-field error messages when any required field is missing or mistyped
    # loading
    -> std.env.lookup
  envcfg.load_from
    @ (schema: env_schema, source: map[string, string]) -> result[env_values, list[string]]
    + returns typed values using source instead of the process environment
    - returns error list on missing or mistyped fields
    ? accepts an explicit source so tests can inject a fake environment
    # loading
  envcfg.get_string
    @ (values: env_values, name: string) -> result[string, string]
    + returns the string value for the named field
    - returns error when the field is not a string or is not present
    # access
  envcfg.get_int
    @ (values: env_values, name: string) -> result[i64, string]
    + returns the integer value for the named field
    - returns error when the field is not an integer or is not present
    # access
  envcfg.get_bool
    @ (values: env_values, name: string) -> result[bool, string]
    + returns the boolean value for the named field
    - returns error when the field is not a boolean or is not present
    # access
