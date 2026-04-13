# Requirement: "a tag-based environment variable configuration loader"

Given a schema of tagged fields, reads values from environment variables and parses them by declared type. Intentionally narrower than a multi-provider loader.

std
  std.env
    std.env.get
      @ (key: string) -> optional[string]
      + returns the environment variable value when set
      # env
  std.strconv
    std.strconv.parse_i64
      @ (s: string) -> result[i64, string]
      + parses a decimal integer
      - returns error on non-digit input
      # parsing
    std.strconv.parse_bool
      @ (s: string) -> result[bool, string]
      + parses true/false/1/0/yes/no
      - returns error on other input
      # parsing

envconf
  envconf.field
    @ (name: string, env_key: string, type_tag: string, default_value: optional[string]) -> field_spec
    + builds a single field spec
    ? type_tag is one of "string", "i64", "bool"
    # schema
  envconf.load
    @ (fields: list[field_spec]) -> result[map[string, typed_value], string]
    + reads env vars, falls back to defaults, and parses each value under its type tag
    - returns error when a required field has no value
    - returns error when a value fails to parse
    # loading
    -> std.env.get
    -> std.strconv.parse_i64
    -> std.strconv.parse_bool
  envconf.get_string
    @ (config: map[string, typed_value], name: string) -> result[string, string]
    + returns the string field
    - returns error when the field is missing or not a string
    # access
  envconf.get_i64
    @ (config: map[string, typed_value], name: string) -> result[i64, string]
    + returns the integer field
    - returns error when the field is missing or not an integer
    # access
