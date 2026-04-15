# Requirement: "a library for finding and removing large dependency directories"

Scans a root for directories matching a target name, reports their sizes and last-modified times, and exposes a deletion primitive.

std
  std.fs
    std.fs.walk_dirs
      fn (root: string) -> result[list[string], string]
      + returns directory paths under root, recursively
      - returns error when root does not exist
      # filesystem
    std.fs.dir_size
      fn (path: string) -> result[i64, string]
      + returns the total size in bytes of a directory tree
      # filesystem
    std.fs.mtime
      fn (path: string) -> result[i64, string]
      + returns the last modified unix time in seconds
      # filesystem
    std.fs.remove_all
      fn (path: string) -> result[void, string]
      + recursively removes a directory and its contents
      - returns error when the path cannot be removed
      # filesystem

dep_sweeper
  dep_sweeper.find
    fn (root: string, target_name: string) -> result[list[dep_entry], string]
    + returns every directory under root whose basename equals target_name, with size and mtime
    - returns error when the root cannot be scanned
    # discovery
    -> std.fs.walk_dirs
    -> std.fs.dir_size
    -> std.fs.mtime
  dep_sweeper.delete
    fn (entry: dep_entry) -> result[i64, string]
    + removes the directory and returns the number of bytes reclaimed
    - returns error when the deletion fails
    # deletion
    -> std.fs.remove_all
