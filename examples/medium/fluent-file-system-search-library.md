# Requirement: "a fluent file system search library"

A small builder that accumulates filters then walks a root directory, returning matches. Filesystem reads go through std primitives.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns entry names (files and subdirectories) of a directory
      - returns error when the path does not exist or is not a directory
      # filesystem
    std.fs.stat
      fn (path: string) -> result[file_info, string]
      + returns size, mtime, and is_dir for a path
      - returns error on missing path
      # filesystem
  std.path
    std.path.join
      fn (a: string, b: string) -> string
      + joins two path segments with the platform separator
      # path
    std.path.extension
      fn (path: string) -> string
      + returns the file extension including the dot, or "" if none
      # path

file_search
  file_search.new
    fn (root: string) -> search_query
    + creates an empty query rooted at the given directory
    # construction
  file_search.with_extension
    fn (q: search_query, ext: string) -> search_query
    + adds an extension filter; only files matching are returned
    # filter
  file_search.with_min_size
    fn (q: search_query, min_bytes: i64) -> search_query
    + adds a minimum file size filter in bytes
    # filter
  file_search.with_max_depth
    fn (q: search_query, depth: i32) -> search_query
    + limits recursion depth from the root
    ? depth 0 means only the root directory itself
    # filter
  file_search.find
    fn (q: search_query) -> result[list[string], string]
    + returns all paths matching all accumulated filters
    - returns error if the root cannot be read
    # execution
    -> std.fs.list_dir
    -> std.fs.stat
    -> std.path.join
    -> std.path.extension
