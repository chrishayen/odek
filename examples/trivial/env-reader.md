# Requirement: "a library for reading environment variables"

Thin typed wrappers over a process environment lookup primitive.

std
  std.os
    std.os.get_env
      fn (name: string) -> optional[string]
      + returns the value of the environment variable, or none when unset
      # environment

envread
  envread.string_or
    fn (name: string, default_value: string) -> string
    + returns the env var value when set, otherwise default_value
    # lookup
    -> std.os.get_env
  envread.int_or
    fn (name: string, default_value: i64) -> result[i64, string]
    + parses the env var as a signed integer when set
    - returns error when the value is set but not numeric
    ? when unset, returns ok(default_value)
    # lookup
    -> std.os.get_env
  envread.bool_or
    fn (name: string, default_value: bool) -> result[bool, string]
    + accepts "true"/"false"/"1"/"0" (case insensitive)
    - returns error on any other value
    # lookup
    -> std.os.get_env
