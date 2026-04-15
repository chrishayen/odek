# Requirement: "a library for expanding multiple glob patterns against a filesystem, with include and exclude support"

Compiles glob patterns into matchers, walks a directory tree, and returns files that match any include pattern and no exclude pattern.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[dir_entry], string]
      + returns the entries under path
      - returns error when path is not a directory
      # filesystem
    std.fs.is_dir
      fn (path: string) -> bool
      + returns true when path is a directory
      # filesystem
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits a string by separator
      # strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings

glob
  glob.compile
    fn (pattern: string) -> result[glob_matcher, string]
    + compiles a pattern supporting '*', '?', '**', and character classes
    - returns error on unbalanced brackets
    # compilation
    -> std.strings.split
  glob.matches
    fn (matcher: glob_matcher, path: string) -> bool
    + returns true when the relative path matches the compiled pattern
    - returns false otherwise
    # matching
    -> std.strings.split
  glob.split_patterns
    fn (patterns: list[string]) -> tuple[list[string], list[string]]
    + partitions into include and exclude lists, where exclude patterns start with '!'
    # parsing
    -> std.strings.starts_with
  glob.walk
    fn (root: string) -> result[list[string], string]
    + recursively lists every file path under root
    - returns error when root cannot be read
    # filesystem
    -> std.fs.list_dir
    -> std.fs.is_dir
  glob.expand
    fn (root: string, patterns: list[string]) -> result[list[string], string]
    + returns all files under root that match any include and no exclude, in stable order
    - returns error when root cannot be walked or any pattern fails to compile
    # orchestration
