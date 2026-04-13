# Requirement: "a typed environment variable reader"

Looks up environment values and coerces them to concrete types. Returns rich errors that name the offending variable.

std
  std.env
    std.env.get
      @ (name: string) -> optional[string]
      + returns the value of the named environment variable
      - returns empty when unset
      # environment

envtyped
  envtyped.require_string
    @ (name: string) -> result[string, string]
    + returns the value when set
    - returns error "<name>: required" when unset
    # lookup
    -> std.env.get
  envtyped.require_i64
    @ (name: string) -> result[i64, string]
    + returns the parsed integer
    - returns error when unset
    - returns error "<name>: not an integer" when the value is non-numeric
    # lookup
    -> std.env.get
  envtyped.require_bool
    @ (name: string) -> result[bool, string]
    + accepts "1", "true", "yes" (case-insensitive) as true
    + accepts "0", "false", "no" as false
    - returns error when unset or unrecognized
    # lookup
    -> std.env.get
  envtyped.optional_string
    @ (name: string, default_val: string) -> string
    + returns the value when set, otherwise default_val
    # lookup
    -> std.env.get
