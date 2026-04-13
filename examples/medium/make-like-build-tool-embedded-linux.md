# Requirement: "a make-like build tool for embedded systems"

Parses recipes that declare source packages with dependencies and build/install phases, then executes them in dependency order in isolated working directories.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's entire contents
      - returns error when the path cannot be read
      # filesystem
    std.fs.make_dir_all
      @ (path: string) -> result[void, string]
      + creates a directory and any missing parents
      # filesystem
  std.process
    std.process.run_in
      @ (cwd: string, command: string, args: list[string], env: map[string, string]) -> result[i32, string]
      + runs a command in a given working directory with a custom environment
      - returns error when the command cannot be launched
      # process

build_tool
  build_tool.new
    @ (work_root: string) -> build_state
    + creates a build state rooted at a working directory
    # construction
  build_tool.load_recipe
    @ (state: build_state, path: string) -> result[build_state, string]
    + parses a recipe file with name, version, depends, sources, do_build, do_install fields
    - returns error on malformed recipe syntax
    # recipes
    -> std.fs.read_all
  build_tool.resolve_order
    @ (state: build_state, target: string) -> result[list[string], string]
    + returns recipes in topologically sorted order for building the target
    - returns error when a dependency cycle is detected
    - returns error when a named dependency has no recipe
    # planning
  build_tool.build_one
    @ (state: build_state, recipe_name: string) -> result[void, string]
    + executes the configure, build, and install phases of a single recipe in its isolated workdir
    - returns error when any phase exits non-zero
    # execution
    -> std.fs.make_dir_all
    -> std.process.run_in
  build_tool.build
    @ (state: build_state, target: string) -> result[void, string]
    + builds the target and all its dependencies in order
    - returns error on the first recipe that fails
    # execution
