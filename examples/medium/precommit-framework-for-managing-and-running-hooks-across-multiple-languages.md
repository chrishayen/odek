# Requirement: "a framework for managing and running pre-commit hooks across multiple languages"

Users declare hooks in a config. For a given set of changed files the framework selects applicable hooks, runs them, and collects results.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file entirely into memory
      - returns error when the file cannot be read
      # filesystem
  std.strings
    std.strings.ends_with
      fn (s: string, suffix: string) -> bool
      + returns true when s ends with suffix
      # text
  std.process
    std.process.run
      fn (command: string, args: list[string], stdin: bytes) -> result[process_result, string]
      + runs a command to completion and returns exit code, stdout, stderr
      - returns error when the command cannot be spawned
      # process

precommit
  precommit.parse_config
    fn (raw: string) -> result[config, string]
    + parses a YAML-like hook configuration into a typed config
    - returns error when a hook entry is missing id or entry fields
    # config
  precommit.load_config
    fn (path: string) -> result[config, string]
    + reads and parses a config file
    # config
    -> std.fs.read_all
  precommit.select_hooks
    fn (cfg: config, changed_files: list[string]) -> list[hook]
    + returns hooks whose file pattern matches at least one changed file
    ? matching is by filename suffix; hooks with no pattern match every file
    # selection
    -> std.strings.ends_with
  precommit.run_hook
    fn (h: hook, files: list[string]) -> hook_result
    + runs the hook's command against the given files and returns pass/fail plus output
    - returns a failing result when the command exits non-zero
    # execution
    -> std.process.run
  precommit.run_all
    fn (cfg: config, changed_files: list[string]) -> list[hook_result]
    + selects and runs every applicable hook in config order
    # execution
  precommit.report
    fn (results: list[hook_result]) -> string
    + builds a human-readable summary of hook outcomes
    # reporting
