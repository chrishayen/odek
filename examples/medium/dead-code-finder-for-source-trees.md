# Requirement: "a library for finding unused code in a source tree"

Parses source files, collects every definition and every reference, and reports definitions that are never referenced.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root is not a directory
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the entire file contents as a string
      # filesystem

dead_code
  dead_code.tokenize
    fn (source: string) -> list[token]
    + returns identifier, keyword, punctuation, and string tokens with line numbers
    ? comments and whitespace are dropped
    # lexing
  dead_code.collect_definitions
    fn (path: string, tokens: list[token]) -> list[definition]
    + returns every function, class, and top-level name defined in the file, with location
    # analysis
  dead_code.collect_references
    fn (tokens: list[token]) -> list[string]
    + returns every identifier that appears in a reference position (not a definition)
    # analysis
  dead_code.scan_tree
    fn (root: string) -> result[scan_result, string]
    + walks the tree, tokenizes every source file, and returns all definitions and references
    - returns error when root is missing
    # orchestration
    -> std.fs.walk
    -> std.fs.read_all
    -> dead_code.tokenize
    -> dead_code.collect_definitions
    -> dead_code.collect_references
  dead_code.find_unused
    fn (scan: scan_result, exported_names: list[string]) -> list[definition]
    + returns every definition whose name never appears in the reference set and is not in exported_names
    ? exported_names lets the caller preserve public API entry points
    # detection
  dead_code.format_report
    fn (unused: list[definition]) -> string
    + returns a human-readable report grouped by file with line numbers
    # reporting
