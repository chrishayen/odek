# Requirement: "a command launcher that matches input against executables on PATH"

Indexes executables from PATH directories and produces ranked completions for a query.

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns file names directly under path
      - returns error when the directory cannot be read
      # filesystem

launcher
  launcher.build_index
    @ (path_dirs: list[string]) -> launcher_index
    + collects executable names from each directory into a deduplicated index
    ? duplicate names across directories keep the first occurrence
    # indexing
    -> std.fs.list_dir
  launcher.match
    @ (index: launcher_index, query: string) -> list[string]
    + returns names starting with query first, then names containing query
    + returns an empty list when query matches nothing
    # matching
