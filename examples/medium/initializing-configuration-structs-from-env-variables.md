# Requirement: "a configuration loader that merges environment, files, flags, and defaults"

Builds a configuration map by layering sources in a declared precedence order.

std
  std.env
    std.env.lookup
      @ (name: string) -> optional[string]
      + returns the environment variable value
      - returns none when unset
      # environment
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when missing
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a json object into string-keyed string values
      - returns error on invalid json
      # serialization

config
  config.declare_field
    @ (name: string, default_value: string, env_name: string, flag_name: string) -> field_spec
    + builds a field specification with its default and source mappings
    # declaration
  config.schema
    @ (fields: list[field_spec]) -> schema
    + groups field specifications into a schema
    # declaration
  config.load_defaults
    @ (s: schema) -> map[string,string]
    + returns a map populated from each field's default value
    # sources
  config.load_env
    @ (s: schema) -> map[string,string]
    + returns a map populated from environment variables named in the schema
    # sources
    -> std.env.lookup
  config.load_file
    @ (s: schema, path: string) -> result[map[string,string], string]
    + parses a json file and returns only keys named in the schema
    # sources
    -> std.fs.read_all
    -> std.json.parse_object
  config.load_flags
    @ (s: schema, argv: list[string]) -> result[map[string,string], string]
    + extracts --name=value arguments matching schema flag names
    - returns error when a flag value is malformed
    # sources
  config.merge
    @ (layers: list[map[string,string]]) -> map[string,string]
    + merges layers with later layers overriding earlier ones
    # composition
  config.resolve
    @ (s: schema, argv: list[string], file_path: optional[string]) -> result[map[string,string], string]
    + returns the final configuration by layering defaults, file, env, flags
    # resolution
