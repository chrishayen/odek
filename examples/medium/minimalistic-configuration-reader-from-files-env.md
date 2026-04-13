# Requirement: "a configuration reader that merges values from files and environment variables"

Configuration is represented as a flat string map. The library loads layers and merges them with a defined precedence.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the text contents of a file
      - returns error when the path does not exist
      # filesystem
  std.env
    std.env.get
      @ (name: string) -> optional[string]
      + returns the value of the named environment variable
      - returns empty when unset
      # environment
    std.env.list_with_prefix
      @ (prefix: string) -> list[tuple[string, string]]
      + returns all (name, value) pairs whose name starts with prefix
      # environment

config
  config.parse_kv_text
    @ (raw: string) -> result[map[string, string], string]
    + parses "key=value" lines into a map, ignoring blank lines and # comments
    - returns error on lines that contain no '='
    # parsing
  config.load_file
    @ (path: string) -> result[map[string, string], string]
    + reads and parses a kv file
    - returns error when the file is missing or malformed
    # loading
    -> std.fs.read_all
  config.load_env
    @ (prefix: string) -> map[string, string]
    + collects env vars with the given prefix, lowercasing keys with the prefix stripped
    # loading
    -> std.env.list_with_prefix
  config.merge
    @ (base: map[string, string], overlay: map[string, string]) -> map[string, string]
    + returns a map where overlay keys override base keys
    # merging
  config.get_required
    @ (cfg: map[string, string], key: string) -> result[string, string]
    + returns the value when present
    - returns error naming the missing key when absent
    # lookup
