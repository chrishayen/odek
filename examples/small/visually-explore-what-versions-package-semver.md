# Requirement: "a semver range matcher that lists which versions satisfy a range"

Given a list of available versions and a semver range expression, return the versions that match.

std
  std.sort
    std.sort.sort_strings
      @ (items: list[string]) -> list[string]
      + returns items sorted ascending
      # sorting

semver_match
  semver_match.parse_version
    @ (text: string) -> result[list[i32], string]
    + parses "1.2.3" into [1, 2, 3]
    - returns error when a segment is not a non-negative integer
    # parsing
  semver_match.parse_range
    @ (text: string) -> result[semver_range, string]
    + parses expressions with operators like ^, ~, >=, <, and ||
    - returns error on unknown operators
    # parsing
  semver_match.satisfies
    @ (version: list[i32], range: semver_range) -> bool
    + returns true when the version falls inside the range
    - returns false when outside any constraint
    # matching
  semver_match.matching_versions
    @ (versions: list[string], range_text: string) -> result[list[string], string]
    + returns versions from the input that match the range, sorted ascending
    - returns error when the range text is invalid
    # query
    -> std.sort.sort_strings
