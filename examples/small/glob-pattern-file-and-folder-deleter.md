# Requirement: "delete files and folders selected by glob patterns"

The project layer expands patterns, applies ignore rules, and deletes the resolved paths. Filesystem primitives live in std.

std
  std.fs
    std.fs.glob
      fn (root: string, pattern: string) -> result[list[string], string]
      + returns absolute paths under root matching the glob
      - returns error on an invalid pattern
      # filesystem
    std.fs.remove
      fn (path: string) -> result[void, string]
      + deletes a file or empty directory
      - returns error when the path does not exist
      # filesystem
    std.fs.remove_tree
      fn (path: string) -> result[void, string]
      + deletes a directory and all its contents
      - returns error when the path does not exist
      # filesystem
    std.fs.is_dir
      fn (path: string) -> bool
      + returns true when the path is a directory
      # filesystem

glob_delete
  glob_delete.resolve
    fn (root: string, patterns: list[string], ignore: list[string]) -> result[list[string], string]
    + returns the set of paths matching any pattern and no ignore pattern
    - returns error on an invalid pattern
    # matching
    -> std.fs.glob
  glob_delete.delete_all
    fn (paths: list[string]) -> result[i32, string]
    + deletes each path, recursing into directories
    + returns the count of successfully deleted entries
    - returns error on the first failing deletion and reports progress so far via the error message
    # deletion
    -> std.fs.is_dir
    -> std.fs.remove
    -> std.fs.remove_tree
