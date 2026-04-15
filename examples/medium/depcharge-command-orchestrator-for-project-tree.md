# Requirement: "a library that orchestrates the execution of commands across many dependencies in a project tree"

Dependencies are declared as a set of directories; commands run sequentially or in parallel across them.

std
  std.process
    std.process.run
      fn (workdir: string, program: string, args: list[string]) -> result[process_result, string]
      + runs a command in workdir and returns (exit_code, stdout, stderr)
      # process
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      # filesystem

depcharge
  depcharge.load_plan
    fn (text: string) -> result[dep_plan, string]
    + parses a plan description (one dependency directory per line) into a plan value
    - lines starting with '#' are treated as comments
    - returns error when no directories remain
    # plan
  depcharge.load_plan_from_file
    fn (path: string) -> result[dep_plan, string]
    + reads a plan from disk and parses it
    # plan
    -> std.fs.read_all
    -> depcharge.load_plan
  depcharge.run_command
    fn (plan: dep_plan, program: string, args: list[string]) -> list[command_result]
    + runs program with args once per dependency, in declaration order
    + each command_result records the dependency directory, exit code, and captured output
    # execution
    -> std.process.run
  depcharge.run_command_parallel
    fn (plan: dep_plan, program: string, args: list[string], workers: i32) -> list[command_result]
    + runs the command across dependencies with up to workers concurrent processes
    ? result order matches plan order regardless of completion order
    # execution
    -> std.process.run
  depcharge.filter_plan
    fn (plan: dep_plan, predicate: fn(dir: string) -> bool) -> dep_plan
    + returns a new plan containing only dependencies for which predicate returns true
    # plan
  depcharge.summarize_results
    fn (results: list[command_result]) -> command_summary
    + counts successes and failures and captures the first failing dependency
    # reporting
