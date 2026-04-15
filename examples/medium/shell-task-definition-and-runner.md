# Requirement: "a library for defining and running shell-oriented tasks"

Users register named tasks (each wrapping a shell command or a callback) and invoke them by name with arguments. The library resolves the task, runs it, and returns the result.

std
  std.process
    std.process.run
      fn (command: string, args: list[string]) -> result[process_result, string]
      + runs the command with the given arguments and returns exit code, stdout, and stderr
      - returns error when the command cannot be launched
      # process
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits on every occurrence of the separator
      # strings

task_runner
  task_runner.new
    fn () -> task_registry
    + creates an empty task registry
    # construction
  task_runner.register_shell
    fn (reg: task_registry, name: string, command_template: string) -> task_registry
    + adds a shell task that substitutes positional arguments into the template before executing
    # registration
  task_runner.register_callback
    fn (reg: task_registry, name: string, handler: fn(list[string]) -> result[string, string]) -> task_registry
    + adds a task whose body is an arbitrary callback receiving the invocation arguments
    # registration
  task_runner.list_tasks
    fn (reg: task_registry) -> list[string]
    + returns the names of all registered tasks in registration order
    # inspection
  task_runner.invoke
    fn (reg: task_registry, name: string, args: list[string]) -> result[task_result, string]
    + runs the named task, substituting arguments into shell templates or passing them to callbacks
    - returns error when the task name is not registered
    - returns error when a shell task exits non-zero
    # execution
    -> std.process.run
    -> std.text.split
