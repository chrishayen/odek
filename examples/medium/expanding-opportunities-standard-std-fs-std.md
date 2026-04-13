# Requirement: "a library of extended filesystem helpers on top of basic file I/O"

Adds recursive copy, recursive delete, and directory walking on top of the primitives that basic I/O provides.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of a regular file
      - returns error when the file does not exist or is not readable
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to a file, creating or truncating as needed
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[dir_entry], string]
      + returns the entries of a directory with their names and kinds
      - returns error when path is not a directory
      # filesystem
    std.fs.make_dir
      @ (path: string) -> result[void, string]
      + creates a directory, erroring if the parent does not exist
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + removes a regular file or empty directory
      # filesystem

fsx
  fsx.copy_file
    @ (src: string, dst: string) -> result[void, string]
    + copies a single file's bytes from src to dst
    # copy
    -> std.fs.read_all
    -> std.fs.write_all
  fsx.copy_tree
    @ (src: string, dst: string) -> result[i32, string]
    + recursively copies a directory tree, returning the number of files copied
    - returns error when src is not a directory
    # recursive_copy
    -> std.fs.list_dir
    -> std.fs.make_dir
  fsx.remove_tree
    @ (path: string) -> result[i32, string]
    + recursively removes a directory and all its contents, returning the number of entries removed
    # recursive_delete
    -> std.fs.list_dir
    -> std.fs.remove
  fsx.walk
    @ (root: string) -> result[list[string], string]
    + returns all file paths beneath root in depth-first order
    - returns error when root does not exist
    # walk
    -> std.fs.list_dir
  fsx.dir_size
    @ (root: string) -> result[i64, string]
    + returns the total byte size of all files beneath root
    # stats
    -> std.fs.read_all
