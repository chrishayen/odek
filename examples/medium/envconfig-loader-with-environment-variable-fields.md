# Requirement: "a configuration loader that reads fields from environment variables"

Describe a schema of expected fields with types and defaults, then fill a typed config map from the environment.

std
  std.env
    std.env.get
      fn (key: string) -> optional[string]
      + returns the value of an environment variable
      - returns none when unset
      # environment

envconfig
  envconfig.field_string
    fn (key: string, default_value: optional[string], required: bool) -> env_field
    + describes a string-valued field
    # schema
  envconfig.field_int
    fn (key: string, default_value: optional[i64], required: bool) -> env_field
    + describes an integer-valued field
    # schema
  envconfig.field_bool
    fn (key: string, default_value: optional[bool], required: bool) -> env_field
    + describes a boolean-valued field, accepting true/false/1/0/yes/no
    # schema
  envconfig.load
    fn (fields: list[env_field]) -> result[map[string, config_value], list[string]]
    + returns a map with one entry per field, using env value or default
    - returns error list containing every required field that is missing
    - returns error list for any field whose value fails to parse to its declared type
    # loading
    -> std.env.get
