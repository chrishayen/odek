# Requirement: "a library that loads configuration from files, environment variables, and flags"

Layers three configuration sources with flags winning over environment winning over file.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when missing
      # filesystem
  std.env
    std.env.get
      fn (name: string) -> optional[string]
      + returns the value when set
      # environment

layered_config
  layered_config.parse_file
    fn (data: bytes) -> result[map[string,string], string]
    + parses key=value lines into a map, ignoring comments and blank lines
    - returns error on a line with no equals sign
    # parsing
  layered_config.load_file
    fn (path: string) -> result[map[string,string], string]
    + reads and parses a configuration file
    - returns error when the file is missing
    # loading
    -> std.fs.read_all
  layered_config.overlay_env
    fn (base: map[string,string], prefix: string) -> map[string,string]
    + returns base with values from environment variables matching prefix_key
    ? lowercases the trailing portion and replaces underscores with dots
    # layering
    -> std.env.get
  layered_config.parse_flags
    fn (argv: list[string]) -> map[string,string]
    + parses --key=value and --key value forms into a map
    # parsing
  layered_config.merge
    fn (file: map[string,string], env: map[string,string], flags: map[string,string]) -> map[string,string]
    + returns the union with flags winning over env winning over file
    # layering
