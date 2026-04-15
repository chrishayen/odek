# Requirement: "a library for writing shell scripts with easy command invocation and output capture"

Provides an ergonomic way to invoke shell commands, capture their output, and chain simple operations.

std
  std.process
    std.process.spawn
      fn (program: string, args: list[string]) -> result[process_handle, string]
      + launches a child process with the given arguments
      - returns error when the program cannot be found
      # process
    std.process.wait
      fn (handle: process_handle) -> result[process_result, string]
      + waits for the process to exit and returns exit code, stdout, and stderr
      - returns error if the process was killed by a signal
      # process
  std.fs
    std.fs.current_dir
      fn () -> string
      + returns the current working directory
      # filesystem
    std.fs.set_current_dir
      fn (path: string) -> result[void, string]
      + changes the current working directory
      - returns error when the path does not exist
      # filesystem

shellkit
  shellkit.run
    fn (command: string) -> result[process_result, string]
    + parses the command string into program and arguments and executes it
    + returns the captured stdout, stderr, and exit code
    - returns error when the command cannot be launched
    # execution
    -> std.process.spawn
    -> std.process.wait
  shellkit.run_checked
    fn (command: string) -> result[string, string]
    + runs a command and returns stdout on success
    - returns error containing stderr when the command exits non-zero
    # execution
    -> shellkit.run
  shellkit.pipe
    fn (commands: list[string]) -> result[string, string]
    + runs commands in sequence, feeding each stdout as the next stdin
    - returns error when any stage fails
    # pipelines
    -> shellkit.run
  shellkit.cd
    fn (path: string) -> result[void, string]
    + changes the working directory for subsequent commands
    - returns error when the path does not exist
    # environment
    -> std.fs.set_current_dir
  shellkit.pwd
    fn () -> string
    + returns the current working directory
    # environment
    -> std.fs.current_dir
