# Requirement: "a mature full-featured testing framework"

A testing framework supports discovery, assertions, fixtures, parameterization, and structured reporting. The std layer provides file walking and pattern matching; the project layer orchestrates the run.

std
  std.fs
    std.fs.walk
      @ (root: string) -> list[string]
      + returns all file paths under root recursively
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full file contents
      - returns error when the file does not exist
      # filesystem
  std.text
    std.text.glob_match
      @ (pattern: string, name: string) -> bool
      + returns true when name matches a glob pattern
      # text
  std.time
    std.time.now_millis
      @ () -> i64
      + returns wall-clock time in milliseconds
      # time

testkit
  testkit.new_suite
    @ (name: string) -> suite_state
    + creates an empty suite with the given name
    # construction
  testkit.register_test
    @ (suite: suite_state, name: string, body: test_fn) -> suite_state
    + appends a test case to the suite
    # registration
  testkit.register_fixture
    @ (suite: suite_state, name: string, setup: fixture_fn, teardown: fixture_fn) -> suite_state
    + appends a fixture with setup and teardown hooks
    # fixtures
  testkit.parameterize
    @ (body: test_fn, params: list[map[string,string]]) -> list[test_fn]
    + expands a single test into one per parameter set
    ? each parameter map is bound into the body scope at invocation
    # parameterization
  testkit.discover
    @ (root: string, file_glob: string, name_glob: string) -> result[list[string], string]
    + returns paths of test files matching file_glob under root
    - returns error when root does not exist
    # discovery
    -> std.fs.walk
    -> std.text.glob_match
  testkit.run_suite
    @ (suite: suite_state) -> suite_report
    + runs all tests, applying matching fixtures
    + captures pass, fail, and error counts with durations
    - marks a test failed when its assertion function reports mismatch
    # execution
    -> std.time.now_millis
  testkit.assert_equal
    @ (got: string, want: string) -> result[void, string]
    + returns ok when got equals want
    - returns error describing the mismatch otherwise
    # assertions
  testkit.assert_contains
    @ (haystack: string, needle: string) -> result[void, string]
    + returns ok when haystack contains needle
    - returns error when it does not
    # assertions
  testkit.format_report
    @ (report: suite_report) -> string
    + renders a human-readable summary with per-test status
    # reporting
  testkit.filter_by_name
    @ (suite: suite_state, name_glob: string) -> suite_state
    + returns a suite containing only tests whose name matches the glob
    # selection
    -> std.text.glob_match
