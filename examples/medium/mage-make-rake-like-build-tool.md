# Requirement: "a make/rake-like build tool library"

Discovers named tasks, resolves their dependencies into a topological order, and runs them with caching based on file modification times.

std
  std.fs
    std.fs.mtime
      @ (path: string) -> result[i64, string]
      + returns the file's last modification time as a unix timestamp
      - returns error when the path does not exist
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.process
    std.process.run
      @ (command: string, args: list[string]) -> result[i32, string]
      + executes a command and returns its exit code
      - returns error when the command cannot be launched
      # process

build_tool
  build_tool.new
    @ () -> build_state
    + creates an empty task registry
    # construction
  build_tool.task
    @ (state: build_state, name: string, deps: list[string], action: string) -> build_state
    + registers a task with its dependencies and an action identifier
    # registration
  build_tool.target
    @ (state: build_state, name: string, inputs: list[string], outputs: list[string]) -> build_state
    + marks a task as producing files from inputs for staleness checks
    # registration
  build_tool.plan
    @ (state: build_state, goal: string) -> result[list[string], string]
    + returns the topologically ordered task list needed to build the goal
    - returns error when the goal is unknown
    - returns error when a dependency cycle is detected
    # planning
  build_tool.is_stale
    @ (state: build_state, task_name: string) -> bool
    + returns true when any input is newer than any output or an output is missing
    # caching
    -> std.fs.mtime
    -> std.fs.exists
  build_tool.run
    @ (state: build_state, goal: string) -> result[void, string]
    + executes the plan, skipping tasks whose outputs are up to date
    - returns error on the first task that fails
    # execution
    -> std.process.run
