# Requirement: "a hierarchical configuration loader with inheritance and key generation"

Loads a configuration with parent-child inheritance: a child config may extend a parent, overriding individual keys while inheriting the rest. Also generates canonical dotted keys for nested values.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when path does not exist
      # io
  std.collections
    std.collections.map_merge
      @ (base: map[string, config_value], overlay: map[string, config_value]) -> map[string, config_value]
      + returns a new map with overlay entries taking precedence
      # collections

config_loader
  config_loader.parse
    @ (raw: string) -> result[map[string, config_value], string]
    + parses a nested key-value document into a map of scalar and nested values
    - returns error on malformed syntax
    # parsing
  config_loader.load_file
    @ (path: string) -> result[map[string, config_value], string]
    + reads and parses a configuration file
    - propagates fs errors
    # loading
    -> std.fs.read_all
  config_loader.resolve_inheritance
    @ (child: map[string, config_value], parent: map[string, config_value]) -> map[string, config_value]
    + recursively merges child over parent so nested child keys override parent keys
    + inherits any key not present in child
    # inheritance
    -> std.collections.map_merge
  config_loader.flatten_keys
    @ (config: map[string, config_value]) -> map[string, config_value]
    + rewrites nested keys as dotted paths ("db.host", "db.port")
    + leaves scalar roots unchanged
    # key_generation
  config_loader.get
    @ (config: map[string, config_value], key: string) -> optional[config_value]
    + looks up a value by dotted key
    - returns none when key is absent
    # querying
