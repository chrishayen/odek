# Requirement: "a library for defining project-specific named tasks and running them with dependencies"

std
  std.fs
    std.fs.read_text
      fn (path: string) -> result[string, string]
      + returns the contents of a text file
      - returns error when the file cannot be read
      # filesystem
  std.process
    std.process.run
      fn (command: string, args: list[string], env: map[string, string]) -> result[i32, string]
      + runs a command and returns its exit status
      - returns error when the command cannot be launched
      # process

task_runner
  task_runner.load
    fn (path: string) -> result[runner_state, string]
    + parses a task file at path into a runner state
    - returns error on syntax errors
    # loading
    -> std.fs.read_text
  task_runner.list_tasks
    fn (state: runner_state) -> list[string]
    + returns the names of every defined task
    # query
  task_runner.describe
    fn (state: runner_state, name: string) -> result[task_description, string]
    + returns the dependencies, parameters, and body of the named task
    - returns error when no task has that name
    # query
  task_runner.resolve_order
    fn (state: runner_state, target: string) -> result[list[string], string]
    + returns the tasks to run, in topological order, to satisfy target
    - returns error when a dependency cycle is detected
    - returns error when a referenced dependency does not exist
    # planning
  task_runner.run
    fn (state: runner_state, target: string, arguments: map[string, string]) -> result[i32, string]
    + runs target and its dependencies in order, returning the exit status of the last command
    - returns error when any command exits non-zero
    # execution
    -> std.process.run
