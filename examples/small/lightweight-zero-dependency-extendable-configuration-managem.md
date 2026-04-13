# Requirement: "a library for extendable configuration management"

Loads configuration from ordered sources (defaults, file, environment, overrides), merges them, and reads typed values.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of the file as a string
      - returns error when the file cannot be read
      # filesystem
  std.env
    std.env.lookup
      @ (key: string) -> optional[string]
      + returns the value of the environment variable or none
      # env

config
  config.new
    @ () -> config_state
    + creates an empty config
    # construction
  config.load_file
    @ (state: config_state, path: string) -> result[config_state, string]
    + parses a key=value file and merges it on top of the current state
    - returns error on read failure or malformed lines
    # loading
    -> std.fs.read_all
  config.load_env
    @ (state: config_state, prefix: string) -> config_state
    + imports every environment variable whose name starts with prefix, stripping the prefix and lowercasing
    # loading
    -> std.env.lookup
  config.set
    @ (state: config_state, key: string, value: string) -> config_state
    + writes an override that wins over every loaded source
    # mutation
  config.get_string
    @ (state: config_state, key: string) -> result[string, string]
    + returns the resolved value for key
    - returns error when the key is not set
    # access
  config.get_int
    @ (state: config_state, key: string) -> result[i64, string]
    + parses the resolved value as a signed integer
    - returns error when the key is missing or not an integer
    # access
  config.get_bool
    @ (state: config_state, key: string) -> result[bool, string]
    + parses "true"/"false"/"1"/"0" case-insensitively
    - returns error when the key is missing or not a boolean literal
    # access
