# Requirement: "a lightweight test suite library with setup and teardown hooks"

A small harness that groups test cases and runs setup/teardown around each.

std: (all units exist)

suite
  suite.new
    @ (name: string) -> suite_state
    + creates an empty named suite
    # construction
  suite.set_before_each
    @ (state: suite_state, hook: callback) -> suite_state
    + registers a function to run before each test
    # hooks
  suite.set_after_each
    @ (state: suite_state, hook: callback) -> suite_state
    + registers a function to run after each test
    # hooks
  suite.add_test
    @ (state: suite_state, name: string, body: callback) -> suite_state
    + appends a test case identified by name
    # registration
  suite.run
    @ (state: suite_state) -> list[test_result]
    + runs before_each, body, after_each for each test in order
    + captures a pass or fail with message for each test
    - after_each still runs when body fails
    # execution
