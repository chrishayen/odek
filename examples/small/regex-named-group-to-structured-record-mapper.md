# Requirement: "a library for mapping regex named groups into a structured record using field tags"

Matches a compiled regex against input and returns a map keyed by the named group. A tag-driven binder maps the map into a caller-supplied record descriptor.

std
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_state, string]
      + compiles a pattern with named capture groups
      - returns error on invalid pattern
      # regex
    std.regex.find_named
      fn (state: regex_state, input: string) -> optional[map[string,string]]
      + returns a map of group_name to captured text on match
      - returns none when the pattern does not match
      # regex

regroup
  regroup.match_into
    fn (pattern: string, input: string) -> result[map[string,string], string]
    + returns a map of named group to captured substring
    - returns error when the pattern fails to compile
    - returns error when the input does not match
    # extraction
    -> std.regex.compile
    -> std.regex.find_named
  regroup.bind_fields
    fn (captures: map[string,string], field_tags: map[string,string]) -> result[map[string,string], string]
    + returns a record map keyed by field_name drawn from field_tags mapping field_name to group_name
    - returns error when a required tag references a missing group
    ? field_tags describes a struct-like record; the binding is value-level only
    # binding
