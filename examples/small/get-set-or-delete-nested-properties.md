# Requirement: "a library for get, set, or delete nested properties of an environment-variable map using a dot path"

Environment variables are flat strings, but the library treats dot-separated keys as nested paths over a string map.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string by separator into segments
      + returns a single-element list when sep is not present
      # strings
    std.strings.join
      @ (parts: list[string], sep: string) -> string
      + joins segments with the given separator
      # strings

dotenv_props
  dotenv_props.get
    @ (env: map[string,string], path: string) -> optional[string]
    + returns the value stored under the joined path
    - returns none when no key matches the path
    # lookup
    -> std.strings.split
    -> std.strings.join
  dotenv_props.set
    @ (env: map[string,string], path: string, value: string) -> map[string,string]
    + returns a new map with the path assigned
    + overwrites any existing value at that path
    # mutation
    -> std.strings.split
    -> std.strings.join
  dotenv_props.delete
    @ (env: map[string,string], path: string) -> map[string,string]
    + returns a new map with the path removed
    + returns the input unchanged when the path was absent
    # mutation
    -> std.strings.split
    -> std.strings.join
  dotenv_props.has
    @ (env: map[string,string], path: string) -> bool
    + returns true when a value is stored under the path
    - returns false when absent
    # lookup
    -> std.strings.split
    -> std.strings.join
