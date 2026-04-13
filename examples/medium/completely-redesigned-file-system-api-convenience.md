# Requirement: "a convenient high-level filesystem api"

A thin layer over standard filesystem primitives that offers recursive listing, safe copy, and atomic writes.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the complete file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parents if needed
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + deletes a file or empty directory
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the immediate entries in a directory
      - returns error when path is not a directory
      # filesystem
    std.fs.stat
      @ (path: string) -> result[file_info, string]
      + returns metadata including size, mtime, and kind
      # filesystem

fs_convenience
  fs_convenience.read_text
    @ (path: string) -> result[string, string]
    + returns file contents decoded as utf-8 text
    - returns error on invalid utf-8
    # convenience
    -> std.fs.read_all
  fs_convenience.write_atomic
    @ (path: string, data: bytes) -> result[void, string]
    + writes to a sibling temp file then renames into place
    ? ensures readers never see a partially written file
    # durability
    -> std.fs.write_all
  fs_convenience.list_tree
    @ (root: string) -> result[list[string], string]
    + returns all file paths below root in depth-first order
    - returns error when root does not exist
    # traversal
    -> std.fs.list_dir
    -> std.fs.stat
  fs_convenience.copy
    @ (src: string, dst: string) -> result[void, string]
    + copies a single file, preserving bytes exactly
    - returns error when source is a directory
    # copy
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.fs.stat
  fs_convenience.remove_tree
    @ (root: string) -> result[void, string]
    + recursively removes a directory and all descendants
    - returns error when root does not exist
    # deletion
    -> std.fs.list_dir
    -> std.fs.remove
    -> std.fs.stat
