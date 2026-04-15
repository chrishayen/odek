# Requirement: "a library to find out which dependencies are slowing you down"

Reads a dependency tree, measures the on-disk size of each dependency subtree, and reports the largest offenders.

std
  std.fs
    std.fs.walk
      fn (root: string) -> list[string]
      + yields every path under root recursively
      # filesystem
    std.fs.stat_size
      fn (path: string) -> result[i64, string]
      + returns the size in bytes of a regular file
      - returns error when path does not exist
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the immediate children of a directory
      - returns error when path is not a directory
      # filesystem
  std.path
    std.path.join
      fn (base: string, name: string) -> string
      + joins base and name with the platform separator
      # path

dep_weigher
  dep_weigher.list_dependencies
    fn (deps_dir: string) -> result[list[string], string]
    + returns the names of each direct dependency directory under deps_dir
    - returns error when deps_dir does not exist
    # discovery
    -> std.fs.list_dir
  dep_weigher.measure_subtree
    fn (root: string) -> result[i64, string]
    + returns the total byte size of all files under root
    + includes nested subdirectories
    - returns error when root cannot be read
    # measurement
    -> std.fs.walk
    -> std.fs.stat_size
  dep_weigher.weigh_all
    fn (deps_dir: string) -> result[list[dep_weight], string]
    + returns one dep_weight per direct dependency with name and total bytes
    - returns error when deps_dir does not exist
    # aggregation
    -> std.path.join
  dep_weigher.sort_by_size
    fn (weights: list[dep_weight]) -> list[dep_weight]
    + returns weights sorted in descending order of bytes
    # ranking
  dep_weigher.top_offenders
    fn (weights: list[dep_weight], n: i32) -> list[dep_weight]
    + returns the n largest dependencies
    + returns all weights when n exceeds list length
    # selection
  dep_weigher.format_report
    fn (weights: list[dep_weight]) -> string
    + renders a human-readable report with name and size per line
    + sizes are shown with unit suffixes (KB, MB, GB)
    # rendering
