# Requirement: "a test framework with async test support"

Tests are registered into suites, run with per-test timeouts, and reported with pass/fail counts.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.async
    std.async.run_with_timeout
      @ (task: fn() -> result[void, string], timeout_ms: i64) -> result[void, string]
      + runs task and returns error if it exceeds the timeout
      - returns the task's own error when it fails before the timeout
      # scheduling

testkit
  testkit.new_suite
    @ (name: string) -> test_suite
    + creates an empty suite with the given name
    # construction
  testkit.add_test
    @ (suite: test_suite, name: string, body: fn() -> result[void, string]) -> test_suite
    + registers a synchronous test case
    # registration
  testkit.add_async_test
    @ (suite: test_suite, name: string, body: fn() -> result[void, string], timeout_ms: i64) -> test_suite
    + registers an async test with a timeout
    # registration
  testkit.run_suite
    @ (suite: test_suite) -> suite_report
    + runs every test and returns per-case results with elapsed time
    # execution
    -> std.time.now_millis
    -> std.async.run_with_timeout
  testkit.format_report
    @ (report: suite_report) -> string
    + produces a human-readable summary with counts and failed case names
    # reporting
