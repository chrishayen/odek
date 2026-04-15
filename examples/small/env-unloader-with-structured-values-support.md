# Requirement: "an environment utility library with support for unmarshaling into structured values"

Read environment variables and decode them into a typed field map. The caller supplies a schema describing expected field names and types.

std
  std.env
    std.env.lookup
      fn (name: string) -> optional[string]
      + returns the value when the variable is set
      - returns none when the variable is not set
      # environment

env_loader
  env_loader.field_spec
    fn (name: string, type_tag: string, required: bool, default_value: optional[string]) -> field_spec
    + builds a schema entry describing one expected variable
    # schema
  env_loader.load
    fn (specs: list[field_spec]) -> result[map[string, string], string]
    + returns a map of resolved values when all required fields are present
    + applies default values for unset optional fields
    - returns error listing missing required fields
    # loading
    -> std.env.lookup
  env_loader.decode_int
    fn (raw: string) -> result[i64, string]
    + parses a signed decimal integer
    - returns error on empty or non-numeric input
    # decoding
  env_loader.decode_bool
    fn (raw: string) -> result[bool, string]
    + accepts "1", "true", "yes" as true and "0", "false", "no" as false (case-insensitive)
    - returns error on unrecognized value
    # decoding
