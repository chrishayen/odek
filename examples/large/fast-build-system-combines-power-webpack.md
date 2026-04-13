# Requirement: "a bundler that resolves, transforms, and packages source modules into a single output"

A dependency graph walker with pluggable loaders and a final emit step.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when path exists
      # filesystem
  std.path
    std.path.join
      @ (base: string, relative: string) -> string
      + joins two path segments
      # paths
    std.path.dirname
      @ (path: string) -> string
      + returns the directory portion of a path
      # paths
    std.path.extension
      @ (path: string) -> string
      + returns the extension including the dot, or empty
      # paths
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex sha256 of data
      # hashing

bundler
  bundler.new
    @ (entry: string, out_path: string) -> bundler_config
    + creates a bundler config with the given entry and output
    # construction
  bundler.register_loader
    @ (config: bundler_config, ext: string, loader: fn(string) -> result[loaded_module, string]) -> bundler_config
    + registers a loader for files with the given extension
    # loader_registration
  bundler.resolve_module
    @ (config: bundler_config, importer: string, request: string) -> result[string, string]
    + resolves a relative or bare import to an absolute path
    - returns error when no candidate exists
    # resolution
    -> std.path.join
    -> std.path.dirname
    -> std.fs.exists
  bundler.load_module
    @ (config: bundler_config, path: string) -> result[loaded_module, string]
    + reads the file and runs the loader matching its extension
    - returns error when no loader handles the extension
    # loading
    -> std.fs.read_all
    -> std.path.extension
  bundler.extract_imports
    @ (module: loaded_module) -> list[string]
    + returns the list of import specifiers found in the module source
    # parsing
  bundler.build_graph
    @ (config: bundler_config) -> result[module_graph, string]
    + walks imports starting from the entry, deduplicating modules by absolute path
    - returns error when any import fails to resolve
    # graph
  bundler.topological_order
    @ (graph: module_graph) -> result[list[string], string]
    + returns modules in dependency-first order
    - returns error when a cycle is detected
    # ordering
  bundler.emit_bundle
    @ (graph: module_graph, order: list[string]) -> string
    + concatenates modules into a single output with module wrappers
    # emit
  bundler.content_hash
    @ (output: string) -> string
    + returns a short content hash for cache busting
    # hashing
    -> std.hash.sha256_hex
  bundler.build
    @ (config: bundler_config) -> result[string, string]
    + runs the full pipeline and writes the bundle, returning the output path
    - returns error when any step fails
    # build
    -> std.fs.write_all
