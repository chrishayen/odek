# Requirement: "a command-line shell library with tokenization, expansion, builtins, and a pipeline executor"

Parses input into commands, expands variables and globs, executes pipelines with external processes and a small set of builtins.

std
  std.fs
    std.fs.glob
      fn (pattern: string, cwd: string) -> result[list[string], string]
      + returns paths matching the glob pattern relative to cwd
      - returns error when the pattern is malformed
      # filesystem
    std.fs.set_cwd
      fn (path: string) -> result[void, string]
      + changes the current working directory of the process
      - returns error when the path does not exist
      # filesystem
    std.fs.get_cwd
      fn () -> string
      + returns the current working directory
      # filesystem
  std.process
    std.process.spawn_pipeline
      fn (stages: list[process_spec]) -> result[list[i32], string]
      + spawns each stage connected by pipes and returns exit codes in order
      - returns error when any stage fails to launch
      # process
  std.env
    std.env.get
      fn (name: string) -> optional[string]
      + returns the environment variable value or none
      # environment
    std.env.set
      fn (name: string, value: string) -> void
      + sets the environment variable
      # environment

shell
  shell.new
    fn () -> shell_state
    + creates a shell with empty history and a default variable scope
    # construction
  shell.tokenize
    fn (line: string) -> result[list[token], string]
    + splits the input into tokens honoring single and double quotes and backslash escapes
    - returns error on an unterminated quoted string
    # lexing
  shell.parse
    fn (tokens: list[token]) -> result[pipeline, string]
    + groups tokens into commands separated by pipe operators
    - returns error on an empty pipeline stage
    # parsing
  shell.expand
    fn (state: shell_state, cmd: command) -> result[command, string]
    + substitutes variable references and expands glob patterns in each argument
    - returns error when a referenced variable is unset in strict mode
    # expansion
    -> std.env.get
    -> std.fs.glob
  shell.run_builtin
    fn (state: shell_state, cmd: command) -> optional[result[shell_state, string]]
    + runs a builtin (cd, export, unset, exit) and returns the updated state
    + returns none when the command is not a builtin
    # builtins
    -> std.fs.set_cwd
    -> std.env.set
  shell.execute
    fn (state: shell_state, line: string) -> result[tuple[i32, shell_state], string]
    + tokenizes, parses, expands, and runs the pipeline; returns the final exit code
    - returns error when any stage fails to launch
    # execution
    -> std.process.spawn_pipeline
  shell.history_append
    fn (state: shell_state, line: string) -> shell_state
    + appends the line to the shell history buffer
    # history
  shell.history_list
    fn (state: shell_state) -> list[string]
    + returns history in insertion order
    # history
