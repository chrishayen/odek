# Requirement: "a unit testing framework"

A collection of named test functions that run and produce a result report.

std: (all units exist)

testing
  testing.new_suite
    fn (name: string) -> suite_state
    + creates an empty suite with the given name
    # construction
  testing.add_test
    fn (suite: suite_state, name: string, thunk_id: string) -> suite_state
    + registers a named test referenced by thunk id
    # registration
  testing.assert_equal
    fn (expected: string, actual: string) -> result[void, string]
    + returns ok when the values are equal
    - returns a formatted error when they differ
    # assertion
  testing.assert_true
    fn (value: bool) -> result[void, string]
    + returns ok when value is true
    - returns error when value is false
    # assertion
  testing.run_one
    fn (suite: suite_state, name: string, outcome: result[void, string]) -> suite_state
    + records the outcome of one test by name
    # execution
  testing.run_all
    fn (suite: suite_state, outcomes: list[result[void, string]]) -> suite_state
    + records outcomes for every registered test in order
    # execution
  testing.report
    fn (suite: suite_state) -> test_report
    + returns a report with counts of passed, failed, and total tests
    + includes the failure message for each failed test
    # reporting
