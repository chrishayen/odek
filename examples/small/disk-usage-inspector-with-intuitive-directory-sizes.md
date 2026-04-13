# Requirement: "a disk-usage inspector that reports directory sizes more intuitively than du"

Walks a directory tree, aggregates sizes, and returns a ranked view suitable for rendering.

std
  std.fs
    std.fs.walk
      @ (path: string) -> list[file_entry]
      + returns all files under path with their size in bytes
      + follows directories but not symbolic links
      # filesystem

dust
  dust.scan
    @ (root: string) -> map[string, u64]
    + returns total bytes per directory, with parents accumulating children
    - returns an empty map when root does not exist
    # scanning
    -> std.fs.walk
  dust.top_n
    @ (sizes: map[string, u64], n: i32) -> list[tuple[string, u64]]
    + returns the n largest entries in descending order
    # ranking
  dust.format_size
    @ (bytes: u64) -> string
    + returns a human-readable size using binary units (KiB, MiB, GiB)
    # formatting
  dust.render_tree
    @ (sizes: map[string, u64], root: string, max_depth: i32) -> string
    + returns an indented tree view with size bars proportional to root
    # rendering
    -> dust.format_size
