# Requirement: "a project command runner driven by a named-command config"

Loads a config mapping task names to shell commands and runs them by name.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of a file as text
      - returns error when the file cannot be read
      # filesystem
  std.process
    std.process.run_shell
      @ (command: string) -> result[i32, string]
      + runs a shell command and returns its exit status
      - returns error when the shell cannot be spawned
      # process

taskrunner
  taskrunner.parse_config
    @ (text: string) -> result[map[string,string], string]
    + parses "name: command" lines into a name-to-command map
    - returns error on a line missing the separator
    # parsing
  taskrunner.load_config
    @ (path: string) -> result[map[string,string], string]
    + reads and parses a config file
    # loading
    -> std.fs.read_all
    -> taskrunner.parse_config
  taskrunner.run_task
    @ (config: map[string,string], name: string) -> result[i32, string]
    + runs the command bound to name and returns its exit status
    - returns error when name is not in the config
    # execution
    -> std.process.run_shell
