# Requirement: "helpers to manage environment variables"

Typed accessors over the process environment with default fallbacks and prefix filtering.

std
  std.env
    std.env.get
      @ (name: string) -> optional[string]
      + returns the value of the environment variable, or none if unset
      # environment
    std.env.all
      @ () -> map[string, string]
      + returns a snapshot of all environment variables
      # environment

envh
  envh.get_string
    @ (name: string, default_value: string) -> string
    + returns the variable value, or default_value when unset or empty
    # accessors
    -> std.env.get
  envh.get_int
    @ (name: string, default_value: i64) -> i64
    + returns the variable parsed as an integer, or default_value when unset or unparsable
    # accessors
    -> std.env.get
  envh.get_bool
    @ (name: string, default_value: bool) -> bool
    + returns true for "1", "true", "yes"; false for "0", "false", "no"
    - returns default_value when the value is not recognized
    # accessors
    -> std.env.get
  envh.require
    @ (name: string) -> result[string, string]
    + returns the value when set and non-empty
    - returns error naming the variable when unset or empty
    # accessors
    -> std.env.get
  envh.with_prefix
    @ (prefix: string) -> map[string, string]
    + returns all variables whose names begin with prefix, with the prefix stripped from keys
    # filtering
    -> std.env.all
