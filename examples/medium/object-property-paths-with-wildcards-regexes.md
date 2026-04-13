# Requirement: "object property paths with wildcards and regex segments"

A path language that matches nested fields in a dynamic value tree, supporting `*`, `**`, and `/regex/` segments.

std
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[regex_handle, string]
      + returns a compiled regex
      - returns error on invalid pattern
      # regex
    std.regex.matches
      @ (re: regex_handle, s: string) -> bool
      + returns true when the string matches the regex
      # regex

path_query
  path_query.parse
    @ (path: string) -> result[list[path_segment], string]
    + parses dot-separated segments recognizing "*", "**", "/re/", and literal keys
    - returns error on unmatched regex delimiters
    - returns error on empty path
    # parsing
  path_query.segment_matches_key
    @ (seg: path_segment, key: string) -> bool
    + literal segments match by equality
    + "*" matches any single key
    + regex segments match keys that satisfy the compiled pattern
    # matching
    -> std.regex.matches
  path_query.list_paths
    @ (segments: list[path_segment], root: dynamic_value) -> list[list[string]]
    + returns all concrete key paths in root that satisfy the segments
    + "**" expands to zero or more intermediate keys
    # expansion
  path_query.get
    @ (segments: list[path_segment], root: dynamic_value) -> list[dynamic_value]
    + returns all values found at matching paths
    - returns an empty list when nothing matches
    # retrieval
  path_query.set
    @ (segments: list[path_segment], root: dynamic_value, value: dynamic_value) -> dynamic_value
    + returns a new root with value assigned at every matching path
    ? creates missing intermediate objects only for literal segments
    # mutation
  path_query.delete
    @ (segments: list[path_segment], root: dynamic_value) -> dynamic_value
    + returns a new root with every matching path removed
    # mutation
