# Requirement: "a library that lets test suites run linters and fail on findings"

Exposes a small facade: register linters (each a function that returns findings), run them over a directory, and aggregate results into a pass/fail decision.

std
  std.fs
    std.fs.walk_files
      fn (root: string) -> result[list[string], string]
      + returns every file path under root recursively
      - returns error when root does not exist
      # filesystem

lintharness
  lintharness.new_harness
    fn () -> harness_state
    + creates a harness with no linters registered
    # construction
  lintharness.register
    fn (state: harness_state, name: string, linter: fn(list[string]) -> list[finding]) -> harness_state
    + registers a linter under a unique name
    # registration
  lintharness.run
    fn (state: harness_state, root: string) -> result[run_report, string]
    + walks the root and invokes each registered linter on the file list, collecting findings per linter
    - propagates filesystem errors
    # execution
    -> std.fs.walk_files
  lintharness.has_failures
    fn (report: run_report) -> bool
    + returns true when any linter reported at least one finding
    # verdict
