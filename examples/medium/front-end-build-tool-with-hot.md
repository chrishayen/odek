# Requirement: "a front-end build tool with module bundling and hot module replacement"

Resolves a dependency graph from an entry module, bundles to an output, and emits an update patch when a source file changes.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.watch
      @ (path: string) -> result[string, string]
      + blocks until a file under the path changes and returns the changed path
      - returns error when the path cannot be watched
      # filesystem

bundler
  bundler.new_project
    @ (root: string) -> project_state
    + returns a project rooted at the given directory
    # construction
  bundler.resolve_imports
    @ (project: project_state, entry: string) -> result[list[string], string]
    + returns the transitive list of module paths reachable from the entry
    - returns error when an import cannot be resolved on disk
    # resolution
    -> std.fs.read_all
  bundler.bundle
    @ (project: project_state, entry: string) -> result[bytes, string]
    + returns the concatenated and wrapped module output for the entry graph
    - returns error on a resolution or read failure
    # bundling
    -> std.fs.read_all
  bundler.compute_update
    @ (project: project_state, changed_path: string) -> result[bytes, string]
    + returns a patch payload that swaps the changed module in a running client
    - returns error when the changed path is not part of the current graph
    # hmr
    -> std.fs.read_all
  bundler.watch_for_updates
    @ (project: project_state) -> result[string, string]
    + blocks until any source file in the project changes and returns its path
    - returns error when watching fails
    # hmr
    -> std.fs.watch
