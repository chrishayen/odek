# Requirement: "a library that finds and removes unnecessary type conversions from source code"

Walks a parsed syntax tree, identifies conversion expressions whose target type already matches the argument's type, and emits a patched source.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file as text
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes the contents to the file, creating it if needed
      # filesystem

unconvert
  unconvert.parse_source
    @ (source: string) -> result[syntax_tree, string]
    + parses the source into a typed syntax tree
    - returns error on syntactic failure
    # parsing
  unconvert.infer_types
    @ (tree: syntax_tree) -> result[typed_tree, string]
    + resolves the type of every expression in the tree
    - returns error on unresolved symbols
    # type_inference
  unconvert.find_redundant
    @ (tree: typed_tree) -> list[source_span]
    + returns the spans of conversion expressions whose target type equals the argument type
    # analysis
  unconvert.rewrite
    @ (source: string, spans: list[source_span]) -> string
    + returns the source with each redundant conversion replaced by its inner expression
    + preserves surrounding whitespace and comments
    # rewriting
  unconvert.process_file
    @ (path: string) -> result[i32, string]
    + rewrites the file in place and returns the number of conversions removed
    - returns error when parsing or type inference fails
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
