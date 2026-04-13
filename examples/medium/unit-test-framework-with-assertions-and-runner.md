# Requirement: "a unit test framework with assertions and a runner"

Tests are registered as named closures; the runner executes each, captures pass/fail with a message, and returns a summary. Assertions are regular functions returning a result.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time

test_framework
  test_framework.new_suite
    @ (name: string) -> suite_state
    + returns an empty suite with the given name
    # construction
  test_framework.register
    @ (state: suite_state, name: string, body: test_fn) -> suite_state
    + adds a named test to the suite
    + a later registration with the same name overwrites the earlier one
    # registration
  test_framework.assert_eq
    @ (expected: value, actual: value) -> result[void, string]
    + returns ok when the two values are equal
    - returns an error message like "expected X, got Y"
    # assertion
  test_framework.assert_true
    @ (condition: bool, msg: string) -> result[void, string]
    + returns ok when condition is true
    - returns the provided message when false
    # assertion
  test_framework.run
    @ (state: suite_state) -> run_report
    + runs every registered test and collects (name, passed, message, duration)
    + records elapsed time per test
    # execution
    -> std.time.now_nanos
  test_framework.summary
    @ (report: run_report) -> string
    + formats a human-readable summary with counts of passed, failed, and skipped tests
    # reporting
