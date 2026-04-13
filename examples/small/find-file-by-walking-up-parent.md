# Requirement: "a library to find a file by walking up parent directories"

Walks from a starting directory upward until a matching filename is found or the root is reached.

std
  std.fs
    std.fs.exists
      @ (path: string) -> bool
      + returns true when a file or directory exists at path
      - returns false when nothing exists at path
      # filesystem
    std.fs.parent_dir
      @ (path: string) -> optional[string]
      + returns the parent directory path
      - returns none when path is the filesystem root
      # filesystem
  std.path
    std.path.join
      @ (base: string, name: string) -> string
      + joins base and name with the platform separator
      # path

find_up
  find_up.find
    @ (start_dir: string, filename: string) -> optional[string]
    + returns the full path to the nearest ancestor containing filename
    + searches start_dir first, then each parent in turn
    - returns none when no ancestor contains filename
    ? symlinks are followed like any other path
    # ancestor_lookup
    -> std.path.join
    -> std.fs.exists
    -> std.fs.parent_dir
