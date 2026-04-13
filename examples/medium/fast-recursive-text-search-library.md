# Requirement: "a fast recursive text search library"

Searches file trees for lines matching a pattern. File walking and pattern matching are the real primitives; the project layer composes them.

std
  std.fs
    std.fs.walk_files
      @ (root: string) -> result[list[string], string]
      + yields every regular file path under root recursively
      - returns error when root does not exist
      # filesystem
    std.fs.read_lines
      @ (path: string) -> result[list[string], string]
      + returns each line of a file without terminators
      - returns error when the file cannot be opened
      # filesystem
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[compiled_regex, string]
      + compiles a pattern for repeated matching
      - returns error on invalid syntax
      # pattern_matching
    std.regex.is_match
      @ (re: compiled_regex, input: string) -> bool
      + returns true when the input contains a match
      # pattern_matching

grep
  grep.search_file
    @ (path: string, re: compiled_regex) -> result[list[search_hit], string]
    + returns one hit per matching line with 1-based line number
    - returns error when the file cannot be read
    # file_search
    -> std.fs.read_lines
    -> std.regex.is_match
  grep.search_tree
    @ (root: string, pattern: string) -> result[list[search_hit], string]
    + walks root and returns every hit across all files
    - returns error when the pattern fails to compile
    ? binary detection and ignore files are out of scope for this core
    # recursive_search
    -> std.fs.walk_files
    -> std.regex.compile
