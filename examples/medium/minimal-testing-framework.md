# Requirement: "a minimal testing framework"

A small test registry with grouped describe/it blocks, assertion helpers, and a runner that reports results. Output formatting goes through thin std primitives.

std
  std.io
    std.io.print_line
      fn (line: string) -> void
      + writes the line followed by a newline to standard output
      # output
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

testing
  testing.new_suite
    fn (name: string) -> suite_state
    + creates an empty suite with the given top-level name
    # construction
  testing.describe
    fn (state: suite_state, group: string, body: fn(suite_state) -> suite_state) -> suite_state
    + registers a nested group; nested describes concatenate names with " > "
    # grouping
  testing.it
    fn (state: suite_state, name: string, body: fn() -> optional[string]) -> suite_state
    + registers a test case; body returns none on success or an error message on failure
    # registration
  testing.assert_equal
    fn (expected: string, actual: string) -> optional[string]
    + returns none when values are equal
    - returns an error message describing both sides when values differ
    # assertion
  testing.run
    fn (state: suite_state) -> run_report
    + executes all registered tests in registration order and collects pass/fail counts and durations
    + reports each test result as it runs
    # execution
    -> std.io.print_line
    -> std.time.now_millis
