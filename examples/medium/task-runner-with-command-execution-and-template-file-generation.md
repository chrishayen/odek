# Requirement: "a task runner that executes commands and generates files from templates"

Loads a task file describing commands and file templates, resolves variable references, then runs the selected task.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file as text
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes text, creating parent directories
      # filesystem
  std.process
    std.process.run
      @ (program: string, args: list[string], env: map[string, string], cwd: string) -> result[process_result, string]
      + runs the program synchronously and captures stdout, stderr, and exit code
      # process
  std.template
    std.template.render
      @ (source: string, vars: map[string, string]) -> result[string, string]
      + expands `{{ name }}` placeholders from the vars map
      - returns error when a referenced name is missing
      # templating

task_runner
  task_runner.load
    @ (path: string) -> result[task_file, string]
    + reads and parses a task description file
    - returns error when the file is missing or malformed
    # loading
    -> std.fs.read_all
  task_runner.list_tasks
    @ (tf: task_file) -> list[string]
    + returns the names of tasks in declaration order
    # inspection
  task_runner.run_task
    @ (tf: task_file, name: string, vars: map[string, string]) -> result[list[process_result], string]
    + executes each command in the named task, returning all results
    - returns error when the task does not exist
    - returns error on the first command with a non-zero exit code
    # execution
    -> std.template.render
    -> std.process.run
  task_runner.generate_file
    @ (tf: task_file, template_name: string, vars: map[string, string], out_path: string) -> result[void, string]
    + renders the named template and writes it to out_path
    - returns error when the template does not exist
    # generation
    -> std.template.render
    -> std.fs.write_all
