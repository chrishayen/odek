# Requirement: "a cross-platform shell command execution library"

Thin wrapper around a std process primitive. Captures stdout, stderr, and exit status.

std
  std.process
    std.process.spawn_and_wait
      fn (program: string, args: list[string], stdin: bytes) -> result[process_result, string]
      + runs the program to completion and returns stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process

shell
  shell.run
    fn (program: string, args: list[string]) -> result[command_output, string]
    + runs a command with no stdin and returns its output
    - returns error when the program cannot be launched
    # execution
    -> std.process.spawn_and_wait
  shell.run_with_input
    fn (program: string, args: list[string], input: bytes) -> result[command_output, string]
    + runs a command piping input to its stdin
    # execution
    -> std.process.spawn_and_wait
  shell.succeeded
    fn (output: command_output) -> bool
    + returns true when the exit code is zero
    # query
