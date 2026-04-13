# Requirement: "a URL pattern matching library that is simpler than raw regex"

Patterns use named segments like ":name" and wildcards "*" and match against input paths, returning extracted values. Compilation is separated from matching so patterns can be reused.

std: (all units exist)

url_pattern
  url_pattern.compile
    @ (pattern: string) -> result[compiled_pattern, string]
    + compiles a pattern with literal, named, and wildcard segments
    - returns error on unbalanced or empty named segments
    # compilation
  url_pattern.match
    @ (compiled: compiled_pattern, input: string) -> optional[map[string, string]]
    + returns the map of extracted named segments when the input matches
    + wildcard segments capture the remaining path under the key "_"
    - returns none when the input does not match
    # matching
  url_pattern.expand
    @ (compiled: compiled_pattern, values: map[string, string]) -> result[string, string]
    + renders a path by substituting named segments with values
    - returns error when a required name is missing from values
    # rendering
  url_pattern.names
    @ (compiled: compiled_pattern) -> list[string]
    + returns the ordered list of named segments in the pattern
    # introspection
