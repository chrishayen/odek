# Requirement: "a library for testing JSON values against expected patterns"

Patterns are JSON-shaped with wildcards like "@string@", "@number@", "@uuid@" that match any value of that type. Comparison returns a human-readable diff on mismatch.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged value
      - returns error on malformed JSON
      # parsing
  std.regex
    std.regex.is_uuid
      fn (s: string) -> bool
      + true when the string matches the UUID v4 format
      # validation

json_match
  json_match.match
    fn (actual: string, pattern: string) -> result[void, string]
    + returns ok when the actual JSON matches the pattern
    - returns error with a path-prefixed diff on mismatch
    - returns error when either input is not valid JSON
    # matching
    -> std.json.parse
  json_match.match_value
    fn (actual: json_value, pattern: json_value, path: string) -> result[void, string]
    + recursively compares values, honoring pattern wildcards
    - returns a path-prefixed error on the first mismatch
    # matching
  json_match.is_wildcard
    fn (pattern: json_value) -> optional[string]
    + returns the wildcard name when the pattern is a wildcard string literal
    ? wildcard literals have the form "@name@"
    # matching
  json_match.check_wildcard
    fn (name: string, actual: json_value) -> result[void, string]
    + checks that the actual value satisfies the named wildcard kind
    - returns error when the type does not match the wildcard
    # matching
    -> std.regex.is_uuid
