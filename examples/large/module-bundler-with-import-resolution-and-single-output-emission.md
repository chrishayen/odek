# Requirement: "a module bundler that resolves imports and emits a single output"

Parse module sources for imports, resolve them, build a dependency graph, and emit in topological order.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full contents of a file as a string
      - returns error when the file cannot be read
      # filesystem
  std.path
    std.path.join
      fn (base: string, rel: string) -> string
      + joins a base directory with a relative path and normalizes
      # paths
    std.path.dirname
      fn (path: string) -> string
      + returns the parent directory portion of a path
      # paths

bundler
  bundler.new
    fn (entry: string) -> bundler_state
    + creates a bundler with the given entry module path
    # construction
  bundler.parse_imports
    fn (source: string) -> list[string]
    + returns all import specifiers found in the source
    - returns empty list when the source has no imports
    # parsing
  bundler.resolve
    fn (from: string, specifier: string) -> result[string, string]
    + returns the absolute module path for a relative specifier
    - returns error when specifier cannot be resolved
    # resolution
    -> std.path.dirname
    -> std.path.join
  bundler.load_module
    fn (state: bundler_state, path: string) -> result[bundler_state, string]
    + reads the module source and records its imports
    - returns error when read fails
    - skips modules already loaded
    # graph_building
    -> std.fs.read_all
  bundler.build_graph
    fn (state: bundler_state) -> result[bundler_state, string]
    + recursively loads the entry and all transitively imported modules
    - returns error on any resolution or read failure
    # graph_building
  bundler.topological_order
    fn (state: bundler_state) -> result[list[string], string]
    + returns module paths in dependency order
    - returns error when a cycle is detected
    # ordering
  bundler.emit
    fn (state: bundler_state) -> result[string, string]
    + concatenates modules in topological order into a single output
    - returns error when the graph has not been built
    # emission
  bundler.tree_shake
    fn (state: bundler_state, roots: list[string]) -> bundler_state
    + removes modules unreachable from the given export roots
    # optimization
  bundler.module_count
    fn (state: bundler_state) -> i32
    + returns the number of modules loaded so far
    # inspection
