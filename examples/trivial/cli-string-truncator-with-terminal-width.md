# Requirement: "a library for truncating a string to a specific display width in the terminal"

A single function that truncates by visual column width, accounting for wide characters and ANSI escape sequences.

std: (all units exist)

cli_truncate
  cli_truncate.truncate
    fn (text: string, columns: i32) -> string
    + returns text unchanged when its display width is at most columns
    + returns a truncated prefix ending in an ellipsis when text is wider than columns
    + counts east-asian wide characters as 2 columns and skips ANSI escape sequences
    - returns an empty string when columns is 0 or negative
    # truncation
