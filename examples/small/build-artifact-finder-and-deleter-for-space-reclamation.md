# Requirement: "a library for finding and deleting software project build artifacts to reclaim disk space"

Walk a directory, classify each child as a known project type, and report (or remove) reclaimable subdirectories.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the child entries of a directory
      - returns error when the path is not a directory
      # io
    std.fs.remove_tree
      fn (path: string) -> result[i64, string]
      + removes a directory recursively and returns the number of bytes freed
      - returns error when the path does not exist
      # io
    std.fs.size_tree
      fn (path: string) -> result[i64, string]
      + returns the total size in bytes of everything under path
      # io

reclaim
  reclaim.classify
    fn (path: string, children: list[string]) -> optional[project_kind]
    + returns a project kind when its marker files are present
    - returns none when no known marker is found
    # classification
  reclaim.artifact_dirs
    fn (kind: project_kind) -> list[string]
    + returns the reclaimable subdirectory names for a project kind
    # classification
  reclaim.scan
    fn (root: string) -> result[list[found_project], string]
    + walks root and returns every detected project with its artifact directories and total size
    # scanning
    -> std.fs.list_dir
    -> std.fs.size_tree
    -> reclaim.classify
    -> reclaim.artifact_dirs
  reclaim.clean
    fn (project: found_project) -> result[i64, string]
    + removes every artifact directory for the project and returns bytes freed
    # cleanup
    -> std.fs.remove_tree
