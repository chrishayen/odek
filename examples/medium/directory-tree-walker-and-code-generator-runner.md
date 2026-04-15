# Requirement: "a code generation runner that walks a directory tree and invokes a generator on matching files"

A path walker, a regex filter, and an external-command runner. Directory source can be either a literal path or an environment variable name.

std
  std.fs
    std.fs.walk_dir
      fn (root: string) -> result[list[string], string]
      + yields every file path under root in depth-first order
      - returns error when root does not exist
      # filesystem
  std.env
    std.env.get
      fn (name: string) -> optional[string]
      + returns the value of an environment variable when set
      # environment
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_handle, string]
      + compiles a regular expression
      - returns error on invalid pattern
      # regex
    std.regex.matches
      fn (re: regex_handle, input: string) -> bool
      + reports whether the input matches
      # regex
  std.process
    std.process.run
      fn (cmd: string, args: list[string], workdir: string) -> result[i32, string]
      + runs a command in a working directory and returns its exit code
      - returns error when the binary cannot be spawned
      # process

generate_runner
  generate_runner.resolve_root
    fn (path_or_env: string) -> result[string, string]
    + returns the literal path when it exists on disk
    + returns the value of the environment variable when the input names one
    - returns error when neither a path nor an env var resolves
    # input_resolution
    -> std.env.get
  generate_runner.collect_targets
    fn (root: string, include_pattern: string) -> result[list[string], string]
    + walks the root and returns files whose path matches the pattern
    - returns error when the pattern is not a valid regex
    # discovery
    -> std.fs.walk_dir
    -> std.regex.compile
    -> std.regex.matches
  generate_runner.run_all
    fn (targets: list[string], command: string, command_args: list[string]) -> result[i32, string]
    + invokes the command in each target's directory and returns the count of successes
    - returns error on the first non-zero exit
    # execution
    -> std.process.run
