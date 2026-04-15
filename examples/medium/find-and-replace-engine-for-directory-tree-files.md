# Requirement: "a find-and-replace engine for files in a directory tree"

Engine layer for an interactive find-and-replace tool. The caller handles the terminal; this library finds matches, previews replacements, and applies them.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all regular file paths under root
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at path
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path
      # filesystem
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_value, string]
      + compiles a regular expression
      - returns error on invalid syntax
      # regex
    std.regex.find_all
      fn (re: regex_value, text: string) -> list[match_span]
      + returns byte ranges for every non-overlapping match
      # regex
    std.regex.replace_all
      fn (re: regex_value, text: string, replacement: string) -> string
      + returns text with every match replaced
      # regex

find_replace
  find_replace.search
    fn (root: string, pattern: string, file_glob: string) -> result[list[match_location], string]
    + returns every match across files whose names match file_glob
    - returns error when the pattern does not compile
    # search
    -> std.fs.walk
    -> std.fs.read_all
    -> std.regex.compile
    -> std.regex.find_all
  find_replace.preview
    fn (location: match_location, replacement: string) -> preview_line
    + returns a before/after pair centered on the match for display
    # preview
  find_replace.apply
    fn (locations: list[match_location], pattern: string, replacement: string) -> result[i32, string]
    + rewrites the affected files and returns how many replacements were made
    - returns error when any target file cannot be written
    ? files touched by multiple locations are read and written once
    # application
    -> std.regex.compile
    -> std.regex.replace_all
    -> std.fs.read_all
    -> std.fs.write_all
