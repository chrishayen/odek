# Requirement: "bundles source modules and assets into browser-loadable artifacts"

A module bundler: resolve a dependency graph from an entry, parse each module's imports, topologically order them, and emit a single bundle. Asset files are hashed and copied through. std supplies file IO, hashing, and path manipulation.

std
  std.fs
    std.fs.read_text
      @ (path: string) -> result[string, string]
      + returns the file contents as UTF-8 text
      - returns error when the file cannot be read
      # filesystem
    std.fs.read_bytes
      @ (path: string) -> result[bytes, string]
      + returns the raw file bytes
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_text
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, creating parent directories as needed
      - returns error when the path cannot be written
      # filesystem
    std.fs.write_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + writes raw bytes to path
      - returns error when the path cannot be written
      # filesystem
  std.path
    std.path.join
      @ (base: string, rel: string) -> string
      + returns the normalized path of rel resolved against base
      # path
    std.path.dir
      @ (path: string) -> string
      + returns the parent directory of path
      # path
    std.path.extension
      @ (path: string) -> string
      + returns the extension including the dot, or "" when absent
      # path
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex-encoded SHA-256 digest
      # cryptography

bundler
  bundler.new_config
    @ (entry: string, output_dir: string) -> bundler_config
    + returns a config with the given entry and output directory
    # construction
  bundler.resolve_import
    @ (importer_path: string, spec: string) -> result[string, string]
    + returns the absolute file path the spec resolves to
    - returns error when the spec cannot be resolved on disk
    # resolution
    -> std.path.join
    -> std.path.dir
  bundler.parse_imports
    @ (source: string) -> list[string]
    + returns every import specifier referenced by the source in order
    + returns an empty list when the source has no imports
    # parsing
  bundler.load_module
    @ (path: string) -> result[bundler_module, string]
    + returns a module whose source, path, and import list are populated
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_text
  bundler.build_graph
    @ (cfg: bundler_config) -> result[bundler_graph, string]
    + returns a graph containing every module reachable from the entry
    - returns error when a transitive import cannot be resolved
    # graph_construction
  bundler.topo_order
    @ (graph: bundler_graph) -> result[list[string], string]
    + returns module paths in dependency-first order
    - returns error when a cycle is detected
    # ordering
  bundler.emit_bundle
    @ (graph: bundler_graph, order: list[string]) -> string
    + returns the concatenated bundle wrapping each module with an id
    # code_generation
  bundler.hash_asset
    @ (path: string) -> result[string, string]
    + returns a content-addressed filename like "logo.<hash8>.png"
    - returns error when the asset cannot be read
    # asset_pipeline
    -> std.fs.read_bytes
    -> std.hash.sha256_hex
    -> std.path.extension
  bundler.copy_asset
    @ (src_path: string, output_dir: string) -> result[string, string]
    + returns the hashed output path after copying the asset
    - returns error when the source cannot be read or destination cannot be written
    # asset_pipeline
    -> std.fs.read_bytes
    -> std.fs.write_bytes
    -> std.path.join
  bundler.build
    @ (cfg: bundler_config) -> result[bundler_output, string]
    + returns the list of emitted files and the bundle path on success
    - returns error when graph construction, ordering, or emission fails
    # orchestration
    -> std.fs.write_text
    -> std.path.join
