# Requirement: "a configuration framework supporting layered config files, overrides, and composition"

Configs are loaded from one or more files, merged in order, then overridden by command-line-style key=value pairs.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file is missing or unreadable
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[config_value, string]
      + parses JSON into a generic config value tree
      - returns error on invalid JSON
      # serialization

hydra
  hydra.load_file
    @ (path: string) -> result[config_value, string]
    + reads and parses a config file
    - returns error when the file cannot be read or parsed
    # loading
    -> std.fs.read_all
    -> std.json.parse
  hydra.merge
    @ (base: config_value, overlay: config_value) -> config_value
    + deep-merges overlay into base; overlay scalars replace base scalars
    + nested maps are merged recursively
    # composition
  hydra.apply_override
    @ (cfg: config_value, key_path: string, value: string) -> result[config_value, string]
    + applies a dotted-key override such as "db.host=localhost"
    - returns error when the key path traverses a non-map
    # override
  hydra.compose
    @ (paths: list[string], overrides: list[string]) -> result[config_value, string]
    + loads each file, merges in order, then applies overrides in sequence
    - returns error when any file fails to load
    # composition
  hydra.get_string
    @ (cfg: config_value, key_path: string) -> optional[string]
    + returns a string value at a dotted key path
    - returns none when the path is missing
    # access
  hydra.get_int
    @ (cfg: config_value, key_path: string) -> optional[i64]
    + returns an integer value at a dotted key path
    - returns none when the path is missing or not numeric
    # access
