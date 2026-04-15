# Requirement: "a library for selecting and executing infrastructure-as-code plan and apply operations across discovered modules"

Discovers modules in a directory tree, lets the caller pick some, and runs plan or apply on each by invoking an external tool.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all file paths under root
      - returns error when root is not a directory
      # filesystem
  std.process
    std.process.run
      fn (argv: list[string], cwd: string) -> result[process_output, string]
      + runs the command and returns stdout, stderr, and exit code
      - returns error when the executable cannot be launched
      # process

iac_runner
  iac_runner.discover_modules
    fn (root: string, marker_file: string) -> result[list[string], string]
    + returns directories under root that contain marker_file
    ? marker_file identifies a module (e.g. "main.tf")
    # discovery
    -> std.fs.walk
  iac_runner.plan
    fn (module_dir: string, tool_path: string) -> result[process_output, string]
    + runs the plan subcommand in module_dir
    - returns error when the tool exits non-zero and stderr is empty
    # execution
    -> std.process.run
  iac_runner.apply
    fn (module_dir: string, tool_path: string) -> result[process_output, string]
    + runs the apply subcommand in module_dir
    # execution
    -> std.process.run
  iac_runner.run_batch
    fn (module_dirs: list[string], action: string, tool_path: string) -> list[batch_result]
    + runs action on each module and collects per-module results
    ? action is "plan" or "apply"; unknown actions mark every module as skipped
    # batch
