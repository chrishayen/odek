# Requirement: "a disk space navigator that scans directories and visualizes usage"

Recursively sizes directories and exposes a zoomable tree so callers can render a treemap or list view.

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[dir_entry], string]
      + returns entries with name, is_dir flag, and size for files
      - returns error when path is not a directory
      # filesystem
    std.fs.stat
      @ (path: string) -> result[file_stat, string]
      + returns size, mtime, and type
      - returns error when path does not exist
      # filesystem

disk_navigator
  disk_navigator.scan
    @ (root: string) -> result[dir_node, string]
    + recursively walks root and returns a tree sized by total bytes
    - returns error when root cannot be read
    ? symlinks are not followed
    # scan
    -> std.fs.list_dir
    -> std.fs.stat
  disk_navigator.sort_children
    @ (node: dir_node, order: i8) -> dir_node
    + sorts children by size (0=desc, 1=asc) or name (2)
    # navigation
  disk_navigator.zoom_into
    @ (tree: dir_node, path: list[string]) -> result[dir_node, string]
    + returns the subtree at path
    - returns error when path does not exist in the tree
    # navigation
  disk_navigator.largest_files
    @ (tree: dir_node, n: i32) -> list[file_entry]
    + returns the n largest files anywhere in the tree
    # analysis
  disk_navigator.compute_percentages
    @ (node: dir_node) -> dir_node
    + annotates each child with its percentage of the parent's total size
    # analysis
