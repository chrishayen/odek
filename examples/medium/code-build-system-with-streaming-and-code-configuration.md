# Requirement: "a streaming build system configured by code"

Tasks are declared in code with source globs, destination, and a transform. The runner wires them into a stream graph, watches sources, and re-runs dependents on change.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every regular file beneath root
      # filesystem
    std.fs.mtime
      fn (path: string) -> result[i64, string]
      + returns the modification time as a unix timestamp
      - returns error when the path is missing
      # filesystem

code_build
  code_build.new
    fn () -> build_graph
    + constructs an empty task graph
    # construction
  code_build.task
    fn (graph: build_graph, name: string, sources: list[string], dest: string, transform: stream_transform) -> build_graph
    + registers a named task in the graph
    # registration
  code_build.depends_on
    fn (graph: build_graph, task: string, upstream: string) -> result[build_graph, string]
    + records that task must run after upstream
    - returns error when either name is unknown
    # dependencies
  code_build.resolve_order
    fn (graph: build_graph) -> result[list[string], string]
    + returns tasks in topological order
    - returns error on cycles
    # ordering
  code_build.expand_sources
    fn (patterns: list[string], root: string) -> result[list[string], string]
    + returns matching files beneath root
    - returns error when a pattern is malformed
    # globbing
    -> std.fs.walk
  code_build.run_task
    fn (graph: build_graph, name: string, root: string) -> result[i32, string]
    + runs the task's transform over its expanded sources to its destination
    + returns the number of files written
    - returns error when the transform reports failure
    # task_run
  code_build.run_all
    fn (graph: build_graph, root: string) -> result[i32, string]
    + runs every task in resolved order
    + returns the total files written
    - returns error when any task fails
    # run_all
  code_build.changed_since
    fn (paths: list[string], reference_time: i64) -> result[list[string], string]
    + returns paths whose mtime is newer than reference_time
    # watcher
    -> std.fs.mtime
  code_build.run_incremental
    fn (graph: build_graph, root: string, last_run: i64) -> result[i32, string]
    + runs only tasks whose sources changed since last_run
    + returns the number of tasks re-run
    # incremental
