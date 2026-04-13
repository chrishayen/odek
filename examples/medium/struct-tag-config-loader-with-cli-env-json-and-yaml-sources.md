# Requirement: "a struct-tag based configuration loader sourcing from command-line arguments, environment variables, JSON, and YAML"

Takes a target struct with tagged fields and populates it from a layered set of sources with a defined precedence.

std
  std.env
    std.env.lookup
      @ (name: string) -> optional[string]
      + returns the value of the environment variable, or none
      # environment
  std.fs
    std.fs.read_text
      @ (path: string) -> result[string, string]
      + returns the file's contents as text
      - returns error when the file does not exist
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a flat string-to-string map
      - returns error on malformed input or non-object root
      # serialization
  std.yaml
    std.yaml.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a YAML mapping into a flat string-to-string map
      - returns error on malformed input or non-mapping root
      # serialization

config_loader
  config_loader.describe
    @ (schema: list[field_spec]) -> schema_state
    + returns a schema_state capturing field names, types, tag keys, and defaults
    # schema
  config_loader.parse_args
    @ (schema: schema_state, argv: list[string]) -> result[map[string, string], string]
    + returns values parsed from POSIX long and short options
    - returns error when a required field has no value after parsing
    # source_cli
  config_loader.from_env
    @ (schema: schema_state, prefix: string) -> map[string, string]
    + returns values read from environment variables matching the schema's env tags
    # source_env
    -> std.env.lookup
  config_loader.from_json
    @ (schema: schema_state, path: string) -> result[map[string, string], string]
    + returns values read from a JSON config file
    - returns error when the file is missing or malformed
    # source_json
    -> std.fs.read_text
    -> std.json.parse_object
  config_loader.from_yaml
    @ (schema: schema_state, path: string) -> result[map[string, string], string]
    + returns values read from a YAML config file
    - returns error when the file is missing or malformed
    # source_yaml
    -> std.fs.read_text
    -> std.yaml.parse_object
  config_loader.merge
    @ (layers: list[map[string, string]]) -> map[string, string]
    + merges layers so later entries override earlier ones
    ? intended order: defaults, file, env, cli
    # merging
  config_loader.coerce
    @ (schema: schema_state, raw: map[string, string]) -> result[map[string, config_value], string]
    + converts each raw string into the typed value declared by the schema
    - returns error when a field's value cannot be parsed as its declared type
    # coercion
