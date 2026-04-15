# Requirement: "a recursive struct configuration loader supporting multiple formats and environment overlays"

Loads a tree of configuration values from several source formats and recursively applies environment overrides.

std
  std.env
    std.env.lookup
      fn (name: string) -> optional[string]
      + returns the environment variable value
      - returns none when unset
      # environment
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when missing
      # filesystem
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary json into a tree
      - returns error on invalid json
      # serialization

swap
  swap.parse_json
    fn (raw: string) -> result[config_node, string]
    + parses json text into a configuration tree
    # parsing
    -> std.json.parse_value
  swap.parse_toml
    fn (raw: string) -> result[config_node, string]
    + parses toml text into a configuration tree
    - returns error on syntax errors
    # parsing
  swap.parse_yaml
    fn (raw: string) -> result[config_node, string]
    + parses yaml text into a configuration tree
    - returns error on syntax errors
    # parsing
  swap.load
    fn (path: string) -> result[config_node, string]
    + reads the file and dispatches to the parser chosen by extension
    - returns error when the extension is not recognized
    # loading
    -> std.fs.read_all
  swap.get_path
    fn (node: config_node, path: list[string]) -> optional[config_node]
    + walks a dotted path through the tree
    - returns none when any segment is missing
    # query
  swap.set_path
    fn (node: config_node, path: list[string], value: config_node) -> config_node
    + replaces or inserts the subtree at the given path
    # write
  swap.apply_env_overrides
    fn (node: config_node, prefix: string) -> config_node
    + recursively overlays environment variables whose names match "PREFIX_PATH"
    # overrides
    -> std.env.lookup
  swap.merge
    fn (base: config_node, overlay: config_node) -> config_node
    + recursively overlays one tree on another
    + overlay values replace scalars and extend maps
    # composition
