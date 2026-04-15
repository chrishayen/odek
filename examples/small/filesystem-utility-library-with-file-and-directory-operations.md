# Requirement: "a filesystem utility library providing higher-level file and directory operations"

Recursive copy, recursive remove, ensure-directory, move, and an atomic write. The project layer composes std.fs primitives.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at path
      - returns error when the file does not exist or is unreadable
      # fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, replacing any existing content
      - returns error when the path is not writable
      # fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the names of entries in the directory
      - returns error when path is not a directory
      # fs
    std.fs.remove
      fn (path: string) -> result[void, string]
      + removes a file or empty directory; idempotent when the path is absent
      # fs
    std.fs.mkdir
      fn (path: string) -> result[void, string]
      + creates the directory; idempotent when it already exists
      # fs
    std.fs.rename
      fn (from: string, to: string) -> result[void, string]
      + renames a path
      - returns error when from does not exist
      # fs
    std.fs.is_dir
      fn (path: string) -> bool
      + returns true when the path exists and is a directory
      # fs

fsx
  fsx.ensure_dir
    fn (path: string) -> result[void, string]
    + creates the directory and all missing parents
    - returns error when a path component exists and is not a directory
    # directories
    -> std.fs.mkdir
    -> std.fs.is_dir
  fsx.copy_tree
    fn (src: string, dst: string) -> result[void, string]
    + recursively copies files and directories from src to dst, creating dst as needed
    - returns error when src does not exist
    # copy
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.fs.list_dir
    -> std.fs.is_dir
  fsx.remove_tree
    fn (path: string) -> result[void, string]
    + recursively removes a directory and everything inside; idempotent when the path is absent
    # remove
    -> std.fs.list_dir
    -> std.fs.remove
    -> std.fs.is_dir
  fsx.move
    fn (src: string, dst: string) -> result[void, string]
    + renames src to dst, falling back to copy-and-remove when rename fails
    - returns error when src does not exist
    # move
    -> std.fs.rename
  fsx.write_atomic
    fn (path: string, data: bytes) -> result[void, string]
    + writes data to a temporary sibling file and renames it over path
    - returns error on write or rename failure, leaving the original path unchanged
    # atomic
    -> std.fs.write_all
    -> std.fs.rename
