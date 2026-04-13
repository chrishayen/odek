# Requirement: "a library for loading environment variables into structured records"

Reads named environment variables, parses them into typed values by a schema, and returns a populated record.

std
  std.env
    std.env.get
      @ (name: string) -> optional[string]
      + returns the value of the environment variable, or none when unset
      # environment

env_loader
  env_loader.define_schema
    @ () -> schema_state
    + returns an empty schema
    # construction
  env_loader.add_string
    @ (schema: schema_state, field: string, var_name: string, default_value: optional[string], required: bool) -> schema_state
    + registers a string-valued field bound to an environment variable
    # schema
  env_loader.add_int
    @ (schema: schema_state, field: string, var_name: string, default_value: optional[i64], required: bool) -> schema_state
    + registers an integer-valued field
    # schema
  env_loader.add_bool
    @ (schema: schema_state, field: string, var_name: string, default_value: optional[bool], required: bool) -> schema_state
    + registers a boolean-valued field
    # schema
  env_loader.load
    @ (schema: schema_state) -> result[map[string, string], string]
    + returns a map of field-to-stringified-value populated from the environment
    - returns error listing every required field that was missing
    - returns error when an int or bool field fails to parse
    # load
    -> std.env.get
