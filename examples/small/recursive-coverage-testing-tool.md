# Requirement: "a library for running the test suite across every nested package and merging the coverage reports"

Walks a project directory for nested test packages, runs each one with coverage enabled, and merges the per-package profiles into a combined report. Project layer owns the walker and merger; std provides filesystem and process primitives.

std
  std.fs
    std.fs.walk_dir
      @ (path: string) -> result[list[string], string]
      + returns all subdirectory paths under path, recursively
      - returns error when path does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file at path as a utf-8 string
      # filesystem
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, creating or truncating the file
      # filesystem
  std.process
    std.process.run
      @ (cmd: string, args: list[string], cwd: string) -> result[process_output, string]
      + runs a command synchronously in cwd and captures stdout/stderr/exit code
      # process

coverage_runner
  coverage_runner.find_test_packages
    @ (root: string) -> result[list[string], string]
    + returns every directory under root that contains test files
    # discovery
    -> std.fs.walk_dir
  coverage_runner.run_package
    @ (pkg_dir: string, profile_path: string) -> result[void, string]
    + runs the test command for pkg_dir and writes its coverage profile to profile_path
    - returns error when tests fail or the runner cannot be invoked
    # execution
    -> std.process.run
  coverage_runner.merge_profiles
    @ (profile_paths: list[string]) -> result[string, string]
    + parses each profile and returns a single merged profile string
    - returns error on malformed profiles
    # aggregation
    -> std.fs.read_all
  coverage_runner.run_all
    @ (root: string, out_path: string) -> result[void, string]
    + discovers, runs, and merges coverage for the whole tree, writing to out_path
    # orchestration
    -> std.fs.write_all
