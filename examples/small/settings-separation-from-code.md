# Requirement: "a library for strict separation of settings from code"

Loads configuration from environment variables or a key/value file, with typed accessors and defaults. Code never hardcodes settings.

std
  std.env
    std.env.lookup
      fn (name: string) -> optional[string]
      + returns the value of an environment variable if set
      # environment
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file fully into a string
      - returns error when the file does not exist
      # filesystem

settings
  settings.load_file
    fn (path: string) -> result[settings_state, string]
    + parses a key=value file (one pair per line, # comments) into a settings store
    - returns error on malformed lines
    # loading
    -> std.fs.read_all
  settings.new_empty
    fn () -> settings_state
    + creates an empty settings store that falls back to environment variables
    # construction
  settings.get_string
    fn (state: settings_state, key: string, default_value: optional[string]) -> result[string, string]
    + returns the value from the file, then from env, then default
    - returns error when key is absent everywhere and no default is given
    # access
    -> std.env.lookup
  settings.get_int
    fn (state: settings_state, key: string, default_value: optional[i64]) -> result[i64, string]
    + parses the resolved value as an integer
    - returns error when the value cannot be parsed as an integer
    # access
  settings.get_bool
    fn (state: settings_state, key: string, default_value: optional[bool]) -> result[bool, string]
    + accepts "1","true","yes","on" and their negatives (case-insensitive)
    - returns error when value is not a recognized boolean literal
    # access
