# Requirement: "a simple opinionated config loader"

Loads config from a key=value source and exposes typed getters.

std: (all units exist)

config
  config.parse
    fn (source: string) -> result[map[string, string], string]
    + parses lines of the form `key=value`, ignoring blank lines and `#` comments
    - returns error on a line missing `=`
    # parsing
  config.get_string
    fn (cfg: map[string, string], key: string, default_value: string) -> string
    + returns the value for key, or default_value when the key is absent
    # lookup
  config.get_int
    fn (cfg: map[string, string], key: string, default_value: i64) -> result[i64, string]
    + parses the value as a signed integer
    - returns error when the value exists but is not numeric
    # lookup
  config.get_bool
    fn (cfg: map[string, string], key: string, default_value: bool) -> result[bool, string]
    + accepts "true", "false", "1", "0", "yes", "no" (case insensitive)
    - returns error on any other string
    # lookup
