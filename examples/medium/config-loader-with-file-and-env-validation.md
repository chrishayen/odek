# Requirement: "a configuration loader that reads from a file and from environment variables, with validation and defaults"

Loads a string-to-string config map from a file, overlays environment variables, applies declared defaults, and runs required-key validation.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the entire file contents as a string
      - returns error when the file is missing or unreadable
      # filesystem
  std.env
    std.env.lookup
      fn (name: string) -> optional[string]
      + returns the value of the environment variable when set
      - returns none when not set
      # environment

config_loader
  config_loader.parse_file
    fn (path: string) -> result[map[string, string], string]
    + parses a "key = value" file into a map, ignoring blank lines and lines starting with '#'
    - returns error when a non-comment line lacks '='
    # parsing
    -> std.fs.read_all
  config_loader.overlay_env
    fn (values: map[string, string], prefix: string) -> map[string, string]
    + for each key, if an environment variable "PREFIX_KEY" exists (uppercased), it replaces the file value
    ? caller supplies the prefix; empty prefix disables the namespace
    # override
    -> std.env.lookup
  config_loader.apply_defaults
    fn (values: map[string, string], defaults: map[string, string]) -> map[string, string]
    + fills in any key from defaults that is absent from values
    + leaves existing values untouched
    # defaults
  config_loader.validate_required
    fn (values: map[string, string], required: list[string]) -> result[void, string]
    + returns ok when every required key is present and non-empty
    - returns error naming the first missing or empty required key
    # validation
  config_loader.load
    fn (path: string, prefix: string, defaults: map[string, string], required: list[string]) -> result[map[string, string], string]
    + parses the file, overlays env, applies defaults, then validates required keys
    - returns error from whichever step fails first
    # orchestration
