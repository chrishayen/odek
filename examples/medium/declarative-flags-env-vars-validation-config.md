# Requirement: "a declarative flags, env vars, validation, and config file loader driven by struct-tag-style schemas"

The project layer resolves each field from multiple sources in priority order (flags > env > file > default) and then validates. Parsing primitives live in std.

std
  std.encoding
    std.encoding.parse_toml
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat key/value config document into a string-to-string map
      - returns error on malformed input
      # serialization
  std.env
    std.env.lookup
      @ (name: string) -> optional[string]
      + returns the value when the environment variable is set
      - returns none when it is not set
      # environment
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as a string
      - returns error when the file does not exist
      # filesystem

config_loader
  config_loader.new_schema
    @ () -> schema_state
    + creates an empty schema accumulator
    # construction
  config_loader.declare_field
    @ (schema: schema_state, name: string, kind: string, default_value: optional[string], required: bool) -> schema_state
    + registers a field with its type name, default, and required flag
    ? kind is one of "string", "int", "bool", "float"
    # schema
  config_loader.load
    @ (schema: schema_state, argv: list[string], config_path: optional[string]) -> result[map[string, string], string]
    + resolves values using priority argv > env > file > default
    + env var names are the uppercased field name
    - returns error when a required field has no value in any source
    - returns error when a value cannot be parsed to its declared kind
    # resolution
    -> std.env.lookup
    -> std.fs.read_all
    -> std.encoding.parse_toml
  config_loader.get_string
    @ (values: map[string, string], name: string) -> result[string, string]
    + returns the resolved string value for a field
    - returns error when the field is missing
    # access
  config_loader.get_int
    @ (values: map[string, string], name: string) -> result[i64, string]
    + returns the resolved integer value for a field
    - returns error when the field is missing or not an integer
    # access
